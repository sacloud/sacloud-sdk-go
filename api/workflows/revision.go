// Copyright 2025- The sacloud/workflows-api-go Authors
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

package workflows

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/workflows/apis/v1"
)

type RevisionAPI interface {
	Create(ctx context.Context, workflowID string, req v1.CreateRevisionReq) (*v1.CreateRevisionCreatedRevision, error)
	List(ctx context.Context, params v1.ListRevisionsParams) (*v1.ListRevisionsOK, error)
	Read(ctx context.Context, workflowID string, revisionNumber int) (*v1.ReadRevisionOKRevision, error)
	UpdateAlias(ctx context.Context, workflowID string, revisionNumber int, req v1.UpdateRevisionAliasReq) (*v1.UpdateRevisionAliasOKRevision, error)
	DeleteAlias(ctx context.Context, workflowID string, revisionNumber int) error
}

var _ RevisionAPI = (*revisionOp)(nil)

type revisionOp struct {
	client *v1.Client
}

func NewRevisionOp(client *v1.Client) RevisionAPI {
	return &revisionOp{client: client}
}

func (op *revisionOp) Create(ctx context.Context, workflowID string, req v1.CreateRevisionReq) (*v1.CreateRevisionCreatedRevision, error) {
	res, err := errorFromDecodedResponse[v1.CreateRevisionCreated]("Revision.Create", func() (any, error) {
		return op.client.CreateRevision(ctx, &req, v1.CreateRevisionParams{ID: workflowID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetRevision()), nil
}

func (op *revisionOp) List(ctx context.Context, params v1.ListRevisionsParams) (*v1.ListRevisionsOK, error) {
	return errorFromDecodedResponse[v1.ListRevisionsOK]("Revision.List", func() (any, error) { return op.client.ListRevisions(ctx, params) })
}

func (op *revisionOp) Read(ctx context.Context, workflowID string, revisionNumber int) (*v1.ReadRevisionOKRevision, error) {
	res, err := errorFromDecodedResponse[v1.ReadRevisionOK]("Revision.Read", func() (any, error) {
		return op.client.ReadRevision(ctx, v1.ReadRevisionParams{
			ID:         workflowID,
			RevisionId: revisionNumber,
		})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetRevision()), nil
}

func (op *revisionOp) UpdateAlias(ctx context.Context, workflowID string, revisionNumber int, req v1.UpdateRevisionAliasReq) (*v1.UpdateRevisionAliasOKRevision, error) {
	res, err := errorFromDecodedResponse[v1.UpdateRevisionAliasOK]("Revision.UpdateAlias", func() (any, error) {
		return op.client.UpdateRevisionAlias(ctx, &req, v1.UpdateRevisionAliasParams{
			ID:         workflowID,
			RevisionId: revisionNumber,
		})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetRevision()), nil
}

func (op *revisionOp) DeleteAlias(ctx context.Context, workflowID string, revisionNumber int) error {
	return errorFromDecodedResponseE[v1.DeleteRevisionAliasOK]("Revision.DeleteAlias", func() (any, error) {
		return op.client.DeleteRevisionAlias(ctx, v1.DeleteRevisionAliasParams{ID: workflowID, RevisionId: revisionNumber})
	})
}
