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

package auth

import (
	"context"

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type AuthAPI interface {
	ReadPasswordPolicy(ctx context.Context) (*v1.PasswordPolicy, error)
	UpdatePasswordPolicy(ctx context.Context, req v1.PasswordPolicy) (*v1.PasswordPolicy, error)

	ReadAuthConditions(ctx context.Context) (*v1.AuthConditions, error)
	UpdateAuthConditions(ctx context.Context, req *v1.AuthConditions) (*v1.AuthConditions, error)

	ReadAuthContext(ctx context.Context) (*v1.GetAuthContextOK, error)
}

type authOp struct {
	client *v1.Client
}

var _ AuthAPI = (*authOp)(nil)

func NewAuthOp(client *v1.Client) AuthAPI { return &authOp{client} }

func (a *authOp) ReadPasswordPolicy(ctx context.Context) (*v1.PasswordPolicy, error) {
	return common.ErrorFromDecodedResponse[v1.PasswordPolicy]("Auth.ReadPasswordPolicy", func() (any, error) {
		return a.client.OrganizationPasswordPolicyGet(ctx)
	})
}

func (a *authOp) UpdatePasswordPolicy(ctx context.Context, req v1.PasswordPolicy) (*v1.PasswordPolicy, error) {
	return common.ErrorFromDecodedResponse[v1.PasswordPolicy]("Auth.UpdatePasswordPolicy", func() (any, error) {
		return a.client.OrganizationPasswordPolicyPut(ctx, &req)
	})
}

func (a *authOp) ReadAuthConditions(ctx context.Context) (*v1.AuthConditions, error) {
	return common.ErrorFromDecodedResponse[v1.AuthConditions]("Auth.ReadAuthConditions", func() (any, error) {
		return a.client.OrganizationAuthConditionsGet(ctx)
	})
}

func (a *authOp) UpdateAuthConditions(ctx context.Context, req *v1.AuthConditions) (*v1.AuthConditions, error) {
	return common.ErrorFromDecodedResponse[v1.AuthConditions]("Auth.UpdateAuthConditions", func() (any, error) {
		return a.client.OrganizationAuthConditionsPut(ctx, req)
	})
}

func (a *authOp) ReadAuthContext(ctx context.Context) (*v1.GetAuthContextOK, error) {
	return common.ErrorFromDecodedResponse[v1.GetAuthContextOK]("Auth.ReadAuthContext", func() (any, error) {
		return a.client.GetAuthContext(ctx)
	})
}
