// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package cluster

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type ClusterAPI interface {
	// List returns the list of Clusters, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, maxItems int64, cursor *v1.ClusterID) (list []ClusterDetail, nextCursor *v1.ClusterID, err error)
	Create(ctx context.Context, params CreateParams) (cluster *v1.ClusterIdentifier, err error)
	Read(ctx context.Context, id v1.ClusterID) (cluster *ClusterDetail, err error)
	Update(ctx context.Context, id v1.ClusterID, params UpdateParams) error
	Delete(ctx context.Context, id v1.ClusterID) error
	ReadLoadBalancer(ctx context.Context, id v1.ClusterID) (settings *v1.ReadClusterLoadBalancer, err error)
	UpdateLoadBalancer(ctx context.Context, id v1.ClusterID, params UpdateLoadBalancerParams) error
}

type ClusterOp struct{ *v1.Client }

func NewClusterOp(client *v1.Client) *ClusterOp { return &ClusterOp{Client: client} }

func (op *ClusterOp) List(ctx context.Context, maxItems int64, cursor *v1.ClusterID) (clusters []ClusterDetail, nextCursor *v1.ClusterID, err error) {
	res, err := common.ErrorFromDecodedResponse("Cluster.List", func() (*v1.ListClustersResponse, error) {
		return op.Client.ListClusters(ctx, v1.ListClustersParams{
			Cursor:   into.Opt[v1.OptClusterID](cursor),
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

func (op *ClusterOp) Create(ctx context.Context, params CreateParams) (cluster *v1.ClusterIdentifier, err error) {
	res, err := common.ErrorFromDecodedResponse("Cluster.Create", func() (*v1.CreateClusterResponse, error) {
		return op.Client.CreateCluster(ctx, new(params.into()))
	})

	if res != nil {
		cluster = &res.Cluster
	}

	return
}

func (op *ClusterOp) Read(ctx context.Context, id v1.ClusterID) (cluster *ClusterDetail, err error) {
	res, err := common.ErrorFromDecodedResponse("Cluster.Read", func() (*v1.ReadClusterResponse, error) {
		return op.Client.ReadCluster(ctx, v1.ReadClusterParams{ClusterID: id})
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
		return op.Client.UpdateCluster(ctx, new(params.into()), v1.UpdateClusterParams{ClusterID: id})
	})
}

func (op *ClusterOp) Delete(ctx context.Context, id v1.ClusterID) error {
	return common.ErrorFromDecodedResponseE("Cluster.Delete", func() error {
		return op.Client.DeleteCluster(ctx, v1.DeleteClusterParams{ClusterID: id})
	})
}

func (op *ClusterOp) ReadLoadBalancer(ctx context.Context, id v1.ClusterID) (settings *v1.ReadClusterLoadBalancer, err error) {
	res, err := common.ErrorFromDecodedResponse("Cluster.ReadLoadBalancer", func() (*v1.ReadClusterLoadBalancerResponse, error) {
		return op.Client.ReadClusterLoadBalancer(ctx, v1.ReadClusterLoadBalancerParams{ClusterID: id})
	})

	if res != nil {
		settings = &res.LoadBalancer
	}

	return
}

func (op *ClusterOp) UpdateLoadBalancer(ctx context.Context, id v1.ClusterID, params UpdateLoadBalancerParams) error {
	return common.ErrorFromDecodedResponseE("Cluster.UpdateLoadBalancer", func() error {
		return op.Client.PatchClusterLoadBalancer(ctx, params.into(), v1.PatchClusterLoadBalancerParams{ClusterID: id})
	})
}

var _ ClusterAPI = (*ClusterOp)(nil)

type CreateParams struct {
	Name               string                       `json:"name"`
	LetsEncryptEmail   *string                      `json:"letsEncryptEmail"`
	Ports              []v1.ClusterLoadBalancerPort `json:"ports"`
	ServicePrincipalID string                       `json:"servicePrincipalID"`
}

func (c CreateParams) into() (ret v1.CreateClusterRequest) {
	ret.SetName(c.Name)
	ret.SetLetsEncryptEmail(into.Opt[v1.OptString](c.LetsEncryptEmail))
	ret.SetPorts(c.Ports)
	ret.SetServicePrincipalID(v1.ServicePrincipalID(c.ServicePrincipalID))

	return
}

type UpdateParams struct {
	LetsEncryptEmail   *string `json:"letsEncryptEmail"`
	ServicePrincipalID string  `json:"servicePrincipalID"`
}

func (u UpdateParams) into() (ret v1.UpdateClusterRequest) {
	ret.SetLetsEncryptEmail(into.Opt[v1.OptString](u.LetsEncryptEmail))
	ret.SetServicePrincipalID(v1.ServicePrincipalID(u.ServicePrincipalID))

	return
}

type ClusterDetail struct {
	ClusterID           v1.ClusterID                 `json:"clusterID"`
	Name                string                       `json:"name"`
	Ports               []v1.ClusterLoadBalancerPort `json:"ports"`
	ServicePrincipalID  string                       `json:"servicePrincipalID"`
	HasLetsEncryptEmail bool                         `json:"hasLetsEncryptEmail"`
	Created             int                          `json:"created"`
}

func (c *ClusterDetail) fromSummary(res *v1.ReadClusterSummary) {
	c.ClusterID = res.GetClusterID()
	c.Name = res.GetName()
	c.Ports = []v1.ClusterLoadBalancerPort(nil)
	c.ServicePrincipalID = ""
	c.HasLetsEncryptEmail = false
	c.Created = res.GetCreated()
}

func (c *ClusterDetail) from(res *v1.ReadCluster) {
	c.ClusterID = res.GetClusterID()
	c.Name = res.GetName()
	c.Ports = res.GetPorts()
	c.ServicePrincipalID = string(res.GetServicePrincipalID())
	c.HasLetsEncryptEmail = res.GetHasLetsEncryptEmail()
	c.Created = res.GetCreated()
}

type UpdateParam[T any] struct {
	// DiscardOldValue indicates whether the old value should be discarded.
	// This is intentionally named as such so that a zero value means "retains old value"
	DiscardOldValue bool `json:"discardOldValue"`
	UpdatingValue   T    `json:"value,omitempty"`
}

type UpdateLoadBalancerPorts struct {
	Port     uint16 `json:"port"`
	Protocol string `json:"protocol"`
}

type UpdateLoadBalancerParams struct {
	LetsEncryptEmail UpdateParam[*string]                   `json:"letsEncryptEmail,omitempty"`
	Ports            UpdateParam[[]UpdateLoadBalancerPorts] `json:"ports,omitempty"`
}

func (u UpdateLoadBalancerParams) into() (ret *v1.PatchClusterLoadBalancerRequest) {
	var le v1.OptNilString
	if !u.LetsEncryptEmail.DiscardOldValue {
		le.Reset()
	} else if v := u.LetsEncryptEmail.UpdatingValue; v == nil {
		le.SetToNull()
	} else {
		le.SetTo(*u.LetsEncryptEmail.UpdatingValue)
	}

	var p []v1.ClusterLoadBalancerPortMergePatchCreateOrUpdateReplaceOnly
	if u.Ports.DiscardOldValue {
		p = make([]v1.ClusterLoadBalancerPortMergePatchCreateOrUpdateReplaceOnly, 0, len(u.Ports.UpdatingValue))

		for _, port := range u.Ports.UpdatingValue {
			q := v1.ClusterLoadBalancerPortMergePatchCreateOrUpdateReplaceOnly{
				Port:     port.Port,
				Protocol: (v1.ClusterLoadBalancerPortMergePatchCreateOrUpdateReplaceOnlyProtocol)(port.Protocol),
			}
			p = append(p, q)
		}
	}

	ret = &v1.PatchClusterLoadBalancerRequest{}
	ret.SetPorts(p)
	ret.SetLetsEncryptEmail(le)
	return
}
