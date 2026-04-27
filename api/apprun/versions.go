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

package apprun

import (
	"context"
	"errors"
	"net/http"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

// ソート順
var VersionSortOrders = []string{
	(string)(v1.ListApplicationVersionsSortOrderAsc),
	(string)(v1.ListApplicationVersionsSortOrderDesc),
}

// バージョンステータス
var VersionStatuses = []string{
	(string)(v1.HandlerGetApplicationVersionOnlyStatusStatusHealthy),
	(string)(v1.HandlerGetApplicationVersionOnlyStatusStatusDeploying),
	(string)(v1.HandlerGetApplicationVersionOnlyStatusStatusUnHealthy),
}

type VersionAPI interface {
	// List アプリケーションバージョン一覧を取得
	List(ctx context.Context, appId string, params *v1.ListApplicationVersionsParams) (*v1.HandlerListVersions, error)
	// Read アプリケーションバージョン詳細を取得
	Read(ctx context.Context, appId, versionId string) (*v1.HandlerGetVersion, error)
	// Delete アプリケーションバージョンを削除
	Delete(ctx context.Context, appId, versionId string) error

	ReadStatus(ctx context.Context, appId, versionId string) (*v1.HandlerGetApplicationVersionOnlyStatus, error)
}

var _ VersionAPI = (*versionOp)(nil)

type versionOp struct {
	client *v1.Client
}

// NewVersionOp アプリケーションバージョン操作関連API
func NewVersionOp(client *v1.Client) VersionAPI {
	return &versionOp{client: client}
}

func (op *versionOp) List(ctx context.Context, appId string, params *v1.ListApplicationVersionsParams) (*v1.HandlerListVersions, error) {
	reqParams := v1.ListApplicationVersionsParams{ID: appId}
	if params != nil {
		reqParams = *params
		reqParams.ID = appId
	}
	const methodName = "Versions.List"
	res, err := op.client.ListApplicationVersions(ctx, reqParams)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerListVersions:
		return result, nil
	case *v1.ListApplicationVersionsBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.ListApplicationVersionsUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.ListApplicationVersionsForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.ListApplicationVersionsNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.ListApplicationVersionsInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *versionOp) Read(ctx context.Context, appId, versionId string) (*v1.HandlerGetVersion, error) {
	const methodName = "Versions.Read"
	res, err := op.client.GetApplicationVersion(ctx, v1.GetApplicationVersionParams{ID: appId, VersionID: versionId})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerGetVersion:
		return result, nil
	case *v1.GetApplicationVersionBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.GetApplicationVersionUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.GetApplicationVersionForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.GetApplicationVersionNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.GetApplicationVersionInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *versionOp) Delete(ctx context.Context, appId, versionId string) error {
	const methodName = "Versions.Delete"
	res, err := op.client.DeleteApplicationVersion(ctx, v1.DeleteApplicationVersionParams{ID: appId, VersionID: versionId})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.DeleteApplicationVersionNoContent:
		return nil
	case *v1.DeleteApplicationVersionBadRequest:
		return apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.DeleteApplicationVersionUnauthorized:
		return apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.DeleteApplicationVersionForbidden:
		return apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.DeleteApplicationVersionNotFound:
		return apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.DeleteApplicationVersionInternalServerError:
		return apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *versionOp) ReadStatus(ctx context.Context, appId, versionId string) (*v1.HandlerGetApplicationVersionOnlyStatus, error) {
	const methodName = "Versions.ReadStatus"
	res, err := op.client.GetApplicationVersionStatus(ctx, v1.GetApplicationVersionStatusParams{ID: appId, VersionID: versionId})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerGetApplicationVersionOnlyStatus:
		return result, nil
	case *v1.GetApplicationVersionStatusBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.GetApplicationVersionStatusUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.GetApplicationVersionStatusForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.GetApplicationVersionStatusNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.GetApplicationVersionStatusInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
