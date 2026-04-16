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
	"fmt"
	"math/rand/v2"
	"sort"
	"strconv"
	"time"

	"github.com/google/uuid"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

func (engine *Engine) ListApplications(param v1.ListApplicationsParams) (*v1.HandlerListApplications, error) {
	defer engine.rLock()()

	if len(engine.Applications) == 0 {
		return nil, newError(
			ErrorTypeNotFound, "application", nil,
			"アプリケーションが見つかりませんでした。")
	}

	// 各Applicationの最新バージョンのみを取り出す
	var apps []*v1.Application
	for id := range engine.appVersionRelations {
		apps = append(apps, engine.latestApplication(id))
	}

	sort.Slice(apps, func(i int, j int) bool {
		switch *param.SortField {
		case "created_at":
			if *param.SortOrder == "desc" {
				return apps[i].CreatedAt.After(apps[j].CreatedAt)
			}
			return apps[i].CreatedAt.Before(apps[j].CreatedAt)
		// sort_fieldのデフォルト値であるcreated_at以外は未サポート
		default:
			return false
		}
	})

	appsLen := len(apps)
	start := (*param.PageNum - 1) * *param.PageSize
	// 範囲外の場合nilを返す
	if start > appsLen {
		return nil, nil
	}

	end := start + *param.PageSize
	if end > appsLen {
		end = appsLen
	}

	var data []v1.HandlerListApplicationsData
	for _, app := range apps[start:end] {
		if app != nil {
			data = append(data, v1.HandlerListApplicationsData{
				Id:        app.Id,
				Name:      app.Name,
				Status:    (v1.HandlerListApplicationsDataStatus)(app.Status),
				PublicUrl: app.PublicUrl,
				CreatedAt: app.CreatedAt,
			})
		}
	}

	meta := v1.HandlerListApplicationsMeta{
		ObjectTotal: appsLen,
	}
	if param.PageNum != nil {
		meta.PageNum = *param.PageNum
	}
	if param.PageSize != nil {
		meta.PageSize = *param.PageSize
	}
	if param.SortField != nil {
		meta.SortField = *param.SortField
	}
	if param.SortOrder != nil {
		so := (*v1.HandlerListApplicationsMetaSortOrder)(param.SortOrder)
		meta.SortOrder = *so
	}

	return &v1.HandlerListApplications{
		Data: data,
		Meta: meta,
	}, nil
}

func (engine *Engine) CreateApplication(reqBody *v1.PostApplicationBody) (*v1.Application, error) {
	defer engine.lock()()

	appId, err := engine.newId()
	if err != nil {
		return nil, newError(
			ErrorTypeUnknown, "application", nil,
			"Application IDの生成に失敗しました。")
	}

	var components []v1.HandlerApplicationComponent
	for _, reqComponent := range reqBody.Components {
		var cr v1.HandlerApplicationComponentDeploySourceContainerRegistry
		if reqComponent.DeploySource.ContainerRegistry != nil {
			cr.Image = reqComponent.DeploySource.ContainerRegistry.Image
			cr.Server = reqComponent.DeploySource.ContainerRegistry.Server
			cr.Username = reqComponent.DeploySource.ContainerRegistry.Username
		}

		var env []v1.HandlerApplicationComponentEnv
		if reqComponent.Env != nil {
			for _, e := range *reqComponent.Env {
				env = append(env, v1.HandlerApplicationComponentEnv(e))
			}
		}

		var probe v1.HandlerApplicationComponentProbe
		if reqComponent.Probe != nil && reqComponent.Probe.HttpGet != nil {
			headers := []v1.HandlerApplicationComponentProbeHttpGetHeader{}
			if reqComponent.Probe.HttpGet.Headers != nil {
				for _, header := range *reqComponent.Probe.HttpGet.Headers {
					headers = append(headers, v1.HandlerApplicationComponentProbeHttpGetHeader(header))
				}
			}

			probe = v1.HandlerApplicationComponentProbe{
				HttpGet: &v1.HandlerApplicationComponentProbeHttpGet{
					Path:    reqComponent.Probe.HttpGet.Path,
					Port:    reqComponent.Probe.HttpGet.Port,
					Headers: &headers,
				},
			}
		}

		var component v1.HandlerApplicationComponent
		component.Name = reqComponent.Name
		component.MaxCpu = string(reqComponent.MaxCpu)
		component.MaxMemory = string(reqComponent.MaxMemory)
		component.DeploySource.ContainerRegistry = &cr
		component.Env = &env
		component.Probe = &probe
		components = append(components, component)
	}

	status := v1.ApplicationStatusHealthy
	url := fmt.Sprintf("https://example.com/apprun/dummy/%s", appId)
	createdAt := time.Now().UTC().Truncate(time.Second)
	app := &v1.Application{
		Id:                     appId,
		Name:                   reqBody.Name,
		TimeoutSeconds:         reqBody.TimeoutSeconds,
		Port:                   reqBody.Port,
		MinScale:               reqBody.MinScale,
		MaxScale:               reqBody.MaxScale,
		ScaleTargetConcurrency: reqBody.ScaleTargetConcurrency,
		Components:             components,
		Status:                 status,
		PublicUrl:              url,
		ResourceId:             engine.newResourceId(),
		CreatedAt:              createdAt,
	}
	engine.Applications = append(engine.Applications, app)

	err = engine.createVersion(app)
	if err != nil {
		return nil, newError(
			ErrorTypeUnknown, "application", nil,
			"Version の生成に失敗しました。")
	}

	engine.initTraffic(app)

	return app, nil
}

