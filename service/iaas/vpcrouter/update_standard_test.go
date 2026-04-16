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
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/setup"
	vpcRouterBuilder "github.com/sacloud/iaas-service-go/vpcrouter/builder"
	"github.com/sacloud/packages-go/pointer"
	"github.com/stretchr/testify/require"
)

func TestVPCRouterService_convertUpdateStandardRequest(t *testing.T) {
	ctx := context.Background()
	caller := testutil.SingletonAPICaller()
	zone := testutil.TestZone()
	name := testutil.ResourceName("vpc-router-service-update")

	// setup
	swOp := iaas.NewSwitchOp(caller)
	additionalSwitch, err := swOp.Create(ctx, zone, &iaas.SwitchCreateRequest{Name: name})
	if err != nil {
		t.Fatal(err)
	}

	createReq := &CreateStandardRequest{
		Zone:        zone,
		Name:        name,
		Description: "desc",
		Tags:        types.Tags{"tag1", "tag2"},
		AdditionalNICSettings: []*vpcRouterBuilder.AdditionalStandardNICSetting{
			{
				SwitchID:       additionalSwitch.ID,
				IPAddress:      "192.168.0.101",
				NetworkMaskLen: 24,
				Index:          1,
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
		iaas.NewVPCRouterOp(caller).Delete(ctx, zone, vpcRouter.ID) //nolint
		swOp.Delete(ctx, zone, additionalSwitch.ID)                 //nolint
	}()

	// test
	cases := []struct {
		in     *UpdateStandardRequest
		expect *ApplyRequest
	}{
		{
			in: &UpdateStandardRequest{
				ID:     vpcRouter.ID,
				Zone:   zone,
				Name:   pointer.NewString(name + "-upd"),
				NoWait: true,
			},
			expect: &ApplyRequest{
				ID:          vpcRouter.ID,
				Zone:        zone,
				Name:        name + "-upd",
				Description: "desc",
				Tags:        types.Tags{"tag1", "tag2"},
				PlanID:      types.VPCRouterPlans.Standard,
				NICSetting:  &vpcRouterBuilder.StandardNICSetting{},
				AdditionalNICSettings: []vpcRouterBuilder.AdditionalNICSettingHolder{
					&vpcRouterBuilder.AdditionalStandardNICSetting{
						SwitchID:       additionalSwitch.ID,
						IPAddress:      "192.168.0.101",
						NetworkMaskLen: 24,
						Index:          1,
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
