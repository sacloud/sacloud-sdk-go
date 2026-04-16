// Copyright 2025- The sacloud/iam-api-go Authors
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

package scim

import (
	"context"

	"github.com/google/uuid"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

// ScimAPI SCIM API
type ScimAPI interface {
	// List ユーザープロビジョニング一覧を取得する
	List(ctx context.Context, params ListParams) (*v1.ScimConfigurationsGetOK, error)
	// Create ユーザープロビジョニングを作成する
	Create(ctx context.Context, params CreateParams) (*v1.ScimConfiguration, error)
	// Read ユーザープロビジョニングを取得する
	Read(ctx context.Context, id string) (*v1.ScimConfigurationBase, error)
	// Update ユーザープロビジョニングを更新する
	Update(ctx context.Context, id string, params UpdateParams) (*v1.ScimConfigurationBase, error)
	// Delete ユーザープロビジョニングを削除する
	Delete(ctx context.Context, id string) error
	// RegenerateToken ユーザープロビジョニングのシークレットトークンを再発行する
	RegenerateToken(ctx context.Context, id string) (*v1.ScimConfigurationsIDRegenerateTokenPostOK, error)
}

// scimOp SCIM APIの実装
type scimOp struct {
	client *v1.Client
}

// NewScimOp SCIM APIのコンストラクタ
func NewScimOp(client *v1.Client) ScimAPI {
	return &scimOp{client: client}
}

// ListParams ユーザープロビジョニング一覧取得パラメータ
type ListParams struct {
	Page    *int `json:"page,omitempty"`
	PerPage *int `json:"per_page,omitempty"`
}

// CreateParams ユーザープロビジョニング作成パラメータ
type CreateParams struct {
	Name string `json:"name"`
}

// UpdateParams ユーザープロビジョニング更新パラメータ
type UpdateParams struct {
	Name string `json:"name"`
}

// List ユーザープロビジョニング一覧を取得する
func (s *scimOp) List(ctx context.Context, params ListParams) (*v1.ScimConfigurationsGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.ScimConfigurationsGetOK]("Scim.List", func() (any, error) {
		return s.client.ScimConfigurationsGet(ctx, v1.ScimConfigurationsGetParams{
			Page:    common.IntoOpt[v1.OptInt](params.Page),
			PerPage: common.IntoOpt[v1.OptInt](params.PerPage),
		})
	})
}

// Create ユーザープロビジョニングを作成する
func (s *scimOp) Create(ctx context.Context, params CreateParams) (*v1.ScimConfiguration, error) {
	return common.ErrorFromDecodedResponse[v1.ScimConfiguration]("Scim.Create", func() (any, error) {
		return s.client.ScimConfigurationsPost(ctx, &v1.ScimConfigurationsPostReq{
			Name: params.Name,
		})
	})
}

// Read ユーザープロビジョニングを取得する
func (s *scimOp) Read(ctx context.Context, id string) (*v1.ScimConfigurationBase, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return common.ErrorFromDecodedResponse[v1.ScimConfigurationBase]("Scim.Read", func() (any, error) {
		return s.client.ScimConfigurationsIDGet(ctx, v1.ScimConfigurationsIDGetParams{
			ID: uuid,
		})
	})
}

// Update ユーザープロビジョニングを更新する
func (s *scimOp) Update(ctx context.Context, id string, params UpdateParams) (*v1.ScimConfigurationBase, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return common.ErrorFromDecodedResponse[v1.ScimConfigurationBase]("Scim.Update", func() (any, error) {
		return s.client.ScimConfigurationsIDPut(ctx, &v1.ScimConfigurationsIDPutReq{
			Name: params.Name,
		}, v1.ScimConfigurationsIDPutParams{
			ID: uuid,
		})
	})
}

// Delete ユーザープロビジョニングを削除する
func (s *scimOp) Delete(ctx context.Context, id string) error {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	_, err = common.ErrorFromDecodedResponse[v1.ScimConfigurationsIDDeleteNoContent]("Scim.Delete", func() (any, error) {
		return s.client.ScimConfigurationsIDDelete(ctx, v1.ScimConfigurationsIDDeleteParams{
			ID: uuid,
		})
	})
	return err
}

// RegenerateToken ユーザープロビジョニングのシークレットトークンを再発行する
func (s *scimOp) RegenerateToken(ctx context.Context, id string) (*v1.ScimConfigurationsIDRegenerateTokenPostOK, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return common.ErrorFromDecodedResponse[v1.ScimConfigurationsIDRegenerateTokenPostOK]("Scim.RegenerateToken", func() (any, error) {
		return s.client.ScimConfigurationsIDRegenerateTokenPost(ctx, v1.ScimConfigurationsIDRegenerateTokenPostParams{
			ID: uuid,
		})
	})
}
