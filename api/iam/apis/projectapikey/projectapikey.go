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

	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type ProjectAPIKeyAPI interface {
	List(ctx context.Context, params ListParams) (*v1.ListApiKeysOK, error)
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
	Ordering *v1.ListApiKeysOrdering
}

func (p *projectApiKeyOp) List(ctx context.Context, params ListParams) (*v1.ListApiKeysOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListApiKeysOK]("ProjectAPIKey.List", func() (any, error) {
		return p.client.ListApiKeys(ctx, v1.ListApiKeysParams{
			Page:     into.Opt[v1.OptInt](params.Page),
			PerPage:  into.Opt[v1.OptInt](params.PerPage),
			Ordering: into.Opt[v1.OptListApiKeysOrdering](params.Ordering),
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
		return p.client.CreateApiKey(ctx, &v1.CreateApiKeyReq{
			ProjectID:        params.ProjectID,
			Name:             params.Name,
			Description:      params.Description,
			ServerResourceID: into.Opt[v1.OptString](params.ServerResourceID),
			IamRoles:         params.IamRoles,
			ZoneID:           into.Opt[v1.OptString](params.Zone),
		})
	})
}

func (p *projectApiKeyOp) Read(ctx context.Context, id int) (*v1.ProjectApiKey, error) {
	return common.ErrorFromDecodedResponse[v1.ProjectApiKey]("ProjectAPIKey.Read", func() (any, error) {
		return p.client.ReadApiKey(ctx, v1.ReadApiKeyParams{ApikeyID: id})
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
		req := v1.UpdateApiKeyReq{
			Name:             params.Name,
			Description:      params.Description,
			ServerResourceID: into.Opt[v1.OptString](params.ServerResourceID),
			IamRoles:         params.IamRoles,
			ZoneID:           into.Opt[v1.OptString](params.Zone),
		}
		param := v1.UpdateApiKeyParams{
			ApikeyID: id,
		}
		return p.client.UpdateApiKey(ctx, &req, param)
	})
}

func (p *projectApiKeyOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteApiKeyNoContent]("ProjectAPIKey.Delete", func() (any, error) {
		return p.client.DeleteApiKey(ctx, v1.DeleteApiKeyParams{ApikeyID: id})
	})

	return err
}
