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
	"github.com/sacloud/iaas-api-go/helper/power"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestMobileGatewayOpCRUD(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURACLOUD_SIM_ICCID", "SAKURACLOUD_SIM_PASSCODE")(t)

	initMobileGatewayVariables()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel: true,

		SetupAPICallerFunc: singletonAPICaller,

		Create: &testutil.CRUDTestFunc{
			Func: testMobileGatewayCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createMobileGatewayExpected,
				IgnoreFields: ignoreMobileGatewayFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testMobileGatewayRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createMobileGatewayExpected,
				IgnoreFields: ignoreMobileGatewayFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testMobileGatewayUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateMobileGatewayExpected,
					IgnoreFields: ignoreMobileGatewayFields,
				}),
			},
			{
				Func: testMobileGatewayUpdateSettings,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateMobileGatewaySettingsExpected,
					IgnoreFields: ignoreMobileGatewayFields,
				}),
			},
			// shutdown(no check)
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					// shutdown
					if err := mgwOp.Shutdown(ctx, testZone, ctx.ID, &iaas.ShutdownOption{Force: true}); err != nil {
						return nil, err
					}

					waiter := iaas.WaiterForDown(func() (interface{}, error) {
						return mgwOp.Read(ctx, testZone, ctx.ID)
					})

					return waiter.WaitForState(ctx)
				},
				SkipExtractID: true,
			},
			// connect to switch
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					// prepare switch
					swOp := iaas.NewSwitchOp(caller)
					sw, err := swOp.Create(ctx, testZone, &iaas.SwitchCreateRequest{
						Name: testutil.ResourceName("switch-for-mobile-gateway"),
					})
					if err != nil {
						return nil, err
					}

					ctx.Values["mobile-gateway/switch"] = sw.ID

					// connect
					mgwOp := iaas.NewMobileGatewayOp(caller)
					if err := mgwOp.ConnectToSwitch(ctx, testZone, ctx.ID, sw.ID); err != nil {
						return nil, err
					}

					return mgwOp.Read(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					mgw := i.(*iaas.MobileGateway)
					return testutil.DoAsserts(
						testutil.AssertLenFunc(t, mgw.Interfaces, 2, "len(MobileGateway.Interfaces)"),
					)
				},
				SkipExtractID: true,
			},
			// set IPAddress to eth1
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					mgs, err := mgwOp.UpdateSettings(ctx, testZone, ctx.ID, &iaas.MobileGatewayUpdateSettingsRequest{
						InterfaceSettings: []*iaas.MobileGatewayInterfaceSetting{
							{
								IPAddress:      []string{"192.168.2.11"},
								NetworkMaskLen: 16,
								Index:          1,
							},
						},
					})
					if err != nil {
						return nil, err
					}
					if err := mgwOp.Config(ctx, testZone, ctx.ID); err != nil {
						return nil, err
					}
					return mgs, nil
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					mgw := i.(*iaas.MobileGateway)
					return testutil.DoAsserts(
						testutil.AssertNotNilFunc(t, mgw.InterfaceSettings, "MobileGateway.Settings.Interfaces"),
						testutil.AssertEqualFunc(t, 1, mgw.InterfaceSettings[0].Index, "MobileGateway.Settings.Interfaces.Index"),
						testutil.AssertEqualFunc(t, "192.168.2.11", mgw.InterfaceSettings[0].IPAddress[0], "MobileGateway.Settings.Interfaces.IPAddress"),
						testutil.AssertEqualFunc(t, 16, mgw.InterfaceSettings[0].NetworkMaskLen, "MobileGateway.Settings.Interfaces.NetworkMaskLen"),
					)
				},
				SkipExtractID: true,
			},

			// Get/Set DNS
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					if err := mgwOp.SetDNS(ctx, testZone, ctx.ID, &iaas.MobileGatewayDNSSetting{
						DNS1: "8.8.8.8",
						DNS2: "8.8.4.4",
					}); err != nil {
						return nil, err
					}
					return mgwOp.GetDNS(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					dns := i.(*iaas.MobileGatewayDNSSetting)
					return testutil.DoAsserts(
						testutil.AssertEqualFunc(t, "8.8.8.8", dns.DNS1, "DNS1"),
						testutil.AssertEqualFunc(t, "8.8.4.4", dns.DNS2, "DNS2"),
					)
				},
				SkipExtractID: true,
			},
			// Add/List SIM
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					simOp := iaas.NewSIMOp(caller)
					sim, err := simOp.Create(ctx, &iaas.SIMCreateRequest{
						Name:     testutil.ResourceName("switch-for-mobile-gateway"),
						ICCID:    iccid,
						PassCode: passcode,
					})
					if err != nil {
						return nil, err
					}

					mgwOp := iaas.NewMobileGatewayOp(caller)
					if err := mgwOp.AddSIM(ctx, testZone, ctx.ID, &iaas.MobileGatewayAddSIMRequest{
						SIMID: sim.ID.String(),
					}); err != nil {
						return nil, err
					}

					ctx.Values["mobile-gateway/sim"] = sim.ID
					return mgwOp.ListSIM(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					sims := i.(iaas.MobileGatewaySIMs)
					return testutil.DoAsserts(
						testutil.AssertLenFunc(t, sims, 1, "len(SIM)"),
					)
				},
				SkipExtractID: true,
			},
			// SIMOp: Assign IP
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					client := iaas.NewSIMOp(caller)
					simID := ctx.Values["mobile-gateway/sim"].(types.ID)
					if err := client.AssignIP(ctx, simID, &iaas.SIMAssignIPRequest{
						IP: "192.168.2.1",
					}); err != nil {
						return nil, err
					}
					return client.Status(ctx, simID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, v interface{}) error {
					simInfo := v.(*iaas.SIMInfo)
					return testutil.DoAsserts(
						testutil.AssertEqualFunc(t, "192.168.2.1", simInfo.IP, "SIMInfo.IP"),
					)
				},
				SkipExtractID: true,
			},
			// SIMOp: clear IP
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					client := iaas.NewSIMOp(caller)
					simID := ctx.Values["mobile-gateway/sim"].(types.ID)
					if err := client.ClearIP(ctx, simID); err != nil {
						return nil, err
					}
					return client.Status(ctx, simID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, v interface{}) error {
					simInfo := v.(*iaas.SIMInfo)
					return testutil.DoAsserts(
						testutil.AssertEmptyFunc(t, simInfo.IP, "SIMInfo.IP"),
					)
				},
				SkipExtractID: true,
			},

			// Get/Set SIMRoutes
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					if err := mgwOp.SetSIMRoutes(ctx, testZone, ctx.ID, []*iaas.MobileGatewaySIMRouteParam{
						{
							ResourceID: ctx.Values["mobile-gateway/sim"].(types.ID).String(),
							Prefix:     "192.168.3.0/24",
						},
					}); err != nil {
						return nil, err
					}
					return mgwOp.GetSIMRoutes(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					routes := i.(iaas.MobileGatewaySIMRoutes)
					simID := ctx.Values["mobile-gateway/sim"].(types.ID)
					return testutil.DoAsserts(
						testutil.AssertLenFunc(t, routes, 1, "len(SIMRoutes)"),
						testutil.AssertEqualFunc(t, "192.168.3.0/24", routes[0].Prefix, "SIMRoute.Prefix"),
						testutil.AssertEqualFunc(t, simID.String(), routes[0].ResourceID, "SIMRoute.ResourceID"),
					)
				},
				SkipExtractID: true,
			},
			// Delete SIMRoutes
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					if err := mgwOp.SetSIMRoutes(ctx, testZone, ctx.ID, []*iaas.MobileGatewaySIMRouteParam{}); err != nil {
						return nil, err
					}
					return mgwOp.GetSIMRoutes(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					routes := i.(iaas.MobileGatewaySIMRoutes)
					return testutil.DoAsserts(
						testutil.AssertLenFunc(t, routes, 0, "len(SIMRoutes)"),
					)
				},
				SkipExtractID: true,
			},

			// Get/Set TrafficConfig
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					if err := mgwOp.SetTrafficConfig(ctx, testZone, ctx.ID, &iaas.MobileGatewayTrafficControl{
						TrafficQuotaInMB:       10,
						BandWidthLimitInKbps:   20,
						EmailNotifyEnabled:     true,
						SlackNotifyEnabled:     true,
						SlackNotifyWebhooksURL: "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX",
						AutoTrafficShaping:     true,
					}); err != nil {
						return nil, err
					}
					return mgwOp.GetTrafficConfig(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					slackURL := "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"
					config := i.(*iaas.MobileGatewayTrafficControl)
					return testutil.DoAsserts(
						testutil.AssertEqualFunc(t, 10, config.TrafficQuotaInMB, "TrafficConfig.TrafficQuotaInMB"),
						testutil.AssertEqualFunc(t, 20, config.BandWidthLimitInKbps, "TrafficConfig.BandWidthLimitInKbps"),
						testutil.AssertEqualFunc(t, true, config.EmailNotifyEnabled, "TrafficConfig.EmailNotifyEnabled"),
						testutil.AssertEqualFunc(t, true, config.SlackNotifyEnabled, "TrafficConfig.SlackNotifyEnabled"),
						testutil.AssertEqualFunc(t, slackURL, config.SlackNotifyWebhooksURL, "TrafficConfig.SlackNotifyWebhooksURL"),
						testutil.AssertEqualFunc(t, true, config.AutoTrafficShaping, "TrafficConfig.AutoTrafficShaping"),
					)
				},
				SkipExtractID: true,
			},
			// Delete TrafficConfig(no check)
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					return nil, mgwOp.DeleteTrafficConfig(ctx, testZone, ctx.ID)
				},
				SkipExtractID: true,
			},

			// Get TrafficStatus
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					return mgwOp.TrafficStatus(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					status := i.(*iaas.MobileGatewayTrafficStatus)
					return testutil.DoAsserts(
						testutil.AssertNotNilFunc(t, status, "TrafficStatus"),
						testutil.AssertEqualFunc(t, types.StringNumber(0), status.UplinkBytes, "TrafficStatus.UplinkBytes"),
						testutil.AssertEqualFunc(t, types.StringNumber(0), status.DownlinkBytes, "TrafficStatus.DownlinkBytes"),
					)
				},
				SkipExtractID: true,
			},

			// Delete SIM
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					simID := ctx.Values["mobile-gateway/sim"].(types.ID)
					mgwOp := iaas.NewMobileGatewayOp(caller)
					if err := mgwOp.DeleteSIM(ctx, testZone, ctx.ID, simID); err != nil {
						return nil, err
					}

					simOp := iaas.NewSIMOp(caller)
					if err := simOp.Delete(ctx, simID); err != nil {
						return nil, err
					}

					return mgwOp.ListSIM(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					sims := i.(iaas.MobileGatewaySIMs)
					return testutil.DoAsserts(
						testutil.AssertLenFunc(t, sims, 0, "len(SIM)"),
					)
				},
				SkipExtractID: true,
			},
			// disconnect from switch
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					mgwOp := iaas.NewMobileGatewayOp(caller)
					if err := mgwOp.DisconnectFromSwitch(ctx, testZone, ctx.ID); err != nil {
						return nil, err
					}

					swID := ctx.Values["mobile-gateway/switch"].(types.ID)
					swOp := iaas.NewSwitchOp(caller)
					if err := swOp.Delete(ctx, testZone, swID); err != nil {
						return nil, err
					}

					return mgwOp.Read(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, i interface{}) error {
					mgw := i.(*iaas.MobileGateway)
					return testutil.DoAsserts(
						testutil.AssertLenFunc(t, mgw.Interfaces, 1, "len(MobileGateway.Interfaces)"),
					)
				},
				SkipExtractID: true,
			},
		},
		Shutdown: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			client := iaas.NewMobileGatewayOp(caller)
			return power.ShutdownMobileGateway(ctx, client, testZone, ctx.ID, true)
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testMobileGatewayDelete,
		},
	})
}

