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
	res, err := errorFromDecodedResponse[v1.CreateExecutionCreated]("Execution.Create", func() (any, error) {
		return op.client.CreateExecution(ctx, req, v1.CreateExecutionParams{ID: workflowID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetExecution()), nil
}

func (op *executionOp) List(ctx context.Context, params v1.ListExecutionsParams) (*v1.ListExecutionsOK, error) {
	return errorFromDecodedResponse[v1.ListExecutionsOK]("Execution.List", func() (any, error) { return op.client.ListExecutions(ctx, params) })
}

func (op *executionOp) Read(ctx context.Context, workflowID, executionID string) (*v1.ReadExecutionOKExecution, error) {
	res, err := errorFromDecodedResponse[v1.ReadExecutionOK]("Execution.Read", func() (any, error) {
		return op.client.ReadExecution(ctx, v1.ReadExecutionParams{ID: workflowID, ExecutionId: executionID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetExecution()), nil
}

func (op *executionOp) Cancel(ctx context.Context, workflowID, executionID string) (*v1.CancelExecutionOKExecution, error) {
	res, err := errorFromDecodedResponse[v1.CancelExecutionOK]("Execution.Cancel", func() (any, error) {
		return op.client.CancelExecution(ctx, v1.CancelExecutionParams{ID: workflowID, ExecutionId: executionID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetExecution()), nil
}

func (op *executionOp) Delete(ctx context.Context, workflowID, executionID string) error {
	return errorFromDecodedResponseE[v1.DeleteExecutionOK]("Execution.Delete", func() (any, error) {
		return op.client.DeleteExecution(ctx, v1.DeleteExecutionParams{ID: workflowID, ExecutionId: executionID})
	})
}

func (op *executionOp) ListHistory(ctx context.Context, params v1.ListExecutionHistoriesParams) (*v1.ListExecutionHistoriesOK, error) {
	return errorFromDecodedResponse[v1.ListExecutionHistoriesOK]("Execution.ListHistory", func() (any, error) { return op.client.ListExecutionHistories(ctx, params) })
}
