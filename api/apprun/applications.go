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

package apprun

import (
	"context"
	"encoding/json"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

// コンポーネントの最大CPU数
var ApplicationMaxCPUs = []string{
	(string)(v1.PostApplicationBodyComponentMaxCpuN05),
	(string)(v1.PostApplicationBodyComponentMaxCpuN1),
	(string)(v1.PostApplicationBodyComponentMaxCpuN2),
}

// コンポーネントの最大メモリ
var ApplicationMaxMemories = []string{
	(string)(v1.PostApplicationBodyComponentMaxMemoryN1Gi),
	(string)(v1.PostApplicationBodyComponentMaxMemoryN2Gi),
	(string)(v1.PostApplicationBodyComponentMaxMemoryN4Gi),
}

// ソート順
var ApplicationSortOrders = []string{
	(string)(v1.ListApplicationsParamsSortOrderAsc),
	(string)(v1.ListApplicationsParamsSortOrderDesc),
}

// アプリケーションステータス
var ApplicationStatuses = []string{
	(string)(v1.ApplicationStatusHealthy),
	(string)(v1.ApplicationStatusDeploying),
	(string)(v1.ApplicationStatusUnHealthy),
}

type ApplicationAPI interface {
	// List アプリケーション一覧を取得
	List(ctx context.Context, params *v1.ListApplicationsParams) (*v1.HandlerListApplications, error)
	// Create アプリケーションを作成
	Create(ctx context.Context, params *v1.PostApplicationBody) (*v1.Application, error)
	// Read アプリケーション詳細を取得
	Read(ctx context.Context, id string) (*v1.Application, error)
	// Update アプリケーションを部分的に変更
	Update(ctx context.Context, id string, params *v1.PatchApplicationBody) (*v1.HandlerPatchApplication, error)
	// Delete アプリケーションを削除
	Delete(ctx context.Context, id string) error
	// ReadStatus アプリケーションステータスを取得
	ReadStatus(ctx context.Context, id string) (*v1.HandlerGetApplicationOnlyStatus, error)
}

var _ ApplicationAPI = (*applicationOp)(nil)

type applicationOp struct {
	client *Client
}

// NewApplicationOp アプリケーション操作関連API
func NewApplicationOp(client *Client) ApplicationAPI {
	return &applicationOp{client: client}
}

func (op *applicationOp) List(ctx context.Context, params *v1.ListApplicationsParams) (*v1.HandlerListApplications, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.ListApplicationsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	applications, err := resp.Result()
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (op *applicationOp) Create(ctx context.Context, params *v1.PostApplicationBody) (*v1.Application, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.PostApplicationWithResponse(ctx, *params)
	if err != nil {
		return nil, err
	}
	created, err := resp.Result()
	if err != nil {
		return nil, err
	}
	// Convert create response to a stable application shape.
	b, err := json.Marshal(created)
	if err != nil {
		return nil, err
	}
	var application v1.Application
	if err := json.Unmarshal(b, &application); err != nil {
		return nil, err
	}
	return &application, nil
}

func (op *applicationOp) Update(ctx context.Context, id string, params *v1.PatchApplicationBody) (*v1.HandlerPatchApplication, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.PatchApplicationWithResponse(ctx, id, *params)
	if err != nil {
		return nil, err
	}
	patchApplication, err := resp.Result()
	if err != nil {
		return nil, err
	}
	return patchApplication, nil
}

func (op *applicationOp) Read(ctx context.Context, id string) (*v1.Application, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.GetApplicationWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}
	application, err := resp.Result()
	if err != nil {
		return nil, err
	}
	return application, nil
}

func (op *applicationOp) Delete(ctx context.Context, id string) error {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return err
	}
	resp, err := apiClient.DeleteApplicationWithResponse(ctx, id)
	if err != nil {
		return err
	}
	return resp.Result()
}

func (op *applicationOp) ReadStatus(ctx context.Context, id string) (*v1.HandlerGetApplicationOnlyStatus, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.GetApplicationStatusWithResponse(ctx, id)
	if err != nil {
		return nil, err
	}
	applicationStatus, err := resp.Result()
	if err != nil {
		return nil, err
	}
	return applicationStatus, nil
}
