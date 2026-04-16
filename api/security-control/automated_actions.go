// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

import (
	"context"
	"errors"

	v1 "github.com/sacloud/security-control-api-go/apis/v1"
)

type AutomatedActionsAPI interface {
	List(ctx context.Context, params v1.AutomatedActionsListParams) (*v1.AutomatedActionsListOK, error)
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

func (op *automatedActionsOp) List(ctx context.Context, params v1.AutomatedActionsListParams) (*v1.AutomatedActionsListOK, error) {
	const methodName = "AutomatedActions.List"

	res, err := op.client.AutomatedActionsList(ctx, params)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AutomatedActionsListOK:
		return r, nil
	case *v1.AutomatedActionsListApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsListApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsListApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, NewAPIError(methodName, r.StatusCode, errors.New(r.Response.Detail.Value))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *automatedActionsOp) Create(ctx context.Context, req *v1.AutomatedActionInput) (*v1.AutomatedActionOutput, error) {
	const methodName = "AutomatedActions.Create"
	res, err := op.client.AutomatedActionsCreate(ctx, req)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AutomatedActionOutput:
		return r, nil
	case *v1.AutomatedActionsCreateApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsCreateApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsCreateApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsCreateApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, NewAPIError(methodName, r.StatusCode, errors.New(r.Response.Detail.Value))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *automatedActionsOp) Read(ctx context.Context, id string) (*v1.AutomatedActionOutput, error) {
	const methodName = "AutomatedActions.Read"

	res, err := op.client.AutomatedActionsRead(ctx, v1.AutomatedActionsReadParams{ActionId: v1.AutomatedActionID(id)})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AutomatedActionOutput:
		return r, nil
	case *v1.AutomatedActionsReadApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsReadApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsReadApplicationJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsReadApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, NewAPIError(methodName, r.StatusCode, errors.New(r.Response.Detail.Value))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *automatedActionsOp) Update(ctx context.Context, id string, req *v1.AutomatedActionInput) (*v1.AutomatedActionOutput, error) {
	const methodName = "AutomatedActions.Update"

	res, err := op.client.AutomatedActionsUpdate(ctx, req, v1.AutomatedActionsUpdateParams{ActionId: v1.AutomatedActionID(id)})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AutomatedActionOutput:
		return r, nil
	case *v1.AutomatedActionsUpdateApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsUpdateApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsUpdateApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsUpdateApplicationJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsUpdateApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, NewAPIError(methodName, r.StatusCode, errors.New(r.Response.Detail.Value))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *automatedActionsOp) Delete(ctx context.Context, id string) error {
	const methodName = "AutomatedActions.Delete"

	res, err := op.client.AutomatedActionsDelete(ctx, v1.AutomatedActionsDeleteParams{ActionId: v1.AutomatedActionID(id)})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AutomatedActionsDeleteNoContent:
		return nil
	case *v1.AutomatedActionsDeleteApplicationJSONUnauthorized:
		return NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsDeleteApplicationJSONForbidden:
		return NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.AutomatedActionsDeleteApplicationJSONTooManyRequests:
		return NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return NewAPIError(methodName, r.StatusCode, errors.New(r.Response.Detail.Value))
	default:
		return NewAPIError(methodName, 0, err)
	}
}
