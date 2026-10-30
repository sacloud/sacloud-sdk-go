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

type ExecutionAPI interface {
	Create(ctx context.Context, workflowID string, req v1.OptCreateExecutionReq) (*v1.CreateExecutionCreatedExecution, error)
	List(ctx context.Context, params v1.ListExecutionsParams) (*v1.ListExecutionsOK, error)
	Read(ctx context.Context, workflowID, executionID string) (*v1.ReadExecutionOKExecution, error)
	Cancel(ctx context.Context, workflowID, executionID string) (*v1.CancelExecutionOKExecution, error)
	Delete(ctx context.Context, workflowID, executionID string) error
	ListHistory(ctx context.Context, params v1.ListExecutionHistoriesParams) (*v1.ListExecutionHistoriesOK, error)
}

var _ ExecutionAPI = (*executionOp)(nil)

type executionOp struct {
	client *v1.Client
}

func NewExecutionOp(client *v1.Client) ExecutionAPI {
	return &executionOp{client: client}
}

func (op *executionOp) Create(ctx context.Context, workflowID string, req v1.OptCreateExecutionReq) (*v1.CreateExecutionCreatedExecution, error) {
	const methodName = "Execution.Create"

	res, err := op.client.CreateExecution(ctx, req, v1.CreateExecutionParams{ID: workflowID})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreateExecutionCreated:
		return &r.Execution, nil
	case *v1.CreateExecutionBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.CreateExecutionUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.CreateExecutionForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.CreateExecutionNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.CreateExecutionConflict:
		return nil, NewAPIError(methodName, http.StatusConflict, errors.New(r.Message))
	case *v1.CreateExecutionInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *executionOp) List(ctx context.Context, params v1.ListExecutionsParams) (*v1.ListExecutionsOK, error) {
	const methodName = "Execution.List"

	res, err := op.client.ListExecutions(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListExecutionsOK:
		return r, nil
	case *v1.ListExecutionsBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListExecutionsUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListExecutionsForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListExecutionsNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListExecutionsInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *executionOp) Read(ctx context.Context, workflowID, executionID string) (*v1.ReadExecutionOKExecution, error) {
	const methodName = "Execution.Read"

	res, err := op.client.ReadExecution(ctx, v1.ReadExecutionParams{
		ID:          workflowID,
		ExecutionId: executionID,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ReadExecutionOK:
		return &r.Execution, nil
	case *v1.ReadExecutionBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ReadExecutionUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ReadExecutionForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ReadExecutionNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ReadExecutionInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *executionOp) Cancel(ctx context.Context, workflowID, executionID string) (*v1.CancelExecutionOKExecution, error) {
	const methodName = "Execution.Cancel"

	res, err := op.client.CancelExecution(ctx, v1.CancelExecutionParams{
		ID:          workflowID,
		ExecutionId: executionID,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CancelExecutionOK:
		return &r.Execution, nil
	case *v1.CancelExecutionBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.CancelExecutionUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.CancelExecutionForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.CancelExecutionNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.CancelExecutionInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}

func (op *executionOp) Delete(ctx context.Context, workflowID, executionID string) error {
	const methodName = "Execution.Delete"

	res, err := op.client.DeleteExecution(ctx, v1.DeleteExecutionParams{
		ID:          workflowID,
		ExecutionId: executionID,
	})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteExecutionOK:
		return nil
	case *v1.DeleteExecutionBadRequest:
		return NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.DeleteExecutionUnauthorized:
		return NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.DeleteExecutionForbidden:
		return NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.DeleteExecutionNotFound:
		return NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.DeleteExecutionInternalServerError:
		return NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return NewAPIError(methodName, 0, err)
	}
}

func (op *executionOp) ListHistory(ctx context.Context, params v1.ListExecutionHistoriesParams) (*v1.ListExecutionHistoriesOK, error) {
	const methodName = "Execution.ListHistory"

	res, err := op.client.ListExecutionHistories(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListExecutionHistoriesOK:
		return r, nil
	case *v1.ListExecutionHistoriesBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListExecutionHistoriesUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListExecutionHistoriesForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListExecutionHistoriesNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListExecutionHistoriesInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, err)
	}
}
