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

type DiskAPI interface {
	CreateSnapshot(ctx context.Context, diskID int64, request *v1.CreateSnapshotRequest) (*v1.DiskSnapshot, error)
	ListSnapshots(ctx context.Context, diskID int64) (*v1.DiskSnapshotsListResponse, error)
	UpdateSnapshot(ctx context.Context, diskID, snapshotID int64, request *v1.UpdateSnapshotRequest) (*v1.DiskSnapshot, error)
	DeleteSnapshot(ctx context.Context, diskID, snapshotID int64) error
	RestoreFromSnapshot(ctx context.Context, diskID, snapshotID int64) error
	Expand(ctx context.Context, diskID int64, request *v1.ExpandDiskRequest) error
}

var _ DiskAPI = (*diskOp)(nil)

type diskOp struct {
	client *v1.Client
}

func NewDiskOp(client *v1.Client) DiskAPI {
	return &diskOp{
		client: client,
	}
}

func (op *diskOp) CreateSnapshot(ctx context.Context, diskID int64, request *v1.CreateSnapshotRequest) (*v1.DiskSnapshot, error) {
	const methodName = "Disk.CreateSnapshot"

	res, err := op.client.DisksCreateSnapshot(ctx, request, v1.DisksCreateSnapshotParams{DiskId: diskID})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return &res.DiskSnapshot, nil
}

func (op *diskOp) ListSnapshots(ctx context.Context, diskID int64) (*v1.DiskSnapshotsListResponse, error) {
	const methodName = "Disk.ListSnapshots"

	res, err := op.client.DisksListSnapshots(ctx, v1.DisksListSnapshotsParams{DiskId: diskID})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return res, nil
}

func (op *diskOp) UpdateSnapshot(ctx context.Context, diskID, snapshotID int64, request *v1.UpdateSnapshotRequest) (*v1.DiskSnapshot, error) {
	const methodName = "Disk.UpdateSnapshot"

	res, err := op.client.DisksUpdateSnapshot(ctx, request, v1.DisksUpdateSnapshotParams{DiskId: diskID, SnapshotId: snapshotID})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return nil, NewAPIError(methodName, 0, err)
	}

	return &res.DiskSnapshot, nil
}

func (op *diskOp) DeleteSnapshot(ctx context.Context, diskID, snapshotID int64) error {
	const methodName = "Disk.DeleteSnapshot"

	err := op.client.DisksDeleteSnapshot(ctx, v1.DisksDeleteSnapshotParams{DiskId: diskID, SnapshotId: snapshotID})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return NewAPIError(methodName, 0, err)
	}
	return nil
}

func (op *diskOp) RestoreFromSnapshot(ctx context.Context, diskID, snapshotID int64) error {
	const methodName = "Disk.RestoreFromSnapshot"

	err := op.client.DisksRestoreSnapshot(ctx, v1.DisksRestoreSnapshotParams{DiskId: diskID, SnapshotId: snapshotID})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return NewAPIError(methodName, 0, err)
	}
	return nil
}

func (op *diskOp) Expand(ctx context.Context, diskID int64, request *v1.ExpandDiskRequest) error {
	const methodName = "Disk.Expand"

	err := op.client.DisksExpand(ctx, request, v1.DisksExpandParams{ID: diskID})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return NewAPIError(methodName, e.StatusCode, errors.New(e.Response.ErrorMsg.Value))
		}
		return NewAPIError(methodName, 0, err)
	}

	return nil
}
