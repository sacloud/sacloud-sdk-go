// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package application

import (
	"context"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
)

type ApplicationAPI interface {
	// List returns the list of Applications, paginated.
	// Pass nil to `cursor` to get the first page, or
	// previously returned `nextCursor` to get the next page.
	List(ctx context.Context, maxItems int64, cursor *string) (list []v1.ReadApplicationDetail, nextCursor *string, err error)
	Create(ctx context.Context, name string, clusterID v1.ClusterID) (app *v1.CreatedApplication, err error)
	Read(ctx context.Context, id v1.ApplicationID) (app *ApplicationDetail, err error)
	Delete(ctx context.Context, id v1.ApplicationID) error
	Containers(ctx context.Context, id v1.ApplicationID) (nodes []Placement, err error)
}

type ApplicationOp struct{ *v1.Client }

func NewApplicationOp(client *v1.Client) *ApplicationOp { return &ApplicationOp{Client: client} }

func (op *ApplicationOp) List(ctx context.Context, maxItems int64, cursor *string) (apps []v1.ReadApplicationDetail, nextCursor *string, err error) {
	res, err := common.ErrorFromDecodedResponse("Application.List", func() (*v1.ListApplicationResponse, error) {
		return op.Client.ListApplications(ctx, v1.ListApplicationsParams{
			Cursor:   common.IntoOpt[v1.OptString](cursor),
			MaxItems: maxItems,
		})
	})

	if res != nil {
		apps = res.Applications
		nextCursor = common.FromOpt(res.NextCursor)
	}

	return
}

func (op *ApplicationOp) Create(ctx context.Context, name string, clusterID v1.ClusterID) (app *v1.CreatedApplication, err error) {
	res, err := common.ErrorFromDecodedResponse("Application.Create", func() (*v1.CreateApplicationResponse, error) {
		return op.Client.CreateApplication(ctx, &v1.CreateApplication{
			Name:      name,
			ClusterID: clusterID,
		})
	})

	if res != nil {
		app = &res.Application
	}

	return
}

type ApplicationDetail struct {
	ApplicationID          v1.ApplicationID
	Name                   string
	ClusterID              v1.ClusterID
	ClusterName            string
	ActiveVersion          *int32
	DesiredCount           *int32
	ScalingCooldownSeconds int32
}

func (op *ApplicationOp) Read(ctx context.Context, id v1.ApplicationID) (app *ApplicationDetail, err error) {
	res, err := common.ErrorFromDecodedResponse("Application.Read", func() (*v1.GetApplicationResponse, error) {
		return op.Client.GetApplication(ctx, v1.GetApplicationParams{
			ApplicationID: id,
		})
	})

	if res != nil {
		tmp := res.GetApplication()
		app = &ApplicationDetail{
			ApplicationID:          tmp.GetApplicationID(),
			Name:                   tmp.GetName(),
			ClusterID:              tmp.GetClusterID(),
			ClusterName:            tmp.GetClusterName(),
			ActiveVersion:          common.FromOpt(tmp.GetActiveVersion()),
			DesiredCount:           common.FromOpt(tmp.GetDesiredCount()),
			ScalingCooldownSeconds: tmp.GetScalingCooldownSeconds(),
		}
	}

	return
}

func (op *ApplicationOp) Delete(ctx context.Context, id v1.ApplicationID) error {
	return common.ErrorFromDecodedResponseE("Application.Delete", func() error {
		return op.Client.DeleteApplication(ctx, v1.DeleteApplicationParams{
			ApplicationID: id,
		})
	})
}

type Placement struct {
	NodeID          string
	ContainersStats v1.ApplicationContainersStats
	Desired         v1.ApplicationPeekDesiredContainersResponse
}

func (op *ApplicationOp) Containers(ctx context.Context, id v1.ApplicationID) (nodes []Placement, err error) {
	res, err := common.ErrorFromDecodedResponse("Application.Containers", func() (*v1.GetApplicationContainersResponse, error) {
		return op.Client.GetApplicationContainers(ctx, v1.GetApplicationContainersParams{
			ApplicationID: id,
		})
	})

	if res != nil {
		nodes = make([]Placement, len(res.Nodes))
		for i, n := range res.Nodes {
			var m Placement
			m.NodeID = n.GetNodeID()
			if s, ok := n.GetContainersStats().Get(); ok {
				m.ContainersStats = s
			}
			if d, ok := n.GetDesired().Get(); ok {
				m.Desired = d
			}
			nodes[i] = m
		}
	}

	return
}

var _ ApplicationAPI = (*ApplicationOp)(nil)
