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
	"testing"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/setup"
)

func getSetupOption() *setup.Options {
	if testutil.IsAccTest() {
		return nil
	}
	return &setup.Options{
		DeleteRetryInterval:       10 * time.Millisecond,
		ProvisioningRetryInterval: 10 * time.Millisecond,
		PollingInterval:           10 * time.Millisecond,
		NICUpdateWaitDuration:     10 * time.Millisecond,
	}
}

func TestDatabaseBuilder_Build(t *testing.T) {
	var switchID types.ID
	var testZone = testutil.TestZone()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			swOp := iaas.NewSwitchOp(caller)

			sw, err := swOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{
				Name: testutil.ResourceName("database-builder"),
			})
			if err != nil {
				return err
			}
			switchID = sw.ID
			return nil
		},
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				builder := &Builder{
					Zone:           testZone,
					PlanID:         types.DatabasePlans.DB10GB,
					SwitchID:       switchID,
					IPAddresses:    []string{"192.168.0.11"},
					NetworkMaskLen: 24,
					DefaultRoute:   "192.168.0.1",
					Conf: &iaas.DatabaseRemarkDBConfCommon{
						DatabaseName: types.RDBMSTypesPostgreSQL.String(),
						DefaultUser:  "builder",
						UserPassword: "builder-password-dummy",
					},
					CommonSetting: &iaas.DatabaseSettingCommon{
						DefaultUser:     "builder",
						UserPassword:    "builder-password-dummy",
						ReplicaUser:     "",
						ReplicaPassword: "",
					},
					BackupSetting: &iaas.DatabaseSettingBackup{
						Rotate:    8,
						Time:      "00:00",
						DayOfWeek: []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
					},
					ReplicationSetting: &iaas.DatabaseReplicationSetting{},
					Parameters: map[string]interface{}{
						"max_connections": 50,
					},
					Name:         testutil.ResourceName("database-builder"),
					Description:  "description",
					Tags:         types.Tags{"tag1", "tag2"},
					SetupOptions: getSetupOption(),
					Client:       NewAPIClient(caller),
				}
				return builder.Build(ctx)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				dbOp := iaas.NewDatabaseOp(caller)
				return dbOp.Read(ctx, testZone, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				db := value.(*iaas.Database)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, db, "Database"),
					testutil.AssertNotNilFunc(t, db.Conf, "Database.Conf"),
				)
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				dbOp := iaas.NewDatabaseOp(caller)
				return dbOp.Delete(ctx, testZone, ctx.ID)
			},
		},
		Shutdown: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			dbOp := iaas.NewDatabaseOp(caller)
			return power.ShutdownDatabase(ctx, dbOp, testZone, ctx.ID, true)
		},
		Cleanup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			swOp := iaas.NewSwitchOp(caller)
			return swOp.Delete(ctx, testZone, switchID)
		},
	})
}
