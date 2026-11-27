// Copyright 2022-2025 The sacloud/dedicated-storage-api-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dedicatedstorage

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/dedicated-storage/apis/v1"
)

type ContractAPI interface {
	Create(ctx context.Context, request v1.CreateDedicatedStorageContractRequest) (*v1.DedicatedStorageContract, error)
	List(ctx context.Context) (*v1.DedicatedStorageContractsListResponse, error)
	Read(ctx context.Context, id int64) (*v1.DedicatedStorageContract, error)
	Update(ctx context.Context, id int64, request v1.UpdateDedicatedStorageContractRequest) (*v1.DedicatedStorageContract, error)
	Delete(ctx context.Context, id int64) error

	PoolUsage(ctx context.Context, id int64) (*v1.PoolUsageResponsePoolUsage, error)
	ListDiskSnapshots(ctx context.Context, id int64) (*v1.DiskSnapshotsListResponse, error)

	ListPlans(ctx context.Context) (*v1.DedicatedStorageContractPlanListResponse, error)
	ReadPlan(ctx context.Context, planID int64) (*v1.DedicatedStorageContractPlan, error)
}

var _ ContractAPI = (*contractOp)(nil)

type contractOp struct {
	client *v1.Client
}

func NewContractOp(client *v1.Client) ContractAPI {
	return &contractOp{client: client}
}

func (op *contractOp) Create(ctx context.Context, req v1.CreateDedicatedStorageContractRequest) (*v1.DedicatedStorageContract, error) {
	res, err := errorFromDecodedResponse("Contract.Create", func() (*v1.DedicatedStorageContractResponse, error) {
		return op.client.DedicatedStorageContractsCreate(ctx, &req)
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetDedicatedStorageContract()), nil
}

func (op *contractOp) List(ctx context.Context) (*v1.DedicatedStorageContractsListResponse, error) {
	return errorFromDecodedResponse("Contract.List", func() (*v1.DedicatedStorageContractsListResponse, error) {
		return op.client.DedicatedStorageContractsList(ctx)
	})
}

func (op *contractOp) Read(ctx context.Context, id int64) (*v1.DedicatedStorageContract, error) {
	res, err := errorFromDecodedResponse("Contract.Read", func() (*v1.DedicatedStorageContractResponse, error) {
		return op.client.DedicatedStorageContractsGet(ctx, v1.DedicatedStorageContractsGetParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetDedicatedStorageContract()), nil
}

func (op *contractOp) Update(ctx context.Context, id int64, request v1.UpdateDedicatedStorageContractRequest) (*v1.DedicatedStorageContract, error) {
	res, err := errorFromDecodedResponse("Contract.Update", func() (*v1.DedicatedStorageContractResponse, error) {
		return op.client.DedicatedStorageContractsUpdate(ctx, &request, v1.DedicatedStorageContractsUpdateParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetDedicatedStorageContract()), nil
}

func (op *contractOp) Delete(ctx context.Context, id int64) error {
	return errorFromDecodedResponseE("Contract.Delete", func() error {
		return op.client.DedicatedStorageContractsDelete(ctx, v1.DedicatedStorageContractsDeleteParams{ID: id})
	})
}

// PoolUsage implements ContractAPI.PoolUsage
func (op *contractOp) PoolUsage(ctx context.Context, id int64) (*v1.PoolUsageResponsePoolUsage, error) {
	res, err := errorFromDecodedResponse("Contract.PoolUsage", func() (*v1.PoolUsageResponse, error) {
		return op.client.DedicatedStorageContractsPoolUsage(ctx, v1.DedicatedStorageContractsPoolUsageParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetPoolUsage()), nil
}

// ListDiskSnapshots implements ContractAPI.DiskSnapshots
func (op *contractOp) ListDiskSnapshots(ctx context.Context, id int64) (*v1.DiskSnapshotsListResponse, error) {
	return errorFromDecodedResponse("Contract.DiskSnapshots", func() (*v1.DiskSnapshotsListResponse, error) {
		return op.client.DedicatedStorageContractsListSnapshotsByContract(ctx, v1.DedicatedStorageContractsListSnapshotsByContractParams{ContractId: id})
	})
}

func (op *contractOp) ListPlans(ctx context.Context) (*v1.DedicatedStorageContractPlanListResponse, error) {
	return errorFromDecodedResponse("Contract.ListPlans", func() (*v1.DedicatedStorageContractPlanListResponse, error) {
		return op.client.ProductPlansListPlans(ctx)
	})
}

func (op *contractOp) ReadPlan(ctx context.Context, planID int64) (*v1.DedicatedStorageContractPlan, error) {
	res, err := errorFromDecodedResponse("Contract.ReadPlan", func() (*v1.DedicatedStorageContractPlanResponse, error) {
		return op.client.ProductPlansGetPlans(ctx, v1.ProductPlansGetPlansParams{ID: planID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetDedicatedStorageContractPlan()), nil
}
