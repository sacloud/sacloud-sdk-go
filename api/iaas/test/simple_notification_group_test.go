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
	"os"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestSimpleNotificationGroupOp_CRUD(t *testing.T) {
	if testutil.IsAccTest() && os.Getenv("SAKURACLOUD_SIMPLE_NOTIFICATION_DESTINATION_ID") == "" {
		t.Skip("SAKURACLOUD_SIMPLE_NOTIFICATION_DESTINATION_ID is required when running the acceptance test")
	}

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: testSimpleNotificationGroupCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createSimpleNotificationGroupExpected,
				IgnoreFields: ignoreSimpleNotificationGroupFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testSimpleNotificationGroupRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createSimpleNotificationGroupExpected,
				IgnoreFields: ignoreSimpleNotificationGroupFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testSimpleNotificationGroupUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateSimpleNotificationGroupExpected,
					IgnoreFields: ignoreSimpleNotificationGroupFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testSimpleNotificationGroupDelete,
		},
	})
}

var (
	ignoreSimpleNotificationGroupFields = []string{
		"ID",
		"Class",
		"SettingsHash",
		"CreatedAt",
		"ModifiedAt",
		"Status",
	}
	createSimpleNotificationGroupParam = &iaas.SimpleNotificationGroupCreateRequest{
		Name:        testutil.ResourceName("simple-notification-destination"),
		Description: "desc",
		Tags:        []string{"tag1", "tag2"},
		Destinations: []string{
			os.Getenv("SAKURACLOUD_SIMPLE_NOTIFICATION_DESTINATION_ID"),
		},
		Disabled: false,
		Sources:  []string{"1"},
	}
	createSimpleNotificationGroupExpected = &iaas.SimpleNotificationGroup{
		Name:         createSimpleNotificationGroupParam.Name,
		Description:  createSimpleNotificationGroupParam.Description,
		Tags:         createSimpleNotificationGroupParam.Tags,
		Availability: types.Availabilities.Available,
		Destinations: createSimpleNotificationGroupParam.Destinations,
		Disabled:     createSimpleNotificationGroupParam.Disabled,
		Sources:      createSimpleNotificationGroupParam.Sources,
	}
	updateSimpleNotificationGroupParam = &iaas.SimpleNotificationGroupUpdateRequest{
		Name:         testutil.ResourceName("auto-scale-upd"),
		Description:  "desc-upd",
		Tags:         []string{"tag1-upd", "tag2-upd"},
		IconID:       testIconID,
		Destinations: createSimpleNotificationGroupParam.Destinations,
		Sources:      createSimpleNotificationGroupParam.Sources,

		Disabled: true,
	}
	updateSimpleNotificationGroupExpected = &iaas.SimpleNotificationGroup{
		Name:         updateSimpleNotificationGroupParam.Name,
		Description:  updateSimpleNotificationGroupParam.Description,
		Tags:         updateSimpleNotificationGroupParam.Tags,
		Availability: types.Availabilities.Available,
		IconID:       testIconID,
		Destinations: createSimpleNotificationGroupParam.Destinations,
		Sources:      createSimpleNotificationGroupParam.Sources,

		Disabled: updateSimpleNotificationGroupParam.Disabled,
	}
)

func testSimpleNotificationGroupCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSimpleNotificationGroupOp(caller)
	return client.Create(ctx, createSimpleNotificationGroupParam)
}

func testSimpleNotificationGroupRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSimpleNotificationGroupOp(caller)
	return client.Read(ctx, ctx.ID)
}

func testSimpleNotificationGroupUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSimpleNotificationGroupOp(caller)
	return client.Update(ctx, ctx.ID, updateSimpleNotificationGroupParam)
}

func testSimpleNotificationGroupDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewSimpleNotificationGroupOp(caller)
	return client.Delete(ctx, ctx.ID)
}
