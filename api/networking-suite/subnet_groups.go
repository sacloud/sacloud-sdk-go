// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/srn"
)

type SubnetGroupsAPI interface {
	List(ctx context.Context) ([]v1.ReadSubnetGroup, error)
	Create(ctx context.Context, req v1.CreateSubnetGroup) (*v1.ReadSubnetGroup, error)
	Read(ctx context.Context, s srn.SRN) (*v1.ReadSubnetGroup, error)
	Update(ctx context.Context, s srn.SRN, req v1.UpdateSubnetGroup) (*v1.ReadSubnetGroup, error)
	Delete(ctx context.Context, s srn.SRN) error
}

var _ SubnetGroupsAPI = (*subnetGroupsOp)(nil)

type subnetGroupsOp struct {
	client *v1.Client
}

func NewSubnetGroupsOp(client *v1.Client) SubnetGroupsAPI {
	return &subnetGroupsOp{client: client}
}

func (op *subnetGroupsOp) List(ctx context.Context) ([]v1.ReadSubnetGroup, error) {
	res, err := errorFromDecodedResponse[v1.SubnetGroupListResponse]("SubnetGroups.List", func() (any, error) { return op.client.ListSubnetGroups(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetSubnetGroups(), nil
}

func (op *subnetGroupsOp) Create(ctx context.Context, req v1.CreateSubnetGroup) (*v1.ReadSubnetGroup, error) {
	res, err := errorFromDecodedResponse[v1.CreatedSubnetGroupResponse]("SubnetGroups.Create", func() (any, error) { return op.client.CreateSubnetGroup(ctx, &req) })
	if err != nil {
		return nil, err
	}

	return new(res.GetSubnetGroup()), nil
}

func (op *subnetGroupsOp) Read(ctx context.Context, s srn.SRN) (*v1.ReadSubnetGroup, error) {
	res, err := errorFromDecodedResponse[v1.SubnetGroupResponse]("SubnetGroups.Read", func() (any, error) {
		return op.client.ReadSubnetGroup(ctx, v1.ReadSubnetGroupParams{Subnetgroupid: s.ID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetSubnetGroup()), nil
}

func (op *subnetGroupsOp) Update(ctx context.Context, s srn.SRN, req v1.UpdateSubnetGroup) (*v1.ReadSubnetGroup, error) {
	res, err := errorFromDecodedResponse[v1.SubnetGroupResponse]("SubnetGroups.Update", func() (any, error) {
		return op.client.UpdateSubnetGroup(ctx, &req, v1.UpdateSubnetGroupParams{Subnetgroupid: s.ID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetSubnetGroup()), nil
}

func (op *subnetGroupsOp) Delete(ctx context.Context, s srn.SRN) error {
	return errorFromDecodedResponseE[v1.DeleteSubnetGroupNoContent]("SubnetGroups.Delete", func() (any, error) {
		return op.client.DeleteSubnetGroup(ctx, v1.DeleteSubnetGroupParams{Subnetgroupid: s.ID})
	})
}
