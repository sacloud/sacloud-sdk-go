// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/srn"
)

type AddressesAPI interface {
	List(ctx context.Context, params ListAddressesParams) ([]v1.ReadAddress, error)
	Read(ctx context.Context, addr srn.SRN) (*v1.ReadAddress, error)
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

func (op *addressesOp) List(ctx context.Context, params ListAddressesParams) ([]v1.ReadAddress, error) {
	var icParam v1.OptString
	if params.InterfaceConnectionSRN != nil {
		icParam.SetTo(params.InterfaceConnectionSRN.String())
	}
	res, err := errorFromDecodedResponse[v1.AddressListResponse]("Addresses.List", func() (any, error) {
		return op.client.ListAddresses(ctx, v1.ListAddressesParams{
			SubnetSRN:              params.SubnetSRN.String(),
			InterfaceConnectionSRN: icParam,
		})
	})
	if err != nil {
		return nil, err
	}

	return res.GetAddresses(), nil
}

func (op *addressesOp) Read(ctx context.Context, addr srn.SRN) (*v1.ReadAddress, error) {
	res, err := errorFromDecodedResponse[v1.AddressResponse]("Addresses.Read", func() (any, error) { return op.client.ReadAddress(ctx, v1.ReadAddressParams{Addressid: addr.ID}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetAddress()), nil
}
