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
	res, err := errorFromDecodedResponse[v1.CreateWorkflowCreated]("Workflow.Create", func() (any, error) { return op.client.CreateWorkflow(ctx, &req) })
	if err != nil {
		return nil, err
	}

	return new(res.GetWorkflow()), nil
}

func (op *workflowOp) List(ctx context.Context, params v1.ListWorkflowsParams) (*v1.ListWorkflowsOK, error) {
	return errorFromDecodedResponse[v1.ListWorkflowsOK]("Workflow.List", func() (any, error) { return op.client.ListWorkflows(ctx, params) })
}

func (op *workflowOp) ListSuggest(ctx context.Context, params v1.ListSuggestionsParams) (*v1.ListSuggestionsOK, error) {
	return errorFromDecodedResponse[v1.ListSuggestionsOK]("Workflow.ListSuggest", func() (any, error) { return op.client.ListSuggestions(ctx, params) })
}

func (op *workflowOp) Read(ctx context.Context, id string) (*v1.ReadWorkflowOKWorkflow, error) {
	res, err := errorFromDecodedResponse[v1.ReadWorkflowOK]("Workflow.Read", func() (any, error) { return op.client.ReadWorkflow(ctx, v1.ReadWorkflowParams{ID: id}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetWorkflow()), nil
}

func (op *workflowOp) Update(ctx context.Context, id string, req v1.PatchWorkflowReq) (*v1.PatchWorkflowOKWorkflow, error) {
	res, err := errorFromDecodedResponse[v1.PatchWorkflowOK]("Workflow.Update", func() (any, error) { return op.client.PatchWorkflow(ctx, &req, v1.PatchWorkflowParams{ID: id}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetWorkflow()), nil
}

func (op *workflowOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[v1.DeleteWorkflowOK]("Workflow.Delete", func() (any, error) { return op.client.DeleteWorkflow(ctx, v1.DeleteWorkflowParams{ID: id}) })
}
