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
	"fmt"
	"os"

	"github.com/sacloud/saclient-go"
	seg "github.com/sacloud/service-endpoint-gateway-api-go"
	v1 "github.com/sacloud/service-endpoint-gateway-api-go/apis/v1"
)

func Example() {
	// setup
	// TODO replace your access token/secret
	os.Setenv("SAKURA_ACCESS_TOKEN", "your-token")         //nolint:errcheck,gosec
	os.Setenv("SAKURA_ACCESS_TOKEN_SECRET", "your-secret") //nolint:errcheck,gosec

	theClient := saclient.Client{}
	client, err := seg.NewClient(&theClient)
	if err != nil {
		panic(err)
	}

	constructAPI := seg.NewConstructOp(client)
	ctx := context.Background()

	// create request
	createRequest := v1.ModelsApplianceApplianceCreateRequest{
		Appliance: v1.ModelsApplianceApplianceCreateBody{
			Remark: v1.ModelsRemarkApplianceCreateRemark{
				Switch: v1.ModelsRemarkSwitchRemark{
					ID: "your-switch-id",
				},
				Network: v1.ModelsRemarkNetworkRemark{
					NetworkMaskLen: 24, // your netmask length for construct instance, should be between 1 and 32
				},
				Servers: []v1.ModelsRemarkServerRemark{
					{
						IPAddress: "your-server-ip-address",
					},
				},
			},
		},
	}

	// create call (auto power on after creation)
	created, err := constructAPI.Create(ctx, createRequest)
	if err != nil {
		panic(err)
	}
	fmt.Println(created)

	// read call
	read, err := constructAPI.Read(ctx, created.Appliance.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(read)

	// list call
	listed, err := constructAPI.List(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(listed)

	// update request
	updateRequest := v1.ModelsApplianceApplianceUpdateRequest{
		Appliance: v1.ModelsApplianceApplianceUpdateBody{
			Settings: v1.ModelsSettingsApplianceSettings{
				ServiceEndpointGateway: v1.ModelsSettingsServiceEndpointGatewaySettings{
					EnabledServices: []v1.ModelsSettingsEnabledService{
						{
							Type: v1.ModelsSettingsEnabledServiceTypeObjectStorage,
							Config: v1.ModelsSettingsServiceConfig{
								Endpoints: []string{
									"objectstorage-endpoint", // "s3.isk01.sakurastorage.jp" etc...
								},
							},
						},
						{
							Type: v1.ModelsSettingsEnabledServiceTypeMonitoringSuite,
							Config: v1.ModelsSettingsServiceConfig{
								Endpoints: []string{
									"monitoring-endpoint", // "XXXXXXXXXX.logs.monitoring.global.api.sacloud.jp"
								},
							},
						},
						{
							Type: v1.ModelsSettingsEnabledServiceTypeContainerRegistry,
							Config: v1.ModelsSettingsServiceConfig{
								Endpoints: []string{
									"container-registry-endpoint", // "XXXXXXXX.sakuracr.jp" etc...
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
			},
		},
	}
	// update call (wait until seg instance is up)
	updated, err := constructAPI.Update(ctx, created.Appliance.ID, updateRequest)
	if err != nil {
		panic(err)
	}
	fmt.Println(updated)

	// apply call(update request is not applied until this call is made)
	err = constructAPI.Apply(ctx, created.Appliance.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(updated)

	// power off before delete (wait until seg instance is up before power off, otherwise power off call will fail)
	powerAPI := seg.NewPowerOp(client)
	err = powerAPI.Delete(ctx, created.Appliance.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println("powered off")
	// delete
	if err := constructAPI.Delete(ctx, created.Appliance.ID); err != nil {
		panic(err)
	}
	fmt.Println("deleted")
}
