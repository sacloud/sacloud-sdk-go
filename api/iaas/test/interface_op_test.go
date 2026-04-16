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
	"github.com/sacloud/packages-go/size"
	"github.com/stretchr/testify/assert"
)

func TestInterface_Operations(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:          true,
		IgnoreStartupWait: true,

		SetupAPICallerFunc: singletonAPICaller,

		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			serverClient := iaas.NewServerOp(caller)
			server, err := serverClient.Create(ctx, testZone, &iaas.ServerCreateRequest{
				CPU:      1,
				MemoryMB: 1 * size.GiB,
				//ConnectedSwitches: []*ConnectedSwitch{
				//	{Scope: types.Scopes.Shared},
				// },
				Commitment: types.Commitments.Standard,
				Name:       testutil.ResourceName("server-with-interface"),
			})
			if !assert.NoError(t, err) {
				return err
			}

			ctx.Values["interface/server"] = server.ID
			createInterfaceParam.ServerID = server.ID
			createInterfaceExpected.ServerID = server.ID
			updateInterfaceExpected.ServerID = server.ID
			return nil
		},

		Create: &testutil.CRUDTestFunc{
			Func: testInterfaceCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createInterfaceExpected,
				IgnoreFields: ignoreInterfaceFields,
			}),
		},
		Read: &testutil.CRUDTestFunc{
			Func: testInterfaceRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createInterfaceExpected,
				IgnoreFields: ignoreInterfaceFields,
			}),
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testInterfaceUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateInterfaceExpected,
					IgnoreFields: ignoreInterfaceFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testInterfaceDelete,
		},

		Cleanup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			serverID, ok := ctx.Values["interface/server"]
			if !ok {
				return nil
			}

			serverClient := iaas.NewServerOp(caller)
			return serverClient.Delete(ctx, testZone, serverID.(types.ID))
		},
	})
}

var (
	ignoreInterfaceFields = []string{
		"ID",
		"MACAddress",
		"IPAddress",
		"CreatedAt",
		"ModifiedAt",
	}
	createInterfaceParam = &iaas.InterfaceCreateRequest{}

	createInterfaceExpected = &iaas.Interface{
		UserIPAddress:  "",
		HostName:       "",
		SwitchID:       types.ID(0),
		PacketFilterID: types.ID(0),
	}
	updateInterfaceParam = &iaas.InterfaceUpdateRequest{
		UserIPAddress: "192.2.0.1",
	}
	updateInterfaceExpected = &iaas.Interface{
		UserIPAddress:  "192.2.0.1",
		HostName:       "",
		SwitchID:       types.ID(0),
		PacketFilterID: types.ID(0),
	}
)

func testInterfaceCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewInterfaceOp(caller)
	return client.Create(ctx, testZone, createInterfaceParam)
}

func testInterfaceRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewInterfaceOp(caller)
	return client.Read(ctx, testZone, ctx.ID)
}

func testInterfaceUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewInterfaceOp(caller)
	return client.Update(ctx, testZone, ctx.ID, updateInterfaceParam)
}

func testInterfaceDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewInterfaceOp(caller)
	return client.Delete(ctx, testZone, ctx.ID)
}
