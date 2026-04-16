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
	"os"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/accessor"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
)

func TestMain(m *testing.M) {
	testZone = testutil.TestZone()

	m.Run()

	skipCleanup := os.Getenv("SKIP_CLEANUP")
	if skipCleanup == "" {
		if err := testutil.CleanupTestResources(context.Background(), singletonAPICaller()); err != nil {
			panic(err)
		}
	}
}

var testZone string
var testIconID = types.ID(112901627749) // テスト用のアイコンID(shared icon)

func singletonAPICaller() iaas.APICaller {
	return testutil.SingletonAPICaller()
}

func isAccTest() bool {
	return testutil.IsAccTest()
}

func setupSwitchFunc(targetResource string, dests ...accessor.SwitchID) func(*testutil.CRUDTestContext, iaas.APICaller) error {
	return func(testContext *testutil.CRUDTestContext, caller iaas.APICaller) error {
		swClient := iaas.NewSwitchOp(caller)
		sw, err := swClient.Create(context.Background(), testZone, &iaas.SwitchCreateRequest{
			Name: testutil.ResourceName("switch-for-" + targetResource),
		})
		if err != nil {
			return err
		}

		testContext.Values[targetResource+"/switch"] = sw.ID
		for _, dest := range dests {
			dest.SetSwitchID(sw.ID)
		}
		return nil
	}
}

func cleanupSwitchFunc(targetResource string) func(*testutil.CRUDTestContext, iaas.APICaller) error {
	return func(testContext *testutil.CRUDTestContext, caller iaas.APICaller) error {
		switchID, ok := testContext.Values[targetResource+"/switch"]
		if !ok {
			return nil
		}

		swClient := iaas.NewSwitchOp(caller)
		return swClient.Delete(context.Background(), testZone, switchID.(types.ID))
	}
}
