// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"

	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
)

type SiteAPI interface {
	List(ctx context.Context) ([]v2.ClusterResponseBody, error)
	Read(ctx context.Context, siteId string) (*v2.ClusterResponseBody, error)

	ListPlans(ctx context.Context) ([]v2.PlanItem, error)
}

var _ SiteAPI = (*siteOp)(nil)

type siteOp struct {
	fedClient  *FedClient
	siteClient *SiteClient // for ListPlans
}

func NewSiteOp(fedClient *FedClient) SiteAPI {
	return NewSiteWithPlansOp(fedClient, nil)
}

func NewSiteWithPlansOp(fedClient *FedClient, siteClient *SiteClient) SiteAPI {
	return &siteOp{fedClient: fedClient, siteClient: siteClient}
}

func (op *siteOp) List(ctx context.Context) ([]v2.ClusterResponseBody, error) {
	res, err := errorFromDecodedResponse[v2.HandlerListClustersResponse]("Site.List", func() (any, error) { return op.fedClient.client.ListClusters(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

func (op *siteOp) Read(ctx context.Context, id string) (*v2.ClusterResponseBody, error) {
	res, err := errorFromDecodedResponse[v2.HandlerReadClusterResponse]("Site.Read", func() (any, error) { return op.fedClient.client.ReadCluster(ctx, v2.ReadClusterParams{ID: id}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *siteOp) ListPlans(ctx context.Context) ([]v2.PlanItem, error) {
	res, err := errorFromDecodedResponse[v2.ListPlansOK]("Site.ListPlans", func() (any, error) { return op.siteClient.client.ListPlans(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}
