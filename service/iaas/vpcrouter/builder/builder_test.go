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
	"github.com/sacloud/iaas-service-go/internet/builder"
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

func TestBuilder_Build(t *testing.T) {
	var switchID types.ID
	var testZone = testutil.TestZone()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			swOp := iaas.NewSwitchOp(caller)

			sw, err := swOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{
				Name: testutil.ResourceName("vpc-router-builder"),
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
					Name:        testutil.ResourceName("vpc-router-builder"),
					Description: "description",
					Tags:        types.Tags{"tag1", "tag2"},
					PlanID:      types.VPCRouterPlans.Standard,
					Version:     1,
					NICSetting:  &StandardNICSetting{},
					AdditionalNICSettings: []AdditionalNICSettingHolder{
						&AdditionalStandardNICSetting{
							SwitchID:       switchID,
							IPAddress:      "192.168.0.1",
							NetworkMaskLen: 24,
							Index:          2,
						},
					},
					RouterSetting: &RouterSetting{
						InternetConnectionEnabled: types.StringTrue,
					},
					SetupOptions: getSetupOption(),
					Client:       iaas.NewVPCRouterOp(caller),
				}
				return builder.Build(ctx)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				vpcRouterOp := iaas.NewVPCRouterOp(caller)
				return vpcRouterOp.Read(ctx, testZone, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				vpcRouter := value.(*iaas.VPCRouter)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, vpcRouter, "VPCRouter"),
					testutil.AssertNotNilFunc(t, vpcRouter.Settings, "VPCRouter.Settings"),
					testutil.AssertLenFunc(t, vpcRouter.Settings.Interfaces, 1, "VPCRouter.Settings.Interfaces"),
				)
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				vpcRouterOp := iaas.NewVPCRouterOp(caller)
				return vpcRouterOp.Delete(ctx, testZone, ctx.ID)
			},
		},
		Cleanup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			swOp := iaas.NewSwitchOp(caller)
			return swOp.Delete(ctx, testZone, switchID)
		},
	})
}

