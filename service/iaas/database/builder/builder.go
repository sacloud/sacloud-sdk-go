// Copyright 2022-2025 The sacloud/iaas-service-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package builder

import (
	"context"
	"errors"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/accessor"
	"github.com/sacloud/iaas-api-go/defaults"
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/types"
	setup2 "github.com/sacloud/iaas-service-go/setup"
)

// Builder データベースの構築を行う
type Builder struct {
	ID   types.ID
	Zone string

	PlanID             types.ID
	SwitchID           types.ID
	IPAddresses        []string
	NetworkMaskLen     int
	DefaultRoute       string
	Conf               *iaas.DatabaseRemarkDBConfCommon
	SourceID           types.ID
	CommonSetting      *iaas.DatabaseSettingCommon
	BackupSetting      *iaas.DatabaseSettingBackup
	Backupv2Setting    *iaas.DatabaseSettingBackupv2
	ReplicationSetting *iaas.DatabaseReplicationSetting
	MonitoringSuite    *iaas.MonitoringSuite
	Disk               *iaas.DatabaseDisk
	Name               string
	Description        string
	Tags               types.Tags
	IconID             types.ID

	// Parameters RDBMS固有のパラメータ設定
	//
	// キーにはiaas.DatabaseParameterMetaのLabelを指定する
	//   - 例: effective_cache_size: 10
	Parameters map[string]interface{}

	SettingsHash string

	NoWait bool

	SetupOptions *setup2.Options
	Client       *APIClient
}

func (b *Builder) init() {
	if b.SetupOptions == nil {
		b.SetupOptions = &setup2.Options{}
	}
	b.SetupOptions.Init()
}

// Validate 設定値の検証
func (b *Builder) Validate(ctx context.Context, zone string) error {
	requiredValues := map[string]bool{
		"PlanID":         b.PlanID.IsEmpty(),
		"SwitchID":       b.SwitchID.IsEmpty(),
		"IPAddresses":    len(b.IPAddresses) == 0,
		"NetworkMaskLen": b.NetworkMaskLen == 0,
		"Conf":           b.Conf == nil,
		"CommonSetting":  b.CommonSetting == nil,
	}
	for key, empty := range requiredValues {
		if empty {
			return fmt.Errorf("%s is required", key)
		}
	}
	return nil
}

// Build データベースアプライアンスの構築
func (b *Builder) Build(ctx context.Context) (*iaas.Database, error) {
	if b.ID.IsEmpty() {
		return b.create(ctx, b.Zone)
	}
	return b.update(ctx, b.Zone, b.ID)
}

