// Copyright 2021-2026 The sacloud/apprun-api-go authors
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
	"github.com/stretchr/testify/require"
)

func TestEngine_Version(t *testing.T) {
	t.Run("list versions", func(t *testing.T) {
		engine := NewEngine()
		req := postApplicationBody()
		createdApp, err := engine.CreateApplication(req)
		require.NoError(t, err)

		timeoutUpdated := 20
		patchedApp, err := engine.UpdateApplication(createdApp.ID, &v1.PatchApplicationBody{
			TimeoutSeconds: v1.NewOptInt(timeoutUpdated),
		})
		require.NoError(t, err)

		pageNum := 1
		pageSize := 2
		sortField := "created_at"
		sortOrder := v1.ListApplicationVersionsSortOrderDesc
		resp, err := engine.ListVersions(patchedApp.ID, v1.ListApplicationVersionsParams{
			PageNum:   v1.NewOptInt(pageNum),
			PageSize:  v1.NewOptInt(pageSize),
			SortField: v1.NewOptString(sortField),
			SortOrder: v1.NewOptListApplicationVersionsSortOrder(sortOrder),
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
				"object_total": 2,
				"page_num": 1,
				"page_size": 2,
				"sort_field": "created_at",
				"sort_order": "desc"
			},
			"data": [
				{
					"id": "` + d0.ID + `",
					"name": "` + d0.Name + `",
					"status": "` + string(d0.Status) + `",
					"created_at": "` + d0.CreatedAt.Format(time.RFC3339) + `"
				},
				{
					"id": "` + d1.ID + `",
					"name": "` + d1.Name + `",
					"status": "` + string(d1.Status) + `",
					"created_at": "` + d1.CreatedAt.Format(time.RFC3339) + `"
				}
			]
		}`
		require.JSONEq(t, expectedJSON, string(respJson))
	})

	t.Run("read version", func(t *testing.T) {
		engine := NewEngine()
		req := postApplicationBody()
		createdApp, err := engine.CreateApplication(req)
		require.NoError(t, err)

		timeoutUpdated := 20
		patchedApp, err := engine.UpdateApplication(createdApp.ID, &v1.PatchApplicationBody{
			TimeoutSeconds: v1.NewOptInt(timeoutUpdated),
		})
		require.NoError(t, err)

		r := engine.appVersionRelations[patchedApp.ID][0]
		resp, err := engine.ReadVersion(r.application.ID, r.version.ID)
		require.NoError(t, err)

		respJson, err := json.Marshal(resp)
		require.NoError(t, err)

		expectedJSON := `
		{
			"id": "` + r.version.ID + `",
			"name": "` + r.version.Name + `",
			"status": "Healthy",
			"timeout_seconds": 20,
			"port": 8081,
			"min_scale": 1,
			"max_scale": 10,
			"scale_target_concurrency": 100,
			"components": [
				{
					"name": "component1",
					"max_cpu": "0.5",
					"max_memory": "1Gi",
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
			"created_at": "` + r.application.CreatedAt.Format(time.RFC3339) + `"
		}`
		require.JSONEq(t, expectedJSON, string(respJson))
	})

	t.Run("read version status", func(t *testing.T) {
		engine := NewEngine()
		req := postApplicationBody()
		createdApp, err := engine.CreateApplication(req)
		require.NoError(t, err)

		r := engine.appVersionRelations[createdApp.ID][0]
		resp, err := engine.ReadVersionStatus(r.application.ID, r.version.ID)
		require.NoError(t, err)

		require.Equal(t, "", resp.Message)
		require.Equal(t, v1.HandlerGetApplicationVersionOnlyStatusStatus(r.version.Status), resp.Status)
	})

	t.Run("delete version", func(t *testing.T) {
		engine := NewEngine()
		req := postApplicationBody()
		createdApp, err := engine.CreateApplication(req)
		require.NoError(t, err)

		timeoutUpdated := 20
		_, err = engine.UpdateApplication(createdApp.ID, &v1.PatchApplicationBody{
			TimeoutSeconds: v1.NewOptInt(timeoutUpdated),
		})
		require.NoError(t, err)
		require.Equal(t, len(engine.appVersionRelations[createdApp.ID]), 2)

		r := engine.appVersionRelations[createdApp.ID][0]
		err = engine.DeleteVersion(r.application.ID, r.version.ID)
		require.NoError(t, err)
		require.Equal(t, len(engine.appVersionRelations[createdApp.ID]), 1)
	})
}
