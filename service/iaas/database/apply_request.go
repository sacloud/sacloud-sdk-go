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

package database

import (
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	builder2 "github.com/sacloud/iaas-service-go/database/builder"
	"github.com/sacloud/packages-go/validate"
)

type ApplyRequest struct {
	Zone string `service:"-" validate:"required"`

	ID                    types.ID `service:"-"`
	Name                  string   `validate:"required"`
	Description           string   `validate:"min=0,max=512"`
	Tags                  types.Tags
	IconID                types.ID
	PlanID                types.ID `validate:"required"`
	SwitchID              types.ID `validate:"required"`
	IPAddresses           []string `validate:"required,min=1,max=2,dive,ipv4"`
	NetworkMaskLen        int      `validate:"required,min=1,max=32"`
	DefaultRoute          string   `validate:"omitempty,ipv4"`
	Port                  int      `validate:"omitempty,min=1,max=65535"`
	SourceNetwork         []string `validate:"dive,cidrv4"`
	DatabaseType          string   `validate:"required,oneof=mariadb postgres"`
	DatabaseVersion       string
	Username              string `validate:"required"`
	Password              string `validate:"required"`
	EnableReplication     bool
	ReplicaUserPassword   string `validate:"required_with=EnableReplication"`
	EnableWebUI           bool
	EnableBackup          bool
	BackupWeekdays        []types.EDayOfTheWeek `validate:"required_with=EnableBackup,max=7"`
	BackupStartTimeHour   int                   `validate:"omitempty,min=0,max=23"`
	BackupStartTimeMinute int                   `validate:"omitempty,oneof=0 15 30 45"`
	Parameters            map[string]interface{}

	EnableBackupv2          bool
	Backupv2Weekdays        []types.EDayOfTheWeek `validate:"required_with=EnableBackupv2,max=7"`
	Backupv2StartTimeHour   int                   `validate:"omitempty,min=0,max=23"`
	Backupv2StartTimeMinute int                   `validate:"omitempty,oneof=0 15 30 45"`
	Backupv2Connect         string                `validate:"required_with=EnableBackupv2"`

	EnableMonitoringSuite bool

	DiskEncryptionAlgorithm types.EDiskEncryptionAlgorithm
	DiskEncryptionKMSKey    types.ID

	NoWait bool
}

func (req *ApplyRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *ApplyRequest) Builder(caller iaas.APICaller) (*builder2.Builder, error) {
	replicaUser := ""
	replicaPassword := ""
	if req.EnableReplication {
		replicaUser = "replica"
		replicaPassword = req.ReplicaUserPassword
	}
	builder := &builder2.Builder{
		ID:   req.ID,
		Zone: req.Zone,

		PlanID:         req.PlanID,
		SwitchID:       req.SwitchID,
		IPAddresses:    req.IPAddresses,
		NetworkMaskLen: req.NetworkMaskLen,
		DefaultRoute:   req.DefaultRoute,
		Conf: &iaas.DatabaseRemarkDBConfCommon{
			DatabaseName:    types.RDBMSTypeFromString(req.DatabaseType).String(),
			DatabaseVersion: req.DatabaseVersion,
		},
		CommonSetting: &iaas.DatabaseSettingCommon{
			WebUI:           types.ToWebUI(req.EnableWebUI),
			ServicePort:     req.Port,
			SourceNetwork:   req.SourceNetwork,
			DefaultUser:     req.Username,
			UserPassword:    req.Password,
			ReplicaUser:     replicaUser,
			ReplicaPassword: replicaPassword,
		},
		MonitoringSuite: &iaas.MonitoringSuite{
			Enabled: req.EnableMonitoringSuite,
		},
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
		IconID:      req.IconID,
		Parameters:  req.Parameters,
		NoWait:      req.NoWait,
		Client:      builder2.NewAPIClient(caller),
	}
	if req.EnableBackup {
		builder.BackupSetting = &iaas.DatabaseSettingBackup{
			Time:      fmt.Sprintf("%02d:%02d", req.BackupStartTimeHour, req.BackupStartTimeMinute),
			DayOfWeek: req.BackupWeekdays,
			Rotate:    8,
		}
	}
	if req.EnableBackupv2 {
		builder.Backupv2Setting = &iaas.DatabaseSettingBackupv2{
			Time:      fmt.Sprintf("%02d:%02d", req.Backupv2StartTimeHour, req.Backupv2StartTimeMinute),
			DayOfWeek: req.Backupv2Weekdays,
			Connect:   req.Backupv2Connect,
			Rotate:    8,
		}
	}
	if req.EnableReplication {
		builder.ReplicationSetting = &iaas.DatabaseReplicationSetting{
			Model: types.DatabaseReplicationModels.MasterSlave,
		}
	}
	if req.DiskEncryptionAlgorithm == types.DiskEncryptionAlgorithms.AES256XTS && !req.DiskEncryptionKMSKey.IsEmpty() {
		builder.Disk = &iaas.DatabaseDisk{
			EncryptionAlgorithm: req.DiskEncryptionAlgorithm,
			EncryptionKeyID:     req.DiskEncryptionKMSKey,
		}
	}
	return builder, nil
}