func (engine *Engine) ReadApplication(id string) (*v1.Application, error) {
	defer engine.rLock()()

	if len(engine.Applications) == 0 {
		return nil, newError(
			ErrorTypeNotFound, "application", nil,
			"アプリケーションが見つかりませんでした。")
	}

	app := engine.latestApplication(id)
	if app != nil && app.Id == id {
		return app, nil
	}

	return nil, newError(
		ErrorTypeNotFound, "application", nil,
		"アプリケーションが見つかりませんでした。")
}

func (engine *Engine) UpdateApplication(id string, reqBody *v1.PatchApplicationBody) (*v1.HandlerPatchApplication, error) {
	defer engine.lock()()

	patchedApp := *(engine.latestApplication(id))
	if reqBody.TimeoutSeconds != nil {
		patchedApp.TimeoutSeconds = *reqBody.TimeoutSeconds
	}
	if reqBody.Port != nil {
		patchedApp.Port = *reqBody.Port
	}
	if reqBody.MinScale != nil {
		patchedApp.MinScale = *reqBody.MinScale
	}
	if reqBody.MaxScale != nil {
		patchedApp.MaxScale = *reqBody.MaxScale
	}
	if reqBody.Components != nil && len(*reqBody.Components) > 0 {
		var components []v1.HandlerApplicationComponent
		for _, reqComponent := range *reqBody.Components {
			var cr v1.HandlerApplicationComponentDeploySourceContainerRegistry
			if reqComponent.DeploySource.ContainerRegistry != nil {
				cr.Image = reqComponent.DeploySource.ContainerRegistry.Image
				cr.Server = reqComponent.DeploySource.ContainerRegistry.Server
				cr.Username = reqComponent.DeploySource.ContainerRegistry.Username
			}

			var env []v1.HandlerApplicationComponentEnv
			if reqComponent.Env != nil {
				for _, e := range *reqComponent.Env {
					env = append(env, v1.HandlerApplicationComponentEnv(e))
				}
			}

			var probe v1.HandlerApplicationComponentProbe
			if reqComponent.Probe != nil && reqComponent.Probe.HttpGet != nil {
				headers := []v1.HandlerApplicationComponentProbeHttpGetHeader{}
				if reqComponent.Probe.HttpGet.Headers != nil {
					for _, header := range *reqComponent.Probe.HttpGet.Headers {
						headers = append(headers, v1.HandlerApplicationComponentProbeHttpGetHeader(header))
					}
				}

				probe = v1.HandlerApplicationComponentProbe{
					HttpGet: &v1.HandlerApplicationComponentProbeHttpGet{
						Path:    reqComponent.Probe.HttpGet.Path,
						Port:    reqComponent.Probe.HttpGet.Port,
						Headers: &headers,
					},
				}
			}

			var component v1.HandlerApplicationComponent
			component.Name = reqComponent.Name
			component.MaxCpu = string(reqComponent.MaxCpu)
			component.MaxMemory = string(reqComponent.MaxMemory)
			component.DeploySource.ContainerRegistry = &cr
			component.Env = &env
			component.Probe = &probe
			components = append(components, component)
		}

		patchedApp.Components = components
	}

	now := time.Now().UTC().Truncate(time.Second)
	patchedApp.CreatedAt = now

	engine.Applications = append(engine.Applications, &patchedApp)
	if err := engine.createVersion(&patchedApp); err != nil {
		return nil, newError(
			ErrorTypeUnknown, "application", nil,
			"Version の生成に失敗しました。")
	}

	if reqBody.AllTrafficAvailable != nil && *reqBody.AllTrafficAvailable {
		engine.initTraffic(&patchedApp)
	}

	return &v1.HandlerPatchApplication{
		Id:                     patchedApp.Id,
		Name:                   patchedApp.Name,
		TimeoutSeconds:         patchedApp.TimeoutSeconds,
		Port:                   patchedApp.Port,
		MinScale:               patchedApp.MinScale,
		MaxScale:               patchedApp.MaxScale,
		ScaleTargetConcurrency: patchedApp.ScaleTargetConcurrency,
		Components:             patchedApp.Components,
		Status:                 (v1.HandlerPatchApplicationStatus)(patchedApp.Status),
		PublicUrl:              patchedApp.PublicUrl,
		ResourceId:             patchedApp.ResourceId,
		UpdatedAt:              now,
	}, nil
}

func (engine *Engine) DeleteApplication(id string) error {
	defer engine.lock()()

	// engine.Applications, engine.Versionにデータは残るがここでは省略する
	delete(engine.appVersionRelations, id)
	return nil
}

func (engine *Engine) latestApplication(id string) *v1.Application {
	var app *v1.Application
	if rs, ok := engine.appVersionRelations[id]; ok {
		// 最新のVersionのApplicationを取得
		for i, r := range rs {
			if i == 0 || r.application.CreatedAt.After(app.CreatedAt) {
				app = r.application
			}
		}
	}

	return app
}

func (engine *Engine) newId() (string, error) {
	id, err := uuid.NewRandom()
	return id.String(), err
}

func (engine *Engine) newResourceId() string {
	id := rand.Int32() //nolint:gosec
	return strconv.FormatInt(int64(id), 10)
}
