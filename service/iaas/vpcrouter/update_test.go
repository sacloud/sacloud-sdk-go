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

package vpcrouter

import (
	"context"
	"testing"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	internetService "github.com/sacloud/iaas-service-go/internet"
	"github.com/sacloud/iaas-service-go/setup"
	vpcRouterBuilder "github.com/sacloud/iaas-service-go/vpcrouter/builder"
	"github.com/sacloud/packages-go/pointer"
	"github.com/stretchr/testify/require"
)

func TestVPCRouterService_convertUpdateRequest(t *testing.T) {
	ctx := context.Background()
	caller := testutil.SingletonAPICaller()
	zone := testutil.TestZone()
	name := testutil.ResourceName("vpc-router-service-update")

	// setup
	internetSvc := internetService.New(caller)
	router, err := internetSvc.CreateWithContext(ctx, &internetService.CreateRequest{
		Name:           name,
		Zone:           zone,
		NetworkMaskLen: 28,
		BandWidthMbps:  100,
	})
	if err != nil {
		t.Fatal(err)
	}

	swOp := iaas.NewSwitchOp(caller)
	sw, err := swOp.Read(ctx, zone, router.Switch.ID)
	if err != nil {
		t.Fatal(err)
	}

	additionalSwitch, err := swOp.Create(ctx, zone, &iaas.SwitchCreateRequest{Name: name})
	if err != nil {
		t.Fatal(err)
	}

	createReq := &CreateRequest{
		Zone:        zone,
		Name:        name,
		Description: "desc",
		Tags:        types.Tags{"tag1", "tag2"},
		PlanID:      types.VPCRouterPlans.Premium,
		NICSetting: &vpcRouterBuilder.PremiumNICSetting{
			SwitchID:         sw.ID,
			IPAddresses:      []string{sw.Subnets[0].GetAssignedIPAddresses()[0], sw.Subnets[0].GetAssignedIPAddresses()[1]},
			VirtualIPAddress: sw.Subnets[0].GetAssignedIPAddresses()[2],
			IPAliases:        []string{sw.Subnets[0].GetAssignedIPAddresses()[3]},
		},
		AdditionalNICSettings: []*vpcRouterBuilder.AdditionalPremiumNICSetting{
			{
				SwitchID:         additionalSwitch.ID,
				IPAddresses:      []string{"192.168.0.101", "192.168.0.102"},
				VirtualIPAddress: "192.168.0.11",
				NetworkMaskLen:   24,
				Index:            1,
			},
		},
		RouterSetting: &RouterSetting{
			VRID:                      1,
			InternetConnectionEnabled: true,
			L2TPIPsecServer: &iaas.VPCRouterL2TPIPsecServer{
				RangeStart:      "192.168.0.250",
				RangeStop:       "192.168.0.254",
				PreSharedSecret: "presharedsecret",
			},
			RemoteAccessUsers: []*iaas.VPCRouterRemoteAccessUser{
				{
					UserName: "username",
					Password: "password",
				},
			},
		},
		NoWait: false,
	}

	builder := createReq.ApplyRequest().Builder(caller)
	if !testutil.IsAccTest() {
		builder.SetupOptions = &setup.Options{
			NICUpdateWaitDuration:     time.Millisecond,
			ProvisioningRetryInterval: time.Millisecond,
			DeleteRetryInterval:       time.Millisecond,
			PollingInterval:           time.Millisecond,
		}
	}
	vpcRouter, err := builder.Build(ctx)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		iaas.NewVPCRouterOp(caller).Delete(ctx, zone, vpcRouter.ID)              //nolint
		swOp.Delete(ctx, zone, additionalSwitch.ID)                              //nolint
		cleanup.DeleteInternet(ctx, iaas.NewInternetOp(caller), zone, router.ID) //nolint
	}()

	// test
	cases := []struct {
		in     *UpdateRequest
		expect *ApplyRequest
	}{
		{
			in: &UpdateRequest{
				ID:     vpcRouter.ID,
				Zone:   zone,
				Name:   pointer.NewString(name + "-upd"),
				NoWait: true,

				NICSetting: &PremiumNICSettingUpdate{
					IPAddresses:      &[]string{sw.Subnets[0].GetAssignedIPAddresses()[1], sw.Subnets[0].GetAssignedIPAddresses()[2]},
					VirtualIPAddress: &sw.Subnets[0].GetAssignedIPAddresses()[3],
					IPAliases:        &[]string{sw.Subnets[0].GetAssignedIPAddresses()[4]},
				},
				AdditionalNICSettings: &[]*AdditionalPremiumNICSettingUpdate{
					{
						SwitchID:         &additionalSwitch.ID,
						IPAddresses:      &[]string{"192.168.0.101", "192.168.0.102"},
						VirtualIPAddress: pointer.NewString("192.168.0.11"),
						NetworkMaskLen:   pointer.NewInt(24),
						Index:            2,
					},
				},
			},
			expect: &ApplyRequest{
				ID:          vpcRouter.ID,
				Zone:        zone,
				Name:        name + "-upd",
				Description: "desc",
				Tags:        types.Tags{"tag1", "tag2"},
				PlanID:      types.VPCRouterPlans.Premium,
				NICSetting: &vpcRouterBuilder.PremiumNICSetting{
					SwitchID:         sw.ID,
					IPAddresses:      []string{sw.Subnets[0].GetAssignedIPAddresses()[1], sw.Subnets[0].GetAssignedIPAddresses()[2]},
					VirtualIPAddress: sw.Subnets[0].GetAssignedIPAddresses()[3],
					IPAliases:        []string{sw.Subnets[0].GetAssignedIPAddresses()[4]},
				},
				AdditionalNICSettings: []vpcRouterBuilder.AdditionalNICSettingHolder{
					&vpcRouterBuilder.AdditionalPremiumNICSetting{
						SwitchID:         additionalSwitch.ID,
						IPAddresses:      []string{"192.168.0.101", "192.168.0.102"},
						VirtualIPAddress: "192.168.0.11",
						NetworkMaskLen:   24,
						Index:            2,
					},
				},
				RouterSetting: &RouterSetting{
					VRID:                      1,
					InternetConnectionEnabled: true,
					L2TPIPsecServer: &iaas.VPCRouterL2TPIPsecServer{
						RangeStart:      "192.168.0.250",
						RangeStop:       "192.168.0.254",
						PreSharedSecret: "presharedsecret",
					},
					RemoteAccessUsers: []*iaas.VPCRouterRemoteAccessUser{
						{
							UserName: "username",
							Password: "password",
						},
					},
				},
				NoWait: true,
			},
		},
	}

	for _, tc := range cases {
		req, err := tc.in.ApplyRequest(ctx, caller)
		require.NoError(t, err)
		require.EqualValues(t, tc.expect, req)
	}
}
