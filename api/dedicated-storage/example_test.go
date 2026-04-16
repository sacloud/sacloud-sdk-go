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

package dedicatedstorage_test

import (
	"context"
	"fmt"
	"os"

	dedicatedstorage "github.com/sacloud/dedicated-storage-api-go"
	v1 "github.com/sacloud/dedicated-storage-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

func Example_dedicatedStorageCRUDL() {
	// setup
	// TODO replace your access token/secret
	os.Setenv("SAKURA_ACCESS_TOKEN", "your-token")  //nolint:errcheck,gosec
	os.Setenv("SAKURA_ACCESS_TOKEN", "your-secret") //nolint:errcheck,gosec

	theClient := saclient.Client{}
	client, err := dedicatedstorage.NewClient(&theClient)
	if err != nil {
		panic(err)
	}

	contractOp := dedicatedstorage.NewContractOp(client)
	ctx := context.Background()

	// list plans & choose plan ID
	plans, err := contractOp.ListPlans(ctx)
	if err != nil {
		panic(err)
	}
	planID := plans.DedicatedStorageContractPlans[0].ID // choose the first plan ID for example

	// create
	created, err := contractOp.Create(ctx, v1.CreateDedicatedStorageContractRequest{
		DedicatedStorageContract: v1.CreateDedicatedStorageContractRequestDedicatedStorageContract{
			Plan: v1.CreateDedicatedStorageContractRequestDedicatedStorageContractPlan{
				ID: planID,
			},
			Name:        "example-name",
			Description: "example-description",
			Tags:        []string{"example1", "example2"},
			// Icon:        v1.NewOptNilIcon(v1.Icon{ID: 111111111111}),
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(created)

	// read
	read, err := contractOp.Read(ctx, created.ID)
	if err != nil {
		panic(err)
	}
	fmt.Println(read)

	// list
	listed, err := contractOp.List(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println(listed)

	// update
	updated, err := contractOp.Update(ctx, created.ID, v1.UpdateDedicatedStorageContractRequest{
		DedicatedStorageContract: v1.UpdateDedicatedStorageContractRequestDedicatedStorageContract{
			Name:        "example-name-updated",
			Description: "example-description-updated",
			Tags:        []string{"example1", "example2"},
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(updated)

	// delete
	if err := contractOp.Delete(ctx, created.ID); err != nil {
		panic(err)
	}
}
