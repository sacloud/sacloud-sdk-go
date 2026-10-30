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
	const methodName = "Revision.Create"

	res, err := op.client.CreateRevision(ctx, &req, v1.CreateRevisionParams{ID: workflowID})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreateRevisionCreated:
		return &r.Revision, nil
	case *v1.CreateRevisionBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.CreateRevisionUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.CreateRevisionPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.CreateRevisionForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.CreateRevisionNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.CreateRevisionConflict:
		return nil, NewAPIError(methodName, http.StatusConflict, errors.New(r.Message))
	case *v1.CreateRevisionInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *revisionOp) List(ctx context.Context, params v1.ListRevisionsParams) (*v1.ListRevisionsOK, error) {
	const methodName = "Revision.List"

	res, err := op.client.ListRevisions(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListRevisionsOK:
		return r, nil
	case *v1.ListRevisionsBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListRevisionsUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListRevisionsPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.ListRevisionsForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListRevisionsNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListRevisionsInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *revisionOp) Read(ctx context.Context, workflowID string, revisionNumber int) (*v1.ReadRevisionOKRevision, error) {
	const methodName = "Revision.Read"

	res, err := op.client.ReadRevision(ctx, v1.ReadRevisionParams{
		ID:         workflowID,
		RevisionId: revisionNumber,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ReadRevisionOK:
		return &r.Revision, nil
	case *v1.ReadRevisionBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ReadRevisionUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ReadRevisionPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.ReadRevisionForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ReadRevisionNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ReadRevisionInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *revisionOp) UpdateAlias(ctx context.Context, workflowID string, revisionNumber int, req v1.UpdateRevisionAliasReq) (*v1.UpdateRevisionAliasOKRevision, error) {
	const methodName = "Revision.UpdateAlias"

	res, err := op.client.UpdateRevisionAlias(ctx, &req, v1.UpdateRevisionAliasParams{
		ID:         workflowID,
		RevisionId: revisionNumber,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.UpdateRevisionAliasOK:
		return &r.Revision, nil
	case *v1.UpdateRevisionAliasBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.UpdateRevisionAliasUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.UpdateRevisionAliasPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.UpdateRevisionAliasForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.UpdateRevisionAliasNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.UpdateRevisionAliasConflict:
		return nil, NewAPIError(methodName, http.StatusConflict, errors.New(r.Message))
	case *v1.UpdateRevisionAliasInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *revisionOp) DeleteAlias(ctx context.Context, workflowID string, revisionNumber int) error {
	const methodName = "Revision.DeleteAlias"

	res, err := op.client.DeleteRevisionAlias(ctx, v1.DeleteRevisionAliasParams{
		ID:         workflowID,
		RevisionId: revisionNumber,
	})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteRevisionAliasOK:
		return nil
	case *v1.DeleteRevisionAliasBadRequest:
		return NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.DeleteRevisionAliasUnauthorized:
		return NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.DeleteRevisionAliasPaymentRequired:
		return NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.DeleteRevisionAliasForbidden:
		return NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.DeleteRevisionAliasNotFound:
		return NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.DeleteRevisionAliasInternalServerError:
		return NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return NewAPIError(methodName, 0, err)
	}
}
