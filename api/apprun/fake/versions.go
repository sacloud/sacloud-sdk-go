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
	"sort"
	"time"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

func (engine *Engine) ListVersions(appId string, param v1.ListApplicationVersionsParams) (*v1.HandlerListVersions, error) {
	defer engine.rLock()()

	var versions []*v1.Version
	for _, r := range engine.appVersionRelations[appId] {
		versions = append(versions, r.version)
	}

	versionsLen := len(versions)
	if versionsLen == 0 {
		return nil, newError(
			ErrorTypeNotFound, "version", nil,
			"アプリケーションが見つかりませんでした。")
	}

	sort.Slice(versions, func(i int, j int) bool {
		switch *param.SortField {
		case "created_at":
			if *param.SortOrder == "desc" {
				return versions[i].CreatedAt.After(versions[j].CreatedAt)
			}
			return versions[i].CreatedAt.Before(versions[j].CreatedAt)
		// sort_fieldのデフォルト値であるcreated_at以外は未サポート
		default:
			return false
		}
	})

	start := (*param.PageNum - 1) * *param.PageSize
	// 範囲外の場合nilを返す
	if start > versionsLen {
		return nil, nil
	}

	end := start + *param.PageSize
	if end > versionsLen {
		end = versionsLen
	}

	var data []v1.Version
	for _, v := range versions[start:end] {
		data = append(data, *v)
	}

	meta := v1.HandlerListVersionsMeta{
		ObjectTotal: versionsLen,
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
		so := (*v1.HandlerListVersionsMetaSortOrder)(param.SortOrder)
		meta.SortOrder = *so
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

	var v v1.HandlerGetVersion
	for _, r := range engine.appVersionRelations[appId] {
		if r.version.Id == versionId {
			v.Id = r.version.Id
			v.Name = r.version.Name
			v.Status = (v1.HandlerGetVersionStatus)(r.version.Status)
			v.TimeoutSeconds = r.application.TimeoutSeconds
			v.Port = r.application.Port
			v.MinScale = r.application.MinScale
			v.MaxScale = r.application.MaxScale
			v.Components = r.application.Components
			v.CreatedAt = r.application.CreatedAt
		}
	}

	return &v, nil
}

func (engine *Engine) ReadVersionStatus(appId string, versionId string) (*v1.HandlerGetApplicationVersionOnlyStatus, error) {
	defer engine.rLock()()

	if _, ok := engine.appVersionRelations[appId]; !ok {
		return nil, newError(
			ErrorTypeNotFound, "version", nil,
			"アプリケーションが見つかりませんでした。")
	}

	for _, r := range engine.appVersionRelations[appId] {
		if r.version.Id == versionId {
			return &v1.HandlerGetApplicationVersionOnlyStatus{
				Status:  v1.HandlerGetVersionStatusStatus(r.version.Status),
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
		if r.version.Id == versionId {
			idx = i
			break
		}
	}

	rs[idx] = rs[len(rs)-1]
	rs = rs[:len(rs)-1]
	engine.appVersionRelations[appId] = rs

	return nil
}

func (engine *Engine) createVersion(app *v1.Application) error {
	versionId, err := engine.newId()
	if err != nil {
		return err
	}
	name := fmt.Sprintf("version-%03d", engine.nextVersionId())
	createdAt := time.Now().UTC().Truncate(time.Second)

	v := v1.Version{
		Id:        versionId,
		Name:      name,
		Status:    (v1.VersionStatus)(app.Status),
		CreatedAt: createdAt,
	}
	engine.Versions = append(engine.Versions, &v)

	// 内部的にVersionとApplicationのリレーションを保持する
	engine.appVersionRelations[app.Id] = append(engine.appVersionRelations[app.Id],
		&appVersionRelation{
			application: app,
			version:     &v,
		},
	)

	return nil
}
