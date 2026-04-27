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
	var apps []*v1.HandlerGetApplication
	for id := range engine.appVersionRelations {
		apps = append(apps, engine.latestApplication(id))
	}

	sortField := "created_at"
	if v, ok := param.SortField.Get(); ok {
		sortField = v
	}
	sortOrder := v1.ListApplicationsSortOrderAsc
	if v, ok := param.SortOrder.Get(); ok {
		sortOrder = v
	}
	pageNum := 1
	if v, ok := param.PageNum.Get(); ok {
		pageNum = v
	}
	pageSize := len(apps)
	if v, ok := param.PageSize.Get(); ok {
		pageSize = v
	}

	sort.Slice(apps, func(i int, j int) bool {
		switch sortField {
		case "created_at":
			if sortOrder == v1.ListApplicationsSortOrderDesc {
				return apps[i].CreatedAt.After(apps[j].CreatedAt)
			}
			return apps[i].CreatedAt.Before(apps[j].CreatedAt)
		// sort_fieldのデフォルト値であるcreated_at以外は未サポート
		default:
			return false
		}
	})

	appsLen := len(apps)
	start := (pageNum - 1) * pageSize
	// 範囲外の場合nilを返す
	if start > appsLen {
		return nil, nil
	}

	end := start + pageSize
	if end > appsLen {
		end = appsLen
	}

	var data []v1.HandlerListApplicationsDataItem
	for _, app := range apps[start:end] {
		if app != nil {
			data = append(data, v1.HandlerListApplicationsDataItem{
				ID:        app.ID,
				Name:      app.Name,
				Status:    v1.HandlerListApplicationsDataItemStatus(app.Status),
				PublicURL: app.PublicURL,
				CreatedAt: app.CreatedAt,
			})
		}
	}

	meta := v1.HandlerListApplicationsMeta{
		ObjectTotal: appsLen,
	}
	if v, ok := param.PageNum.Get(); ok {
		meta.PageNum = v
	}
	if v, ok := param.PageSize.Get(); ok {
		meta.PageSize = v
	}
	if v, ok := param.SortField.Get(); ok {
		meta.SortField = v
	}
	if v, ok := param.SortOrder.Get(); ok {
		meta.SortOrder = v1.HandlerListApplicationsMetaSortOrder(v)
	}

	return &v1.HandlerListApplications{
		Data: data,
		Meta: meta,
	}, nil
}

func (engine *Engine) CreateApplication(reqBody *v1.PostApplicationBody) (*v1.HandlerPostApplication, error) {
	defer engine.lock()()

	appID, err := engine.newId()
	if err != nil {
		return nil, newError(
			ErrorTypeUnknown, "application", nil,
			"Application IDの生成に失敗しました。")
	}

	components := convertPostComponents(reqBody.Components)
	status := v1.HandlerGetApplicationStatusHealthy
	url := fmt.Sprintf("https://example.com/apprun/dummy/%s", appID)
	createdAt := time.Now().UTC().Truncate(time.Second)
	app := &v1.HandlerGetApplication{
		ID:                     appID,
		Name:                   reqBody.Name,
		TimeoutSeconds:         reqBody.TimeoutSeconds,
		Port:                   reqBody.Port,
		MinScale:               reqBody.MinScale,
		MaxScale:               reqBody.MaxScale,
		ScaleTargetConcurrency: reqBody.ScaleTargetConcurrency,
		Components:             components,
		Status:                 status,
		PublicURL:              url,
		ResourceID:             engine.newResourceId(),
		CreatedAt:              createdAt,
	}
	engine.Applications = append(engine.Applications, app)

	if err := engine.createVersion(app); err != nil {
		return nil, newError(
			ErrorTypeUnknown, "application", nil,
			"Version の生成に失敗しました。")
	}

	engine.initTraffic(app)

	return toHandlerPostApplication(app), nil
}

func (engine *Engine) ReadApplication(id string) (*v1.HandlerGetApplication, error) {
	defer engine.rLock()()

	if len(engine.Applications) == 0 {
		return nil, newError(
			ErrorTypeNotFound, "application", nil,
			"アプリケーションが見つかりませんでした。")
	}

	app := engine.latestApplication(id)
	if app != nil && app.ID == id {
		return app, nil
	}

	return nil, newError(
		ErrorTypeNotFound, "application", nil,
		"アプリケーションが見つかりませんでした。")
}

