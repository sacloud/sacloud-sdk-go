// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"errors"
	"time"

	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
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
	res, err := op.client.client.ReadStatus(ctx)
	if err != nil {
		return nil, NewAPIError("SiteStatus.Read", 0, err)
	}

	switch r := res.(type) {
	case *v2.Status:
		return &r.Data, nil
	case *v2.ErrorResponse:
		return nil, NewAPIError("SiteStatus.Read", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ErrorResponseStatusCode:
		return nil, NewAPIError("SiteStatus.Read", r.StatusCode, errors.New(string(r.Response.Error.Message)))
	default:
		return nil, unknownError("SiteStatus.Read", r)
	}
}

func (op *siteStatusOp) ReadQuota(ctx context.Context) (*v2.QuotaData, error) {
	res, err := op.client.client.ReadQuota(ctx)
	if err != nil {
		return nil, NewAPIError("SiteStatus.ReadQuota", 0, err)
	}

	switch r := res.(type) {
	case *v2.Quota:
		return &r.Data, nil
	case *v2.ErrorResponse:
		return nil, NewAPIError("SiteStatus.ReadQuota", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ErrorResponseStatusCode:
		return nil, NewAPIError("SiteStatus.ReadQuota", r.StatusCode, errors.New(string(r.Response.Error.Message)))
	default:
		return nil, unknownError("SiteStatus.ReadQuota", r)
	}
}

func (op *siteStatusOp) ReadBucketMetering(ctx context.Context, bucketName string, from, to time.Time) ([]v2.BucketBillingItem, error) {
	res, err := op.client.client.ListBucketMetering(ctx, v2.ListBucketMeteringParams{
		Name: v2.BucketName(bucketName),
		From: v2.YearMonth(from.Format("2006-01")),
		To:   v2.YearMonth(to.Format("2006-01")),
	})
	if err != nil {
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", 0, err)
	}

	switch r := res.(type) {
	case *v2.ListBucketMeteringOK:
		return r.Data, nil
	case *v2.ListBucketMeteringBadRequest:
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ListBucketMeteringUnauthorized:
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ErrorResponseStatusCode:
		return nil, NewAPIError("SiteStatus.ReadBucketMetering", r.StatusCode, errors.New(string(r.Response.Error.Message)))
	default:
		return nil, unknownError("SiteStatus.ReadBucketMetering", r)
	}
}
