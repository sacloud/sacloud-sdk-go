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
	"strconv"
	"testing"
	"time"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	service_endpoint_gateway "github.com/sacloud/service-endpoint-gateway-api-go"
	v1 "github.com/sacloud/service-endpoint-gateway-api-go/apis/v1"
	seg_testutil "github.com/sacloud/service-endpoint-gateway-api-go/testutil"
)

func constructAPISetup(t *testing.T) (ctx context.Context, api service_endpoint_gateway.ConstructAPI) {
	ctx = t.Context()
	var saClient saclient.Client

	client, err := service_endpoint_gateway.NewClient(&saClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	api = service_endpoint_gateway.NewConstructOp(client)

	return ctx, api
}

func TestConstructOpFULL(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)
	ctx, constructAPI := constructAPISetup(t)
	id := ""
	switchID := os.Getenv("SAKURA_SEG_SWITCH_ID")
	networkMaskLen, err := strconv.ParseInt(os.Getenv("SAKURA_SEG_NETMASK_LEN"), 10, 32)
	int32NetworkMaskLen := int32(networkMaskLen)
	if err != nil {
		t.Fatalf("invalid SAKURA_SEG_NETMASK_LEN(valid:1-32): %v", err)
	}
	serverIPAddress := os.Getenv("SAKURA_SEG_SERVER_IP")

	result := t.Run("Create", func(t *testing.T) {
		request := v1.ModelsApplianceApplianceCreateRequest{
			Appliance: v1.ModelsApplianceApplianceCreateBody{
				Remark: v1.ModelsRemarkApplianceCreateRemark{
					Switch: v1.ModelsRemarkSwitchRemark{
						ID: switchID,
					},
					Network: v1.ModelsRemarkNetworkRemark{
						NetworkMaskLen: int32NetworkMaskLen,
					},
					Servers: []v1.ModelsRemarkServerRemark{
						{
							IPAddress: serverIPAddress,
						},
					},
				},
			},
		}
		resp, err := constructAPI.Create(ctx, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		id = resp.Appliance.ID

		// wait until instance is up
		waitDownCtx := context.Background()
		withTimeout, cancel := context.WithTimeout(waitDownCtx, 3*time.Minute)
		defer cancel()
		checkInterval := 5 * time.Second

		err = seg_testutil.WaitUntil(withTimeout, checkInterval, func(ctx context.Context) (bool, error) {
			return checkInstanceStatus(ctx, constructAPI, id, v1.ModelsInstanceInstanceStatusUp)
		})
		if err != nil {
			t.Fatalf("instance did not become up in time: %v", err)
		}
	})

	defer func() {
		if id != "" {
			err := deleteConstruct(t, ctx, constructAPI, id)
			if err != nil {
				t.Fatalf("unexpected error on delete: %v", err)
			}
			t.Log("Defer PowerOp.Delete succeeded")
		}
	}()

	if !result {
		t.Fatal("skipping rest of tests due to Create failure")
	}

	t.Run("List", func(t *testing.T) {
		resp, err := constructAPI.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
	})
	t.Run("Read", func(t *testing.T) {
		resp, err := constructAPI.Read(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
	})

	t.Run("Update and Apply", func(t *testing.T) {
		crEndpoints := os.Getenv("SAKURA_SEG_CR_ENDPOINTS")
		monitorEndpoint := os.Getenv("SAKURA_SEG_MONITORING_ENDPOINTS")
		settings := v1.ModelsSettingsApplianceSettings{
			ServiceEndpointGateway: v1.ModelsSettingsServiceEndpointGatewaySettings{
				EnabledServices: []v1.ModelsSettingsEnabledService{
					{
						Type: v1.ModelsSettingsEnabledServiceTypeObjectStorage,
						Config: v1.ModelsSettingsServiceConfig{
							Endpoints: []string{
								"s3.isk01.sakurastorage.jp",
								"s3.tky01.sakurastorage.jp",
								"s3.arc02.sakurastorage.jp",
							},
						},
					},
					{
						Type: v1.ModelsSettingsEnabledServiceTypeMonitoringSuite,
						Config: v1.ModelsSettingsServiceConfig{
							Endpoints: []string{
								monitorEndpoint,
							},
						},
					},
					{
						Type: v1.ModelsSettingsEnabledServiceTypeContainerRegistry,
						Config: v1.ModelsSettingsServiceConfig{
							Endpoints: []string{
								crEndpoints,
							},
						},
					},
					{
						Type: v1.ModelsSettingsEnabledServiceTypeAppRunDedicatedControlPlane,
						Config: v1.ModelsSettingsServiceConfig{
							Mode: v1.OptModelsSettingsServiceConfigMode{
								Value: v1.ModelsSettingsServiceConfigModeManaged,
								Set:   true,
							},
						},
					},
				},
			},
		}

		request := v1.ModelsApplianceApplianceUpdateRequest{
			Appliance: v1.ModelsApplianceApplianceUpdateBody{
				Settings: settings,
			},
		}
		resp, err := constructAPI.Update(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}

		err = constructAPI.Apply(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error on apply: %v", err)
		}
	})

	t.Run("Delete", func(t *testing.T) {
		err := deleteConstruct(t, ctx, constructAPI, id)
		if err != nil {
			t.Fatalf("unexpected error on delete: %v", err)
		}
		id = "" // prevent double delete in defer
	})
}

func deleteConstruct(t *testing.T, ctx context.Context, api service_endpoint_gateway.ConstructAPI, id string) error {
	_, powerAPI := powerAPISetup(t)
	err := powerAPI.Delete(ctx, id)
	if err != nil {
		t.Fatalf("unexpected error on power delete: %v", err)
	}

	waitDownCtx := context.Background()
	withTimeout, cancel := context.WithTimeout(waitDownCtx, 3*time.Minute)
	defer cancel()

	checkInterval := 5 * time.Second
	err = seg_testutil.WaitUntil(withTimeout, checkInterval, func(ctx context.Context) (bool, error) {
		return checkInstanceStatus(ctx, api, id, v1.ModelsInstanceInstanceStatusDown)
	})
	if err != nil {
		t.Fatalf("instance did not become down in time: %v", err)
	}
	return api.Delete(ctx, id)
}

func checkInstanceStatus(ctx context.Context, api service_endpoint_gateway.ConstructAPI, id string,
	status v1.ModelsInstanceInstanceStatus) (bool, error) {
	resp, err := api.Read(ctx, id)
	if err != nil {
		return false, err
	}
	if resp == nil {
		return false, nil
	}
	currentStatus, set := resp.Appliance.Instance.Status.Get()
	if !set {
		return false, nil
	}
	return currentStatus == status, nil
}
