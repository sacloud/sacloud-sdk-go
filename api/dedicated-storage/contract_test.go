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
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
)

func TestContract_CRUDL(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)
	if os.Getenv("TESTACC") != "1" {
		t.SkipNow()
	}
	ctx := t.Context()
	apiRootURL := DefaultAPIRootURL
	if os.Getenv("SAKURA_API_ROOT_URL") != "" {
		apiRootURL = os.Getenv("SAKURA_API_ROOT_URL")
	}

	var theClient saclient.Client
	client, err := NewClientWithAPIRootURL(&theClient, apiRootURL)
	if err != nil {
		t.Fatal(err)
	}

	contractOp := NewContractOp(client)
	var planID int64
	var targetStorage *v1.DedicatedStorageContract

	tests := []struct {
		name string
		test func(t *testing.T)
	}{
		{
			name: "list plans",
			test: func(t *testing.T) {
				plans, err := contractOp.ListPlans(ctx)
				if err != nil {
					t.Fatal(err)
				}

				if len(plans.DedicatedStorageContractPlans) < 1 {
					t.Fatal("got no plans from ContractOp::ListPlans()")
				}
				planID = plans.DedicatedStorageContractPlans[0].ID
			},
		},
		{
			name: "create",
			test: func(t *testing.T) {
				created, err := contractOp.Create(ctx, v1.CreateDedicatedStorageContractRequest{
					DedicatedStorageContract: v1.CreateDedicatedStorageContractRequestDedicatedStorageContract{
						Plan: v1.CreateDedicatedStorageContractRequestDedicatedStorageContractPlan{
							ID: planID,
						},
						Name:        "from-dedicated-storage-api-go",
						Description: "description",
						Tags:        []string{"tag1", "tag2"},
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				if created == nil {
					t.Fatal("got unexpected nil contract from ContractOp::Create()")
				}
				targetStorage = created
			},
		},
		{
			name: "read",
			test: func(t *testing.T) {
				read, err := contractOp.Read(ctx, targetStorage.ID)
				if err != nil {
					t.Fatal(err)
				}
				if read == nil {
					t.Fatal("got unexpected nil contract from ContractOp::Read()")
				}
			},
		},
		{
			name: "list",
			test: func(t *testing.T) {
				list, err := contractOp.List(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if list == nil {
					t.Fatal("got unexpected nil contract from ContractOp::List()")
				}
				found := false
				for _, v := range list.DedicatedStorageContracts {
					if v.ID == targetStorage.ID {
						found = true
						break
					}
				}
				if !found {
					t.Fatalf("created contract(ID: %d) not found in ContractOp::List()", targetStorage.ID)
				}
			},
		},
		{
			name: "update",
			test: func(t *testing.T) {
				updated, err := contractOp.Update(ctx, targetStorage.ID, v1.UpdateDedicatedStorageContractRequest{
					DedicatedStorageContract: v1.UpdateDedicatedStorageContractRequestDedicatedStorageContract{
						Name:        targetStorage.Name + "-updated",
						Description: targetStorage.Description + "-updated",
						Tags:        []string{"tag1-updated", "tag2-updated"},
						Icon:        v1.OptNilIcon{},
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				if updated == nil {
					t.Fatal("got unexpected nil contract from ContractOp::Update()")
				}
				targetStorage = updated
			},
		},
		{
			name: "delete",
			test: func(t *testing.T) {
				if err := contractOp.Delete(ctx, targetStorage.ID); err != nil {
					t.Fatal(err)
				}

				// 削除できたか一覧を取得して確認
				list, err := contractOp.List(ctx)
				if err != nil {
					t.Fatal(err)
				}
				if list == nil {
					t.Fatal("got unexpected nil contract from ContractOp::List()")
				}
				found := false
				for _, v := range list.DedicatedStorageContracts {
					if v.ID == targetStorage.ID {
						found = true
						break
					}
				}
				if found {
					t.Fatalf("deleted contract(ID: %d) found in ContractOp::List()", targetStorage.ID)
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

func TestContract_PoolUsage(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_DEDICATED_STORAGE_ID")(t)
	if os.Getenv("TESTACC") != "1" {
		t.SkipNow()
	}

	ctx := t.Context()

	apiRootURL := DefaultAPIRootURL
	if os.Getenv("SAKURA_API_ROOT_URL") != "" {
		apiRootURL = os.Getenv("SAKURA_API_ROOT_URL")
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

	contractOp := NewContractOp(client)
	usage, err := contractOp.PoolUsage(ctx, int64(dedicatedStorageID))
	if err != nil {
		t.Fatal(err)
	}
	if usage == nil {
		t.Fatal("got unexpected nil contract from ContractOp::PoolUsage()")
	}
}

func TestContract_ListDiskSnapshots(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_DEDICATED_STORAGE_ID")(t)
	if os.Getenv("TESTACC") != "1" {
		t.SkipNow()
	}

	ctx := t.Context()

	apiRootURL := DefaultAPIRootURL
	if os.Getenv("SAKURA_API_ROOT_URL") != "" {
		apiRootURL = os.Getenv("SAKURA_API_ROOT_URL")
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

	contractOp := NewContractOp(client)
	snapshots, err := contractOp.ListDiskSnapshots(ctx, int64(dedicatedStorageID))
	if err != nil {
		t.Fatal(err)
	}
	if snapshots == nil {
		t.Fatal("got unexpected nil contract from ContractOp::ListDiskSnapshots()")
	}
}
