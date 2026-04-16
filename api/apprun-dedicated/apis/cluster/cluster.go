// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"context"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
	"github.com/sacloud/saclient-go"
)

type ClusterAPI interface {
	// List returns the list of Clusters, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, maxItems int64, cursor *v1.ClusterID) (list []ClusterDetail, nextCursor *v1.ClusterID, err error)
	Create(ctx context.Context, params CreateParams) (cluster *v1.CreatedCluster, err error)
	Read(ctx context.Context, id v1.ClusterID) (cluster *ClusterDetail, err error)
	Update(ctx context.Context, id v1.ClusterID, params UpdateParams) error
	Delete(ctx context.Context, id v1.ClusterID) error
}

type ClusterOp struct{ *v1.Client }

func NewClusterOp(client *v1.Client) *ClusterOp { return &ClusterOp{Client: client} }

func (op *ClusterOp) List(ctx context.Context, maxItems int64, cursor *v1.ClusterID) (clusters []ClusterDetail, nextCursor *v1.ClusterID, err error) {
	res, err := common.ErrorFromDecodedResponse("Cluster.List", func() (*v1.ListClusterResponse, error) {
		return op.Client.ListClusters(ctx, v1.ListClustersParams{
			Cursor:   common.IntoOpt[v1.OptClusterID](cursor),
			MaxItems: maxItems,
		})
	})

	if res != nil {
		nextCursor = common.FromOpt(res.NextCursor)
		clusters = common.MapSlice(res.Clusters, func(c v1.ReadClusterSummary) ClusterDetail {
			var detail ClusterDetail
			detail.fromSummary(&c)
			return detail
		})
	}

	return
}

func (op *ClusterOp) Create(ctx context.Context, params CreateParams) (cluster *v1.CreatedCluster, err error) {
	res, err := common.ErrorFromDecodedResponse("Cluster.Create", func() (*v1.CreateClusterResponse, error) {
		return op.Client.CreateCluster(ctx, saclient.Ptr(params.into()))
	})

	if res != nil {
		cluster = &res.Cluster
	}

	return
}

func (op *ClusterOp) Read(ctx context.Context, id v1.ClusterID) (cluster *ClusterDetail, err error) {
	res, err := common.ErrorFromDecodedResponse("Cluster.Read", func() (*v1.GetClusterResponse, error) {
		return op.Client.GetCluster(ctx, v1.GetClusterParams{ClusterID: id})
	})

	if res != nil {
		var detail ClusterDetail
		detail.from(&res.Cluster)
		cluster = &detail
	}

	return
}

func (op *ClusterOp) Update(ctx context.Context, id v1.ClusterID, params UpdateParams) error {
	return common.ErrorFromDecodedResponseE("Cluster.Update", func() error {
		return op.Client.UpdateCluster(ctx, saclient.Ptr(params.into()), v1.UpdateClusterParams{ClusterID: id})
	})
}

func (op *ClusterOp) Delete(ctx context.Context, id v1.ClusterID) error {
	return common.ErrorFromDecodedResponseE("Cluster.Delete", func() error {
		return op.Client.DeleteCluster(ctx, v1.DeleteClusterParams{ClusterID: id})
	})
}

var _ ClusterAPI = (*ClusterOp)(nil)

type CreateParams struct {
	Name               string
	LetsEncryptEmail   *string
	Ports              []v1.CreateLoadBalancerPort
	ServicePrincipalID string
}

func (c CreateParams) into() (ret v1.CreateCluster) {
	ret.SetName(c.Name)
	ret.SetLetsEncryptEmail(common.IntoOpt[v1.OptString](c.LetsEncryptEmail))
	ret.SetPorts(c.Ports)
	ret.SetServicePrincipalID(c.ServicePrincipalID)

	return
}

type UpdateParams struct {
	LetsEncryptEmail   *string
	ServicePrincipalID string
}

func (u UpdateParams) into() (ret v1.UpdateCluster) {
	ret.SetLetsEncryptEmail(common.IntoOpt[v1.OptString](u.LetsEncryptEmail))
	ret.SetServicePrincipalID(u.ServicePrincipalID)

	return
}

type ClusterDetail struct {
	ClusterID           v1.ClusterID
	Name                string
	Ports               []v1.ReadLoadBalancerPort
	ServicePrincipalID  string
	HasLetsEncryptEmail bool
	Created             int
}

func (c *ClusterDetail) fromSummary(res *v1.ReadClusterSummary) {
	c.ClusterID = res.GetClusterID()
	c.Name = res.GetName()
	c.Ports = []v1.ReadLoadBalancerPort(nil)
	c.ServicePrincipalID = ""
	c.HasLetsEncryptEmail = false
	c.Created = res.GetCreated()
}

func (c *ClusterDetail) from(res *v1.ReadClusterDetail) {
	c.ClusterID = res.GetClusterID()
	c.Name = res.GetName()
	c.Ports = res.GetPorts()
	c.ServicePrincipalID = res.GetServicePrincipalID()
	c.HasLetsEncryptEmail = res.GetHasLetsEncryptEmail()
	c.Created = res.GetCreated()
}
