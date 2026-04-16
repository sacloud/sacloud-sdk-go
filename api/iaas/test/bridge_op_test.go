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
)

func TestBridgeOpCRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel: true,

		SetupAPICallerFunc: singletonAPICaller,

		Create: &testutil.CRUDTestFunc{
			Func: testBridgeCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createBridgeExpected,
				IgnoreFields: ignoreBridgeFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testBridgeRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createBridgeExpected,
				IgnoreFields: ignoreBridgeFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testBridgeUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateBridgeExpected,
					IgnoreFields: ignoreBridgeFields,
				}),
			},
			{
				Func: testBridgeUpdateToMin,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateBridgeToMinExpected,
					IgnoreFields: ignoreBridgeFields,
				}),
			},
		},

		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testBridgeDelete,
		},
	})
}

var (
	ignoreBridgeFields = []string{
		"ID",
		"CreatedAt",
		"Region",
		"SwitchInZone",
		"BridgeInfo",
	}

	createBridgeParam = &iaas.BridgeCreateRequest{
		Name:        testutil.ResourceName("bridge"),
		Description: "desc",
	}
	createBridgeExpected = &iaas.Bridge{
		Name:        createBridgeParam.Name,
		Description: createBridgeParam.Description,
	}
	updateBridgeParam = &iaas.BridgeUpdateRequest{
		Name:        testutil.ResourceName("bridge-upd"),
		Description: "desc-upd",
	}
	updateBridgeExpected = &iaas.Bridge{
		Name:        updateBridgeParam.Name,
		Description: updateBridgeParam.Description,
	}
	updateBridgeToMinParam = &iaas.BridgeUpdateRequest{
		Name: testutil.ResourceName("bridge-to-min"),
	}
	updateBridgeToMinExpected = &iaas.Bridge{
		Name: updateBridgeToMinParam.Name,
	}
)

func testBridgeCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewBridgeOp(caller)
	return client.Create(ctx, testZone, createBridgeParam)
}

func testBridgeRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewBridgeOp(caller)
	return client.Read(ctx, testZone, ctx.ID)
}

func testBridgeUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewBridgeOp(caller)
	return client.Update(ctx, testZone, ctx.ID, updateBridgeParam)
}

func testBridgeUpdateToMin(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewBridgeOp(caller)
	return client.Update(ctx, testZone, ctx.ID, updateBridgeToMinParam)
}

func testBridgeDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewBridgeOp(caller)
	return client.Delete(ctx, testZone, ctx.ID)
}
