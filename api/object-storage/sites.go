// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"errors"

	v2 "github.com/sacloud/object-storage-api-go/apis/v2"
)

type SiteAPI interface {
	List(ctx context.Context) ([]v2.ModelCluster, error)
	Read(ctx context.Context, siteId string) (*v2.ModelCluster, error)

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

func (op *siteOp) List(ctx context.Context) ([]v2.ModelCluster, error) {
	res, err := op.fedClient.client.GetClusters(ctx)
	if err != nil {
		return nil, NewAPIError("Site.List", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerListClustersRes:
		return r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Site.List", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Site.List", 0, errors.New("unknown error"))
	}
}

func (op *siteOp) Read(ctx context.Context, id string) (*v2.ModelCluster, error) {
	res, err := op.fedClient.client.GetCluster(ctx, v2.GetClusterParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Site.Read", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerGetClusterRes:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("Site.Read", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Site.Read", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Site.Read", 0, errors.New("unknown error"))
	}
}

func (op *siteOp) ListPlans(ctx context.Context) ([]v2.PlanItem, error) {
	res, err := op.siteClient.client.GetPlans(ctx)
	if err != nil {
		return nil, NewAPIError("Site.ListPlans", 0, err)
	}

	switch r := res.(type) {
	case *v2.GetPlansOK:
		return r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Site.ListPlans", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("Site.ListPlans", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Site.ListPlans", 0, errors.New("unknown error"))
	}
}
