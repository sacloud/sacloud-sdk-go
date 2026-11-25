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

	return errorFromDecodedResponse[v1.HandlerListApplications]("Applications.List", func() (any, error) {
		return op.client.ListApplications(ctx, reqParams)
	})
}

func (op *applicationOp) Create(ctx context.Context, params *v1.CreateApplicationBody) (*v1.HandlerCreateApplication, error) {
	const methodName = "Applications.Create"
	if params == nil {
		return nil, NewError(methodName, errors.New("params is nil"))
	}

	return errorFromDecodedResponse[v1.HandlerCreateApplication](methodName, func() (any, error) {
		return op.client.CreateApplication(ctx, params)
	})
}

func (op *applicationOp) Update(ctx context.Context, id string, params *v1.PatchApplicationBody) (*v1.HandlerPatchApplication, error) {
	const methodName = "Applications.Update"
	if params == nil {
		return nil, NewError(methodName, errors.New("params is nil"))
	}

	return errorFromDecodedResponse[v1.HandlerPatchApplication](methodName, func() (any, error) {
		return op.client.PatchApplication(ctx, params, v1.PatchApplicationParams{ID: id})
	})
}

func (op *applicationOp) Read(ctx context.Context, id string) (*v1.HandlerReadApplication, error) {
	return errorFromDecodedResponse[v1.HandlerReadApplication]("Applications.Read", func() (any, error) {
		return op.client.ReadApplication(ctx, v1.ReadApplicationParams{ID: id})
	})
}

func (op *applicationOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[v1.DeleteApplicationNoContent]("Applications.Delete", func() (any, error) {
		return op.client.DeleteApplication(ctx, v1.DeleteApplicationParams{ID: id})
	})
}

func (op *applicationOp) ReadStatus(ctx context.Context, id string) (*v1.HandlerReadApplicationOnlyStatus, error) {
	return errorFromDecodedResponse[v1.HandlerReadApplicationOnlyStatus]("Applications.ReadStatus", func() (any, error) {
		return op.client.ReadApplicationStatus(ctx, v1.ReadApplicationStatusParams{ID: id})
	})
}
