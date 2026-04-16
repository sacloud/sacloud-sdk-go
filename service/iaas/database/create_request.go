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
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/validate"
)

type CreateRequest struct {
	Zone string `service:"-" validate:"required"`

	Name                  string `validate:"required"`
	Description           string `validate:"min=0,max=512"`
	Tags                  types.Tags
	IconID                types.ID
	PlanID                types.ID `validate:"required"`
	SwitchID              types.ID `validate:"required"`
	IPAddresses           []string `validate:"required,min=1,max=2,dive,ipv4"`
	NetworkMaskLen        int      `validate:"required,min=1,max=32"`
	DefaultRoute          string   `validate:"omitempty,ipv4"`
	Port                  int      `validate:"omitempty,min=1,max=65535"`
	SourceNetwork         []string `validate:"omitempty,dive,cidrv4"`
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

func (req *CreateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *CreateRequest) ApplyRequest() *ApplyRequest {
	return &ApplyRequest{
		Zone:                    req.Zone,
		Name:                    req.Name,
		Description:             req.Description,
		Tags:                    req.Tags,
		IconID:                  req.IconID,
		PlanID:                  req.PlanID,
		SwitchID:                req.SwitchID,
		IPAddresses:             req.IPAddresses,
		NetworkMaskLen:          req.NetworkMaskLen,
		DefaultRoute:            req.DefaultRoute,
		Port:                    req.Port,
		SourceNetwork:           req.SourceNetwork,
		DatabaseType:            req.DatabaseType,
		DatabaseVersion:         req.DatabaseVersion,
		Username:                req.Username,
		Password:                req.Password,
		EnableReplication:       req.EnableBackup,
		ReplicaUserPassword:     req.ReplicaUserPassword,
		EnableWebUI:             req.EnableWebUI,
		EnableBackup:            req.EnableBackup,
		BackupWeekdays:          req.BackupWeekdays,
		BackupStartTimeHour:     req.BackupStartTimeHour,
		BackupStartTimeMinute:   req.BackupStartTimeMinute,
		Parameters:              req.Parameters,
		EnableBackupv2:          req.EnableBackupv2,
		Backupv2Weekdays:        req.Backupv2Weekdays,
		Backupv2StartTimeHour:   req.Backupv2StartTimeHour,
		Backupv2StartTimeMinute: req.Backupv2StartTimeMinute,
		Backupv2Connect:         req.Backupv2Connect,
		EnableMonitoringSuite:   req.EnableMonitoringSuite,
		DiskEncryptionAlgorithm: req.DiskEncryptionAlgorithm,
		DiskEncryptionKMSKey:    req.DiskEncryptionKMSKey,
		NoWait:                  req.NoWait,
	}
}
