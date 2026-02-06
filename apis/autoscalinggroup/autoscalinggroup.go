// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package autoscalinggroup

import (
	"context"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
	"github.com/sacloud/saclient-go"
)

type AutoScalingGroupAPI interface {
	// List returns the list of AutoScalingGroups, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, elems int64, cursor *v1.AutoScalingGroupID) (list []v1.ReadAutoScalingGroupDetail, nextCursor *v1.AutoScalingGroupID, err error)

	// Create creates a new AutoScalingGroup.
	Create(ctx context.Context, params CreateParams) (group *v1.CreatedAutoScalingGroup, err error)

	// Read retrieves an AutoScalingGroup by its ID.
	Read(ctx context.Context, id v1.AutoScalingGroupID) (group *AutoScalingGroupDetail, err error)

	// Delete deletes an AutoScalingGroup by its ID.
	Delete(ctx context.Context, id v1.AutoScalingGroupID) error
}

type AutoScalingGroupOp struct {
	client    *v1.Client
	clusterID v1.ClusterID
}

func NewAutoScalingGroupOp(client *v1.Client, clusterID v1.ClusterID) *AutoScalingGroupOp {
	return &AutoScalingGroupOp{
		client:    client,
		clusterID: clusterID,
	}
}

func (op *AutoScalingGroupOp) List(
	ctx context.Context,
	maxItems int64,
	cursor *v1.AutoScalingGroupID,
) (
	groups []v1.ReadAutoScalingGroupDetail,
	nextCursor *v1.AutoScalingGroupID,
	err error,
) {
	res, err := common.ErrorFromDecodedResponse("AutoScalingGroup.List", func() (*v1.ListAutoScalingGroupResponse, error) {
		return op.client.ListAutoScalingGroups(ctx, v1.ListAutoScalingGroupsParams{
			ClusterID: op.clusterID,
			Cursor:    common.IntoOpt[v1.OptAutoScalingGroupID](cursor),
			MaxItems:  maxItems,
		})
	})

	if res != nil {
		groups = res.AutoScalingGroups
		nextCursor = common.FromOpt(res.NextCursor)
	}

	return
}

func (op *AutoScalingGroupOp) Create(
	ctx context.Context,
	params CreateParams,
) (
	group *v1.CreatedAutoScalingGroup,
	err error,
) {
	res, err := common.ErrorFromDecodedResponse("AutoScalingGroup.Create", func() (*v1.CreateAutoScalingGroupResponse, error) {
		return op.client.CreateAutoScalingGroup(ctx, saclient.Ptr(params.into()), v1.CreateAutoScalingGroupParams{
			ClusterID: op.clusterID,
		})
	})

	if res != nil {
		group = &res.AutoScalingGroup
	}

	return
}

func (op *AutoScalingGroupOp) Read(
	ctx context.Context,
	id v1.AutoScalingGroupID,
) (
	group *AutoScalingGroupDetail,
	err error,
) {
	res, err := common.ErrorFromDecodedResponse("AutoScalingGroup.Read", func() (*v1.GetAutoScalingGroupResponse, error) {
		return op.client.GetAutoScalingGroup(ctx, v1.GetAutoScalingGroupParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: id,
		})
	})

	if res != nil {
		var detail AutoScalingGroupDetail
		detail.from(&res.AutoScalingGroup)
		group = &detail
	}

	return
}

func (op *AutoScalingGroupOp) Delete(ctx context.Context, id v1.AutoScalingGroupID) error {
	return common.ErrorFromDecodedResponseE("AutoScalingGroup.Delete", func() error {
		return op.client.DeleteAutoScalingGroup(ctx, v1.DeleteAutoScalingGroupParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: id,
		})
	})
}

var _ AutoScalingGroupAPI = (*AutoScalingGroupOp)(nil)

type NodeInterface struct {
	InterfaceIndex int16
	Upstream       string
	IpPool         []v1.IpRange
	NetmaskLen     *int16
	DefaultGateway *string
	PacketFilterID *string
	ConnectsToLB   bool
}

func (n *NodeInterface) into() (ret v1.AutoScalingGroupNodeInterface) {
	ret.SetInterfaceIndex(n.InterfaceIndex)
	ret.SetUpstream(n.Upstream)
	ret.SetIpPool(n.IpPool)
	ret.SetNetmaskLen(common.IntoOpt[v1.OptInt16](n.NetmaskLen))
	ret.SetDefaultGateway(common.IntoOpt[v1.OptString](n.DefaultGateway))
	ret.SetPacketFilterID(common.IntoOpt[v1.OptString](n.PacketFilterID))
	ret.SetConnectsToLB(n.ConnectsToLB)

	return
}

func (n *NodeInterface) from(res *v1.AutoScalingGroupNodeInterface) {
	n.InterfaceIndex = res.GetInterfaceIndex()
	n.Upstream = res.GetUpstream()
	n.IpPool = res.GetIpPool()
	n.NetmaskLen = common.FromOpt(res.GetNetmaskLen())
	n.DefaultGateway = common.FromOpt(res.GetDefaultGateway())
	n.PacketFilterID = common.FromOpt(res.GetPacketFilterID())
	n.ConnectsToLB = res.GetConnectsToLB()
}

type CreateParams struct {
	Name                   string
	Zone                   string
	NameServers            []v1.IPv4
	WorkerServiceClassPath string
	MinNodes               int32
	MaxNodes               int32
	Interfaces             []NodeInterface
}

func (c *CreateParams) into() (ret v1.CreateAutoScalingGroup) {
	ret.SetName(c.Name)
	ret.SetZone(c.Zone)
	ret.SetNameServers(c.NameServers)
	ret.SetWorkerServiceClassPath(c.WorkerServiceClassPath)
	ret.SetMinNodes(c.MinNodes)
	ret.SetMaxNodes(c.MaxNodes)
	ret.SetInterfaces(common.MapSlice(c.Interfaces, func(n NodeInterface) v1.AutoScalingGroupNodeInterface { return n.into() }))

	return
}

type AutoScalingGroupDetail struct {
	AutoScalingGroupID     v1.AutoScalingGroupID
	Name                   string
	Zone                   string
	NameServers            []v1.IPv4
	WorkerServiceClassPath string
	MinNodes               int32
	MaxNodes               int32
	WorkerNodeCount        int32
	Deleting               bool
	Interfaces             []NodeInterface
}

func (a *AutoScalingGroupDetail) from(res *v1.ReadAutoScalingGroupDetail) {
	a.AutoScalingGroupID = res.GetAutoScalingGroupID()
	a.Name = res.GetName()
	a.Zone = res.GetZone()
	a.NameServers = res.GetNameServers()
	a.WorkerServiceClassPath = res.GetWorkerServiceClassPath()
	a.MinNodes = res.GetMinNodes()
	a.MaxNodes = res.GetMaxNodes()
	a.WorkerNodeCount = res.GetWorkerNodeCount()
	a.Deleting = res.GetDeleting()
	a.Interfaces = common.MapSlice(res.GetInterfaces(), func(n v1.AutoScalingGroupNodeInterface) (m NodeInterface) {
		m.from(&n)
		return
	})
}
