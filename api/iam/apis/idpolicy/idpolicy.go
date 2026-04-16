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

// Package idpolicy provides the IdPolicyAPI that wraps the generated v1 client.
package idpolicy

import (
	"context"

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type IDPolicyAPI interface {
	ReadOrganizationIdPolicy(ctx context.Context) ([]v1.IdPolicy, error)
	UpdateOrganizationIdPolicy(ctx context.Context, bindings []v1.IdPolicy) ([]v1.IdPolicy, error)
}

type idPolicyOp struct {
	client *v1.Client
}

func NewIDPolicyOp(client *v1.Client) IDPolicyAPI { return &idPolicyOp{client: client} }

func (o *idPolicyOp) ReadOrganizationIdPolicy(ctx context.Context) ([]v1.IdPolicy, error) {
	if ret, err := common.ErrorFromDecodedResponse[v1.OrganizationIDPolicyGetOK]("IdPolicy.ReadOrganizationIdPolicy", func() (any, error) {
		return o.client.OrganizationIDPolicyGet(ctx)
	}); err != nil {
		return nil, err
	} else {
		return ret.GetBindings(), nil
	}
}

func (o *idPolicyOp) UpdateOrganizationIdPolicy(ctx context.Context, bindings []v1.IdPolicy) ([]v1.IdPolicy, error) {
	if ret, err := common.ErrorFromDecodedResponse[v1.OrganizationIDPolicyPutOK]("IdPolicy.UpdateOrganizationIdPolicy", func() (any, error) {
		return o.client.OrganizationIDPolicyPut(ctx, &v1.OrganizationIDPolicyPutReq{Bindings: bindings})
	}); err != nil {
		return nil, err
	} else {
		return ret.GetBindings(), nil
	}
}
