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
	"crypto/rand"
	"encoding/hex"
	"testing"
	"time"

	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/helper/power"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/helper/wait"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/testutil"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
	builder2 "github.com/sacloud/sacloud-sdk-go/service/iaas/database/builder"
	"github.com/sacloud/sacloud-sdk-go/service/iaas/setup"
	"github.com/stretchr/testify/require"
)

func TestDatabaseService_UpdateWithContext(t *testing.T) {
	if !testutil.IsAccTest() {
		t.Skip("This test runs only with TESTACC=1")
	}

	passwordBytes := make([]byte, 15)
	if _, err := rand.Read(passwordBytes); err != nil {
		t.Fatal(err)
	}
	newPassword := hex.EncodeToString(passwordBytes)

	zone := testutil.TestZone()
	name := testutil.ResourceName("database-update-password")
	updateName := name + "-upd"
	var switchID types.ID
	svc := New(testutil.SingletonAPICaller())
	weekdays := make([]types.EDayOfTheWeek, 0)

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			if err := testutil.ComposeCleanupResourceFunc(
				name,
				testutil.CleanupTargets.Database,
				testutil.CleanupTargets.Switch,
			)(ctx, caller); err != nil {
				return err
			}
			sw, err := iaas.NewSwitchOp(caller).Create(ctx, zone, &iaas.SwitchCreateRequest{Name: name})
			if err != nil {
				return err
			}
			switchID = sw.ID
			return nil
		},
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (any, error) {
				return (&builder2.Builder{
					Zone:           zone,
					PlanID:         types.DatabasePlans.DB10GB,
					SwitchID:       switchID,
					IPAddresses:    []string{"192.168.0.101"},
					NetworkMaskLen: 24,
					DefaultRoute:   "192.168.0.1",
					Conf: &iaas.DatabaseRemarkDBConfCommon{
						DatabaseName: types.RDBMSTypesPostgreSQL.String(),
						DefaultUser:  "default",
						UserPassword: "initial-password",
					},
					CommonSetting: &iaas.DatabaseSettingCommon{
						ServicePort:  5432,
						DefaultUser:  "default",
						UserPassword: "initial-password",
					},
					Name:   name,
					Client: builder2.NewAPIClient(caller),
				}).Build(ctx)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (any, error) {
				return iaas.NewDatabaseOp(caller).Read(ctx, zone, ctx.ID)
			},
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (any, error) {
					t.Log("updating database without setting Password")
					return svc.UpdateWithContext(ctx, &UpdateRequest{
						Zone:             zone,
						ID:               ctx.ID,
						Name:             new(updateName),
						Backupv2Weekdays: &weekdays,
						NoWait:           true,
					})
				},
				CheckFunc: func(t testutil.TestT, _ *testutil.CRUDTestContext, value any) error {
					db := value.(*iaas.Database)
					return testutil.AssertEqual(t, updateName, db.Name, "Database.Name")
				},
			},
			{
				Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (any, error) {
					t.Log("updating database with a new Password")
					return svc.UpdateWithContext(ctx, &UpdateRequest{
						Zone:             zone,
						ID:               ctx.ID,
						Password:         &newPassword,
						Backupv2Weekdays: &weekdays,
						NoWait:           true,
					})
				},
				CheckFunc: func(t testutil.TestT, _ *testutil.CRUDTestContext, value any) error {
					db := value.(*iaas.Database)
					return testutil.AssertEqual(t, updateName, db.Name, "Database.Name")
				},
			},
		},
		Shutdown: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			return power.ShutdownDatabase(ctx, iaas.NewDatabaseOp(caller), zone, ctx.ID, true)
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				return iaas.NewDatabaseOp(caller).Delete(ctx, zone, ctx.ID)
			},
		},
		Cleanup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			return testutil.ComposeCleanupResourceFunc(
				name,
				testutil.CleanupTargets.Database,
				testutil.CleanupTargets.Switch,
			)(ctx, caller)
		},
	})
}

func TestDatabaseService_convertUpdateRequest(t *testing.T) {
	if testutil.IsAccTest() {
		t.Skip("This test runs only without TESTACC=1")
	}
	for _, tc := range []struct {
		databaseName string
		databaseType string
	}{
		{databaseName: "MariaDB", databaseType: "mariadb"},
		{databaseName: "PostgreSQL", databaseType: "postgres"},
		{databaseName: "postgres", databaseType: "postgres"},
	} {
		t.Run(tc.databaseName, func(t *testing.T) {
			testDatabaseServiceConvertUpdateRequest(t, tc.databaseName, tc.databaseType)
		})
	}
}

func testDatabaseServiceConvertUpdateRequest(t *testing.T, databaseName, databaseType string) {
	t.Helper()
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
			DatabaseName: databaseName,
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
		Parameters: map[string]any{
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
				Name:   new(db.Name + "-upd"),
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
				DatabaseType:          databaseType,
				Username:              "default",
				Password:              "",
				EnableReplication:     true,
				ReplicaUserPassword:   "password2",
				EnableWebUI:           true,
				EnableBackup:          true,
				BackupWeekdays:        []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
				BackupStartTimeHour:   10,
				BackupStartTimeMinute: 0,
				Parameters: map[string]any{
					"max_connections": float64(100),
				},
				NoWait: true,
			},
		},
		{
			in: &UpdateRequest{
				Zone:              zone,
				ID:                db.ID,
				Password:          new("new-password"),
				EnableReplication: new(false),
				EnableBackup:      new(false),
				Parameters: &map[string]any{
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
				DatabaseType:          databaseType,
				Username:              "default",
				Password:              "new-password",
				EnableReplication:     false,
				ReplicaUserPassword:   "password2",
				EnableWebUI:           true,
				EnableBackup:          false,
				BackupWeekdays:        []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
				BackupStartTimeHour:   10,
				BackupStartTimeMinute: 0,
				Parameters: map[string]any{
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
		require.NoError(t, got.Validate())

		// データベースバージョンは比較対象としない
		got.DatabaseVersion = ""
		require.EqualValues(t, tc.expect, got)
	}
}
