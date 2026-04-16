// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package test

import (
	"context"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/stretchr/testify/assert"
)

func TestDiskPlanOp_Find(t *testing.T) {
	t.Parallel()

	client := iaas.NewDiskPlanOp(singletonAPICaller())

	searched, err := client.Find(context.Background(), iaas.APIDefaultZone, &iaas.FindCondition{Count: 1})
	assert.NoError(t, err)

	err = testutil.DoAsserts(
		testutil.AssertLenFunc(t, searched.DiskPlans, 1, "DiskPlans"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].ID, "DiskPlans.ID"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].Name, "DiskPlans.Name"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].Availability, "DiskPlans.Availability"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].StorageClass, "DiskPlans.StorageClass"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].Size, "DiskPlans.Size"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].Size[0].Availability, "DiskPlans.Size.Availability"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].Size[0].DisplaySize, "DiskPlans.Size.DisplaySize"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].Size[0].DisplaySuffix, "DiskPlans.Size.DisplaySuffix"),
		testutil.AssertNotEmptyFunc(t, searched.DiskPlans[0].Size[0].SizeMB, "DiskPlans.Size.SizeMB"),
	)
	assert.NoError(t, err)
}
