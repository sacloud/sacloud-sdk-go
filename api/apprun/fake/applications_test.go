// Copyright 2021-2024 The sacloud/apprun-api-go authors
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

package fake

import (
	"encoding/json"
	"testing"
	"time"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func TestEngine_Application(t *testing.T) {
	t.Run("list applications", func(t *testing.T) {
		engine := NewEngine()
		for i := 0; i < 3; i++ {
			req := postApplicationBody()
			_, err := engine.CreateApplication(req)
			require.NoError(t, err)
		}

		pageNum := 1
		pageSize := 2
		sortField := "created_at"
		sortOrder := v1.ListApplicationsParamsSortOrderDesc
		resp, err := engine.ListApplications(v1.ListApplicationsParams{
			PageNum:   &pageNum,
			PageSize:  &pageSize,
			SortField: &sortField,
			SortOrder: &sortOrder,
		})
		require.NoError(t, err)

		d := resp.Data
		d0 := d[0]
		d1 := d[1]

		respJson, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedJSON := `
		{
			"meta": {
				"object_total": 3,
				"page_num": 1,
				"page_size": 2,
				"sort_field": "created_at",
				"sort_order": "desc"
			},
			"data": [
				{
					"id": "` + d0.Id + `",
					"name": "` + d0.Name + `",
					"status": "` + string(d0.Status) + `",
					"public_url": "` + d0.PublicUrl + `",
					"created_at": "` + d0.CreatedAt.Format(time.RFC3339) + `"
				},
				{
					"id": "` + d1.Id + `",
					"name": "` + d1.Name + `",
					"status": "` + string(d1.Status) + `",
					"public_url": "` + d1.PublicUrl + `",
					"created_at": "` + d1.CreatedAt.Format(time.RFC3339) + `"
				}
			]
		}`
		require.JSONEq(t, expectedJSON, string(respJson))
	})

	t.Run("create application", func(t *testing.T) {
		engine := NewEngine()
		req := postApplicationBody()
		resp, err := engine.CreateApplication(req)
		require.NoError(t, err)

		respJson, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedJSON := `
		{
			"id": "` + resp.Id + `",
			"name": "app1",
			"timeout_seconds": 20,
			"port": 8081,
			"min_scale": 1,
			"max_scale": 10,
			"scale_target_concurrency": 100,
			"components": [
				{
					"name": "component1",
					"max_cpu": "0.2",
					"max_memory": "512Mi",
					"deploy_source": {
						"container_registry": {
							"image": "apprun-example.sakuracr.jp/helloworld:latest",
							"server": "apprun-example.sakuracr.jp",
							"username": "apprun"
						}
					},
					"env": [
						{
							"key": "envkey",
							"value": "envvalue"
						}
					],
					"probe": {
						"http_get": {
							"path": "/healthz",
							"port": 8080,
							"headers": [
								{
									"name": "Custom-Header",
									"value": "Awesome"
								}
							]
						}
					}
				}
			],
			"status": "Healthy",
			"public_url": "` + resp.PublicUrl + `",
			"resource_id": "` + resp.ResourceId + `",
			"created_at": "` + resp.CreatedAt.Format(time.RFC3339) + `"
		}`
		require.JSONEq(t, expectedJSON, string(respJson))
	})

	t.Run("read application", func(t *testing.T) {
		engine := NewEngine()
		req := postApplicationBody()
		createResp, err := engine.CreateApplication(req)
		require.NoError(t, err)

		readResp, err := engine.ReadApplication(createResp.Id)
		require.NoError(t, err)

		respJson, err := json.Marshal(readResp)
		require.NoError(t, err)

		expectedJSON := `
		{
			"id": "` + readResp.Id + `",
			"name": "app1",
			"timeout_seconds": 20,
			"port": 8081,
			"min_scale": 1,
			"max_scale": 10,
			"scale_target_concurrency": 100,
			"components": [
				{
					"name": "component1",
					"max_cpu": "0.2",
					"max_memory": "512Mi",
					"deploy_source": {
						"container_registry": {
							"image": "apprun-example.sakuracr.jp/helloworld:latest",
							"server": "apprun-example.sakuracr.jp",
							"username": "apprun"
						}
					},
					"env": [
						{
							"key": "envkey",
							"value": "envvalue"
						}
					],
					"probe": {
						"http_get": {
							"path": "/healthz",
							"port": 8080,
							"headers": [
								{
									"name": "Custom-Header",
									"value": "Awesome"
								}
							]
						}
					}
				}
			],
			"status": "Healthy",
			"public_url": "` + readResp.PublicUrl + `",
			"resource_id": "` + readResp.ResourceId + `",
			"created_at": "` + readResp.CreatedAt.Format(time.RFC3339) + `"
		}`
		require.JSONEq(t, expectedJSON, string(respJson))
	})

	t.Run("update application", func(t *testing.T) {
		engine := NewEngine()
		req := postApplicationBody()
		createdApp, err := engine.CreateApplication(req)
		require.NoError(t, err)

		timeoutUpdated := 20
		patchedApp, err := engine.UpdateApplication(createdApp.Id, &v1.PatchApplicationBody{
			TimeoutSeconds: &timeoutUpdated,
		})
		require.NoError(t, err)
		require.Equal(t, timeoutUpdated, patchedApp.TimeoutSeconds)

		require.Equal(t, len(engine.Versions), 2)
	})

	t.Run("delete application", func(t *testing.T) {
		engine := NewEngine()
		for i := 0; i < 3; i++ {
			req := postApplicationBody()
			_, err := engine.CreateApplication(req)
			require.NoError(t, err)
		}

		id := engine.Applications[0].Id
		err := engine.DeleteApplication(id)
		require.NoError(t, err)

		require.Equal(t, len(engine.appVersionRelations), 2)
	})
}

func postApplicationBody() *v1.PostApplicationBody {
	server, userName, password := "apprun-example.sakuracr.jp", "apprun", "apprun" //nolint:gosec
	envKey, envValue := "envkey", "envvalue"
	headerName, headerValue := "Custom-Header", "Awesome"
	probe := v1.PostApplicationBodyComponentProbe{
		HttpGet: &v1.PostApplicationBodyComponentProbeHttpGet{
			Path: "/healthz",
			Port: 8080,
			Headers: &[]v1.PostApplicationBodyComponentProbeHttpGetHeader{
				{
					Name:  &headerName,
					Value: &headerValue,
				},
			},
		},
	}
	req := &v1.PostApplicationBody{
		Name:                   "app1",
		Port:                   8081,
		MinScale:               1,
		MaxScale:               10,
		ScaleTargetConcurrency: saclient.Ptr(100),
		Components: []v1.PostApplicationBodyComponent{
			{
				Name:      "component1",
				MaxCpu:    "0.2",
				MaxMemory: "512Mi",
				DeploySource: v1.PostApplicationBodyComponentDeploySource{
					ContainerRegistry: &v1.PostApplicationBodyComponentDeploySourceContainerRegistry{
						Image:    "apprun-example.sakuracr.jp/helloworld:latest",
						Server:   &server,
						Username: &userName,
						Password: &password,
					},
				},
				Env: &[]v1.PostApplicationBodyComponentEnv{
					{
						Key:   &envKey,
						Value: &envValue,
					},
				},
				Probe: &probe,
			},
		},
		TimeoutSeconds: 20,
	}

	return req
}
