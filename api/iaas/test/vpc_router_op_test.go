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
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestVPCRouterOp_CRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: testVPCRouterCreate(createVPCRouterParam),
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createVPCRouterExpected,
				IgnoreFields: ignoreVPCRouterFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testVPCRouterRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createVPCRouterExpected,
				IgnoreFields: ignoreVPCRouterFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testVPCRouterUpdate(updateVPCRouterParam),
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateVPCRouterExpected,
					IgnoreFields: ignoreVPCRouterFields,
				}),
			},
		},

		Shutdown: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			client := iaas.NewVPCRouterOp(caller)
			return power.ShutdownVPCRouter(ctx, client, testZone, ctx.ID, true)
		},

		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testVPCRouterDelete,
		},
	})
}

var (
	ignoreVPCRouterFields = []string{
		"ID",
		"Availability",
		"Class",
		"CreatedAt",
		"SettingsHash",
		"Settings",
		"InstanceHostName",
		"InstanceHostInfoURL",
		"InstanceStatus",
		"InstanceStatusChangedAt",
		"Interfaces",
		"ZoneID",
	}

	createVPCRouterParam = &iaas.VPCRouterCreateRequest{
		PlanID: types.VPCRouterPlans.Standard,
		Switch: &iaas.ApplianceConnectedSwitch{
			Scope: types.Scopes.Shared,
		},
		Name:        testutil.ResourceName("vpc-router"),
		Description: "desc",
		Tags:        []string{"tag1", "tag2"},
		Settings:    &iaas.VPCRouterSetting{},
	}
	createVPCRouterExpected = &iaas.VPCRouter{
		Class:          "vpcrouter",
		Name:           createVPCRouterParam.Name,
		Description:    createVPCRouterParam.Description,
		Tags:           createVPCRouterParam.Tags,
		Availability:   types.Availabilities.Available,
		InstanceStatus: types.ServerInstanceStatuses.Up,
		PlanID:         createVPCRouterParam.PlanID,
		Version:        2,
		Settings:       createVPCRouterParam.Settings,
	}
	updateVPCRouterParam = &iaas.VPCRouterUpdateRequest{
		Name:        testutil.ResourceName("vpc-router-upd"),
		Tags:        []string{"tag1-upd", "tag2-upd"},
		Description: "desc-upd",
	}
	updateVPCRouterExpected = &iaas.VPCRouter{
		Class:          "vpcrouter",
		Name:           updateVPCRouterParam.Name,
		Description:    updateVPCRouterParam.Description,
		Tags:           updateVPCRouterParam.Tags,
		Availability:   types.Availabilities.Available,
		InstanceStatus: types.ServerInstanceStatuses.Up,
		PlanID:         createVPCRouterParam.PlanID,
		Version:        2,
	}
)

func testVPCRouterCreate(createParam *iaas.VPCRouterCreateRequest) func(*testutil.CRUDTestContext, iaas.APICaller) (interface{}, error) {
	return func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
		client := iaas.NewVPCRouterOp(caller)
		vpcRouter, err := client.Create(ctx, testZone, createParam)
		if err != nil {
			return nil, err
		}

		n, err := iaas.WaiterForReady(func() (interface{}, error) {
			return client.Read(ctx, testZone, vpcRouter.ID)
		}).WaitForState(ctx)
		if err != nil {
			return nil, err
		}

		if err := client.Boot(ctx, testZone, vpcRouter.ID); err != nil {
			return nil, err
		}

		return n.(*iaas.VPCRouter), nil
	}
}

func testVPCRouterRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewVPCRouterOp(caller)
	return client.Read(ctx, testZone, ctx.ID)
}

func testVPCRouterUpdate(updateParam *iaas.VPCRouterUpdateRequest) func(*testutil.CRUDTestContext, iaas.APICaller) (interface{}, error) {
	return func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
		client := iaas.NewVPCRouterOp(caller)
		return client.Update(ctx, testZone, ctx.ID, updateParam)
	}
}

func testVPCRouterDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewVPCRouterOp(caller)
	return client.Delete(ctx, testZone, ctx.ID)
}