func (engine *Engine) UpdateApplication(id string, reqBody *v1.PatchApplicationBody) (*v1.HandlerPatchApplication, error) {
	defer engine.lock()()

	patchedApp := *(engine.latestApplication(id))
	if v, ok := reqBody.TimeoutSeconds.Get(); ok {
		patchedApp.TimeoutSeconds = v
	}
	if v, ok := reqBody.Port.Get(); ok {
		patchedApp.Port = v
	}
	if v, ok := reqBody.MinScale.Get(); ok {
		patchedApp.MinScale = v
	}
	if v, ok := reqBody.MaxScale.Get(); ok {
		patchedApp.MaxScale = v
	}
	if v, ok := reqBody.ScaleTargetConcurrency.Get(); ok {
		patchedApp.ScaleTargetConcurrency = v1.NewOptInt(v)
	}
	if len(reqBody.Components) > 0 {
		patchedApp.Components = convertPatchComponents(reqBody.Components)
	}

	now := time.Now().UTC().Truncate(time.Second)
	patchedApp.CreatedAt = now

	engine.Applications = append(engine.Applications, &patchedApp)
	if err := engine.createVersion(&patchedApp); err != nil {
		return nil, newError(
			ErrorTypeUnknown, "application", nil,
			"Version の生成に失敗しました。")
	}

	if v, ok := reqBody.AllTrafficAvailable.Get(); ok && v {
		engine.initTraffic(&patchedApp)
	}

	return toHandlerPatchApplication(&patchedApp, now), nil
}

func (engine *Engine) DeleteApplication(id string) error {
	defer engine.lock()()

	// engine.Applications, engine.Versionにデータは残るがここでは省略する
	delete(engine.appVersionRelations, id)
	return nil
}

func (engine *Engine) latestApplication(id string) *v1.HandlerGetApplication {
	var app *v1.HandlerGetApplication
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

func convertPostComponents(components []v1.PostApplicationBodyComponentsItem) []v1.HandlerGetApplicationComponentsItem {
	if len(components) == 0 {
		return nil
	}
	payload, err := json.Marshal(components)
	if err != nil {
		return nil
	}
	var out []v1.HandlerGetApplicationComponentsItem
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil
	}
	return out
}

func convertPatchComponents(components []v1.PatchApplicationBodyComponentsItem) []v1.HandlerGetApplicationComponentsItem {
	if len(components) == 0 {
		return nil
	}
	payload, err := json.Marshal(components)
	if err != nil {
		return nil
	}
	var out []v1.HandlerGetApplicationComponentsItem
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil
	}
	return out
}

func toHandlerPostApplication(app *v1.HandlerGetApplication) *v1.HandlerPostApplication {
	if app == nil {
		return nil
	}
	payload, err := json.Marshal(app)
	if err != nil {
		return nil
	}
	var out v1.HandlerPostApplication
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil
	}
	out.Status = v1.HandlerPostApplicationStatus(app.Status)
	return &out
}

func toHandlerPatchApplication(app *v1.HandlerGetApplication, updatedAt time.Time) *v1.HandlerPatchApplication {
	if app == nil {
		return nil
	}
	payload, err := json.Marshal(app)
	if err != nil {
		return nil
	}

	// HandlerPatchApplicationはupdated_atをrequiredフィールドとして持つため、payloadからcreated_atを削除し、updated_atを追加する
	tempJson := make(map[string]interface{})
	if err := json.Unmarshal(payload, &tempJson); err != nil {
		return nil
	}
	delete(tempJson, "created_at")
	tempJson["updated_at"] = updatedAt.Format(time.RFC3339)
	payload, err = json.Marshal(tempJson)
	if err != nil {
		return nil
	}

	var out v1.HandlerPatchApplication
	if err := json.Unmarshal(payload, &out); err != nil {
		fmt.Printf("patched error: %s\n", err)
		return nil
	}
	out.Status = v1.HandlerPatchApplicationStatus(app.Status)
	return &out
}

func (engine *Engine) newId() (string, error) {
	id, err := uuid.NewRandom()
	return id.String(), err
}

func (engine *Engine) newResourceId() string {
	id := rand.Int32() //nolint:gosec
	return strconv.FormatInt(int64(id), 10)
}
