// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite

import (
	"context"
	"errors"

	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
	"github.com/sacloud/sacloud-sdk-go/srn"
)

type AddressesAPI interface {
	List(ctx context.Context, params ListAddressesParams) ([]v1.Address, error)
	Create(ctx context.Context, params CreateAddressParams) (*v1.Address, error)
	Read(ctx context.Context, addr srn.SRN) (*v1.Address, error)
	Delete(ctx context.Context, addr srn.SRN) error
}

var _ AddressesAPI = (*addressesOp)(nil)

type addressesOp struct {
	client *v1.Client
}

func NewAddressesOp(client *v1.Client) AddressesAPI {
	return &addressesOp{client: client}
}

type ListAddressesParams struct {
	SubnetSRN              srn.SRN
	InterfaceConnectionSRN *srn.SRN
}

type CreateAddressParams struct {
	SubnetSRN srn.SRN
	IPAddress *string
}

func (op *addressesOp) List(ctx context.Context, params ListAddressesParams) ([]v1.Address, error) {
	const methodName = "Addresses.List"

	var icParam v1.OptString
	if params.InterfaceConnectionSRN != nil {
		icParam.SetTo(params.InterfaceConnectionSRN.String())
	}
	res, err := op.client.ListAddresses(ctx, v1.ListAddressesParams{
		SubnetSRN:              params.SubnetSRN.String(),
		InterfaceConnectionSRN: icParam,
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	return res.Addresses, nil
}

func (op *addressesOp) Create(ctx context.Context, params CreateAddressParams) (*v1.Address, error) {
	const methodName = "Addresses.Create"

	res, err := op.client.CreateAddress(ctx, &v1.CreateAddress{
		AddressType: v1.CreateAddressAddressTypeSERVICEENDPOINTGATEWAYADDRESS,
		IPVersion:   v1.CreateAddressIPVersionIPv4,
		IPAddress:   into.Opt[v1.OptString](params.IPAddress),
		Subnet:      v1.SakuraResourceNameRef{SRN: params.SubnetSRN.String()},
	})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreatedSubnetAddressResponse:
		return &r.Address, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *addressesOp) Read(ctx context.Context, addr srn.SRN) (*v1.Address, error) {
	const methodName = "Addresses.Read"

	res, err := op.client.ReadAddress(ctx, v1.ReadAddressParams{Addressid: addr.ID})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.AddressResponse:
		return &r.Address, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *addressesOp) Delete(ctx context.Context, addr srn.SRN) error {
	const methodName = "Addresses.Delete"

	res, err := op.client.DeleteAddress(ctx, v1.DeleteAddressParams{Addressid: addr.ID})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteAddressNoContent:
		return nil
	case *v1.ApiErrorStatusCode:
		return newGeneratedAPIError(methodName, r)
	default:
		return NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