var fakeWireGuardPublicKey = `fqxOlS2X0Jtg4P9zVf8D3BAUtJmrp+z2mjzUmgxxxxx=`

func TestVPCRouterOp_WithRouterCRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		SetupAPICallerFunc: singletonAPICaller,

		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			routerOp := iaas.NewInternetOp(caller)
			created, err := routerOp.Create(ctx, testZone, &iaas.InternetCreateRequest{
				Name:           testutil.ResourceName("internet-for-vpc-router"),
				BandWidthMbps:  100,
				NetworkMaskLen: 28,
			})
			if err != nil {
				return err
			}

			ctx.Values["vpcrouter/internet"] = created.ID
			max := 30
			for {
				if max == 0 {
					break
				}
				_, err := routerOp.Read(ctx, testZone, created.ID)
				if err != nil || iaas.IsNotFoundError(err) {
					max--
					time.Sleep(3 * time.Second)
					continue
				}
				break
			}

			swOp := iaas.NewSwitchOp(caller)
			sw, err := swOp.Read(ctx, testZone, created.Switch.ID)
			if err != nil {
				return err
			}

			ipaddresses := sw.Subnets[0].GetAssignedIPAddresses()
			p := withRouterCreateVPCRouterParam
			p.Switch = &iaas.ApplianceConnectedSwitch{
				ID: sw.ID,
			}
			p.IPAddresses = []string{ipaddresses[1], ipaddresses[2]}
			p.Settings = &iaas.VPCRouterSetting{
				VRID:                      100,
				InternetConnectionEnabled: true,
				Interfaces: []*iaas.VPCRouterInterfaceSetting{
					{
						VirtualIPAddress: ipaddresses[0],
						IPAddress:        []string{ipaddresses[1], ipaddresses[2]},
						IPAliases:        []string{ipaddresses[3]},
						NetworkMaskLen:   sw.Subnets[0].NetworkMaskLen,
					},
				},
			}

			withRouterCreateVPCRouterExpected.Settings = p.Settings
			return nil
		},
		Create: &testutil.CRUDTestFunc{
			Func: testVPCRouterCreate(withRouterCreateVPCRouterParam),
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  withRouterCreateVPCRouterExpected,
				IgnoreFields: ignoreVPCRouterFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testVPCRouterRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  withRouterCreateVPCRouterExpected,
				IgnoreFields: ignoreVPCRouterFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					if isAccTest() {
						// 起動直後だとシャットダウンできない場合があるため20秒ほど待つ
						time.Sleep(20 * time.Second)
					}

					vpcOp := iaas.NewVPCRouterOp(caller)
					// shutdown
					if err := vpcOp.Shutdown(ctx, testZone, ctx.ID, &iaas.ShutdownOption{Force: true}); err != nil {
						return nil, err
					}
					_, err := iaas.WaiterForDown(func() (interface{}, error) {
						return vpcOp.Read(ctx, testZone, ctx.ID)
					}).WaitForState(ctx)
					if err != nil {
						return nil, err
					}

					swOp := iaas.NewSwitchOp(caller)
					sw, err := swOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{
						Name: testutil.ResourceName("switch-for-vpc-router"),
					})
					if err != nil {
						return nil, err
					}
					ctx.Values["vpcrouter/switch"] = sw.ID

					// connect to switch
					if err := vpcOp.ConnectToSwitch(ctx, testZone, ctx.ID, 2, sw.ID); err != nil {
						return nil, err
					}

					// setup update param
					p := withRouterUpdateVPCRouterParam
					p.Settings = &iaas.VPCRouterSetting{
						VRID:                      10,
						SyslogHost:                "192.168.2.199",
						InternetConnectionEnabled: true,
						Interfaces: []*iaas.VPCRouterInterfaceSetting{
							withRouterCreateVPCRouterParam.Settings.Interfaces[0],
							{
								VirtualIPAddress: "192.168.2.1",
								IPAddress:        []string{"192.168.2.11", "192.168.2.12"},
								NetworkMaskLen:   24,
								Index:            2,
							},
						},
						StaticNAT: []*iaas.VPCRouterStaticNAT{
							{
								GlobalAddress:  withRouterCreateVPCRouterParam.Settings.Interfaces[0].IPAliases[0],
								PrivateAddress: "192.168.2.1",
							},
						},
						PortForwarding: []*iaas.VPCRouterPortForwarding{
							{
								Protocol:       types.VPCRouterPortForwardingProtocols.TCP,
								GlobalPort:     22,
								PrivateAddress: "192.168.2.2",
								PrivatePort:    10022,
								Description:    "port forwarding",
							},
						},
						DHCPServer: []*iaas.VPCRouterDHCPServer{
							{
								Interface:  "eth2",
								RangeStart: "192.168.2.51",
								RangeStop:  "192.168.2.60",
							},
						},
						DHCPStaticMapping: []*iaas.VPCRouterDHCPStaticMapping{
							{
								MACAddress: "aa:bb:cc:dd:ee:ff",
								IPAddress:  "192.168.2.21",
							},
						},
						DNSForwarding: &iaas.VPCRouterDNSForwarding{
							Interface:  "eth2",
							DNSServers: []string{"133.242.0.3", "133.242.0.4"},
						},
						PPTPServer: &iaas.VPCRouterPPTPServer{
							RangeStart: "192.168.2.61",
							RangeStop:  "192.168.2.70",
						},
						PPTPServerEnabled: true,
						L2TPIPsecServer: &iaas.VPCRouterL2TPIPsecServer{
							RangeStart:      "192.168.2.71",
							RangeStop:       "192.168.2.80",
							PreSharedSecret: "presharedsecret",
						},
						L2TPIPsecServerEnabled: true,
						WireGuard: &iaas.VPCRouterWireGuard{
							IPAddress: "192.168.3.1/24",
							Peers: []*iaas.VPCRouterWireGuardPeer{
								{
									Name:      "foobar",
									IPAddress: "192.168.3.11",
									PublicKey: fakeWireGuardPublicKey,
								},
							},
						},
						WireGuardEnabled: true,
						RemoteAccessUsers: []*iaas.VPCRouterRemoteAccessUser{
							{
								UserName: "user1",
								Password: "password1",
							},
						},
						SiteToSiteIPsecVPN: &iaas.VPCRouterSiteToSiteIPsecVPN{
							Config: []*iaas.VPCRouterSiteToSiteIPsecVPNConfig{
								{
									Peer:            "1.2.3.4",
									PreSharedSecret: "presharedsecret",
									RemoteID:        "1.2.3.4",
									Routes:          []string{"10.0.0.0/24"},
									LocalPrefix:     []string{"192.168.2.0/24"},
								},
							},
							IKE: &iaas.VPCRouterSiteToSiteIPsecVPNIKE{
								Lifetime: 28801,
								DPD: &iaas.VPCRouterSiteToSiteIPsecVPNIKEDPD{
									Interval: 16,
									Timeout:  31,
								},
							},
							ESP: &iaas.VPCRouterSiteToSiteIPsecVPNESP{
								Lifetime: 1801,
							},
							EncryptionAlgo: types.VPCRouterSiteToSiteVPNEncryptionAlgoAES256,
							HashAlgo:       types.VPCRouterSiteToSiteVPNHashAlgoSHA256,
							DHGroup:        types.VPCRouterSiteToSiteVPNDHGroupModp2048,
						},
						StaticRoute: []*iaas.VPCRouterStaticRoute{
							{
								Prefix:  "172.16.0.0/16",
								NextHop: "192.168.2.11",
							},
						},
						ScheduledMaintenance: &iaas.VPCRouterScheduledMaintenance{
							DayOfWeek: 1,
							Hour:      2,
						},
					}

					withRouterUpdateVPCRouterExpected.Settings = p.Settings
					return testVPCRouterUpdate(withRouterUpdateVPCRouterParam)(ctx, caller)
				},
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  withRouterUpdateVPCRouterExpected,
					IgnoreFields: ignoreVPCRouterFields,
				}),
			},
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					// setup update param
					p := withRouterUpdateVPCRouterToMinParam
					p.Settings = &iaas.VPCRouterSetting{
						VRID:                      10,
						InternetConnectionEnabled: false,
						Interfaces: []*iaas.VPCRouterInterfaceSetting{
							withRouterCreateVPCRouterParam.Settings.Interfaces[0],
						},
					}

					withRouterUpdateVPCRouterToMinExpected.Settings = p.Settings
					return testVPCRouterUpdate(withRouterUpdateVPCRouterToMinParam)(ctx, caller)
				},
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  withRouterUpdateVPCRouterToMinExpected,
					IgnoreFields: ignoreVPCRouterFields,
				}),
			},
		},

		Shutdown: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			client := iaas.NewVPCRouterOp(caller)
			return power.ShutdownVPCRouter(ctx, client, testZone, ctx.ID, true)
		},

		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testVPCRouterDelete,
		},

		Cleanup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			routerOp := iaas.NewInternetOp(caller)
			routerID, ok := ctx.Values["vpcrouter/internet"]
			if ok {
				if err := routerOp.Delete(ctx, testZone, routerID.(types.ID)); err != nil {
					return err
				}
			}

			swOp := iaas.NewSwitchOp(caller)
			switchID, ok := ctx.Values["vpcrouter/switch"]
			if ok {
				if err := swOp.Delete(ctx, testZone, switchID.(types.ID)); err != nil {
					return err
				}
			}
			return nil
		},
	})
}

