// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package test

import (
	"fmt"
	"os"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/size"
)

var autoScaleTestServerName = testutil.ResourceName("auto-scale")

func TestAutoScaleOp_CRUD(t *testing.T) {
	if testutil.IsAccTest() && os.Getenv("SAKURACLOUD_API_KEY_ID") == "" {
		t.Skip("SAKURACLOUD_API_KEY_ID is required when running the acceptance test")
	}

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		SetupAPICallerFunc: singletonAPICaller,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			serverOp := iaas.NewServerOp(caller)
			// ディスクレスサーバを作成
			_, err := serverOp.Create(ctx, testutil.TestZone(), &iaas.ServerCreateRequest{
				CPU:      1,
				MemoryMB: 2 * size.GiB,
				Name:     autoScaleTestServerName,
			})
			return err
		},
		Create: &testutil.CRUDTestFunc{
			Func: testAutoScaleCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createAutoScaleExpected,
				IgnoreFields: ignoreAutoScaleFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testAutoScaleRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createAutoScaleExpected,
				IgnoreFields: ignoreAutoScaleFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testAutoScaleUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateAutoScaleExpected,
					IgnoreFields: ignoreAutoScaleFields,
				}),
			},
			{
				Func: testAutoScaleUpdateTriggerType,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateAutoScaleTriggerTypeExpected,
					IgnoreFields: ignoreAutoScaleFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testAutoScaleDelete,
		},
	})
}

