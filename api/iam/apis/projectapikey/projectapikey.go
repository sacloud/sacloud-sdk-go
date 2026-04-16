// Copyright 2025- The sacloud/iam-api-go authors
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

package projectapikey

import (
	"context"

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type ProjectAPIKeyAPI interface {
	List(ctx context.Context, params ListParams) (*v1.CompatAPIKeysGetOK, error)
	Create(ctx context.Context, params CreateParams) (*v1.ProjectApiKeyWithSecret, error)
	Read(ctx context.Context, id int) (*v1.ProjectApiKey, error)
	Update(ctx context.Context, id int, params UpdateParams) (*v1.ProjectApiKey, error)
	Delete(ctx context.Context, id int) error
}

type projectApiKeyOp struct {
	client *v1.Client
}

func NewProjectAPIKeyOp(client *v1.Client) ProjectAPIKeyAPI {
	return &projectApiKeyOp{client: client}
}

type ListParams struct {
	Page     *int
	PerPage  *int
	Ordering *v1.CompatAPIKeysGetOrdering
}

func (p *projectApiKeyOp) List(ctx context.Context, params ListParams) (*v1.CompatAPIKeysGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.CompatAPIKeysGetOK]("ProjectAPIKey.List", func() (any, error) {
		return p.client.CompatAPIKeysGet(ctx, v1.CompatAPIKeysGetParams{
			Page:     common.IntoOpt[v1.OptInt](params.Page),
			PerPage:  common.IntoOpt[v1.OptInt](params.PerPage),
			Ordering: common.IntoOpt[v1.OptCompatAPIKeysGetOrdering](params.Ordering),
		})
	})
}

type CreateParams struct {
	ProjectID        int
	Name             string
	Description      string
	ServerResourceID *string
	IamRoles         []string
	Zone             *string
}

func (p *projectApiKeyOp) Create(ctx context.Context, params CreateParams) (*v1.ProjectApiKeyWithSecret, error) {
	return common.ErrorFromDecodedResponse[v1.ProjectApiKeyWithSecret]("ProjectAPIKey.Create", func() (any, error) {
		return p.client.CompatAPIKeysPost(ctx, &v1.CompatAPIKeysPostReq{
			ProjectID:        params.ProjectID,
			Name:             params.Name,
			Description:      params.Description,
			ServerResourceID: common.IntoOpt[v1.OptString](params.ServerResourceID),
			IamRoles:         params.IamRoles,
			ZoneID:           common.IntoOpt[v1.OptString](params.Zone),
		})
	})
}

func (p *projectApiKeyOp) Read(ctx context.Context, id int) (*v1.ProjectApiKey, error) {
	return common.ErrorFromDecodedResponse[v1.ProjectApiKey]("ProjectAPIKey.Read", func() (any, error) {
		return p.client.CompatAPIKeysApikeyIDGet(ctx, v1.CompatAPIKeysApikeyIDGetParams{ApikeyID: id})
	})
}

type UpdateParams struct {
	Name             string
	Description      string
	ServerResourceID *string
	IamRoles         []string
	Zone             *string
}

func (p *projectApiKeyOp) Update(ctx context.Context, id int, params UpdateParams) (*v1.ProjectApiKey, error) {
	return common.ErrorFromDecodedResponse[v1.ProjectApiKey]("ProjectAPIKey.Update", func() (any, error) {
		req := v1.CompatAPIKeysApikeyIDPutReq{
			Name:             params.Name,
			Description:      params.Description,
			ServerResourceID: common.IntoOpt[v1.OptString](params.ServerResourceID),
			IamRoles:         params.IamRoles,
			ZoneID:           common.IntoOpt[v1.OptString](params.Zone),
		}
		param := v1.CompatAPIKeysApikeyIDPutParams{
			ApikeyID: id,
		}
		return p.client.CompatAPIKeysApikeyIDPut(ctx, &req, param)
	})
}

func (p *projectApiKeyOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.CompatAPIKeysApikeyIDDeleteNoContent]("ProjectAPIKey.Delete", func() (any, error) {
		return p.client.CompatAPIKeysApikeyIDDelete(ctx, v1.CompatAPIKeysApikeyIDDeleteParams{ApikeyID: id})
	})

	return err
}
