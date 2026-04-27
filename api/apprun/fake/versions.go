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
	"sort"
	"time"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

func (engine *Engine) ListVersions(appId string, param v1.ListApplicationVersionsParams) (*v1.HandlerListVersions, error) {
	defer engine.rLock()()

	var versions []*v1.HandlerGetVersion
	for _, r := range engine.appVersionRelations[appId] {
		versions = append(versions, r.version)
	}

	versionsLen := len(versions)
	if versionsLen == 0 {
		return nil, newError(
			ErrorTypeNotFound, "version", nil,
			"アプリケーションが見つかりませんでした。")
	}

	sortField := "created_at"
	if v, ok := param.SortField.Get(); ok {
		sortField = v
	}
	sortOrder := v1.ListApplicationVersionsSortOrderAsc
	if v, ok := param.SortOrder.Get(); ok {
		sortOrder = v
	}
	pageNum := 1
	if v, ok := param.PageNum.Get(); ok {
		pageNum = v
	}
	pageSize := len(versions)
	if v, ok := param.PageSize.Get(); ok {
		pageSize = v
	}

	sort.Slice(versions, func(i int, j int) bool {
		switch sortField {
		case "created_at":
			if sortOrder == v1.ListApplicationVersionsSortOrderDesc {
				return versions[i].CreatedAt.After(versions[j].CreatedAt)
			}
			return versions[i].CreatedAt.Before(versions[j].CreatedAt)
		// sort_fieldのデフォルト値であるcreated_at以外は未サポート
		default:
			return false
		}
	})

	start := (pageNum - 1) * pageSize
	// 範囲外の場合nilを返す
	if start > versionsLen {
		return nil, nil
	}

	end := start + pageSize
	if end > versionsLen {
		end = versionsLen
	}

	var data []v1.HandlerListVersionsDataItem
	for _, v := range versions[start:end] {
		data = append(data, v1.HandlerListVersionsDataItem{
			ID:        v.ID,
			Name:      v.Name,
			Status:    v1.HandlerListVersionsDataItemStatus(v.Status),
			CreatedAt: v.CreatedAt,
		})
	}

	meta := v1.HandlerListVersionsMeta{
		ObjectTotal: versionsLen,
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
		meta.SortOrder = v1.HandlerListVersionsMetaSortOrder(v)
	}

	return &v1.HandlerListVersions{
		Data: data,
		Meta: meta,
	}, nil
}

func (engine *Engine) ReadVersion(appId string, versionId string) (*v1.HandlerGetVersion, error) {
	defer engine.rLock()()

	if _, ok := engine.appVersionRelations[appId]; !ok {
		return nil, newError(
			ErrorTypeNotFound, "version", nil,
			"アプリケーションが見つかりませんでした。")
	}

	var v *v1.HandlerGetVersion
	for _, r := range engine.appVersionRelations[appId] {
		if r.version.ID == versionId {
			v = r.version
			break
		}
	}

	if v == nil {
		return nil, newError(
			ErrorTypeNotFound, "version", nil,
			"アプリケーションが見つかりませんでした。")
	}

	return v, nil
}

func (engine *Engine) ReadVersionStatus(appId string, versionId string) (*v1.HandlerGetApplicationVersionOnlyStatus, error) {
	defer engine.rLock()()

	if _, ok := engine.appVersionRelations[appId]; !ok {
		return nil, newError(
			ErrorTypeNotFound, "version", nil,
			"アプリケーションが見つかりませんでした。")
	}

	for _, r := range engine.appVersionRelations[appId] {
		if r.version.ID == versionId {
			return &v1.HandlerGetApplicationVersionOnlyStatus{
				Status:  v1.HandlerGetApplicationVersionOnlyStatusStatus(r.version.Status),
				Message: "",
			}, nil
		}
	}

	return nil, newError(
		ErrorTypeNotFound, "version", nil,
		"アプリケーションが見つかりませんでした。")
}

func (engine *Engine) DeleteVersion(appId string, versionId string) error {
	defer engine.lock()()

	if _, ok := engine.appVersionRelations[appId]; !ok {
		return newError(
			ErrorTypeNotFound, "version", nil,
			"アプリケーションが見つかりませんでした。")
	}

	var idx int
	rs := engine.appVersionRelations[appId]
	for i, r := range rs {
		if r.version.ID == versionId {
			idx = i
			break
		}
	}

	rs[idx] = rs[len(rs)-1]
	rs = rs[:len(rs)-1]
	engine.appVersionRelations[appId] = rs

	return nil
}

func (engine *Engine) createVersion(app *v1.HandlerGetApplication) error {
	versionId, err := engine.newId()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("version-%03d", engine.nextVersionId())
	createdAt := time.Now().UTC().Truncate(time.Second)
	components := convertVersionComponents(app.Components)

	v := &v1.HandlerGetVersion{
		ID:                     versionId,
		Name:                   name,
		Status:                 v1.HandlerGetVersionStatus(app.Status),
		TimeoutSeconds:         app.TimeoutSeconds,
		Port:                   app.Port,
		MinScale:               app.MinScale,
		MaxScale:               app.MaxScale,
		ScaleTargetConcurrency: app.ScaleTargetConcurrency,
		Components:             components,
		CreatedAt:              createdAt,
	}
	engine.Versions = append(engine.Versions, v)

	// 内部的にVersionとApplicationのリレーションを保持する
	engine.appVersionRelations[app.ID] = append(engine.appVersionRelations[app.ID],
		&appVersionRelation{
			application: app,
			version:     v,
		},
	)

	return nil
}

func convertVersionComponents(components []v1.HandlerGetApplicationComponentsItem) []v1.HandlerGetVersionComponentsItem {
	if len(components) == 0 {
		return nil
	}
	payload, err := json.Marshal(components)
	if err != nil {
		return nil
	}
	var out []v1.HandlerGetVersionComponentsItem
	if err := json.Unmarshal(payload, &out); err != nil {
		return nil
	}
	return out
}
