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
	"errors"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/assert"
)

func TestInternetOp_CRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel: true,

		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: testInternetCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createInternetExpected,
				IgnoreFields: ignoreInternetFields,
			}),
		},
		Read: &testutil.CRUDTestFunc{
			Func: testInternetRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createInternetExpected,
				IgnoreFields: ignoreInternetFields,
			}),
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testInternetUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateInternetExpected,
					IgnoreFields: ignoreInternetFields,
				}),
			},
			{
				Func: testInternetUpdateToMin,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateInternetToMinExpected,
					IgnoreFields: ignoreInternetFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testInternetDelete,
		},
	})
}

var (
	ignoreInternetFields = []string{
		"ID",
		"CreatedAt",
		"Switch",
	}
	createInternetParam = &iaas.InternetCreateRequest{
		Name:           testutil.ResourceName("internet"),
		Description:    "desc",
		Tags:           []string{"tag1", "tag2"},
		NetworkMaskLen: 28,
		BandWidthMbps:  100,
	}
	createInternetExpected = &iaas.Internet{
		Name:           createInternetParam.Name,
		Description:    createInternetParam.Description,
		Tags:           createInternetParam.Tags,
		NetworkMaskLen: createInternetParam.NetworkMaskLen,
		BandWidthMbps:  createInternetParam.BandWidthMbps,
	}
	updateInternetParam = &iaas.InternetUpdateRequest{
		Name:        testutil.ResourceName("internet-upd"),
		Tags:        []string{"tag1-upd", "tag2-upd"},
		Description: "desc-upd",
		IconID:      testIconID,
	}
	updateInternetExpected = &iaas.Internet{
		Name:           updateInternetParam.Name,
		Description:    updateInternetParam.Description,
		Tags:           updateInternetParam.Tags,
		NetworkMaskLen: createInternetParam.NetworkMaskLen,
		BandWidthMbps:  createInternetParam.BandWidthMbps,
		IconID:         testIconID,
	}
	updateInternetToMinParam = &iaas.InternetUpdateRequest{
		Name: testutil.ResourceName("internet-to-min"),
	}
	updateInternetToMinExpected = &iaas.Internet{
		Name:           updateInternetToMinParam.Name,
		NetworkMaskLen: createInternetParam.NetworkMaskLen,
		BandWidthMbps:  createInternetParam.BandWidthMbps,
	}
)

func testInternetCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewInternetOp(caller)
	return client.Create(ctx, testZone, createInternetParam)
}

func testInternetRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewInternetOp(caller)
	return client.Read(ctx, testZone, ctx.ID)
}

func testInternetUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewInternetOp(caller)
	return client.Update(ctx, testZone, ctx.ID, updateInternetParam)
}

func testInternetUpdateToMin(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewInternetOp(caller)
	return client.Update(ctx, testZone, ctx.ID, updateInternetToMinParam)
}

func testInternetDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewInternetOp(caller)
	return client.Delete(ctx, testZone, ctx.ID)
}

func TestInternetOp_Subnet(t *testing.T) {
	client := iaas.NewInternetOp(singletonAPICaller())
	var minIP, maxIP string
	var subnetID types.ID

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		IgnoreStartupWait:  true,
		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
				var internet *iaas.Internet
				internet, err := client.Create(ctx, testZone, createInternetParam)
				if err != nil {
					return nil, err
				}
				waiter := iaas.WaiterForApplianceUp(func() (interface{}, error) {
					return client.Read(ctx, testZone, internet.ID)
				}, 100)
				if _, err := waiter.WaitForState(ctx); err != nil {
					t.Error("WaitForUp is failed: ", err)
					return nil, err
				}

				internet, err = client.Read(ctx, testZone, internet.ID)
				if err != nil {
					return nil, err
				}

				swOp := iaas.NewSwitchOp(singletonAPICaller())
				sw, err := swOp.Read(ctx, testZone, internet.Switch.ID)
				if err != nil {
					return nil, err
				}
				minIP = sw.Subnets[0].AssignedIPAddressMin
				maxIP = sw.Subnets[0].AssignedIPAddressMax

				return internet, nil
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: testInternetRead,
		},
		Updates: []*testutil.CRUDTestFunc{
			// add subnet
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					// add subnet
					subnet, err := client.AddSubnet(ctx, testZone, ctx.ID, &iaas.InternetAddSubnetRequest{
						NetworkMaskLen: 28,
						NextHop:        minIP,
					})
					if err != nil {
						return nil, err
					}

					if !assert.Len(t, subnet.IPAddresses, 16) {
						return nil, errors.New("unexpected state: Subnet.IPAddresses")
					}
					if !assert.Equal(t, minIP, subnet.NextHop) {
						return nil, errors.New("unexpected state: Subnet.NextHop")
					}
					subnetID = subnet.ID
					return subnet, nil
				},
				SkipExtractID: true,
			},
			// update subnet
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					subnet, err := client.UpdateSubnet(ctx, testZone, ctx.ID, subnetID, &iaas.InternetUpdateSubnetRequest{
						NextHop: maxIP,
					})
					if err != nil {
						return nil, err
					}

					if !assert.Len(t, subnet.IPAddresses, 16) {
						return nil, errors.New("unexpected state: Subnet.IPAddresses")
					}
					if !assert.Equal(t, maxIP, subnet.NextHop) {
						return nil, errors.New("unexpected state: Subnet.NextHop")
					}
					return subnet, nil
				},
				SkipExtractID: true,
			},
			// delete subnet
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					return nil, client.DeleteSubnet(ctx, testZone, ctx.ID, subnetID)
				},
				SkipExtractID: true,
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testInternetDelete,
		},
	})
}
