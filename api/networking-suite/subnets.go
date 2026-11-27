// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
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
	res, err := errorFromDecodedResponse[v1.SubnetListResponse]("Subnets.List", func() (any, error) {
		return op.client.ListSubnets(ctx, v1.ListSubnetsParams{SubnetGroupSRN: subnetGroupSRN.String()})
	})
	if err != nil {
		return nil, err
	}

	return res.GetSubnets(), nil
}

func (op *subnetsOp) Create(ctx context.Context, req v1.CreateSubnet) (*v1.ReadSubnet, error) {
	res, err := errorFromDecodedResponse[v1.CreatedSubnetResponse]("Subnets.Create", func() (any, error) { return op.client.CreateSubnet(ctx, &req) })
	if err != nil {
		return nil, err
	}

	return new(res.GetSubnet()), nil
}

func (op *subnetsOp) Read(ctx context.Context, s srn.SRN) (*v1.ReadSubnet, error) {
	res, err := errorFromDecodedResponse[v1.SubnetResponse]("Subnets.Read", func() (any, error) { return op.client.ReadSubnet(ctx, v1.ReadSubnetParams{Subnetid: s.ID}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetSubnet()), nil
}

func (op *subnetsOp) Update(ctx context.Context, s srn.SRN, req v1.UpdateSubnet) (*v1.ReadSubnet, error) {
	res, err := errorFromDecodedResponse[v1.SubnetResponse]("Subnets.Update", func() (any, error) { return op.client.UpdateSubnet(ctx, &req, v1.UpdateSubnetParams{Subnetid: s.ID}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetSubnet()), nil
}

func (op *subnetsOp) Delete(ctx context.Context, s srn.SRN) error {
	return errorFromDecodedResponseE[v1.DeleteSubnetNoContent]("Subnets.Delete", func() (any, error) { return op.client.DeleteSubnet(ctx, v1.DeleteSubnetParams{Subnetid: s.ID}) })
}

type CreateInterfaceConnectionParams struct {
	InterfaceSRN srn.SRN
	SubnetSRN    srn.SRN
	IPAddress    *string
}

func (op *interfaceConnectionOp) Create(ctx context.Context, params CreateInterfaceConnectionParams) (*v1.ReadInterfaceConnection, error) {
	res, err := errorFromDecodedResponse[v1.CreatedInterfaceConnectionResponse]("InterfaceConnection.Create", func() (any, error) {
		return op.client.CreateInterfaceConnection(ctx, &v1.CreateInterfaceConnection{
			EphemeralIPv4Address: into.Opt[v1.OptString](params.IPAddress),
			Interface:            v1.SakuraResourceNameRef{SRN: params.InterfaceSRN.String()},
			Subnet:               v1.SakuraResourceNameRef{SRN: params.SubnetSRN.String()},
		})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetInterfaceConnection()), nil
}

func (op *interfaceConnectionOp) Delete(ctx context.Context, s srn.SRN) error {
	return errorFromDecodedResponseE[v1.DeleteInterfaceConnectionNoContent]("InterfaceConnection.Delete", func() (any, error) {
		return op.client.DeleteInterfaceConnection(ctx, v1.DeleteInterfaceConnectionParams{
			Interfaceconnectionid: s.ID,
		})
	})
}
