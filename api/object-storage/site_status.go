// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
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
	res, err := errorFromDecodedResponse[v2.Status]("SiteStatus.Read", func() (any, error) { return op.client.client.ReadStatus(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *siteStatusOp) ReadQuota(ctx context.Context) (*v2.QuotaData, error) {
	res, err := errorFromDecodedResponse[v2.Quota]("SiteStatus.ReadQuota", func() (any, error) { return op.client.client.ReadQuota(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *siteStatusOp) ReadBucketMetering(ctx context.Context, bucketName string, from, to time.Time) ([]v2.BucketBillingItem, error) {
	res, err := errorFromDecodedResponse[v2.ListBucketMeteringOK]("SiteStatus.ReadBucketMetering", func() (any, error) {
		return op.client.client.ListBucketMetering(ctx, v2.ListBucketMeteringParams{
			Name: v2.BucketName(bucketName),
			From: v2.YearMonth(from.Format("2006-01")),
			To:   v2.YearMonth(to.Format("2006-01")),
		})
	})
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}
