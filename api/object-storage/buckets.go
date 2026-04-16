// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"errors"
	"fmt"

	v2 "github.com/sacloud/object-storage-api-go/apis/v2"
)

type BucketAPI interface {
	List(ctx context.Context) ([]v2.BucketListDataItem, error)
	Create(ctx context.Context, params *BucketCreateParams) (*v2.ModelBucket, error)
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

func (op *bucketOp) List(ctx context.Context) ([]v2.BucketListDataItem, error) {
	res, err := op.siteClient.client.ListBuckets(ctx)
	if err != nil {
		return nil, NewAPIError("Buckets.List", 0, err)
	}

	switch r := res.(type) {
	case *v2.BucketList:
		return r.Data, nil
	case *v2.Error401:
		return nil, NewAPIError("Buckets.List", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Buckets.List", 0, errors.New("unknown error"))
	}
}

type BucketCreateParams struct {
	Bucket string
	SiteId string
	Plan   string
}

func createRequest(params *BucketCreateParams) *v2.HandlerPutBucketReqBody {
	res := &v2.HandlerPutBucketReqBody{
		ClusterID: params.SiteId,
	}

	switch params.SiteId {
	case "isk01", "tky01":
		res.Plan = v2.NewOptNilHandlerPutBucketReqBodyPlan(v2.HandlerPutBucketReqBodyPlan{
			Type:             v2.NewOptModelPlanType(v2.ModelPlanType("standard")),
			ServiceClassPath: v2.NewOptServiceClassPath(v2.ServiceClassPath(fmt.Sprintf("objectstorage/%s/bucket", params.SiteId))),
		})
	case "arc02":
		res.Plan = v2.NewOptNilHandlerPutBucketReqBodyPlan(v2.HandlerPutBucketReqBodyPlan{
			Type:             v2.NewOptModelPlanType(v2.ModelPlanType("archive")),
			ServiceClassPath: v2.NewOptServiceClassPath(v2.ServiceClassPath(fmt.Sprintf("objectstorage/%s/bucket/%s", params.SiteId, params.Plan))),
		})
	}

	return res
}

func (op *bucketOp) Create(ctx context.Context, params *BucketCreateParams) (*v2.ModelBucket, error) {
	res, err := op.fedClient.client.CreateBucket(ctx, createRequest(params), v2.CreateBucketParams{Name: params.Bucket})
	if err != nil {
		return nil, NewAPIError("Buckets.Create", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerPutBucketRes:
		return &r.Data.Value, nil
	case *v2.Error400:
		return nil, NewAPIError("Buckets.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("Buckets.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error409:
		return nil, NewAPIError("Buckets.Create", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("Buckets.Create", 0, errors.New("unknown error"))
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
	case *v2.Error400:
		return NewAPIError("Buckets.Delete", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error409:
		return NewAPIError("Buckets.Delete", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return NewAPIError("Buckets.Delete", 0, errors.New("unknown error"))
	}
}

type BucketExtraAPI interface {
	ReadEncryption(ctx context.Context) (*v2.HandlerEncryptionConfigRes, error)
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

func (op *bucketExtraOp) ReadEncryption(ctx context.Context) (*v2.HandlerEncryptionConfigRes, error) {
	res, err := op.siteClient.client.GetBucketEncryption(ctx, v2.GetBucketEncryptionParams{Name: op.bucket})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadEncryption", 0, err)
	}

	switch r := res.(type) {
	case *v2.GetBucketEncryptionOK:
		return &r.Data, nil
	case *v2.Error404:
		return nil, NewAPIError("BucketExtra.ReadEncryption", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("BucketExtra.ReadEncryption", 0, errors.New("unknown error"))
	}
}

func (op *bucketExtraOp) EnableEncryption(ctx context.Context, keyId string) error {
	res, err := op.siteClient.client.PutBucketEncryption(ctx, &v2.HandlerEncryptionConfigReqBody{
		KmsKeyID: v2.ResourceID(keyId),
	}, v2.PutBucketEncryptionParams{Name: op.bucket})
	if err != nil {
		return NewAPIError("BucketExtra.EnableEncryption", 0, err)
	}

	switch r := res.(type) {
	case *v2.PutBucketEncryptionOK:
		return nil
	case *v2.Error400:
		return NewAPIError("BucketExtra.EnableEncryption", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return NewAPIError("BucketExtra.EnableEncryption", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return NewAPIError("BucketExtra.EnableEncryption", 0, errors.New("unknown error"))
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
	case *v2.Error404:
		return NewAPIError("BucketExtra.DisableEncryption", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return NewAPIError("BucketExtra.DisableEncryption", 0, errors.New("unknown error"))
	}
}

func (op *bucketExtraOp) ReadReplication(ctx context.Context) (*v2.ModelReplication, error) {
	res, err := op.fedClient.client.GetBucketReplication(ctx, v2.GetBucketReplicationParams{Name: op.bucket})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadReplication", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerGetReplicationRes:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("BucketExtra.ReadReplication", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("BucketExtra.ReadReplication", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("BucketExtra.ReadReplication", 0, errors.New("unknown error"))
	}
}

func (op *bucketExtraOp) EnableReplication(ctx context.Context, targetBucket string) (*v2.ModelReplication, error) {
	res, err := op.fedClient.client.PostBucketReplication(ctx, &v2.PostBucketReplicationReq{
		DestBucket: targetBucket,
	}, v2.PostBucketReplicationParams{Name: op.bucket})
	if err != nil {
		return nil, NewAPIError("BucketExtra.EnableReplication", 0, err)
	}

	switch r := res.(type) {
	case *v2.HandlerPostReplicationRes:
		return &r.Data.Value, nil
	case *v2.Error400:
		return nil, NewAPIError("BucketExtra.EnableReplication", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error401:
		return nil, NewAPIError("BucketExtra.EnableReplication", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("BucketExtra.EnableReplication", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("BucketExtra.EnableReplication", 0, errors.New("unknown error"))
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
	case *v2.Error401:
		return NewAPIError("BucketExtra.DisableReplication", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return NewAPIError("BucketExtra.DisableReplication", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return NewAPIError("BucketExtra.DisableReplication", 0, errors.New("unknown error"))
	}
}

func (op *bucketExtraOp) ReadPenalty(ctx context.Context) (*v2.BucketPenaltyData, error) {
	res, err := op.siteClient.client.GetBucketPenalty(ctx, v2.GetBucketPenaltyParams{Name: v2.BucketName(op.bucket)})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadPenalty", 0, err)
	}

	switch r := res.(type) {
	case *v2.BucketPenalty:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("BucketExtra.ReadPenalty", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("BucketExtra.ReadPenalty", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("BucketExtra.ReadPenalty", 0, errors.New("unknown error"))
	}
}

func (op *bucketExtraOp) ReadUsage(ctx context.Context) (*v2.BucketUsageData, error) {
	res, err := op.siteClient.client.GetBucketUsage(ctx, v2.GetBucketUsageParams{Name: v2.BucketName(op.bucket)})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadUsage", 0, err)
	}

	switch r := res.(type) {
	case *v2.BucketUsage:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("BucketExtra.ReadUsage", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("BucketExtra.ReadUsage", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("BucketExtra.ReadUsage", 0, errors.New("unknown error"))
	}
}

func (op *bucketExtraOp) ReadQuota(ctx context.Context) (*v2.BucketQuotaData, error) {
	res, err := op.siteClient.client.GetBucketQuota(ctx, v2.GetBucketQuotaParams{Name: v2.BucketName(op.bucket)})
	if err != nil {
		return nil, NewAPIError("BucketExtra.ReadQuota", 0, err)
	}

	switch r := res.(type) {
	case *v2.BucketQuota:
		return &r.Data.Value, nil
	case *v2.Error401:
		return nil, NewAPIError("BucketExtra.ReadQuota", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	case *v2.Error404:
		return nil, NewAPIError("BucketExtra.ReadQuota", int(r.Error.Value.Code.Value), errors.New(string(r.Error.Value.Message.Value)))
	default:
		return nil, NewAPIError("BucketExtra.ReadQuota", 0, errors.New("unknown error"))
	}
}