func (b *Builder) create(ctx context.Context, zone string) (*iaas.Database, error) {
	b.init()

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	builder := &setup2.RetryableSetup{
		Create: func(ctx context.Context, zone string) (accessor.ID, error) {
			return b.Client.Database.Create(ctx, zone, &iaas.DatabaseCreateRequest{
				PlanID:             b.PlanID,
				SwitchID:           b.SwitchID,
				IPAddresses:        b.IPAddresses,
				NetworkMaskLen:     b.NetworkMaskLen,
				DefaultRoute:       b.DefaultRoute,
				Conf:               b.Conf,
				SourceID:           b.SourceID,
				CommonSetting:      b.CommonSetting,
				BackupSetting:      b.BackupSetting,
				Backupv2Setting:    b.Backupv2Setting,
				ReplicationSetting: b.ReplicationSetting,
				MonitoringSuite:    b.MonitoringSuite,
				Disk:               b.Disk,
				Name:               b.Name,
				Description:        b.Description,
				Tags:               b.Tags,
				IconID:             b.IconID,
			})
		},
		Delete: func(ctx context.Context, zone string, id types.ID) error {
			return b.Client.Database.Delete(ctx, zone, id)
		},
		Read: func(ctx context.Context, zone string, id types.ID) (interface{}, error) {
			return b.Client.Database.Read(ctx, zone, id)
		},
		ProvisionBeforeUp: func(ctx context.Context, zone string, id types.ID, _ interface{}) error {
			if b.NoWait {
				return nil
			}

			// [HACK] データベースアプライアンス場合のみ/appliance/:id/statusも考慮する
			waiter := iaas.WaiterForUp(func() (interface{}, error) {
				return b.Client.Database.Status(ctx, zone, id)
			})
			waiter.(*iaas.StatePollingWaiter).Interval = defaults.DefaultDBStatusPollingInterval // HACK 現状は決め打ち、ユースケースが出たら修正する

			_, err := waiter.WaitForState(ctx)
			if err != nil {
				return err
			}

			if err := b.reconcileDatabaseParameters(ctx, zone, id); err != nil {
				return err
			}

			return b.Client.Database.Config(ctx, zone, id)
		},
		IsWaitForCopy: !b.NoWait,
		IsWaitForUp:   !b.NoWait,
		Options:       b.SetupOptions,
	}

	result, err := builder.Setup(ctx, zone)
	var db *iaas.Database
	if result != nil {
		db = result.(*iaas.Database)
	}
	if err != nil {
		return db, err
	}

	// refresh
	db, err = b.Client.Database.Read(ctx, zone, db.ID)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Update データベースの更新
func (b *Builder) update(ctx context.Context, zone string, id types.ID) (*iaas.Database, error) {
	b.init()

	if err := b.Validate(ctx, zone); err != nil {
		return nil, err
	}

	// check Database is exists
	db, err := b.Client.Database.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	isNeedShutdown, err := b.collectUpdateInfo(db)
	if err != nil {
		return nil, err
	}

	isNeedRestart := false
	if db.InstanceStatus.IsUp() && isNeedShutdown {
		if b.NoWait {
			return nil, errors.New("NoWait option is not available due to the need to shut down")
		}

		isNeedRestart = true
		if err := power.ShutdownDatabase(ctx, b.Client.Database, zone, id, false); err != nil {
			return nil, err
		}
	}

	_, err = b.Client.Database.Update(ctx, zone, id, &iaas.DatabaseUpdateRequest{
		Name:               b.Name,
		Description:        b.Description,
		Tags:               b.Tags,
		IconID:             b.IconID,
		CommonSetting:      b.CommonSetting,
		BackupSetting:      b.BackupSetting,
		Backupv2Setting:    b.Backupv2Setting,
		ReplicationSetting: b.ReplicationSetting,
		MonitoringSuite:    b.MonitoringSuite,
		SettingsHash:       b.SettingsHash,
	})
	if err != nil {
		return nil, err
	}
	if err := b.reconcileDatabaseParameters(ctx, zone, id); err != nil {
		return nil, err
	}
	if err := b.Client.Database.Config(ctx, zone, id); err != nil {
		return nil, err
	}
	if isNeedRestart {
		if err := power.BootDatabase(ctx, b.Client.Database, zone, id); err != nil {
			return nil, err
		}
	}

	// refresh
	db, err = b.Client.Database.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}
	return db, err
}

func (b *Builder) collectUpdateInfo(db *iaas.Database) (isNeedShutdown bool, err error) {
	isNeedShutdown = b.CommonSetting.ReplicaPassword != db.CommonSetting.ReplicaPassword
	return
}

func (b *Builder) reconcileDatabaseParameters(ctx context.Context, zone string, id types.ID) error {
	parameters, err := b.Client.Database.GetParameter(ctx, zone, id)
	if err != nil {
		return err
	}

	newParameters := make(map[string]interface{})
	// 既存のパラメータは一旦nullに
	for k := range parameters.Settings {
		newParameters[k] = nil
	}

	for k, v := range b.Parameters {
		found := false
		for _, meta := range parameters.MetaInfo {
			if k == meta.Label {
				newParameters[meta.Name] = v
				found = true
				break
			}
		}
		// kvのキーがラベルではなかったらそのままnmへ
		if !found {
			newParameters[k] = v
		}
	}
	if len(newParameters) > 0 {
		// DatabaseAPI.Configはあとで呼ぶ
		return b.Client.Database.SetParameter(ctx, zone, id, newParameters)
	}

	return nil
}
