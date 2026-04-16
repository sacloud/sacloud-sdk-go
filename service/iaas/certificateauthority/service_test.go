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

package certificateauthority

import (
	"testing"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/certificateauthority/builder"
)

func TestCertificateAuthorityService_CRUD(t *testing.T) {
	if !testutil.IsAccTest() {
		t.Skip("TestCertificateAuthorityService_CRUD only exec when running an Acceptance Test")
	}

	svc := New(testutil.SingletonAPICaller())
	name := testutil.ResourceName("ca")

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		PreCheck:           nil,
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Setup:              nil,
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
				return svc.Create(&CreateRequest{
					Name:             name,
					Description:      "test",
					Tags:             types.Tags{"tag1", "tag2"},
					Country:          "JP",
					Organization:     "usacloud",
					OrganizationUnit: []string{"ou1", "ou2"},
					CommonName:       "www.usacloud.jp",
					NotAfter:         time.Now().Add(365 * 24 * time.Hour),
					Clients: []*builder.ClientCert{
						{
							Country:        "JP",
							Organization:   "usacloud",
							CommonName:     "client.usacloud.jp",
							NotAfter:       time.Now().Add(365 * 24 * time.Hour),
							IssuanceMethod: types.CertificateAuthorityIssuanceMethods.URL,
						},
					},
				})
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
				return svc.Read(&ReadRequest{ID: ctx.ID})
			},
		},
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) error {
				return svc.Delete(&DeleteRequest{ID: ctx.ID})
			},
		},
	})
}
