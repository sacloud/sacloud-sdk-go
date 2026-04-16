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

package idrole

import (
	"context"

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type IDRoleAPI interface {
	List(ctx context.Context, page, perPage *int) (*v1.IDRolesGetOK, error)
	Read(ctx context.Context, id string) (*v1.IdRole, error)
}

type idRoleOp struct {
	client *v1.Client
}

func NewIdRoleOp(client *v1.Client) IDRoleAPI { return &idRoleOp{client: client} }

func (i *idRoleOp) List(ctx context.Context, page, perPage *int) (*v1.IDRolesGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.IDRolesGetOK]("IdRole.List", func() (any, error) {
		return i.client.IDRolesGet(ctx, v1.IDRolesGetParams{
			Page:    common.IntoOpt[v1.OptInt](page),
			PerPage: common.IntoOpt[v1.OptInt](perPage),
		})
	})
}

func (i *idRoleOp) Read(ctx context.Context, id string) (*v1.IdRole, error) {
	return common.ErrorFromDecodedResponse[v1.IdRole]("IdRole.Read", func() (any, error) {
		return i.client.IDRolesIDRoleIDGet(ctx, v1.IDRolesIDRoleIDGetParams{IDRoleID: id})
	})
}
