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
	"testing"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/wait"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	builder2 "github.com/sacloud/iaas-service-go/database/builder"
	"github.com/sacloud/iaas-service-go/setup"
	"github.com/sacloud/packages-go/pointer"
	"github.com/stretchr/testify/require"
)

func TestDatabaseService_convertUpdateRequest(t *testing.T) {
	if testutil.IsAccTest() {
		t.Skip("This test runs only without TESTACC=1")
	}
	ctx := context.Background()
	zone := testutil.TestZone()
	name := testutil.ResourceName("database-service-update")
	caller := testutil.SingletonAPICaller()

	sw, err := iaas.NewSwitchOp(caller).Create(ctx, zone, &iaas.SwitchCreateRequest{Name: name})
	if err != nil {
		t.Fatal(err)
	}
	builder := &builder2.Builder{
		Zone:           zone,
		PlanID:         types.DatabasePlans.DB10GB,
		SwitchID:       sw.ID,
		IPAddresses:    []string{"192.168.0.101"},
		NetworkMaskLen: 24,
		DefaultRoute:   "192.168.0.1",
		Conf: &iaas.DatabaseRemarkDBConfCommon{
			DatabaseName: types.RDBMSTypesPostgreSQL.String(),
		},
		SourceID: 0,
		CommonSetting: &iaas.DatabaseSettingCommon{
			WebUI:           types.ToWebUI(true),
			ServicePort:     5432,
			SourceNetwork:   []string{"192.168.0.0/24"},
			DefaultUser:     "default",
			UserPassword:    "password1",
			ReplicaUser:     "replica",
			ReplicaPassword: "password2",
		},
		BackupSetting: &iaas.DatabaseSettingBackup{
			Rotate:    0,
			Time:      "10:00",
			DayOfWeek: []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
		},
		ReplicationSetting: &iaas.DatabaseReplicationSetting{
			Model: types.DatabaseReplicationModels.MasterSlave,
		},
		Name:         name,
		Description:  "description",
		Tags:         types.Tags{"tag1", "tag2"},
		SettingsHash: "",
		NoWait:       false,
		Parameters: map[string]interface{}{
			"max_connections": float64(100),
		},
		SetupOptions: &setup.Options{
			BootAfterBuild:        true,
			NICUpdateWaitDuration: time.Millisecond,
			PollingInterval:       time.Millisecond,
		},
		Client: builder2.NewAPIClient(caller),
	}
	db, err := builder.Build(ctx)
	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		dbOp := iaas.NewDatabaseOp(caller)
		if err := dbOp.Shutdown(ctx, zone, db.ID, &iaas.ShutdownOption{Force: true}); err != nil {
			return
		}
		if _, err := wait.UntilDatabaseIsDown(ctx, dbOp, zone, db.ID); err != nil {
			return
		}
		if err := dbOp.Delete(ctx, zone, db.ID); err != nil {
			return
		}
		iaas.NewSwitchOp(caller).Delete(ctx, zone, sw.ID) //nolint
	}()

	cases := []struct {
		in     *UpdateRequest
		expect *ApplyRequest
	}{
		{
			in: &UpdateRequest{
				Zone:   zone,
				ID:     db.ID,
				Name:   pointer.NewString(db.Name + "-upd"),
				NoWait: true,
			},
			expect: &ApplyRequest{
				Zone:                  zone,
				ID:                    db.ID,
				Name:                  db.Name + "-upd",
				Description:           "description",
				Tags:                  types.Tags{"tag1", "tag2"},
				PlanID:                types.DatabasePlans.DB10GB,
				SwitchID:              sw.ID,
				IPAddresses:           []string{"192.168.0.101"},
				NetworkMaskLen:        24,
				DefaultRoute:          "192.168.0.1",
				Port:                  5432,
				SourceNetwork:         []string{"192.168.0.0/24"},
				DatabaseType:          types.RDBMSTypesPostgreSQL.String(),
				Username:              "default",
				Password:              "password1",
				EnableReplication:     true,
				ReplicaUserPassword:   "password2",
				EnableWebUI:           true,
				EnableBackup:          true,
				BackupWeekdays:        []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
				BackupStartTimeHour:   10,
				BackupStartTimeMinute: 0,
				Parameters: map[string]interface{}{
					"max_connections": float64(100),
				},
				NoWait: true,
			},
		},
		{
			in: &UpdateRequest{
				Zone:              zone,
				ID:                db.ID,
				EnableReplication: pointer.NewBool(false),
				EnableBackup:      pointer.NewBool(false),
				Parameters: &map[string]interface{}{
					"work_mem": float64(4096),
				},
				NoWait: true,
			},
			expect: &ApplyRequest{
				Zone:                  zone,
				ID:                    db.ID,
				Name:                  db.Name,
				Description:           "description",
				Tags:                  types.Tags{"tag1", "tag2"},
				PlanID:                types.DatabasePlans.DB10GB,
				SwitchID:              sw.ID,
				IPAddresses:           []string{"192.168.0.101"},
				NetworkMaskLen:        24,
				DefaultRoute:          "192.168.0.1",
				Port:                  5432,
				SourceNetwork:         []string{"192.168.0.0/24"},
				DatabaseType:          types.RDBMSTypesPostgreSQL.String(),
				Username:              "default",
				Password:              "password1",
				EnableReplication:     false,
				ReplicaUserPassword:   "password2",
				EnableWebUI:           true,
				EnableBackup:          false,
				BackupWeekdays:        []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
				BackupStartTimeHour:   10,
				BackupStartTimeMinute: 0,
				Parameters: map[string]interface{}{
					"max_connections": float64(100),
					"work_mem":        float64(4096),
				},
				NoWait: true,
			},
		},
	}

	for _, tc := range cases {
		got, err := tc.in.ApplyRequest(ctx, caller)
		require.NoError(t, err)

		// データベースバージョンは比較対象としない
		got.DatabaseVersion = ""
		require.EqualValues(t, tc.expect, got)
	}
}
