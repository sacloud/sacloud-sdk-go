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
	"fmt"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	sacloudtestutil "github.com/sacloud/packages-go/testutil"
)

func TestEnhancedDBOp_CRUD(t *testing.T) {
	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		SetupAPICallerFunc: singletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: testEnhancedDBCreate,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createEnhancedDBExpected,
				IgnoreFields: ignoreEnhancedDBFields,
			}),
		},

		Read: &testutil.CRUDTestFunc{
			Func: testEnhancedDBRead,
			CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
				ExpectValue:  createEnhancedDBExpected,
				IgnoreFields: ignoreEnhancedDBFields,
			}),
		},

		Updates: []*testutil.CRUDTestFunc{
			{
				Func: testEnhancedDBUpdate,
				CheckFunc: testutil.AssertEqualWithExpected(&testutil.CRUDTestExpect{
					ExpectValue:  updateEnhancedDBExpected,
					IgnoreFields: ignoreEnhancedDBFields,
				}),
			},
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					edbOp := iaas.NewEnhancedDBOp(caller)
					return nil, edbOp.SetConfig(ctx, ctx.ID, &iaas.EnhancedDBSetConfigRequest{
						AllowedNetworks: []string{"192.0.2.1/32"},
					})
				},
				SkipExtractID: true,
			},
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					edbOp := iaas.NewEnhancedDBOp(caller)
					config, err := edbOp.GetConfig(ctx, ctx.ID)
					if err != nil {
						return nil, err
					}
					if config.MaxConnections != 50 {
						return nil, fmt.Errorf("got unexpected value: MaxConnections: expect: %d actual: %d", 50, config.MaxConnections)
					}
					if testutil.IsAccTest() {
						if len(config.AllowedNetworks) != 1 || config.AllowedNetworks[0] != "192.0.2.1/32" {
							return nil, fmt.Errorf("got unexpected value: AllowedNetworks: expect: %s actual: %s", "[192.0.2.1/32]", config.AllowedNetworks)
						}
					}
					return nil, nil
				},
				SkipExtractID: true,
			},
			{
				Func: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
					edbOp := iaas.NewEnhancedDBOp(caller)
					return nil, edbOp.SetPassword(ctx, ctx.ID, &iaas.EnhancedDBSetPasswordRequest{
						Password: "password",
					})
				},
				SkipExtractID: true,
			},
		},

		Delete: &testutil.CRUDTestDeleteFunc{
			Func: testEnhancedDBDelete,
		},
	})
}

var (
	ignoreEnhancedDBFields = []string{
		"ID",
		"Class",
		"SettingsHash",
		"CreatedAt",
		"ModifiedAt",
	}
	createEnhancedDBParam = &iaas.EnhancedDBCreateRequest{
		Name:         testutil.ResourceName("enhanced-db"),
		Description:  "desc",
		Tags:         []string{"tag1", "tag2"},
		DatabaseName: sacloudtestutil.RandomName("", 32, sacloudtestutil.CharSetAlpha),
		DatabaseType: types.EnhancedDBTypesMariaDB,
		Region:       types.EnhancedDBRegionsTk1,
	}
	createEnhancedDBExpected = &iaas.EnhancedDB{
		Name:         createEnhancedDBParam.Name,
		Description:  createEnhancedDBParam.Description,
		Tags:         createEnhancedDBParam.Tags,
		Availability: types.Availabilities.Available,
		DatabaseName: createEnhancedDBParam.DatabaseName,
		DatabaseType: types.EnhancedDBTypesMariaDB,
		Region:       types.EnhancedDBRegionsTk1,
		HostName:     createEnhancedDBParam.DatabaseName + ".mariadb-tk1.db.sakurausercontent.com",
		Port:         3306,
	}
	updateEnhancedDBParam = &iaas.EnhancedDBUpdateRequest{
		Name:        testutil.ResourceName("enhanced-db-upd"),
		Description: "desc-upd",
		Tags:        []string{"tag1-upd", "tag2-upd"},
		IconID:      testIconID,
	}
	updateEnhancedDBExpected = &iaas.EnhancedDB{
		Name:         updateEnhancedDBParam.Name,
		Description:  updateEnhancedDBParam.Description,
		Tags:         updateEnhancedDBParam.Tags,
		Availability: types.Availabilities.Available,
		IconID:       testIconID,
		DatabaseName: createEnhancedDBParam.DatabaseName,
		DatabaseType: types.EnhancedDBTypesMariaDB,
		Region:       types.EnhancedDBRegionsTk1,
		HostName:     createEnhancedDBParam.DatabaseName + ".mariadb-tk1.db.sakurausercontent.com",
		Port:         3306,
	}
)

func testEnhancedDBCreate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewEnhancedDBOp(caller)
	return client.Create(ctx, createEnhancedDBParam)
}

func testEnhancedDBRead(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewEnhancedDBOp(caller)
	return client.Read(ctx, ctx.ID)
}

func testEnhancedDBUpdate(ctx *testutil.CRUDTestContext, caller iaas.APICaller) (interface{}, error) {
	client := iaas.NewEnhancedDBOp(caller)
	return client.Update(ctx, ctx.ID, updateEnhancedDBParam)
}

func testEnhancedDBDelete(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
	client := iaas.NewEnhancedDBOp(caller)
	return client.Delete(ctx, ctx.ID)
}
