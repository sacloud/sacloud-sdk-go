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

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

// ソート順
var VersionSortOrders = []string{
	(string)(v1.ListApplicationVersionsParamsSortOrderAsc),
	(string)(v1.ListApplicationVersionsParamsSortOrderDesc),
}

// バージョンステータス
var VersionStatuses = []string{
	(string)(v1.HandlerGetVersionStatusHealthy),
	(string)(v1.HandlerGetVersionStatusDeploying),
	(string)(v1.HandlerGetVersionStatusUnHealthy),
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
	client *Client
}

// NewVersionOp アプリケーションバージョン操作関連API
func NewVersionOp(client *Client) VersionAPI {
	return &versionOp{client: client}
}

func (op *versionOp) List(ctx context.Context, appId string, params *v1.ListApplicationVersionsParams) (*v1.HandlerListVersions, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.ListApplicationVersionsWithResponse(ctx, appId, params)
	if err != nil {
		return nil, err
	}
	versions, err := resp.Result()
	if err != nil {
		return nil, err
	}
	return versions, nil
}

func (op *versionOp) Read(ctx context.Context, appId, versionId string) (*v1.HandlerGetVersion, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.GetApplicationVersionWithResponse(ctx, appId, versionId)
	if err != nil {
		return nil, err
	}
	version, err := resp.Result()
	if err != nil {
		return nil, err
	}
	return version, nil
}

func (op *versionOp) Delete(ctx context.Context, appId, versionId string) error {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return err
	}
	resp, err := apiClient.DeleteApplicationVersionWithResponse(ctx, appId, versionId)
	if err != nil {
		return err
	}
	return resp.Result()
}

func (op *versionOp) ReadStatus(ctx context.Context, appId, versionId string) (*v1.HandlerGetApplicationVersionOnlyStatus, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.GetApplicationVersionStatusWithResponse(ctx, appId, versionId)
	if err != nil {
		return nil, err
	}
	status, err := resp.Result()
	if err != nil {
		return nil, err
	}
	return status, nil
}
