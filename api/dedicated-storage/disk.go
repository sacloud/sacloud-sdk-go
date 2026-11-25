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
	res, err := errorFromDecodedResponse("Disk.CreateSnapshot", func() (*v1.DiskSnapshotResponse, error) {
		return op.client.DisksCreateSnapshot(ctx, request, v1.DisksCreateSnapshotParams{DiskId: diskID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetDiskSnapshot()), nil
}

func (op *diskOp) ListSnapshots(ctx context.Context, diskID int64) (*v1.DiskSnapshotsListResponse, error) {
	return errorFromDecodedResponse("Disk.ListSnapshots", func() (*v1.DiskSnapshotsListResponse, error) {
		return op.client.DisksListSnapshots(ctx, v1.DisksListSnapshotsParams{DiskId: diskID})
	})
}

func (op *diskOp) UpdateSnapshot(ctx context.Context, diskID, snapshotID int64, request *v1.UpdateSnapshotRequest) (*v1.DiskSnapshot, error) {
	res, err := errorFromDecodedResponse("Disk.UpdateSnapshot", func() (*v1.DiskSnapshotResponse, error) {
		return op.client.DisksUpdateSnapshot(ctx, request, v1.DisksUpdateSnapshotParams{DiskId: diskID, SnapshotId: snapshotID})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetDiskSnapshot()), nil
}

func (op *diskOp) DeleteSnapshot(ctx context.Context, diskID, snapshotID int64) error {
	return errorFromDecodedResponseE("Disk.DeleteSnapshot", func() error {
		return op.client.DisksDeleteSnapshot(ctx, v1.DisksDeleteSnapshotParams{DiskId: diskID, SnapshotId: snapshotID})
	})
}

func (op *diskOp) RestoreFromSnapshot(ctx context.Context, diskID, snapshotID int64) error {
	return errorFromDecodedResponseE("Disk.RestoreFromSnapshot", func() error {
		return op.client.DisksRestoreSnapshot(ctx, v1.DisksRestoreSnapshotParams{DiskId: diskID, SnapshotId: snapshotID})
	})
}

func (op *diskOp) Expand(ctx context.Context, diskID int64, request *v1.ExpandDiskRequest) error {
	return errorFromDecodedResponseE("Disk.Expand", func() error {
		return op.client.DisksExpand(ctx, request, v1.DisksExpandParams{ID: diskID})
	})
}
