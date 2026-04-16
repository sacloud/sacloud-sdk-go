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

package servicepolicy

import (
	"context"

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type ServicePolicyAPI interface {
	Enable(ctx context.Context) error
	Disable(ctx context.Context) error
	IsEnabled(ctx context.Context) (bool, error)
	ListRuleTemplates(ctx context.Context, params ListRuleTemplatesParams) (*v1.ServicePolicyRuleTemplatesGetOK, error)
}

type servicePolicyOp struct {
	client *v1.Client
}

var _ ServicePolicyAPI = (*servicePolicyOp)(nil)

func NewServicePolicyOp(client *v1.Client) ServicePolicyAPI { return &servicePolicyOp{client} }

func (s *servicePolicyOp) Enable(ctx context.Context) error {
	_, err := common.ErrorFromDecodedResponse[v1.EnableServicePolicyPostNoContent]("ServicePolicy.Enable", func() (any, error) {
		return s.client.EnableServicePolicyPost(ctx)
	})
	return err
}

func (s *servicePolicyOp) Disable(ctx context.Context) error {
	_, err := common.ErrorFromDecodedResponse[v1.DisableServicePolicyPostNoContent]("ServicePolicy.Disable", func() (any, error) {
		return s.client.DisableServicePolicyPost(ctx)
	})
	return err
}

func (s *servicePolicyOp) IsEnabled(ctx context.Context) (bool, error) {
	if ret, err := common.ErrorFromDecodedResponse[v1.ServicePolicyStatusGetOK]("ServicePolicy.IsEnabled", func() (any, error) {
		return s.client.ServicePolicyStatusGet(ctx)
	}); err != nil {
		return false, err
	} else {
		return ret.Enabled, nil
	}
}

type ListRuleTemplatesParams struct {
	Page    *int
	PerPage *int
	Name    *string
	Code    *string
	Type    *v1.ServicePolicyRuleTemplatesGetType
}

func (s *servicePolicyOp) ListRuleTemplates(ctx context.Context, params ListRuleTemplatesParams) (*v1.ServicePolicyRuleTemplatesGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.ServicePolicyRuleTemplatesGetOK]("ServicePolicy.ListRuleTemplates", func() (any, error) {
		return s.client.ServicePolicyRuleTemplatesGet(ctx, v1.ServicePolicyRuleTemplatesGetParams{
			Page:    common.IntoOpt[v1.OptInt](params.Page),
			PerPage: common.IntoOpt[v1.OptInt](params.PerPage),
			Name:    common.IntoOpt[v1.OptString](params.Name),
			Code:    common.IntoOpt[v1.OptString](params.Code),
			Type:    common.IntoOpt[v1.OptServicePolicyRuleTemplatesGetType](params.Type),
		})
	})
}
