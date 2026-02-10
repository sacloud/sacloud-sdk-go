// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package loadbalancer

import (
	"context"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
	"github.com/sacloud/saclient-go"
)

type LoadBalancerAPI interface {
	// List returns the list of LoadBalancers, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, elems int64, cursor *v1.LoadBalancerID) (list []v1.ReadLoadBalancerSummary, nextCursor *v1.LoadBalancerID, err error)

	// Create creates a new LoadBalancer.
	Create(ctx context.Context, params CreateParams) (lb *v1.CreatedLoadBalancer, err error)

	// Read retrieves a LoadBalancer by its ID.
	Read(ctx context.Context, id v1.LoadBalancerID) (lb *LoadBalancerDetail, err error)

	// Delete deletes a LoadBalancer by its ID.
	Delete(ctx context.Context, id v1.LoadBalancerID) error

	// ListNode returns the list of LoadBalancerNodes, paginated.
	// Pass nil to `cursor` to get the first page.
	ListNode(ctx context.Context, lbID v1.LoadBalancerID, elems int64, cursor *v1.LoadBalancerID) (list []v1.ReadLoadBalancerNodeSummary, err error)

	// ReadNode retrieves a LoadBalancerNode by its ID.
	ReadNode(ctx context.Context, lbID v1.LoadBalancerID, nodeID v1.LoadBalancerNodeID) (node *LoadBalancerNodeDetail, err error)
}

type LoadBalancerOp struct {
	client             *v1.Client
	clusterID          v1.ClusterID
	autoScalingGroupID v1.AutoScalingGroupID
}

func NewLoadBalancerOp(client *v1.Client, clusterID v1.ClusterID, autoScalingGroupID v1.AutoScalingGroupID) *LoadBalancerOp {
	return &LoadBalancerOp{
		client:             client,
		clusterID:          clusterID,
		autoScalingGroupID: autoScalingGroupID,
	}
}

func (op *LoadBalancerOp) List(
	ctx context.Context,
	maxItems int64,
	cursor *v1.LoadBalancerID,
) (
	lbs []v1.ReadLoadBalancerSummary,
	nextCursor *v1.LoadBalancerID,
	err error,
) {
	res, err := common.ErrorFromDecodedResponse("LoadBalancer.List", func() (*v1.ListLoadBalancersResponse, error) {
		return op.client.ListLoadBalancers(ctx, v1.ListLoadBalancersParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
			Cursor:             common.IntoOpt[v1.OptLoadBalancerID](cursor),
			MaxItems:           maxItems,
		})
	})

	if res != nil {
		lbs = res.LoadBalancers
		nextCursor = common.FromOpt(res.NextCursor)
	}

	return
}

func (op *LoadBalancerOp) Create(
	ctx context.Context,
	params CreateParams,
) (
	lb *v1.CreatedLoadBalancer,
	err error,
) {
	res, err := common.ErrorFromDecodedResponse("LoadBalancer.Create", func() (*v1.CreateLoadBalancerResponse, error) {
		return op.client.CreateLoadBalancer(ctx, saclient.Ptr(params.into()), v1.CreateLoadBalancerParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
		})
	})

	if res != nil {
		lb = &res.LoadBalancer
	}

	return
}

func (op *LoadBalancerOp) Read(
	ctx context.Context,
	id v1.LoadBalancerID,
) (
	lb *LoadBalancerDetail,
	err error,
) {
	res, err := common.ErrorFromDecodedResponse("LoadBalancer.Read", func() (*v1.GetLoadBalancerResponse, error) {
		return op.client.GetLoadBalancer(ctx, v1.GetLoadBalancerParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
			LoadBalancerID:     id,
		})
	})

	if res != nil {
		var detail LoadBalancerDetail
		detail.from(&res.LoadBalancer)
		lb = &detail
	}

	return
}

func (op *LoadBalancerOp) Delete(ctx context.Context, id v1.LoadBalancerID) error {
	return common.ErrorFromDecodedResponseE("LoadBalancer.Delete", func() error {
		return op.client.DeleteLoadBalancer(ctx, v1.DeleteLoadBalancerParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
			LoadBalancerID:     id,
		})
	})
}

func (op *LoadBalancerOp) ListNode(
	ctx context.Context,
	lbID v1.LoadBalancerID,
	maxItems int64,
	cursor *v1.LoadBalancerID,
) (
	nodes []v1.ReadLoadBalancerNodeSummary,
	err error,
) {
	res, err := common.ErrorFromDecodedResponse("LoadBalancer.ListNode", func() (*v1.ListLoadBalancerNodesResponse, error) {
		return op.client.ListLoadBalancerNodes(ctx, v1.ListLoadBalancerNodesParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
			LoadBalancerID:     lbID,
			Cursor:             common.IntoOpt[v1.OptLoadBalancerID](cursor),
			MaxItems:           maxItems,
		})
	})

	if res != nil {
		nodes = res.LoadBalancerNodes
	}

	return
}

