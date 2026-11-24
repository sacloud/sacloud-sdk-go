// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

import (
	"context"
	"errors"

	v1 "github.com/sacloud/sacloud-sdk-go/api/security-control/apis/v1"
)

type ActivationAPI interface {
	Create(ctx context.Context, servicePrincipalID string) (*v1.ActivationOutput, error)
	Read(ctx context.Context) (*v1.ActivationOutput, error)
	Update(ctx context.Context, servicePrincipalID string, isActive bool) (*v1.ActivationOutput, error)
}

// Deprecated: use ActivationAPI instead
type ActivationsAPI interface {
	ActivationAPI
}

var _ ActivationAPI = (*activationOp)(nil)
var _ ActivationsAPI = (*activationOp)(nil)

type activationOp struct {
	client *v1.Client
}

func NewActivationOp(client *v1.Client) ActivationAPI {
	return &activationOp{client: client}
}

func NewActivationsOp(client *v1.Client) ActivationAPI {
	return &activationOp{client: client}
}

func (op *activationOp) Create(ctx context.Context, servicePrincipalID string) (*v1.ActivationOutput, error) {
	const methodName = "Activation.Create"
	res, err := op.client.CreateActivation(ctx, &v1.ActivationCreateInput{ServicePrincipalId: servicePrincipalID})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ActivationOutput:
		return r, nil
	case *v1.CreateActivationApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateActivationApplicationProblemJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateActivationApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateActivationApplicationProblemJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateActivationApplicationJSONConflict:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateActivationApplicationProblemJSONConflict:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateActivationApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateActivationApplicationProblemJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.CreateActivationApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.CreateActivationApplicationProblemJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, responseError(methodName, r)
	default:
		return nil, responseError(methodName, res)
	}
}

func (op *activationOp) Read(ctx context.Context) (*v1.ActivationOutput, error) {
	const methodName = "Activation.Read"

	res, err := op.client.ReadActivation(ctx)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ActivationOutput:
		return r, nil
	case *v1.ReadActivationApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadActivationApplicationProblemJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadActivationApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadActivationApplicationProblemJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadActivationApplicationJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadActivationApplicationProblemJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.ReadActivationApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.ReadActivationApplicationProblemJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, responseError(methodName, r)
	default:
		return nil, responseError(methodName, res)
	}
}

func (op *activationOp) Update(ctx context.Context, servicePrincipalID string, isActive bool) (*v1.ActivationOutput, error) {
	const methodName = "Activation.Update"

	res, err := op.client.UpdateActivation(ctx, &v1.ActivationUpdateInput{
		ServicePrincipalId: servicePrincipalID,
		IsActive:           isActive,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ActivationOutput:
		return r, nil
	case *v1.UpdateActivationApplicationJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateActivationApplicationProblemJSONBadRequest:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateActivationApplicationJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateActivationApplicationProblemJSONUnauthorized:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateActivationApplicationJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateActivationApplicationProblemJSONForbidden:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateActivationApplicationJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateActivationApplicationProblemJSONNotFound:
		return nil, NewAPIError(methodName, int(r.Status.Value), errors.New(string(r.Detail.Value)))
	case *v1.UpdateActivationApplicationJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UpdateActivationApplicationProblemJSONTooManyRequests:
		return nil, NewAPIError(methodName, int(r.Response.Status.Value), errors.New(string(r.Response.Detail.Value)))
	case *v1.UnexpectedErrorStatusCode:
		return nil, responseError(methodName, r)
	default:
		return nil, responseError(methodName, res)
	}
}
