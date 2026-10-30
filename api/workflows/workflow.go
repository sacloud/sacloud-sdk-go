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

type WorkflowAPI interface {
	Create(ctx context.Context, request v1.CreateWorkflowReq) (*v1.CreateWorkflowCreatedWorkflow, error)
	List(ctx context.Context, parameter v1.ListWorkflowsParams) (*v1.ListWorkflowsOK, error)
	ListSuggest(ctx context.Context, parameter v1.ListSuggestionsParams) (*v1.ListSuggestionsOK, error)
	Read(ctx context.Context, id string) (*v1.ReadWorkflowOKWorkflow, error)
	Update(ctx context.Context, id string, request v1.PatchWorkflowReq) (*v1.PatchWorkflowOKWorkflow, error)
	Delete(ctx context.Context, id string) error
}

var _ WorkflowAPI = (*workflowOp)(nil)

type workflowOp struct {
	client *v1.Client
}

func NewWorkflowOp(client *v1.Client) WorkflowAPI {
	return &workflowOp{client: client}
}

func (op *workflowOp) Create(ctx context.Context, req v1.CreateWorkflowReq) (*v1.CreateWorkflowCreatedWorkflow, error) {
	const methodName = "Workflow.Create"

	res, err := op.client.CreateWorkflow(ctx, &req)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreateWorkflowCreated:
		return &r.Workflow, nil
	case *v1.CreateWorkflowBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.CreateWorkflowUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.CreateWorkflowPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.CreateWorkflowForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.CreateWorkflowNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.CreateWorkflowInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *workflowOp) List(ctx context.Context, params v1.ListWorkflowsParams) (*v1.ListWorkflowsOK, error) {
	const methodName = "Workflow.List"

	res, err := op.client.ListWorkflows(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListWorkflowsOK:
		return r, nil
	case *v1.ListWorkflowsBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListWorkflowsUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListWorkflowsPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.ListWorkflowsForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListWorkflowsNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListWorkflowsInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *workflowOp) ListSuggest(ctx context.Context, params v1.ListSuggestionsParams) (*v1.ListSuggestionsOK, error) {
	const methodName = "Workflow.ListSuggest"

	res, err := op.client.ListSuggestions(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListSuggestionsOK:
		return r, nil
	case *v1.ListSuggestionsBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListSuggestionsUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListSuggestionsPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.ListSuggestionsForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListSuggestionsNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListSuggestionsInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *workflowOp) Read(ctx context.Context, id string) (*v1.ReadWorkflowOKWorkflow, error) {
	const methodName = "Workflow.Read"

	res, err := op.client.ReadWorkflow(ctx, v1.ReadWorkflowParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ReadWorkflowOK:
		return &r.Workflow, nil
	case *v1.ReadWorkflowBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ReadWorkflowUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ReadWorkflowPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.ReadWorkflowForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ReadWorkflowNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ReadWorkflowInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *workflowOp) Update(ctx context.Context, id string, req v1.PatchWorkflowReq) (*v1.PatchWorkflowOKWorkflow, error) {
	const methodName = "Workflow.Update"

	res, err := op.client.PatchWorkflow(ctx, &req, v1.PatchWorkflowParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.PatchWorkflowOK:
		return &r.Workflow, nil
	case *v1.PatchWorkflowBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.PatchWorkflowUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.PatchWorkflowPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.PatchWorkflowForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.PatchWorkflowNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.PatchWorkflowInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *workflowOp) Delete(ctx context.Context, id string) error {
	const methodName = "Workflow.Delete"

	res, err := op.client.DeleteWorkflow(ctx, v1.DeleteWorkflowParams{ID: id})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteWorkflowOK:
		return nil
	case *v1.DeleteWorkflowBadRequest:
		return NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.DeleteWorkflowUnauthorized:
		return NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.DeleteWorkflowPaymentRequired:
		return NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.DeleteWorkflowForbidden:
		return NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.DeleteWorkflowNotFound:
		return NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.DeleteWorkflowConflict:
		return NewAPIError(methodName, http.StatusConflict, errors.New(r.Message))
	case *v1.DeleteWorkflowInternalServerError:
		return NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return NewAPIError(methodName, 0, err)
	}
}
