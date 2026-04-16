// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"errors"
	"time"

	v2 "github.com/sacloud/object-storage-api-go/apis/v2"
)

type SiteStatusAPI interface {
	Read(ctx context.Context) (*v2.StatusData, error)
	ReadQuota(ctx context.Context) (*v2.QuotaData, error)
	ReadBucketMetering(ctx context.Context, bucketName string, from, to time.Time) ([]v2.BucketBillingItem, error)
}

var _ SiteStatusAPI = (*siteStatusOp)(nil)

type siteStatusOp struct {
	client *SiteClient
}

func NewSiteStatusOp(client *SiteClient) SiteStatusAPI {
	return &siteStatusOp{client: client}
}

func (op *siteStatusOp) Read(ctx context.Context) (*v2.StatusData, error) {
	res, err := op.client.client.GetStatus(ctx)
	if err != nil {
		return nil, NewAPIError("SiteStatus.Read", 0, err)
	}

	switch r := res.(type) {
	case *v2.Status:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("SiteStatus.Read", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("SiteStatus.Read", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("SiteStatus.Read", 0, errors.New("unknown error"))
	}
}

func (op *siteStatusOp) ReadQuota(ctx context.Context) (*v2.QuotaData, error) {
	res, err := op.client.client.GetQuota(ctx)
	if err != nil {
		return nil, NewAPIError("SiteStatus.ReadQuota", 0, err)
	}

	switch r := res.(type) {
	case *v2.Quota:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("SiteStatus.ReadQuota", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("SiteStatus.ReadQuota", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("SiteStatus.ReadQuota", 0, errors.New("unknown error"))
	}
}

func (op *siteStatusOp) ReadBucketMetering(ctx context.Context, bucketName string, from, to time.Time) ([]v2.BucketBillingItem, error) {
	res, err := op.client.client.GetBucketMetering(ctx, v2.GetBucketMeteringParams{Name: v2.BucketName(bucketName), From: from, To: to})
	if err != nil {
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", 0, err)
	}

	switch r := res.(type) {
	case *v2.GetBucketMeteringOK:
		return r.Data, nil
	case *v2.Error400:
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error401:
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.ErrorDefaultStatusCode:
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", r.StatusCode, errors.New(string(r.Response.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", 0, errors.New("unknown error"))
	}
}
