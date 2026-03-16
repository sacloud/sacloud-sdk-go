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
	"errors"
	"net/http"

	v1 "github.com/sacloud/workflows-api-go/apis/v1"
)

type RevisionAPI interface {
	Create(ctx context.Context, workflowID string, req v1.CreateWorkflowRevisionReq) (*v1.CreateWorkflowRevisionCreatedRevision, error)
	List(ctx context.Context, params v1.ListWorkflowRevisionsParams) (*v1.ListWorkflowRevisionsOK, error)
	Read(ctx context.Context, workflowID string, revisionNumber int) (*v1.GetWorkflowRevisionsOKRevision, error)
	UpdateAlias(ctx context.Context, workflowID string, revisionNumber int, req v1.UpdateWorkflowRevisionAliasReq) (*v1.UpdateWorkflowRevisionAliasOKRevision, error)
	DeleteAlias(ctx context.Context, workflowID string, revisionNumber int) error
}

var _ RevisionAPI = (*revisionOp)(nil)

type revisionOp struct {
	client *v1.Client
}

func NewRevisionOp(client *v1.Client) RevisionAPI {
	return &revisionOp{client: client}
}

func (op *revisionOp) Create(ctx context.Context, workflowID string, req v1.CreateWorkflowRevisionReq) (*v1.CreateWorkflowRevisionCreatedRevision, error) {
	const methodName = "Revision.Create"

	res, err := op.client.CreateWorkflowRevision(ctx, &req, v1.CreateWorkflowRevisionParams{ID: workflowID})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreateWorkflowRevisionCreated:
		return &r.Revision, nil
	case *v1.CreateWorkflowRevisionBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.CreateWorkflowRevisionUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.CreateWorkflowRevisionPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.CreateWorkflowRevisionForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.CreateWorkflowRevisionNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.CreateWorkflowRevisionConflict:
		return nil, NewAPIError(methodName, http.StatusConflict, errors.New(r.Message))
	case *v1.CreateWorkflowRevisionInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *revisionOp) List(ctx context.Context, params v1.ListWorkflowRevisionsParams) (*v1.ListWorkflowRevisionsOK, error) {
	const methodName = "Revision.List"

	res, err := op.client.ListWorkflowRevisions(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListWorkflowRevisionsOK:
		return r, nil
	case *v1.ListWorkflowRevisionsBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListWorkflowRevisionsUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListWorkflowRevisionsPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.ListWorkflowRevisionsForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListWorkflowRevisionsNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListWorkflowRevisionsInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *revisionOp) Read(ctx context.Context, workflowID string, revisionNumber int) (*v1.GetWorkflowRevisionsOKRevision, error) {
	const methodName = "Revision.Read"

	res, err := op.client.GetWorkflowRevisions(ctx, v1.GetWorkflowRevisionsParams{
		ID:         workflowID,
		RevisionId: revisionNumber,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.GetWorkflowRevisionsOK:
		return &r.Revision, nil
	case *v1.GetWorkflowRevisionsBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.GetWorkflowRevisionsUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.GetWorkflowRevisionsPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.GetWorkflowRevisionsForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.GetWorkflowRevisionsNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.GetWorkflowRevisionsInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *revisionOp) UpdateAlias(ctx context.Context, workflowID string, revisionNumber int, req v1.UpdateWorkflowRevisionAliasReq) (*v1.UpdateWorkflowRevisionAliasOKRevision, error) {
	const methodName = "Revision.UpdateAlias"

	res, err := op.client.UpdateWorkflowRevisionAlias(ctx, &req, v1.UpdateWorkflowRevisionAliasParams{
		ID:         workflowID,
		RevisionId: revisionNumber,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.UpdateWorkflowRevisionAliasOK:
		return &r.Revision, nil
	case *v1.UpdateWorkflowRevisionAliasBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.UpdateWorkflowRevisionAliasUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.UpdateWorkflowRevisionAliasPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.UpdateWorkflowRevisionAliasForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.UpdateWorkflowRevisionAliasNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.UpdateWorkflowRevisionAliasConflict:
		return nil, NewAPIError(methodName, http.StatusConflict, errors.New(r.Message))
	case *v1.UpdateWorkflowRevisionAliasInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *revisionOp) DeleteAlias(ctx context.Context, workflowID string, revisionNumber int) error {
	const methodName = "Revision.DeleteAlias"

	res, err := op.client.DeleteWorkflowRevisionAlias(ctx, v1.DeleteWorkflowRevisionAliasParams{
		ID:         workflowID,
		RevisionId: revisionNumber,
	})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteWorkflowRevisionAliasOK:
		return nil
	case *v1.DeleteWorkflowRevisionAliasBadRequest:
		return NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.DeleteWorkflowRevisionAliasUnauthorized:
		return NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.DeleteWorkflowRevisionAliasPaymentRequired:
		return NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.DeleteWorkflowRevisionAliasForbidden:
		return NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.DeleteWorkflowRevisionAliasNotFound:
		return NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.DeleteWorkflowRevisionAliasInternalServerError:
		return NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return NewAPIError(methodName, 0, err)
	}
}
