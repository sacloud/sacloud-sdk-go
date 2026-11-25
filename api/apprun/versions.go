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

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun/apis/v1"
)

// ソート順
var VersionSortOrders = []string{
	(string)(v1.ListApplicationVersionsSortOrderAsc),
	(string)(v1.ListApplicationVersionsSortOrderDesc),
}

// バージョンステータス
var VersionStatuses = []string{
	(string)(v1.HandlerReadApplicationVersionOnlyStatusStatusHealthy),
	(string)(v1.HandlerReadApplicationVersionOnlyStatusStatusDeploying),
	(string)(v1.HandlerReadApplicationVersionOnlyStatusStatusUnHealthy),
}

type VersionAPI interface {
	// List アプリケーションバージョン一覧を取得
	List(ctx context.Context, appId string, params *v1.ListApplicationVersionsParams) (*v1.HandlerListVersions, error)
	// Read アプリケーションバージョン詳細を取得
	Read(ctx context.Context, appId, versionId string) (*v1.HandlerReadVersion, error)
	// Delete アプリケーションバージョンを削除
	Delete(ctx context.Context, appId, versionId string) error

	ReadStatus(ctx context.Context, appId, versionId string) (*v1.HandlerReadApplicationVersionOnlyStatus, error)
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
	return errorFromDecodedResponse[v1.HandlerListVersions]("Versions.List", func() (any, error) {
		return op.client.ListApplicationVersions(ctx, reqParams)
	})
}

func (op *versionOp) Read(ctx context.Context, appId, versionId string) (*v1.HandlerReadVersion, error) {
	return errorFromDecodedResponse[v1.HandlerReadVersion]("Versions.Read", func() (any, error) {
		return op.client.ReadApplicationVersion(ctx, v1.ReadApplicationVersionParams{ID: appId, VersionID: versionId})
	})
}

func (op *versionOp) Delete(ctx context.Context, appId, versionId string) error {
	return errorFromDecodedResponseE[v1.DeleteApplicationVersionNoContent]("Versions.Delete", func() (any, error) {
		return op.client.DeleteApplicationVersion(ctx, v1.DeleteApplicationVersionParams{ID: appId, VersionID: versionId})
	})
}

func (op *versionOp) ReadStatus(ctx context.Context, appId, versionId string) (*v1.HandlerReadApplicationVersionOnlyStatus, error) {
	return errorFromDecodedResponse[v1.HandlerReadApplicationVersionOnlyStatus]("Versions.ReadStatus", func() (any, error) {
		return op.client.ReadApplicationVersionStatus(ctx, v1.ReadApplicationVersionStatusParams{ID: appId, VersionID: versionId})
	})
}
