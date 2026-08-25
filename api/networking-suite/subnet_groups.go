// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite

import (
	"context"
	"errors"

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
	const methodName = "SubnetGroups.List"

	res, err := op.client.ListSubnetGroups(ctx)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	return res.SubnetGroups, nil
}

func (op *subnetGroupsOp) Create(ctx context.Context, req v1.CreateSubnetGroup) (*v1.ReadSubnetGroup, error) {
	const methodName = "SubnetGroups.Create"

	res, err := op.client.CreateSubnetGroup(ctx, &req)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreatedSubnetGroupResponse:
		sg := r.GetSubnetGroup()
		return &sg, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subnetGroupsOp) Read(ctx context.Context, s srn.SRN) (*v1.ReadSubnetGroup, error) {
	const methodName = "SubnetGroups.Read"

	res, err := op.client.ReadSubnetGroup(ctx, v1.ReadSubnetGroupParams{Subnetgroupid: s.ID})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.SubnetGroupResponse:
		sg := r.GetSubnetGroup()
		return &sg, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subnetGroupsOp) Update(ctx context.Context, s srn.SRN, req v1.UpdateSubnetGroup) (*v1.ReadSubnetGroup, error) {
	const methodName = "SubnetGroups.Update"

	res, err := op.client.UpdateSubnetGroup(ctx, &req, v1.UpdateSubnetGroupParams{Subnetgroupid: s.ID})
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.SubnetGroupResponse:
		sg := r.GetSubnetGroup()
		return &sg, nil
	case *v1.ApiErrorStatusCode:
		return nil, newGeneratedAPIError(methodName, r)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subnetGroupsOp) Delete(ctx context.Context, s srn.SRN) error {
	const methodName = "SubnetGroups.Delete"

	res, err := op.client.DeleteSubnetGroup(ctx, v1.DeleteSubnetGroupParams{Subnetgroupid: s.ID})
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteSubnetGroupNoContent:
		return nil
	case *v1.ApiErrorStatusCode:
		return newGeneratedAPIError(methodName, r)
	default:
		return NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func newGeneratedAPIError(methodName string, errRes *v1.ApiErrorStatusCode) error {
	msg := errRes.Response.ErrorMsg.Or("unknown error")
	if code := errRes.Response.ErrorCode.Or(""); code != "" {
		msg = code + ": " + msg
	}

	return NewAPIError(methodName, errRes.StatusCode, errors.New(msg))
}
