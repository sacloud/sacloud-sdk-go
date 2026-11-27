// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

import (
	"context"

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
	return errorFromDecodedResponse[v1.ActivationOutput]("Activation.Create", func() (any, error) {
		return op.client.CreateActivation(ctx, &v1.ActivationCreateInput{ServicePrincipalId: servicePrincipalID})
	})
}

func (op *activationOp) Read(ctx context.Context) (*v1.ActivationOutput, error) {
	return errorFromDecodedResponse[v1.ActivationOutput]("Activation.Read", func() (any, error) { return op.client.ReadActivation(ctx) })
}

func (op *activationOp) Update(ctx context.Context, servicePrincipalID string, isActive bool) (*v1.ActivationOutput, error) {
	return errorFromDecodedResponse[v1.ActivationOutput]("Activation.Update", func() (any, error) {
		return op.client.UpdateActivation(ctx, &v1.ActivationUpdateInput{
			ServicePrincipalId: servicePrincipalID,
			IsActive:           isActive,
		})
	})
}
