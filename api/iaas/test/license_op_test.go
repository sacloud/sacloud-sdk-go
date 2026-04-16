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
)

func TestLicenseOpCRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		IgnoreStartupWait:  true,
		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: testLicenseCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createLicenseExpected,
				IgnoreFields: ignoreLicenseFields,
			}),
		},
		Read: &testutil.CRUDTestFunc{
			Func: testLicenseRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createLicenseExpected,
				IgnoreFields: ignoreLicenseFields,
			}),
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testLicenseUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateLicenseExpected,
					IgnoreFields: ignoreLicenseFields,
				}),
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testLicenseDelete,
		},
	})
}

var (
	ignoreLicenseFields = []string{
		"ID",
		"CreatedAt",
		"ModifiedAt",
	}

	createLicenseParam = &iaas.LicenseCreateRequest{
		Name:          testutil.ResourceName("license"),
		LicenseInfoID: types.ID(10001),
	}
	createLicenseExpected = &iaas.License{
		Name:            createLicenseParam.Name,
		LicenseInfoID:   createLicenseParam.LicenseInfoID,
		LicenseInfoName: "Windows RDS SAL",
	}
	updateLicenseParam = &iaas.LicenseUpdateRequest{
		Name: testutil.ResourceName("license-upd"),
	}
	updateLicenseExpected = &iaas.License{
		Name:            updateLicenseParam.Name,
		LicenseInfoID:   createLicenseParam.LicenseInfoID,
		LicenseInfoName: "Windows RDS SAL",
	}
)

func testLicenseCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewLicenseOp(caller)
	return client.Create(ctx, createLicenseParam)
}

func testLicenseRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewLicenseOp(caller)
	return client.Read(ctx, ctx.ID)
}

func testLicenseUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewLicenseOp(caller)
	return client.Update(ctx, ctx.ID, updateLicenseParam)
}

func testLicenseDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewLicenseOp(caller)
	return client.Delete(ctx, ctx.ID)
}
