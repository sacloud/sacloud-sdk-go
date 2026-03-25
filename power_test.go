// Copyright 2026- The sacloud/service-endpoint-gateway-api-go Authors
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

package service_endpoint_gateway_test

import (
	"context"
	"os"
	"testing"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	service_endpoint_gateway "github.com/sacloud/service-endpoint-gateway-api-go"
	v1 "github.com/sacloud/service-endpoint-gateway-api-go/apis/v1"
)

func powerAPISetup(t *testing.T) (ctx context.Context, api service_endpoint_gateway.PowerAPI) {
	ctx = t.Context()
	var saClient saclient.Client

	client, err := service_endpoint_gateway.NewClient(&saClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	api = service_endpoint_gateway.NewPowerOp(client)

	return ctx, api
}
func TestPowerOp(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_SEG_ID")(t)
	ctx, powerAPI := powerAPISetup(t)
	id := os.Getenv("SAKURA_SEG_ID")
	status := ""
	t.Run("Read", func(t *testing.T) {
		res, err := powerAPI.Read(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil {
			t.Fatalf("failed to read appliance with id %s", id)
		}
		t.Logf("Current power status: %+v", res)
		status = string(res.Instance.Status)
	})

	// Note: Update is used to power on the appliance, before power off.
	t.Run("Update", func(t *testing.T) {
		if status != string(v1.ModelsInstanceInstanceForPowerStatusDown) {
			t.Logf("Skipping Delete test because appliance is not powered on (current status: %s)", status)
			return
		}
		res, err := powerAPI.Update(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if res == nil {
			t.Fatalf("failed to read appliance with id %s", id)
		}
		t.Logf("Current power status: %+v", res)
	})
	t.Run("Delete", func(t *testing.T) {
		if status != string(v1.ModelsInstanceInstanceForPowerStatusUp) {
			t.Logf("Skipping Delete test because appliance is not powered on (current status: %s)", status)
			return
		}
		err := powerAPI.Delete(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		t.Logf("Successfully deleted appliance with id %s", id)
	})
	t.Run("Reset", func(t *testing.T) {
		err := powerAPI.Reset(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		t.Logf("Successfully reset appliance with id %s", id)
	})
}
