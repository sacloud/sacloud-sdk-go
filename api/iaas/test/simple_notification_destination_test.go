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
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestSimpleNotificationDestinationOp_CRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: testSimpleNotificationDestinationCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createSimpleNotificationDestinationExpected,
				IgnoreFields: ignoreSimpleNotificationDestinationFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testSimpleNotificationDestinationRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createSimpleNotificationDestinationExpected,
				IgnoreFields: ignoreSimpleNotificationDestinationFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testSimpleNotificationDestinationUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateSimpleNotificationDestinationExpected,
					IgnoreFields: ignoreSimpleNotificationDestinationFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testSimpleNotificationDestinationDelete,
		},
	})
}

var (
	ignoreSimpleNotificationDestinationFields = []string{
		"ID",
		"Class",
		"SettingsHash",
		"CreatedAt",
		"ModifiedAt",
		"Status",
	}
	createSimpleNotificationDestinationParam = &iaas.SimpleNotificationDestinationCreateRequest{
		Name:        testutil.ResourceName("simple-notification-destination"),
		Description: "desc",
		Tags:        []string{"tag1", "tag2"},
		Type:        types.SimpleNotificationDestinationTypes.EMail,
		Disabled:    false,
		Value:       "foobar@exaple.com",
	}
	createSimpleNotificationDestinationExpected = &iaas.SimpleNotificationDestination{
		Name:         createSimpleNotificationDestinationParam.Name,
		Description:  createSimpleNotificationDestinationParam.Description,
		Tags:         createSimpleNotificationDestinationParam.Tags,
		Availability: types.Availabilities.Available,
		Type:         createSimpleNotificationDestinationParam.Type,
		Disabled:     createSimpleNotificationDestinationParam.Disabled,
		Value:        createSimpleNotificationDestinationParam.Value,
	}
	updateSimpleNotificationDestinationParam = &iaas.SimpleNotificationDestinationUpdateRequest{
		Name:        testutil.ResourceName("auto-scale-upd"),
		Description: "desc-upd",
		Tags:        []string{"tag1-upd", "tag2-upd"},
		IconID:      testIconID,

		Disabled: true,
	}
	updateSimpleNotificationDestinationExpected = &iaas.SimpleNotificationDestination{
		Name:         updateSimpleNotificationDestinationParam.Name,
		Description:  updateSimpleNotificationDestinationParam.Description,
		Tags:         updateSimpleNotificationDestinationParam.Tags,
		Availability: types.Availabilities.Available,
		IconID:       testIconID,
		Type:         createSimpleNotificationDestinationParam.Type,
		Value:        createSimpleNotificationDestinationParam.Value,

		Disabled: updateSimpleNotificationDestinationParam.Disabled,
	}
)

func testSimpleNotificationDestinationCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSimpleNotificationDestinationOp(caller)
	return client.Create(ctx, createSimpleNotificationDestinationParam)
}

func testSimpleNotificationDestinationRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSimpleNotificationDestinationOp(caller)
	return client.Read(ctx, ctx.ID)
}

func testSimpleNotificationDestinationUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewSimpleNotificationDestinationOp(caller)
	return client.Update(ctx, ctx.ID, updateSimpleNotificationDestinationParam)
}

func testSimpleNotificationDestinationDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewSimpleNotificationDestinationOp(caller)
	return client.Delete(ctx, ctx.ID)
}
