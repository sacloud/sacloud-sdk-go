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

package iamrole

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type IAMRoleAPI interface {
	List(ctx context.Context, page, perPage *int) (*v1.ListIamRolesOK, error)
	Read(ctx context.Context, id string) (*v1.IamRole, error)
}

type iamRoleOp struct {
	client *v1.Client
}

func NewIAMRoleOp(client *v1.Client) IAMRoleAPI { return &iamRoleOp{client: client} }

func (i *iamRoleOp) List(ctx context.Context, page, perPage *int) (*v1.ListIamRolesOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListIamRolesOK]("IAMRole.List", func() (any, error) {
		return i.client.ListIamRoles(ctx, v1.ListIamRolesParams{
			Page:    into.Opt[v1.OptInt](page),
			PerPage: into.Opt[v1.OptInt](perPage),
		})
	})
}

func (i *iamRoleOp) Read(ctx context.Context, id string) (*v1.IamRole, error) {
	return common.ErrorFromDecodedResponse[v1.IamRole]("IAMRole.Read", func() (any, error) {
		return i.client.ReadIamRole(ctx, v1.ReadIamRoleParams{IamRoleID: id})
	})
}