var (
	ignoreAutoScaleFields = []string{
		"ID",
		"Class",
		"SettingsHash",
		"CreatedAt",
		"ModifiedAt",
		"ScheduleScaling",
	}
	createAutoScaleParam = &iaas.AutoScaleCreateRequest{
		Name:        testutil.ResourceName("auto-scale"),
		Description: "desc",
		Tags:        []string{"tag1", "tag2"},

		Config:      fmt.Sprintf(autoScaleConfigTemplate, autoScaleTestServerName, testutil.TestZone()),
		Zones:       []string{testutil.TestZone()},
		TriggerType: types.AutoScaleTriggerTypes.CPU,
		Disabled:    true,
		CPUThresholdScaling: &iaas.AutoScaleCPUThresholdScaling{
			ServerPrefix: autoScaleTestServerName,
			Up:           80,
			Down:         50,
		},
		APIKeyID: os.Getenv("SAKURACLOUD_API_KEY_ID"),
	}
	createAutoScaleExpected = &iaas.AutoScale{
		Name:         createAutoScaleParam.Name,
		Description:  createAutoScaleParam.Description,
		Tags:         createAutoScaleParam.Tags,
		Availability: types.Availabilities.Available,

		TriggerType: types.AutoScaleTriggerTypes.CPU,
		Disabled:    true,
		Config:      fmt.Sprintf(autoScaleConfigTemplate, autoScaleTestServerName, testutil.TestZone()),
		Zones:       []string{testutil.TestZone()},
		CPUThresholdScaling: &iaas.AutoScaleCPUThresholdScaling{
			ServerPrefix: autoScaleTestServerName,
			Up:           80,
			Down:         50,
		},
		APIKeyID: os.Getenv("SAKURACLOUD_API_KEY_ID"),
	}
	updateAutoScaleParam = &iaas.AutoScaleUpdateRequest{
		Name:        testutil.ResourceName("auto-scale-upd"),
		Description: "desc-upd",
		Tags:        []string{"tag1-upd", "tag2-upd"},
		IconID:      testIconID,

		Config:      fmt.Sprintf(autoScaleConfigTemplateUpd, autoScaleTestServerName, testutil.TestZone()),
		Zones:       []string{testutil.TestZone()},
		TriggerType: types.AutoScaleTriggerTypes.CPU,
		Disabled:    false,
		CPUThresholdScaling: &iaas.AutoScaleCPUThresholdScaling{
			ServerPrefix: autoScaleTestServerName,
			Up:           81,
			Down:         51,
		},
	}
	updateAutoScaleExpected = &iaas.AutoScale{
		Name:         updateAutoScaleParam.Name,
		Description:  updateAutoScaleParam.Description,
		Tags:         updateAutoScaleParam.Tags,
		Availability: types.Availabilities.Available,
		IconID:       testIconID,

		Config:      fmt.Sprintf(autoScaleConfigTemplateUpd, autoScaleTestServerName, testutil.TestZone()),
		Zones:       []string{testutil.TestZone()},
		TriggerType: types.AutoScaleTriggerTypes.CPU,
		Disabled:    false,
		CPUThresholdScaling: &iaas.AutoScaleCPUThresholdScaling{
			ServerPrefix: autoScaleTestServerName,
			Up:           81,
			Down:         51,
		},
		APIKeyID: os.Getenv("SAKURACLOUD_API_KEY_ID"),
	}
	updateAutoScaleTriggerTypeParam = &iaas.AutoScaleUpdateRequest{
		Name:        testutil.ResourceName("auto-scale-upd"),
		Description: "desc-upd",
		Tags:        []string{"tag1-upd", "tag2-upd"},
		IconID:      testIconID,

		Config:      fmt.Sprintf(autoScaleConfigTemplateUpd, autoScaleTestServerName, testutil.TestZone()),
		Zones:       []string{testutil.TestZone()},
		TriggerType: types.AutoScaleTriggerTypes.Schedule,
		Disabled:    false,
		ScheduleScaling: []*iaas.AutoScaleScheduleScaling{
			{
				Action:    types.AutoScaleActions.Up,
				Hour:      10,
				Minute:    15,
				DayOfWeek: []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
			},
			{
				Action:    types.AutoScaleActions.Down,
				Hour:      18,
				Minute:    15,
				DayOfWeek: []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
			},
		},
	}
	updateAutoScaleTriggerTypeExpected = &iaas.AutoScale{
		Name:         updateAutoScaleParam.Name,
		Description:  updateAutoScaleParam.Description,
		Tags:         updateAutoScaleParam.Tags,
		Availability: types.Availabilities.Available,
		IconID:       testIconID,

		Config:      fmt.Sprintf(autoScaleConfigTemplateUpd, autoScaleTestServerName, testutil.TestZone()),
		Zones:       []string{testutil.TestZone()},
		TriggerType: types.AutoScaleTriggerTypes.Schedule,
		Disabled:    false,
		ScheduleScaling: []*iaas.AutoScaleScheduleScaling{
			{
				Action:    types.AutoScaleActions.Up,
				Hour:      10,
				Minute:    15,
				DayOfWeek: []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
			},
			{
				Action:    types.AutoScaleActions.Down,
				Hour:      18,
				Minute:    15,
				DayOfWeek: []types.EDayOfTheWeek{types.DaysOfTheWeek.Monday},
			},
		},
		APIKeyID: os.Getenv("SAKURACLOUD_API_KEY_ID"),
	}
)

func testAutoScaleCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewAutoScaleOp(caller)
	return client.Create(ctx, createAutoScaleParam)
}

func testAutoScaleRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewAutoScaleOp(caller)
	return client.Read(ctx, ctx.ID)
}

func testAutoScaleUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewAutoScaleOp(caller)
	return client.Update(ctx, ctx.ID, updateAutoScaleParam)
}

func testAutoScaleUpdateTriggerType(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewAutoScaleOp(caller)
	return client.Update(ctx, ctx.ID, updateAutoScaleTriggerTypeParam)
}

func testAutoScaleDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewAutoScaleOp(caller)
	return client.Delete(ctx, ctx.ID)
}

var autoScaleConfigTemplate = `
resources:
  - type: Server
    selector:
      names: ["%s"]
      zones: ["%s"]

    shutdown_force: true
`

var autoScaleConfigTemplateUpd = `
resources:
  - type: Server
    selector:
      names: ["%s"]
      zones: ["%s"]

    shutdown_force: true

autoscaler:
  cooldown: 300
`
