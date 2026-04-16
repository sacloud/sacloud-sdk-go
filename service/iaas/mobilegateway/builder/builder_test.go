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

func TestMobileGatewayBuilder_Build(t *testing.T) {
	var switchID types.ID
	var testZone = testutil.TestZone()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			swOp := iaas.NewSwitchOp(caller)

			sw, err := swOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{
				Name: testutil.ResourceName("mobile-gateway-builder"),
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
					Zone:        testZone,
					Name:        testutil.ResourceName("mobile-gateway-builder"),
					Description: "description",
					Tags:        types.Tags{"tag1", "tag2"},
					PrivateInterface: &PrivateInterfaceSetting{
						SwitchID:       switchID,
						IPAddress:      "192.168.0.1",
						NetworkMaskLen: 24,
					},
					StaticRoutes: []*iaas.MobileGatewayStaticRoute{
						{
							Prefix:  "192.168.1.0/24",
							NextHop: "192.168.0.1",
						},
						{
							Prefix:  "192.168.2.0/24",
							NextHop: "192.168.0.1",
						},
					},
					SIMs:                            nil,
					SIMRoutes:                       nil,
					InternetConnectionEnabled:       true,
					InterDeviceCommunicationEnabled: true,
					DNS: &iaas.MobileGatewayDNSSetting{
						DNS1: "1.1.1.1",
						DNS2: "2.2.2.2",
					},
					TrafficConfig: &iaas.MobileGatewayTrafficControl{
						TrafficQuotaInMB:     1024,
						BandWidthLimitInKbps: 128,
						EmailNotifyEnabled:   true,
						AutoTrafficShaping:   true,
					},
					SetupOptions: getSetupOption(),
					Client:       NewAPIClient(caller),
				}
				return builder.Build(ctx)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				mgwOp := iaas.NewMobileGatewayOp(caller)
				return mgwOp.Read(ctx, testZone, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				mgw := value.(*iaas.MobileGateway)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, mgw, "MobileGateway"),
					testutil.AssertLenFunc(t, mgw.InterfaceSettings, 1, "MobileGateway.InterfaceSettings"),
				)
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				mgwOp := iaas.NewMobileGatewayOp(caller)
				return mgwOp.Delete(ctx, testZone, ctx.ID)
			},
		},
		Cleanup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			swOp := iaas.NewSwitchOp(caller)
			return swOp.Delete(ctx, testZone, switchID)
		},
	})
}