func (op *LoadBalancerOp) ReadNode(
	ctx context.Context,
	lbID v1.LoadBalancerID,
	nodeID v1.LoadBalancerNodeID,
) (
	node *LoadBalancerNodeDetail,
	err error,
) {
	res, err := common.ErrorFromDecodedResponse("LoadBalancer.ReadNode", func() (*v1.GetLoadBalancerNodeResponse, error) {
		return op.client.GetLoadBalancerNode(ctx, v1.GetLoadBalancerNodeParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
			LoadBalancerID:     lbID,
			LoadBalancerNodeID: nodeID,
		})
	})

	if res != nil {
		var detail LoadBalancerNodeDetail
		detail.from(&res.LoadBalancerNode)
		node = &detail
	}

	return
}

var _ LoadBalancerAPI = (*LoadBalancerOp)(nil)

type LoadBalancerInterface struct {
	InterfaceIndex  int16
	Upstream        string
	IpPool          []v1.IpRange
	NetmaskLen      *int16
	DefaultGateway  *string
	Vip             *string
	VirtualRouterID *int16
	PacketFilterID  *string
}

func (l LoadBalancerInterface) into() (ret v1.LoadBalancerInterface) {
	ret.SetInterfaceIndex(l.InterfaceIndex)
	ret.SetUpstream(l.Upstream)
	ret.SetIpPool(l.IpPool)
	ret.SetNetmaskLen(common.IntoOpt[v1.OptInt16](l.NetmaskLen))
	ret.SetDefaultGateway(common.IntoOpt[v1.OptString](l.DefaultGateway))
	ret.SetVip(common.IntoOpt[v1.OptString](l.Vip))
	ret.SetVirtualRouterID(common.IntoOpt[v1.OptInt16](l.VirtualRouterID))
	ret.SetPacketFilterID(common.IntoOpt[v1.OptString](l.PacketFilterID))

	return
}

func (l *LoadBalancerInterface) From(res *v1.LoadBalancerInterface) {
	l.InterfaceIndex = res.GetInterfaceIndex()
	l.Upstream = res.GetUpstream()
	l.IpPool = res.GetIpPool()
	l.NetmaskLen = common.FromOpt(res.GetNetmaskLen())
	l.DefaultGateway = common.FromOpt(res.GetDefaultGateway())
	l.Vip = common.FromOpt(res.GetVip())
	l.VirtualRouterID = common.FromOpt(res.GetVirtualRouterID())
	l.PacketFilterID = common.FromOpt(res.GetPacketFilterID())
}

type NodeInterface struct {
	InterfaceIndex int16
	Addresses      []NodeInterfaceAddress
}

func (n *NodeInterface) From(res *v1.ReadLoadBalancerNodeInterface) {
	n.InterfaceIndex = res.GetInterfaceIndex()
	n.Addresses = common.MapSlice(res.GetAddresses(), common.ConvertFrom[v1.ReadLoadBalancerNodeInterfaceAddress, NodeInterfaceAddress]())
}

type NodeInterfaceAddress struct {
	Address string
	Vip     bool
}

func (n *NodeInterfaceAddress) From(res *v1.ReadLoadBalancerNodeInterfaceAddress) {
	n.Address = res.GetAddress()
	n.Vip = res.GetVip()
}

type CreateParams struct {
	Name             string
	ServiceClassPath string
	NameServers      []v1.IPv4
	Interfaces       []LoadBalancerInterface
}

func (c CreateParams) into() (ret v1.CreateLoadBalancer) {
	ret.SetName(c.Name)
	ret.SetServiceClassPath(c.ServiceClassPath)
	ret.SetNameServers(c.NameServers)
	ret.SetInterfaces(common.MapSlice(c.Interfaces, LoadBalancerInterface.into))

	return
}

type LoadBalancerDetail struct {
	LoadBalancerID   v1.LoadBalancerID
	Name             string
	ServiceClassPath string
	NameServers      []v1.IPv4
	Interfaces       []LoadBalancerInterface
	Created          int
	Deleting         bool
}

func (l *LoadBalancerDetail) from(res *v1.ReadLoadBalancerDetail) {
	l.LoadBalancerID = res.GetLoadBalancerID()
	l.Name = res.GetName()
	l.ServiceClassPath = res.GetServiceClassPath()
	l.NameServers = res.GetNameServers()
	l.Interfaces = common.MapSlice(res.GetInterfaces(), common.ConvertFrom[v1.LoadBalancerInterface, LoadBalancerInterface]())
	l.Created = res.GetCreated()
	l.Deleting = res.GetDeleting()
}

type LoadBalancerNodeDetail struct {
	LoadBalancerNodeID v1.LoadBalancerNodeID
	ResourceID         *string
	Interfaces         []NodeInterface
	Status             v1.LoadBalancerNodeStatus
	ArchiveVersion     *string
	CreateErrorMessage *string
	Created            int
}

func (l *LoadBalancerNodeDetail) from(res *v1.ReadLoadBalancerNode) {
	l.LoadBalancerNodeID = res.GetLoadBalancerNodeID()
	l.ResourceID = common.FromOpt(res.GetResourceID())
	l.Interfaces = common.MapSlice(res.GetInterfaces(), common.ConvertFrom[v1.ReadLoadBalancerNodeInterface, NodeInterface]())
	l.Status = res.GetStatus()
	l.ArchiveVersion = common.FromOpt(res.GetArchiveVersion())
	l.CreateErrorMessage = common.FromOpt(res.GetCreateErrorMessage())
	l.Created = res.GetCreated()
}
