// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"errors"

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
	res, err := op.fedClient.client.ListClusters(ctx)
	if err != nil {
		return nil, NewAPIError("Site.List", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerListClustersResponse:
		return r.Data, nil
	case *v2.ErrorResponse:
		return nil, NewAPIError("Site.List", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return nil, unknownError("Site.List", r)
	}
}

func (op *siteOp) Read(ctx context.Context, id string) (*v2.ClusterResponseBody, error) {
	res, err := op.fedClient.client.ReadCluster(ctx, v2.ReadClusterParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Site.Read", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerReadClusterResponse:
		return &r.Data, nil
	case *v2.ReadClusterUnauthorized:
		return nil, NewAPIError("Site.Read", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ReadClusterNotFound:
		return nil, NewAPIError("Site.Read", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return nil, unknownError("Site.Read", r)
	}
}

func (op *siteOp) ListPlans(ctx context.Context) ([]v2.PlanItem, error) {
	res, err := op.siteClient.client.ListPlans(ctx)
	if err != nil {
		return nil, NewAPIError("Site.ListPlans", 0, err)
	}

	switch r := res.(type) {
	case *v2.ListPlansOK:
		return r.Data, nil
	case *v2.ErrorResponse:
		return nil, NewAPIError("Site.ListPlans", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ErrorResponseStatusCode:
		return nil, NewAPIError("Site.ListPlans", r.StatusCode, errors.New(string(r.Response.Error.Message)))
	default:
		return nil, unknownError("Site.ListPlans", r)
	}
}
