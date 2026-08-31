// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite

import (
	"context"
	"errors"

	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/srn"
)

type SubnetsAPI interface {
	List(ctx context.Context, subnetGroupSRN srn.SRN) ([]v1.ReadSubnet, error)
	Create(ctx context.Context, req v1.CreateSubnet) (*v1.ReadSubnet, error)
	Read(ctx context.Context, s srn.SRN) (*v1.ReadSubnet, error)
	Update(ctx context.Context, s srn.SRN, req v1.UpdateSubnet) (*v1.ReadSubnet, error)
	Delete(ctx context.Context, s srn.SRN) error
}

type InterfaceConnectionAPI interface {
	Create(ctx context.Context, params CreateInterfaceConnectionParams) (*v1.ReadInterfaceConnection, error)
	Delete(ctx context.Context, s srn.SRN) error
}

var _ SubnetsAPI = (*subnetsOp)(nil)
var _ InterfaceConnectionAPI = (*interfaceConnectionOp)(nil)

type subnetsOp struct {
	client *v1.Client
}

type interfaceConnectionOp struct {
	client *v1.Client
}

func NewSubnetsOp(client *v1.Client) SubnetsAPI {
	return &subnetsOp{client: client}
}

func NewInterfaceConnectionOp(client *v1.Client) InterfaceConnectionAPI {
	return &interfaceConnectionOp{client: client}
}

func (op *subnetsOp) List(ctx context.Context, subnetGroupSRN srn.SRN) ([]v1.ReadSubnet, error) {
	const methodName = "Subnets.List"

	res, err := op.client.ListSubnets(ctx, v1.ListSubnetsParams{SubnetGroupSRN: subnetGroupSRN.String()})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	return res.Subnets, nil
}

func (op *subnetsOp) Create(ctx context.Context, req v1.CreateSubnet) (*v1.ReadSubnet, error) {
	const methodName = "Subnets.Create"

	res, err := op.client.CreateSubnet(ctx, &req)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreatedSubnetResponse:
		return &r.Subnet, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subnetsOp) Read(ctx context.Context, s srn.SRN) (*v1.ReadSubnet, error) {
	const methodName = "Subnets.Read"

	res, err := op.client.ReadSubnet(ctx, v1.ReadSubnetParams{
		Subnetid: s.ID,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.SubnetResponse:
		return &r.Subnet, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subnetsOp) Update(ctx context.Context, s srn.SRN, req v1.UpdateSubnet) (*v1.ReadSubnet, error) {
	const methodName = "Subnets.Update"

	res, err := op.client.UpdateSubnet(ctx, &req, v1.UpdateSubnetParams{
		Subnetid: s.ID,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.SubnetResponse:
		return &r.Subnet, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subnetsOp) Delete(ctx context.Context, s srn.SRN) error {
	const methodName = "Subnets.Delete"

	res, err := op.client.DeleteSubnet(ctx, v1.DeleteSubnetParams{
		Subnetid: s.ID,
	})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteSubnetNoContent:
		return nil
	case *v1.ApiErrorStatusCode:
		return newGeneratedAPIError(methodName, r)
	default:
		return NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

type CreateInterfaceConnectionParams struct {
	InterfaceSRN srn.SRN
	SubnetSRN    srn.SRN
	IPAddress    *string
}

func (op *interfaceConnectionOp) Create(ctx context.Context, params CreateInterfaceConnectionParams) (*v1.ReadInterfaceConnection, error) {
	const methodName = "InterfaceConnection.Create"

	var ia v1.OptString
	if params.IPAddress != nil {
		ia.SetTo(*params.IPAddress)
	}
	res, err := op.client.CreateInterfaceConnection(ctx, &v1.CreateInterfaceConnection{
		EphemeralIPv4Address: ia,
		Interface:            v1.SakuraResourceNameRef{SRN: params.InterfaceSRN.String()},
		Subnet:               v1.SakuraResourceNameRef{SRN: params.SubnetSRN.String()},
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreatedInterfaceConnectionResponse:
		return &r.InterfaceConnection, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *interfaceConnectionOp) Delete(ctx context.Context, s srn.SRN) error {
	const methodName = "InterfaceConnection.Delete"

	res, err := op.client.DeleteInterfaceConnection(ctx, v1.DeleteInterfaceConnectionParams{
		Interfaceconnectionid: s.ID,
	})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteInterfaceConnectionNoContent:
		return nil
	case *v1.ApiErrorStatusCode:
		return newGeneratedAPIError(methodName, r)
	default:
		return NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
