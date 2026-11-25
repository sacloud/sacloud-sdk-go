// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"errors"
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
	res, err := op.siteClient.client.ListBuckets(ctx)
	if err != nil {
		return nil, NewAPIError("Buckets.List", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerListBucketsResponse:
		return r.Data, nil
	case *v2.ErrorResponse:
		return nil, NewAPIError("Buckets.List", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return nil, unknownError("Buckets.List", r)
	}
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
	res, err := op.fedClient.client.CreateBucket(ctx, createRequest(params), v2.CreateBucketParams{Name: params.Bucket})
	if err != nil {
		return nil, NewAPIError("Buckets.Create", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerCreateBucketResponse:
		return &r.Data, nil
	case *v2.CreateBucketBadRequest:
		return nil, NewAPIError("Buckets.Create", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.CreateBucketNotFound:
		return nil, NewAPIError("Buckets.Create", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.CreateBucketConflict:
		return nil, NewAPIError("Buckets.Create", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return nil, unknownError("Buckets.Create", r)
	}
}

func (op *bucketOp) Delete(ctx context.Context, bucketName string) error {
	res, err := op.fedClient.client.DeleteBucket(ctx, v2.DeleteBucketParams{Name: bucketName})
	if err != nil {
		return NewAPIError("Buckets.Delete", 0, err)
	}

	switch r := res.(type) {
	case *v2.DeleteBucketNoContent:
		return nil
	case *v2.DeleteBucketBadRequest:
		return NewAPIError("Buckets.Delete", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.DeleteBucketConflict:
		return NewAPIError("Buckets.Delete", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return unknownError("Buckets.Delete", r)
	}
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
	res, err := op.siteClient.client.ReadBucketEncryption(ctx, v2.ReadBucketEncryptionParams{Name: op.bucket})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadEncryption", 0, err)
	}

	switch r := res.(type) {
	case *v2.ReadBucketEncryptionOK:
		return &r.Data, nil
	case *v2.ReadBucketEncryptionUnauthorized:
		return nil, NewAPIError("BucketExtra.ReadEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ReadBucketEncryptionNotFound:
		return nil, NewAPIError("BucketExtra.ReadEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ReadBucketEncryptionMethodNotAllowed:
		return nil, NewAPIError("BucketExtra.ReadEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return nil, unknownError("BucketExtra.ReadEncryption", r)
	}
}

func (op *bucketExtraOp) EnableEncryption(ctx context.Context, keyId string) error {
	res, err := op.siteClient.client.UpdateBucketEncryption(ctx, &v2.HandlerUpdateBucketEncryptionRequestBody{
		KmsKeyID: v2.ResourceID(keyId),
	}, v2.UpdateBucketEncryptionParams{Name: op.bucket})
	if err != nil {
		return NewAPIError("BucketExtra.EnableEncryption", 0, err)
	}

	switch r := res.(type) {
	case *v2.UpdateBucketEncryptionOK:
		return nil
	case *v2.UpdateBucketEncryptionBadRequest:
		return NewAPIError("BucketExtra.EnableEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.UpdateBucketEncryptionUnauthorized:
		return NewAPIError("BucketExtra.EnableEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.UpdateBucketEncryptionNotFound:
		return NewAPIError("BucketExtra.EnableEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.UpdateBucketEncryptionMethodNotAllowed:
		return NewAPIError("BucketExtra.EnableEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return unknownError("BucketExtra.EnableEncryption", r)
	}
}

func (op *bucketExtraOp) DisableEncryption(ctx context.Context) error {
	res, err := op.siteClient.client.DeleteBucketEncryption(ctx, v2.DeleteBucketEncryptionParams{Name: op.bucket})
	if err != nil {
		return NewAPIError("BucketExtra.DisableEncryption", 0, err)
	}

	switch r := res.(type) {
	case *v2.DeleteBucketEncryptionNoContent:
		return nil
	case *v2.DeleteBucketEncryptionUnauthorized:
		return NewAPIError("BucketExtra.DisableEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.DeleteBucketEncryptionNotFound:
		return NewAPIError("BucketExtra.DisableEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.DeleteBucketEncryptionMethodNotAllowed:
		return NewAPIError("BucketExtra.DisableEncryption", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return unknownError("BucketExtra.DisableEncryption", r)
	}
}

func (op *bucketExtraOp) ReadReplication(ctx context.Context) (*v2.ModelReplication, error) {
	res, err := op.fedClient.client.ReadBucketReplication(ctx, v2.ReadBucketReplicationParams{Name: op.bucket})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadReplication", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerReadBucketReplicationResponse:
		return &r.Data, nil
	case *v2.ReadBucketReplicationUnauthorized:
		return nil, NewAPIError("BucketExtra.ReadReplication", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ReadBucketReplicationNotFound:
		return nil, NewAPIError("BucketExtra.ReadReplication", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return nil, unknownError("BucketExtra.ReadReplication", r)
	}
}

func (op *bucketExtraOp) EnableReplication(ctx context.Context, targetBucket string) (*v2.ModelReplication, error) {
	res, err := op.fedClient.client.CreateBucketReplication(ctx, &v2.CreateBucketReplicationReq{
		DestBucket: targetBucket,
	}, v2.CreateBucketReplicationParams{Name: op.bucket})
	if err != nil {
		return nil, NewAPIError("BucketExtra.EnableReplication", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerCreateBucketReplicationResponse:
		return &r.Data, nil
	case *v2.CreateBucketReplicationBadRequest:
		return nil, NewAPIError("BucketExtra.EnableReplication", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.CreateBucketReplicationUnauthorized:
		return nil, NewAPIError("BucketExtra.EnableReplication", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.CreateBucketReplicationNotFound:
		return nil, NewAPIError("BucketExtra.EnableReplication", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return nil, unknownError("BucketExtra.EnableReplication", r)
	}
}

func (op *bucketExtraOp) DisableReplication(ctx context.Context) error {
	res, err := op.fedClient.client.DeleteBucketReplication(ctx, v2.DeleteBucketReplicationParams{Name: op.bucket})
	if err != nil {
		return NewAPIError("BucketExtra.DisableReplication", 0, err)
	}

	switch r := res.(type) {
	case *v2.DeleteBucketReplicationNoContent:
		return nil
	case *v2.DeleteBucketReplicationUnauthorized:
		return NewAPIError("BucketExtra.DisableReplication", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.DeleteBucketReplicationNotFound:
		return NewAPIError("BucketExtra.DisableReplication", int(r.Error.Code), errors.New(string(r.Error.Message)))
	default:
		return unknownError("BucketExtra.DisableReplication", r)
	}
}

func (op *bucketExtraOp) ReadPenalty(ctx context.Context) (*v2.BucketPenaltyData, error) {
	res, err := op.siteClient.client.ReadBucketPenalty(ctx, v2.ReadBucketPenaltyParams{Name: v2.BucketName(op.bucket)})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadPenalty", 0, err)
	}

	switch r := res.(type) {
	case *v2.BucketPenalty:
		return &r.Data, nil
	case *v2.ReadBucketPenaltyUnauthorized:
		return nil, NewAPIError("BucketExtra.ReadPenalty", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ReadBucketPenaltyNotFound:
		return nil, NewAPIError("BucketExtra.ReadPenalty", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ErrorResponseStatusCode:
		return nil, NewAPIError("BucketExtra.ReadPenalty", r.StatusCode, errors.New(string(r.Response.Error.Message)))
	default:
		return nil, unknownError("BucketExtra.ReadPenalty", r)
	}
}

func (op *bucketExtraOp) ReadUsage(ctx context.Context) (*v2.BucketUsageData, error) {
	res, err := op.siteClient.client.ReadBucketUsage(ctx, v2.ReadBucketUsageParams{Name: v2.BucketName(op.bucket)})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadUsage", 0, err)
	}

	switch r := res.(type) {
	case *v2.BucketUsage:
		return &r.Data, nil
	case *v2.ReadBucketUsageUnauthorized:
		return nil, NewAPIError("BucketExtra.ReadUsage", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ReadBucketUsageNotFound:
		return nil, NewAPIError("BucketExtra.ReadUsage", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ErrorResponseStatusCode:
		return nil, NewAPIError("BucketExtra.ReadUsage", r.StatusCode, errors.New(string(r.Response.Error.Message)))
	default:
		return nil, unknownError("BucketExtra.ReadUsage", r)
	}
}

func (op *bucketExtraOp) ReadQuota(ctx context.Context) (*v2.BucketQuotaData, error) {
	res, err := op.siteClient.client.ReadBucketQuota(ctx, v2.ReadBucketQuotaParams{Name: v2.BucketName(op.bucket)})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadQuota", 0, err)
	}

	switch r := res.(type) {
	case *v2.BucketQuota:
		return &r.Data, nil
	case *v2.ReadBucketQuotaBadRequest:
		return nil, NewAPIError("BucketExtra.ReadQuota", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ReadBucketQuotaUnauthorized:
		return nil, NewAPIError("BucketExtra.ReadQuota", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ReadBucketQuotaNotFound:
		return nil, NewAPIError("BucketExtra.ReadQuota", int(r.Error.Code), errors.New(string(r.Error.Message)))
	case *v2.ErrorResponseStatusCode:
		return nil, NewAPIError("BucketExtra.ReadQuota", r.StatusCode, errors.New(string(r.Response.Error.Message)))
	default:
		return nil, unknownError("BucketExtra.ReadQuota", r)
	}
}
