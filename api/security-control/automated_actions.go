// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/security-control/apis/v1"
)

type AutomatedActionsAPI interface {
	List(ctx context.Context, params v1.ListAutomatedActionsParams) (*v1.ListAutomatedActionsOK, error)
	Create(ctx context.Context, req *v1.AutomatedActionInput) (*v1.AutomatedActionOutput, error)
	Read(ctx context.Context, id string) (*v1.AutomatedActionOutput, error)
	Update(ctx context.Context, id string, req *v1.AutomatedActionInput) (*v1.AutomatedActionOutput, error)
	Delete(ctx context.Context, id string) error
}

var _ AutomatedActionsAPI = (*automatedActionsOp)(nil)

type automatedActionsOp struct {
	client *v1.Client
}

func NewAutomatedActionsOp(client *v1.Client) AutomatedActionsAPI {
	return &automatedActionsOp{client: client}
}

func (op *automatedActionsOp) List(ctx context.Context, params v1.ListAutomatedActionsParams) (*v1.ListAutomatedActionsOK, error) {
	return errorFromDecodedResponse[v1.ListAutomatedActionsOK]("AutomatedActions.List", func() (any, error) { return op.client.ListAutomatedActions(ctx, params) })
}

func (op *automatedActionsOp) Create(ctx context.Context, req *v1.AutomatedActionInput) (*v1.AutomatedActionOutput, error) {
	return errorFromDecodedResponse[v1.AutomatedActionOutput]("AutomatedActions.Create", func() (any, error) { return op.client.CreateAutomatedAction(ctx, req) })
}

func (op *automatedActionsOp) Read(ctx context.Context, id string) (*v1.AutomatedActionOutput, error) {
	return errorFromDecodedResponse[v1.AutomatedActionOutput]("AutomatedActions.Read", func() (any, error) {
		return op.client.ReadAutomatedAction(ctx, v1.ReadAutomatedActionParams{ActionId: v1.AutomatedActionID(id)})
	})
}

func (op *automatedActionsOp) Update(ctx context.Context, id string, req *v1.AutomatedActionInput) (*v1.AutomatedActionOutput, error) {
	return errorFromDecodedResponse[v1.AutomatedActionOutput]("AutomatedActions.Update", func() (any, error) {
		return op.client.UpdateAutomatedAction(ctx, req, v1.UpdateAutomatedActionParams{ActionId: v1.AutomatedActionID(id)})
	})
}

func (op *automatedActionsOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[v1.DeleteAutomatedActionNoContent]("AutomatedActions.Delete", func() (any, error) {
		return op.client.DeleteAutomatedAction(ctx, v1.DeleteAutomatedActionParams{ActionId: v1.AutomatedActionID(id)})
	})
}
