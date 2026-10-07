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

package bridge

import (
	"testing"

	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/testutil"
)

func TestBridgeService_CRUD(t *testing.T) {
	svc := New(testutil.SingletonAPICaller())
	name := testutil.ResourceName("bridge")
	testZone := testutil.TestZone()

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		IgnoreStartupWait:  true,
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (any, error) {
				return svc.Create(&CreateRequest{
					Zone:        testZone,
					Name:        name,
					Description: "desc",
				})
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (any, error) {
				return svc.Read(&ReadRequest{ID: ctx.ID, Zone: testZone})
			},
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (any, error) {
					return svc.Update(&UpdateRequest{
						ID:          ctx.ID,
						Name:        new(name + "-upd"),
						Description: new("test-upd"),
						Zone:        testZone,
					})
				},
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) error {
				return svc.Delete(&DeleteRequest{ID: ctx.ID, Zone: testZone})
			},
		},
	})
}