func TestBuilder_BuildWithRouter(t *testing.T) {
	var routerID, routerSwitchID, switchID, updSwitchID types.ID
	var addresses []string
	var testZone = testutil.TestZone()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			routerBuilder := &builder.Builder{
				Name:           testutil.ResourceName("vpc-router-builder"),
				NetworkMaskLen: 28,
				BandWidthMbps:  100,
				Client:         builder.NewAPIClient(caller),
			}

			created, err := routerBuilder.Build(ctx, testZone)
			if err != nil {
				return err
			}

			routerID = created.ID
			routerSwitchID = created.Switch.ID

			swOp := iaas.NewSwitchOp(caller)
			sw, err := swOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{
				Name: testutil.ResourceName("vpc-router-builder"),
			})
			if err != nil {
				return err
			}
			switchID = sw.ID

			updSwitch, err := swOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{
				Name: testutil.ResourceName("vpc-router-builder-upd"),
			})
			if err != nil {
				return err
			}
			updSwitchID = updSwitch.ID

			routerSwitch, err := swOp.Read(ctx, testZone, created.Switch.ID)
			if err != nil {
				return err
			}
			addresses = routerSwitch.Subnets[0].GetAssignedIPAddresses()
			return nil
		},
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				builder := &Builder{
					Zone:        testZone,
					Name:        testutil.ResourceName("vpc-router-builder"),
					Description: "description",
					Tags:        types.Tags{"tag1", "tag2"},
					PlanID:      types.VPCRouterPlans.Premium,
					NICSetting: &PremiumNICSetting{
						SwitchID:         routerSwitchID,
						VirtualIPAddress: addresses[0],
						IPAddresses:      []string{addresses[1], addresses[2]},
						IPAliases:        []string{addresses[3], addresses[4]},
					},
					AdditionalNICSettings: []AdditionalNICSettingHolder{
						&AdditionalPremiumNICSetting{
							SwitchID:         switchID,
							IPAddresses:      []string{"192.168.0.11", "192.168.0.12"},
							VirtualIPAddress: "192.168.0.1",
							NetworkMaskLen:   24,
							Index:            2,
						},
					},
					RouterSetting: &RouterSetting{
						VRID:                      1,
						InternetConnectionEnabled: types.StringTrue,
					},
					SetupOptions: getSetupOption(),
					Client:       iaas.NewVPCRouterOp(caller),
				}
				return builder.Build(ctx)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				vpcRouter := value.(*iaas.VPCRouter)
				found := false
				for _, iface := range vpcRouter.Interfaces {
					if iface.Index == 2 {
						found = true
						if err := testutil.AssertEqual(t, switchID, iface.SwitchID, "VPCRouter.Interfaces[index=2].SwitchID"); err != nil {
							return err
						}
					}
				}
				return testutil.AssertTrue(t, found, "VPCRouter.Interfaces[index=2]")
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				vpcRouterOp := iaas.NewVPCRouterOp(caller)
				return vpcRouterOp.Read(ctx, testZone, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				vpcRouter := value.(*iaas.VPCRouter)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, vpcRouter, "VPCRouter"),
					testutil.AssertNotNilFunc(t, vpcRouter.Settings, "VPCRouter.Settings"),
					testutil.AssertLenFunc(t, vpcRouter.Settings.Interfaces, 2, "VPCRouter.Settings.Interfaces"),
				)
			},
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					builder := &Builder{
						ID:   ctx.ID,
						Zone: testZone,

						Name:        testutil.ResourceName("vpc-router-builder"),
						Description: "description",
						Tags:        types.Tags{"tag1", "tag2"},
						PlanID:      types.VPCRouterPlans.Premium,
						NICSetting: &PremiumNICSetting{
							SwitchID:         routerSwitchID,
							VirtualIPAddress: addresses[0],
							IPAddresses:      []string{addresses[1], addresses[2]},
							IPAliases:        []string{addresses[3], addresses[4]},
						},
						AdditionalNICSettings: []AdditionalNICSettingHolder{
							&AdditionalPremiumNICSetting{
								SwitchID:         updSwitchID,
								VirtualIPAddress: "192.168.0.5",
								IPAddresses:      []string{"192.168.0.6", "192.168.0.7"},
								NetworkMaskLen:   28,
								Index:            3,
							},
						},
						RouterSetting: &RouterSetting{
							VRID:                      1,
							InternetConnectionEnabled: types.StringTrue,
						},
						SetupOptions: getSetupOption(),
						Client:       iaas.NewVPCRouterOp(caller),
					}
					return builder.Build(ctx)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
					vpcRouter := value.(*iaas.VPCRouter)
					found := false
					for _, iface := range vpcRouter.Interfaces {
						if iface.Index == 3 {
							found = true
							if err := testutil.AssertEqual(t, updSwitchID, iface.SwitchID, "VPCRouter.Interfaces[index=2].SwitchID"); err != nil {
								return err
							}
						}
					}
					if err := testutil.AssertTrue(t, found, "VPCRouter.Interfaces[index=2]"); err != nil {
						return err
					}

					found = false
					for _, nicSetting := range vpcRouter.Settings.Interfaces {
						if nicSetting.Index == 3 {
							found = true
							err := testutil.DoAsserts(
								testutil.AssertEqualFunc(t, "192.168.0.5", nicSetting.VirtualIPAddress, "VPCRouter.Settings.Interfaces.VirtualIPAddress"),
								testutil.AssertEqualFunc(t, []string{"192.168.0.6", "192.168.0.7"}, nicSetting.IPAddress, "VPCRouter.Settings.Interfaces.IPAddress"),
								testutil.AssertEqualFunc(t, 28, nicSetting.NetworkMaskLen, "VPCRouter.Settings.Interfaces.NetworkMaskLen"),
							)
							if err != nil {
								return err
							}
						}
					}
					return testutil.AssertTrue(t, found, "VPCRouter.Setting.Interfaces[index=2]")
				},
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				vpcRouterOp := iaas.NewVPCRouterOp(caller)
				if err := vpcRouterOp.Delete(ctx, testZone, ctx.ID); err != nil {
					return err
				}

				internetOp := iaas.NewInternetOp(caller)
				if err := internetOp.Delete(ctx, testZone, routerID); err != nil {
					return err
				}
				swOp := iaas.NewSwitchOp(caller)
				if err := swOp.Delete(ctx, testZone, switchID); err != nil {
					return err
				}
				return swOp.Delete(ctx, testZone, updSwitchID)
			},
		},
	})
}
