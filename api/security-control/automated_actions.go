// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

import (
	"context"
	"errors"

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
	const methodName = "AutomatedActions.List"

	res, err := op.client.ListAutomatedActions(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListAutomatedActionsOK:
		return r, nil
	case *v1.ListAutomatedActionsApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ListAutomatedActionsApplicationProblemJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ListAutomatedActionsApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ListAutomatedActionsApplicationProblemJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ListAutomatedActionsApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ListAutomatedActionsApplicationProblemJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, responseError(methodName, r)
	default:
		return nil, responseError(methodName, res)
	}
}

func (op *automatedActionsOp) Create(ctx context.Context, req *v1.AutomatedActionInput) (*v1.AutomatedActionOutput, error) {
	const methodName = "AutomatedActions.Create"
	res, err := op.client.CreateAutomatedAction(ctx, req)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AutomatedActionOutput:
		return r, nil
	case *v1.CreateAutomatedActionApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateAutomatedActionApplicationProblemJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateAutomatedActionApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateAutomatedActionApplicationProblemJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateAutomatedActionApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateAutomatedActionApplicationProblemJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateAutomatedActionApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.CreateAutomatedActionApplicationProblemJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, responseError(methodName, r)
	default:
		return nil, responseError(methodName, res)
	}
}

func (op *automatedActionsOp) Read(ctx context.Context, id string) (*v1.AutomatedActionOutput, error) {
	const methodName = "AutomatedActions.Read"

	res, err := op.client.ReadAutomatedAction(ctx, v1.ReadAutomatedActionParams{ActionId: v1.AutomatedActionID(id)})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AutomatedActionOutput:
		return r, nil
	case *v1.ReadAutomatedActionApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadAutomatedActionApplicationProblemJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadAutomatedActionApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadAutomatedActionApplicationProblemJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadAutomatedActionApplicationJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadAutomatedActionApplicationProblemJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadAutomatedActionApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.ReadAutomatedActionApplicationProblemJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, responseError(methodName, r)
	default:
		return nil, responseError(methodName, res)
	}
}

func (op *automatedActionsOp) Update(ctx context.Context, id string, req *v1.AutomatedActionInput) (*v1.AutomatedActionOutput, error) {
	const methodName = "AutomatedActions.Update"

	res, err := op.client.UpdateAutomatedAction(ctx, req, v1.UpdateAutomatedActionParams{ActionId: v1.AutomatedActionID(id)})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AutomatedActionOutput:
		return r, nil
	case *v1.UpdateAutomatedActionApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationProblemJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationProblemJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationProblemJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationProblemJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UpdateAutomatedActionApplicationProblemJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, responseError(methodName, r)
	default:
		return nil, responseError(methodName, res)
	}
}

func (op *automatedActionsOp) Delete(ctx context.Context, id string) error {
	const methodName = "AutomatedActions.Delete"

	res, err := op.client.DeleteAutomatedAction(ctx, v1.DeleteAutomatedActionParams{ActionId: v1.AutomatedActionID(id)})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteAutomatedActionNoContent:
		return nil
	case *v1.DeleteAutomatedActionApplicationJSONUnauthorized:
		return NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.DeleteAutomatedActionApplicationProblemJSONUnauthorized:
		return responseError(methodName, r)
	case *v1.DeleteAutomatedActionApplicationJSONForbidden:
		return NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.DeleteAutomatedActionApplicationProblemJSONForbidden:
		return responseError(methodName, r)
	case *v1.DeleteAutomatedActionApplicationJSONTooManyRequests:
		return NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.DeleteAutomatedActionApplicationProblemJSONTooManyRequests:
		return responseError(methodName, r)
	default:
		return responseError(methodName, res)
	}
}
