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

package seg_test

import (
	"context"
	"errors"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	seg "github.com/sacloud/service-endpoint-gateway-api-go"
	v1 "github.com/sacloud/service-endpoint-gateway-api-go/apis/v1"
	"github.com/stretchr/testify/assert"
)

func segAPISetup(t *testing.T) (ctx context.Context, api seg.ServiceEndpointGatewayAPI) {
	ctx = t.Context()
	var saClient saclient.Client

	client, err := seg.NewClient(&saClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	api = seg.NewServiceEndpointGatewayOp(client)

	return ctx, api
}

func TestOpFULL(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET",
		"SAKURA_SEG_SWITCH_ID", "SAKURA_SEG_SERVER_IP", "SAKURA_SEG_NETMASK_LEN",
		"SAKURA_SEG_CR_ENDPOINTS", "SAKURA_SEG_MONITORING_ENDPOINTS",
		"SAKURA_SEG_DNS_PRIVATEZONE", "SAKURA_SEG_DNS_UPSTREAM_SERVER_1", "SAKURA_SEG_DNS_UPSTREAM_SERVER_2")(t)
	ctx, segAPI := segAPISetup(t)

	// prepare for create (by iaas API)
	switchID := os.Getenv("SAKURA_SEG_SWITCH_ID")
	serverIPAddress := os.Getenv("SAKURA_SEG_SERVER_IP")

	// netmask length for construct instance, should be between 1 and 32
	networkMaskLen, err := strconv.ParseInt(os.Getenv("SAKURA_SEG_NETMASK_LEN"), 10, 32)
	int32NetworkMaskLen := int32(networkMaskLen)
	if err != nil {
		t.Fatalf("invalid SAKURA_SEG_NETMASK_LEN(valid:1-32): %v", err)
	}

	// create seg instance id for control and cleanup in defer
	id := ""

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
		resp, err := segAPI.Create(ctx, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		id = resp.Appliance.ID

		err = waitForInstanceStatus(t, ctx, segAPI, id, v1.ModelsInstanceInstanceStatusUp)
		if err != nil {
			t.Fatalf("instance did not become up in time: %v", err)
		}
	})

	defer func() {
		if id != "" {
			err := delete(t, ctx, segAPI, id)
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
		resp, err := segAPI.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
	})
	t.Run("Read", func(t *testing.T) {
		resp, err := segAPI.Read(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
	})
	t.Run("ReadInterface", func(t *testing.T) {
		interfaceID := "1"
		resp, err := segAPI.ReadInterface(ctx, id, interfaceID)
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
		dnsPrivateZone := os.Getenv("SAKURA_SEG_DNS_PRIVATEZONE")
		dnsUpstreamServer1 := os.Getenv("SAKURA_SEG_DNS_UPSTREAM_SERVER_1")
		dnsUpstreamServer2 := os.Getenv("SAKURA_SEG_DNS_UPSTREAM_SERVER_2")
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
				MonitoringSuite: v1.OptModelsSettingsMonitoringSuiteSettings{
					Value: v1.ModelsSettingsMonitoringSuiteSettings{
						Enabled: v1.ModelsSettingsMonitoringSuiteSettingsEnabledTrue,
					},
					Set: true,
				},
				DNSForwarding: v1.OptModelsSettingsDNSForwardingSettings{
					Value: v1.ModelsSettingsDNSForwardingSettings{
						Enabled:           v1.ModelsSettingsDNSForwardingSettingsEnabledTrue,
						PrivateHostedZone: dnsPrivateZone,
						UpstreamDNS1:      dnsUpstreamServer1,
						UpstreamDNS2:      dnsUpstreamServer2,
					},
					Set: true,
				},
			},
		}

		request := v1.ModelsApplianceApplianceUpdateRequest{
			Appliance: v1.ModelsApplianceApplianceUpdateBody{
				Settings: settings,
			},
		}
		resp, err := segAPI.Update(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}

		err = segAPI.Apply(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error on apply: %v", err)
		}
		res, err := segAPI.Read(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error on read: %v", err)
		}
		if res == nil {
			t.Fatal("expected response but got nil")
		}
		assert.Equal(t, settings.ServiceEndpointGateway.EnabledServices, res.Appliance.Settings.Value.ServiceEndpointGateway.EnabledServices)
	})
	t.Run("Shutdown and PowerOn", func(t *testing.T) {
		err := segAPI.Shutdown(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error on shutdown: %v", err)
		}
		err = waitForInstanceStatus(t, ctx, segAPI, id, v1.ModelsInstanceInstanceStatusDown)
		if err != nil {
			t.Fatalf("instance did not become down in time: %v", err)
		}
		_, err = segAPI.PowerOn(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error on power on: %v", err)
		}
		err = waitForInstanceStatus(t, ctx, segAPI, id, v1.ModelsInstanceInstanceStatusUp)
		if err != nil {
			t.Fatalf("instance did not become up in time: %v", err)
		}
	})

	t.Run("Reset", func(t *testing.T) {
		err := segAPI.Reset(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error on reset: %v", err)
		}
	})

	t.Run("GetPowerStatus", func(t *testing.T) {
		resp, err := segAPI.ReadPowerStatus(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error on get power status: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		assert.Equal(t, v1.ModelsInstanceInstanceForPowerStatusUp, resp.Instance.Status)
	})

	t.Run("Delete", func(t *testing.T) {
		err := delete(t, ctx, segAPI, id)
		if err != nil {
			t.Fatalf("unexpected error on delete: %v", err)
		}
		id = "" // prevent double delete in defer
	})
}

func delete(t *testing.T, ctx context.Context, api seg.ServiceEndpointGatewayAPI, id string) error {
	needShutdown, err := checkInstanceStatus(t, ctx, api, id, v1.ModelsInstanceInstanceStatusUp)
	if err != nil {
		t.Logf("Failed to check instance status before delete, attempting shutdown just in case: %v", err)
		return err
	}
	if needShutdown {
		err := api.Shutdown(ctx, id)
		if err != nil {
			t.Logf("Failed to shutdown instance before delete: %v", err)
			return err
		}
		err = waitForInstanceStatus(t, ctx, api, id, v1.ModelsInstanceInstanceStatusDown)
		if err != nil {
			t.Logf("Instance did not become down in time before delete: %v", err)
			return err
		}
	}
	return api.Delete(ctx, id)
}

func waitForInstanceStatus(t *testing.T, ctx context.Context, api seg.ServiceEndpointGatewayAPI, id string, status v1.ModelsInstanceInstanceStatus) error {
	withTimeout, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		ok, err := checkInstanceStatus(t, ctx, api, id, status)
		if err != nil {
			return err
		}
		if ok {
			return nil // desired status reached
		}
		select {
		case <-withTimeout.Done():
			return errors.New("timeout waiting for condition")
		case <-ticker.C:
			// retry
		}
	}
}

func checkInstanceStatus(t *testing.T, ctx context.Context, api seg.ServiceEndpointGatewayAPI, id string,
	requestStatus v1.ModelsInstanceInstanceStatus) (bool, error) {
	resp, err := api.Read(ctx, id)
	if err != nil {
		t.Errorf("failed to read instance status: %v", err)
		return false, err
	}
	if resp == nil {
		return false, nil
	}
	currentStatus, set := resp.Appliance.Instance.Status.Get()
	if !set {
		return false, nil
	}
	return currentStatus == requestStatus, nil
}
