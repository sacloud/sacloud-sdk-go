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

type WorkflowAPI interface {
	Create(ctx context.Context, request v1.CreateWorkflowReq) (*v1.CreateWorkflowCreatedWorkflow, error)
	List(ctx context.Context, parameter v1.ListWorkflowParams) (*v1.ListWorkflowOK, error)
	ListSuggest(ctx context.Context, parameter v1.ListWorkflowSuggestParams) (*v1.ListWorkflowSuggestOK, error)
	Read(ctx context.Context, id string) (*v1.GetWorkflowOKWorkflow, error)
	Update(ctx context.Context, id string, request v1.UpdateWorkflowReq) (*v1.UpdateWorkflowOKWorkflow, error)
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

func (op *workflowOp) List(ctx context.Context, params v1.ListWorkflowParams) (*v1.ListWorkflowOK, error) {
	const methodName = "Workflow.List"

	res, err := op.client.ListWorkflow(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListWorkflowOK:
		return r, nil
	case *v1.ListWorkflowBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListWorkflowUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListWorkflowPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.ListWorkflowForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListWorkflowNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListWorkflowInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *workflowOp) ListSuggest(ctx context.Context, params v1.ListWorkflowSuggestParams) (*v1.ListWorkflowSuggestOK, error) {
	const methodName = "Workflow.ListSuggest"

	res, err := op.client.ListWorkflowSuggest(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListWorkflowSuggestOK:
		return r, nil
	case *v1.ListWorkflowSuggestBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListWorkflowSuggestUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListWorkflowSuggestPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.ListWorkflowSuggestForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListWorkflowSuggestNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListWorkflowSuggestInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *workflowOp) Read(ctx context.Context, id string) (*v1.GetWorkflowOKWorkflow, error) {
	const methodName = "Workflow.Read"

	res, err := op.client.GetWorkflow(ctx, v1.GetWorkflowParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.GetWorkflowOK:
		return &r.Workflow, nil
	case *v1.GetWorkflowBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.GetWorkflowUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.GetWorkflowPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.GetWorkflowForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.GetWorkflowNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.GetWorkflowInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *workflowOp) Update(ctx context.Context, id string, req v1.UpdateWorkflowReq) (*v1.UpdateWorkflowOKWorkflow, error) {
	const methodName = "Workflow.Update"

	res, err := op.client.UpdateWorkflow(ctx, &req, v1.UpdateWorkflowParams{ID: id})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.UpdateWorkflowOK:
		return &r.Workflow, nil
	case *v1.UpdateWorkflowBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.UpdateWorkflowUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.UpdateWorkflowPaymentRequired:
		return nil, NewAPIError(methodName, http.StatusPaymentRequired, errors.New(r.Message))
	case *v1.UpdateWorkflowForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.UpdateWorkflowNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.UpdateWorkflowInternalServerError:
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