func initMobileGatewayVariables() {
	iccid = os.Getenv("SAKURACLOUD_SIM_ICCID")
	passcode = os.Getenv("SAKURACLOUD_SIM_PASSCODE")

	createMobileGatewayParam = &iaas.MobileGatewayCreateRequest{
		Name:                            testutil.ResourceName("mobile-gateway"),
		Description:                     "desc",
		Tags:                            []string{"tag1", "tag2"},
		InternetConnectionEnabled:       true,
		InterDeviceCommunicationEnabled: true,
	}
	createMobileGatewayExpected = &iaas.MobileGateway{
		Name:                            createMobileGatewayParam.Name,
		Description:                     createMobileGatewayParam.Description,
		Tags:                            createMobileGatewayParam.Tags,
		Availability:                    types.Availabilities.Available,
		InternetConnectionEnabled:       true,
		InterDeviceCommunicationEnabled: true,
	}
	updateMobileGatewayParam = &iaas.MobileGatewayUpdateRequest{
		Name:                            testutil.ResourceName("mobile-gateway-upd"),
		Description:                     "desc-upd",
		Tags:                            []string{"tag1-upd", "tag2-upd"},
		InternetConnectionEnabled:       false,
		InterDeviceCommunicationEnabled: false,
	}
	updateMobileGatewayExpected = &iaas.MobileGateway{
		Name:                            updateMobileGatewayParam.Name,
		Description:                     updateMobileGatewayParam.Description,
		Tags:                            updateMobileGatewayParam.Tags,
		Availability:                    types.Availabilities.Available,
		InternetConnectionEnabled:       false,
		InterDeviceCommunicationEnabled: false,
	}
	updateMobileGatewaySettingsParam = &iaas.MobileGatewayUpdateSettingsRequest{
		InternetConnectionEnabled:       true,
		InterDeviceCommunicationEnabled: true,
	}
	updateMobileGatewaySettingsExpected = &iaas.MobileGateway{
		Name:                            updateMobileGatewayParam.Name,
		Description:                     updateMobileGatewayParam.Description,
		Tags:                            updateMobileGatewayParam.Tags,
		Availability:                    types.Availabilities.Available,
		InternetConnectionEnabled:       true,
		InterDeviceCommunicationEnabled: true,
	}
}

