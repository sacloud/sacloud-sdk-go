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

// コンポーネントの最大CPU数
var ApplicationMaxCPUs = []string{
	(string)(v1.PostApplicationBodyComponentsItemMaxCPU05),
	(string)(v1.PostApplicationBodyComponentsItemMaxCPU1),
	(string)(v1.PostApplicationBodyComponentsItemMaxCPU2),
}

// コンポーネントの最大メモリ
var ApplicationMaxMemories = []string{
	(string)(v1.PostApplicationBodyComponentsItemMaxMemory1Gi),
	(string)(v1.PostApplicationBodyComponentsItemMaxMemory2Gi),
	(string)(v1.PostApplicationBodyComponentsItemMaxMemory4Gi),
}

// ソート順
var ApplicationSortOrders = []string{
	(string)(v1.ListApplicationsSortOrderAsc),
	(string)(v1.ListApplicationsSortOrderDesc),
}

// アプリケーションステータス
var ApplicationStatuses = []string{
	(string)(v1.HandlerGetApplicationOnlyStatusStatusHealthy),
	(string)(v1.HandlerGetApplicationOnlyStatusStatusDeploying),
	(string)(v1.HandlerGetApplicationOnlyStatusStatusUnHealthy),
}

type ApplicationAPI interface {
	// List アプリケーション一覧を取得
	List(ctx context.Context, params *v1.ListApplicationsParams) (*v1.HandlerListApplications, error)
	// Create アプリケーションを作成
	Create(ctx context.Context, params *v1.PostApplicationBody) (*v1.HandlerPostApplication, error)
	// Read アプリケーション詳細を取得
	Read(ctx context.Context, id string) (*v1.HandlerGetApplication, error)
	// Update アプリケーションを部分的に変更
	Update(ctx context.Context, id string, params *v1.PatchApplicationBody) (*v1.HandlerPatchApplication, error)
	// Delete アプリケーションを削除
	Delete(ctx context.Context, id string) error
	// ReadStatus アプリケーションステータスを取得
	ReadStatus(ctx context.Context, id string) (*v1.HandlerGetApplicationOnlyStatus, error)
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

func (op *applicationOp) Create(ctx context.Context, params *v1.PostApplicationBody) (*v1.HandlerPostApplication, error) {
	const methodName = "Applications.Create"
	if params == nil {
		return nil, NewError(methodName, errors.New("params is nil"))
	}

	res, err := op.client.PostApplication(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerPostApplication:
		return result, nil
	case *v1.PostApplicationBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.PostApplicationUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.PostApplicationForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.PostApplicationConflict:
		return nil, apiErrorFromModel(methodName, http.StatusConflict, result)
	case *v1.PostApplicationInternalServerError:
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

func (op *applicationOp) Read(ctx context.Context, id string) (*v1.HandlerGetApplication, error) {
	const methodName = "Applications.Read"
	res, err := op.client.GetApplication(ctx, v1.GetApplicationParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerGetApplication:
		return result, nil
	case *v1.GetApplicationBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.GetApplicationUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.GetApplicationForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.GetApplicationNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.GetApplicationInternalServerError:
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

func (op *applicationOp) ReadStatus(ctx context.Context, id string) (*v1.HandlerGetApplicationOnlyStatus, error) {
	const methodName = "Applications.ReadStatus"
	res, err := op.client.GetApplicationStatus(ctx, v1.GetApplicationStatusParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerGetApplicationOnlyStatus:
		return result, nil
	case *v1.GetApplicationStatusBadRequest:
		return nil, apiErrorFromModel(methodName, http.StatusBadRequest, result)
	case *v1.GetApplicationStatusUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.GetApplicationStatusForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.GetApplicationStatusNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.GetApplicationStatusInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
