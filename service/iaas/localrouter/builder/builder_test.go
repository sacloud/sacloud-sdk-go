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

package localrouter

import (
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestLocalRouterBuilder_Build(t *testing.T) {
	var testZone = testutil.TestZone()
	var peerLocalRouter *iaas.LocalRouter
	var sw *iaas.Switch

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			lrOp := iaas.NewLocalRouterOp(caller)
			lr, err := lrOp.Create(ctx, &iaas.LocalRouterCreateRequest{
				Name: testutil.ResourceName("local-router-builder"),
			})
			if err != nil {
				return err
			}
			peerLocalRouter = lr

			swOp := iaas.NewSwitchOp(caller)
			sw, err = swOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{
				Name: testutil.ResourceName("local-router-builder"),
			})
			return err
		},
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				builder := &Builder{
					Name:        testutil.ResourceName("local-router-builder"),
					Description: "description",
					Tags:        types.Tags{"tag1", "tag2"},
					Switch: &iaas.LocalRouterSwitch{
						Code:     sw.ID.String(),
						Category: "cloud",
						ZoneID:   testZone,
					},
					Interface: &iaas.LocalRouterInterface{
						VirtualIPAddress: "192.168.0.1",
						IPAddress:        []string{"192.168.0.11", "192.168.0.12"},
						NetworkMaskLen:   24,
						VRID:             101,
					},
					Peers: []*iaas.LocalRouterPeer{
						{
							ID:        peerLocalRouter.ID,
							SecretKey: peerLocalRouter.SecretKeys[0],
							Enabled:   true,
						},
					},
					StaticRoutes: []*iaas.LocalRouterStaticRoute{
						{
							Prefix:  "192.168.1.0/24",
							NextHop: "192.168.0.101",
						},
					},
					Client: NewAPIClient(caller),
				}
				return builder.Build(ctx)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				return iaas.NewLocalRouterOp(caller).Read(ctx, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				lr := value.(*iaas.LocalRouter)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, lr, "LocalRouter"),
					testutil.AssertLenFunc(t, lr.Peers, 1, "LocalRouter.Peers"),
				)
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				lrOp := iaas.NewLocalRouterOp(caller)
				if err := lrOp.Delete(ctx, ctx.ID); err != nil {
					return err
				}
				if err := lrOp.Delete(ctx, peerLocalRouter.ID); err != nil {
					return err
				}
				iaas.NewSwitchOp(caller).Delete(ctx, testZone, sw.ID) //nolint
				return nil
			},
		},
	})
}

func TestLocalRouterBuilder_minimum(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				builder := &Builder{
					Name:   testutil.ResourceName("local-router-builder"),
					Client: NewAPIClient(caller),
				}
				return builder.Build(ctx)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				return iaas.NewLocalRouterOp(caller).Read(ctx, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				lr := value.(*iaas.LocalRouter)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, lr, "LocalRouter"),
				)
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				lrOp := iaas.NewLocalRouterOp(caller)
				return lrOp.Delete(ctx, ctx.ID)
			},
		},
	})
}
