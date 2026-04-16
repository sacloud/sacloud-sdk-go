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
	"errors"

	v1 "github.com/sacloud/dedicated-storage-api-go/apis/v1"
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
	const methodName = "Contract.Create"

	res, err := op.client.DedicatedStorageContractsCreate(ctx, &req)
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return &res.DedicatedStorageContract, nil
}

func (op *contractOp) List(ctx context.Context) (*v1.DedicatedStorageContractsListResponse, error) {
	const methodName = "Contract.List"

	res, err := op.client.DedicatedStorageContractsList(ctx)
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return res, nil
}

func (op *contractOp) Read(ctx context.Context, id int64) (*v1.DedicatedStorageContract, error) {
	const methodName = "Contract.Read"

	res, err := op.client.DedicatedStorageContractsGet(ctx, v1.DedicatedStorageContractsGetParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return &res.DedicatedStorageContract, nil
}

func (op *contractOp) Update(ctx context.Context, id int64, request v1.UpdateDedicatedStorageContractRequest) (*v1.DedicatedStorageContract, error) {
	const methodName = "Contract.Update"

	res, err := op.client.DedicatedStorageContractsUpdate(ctx, &request, v1.DedicatedStorageContractsUpdateParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return &res.DedicatedStorageContract, nil
}

func (op *contractOp) Delete(ctx context.Context, id int64) error {
	const methodName = "Contract.Delete"

	err := op.client.DedicatedStorageContractsDelete(ctx, v1.DedicatedStorageContractsDeleteParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return NewAPIError(methodName, 0, err)
	}
	return nil
}

// PoolUsage implements ContractAPI.PoolUsage
func (op *contractOp) PoolUsage(ctx context.Context, id int64) (*v1.PoolUsageResponsePoolUsage, error) {
	const methodName = "Contract.PoolUsage"

	res, err := op.client.DedicatedStorageContractsPoolUsage(ctx, v1.DedicatedStorageContractsPoolUsageParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return &res.PoolUsage, nil
}

// ListDiskSnapshots implements ContractAPI.DiskSnapshots
func (op *contractOp) ListDiskSnapshots(ctx context.Context, id int64) (*v1.DiskSnapshotsListResponse, error) {
	const methodName = "Contract.DiskSnapshots"

	res, err := op.client.DedicatedStorageContractsListSnapshotsByContract(ctx, v1.DedicatedStorageContractsListSnapshotsByContractParams{ContractId: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return res, nil
}

func (op *contractOp) ListPlans(ctx context.Context) (*v1.DedicatedStorageContractPlanListResponse, error) {
	const methodName = "Contract.ListPlans"

	res, err := op.client.ProductPlansListPlans(ctx)
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return res, nil
}
func (op *contractOp) ReadPlan(ctx context.Context, planID int64) (*v1.DedicatedStorageContractPlan, error) {
	const methodName = "Contract.ReadPlan"

	res, err := op.client.ProductPlansGetPlans(ctx, v1.ProductPlansGetPlansParams{ID: planID})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return &res.DedicatedStorageContractPlan, nil
}
