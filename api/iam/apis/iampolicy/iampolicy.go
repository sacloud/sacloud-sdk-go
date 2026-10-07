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

package iampolicy

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
)

type IAMPolicyAPI interface {
	ReadOrganizationPolicy(ctx context.Context) ([]v1.IamPolicy, error)
	UpdateOrganizationPolicy(ctx context.Context, bindings []v1.IamPolicy) ([]v1.IamPolicy, error)

	ReadProjectPolicy(ctx context.Context, projectID int) ([]v1.IamPolicy, error)
	UpdateProjectPolicy(ctx context.Context, projectID int, bindings []v1.IamPolicy) ([]v1.IamPolicy, error)

	ReadFolderPolicy(ctx context.Context, folderID int) ([]v1.IamPolicy, error)
	UpdateFolderPolicy(ctx context.Context, folderID int, bindings []v1.IamPolicy) ([]v1.IamPolicy, error)
}

type iamPolicyOp struct {
	client *v1.Client
}

func NewIAMPolicyOp(client *v1.Client) IAMPolicyAPI { return &iamPolicyOp{client: client} }

func (op *iamPolicyOp) ReadOrganizationPolicy(ctx context.Context) ([]v1.IamPolicy, error) {
	if res, err := common.ErrorFromDecodedResponse[v1.ReadOrganizationIamPolicyOK]("IamPolicy.ReadOrganizationPolicy", func() (any, error) {
		return op.client.ReadOrganizationIamPolicy(ctx)
	}); err != nil {
		return nil, err
	} else {
		return res.Bindings, nil
	}
}

func (op *iamPolicyOp) UpdateOrganizationPolicy(ctx context.Context, bindings []v1.IamPolicy) ([]v1.IamPolicy, error) {
	if res, err := common.ErrorFromDecodedResponse[v1.UpdateOrganizationIamPolicyOK]("IamPolicy.UpdateOrganizationPolicy", func() (any, error) {
		return op.client.UpdateOrganizationIamPolicy(ctx, &v1.UpdateOrganizationIamPolicyReq{Bindings: bindings})
	}); err != nil {
		return nil, err
	} else {
		return res.Bindings, nil
	}
}

func (op *iamPolicyOp) ReadProjectPolicy(ctx context.Context, id int) ([]v1.IamPolicy, error) {
	if res, err := common.ErrorFromDecodedResponse[v1.ReadProjectIamPolicyOK]("IamPolicy.ReadProjectPolicy", func() (any, error) {
		return op.client.ReadProjectIamPolicy(ctx, v1.ReadProjectIamPolicyParams{ProjectID: id})
	}); err != nil {
		return nil, err
	} else {
		return res.Bindings, nil
	}
}

func (op *iamPolicyOp) UpdateProjectPolicy(ctx context.Context, id int, bindings []v1.IamPolicy) ([]v1.IamPolicy, error) {
	if res, err := common.ErrorFromDecodedResponse[v1.UpdateProjectIamPolicyOK]("IamPolicy.UpdateProjectPolicy", func() (any, error) {
		request := v1.UpdateProjectIamPolicyReq{Bindings: bindings}
		params := v1.UpdateProjectIamPolicyParams{ProjectID: id}
		return op.client.UpdateProjectIamPolicy(ctx, &request, params)
	}); err != nil {
		return nil, err
	} else {
		return res.Bindings, nil
	}
}

func (op *iamPolicyOp) ReadFolderPolicy(ctx context.Context, id int) ([]v1.IamPolicy, error) {
	if res, err := common.ErrorFromDecodedResponse[v1.ReadFolderIamPolicyOK]("IamPolicy.ReadFolderPolicy", func() (any, error) {
		return op.client.ReadFolderIamPolicy(ctx, v1.ReadFolderIamPolicyParams{FolderID: id})
	}); err != nil {
		return nil, err
	} else {
		return res.Bindings, nil
	}
}

func (op *iamPolicyOp) UpdateFolderPolicy(ctx context.Context, id int, bindings []v1.IamPolicy) ([]v1.IamPolicy, error) {
	if res, err := common.ErrorFromDecodedResponse[v1.UpdateFolderIamPolicyOK]("IamPolicy.UpdateFolderPolicy", func() (any, error) {
		request := v1.UpdateFolderIamPolicyReq{Bindings: bindings}
		params := v1.UpdateFolderIamPolicyParams{FolderID: id}
		return op.client.UpdateFolderIamPolicy(ctx, &request, params)
	}); err != nil {
		return nil, err
	} else {
		return res.Bindings, nil
	}
}
