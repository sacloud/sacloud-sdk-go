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

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/database/builder"
)

func BuilderFromResource(ctx context.Context, caller iaas.APICaller, zone string, id types.ID) (*builder.Builder, error) {
	client := iaas.NewDatabaseOp(caller)
	current, err := client.Read(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	parameters, err := client.GetParameter(ctx, zone, id)
	if err != nil {
		return nil, err
	}

	var backupv2Settings *iaas.DatabaseSettingBackupv2
	if current.Backupv2Setting != nil {
		backupv2Settings = &iaas.DatabaseSettingBackupv2{
			Rotate:    current.Backupv2Setting.Rotate,
			Time:      current.Backupv2Setting.Time,
			DayOfWeek: current.Backupv2Setting.DayOfWeek,
			Connect:   current.Backupv2Setting.Connect,
		}
	}

	return &builder.Builder{
		ID:                 current.ID,
		Zone:               zone,
		PlanID:             current.PlanID,
		SwitchID:           current.SwitchID,
		IPAddresses:        current.IPAddresses,
		NetworkMaskLen:     current.NetworkMaskLen,
		DefaultRoute:       current.DefaultRoute,
		Conf:               current.Conf,
		CommonSetting:      current.CommonSetting,
		BackupSetting:      current.BackupSetting,
		Backupv2Setting:    backupv2Settings,
		ReplicationSetting: current.ReplicationSetting,
		MonitoringSuite:    current.MonitoringSuite,
		Disk:               current.Disk,
		Name:               current.Name,
		Description:        current.Description,
		Tags:               current.Tags,
		IconID:             current.IconID,
		Parameters:         parameters.Settings,
		SettingsHash:       current.SettingsHash,
		Client:             builder.NewAPIClient(caller),
	}, nil
}
