// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package workernode

import (
	"context"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
)

type WorkerNodeAPI interface {
	// List returns the list of WorkerNodes, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, maxItems int64, cursor *v1.WorkerNodeID) (list []WorkerNodeDetail, nextCursor *v1.WorkerNodeID, err error)
	Read(ctx context.Context, id v1.WorkerNodeID) (node *WorkerNodeDetail, err error)
	Update(ctx context.Context, id v1.WorkerNodeID, draining bool) error
}

type WorkerNodeOp struct {
	client             *v1.Client
	clusterID          v1.ClusterID
	autoScalingGroupID v1.AutoScalingGroupID
}

func NewWorkerNodeOp(client *v1.Client, clusterID v1.ClusterID, autoScalingGroupID v1.AutoScalingGroupID) *WorkerNodeOp {
	return &WorkerNodeOp{
		client:             client,
		clusterID:          clusterID,
		autoScalingGroupID: autoScalingGroupID,
	}
}

func (op *WorkerNodeOp) List(ctx context.Context, maxItems int64, cursor *v1.WorkerNodeID) (nodes []WorkerNodeDetail, nextCursor *v1.WorkerNodeID, err error) {
	res, err := common.ErrorFromDecodedResponse("WorkerNode.List", func() (*v1.ListWorkerNodesResponse, error) {
		return op.client.ListWorkerNodes(ctx, v1.ListWorkerNodesParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
			Cursor:             common.IntoOpt[v1.OptWorkerNodeID](cursor),
			MaxItems:           maxItems,
		})
	})

	if res != nil {
		nextCursor = common.FromOpt(res.NextCursor)
		nodes = common.MapSlice(res.WorkerNodes, func(w v1.ReadWorkerNodeSummary) WorkerNodeDetail {
			var detail WorkerNodeDetail
			detail.fromSummary(&w)
			return detail
		})
	}

	return
}

func (op *WorkerNodeOp) Read(ctx context.Context, id v1.WorkerNodeID) (node *WorkerNodeDetail, err error) {
	res, err := common.ErrorFromDecodedResponse("WorkerNode.Read", func() (*v1.GetWorkerNodeResponse, error) {
		return op.client.GetWorkerNode(ctx, v1.GetWorkerNodeParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
			WorkerNodeID:       id,
		})
	})

	if res != nil {
		var detail WorkerNodeDetail
		detail.from(&res.WorkerNode)
		node = &detail
	}

	return
}

func (op *WorkerNodeOp) Update(ctx context.Context, id v1.WorkerNodeID, draining bool) error {
	return common.ErrorFromDecodedResponseE("WorkerNode.Update", func() error {
		var req v1.UpdateWorkerNodeDrainingRequest
		req.SetDraining(draining)
		return op.client.UpdateWorkerNodeDrainingState(ctx, &req, v1.UpdateWorkerNodeDrainingStateParams{
			ClusterID:          op.clusterID,
			AutoScalingGroupID: op.autoScalingGroupID,
			WorkerNodeID:       id,
		})
	})
}

var _ WorkerNodeAPI = (*WorkerNodeOp)(nil)

type WorkerNodeDetail struct {
	WorkerNodeID       v1.WorkerNodeID
	ResourceID         *string
	Draining           bool
	Status             v1.WorkerNodeStatus
	Healthy            bool
	Creating           bool
	Created            int
	RunningContainers  []v1.RunningContainer
	NetworkInterfaces  []WorkerNodeNetworkInterface
	ArchiveVersion     *string
	CreateErrorMessage *string
}

func (w *WorkerNodeDetail) fromSummary(res *v1.ReadWorkerNodeSummary) {
	w.WorkerNodeID = res.GetWorkerNodeID()
	w.ResourceID = common.FromOpt(res.GetResourceID())
	w.Draining = res.GetDraining()
	w.Status = res.GetStatus()
	w.Healthy = false
	w.Creating = false
	w.Created = res.GetCreated()
	w.RunningContainers = []v1.RunningContainer(nil)
	w.NetworkInterfaces = common.MapSlice(res.GetNetworkInterfaces(), common.ConvertFrom[v1.ReadWorkerNodeNetworkInterface, WorkerNodeNetworkInterface]())
	w.ArchiveVersion = common.FromOpt(res.GetArchiveVersion())
	w.CreateErrorMessage = common.FromOpt(res.GetCreateErrorMessage())
}

func (w *WorkerNodeDetail) from(res *v1.ReadWorkerNodeDetail) {
	w.WorkerNodeID = res.GetWorkerNodeID()
	w.ResourceID = common.FromOpt(res.GetResourceID())
	w.Draining = res.GetDraining()
	w.Status = res.GetStatus()
	w.Healthy = res.GetHealthy()
	w.Creating = res.GetCreating()
	w.Created = res.GetCreated()
	w.RunningContainers = res.GetRunningContainers()
	w.NetworkInterfaces = common.MapSlice(res.GetNetworkInterfaces(), common.ConvertFrom[v1.ReadWorkerNodeNetworkInterface, WorkerNodeNetworkInterface]())
	w.ArchiveVersion = common.FromOpt(res.GetArchiveVersion())
	w.CreateErrorMessage = common.FromOpt(res.GetCreateErrorMessage())
}

type WorkerNodeNetworkInterface struct {
	InterfaceIndex int16
	Addresses      []string
}

func (w *WorkerNodeNetworkInterface) From(res *v1.ReadWorkerNodeNetworkInterface) {
	w.InterfaceIndex = res.GetInterfaceIndex()
	w.Addresses = common.MapSlice(res.GetAddresses(), func(a v1.ReadWorkerNodeInterfaceAddress) string {
		return a.GetAddress()
	})
}
