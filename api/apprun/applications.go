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

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun/apis/v1"
)

// コンポーネントの最大CPU数
var ApplicationMaxCPUs = []string{
	(string)(v1.CreateApplicationBodyComponentsItemMaxCPU05),
	(string)(v1.CreateApplicationBodyComponentsItemMaxCPU1),
	(string)(v1.CreateApplicationBodyComponentsItemMaxCPU2),
}

// コンポーネントの最大メモリ
var ApplicationMaxMemories = []string{
	(string)(v1.CreateApplicationBodyComponentsItemMaxMemory1Gi),
	(string)(v1.CreateApplicationBodyComponentsItemMaxMemory2Gi),
	(string)(v1.CreateApplicationBodyComponentsItemMaxMemory4Gi),
}

// ソート順
var ApplicationSortOrders = []string{
	(string)(v1.ListApplicationsSortOrderAsc),
	(string)(v1.ListApplicationsSortOrderDesc),
}

// アプリケーションステータス
var ApplicationStatuses = []string{
	(string)(v1.HandlerReadApplicationOnlyStatusStatusHealthy),
	(string)(v1.HandlerReadApplicationOnlyStatusStatusDeploying),
	(string)(v1.HandlerReadApplicationOnlyStatusStatusUnHealthy),
}

type ApplicationAPI interface {
	// List アプリケーション一覧を取得
	List(ctx context.Context, params *v1.ListApplicationsParams) (*v1.HandlerListApplications, error)
	// Create アプリケーションを作成
	Create(ctx context.Context, params *v1.CreateApplicationBody) (*v1.HandlerCreateApplication, error)
	// Read アプリケーション詳細を取得
	Read(ctx context.Context, id string) (*v1.HandlerReadApplication, error)
	// Update アプリケーションを部分的に変更
	Update(ctx context.Context, id string, params *v1.PatchApplicationBody) (*v1.HandlerPatchApplication, error)
	// Delete アプリケーションを削除
	Delete(ctx context.Context, id string) error
	// ReadStatus アプリケーションステータスを取得
	ReadStatus(ctx context.Context, id string) (*v1.HandlerReadApplicationOnlyStatus, error)
}

var _ ApplicationAPI = (*applicationOp)(nil)

type applicationOp struct {
	client *v1.Client
}

// NewApplicationOp アプリケーション操作関連API
func NewApplicationOp(client *v1.Client) ApplicationAPI {
	return &applicationOp{client: client}
}

func (op *applicationOp) List(ctx context.Context, params *v1.ListApplicationsParams) (*v1.HandlerListApplications, error) {
	reqParams := v1.ListApplicationsParams{}
	if params != nil {
		reqParams = *params
	}
	const methodName = "Applications.List"
	res, err := op.client.ListApplications(ctx, reqParams)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerListApplications:
		return result, nil
	case *v1.ListApplicationsBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.ListApplicationsUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.ListApplicationsForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.ListApplicationsInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *applicationOp) Create(ctx context.Context, params *v1.CreateApplicationBody) (*v1.HandlerCreateApplication, error) {
	const methodName = "Applications.Create"
	if params == nil {
		return nil, NewError(methodName, errors.New("params is nil"))
	}

	res, err := op.client.CreateApplication(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerCreateApplication:
		return result, nil
	case *v1.CreateApplicationBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.CreateApplicationUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.CreateApplicationForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.CreateApplicationConflict:
		return nil, apiErrorFromModel(methodName, http.StatusConflict, result)
	case *v1.CreateApplicationInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *applicationOp) Update(ctx context.Context, id string, params *v1.PatchApplicationBody) (*v1.HandlerPatchApplication, error) {
	const methodName = "Applications.Update"
	if params == nil {
		return nil, NewError(methodName, errors.New("params is nil"))
	}

	res, err := op.client.PatchApplication(ctx, params, v1.PatchApplicationParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerPatchApplication:
		return result, nil
	case *v1.PatchApplicationBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.PatchApplicationUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.PatchApplicationForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.PatchApplicationNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.PatchApplicationConflict:
		return nil, apiErrorFromModel(methodName, http.StatusConflict, result)
	case *v1.PatchApplicationInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *applicationOp) Read(ctx context.Context, id string) (*v1.HandlerReadApplication, error) {
	const methodName = "Applications.Read"
	res, err := op.client.ReadApplication(ctx, v1.ReadApplicationParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerReadApplication:
		return result, nil
	case *v1.ReadApplicationBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.ReadApplicationUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.ReadApplicationForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.ReadApplicationNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.ReadApplicationInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *applicationOp) Delete(ctx context.Context, id string) error {
	const methodName = "Applications.Delete"
	res, err := op.client.DeleteApplication(ctx, v1.DeleteApplicationParams{ID: id})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.DeleteApplicationNoContent:
		return nil
	case *v1.DeleteApplicationBadRequest:
		return apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.DeleteApplicationUnauthorized:
		return apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.DeleteApplicationForbidden:
		return apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.DeleteApplicationNotFound:
		return apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.DeleteApplicationInternalServerError:
		return apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *applicationOp) ReadStatus(ctx context.Context, id string) (*v1.HandlerReadApplicationOnlyStatus, error) {
	const methodName = "Applications.ReadStatus"
	res, err := op.client.ReadApplicationStatus(ctx, v1.ReadApplicationStatusParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerReadApplicationOnlyStatus:
		return result, nil
	case *v1.ReadApplicationStatusBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.ReadApplicationStatusUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.ReadApplicationStatusForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.ReadApplicationStatusNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.ReadApplicationStatusInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
