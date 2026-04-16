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

func TestPrivateHostOp_CRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		IgnoreStartupWait:  true,
		SetupAPICallerFunc: singletonAPICaller,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			planOp := iaas.NewPrivateHostPlanOp(caller)
			searched, err := planOp.Find(ctx, privateHostTestZone, nil)
			if err != nil {
				return err
			}
			planID := searched.PrivateHostPlans[0].ID
			createPrivateHostParam.PlanID = planID
			createPrivateHostExpected.PlanID = planID
			updatePrivateHostExpected.PlanID = planID
			updatePrivateHostToMinExpected.PlanID = planID
			return nil
		},
		Create: &testutil.CRUDTestFunc{
			Func: testPrivateHostCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createPrivateHostExpected,
				IgnoreFields: ignorePrivateHostFields,
			}),
		},
		Read: &testutil.CRUDTestFunc{
			Func: testPrivateHostRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createPrivateHostExpected,
				IgnoreFields: ignorePrivateHostFields,
			}),
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testPrivateHostUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updatePrivateHostExpected,
					IgnoreFields: ignorePrivateHostFields,
				}),
			},
			{
				Func: testPrivateHostUpdateToMin,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updatePrivateHostToMinExpected,
					IgnoreFields: ignorePrivateHostFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testPrivateHostDelete,
		},
	})
}

var (
	privateHostTestZone = "tk1a"

	ignorePrivateHostFields = []string{
		"ID",
		"CreatedAt",
		"PlanName",
		"PlanClass",
		"HostName",
		"CPU",
		"MemoryMB",
	}

	createPrivateHostParam = &iaas.PrivateHostCreateRequest{
		Name:        testutil.ResourceName("private-host"),
		Description: "libsacloud-private-host",
		Tags:        []string{"tag1", "tag2"},
	}
	createPrivateHostExpected = &iaas.PrivateHost{
		Name:             createPrivateHostParam.Name,
		Description:      createPrivateHostParam.Description,
		Tags:             createPrivateHostParam.Tags,
		CPU:              224,
		AssignedCPU:      0,
		AssignedMemoryMB: 0,
	}
	updatePrivateHostParam = &iaas.PrivateHostUpdateRequest{
		Name:        testutil.ResourceName("private-host-upd"),
		Description: "libsacloud-private-host-upd",
		Tags:        []string{"tag1-upd", "tag2-upd"},
		IconID:      testIconID,
	}
	updatePrivateHostExpected = &iaas.PrivateHost{
		Name:             updatePrivateHostParam.Name,
		Description:      updatePrivateHostParam.Description,
		Tags:             updatePrivateHostParam.Tags,
		CPU:              224,
		AssignedCPU:      0,
		AssignedMemoryMB: 0,
		IconID:           testIconID,
	}
	updatePrivateHostToMinParam = &iaas.PrivateHostUpdateRequest{
		Name: testutil.ResourceName("private-host-to-min"),
	}
	updatePrivateHostToMinExpected = &iaas.PrivateHost{
		Name:             updatePrivateHostToMinParam.Name,
		CPU:              224,
		AssignedCPU:      0,
		AssignedMemoryMB: 0,
	}
)

func testPrivateHostCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewPrivateHostOp(caller)
	return client.Create(ctx, privateHostTestZone, createPrivateHostParam)
}

func testPrivateHostRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewPrivateHostOp(caller)
	return client.Read(ctx, privateHostTestZone, ctx.ID)
}

func testPrivateHostUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewPrivateHostOp(caller)
	return client.Update(ctx, privateHostTestZone, ctx.ID, updatePrivateHostParam)
}

func testPrivateHostUpdateToMin(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewPrivateHostOp(caller)
	return client.Update(ctx, privateHostTestZone, ctx.ID, updatePrivateHostToMinParam)
}

func testPrivateHostDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewPrivateHostOp(caller)
	return client.Delete(ctx, privateHostTestZone, ctx.ID)
}
