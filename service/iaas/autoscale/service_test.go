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

package autoscale

import (
	"fmt"
	"os"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/pointer"
	"github.com/sacloud/packages-go/size"
)

func TestAutoScaleService_CRUD(t *testing.T) {
	apiKey := os.Getenv("SAKURACLOUD_API_KEY_ID")
	if testutil.IsAccTest() && apiKey == "" {
		t.Skip("SAKURACLOUD_API_KEY_ID is required when running the acceptance test")
	}
	if apiKey == "" {
		apiKey = "dummy"
	}

	svc := New(testutil.SingletonAPICaller())
	name := testutil.ResourceName("auto-scale")
	var server *iaas.Server

	testutil.RunCRUD(t, &testutil.CRUDTestCase{
		Parallel:           true,
		PreCheck:           nil,
		SetupAPICallerFunc: testutil.SingletonAPICaller,
		Setup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			serverOp := iaas.NewServerOp(caller)
			// ディスクレスサーバを作成
			created, err := serverOp.Create(ctx, testutil.TestZone(), &iaas.ServerCreateRequest{
				CPU:      1,
				MemoryMB: 2 * size.GiB,
				Name:     testServerName,
			})
			server = created
			return err
		},
		Create: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
				return svc.Create(&CreateRequest{
					Name:        name,
					Description: "test",
					Tags:        types.Tags{"tag1", "tag2"},
					Zones:       []string{testutil.TestZone()},
					Config:      fmt.Sprintf(autoScaleConfigTemplate, testServerName, testutil.TestZone()),
					CPUThresholdScaling: &CreateCPUThresholdScaling{
						ServerPrefix: testServerName,
						Up:           80,
						Down:         50,
					},
					APIKeyID: apiKey,
				})
			},
		},
		Read: &testutil.CRUDTestFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
				return svc.Read(&ReadRequest{ID: ctx.ID})
			},
		},
		Updates: []*testutil.CRUDTestFunc{
			{
				Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) (interface{}, error) {
					return svc.Update(&UpdateRequest{
						ID:          ctx.ID,
						Name:        pointer.NewString(name + "-upd"),
						Description: pointer.NewString("test-upd"),
						Tags:        &types.Tags{"tag1-upd", "tag2-upd"},
						Zones:       &[]string{testutil.TestZone()},
						Config:      pointer.NewString(fmt.Sprintf(autoScaleConfigTemplateUpd, testServerName, testutil.TestZone())),
						CPUThresholdScaling: &UpdateCPUThresholdScaling{
							Up:   pointer.NewInt(80),
							Down: pointer.NewInt(50),
						},
					})
				},
			},
		},
		Shutdown: nil,
		Delete: &testutil.CRUDTestDeleteFunc{
			Func: func(ctx *testutil.CRUDTestContext, _ iaas.APICaller) error {
				return svc.Delete(&DeleteRequest{ID: ctx.ID})
			},
		},
		Cleanup: func(ctx *testutil.CRUDTestContext, caller iaas.APICaller) error {
			serverOp := iaas.NewServerOp(caller)
			return serverOp.Delete(ctx, testutil.TestZone(), server.ID)
		},
	})
}

var (
	testServerName          = testutil.ResourceName("auto-scale")
	autoScaleConfigTemplate = `
resources:
  - type: Server
    selector:
      names: ["%s"]
      zones: ["%s"]

    shutdown_force: true
`

	autoScaleConfigTemplateUpd = `
resources:
  - type: Server
    selector:
      names: ["%s"]
      zones: ["%s"]

    shutdown_force: true

autoscaler:
  cooldown: 300
`
)
