// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"fmt"

	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
)

type BucketAPI interface {
	List(ctx context.Context) ([]v2.BucketListResponseBody, error)
	Create(ctx context.Context, params *BucketCreateParams) (*v2.BucketResponseBody, error)
	Delete(ctx context.Context, bucketName string) error
}

var _ BucketAPI = (*bucketOp)(nil)

type bucketOp struct {
	fedClient  *FedClient
	siteClient *SiteClient // for List
}

func NewBucketOp(fedClient *FedClient, siteClient *SiteClient) BucketAPI {
	return &bucketOp{fedClient: fedClient, siteClient: siteClient}
}

func (op *bucketOp) List(ctx context.Context) ([]v2.BucketListResponseBody, error) {
	res, err := errorFromDecodedResponse[v2.HandlerListBucketsResponse]("Buckets.List", func() (any, error) { return op.siteClient.client.ListBuckets(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetData(), nil
}

type BucketCreateParams struct {
	Bucket string
	SiteId string
	Plan   string
}

func createRequest(params *BucketCreateParams) *v2.HandlerCreateBucketRequestBody {
	res := &v2.HandlerCreateBucketRequestBody{
		ClusterID: params.SiteId,
	}

	switch params.SiteId {
	case "isk01", "tky01":
		res.Plan = v2.NewOptNilHandlerCreateBucketRequestBodyPlan(v2.HandlerCreateBucketRequestBodyPlan{
			Type:             v2.NewOptModelPlanType(v2.ModelPlanType("standard")),
			ServiceClassPath: v2.NewOptServiceClassPath(v2.ServiceClassPath(fmt.Sprintf("objectstorage/%s/bucket", params.SiteId))),
		})
	case "arc02":
		res.Plan = v2.NewOptNilHandlerCreateBucketRequestBodyPlan(v2.HandlerCreateBucketRequestBodyPlan{
			Type:             v2.NewOptModelPlanType(v2.ModelPlanType("archive")),
			ServiceClassPath: v2.NewOptServiceClassPath(v2.ServiceClassPath(fmt.Sprintf("objectstorage/%s/bucket/%s", params.SiteId, params.Plan))),
		})
	}

	return res
}

func (op *bucketOp) Create(ctx context.Context, params *BucketCreateParams) (*v2.BucketResponseBody, error) {
	res, err := errorFromDecodedResponse[v2.HandlerCreateBucketResponse]("Buckets.Create", func() (any, error) {
		return op.fedClient.client.CreateBucket(ctx, createRequest(params), v2.CreateBucketParams{Name: params.Bucket})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *bucketOp) Delete(ctx context.Context, bucketName string) error {
	return errorFromDecodedResponseE[v2.DeleteBucketNoContent]("Buckets.Delete", func() (any, error) {
		return op.fedClient.client.DeleteBucket(ctx, v2.DeleteBucketParams{Name: bucketName})
	})
}

type BucketExtraAPI interface {
	ReadEncryption(ctx context.Context) (*v2.HandlerReadBucketEncryptionResponse, error)
	EnableEncryption(ctx context.Context, KMSKeyID string) error
	DisableEncryption(ctx context.Context) error

	ReadReplication(ctx context.Context) (*v2.ModelReplication, error)
	EnableReplication(ctx context.Context, targetBucket string) (*v2.ModelReplication, error)
	DisableReplication(ctx context.Context) error

	ReadPenalty(ctx context.Context) (*v2.BucketPenaltyData, error)
	ReadUsage(ctx context.Context) (*v2.BucketUsageData, error)
	ReadQuota(ctx context.Context) (*v2.BucketQuotaData, error)
}

var _ BucketExtraAPI = (*bucketExtraOp)(nil)

type bucketExtraOp struct {
	siteClient *SiteClient
	fedClient  *FedClient // for replication settings
	bucket     string
}

func NewBucketExtraOp(siteClient *SiteClient, fedClient *FedClient, bucket string) BucketExtraAPI {
	return &bucketExtraOp{siteClient: siteClient, fedClient: fedClient, bucket: bucket}
}

func (op *bucketExtraOp) ReadEncryption(ctx context.Context) (*v2.HandlerReadBucketEncryptionResponse, error) {
	res, err := errorFromDecodedResponse[v2.ReadBucketEncryptionOK]("BucketExtra.ReadEncryption", func() (any, error) {
		return op.siteClient.client.ReadBucketEncryption(ctx, v2.ReadBucketEncryptionParams{Name: op.bucket})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *bucketExtraOp) EnableEncryption(ctx context.Context, keyId string) error {
	return errorFromDecodedResponseE[v2.UpdateBucketEncryptionOK]("BucketExtra.EnableEncryption", func() (any, error) {
		return op.siteClient.client.UpdateBucketEncryption(ctx, &v2.HandlerUpdateBucketEncryptionRequestBody{KmsKeyID: v2.ResourceID(keyId)}, v2.UpdateBucketEncryptionParams{Name: op.bucket})
	})
}

func (op *bucketExtraOp) DisableEncryption(ctx context.Context) error {
	return errorFromDecodedResponseE[v2.DeleteBucketEncryptionNoContent]("BucketExtra.DisableEncryption", func() (any, error) {
		return op.siteClient.client.DeleteBucketEncryption(ctx, v2.DeleteBucketEncryptionParams{Name: op.bucket})
	})
}

func (op *bucketExtraOp) ReadReplication(ctx context.Context) (*v2.ModelReplication, error) {
	res, err := errorFromDecodedResponse[v2.HandlerReadBucketReplicationResponse]("BucketExtra.ReadReplication", func() (any, error) {
		return op.fedClient.client.ReadBucketReplication(ctx, v2.ReadBucketReplicationParams{Name: op.bucket})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *bucketExtraOp) EnableReplication(ctx context.Context, targetBucket string) (*v2.ModelReplication, error) {
	res, err := errorFromDecodedResponse[v2.HandlerCreateBucketReplicationResponse]("BucketExtra.EnableReplication", func() (any, error) {
		return op.fedClient.client.CreateBucketReplication(ctx, &v2.CreateBucketReplicationReq{DestBucket: targetBucket}, v2.CreateBucketReplicationParams{Name: op.bucket})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *bucketExtraOp) DisableReplication(ctx context.Context) error {
	return errorFromDecodedResponseE[v2.DeleteBucketReplicationNoContent]("BucketExtra.DisableReplication", func() (any, error) {
		return op.fedClient.client.DeleteBucketReplication(ctx, v2.DeleteBucketReplicationParams{Name: op.bucket})
	})
}

func (op *bucketExtraOp) ReadPenalty(ctx context.Context) (*v2.BucketPenaltyData, error) {
	res, err := errorFromDecodedResponse[v2.BucketPenalty]("BucketExtra.ReadPenalty", func() (any, error) {
		return op.siteClient.client.ReadBucketPenalty(ctx, v2.ReadBucketPenaltyParams{Name: v2.BucketName(op.bucket)})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *bucketExtraOp) ReadUsage(ctx context.Context) (*v2.BucketUsageData, error) {
	res, err := errorFromDecodedResponse[v2.BucketUsage]("BucketExtra.ReadUsage", func() (any, error) {
		return op.siteClient.client.ReadBucketUsage(ctx, v2.ReadBucketUsageParams{Name: v2.BucketName(op.bucket)})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}

func (op *bucketExtraOp) ReadQuota(ctx context.Context) (*v2.BucketQuotaData, error) {
	res, err := errorFromDecodedResponse[v2.BucketQuota]("BucketExtra.ReadQuota", func() (any, error) {
		return op.siteClient.client.ReadBucketQuota(ctx, v2.ReadBucketQuotaParams{Name: v2.BucketName(op.bucket)})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetData()), nil
}
