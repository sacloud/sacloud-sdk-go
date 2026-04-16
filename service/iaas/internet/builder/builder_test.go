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

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/cleanup"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestBuilder_Build(t *testing.T) {
	var testZone = testutil.TestZone()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Parallel:           true,
		IgnoreStartupWait:  true,
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				builder := &Builder{
					Name:           testutil.ResourceName("internet-builder"),
					Description:    "description",
					Tags:           types.Tags{"tag1", "tag2"},
					NetworkMaskLen: 28,
					BandWidthMbps:  100,
					EnableIPv6:     true,
					Client:         NewAPIClient(caller),
				}
				return builder.Build(ctx, testZone)
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				return iaas.NewInternetOp(caller).Read(ctx, testZone, ctx.ID)
			},
			CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
				internet := value.(*iaas.Internet)
				return testutil.DoAsserts(
					testutil.AssertNotNilFunc(t, internet, "Internet"),
					testutil.AssertEqualFunc(t, 28, internet.NetworkMaskLen, "Internet.NetworkMaskLen"),
					testutil.AssertEqualFunc(t, 100, internet.BandWidthMbps, "Internet.BandWidthMbps"),
					testutil.AssertTrueFunc(t, len(internet.Switch.IPv6Nets) > 0, "Internet.Switch.IPv6Nets"),
				)
			},
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					builder := &Builder{
						Name:           testutil.ResourceName("internet-builder"),
						Description:    "description",
						Tags:           types.Tags{"tag1", "tag2"},
						NetworkMaskLen: 28,
						BandWidthMbps:  500,
						EnableIPv6:     false,
						Client:         NewAPIClient(caller),
					}
					return builder.Update(ctx, testZone, ctx.ID)
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
					internet := value.(*iaas.Internet)
					return testutil.DoAsserts(
						testutil.AssertNotNilFunc(t, internet, "Internet"),
						testutil.AssertEqualFunc(t, 28, internet.NetworkMaskLen, "Internet.NetworkMaskLen"),
						testutil.AssertEqualFunc(t, 500, internet.BandWidthMbps, "Internet.BandWidthMbps"),
						testutil.AssertTrueFunc(t, len(internet.Switch.IPv6Nets) == 0, "Internet.Switch.IPv6Nets"),
					)
				},
			},
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					internetOp := iaas.NewInternetOp(caller)
					swOp := iaas.NewSwitchOp(caller)

					internet, err := internetOp.Read(ctx, testZone, ctx.ID)
					if err != nil {
						return nil, err
					}
					sw, err := swOp.Read(ctx, testZone, internet.Switch.ID)
					if err != nil {
						return nil, err
					}

					_, err = internetOp.AddSubnet(ctx, testZone, ctx.ID, &iaas.InternetAddSubnetRequest{
						NetworkMaskLen: 28,
						NextHop:        sw.Subnets[0].AssignedIPAddressMin,
					})
					return nil, err
				},
				CheckFunc: func(t testutil.TestT, ctx *testutil.CRUDTestContext, value interface{}) error {
					internet := value.(*iaas.Internet)
					return testutil.DoAsserts(
						testutil.AssertNotNilFunc(t, internet, "Internet"),
						testutil.AssertTrueFunc(t, len(internet.Switch.Subnets) == 2, "Internet.Switch.Subnets"),
					)
				},
				SkipExtractID: true,
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
				return cleanup.DeleteInternet(ctx, iaas.NewInternetOp(caller), testZone, ctx.ID)
			},
		},
	})
}
