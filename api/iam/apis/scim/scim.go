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
	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

// ScimAPI SCIM API
type ScimAPI interface {
	// List ユーザープロビジョニング一覧を取得する
	List(ctx context.Context, params ListParams) (*v1.ListScimConfigurationsOK, error)
	// Create ユーザープロビジョニングを作成する
	Create(ctx context.Context, params CreateParams) (*v1.ScimConfiguration, error)
	// Read ユーザープロビジョニングを取得する
	Read(ctx context.Context, id string) (*v1.ScimConfigurationBase, error)
	// Update ユーザープロビジョニングを更新する
	Update(ctx context.Context, id string, params UpdateParams) (*v1.ScimConfigurationBase, error)
	// Delete ユーザープロビジョニングを削除する
	Delete(ctx context.Context, id string) error
	// RegenerateToken ユーザープロビジョニングのシークレットトークンを再発行する
	RegenerateToken(ctx context.Context, id string) (*v1.RegenerateScimConfigurationTokenOK, error)
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
func (s *scimOp) List(ctx context.Context, params ListParams) (*v1.ListScimConfigurationsOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListScimConfigurationsOK]("Scim.List", func() (any, error) {
		return s.client.ListScimConfigurations(ctx, v1.ListScimConfigurationsParams{
			Page:    into.Opt[v1.OptInt](params.Page),
			PerPage: into.Opt[v1.OptInt](params.PerPage),
		})
	})
}

// Create ユーザープロビジョニングを作成する
func (s *scimOp) Create(ctx context.Context, params CreateParams) (*v1.ScimConfiguration, error) {
	return common.ErrorFromDecodedResponse[v1.ScimConfiguration]("Scim.Create", func() (any, error) {
		return s.client.CreateScimConfiguration(ctx, &v1.CreateScimConfigurationReq{
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
		return s.client.ReadScimConfiguration(ctx, v1.ReadScimConfigurationParams{
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
		return s.client.UpdateScimConfiguration(ctx, &v1.UpdateScimConfigurationReq{
			Name: params.Name,
		}, v1.UpdateScimConfigurationParams{
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
	_, err = common.ErrorFromDecodedResponse[v1.DeleteScimConfigurationNoContent]("Scim.Delete", func() (any, error) {
		return s.client.DeleteScimConfiguration(ctx, v1.DeleteScimConfigurationParams{
			ID: uuid,
		})
	})
	return err
}

// RegenerateToken ユーザープロビジョニングのシークレットトークンを再発行する
func (s *scimOp) RegenerateToken(ctx context.Context, id string) (*v1.RegenerateScimConfigurationTokenOK, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return common.ErrorFromDecodedResponse[v1.RegenerateScimConfigurationTokenOK]("Scim.RegenerateToken", func() (any, error) {
		return s.client.RegenerateScimConfigurationToken(ctx, v1.RegenerateScimConfigurationTokenParams{
			ID: uuid,
		})
	})
}