var (
	ignoreMobileGatewayFields = []string{
		"ID",
		"Class",
		"IconID",
		"CreatedAt",
		"Availability",
		"InstanceHostName",
		"InstanceHostInfoURL",
		"InstanceStatus",
		"InstanceStatusChangedAt",
		"Interfaces",
		"ZoneID",
		"SettingsHash",
	}
	iccid                               string
	passcode                            string
	createMobileGatewayParam            *iaas.MobileGatewayCreateRequest
	createMobileGatewayExpected         *iaas.MobileGateway
	updateMobileGatewayParam            *iaas.MobileGatewayUpdateRequest
	updateMobileGatewayExpected         *iaas.MobileGateway
	updateMobileGatewaySettingsParam    *iaas.MobileGatewayUpdateSettingsRequest
	updateMobileGatewaySettingsExpected *iaas.MobileGateway
)

func testMobileGatewayCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewMobileGatewayOp(caller)
	v, err := client.Create(ctx, testZone, createMobileGatewayParam)
	if err != nil {
		return nil, err
	}
	value, err := iaas.WaiterForReady(func() (interface{}, error) {
		return client.Read(ctx, testZone, v.ID)
	}).WaitForState(ctx)
	if err != nil {
		return nil, err
	}
	if err := client.Boot(ctx, testZone, v.ID); err != nil {
		return nil, err
	}
	return value, nil
}

func testMobileGatewayRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewMobileGatewayOp(caller)
	return client.Read(ctx, testZone, ctx.ID)
}

func testMobileGatewayUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewMobileGatewayOp(caller)
	return client.Update(ctx, testZone, ctx.ID, updateMobileGatewayParam)
}

func testMobileGatewayUpdateSettings(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewMobileGatewayOp(caller)
	return client.UpdateSettings(ctx, testZone, ctx.ID, updateMobileGatewaySettingsParam)
}

func testMobileGatewayDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewMobileGatewayOp(caller)
	return client.Delete(ctx, testZone, ctx.ID)
}
