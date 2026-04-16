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
	"os"
	"strconv"
	"testing"

	v1 "github.com/sacloud/dedicated-storage-api-go/apis/v1"
	"github.com/sacloud/packages-go/size"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
)

func TestDisk_SnapShotCRUDL(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_DISK_ID", "SAKURA_DEDICATED_STORAGE_ID")(t)
	if os.Getenv("TESTACC") != "1" {
		t.SkipNow()
	}
	ctx := t.Context()

	apiRootURL := DefaultAPIRootURL
	if os.Getenv("SAKURA_API_ROOT_URL") != "" {
		apiRootURL = os.Getenv("SAKURA_API_ROOT_URL")
	}

	diskID, err := strconv.Atoi(os.Getenv("SAKURA_DISK_ID"))
	if err != nil {
		t.Fatal(err)
	}
	dedicatedStorageID, err := strconv.Atoi(os.Getenv("SAKURA_DEDICATED_STORAGE_ID"))
	if err != nil {
		t.Fatal(err)
	}

	var theClient saclient.Client
	client, err := NewClientWithAPIRootURL(&theClient, apiRootURL)
	if err != nil {
		t.Fatal(err)
	}

	diskOp := NewDiskOp(client)
	var snapshot *v1.DiskSnapshot

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "create",
			test: func(t *testing.T) {
				created, err := diskOp.CreateSnapshot(ctx, int64(diskID), &v1.CreateSnapshotRequest{
					DiskSnapshot: v1.CreateSnapshotRequestDiskSnapshot{
						DedicatedStorageContract: v1.CreateSnapshotRequestDiskSnapshotDedicatedStorageContract{
							ID: int64(dedicatedStorageID),
						},
						Name:        "from-dedicated-storage-api-go",
						Description: "description",
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				if created == nil {
					t.Fatal("got unexpected nil contract from SnapshotOp::Create()")
				}
				snapshot = created
			},
		},
		{
			name: "list",
			test: func(t *testing.T) {
				list, err := diskOp.ListSnapshots(ctx, int64(diskID))
				if err != nil {
					t.Fatal(err)
				}
				if list == nil {
					t.Fatal("got unexpected nil contract from SnapshotOp::List()")
				}
				found := false
				for _, v := range list.DiskSnapshots {
					if v.ID == snapshot.ID {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("created snapshot(ID: %d) not found in SnapshotOp::List()", snapshot.ID)
				}
			},
		},
		{
			name: "update",
			test: func(t *testing.T) {
				updated, err := diskOp.UpdateSnapshot(ctx, int64(diskID), snapshot.ID, &v1.UpdateSnapshotRequest{
					DiskSnapshot: v1.UpdateSnapshotRequestDiskSnapshot{
						Name:        snapshot.Name + "-updated",
						Description: snapshot.Description + "-updated",
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				if updated == nil {
					t.Fatal("got unexpected nil contract from SnapshotOp::Update()")
				}
				snapshot = updated
			},
		},
		{
			name: "restore",
			test: func(t *testing.T) {
				if err := diskOp.RestoreFromSnapshot(ctx, int64(diskID), snapshot.ID); err != nil {
					t.Fatal(err)
				}
			},
		},
		{
			name: "delete",
			test: func(t *testing.T) {
				if err := diskOp.DeleteSnapshot(ctx, int64(diskID), snapshot.ID); err != nil {
					t.Fatal(err)
				}

				// 削除できたか一覧を取得して確認
				list, err := diskOp.ListSnapshots(ctx, int64(diskID))
				if err != nil {
					t.Fatal(err)
				}
				if list == nil {
					t.Fatal("got unexpected nil contract from SnapshotOp::List()")
				}
				found := false
				for _, v := range list.DiskSnapshots {
					if v.ID == snapshot.ID {
						found = true
						break
					}
				}
				if found {
					t.Fatalf("deleted contract(ID: %d) found in SnapshotOp::List()", snapshot.ID)
				}
			},
		},
	}

	for _, tt := range tests {
		if !t.Run(tt.name, tt.test) {
			t.FailNow()
		}
	}
}

func TestDisk_Expand(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_DISK_ID")(t)
	if os.Getenv("TESTACC") != "1" {
		t.SkipNow()
	}
	ctx := t.Context()

	apiRootURL := DefaultAPIRootURL
	if os.Getenv("SAKURA_API_ROOT_URL") != "" {
		apiRootURL = os.Getenv("SAKURA_API_ROOT_URL")
	}

	diskID, err := strconv.Atoi(os.Getenv("SAKURA_DISK_ID"))
	if err != nil {
		t.Fatal(err)
	}

	var theClient saclient.Client
	client, err := NewClientWithAPIRootURL(&theClient, apiRootURL)
	if err != nil {
		t.Fatal(err)
	}

	diskOp := NewDiskOp(client)

	err = diskOp.Expand(ctx, int64(diskID), &v1.ExpandDiskRequest{
		ExpanedSizeMB: 40 * size.GiB,
	})
	if err != nil {
		t.Fatal(err)
	}
}