var (
	withRouterCreateVPCRouterParam = &iaas.VPCRouterCreateRequest{
		PlanID:      types.VPCRouterPlans.Premium,
		Name:        testutil.ResourceName("vpc-router"),
		Description: "desc",
		Tags:        []string{"tag1", "tag2"},
	}
	withRouterCreateVPCRouterExpected = &iaas.VPCRouter{
		Class:          "vpcrouter",
		Name:           withRouterCreateVPCRouterParam.Name,
		Description:    withRouterCreateVPCRouterParam.Description,
		Tags:           withRouterCreateVPCRouterParam.Tags,
		Availability:   types.Availabilities.Available,
		InstanceStatus: types.ServerInstanceStatuses.Up,
		PlanID:         withRouterCreateVPCRouterParam.PlanID,
		Version:        2,
		Settings:       withRouterCreateVPCRouterParam.Settings,
	}
	withRouterUpdateVPCRouterParam = &iaas.VPCRouterUpdateRequest{
		Name:        testutil.ResourceName("vpc-router-upd"),
		Tags:        []string{"tag1-upd", "tag2-upd"},
		Description: "desc-upd",
		IconID:      testIconID,
	}
	withRouterUpdateVPCRouterExpected = &iaas.VPCRouter{
		Class:          "vpcrouter",
		Name:           withRouterUpdateVPCRouterParam.Name,
		Description:    withRouterUpdateVPCRouterParam.Description,
		Tags:           withRouterUpdateVPCRouterParam.Tags,
		Availability:   types.Availabilities.Available,
		InstanceStatus: types.ServerInstanceStatuses.Up,
		PlanID:         withRouterCreateVPCRouterParam.PlanID,
		Version:        2,
		Settings:       withRouterUpdateVPCRouterParam.Settings,
		IconID:         testIconID,
	}
	withRouterUpdateVPCRouterToMinParam = &iaas.VPCRouterUpdateRequest{
		Name: testutil.ResourceName("vpc-router-to-min"),
	}
	withRouterUpdateVPCRouterToMinExpected = &iaas.VPCRouter{
		Class:          "vpcrouter",
		Name:           withRouterUpdateVPCRouterToMinParam.Name,
		Availability:   types.Availabilities.Available,
		InstanceStatus: types.ServerInstanceStatuses.Up,
		PlanID:         withRouterCreateVPCRouterParam.PlanID,
		Version:        2,
		Settings:       withRouterUpdateVPCRouterToMinParam.Settings,
	}
)
