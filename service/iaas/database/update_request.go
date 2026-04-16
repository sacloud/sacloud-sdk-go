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
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	Zone string   `validate:"required"`
	ID   types.ID `validate:"required"`

	Name        *string     `service:",omitempty" validate:"omitempty,min=1"`
	Description *string     `service:",omitempty" validate:"omitempty,max=512"`
	Tags        *types.Tags `service:",omitempty"`
	IconID      *types.ID   `service:",omitempty"`

	SourceNetwork         *[]string               `service:",omitempty" validate:"omitempty,dive,cidrv4"`
	EnableReplication     *bool                   `service:",omitempty"`
	ReplicaUserPassword   *string                 `service:",omitempty" validate:"omitempty,required_with=EnableReplication"`
	EnableWebUI           *bool                   `service:",omitempty"`
	EnableBackup          *bool                   `service:",omitempty"`
	BackupWeekdays        *[]types.EDayOfTheWeek  `service:",omitempty" validate:"omitempty,required_with=EnableBackup,max=7"`
	BackupStartTimeHour   *int                    `service:",omitempty" validate:"omitempty,min=0,max=23"`
	BackupStartTimeMinute *int                    `service:",omitempty" validate:"omitempty,oneof=0 15 30 45"`
	Parameters            *map[string]interface{} `service:",omitempty"`

	EnableBackupv2          *bool
	Backupv2Weekdays        *[]types.EDayOfTheWeek `validate:"required_with=EnableBackupv2,max=7"`
	Backupv2StartTimeHour   *int                   `validate:"omitempty,min=0,max=23"`
	Backupv2StartTimeMinute *int                   `validate:"omitempty,oneof=0 15 30 45"`
	Backupv2Connect         *string                `validate:"required_with=EnableBackupv2"`

	EnableMonitoringSuite *bool

	SettingsHash string
	NoWait       bool
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) ApplyRequest(ctx context.Context, caller iaas.APICaller) (*ApplyRequest, error) {
	dbOp := iaas.NewDatabaseOp(caller)
	current, err := dbOp.Read(ctx, req.Zone, req.ID)
	if err != nil {
		return nil, err
	}

	if current.Availability != types.Availabilities.Available {
		return nil, fmt.Errorf("target has invalid Availability: Zone=%s ID=%s Availability=%v", req.Zone, req.ID.String(), current.Availability)
	}

	var bkHour, bkMinute int
	var bkWeekdays []types.EDayOfTheWeek
	if current.BackupSetting != nil {
		bkWeekdays = current.BackupSetting.DayOfWeek
		if current.BackupSetting.Time != "" {
			timeStrings := strings.Split(current.BackupSetting.Time, ":")
			if len(timeStrings) == 2 {
				hour, err := strconv.ParseInt(timeStrings[0], 10, 64)
				if err != nil {
					return nil, err
				}
				bkHour = int(hour)

				minute, err := strconv.ParseInt(timeStrings[1], 10, 64)
				if err != nil {
					return nil, err
				}
				bkMinute = int(minute)
			}
		}
	}

	var bkv2Hour, bkv2Minute int
	var bkv2Weekdays []types.EDayOfTheWeek
	var bkv2Connect string
	if current.Backupv2Setting != nil {
		bkv2Weekdays = current.Backupv2Setting.DayOfWeek
		bkv2Connect = current.Backupv2Setting.Connect
		if current.Backupv2Setting.Time != "" {
			timeStrings := strings.Split(current.Backupv2Setting.Time, ":")
			if len(timeStrings) == 2 {
				hour, err := strconv.ParseInt(timeStrings[0], 10, 64)
				if err != nil {
					return nil, err
				}
				bkv2Hour = int(hour)

				minute, err := strconv.ParseInt(timeStrings[1], 10, 64)
				if err != nil {
					return nil, err
				}
				bkv2Minute = int(minute)
			}
		}
	}

	applyRequest := &ApplyRequest{
		Zone:                  req.Zone,
		ID:                    req.ID,
		Name:                  current.Name,
		Description:           current.Description,
		Tags:                  current.Tags,
		IconID:                current.IconID,
		PlanID:                current.PlanID,
		SwitchID:              current.SwitchID,
		IPAddresses:           current.IPAddresses,
		NetworkMaskLen:        current.NetworkMaskLen,
		DefaultRoute:          current.DefaultRoute,
		Port:                  current.CommonSetting.ServicePort,
		SourceNetwork:         current.CommonSetting.SourceNetwork,
		DatabaseType:          current.Conf.DatabaseName,
		DatabaseVersion:       current.Conf.DatabaseVersion,
		Username:              current.CommonSetting.DefaultUser,
		Password:              current.CommonSetting.UserPassword,
		EnableReplication:     current.ReplicationSetting != nil,
		ReplicaUserPassword:   current.CommonSetting.ReplicaPassword,
		EnableWebUI:           current.CommonSetting.WebUI.Bool(),
		EnableBackup:          current.BackupSetting != nil,
		BackupWeekdays:        bkWeekdays,
		BackupStartTimeHour:   bkHour,
		BackupStartTimeMinute: bkMinute,

		EnableBackupv2:          current.Backupv2Setting != nil,
		Backupv2Weekdays:        bkv2Weekdays,
		Backupv2StartTimeHour:   bkv2Hour,
		Backupv2StartTimeMinute: bkv2Minute,
		Backupv2Connect:         bkv2Connect,

		EnableMonitoringSuite: current.MonitoringSuite != nil && current.MonitoringSuite.Enabled,

		NoWait: false,
	}

	if err := serviceutil.RequestConvertTo(req, applyRequest); err != nil {
		return nil, err
	}
	// パラメータは手動マージ
	parameter, err := dbOp.GetParameter(ctx, req.Zone, req.ID)
	if err != nil {
		return nil, err
	}
	// パラメータ設定をLabelをキーにするように正規化
	ps := make(map[string]interface{})
	for k, v := range parameter.Settings {
		for _, meta := range parameter.MetaInfo {
			if meta.Name == k {
				ps[meta.Label] = v
			}
		}
	}
	if req.Parameters != nil {
		for k, v := range *req.Parameters {
			key := k
			for _, meta := range parameter.MetaInfo {
				if meta.Name == key {
					key = meta.Label
					break
				}
			}
			ps[key] = v
		}
	}
	applyRequest.Parameters = ps

	return applyRequest, nil
}
