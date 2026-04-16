// Copyright 2025- The sacloud/workflows-api-go Authors
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

package workflows_test

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/sacloud/workflows-api-go"
	v1 "github.com/sacloud/workflows-api-go/apis/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	// NOTE: 課金プランが設定されていないと多くのAPIが402を返すため、E2Eテストの前に設定しておく。
	isE2ETest := os.Getenv("SAKURA_ACCESS_TOKEN") != "" && os.Getenv("SAKURA_ACCESS_TOKEN_SECRET") != ""
	if isE2ETest {
		ctx := context.Background()

		var theClient saclient.Client
		client, err := workflows.NewClient(&theClient)
		if err != nil {
			log.Fatalf("Error in TestMain. create saclient failed: %v", err)
		}
		subscriptionAPI := workflows.NewSubscriptionOp(client)

		respListPlans, err := subscriptionAPI.ListPlans(ctx)
		if err != nil {
			log.Fatalf("Error in TestMain. list plans failed: %v", err)
		}
		if respListPlans == nil || len(respListPlans.Plans) == 0 {
			log.Fatalf("Error in TestMain. list plans returned empty list: %v", err)
		}

		if err := subscriptionAPI.Create(ctx, v1.CreateSubscriptionReq{PlanId: respListPlans.Plans[0].ID}); err != nil {
			log.Fatalf("Error in TestMain. set Plan failed: %v", err)
		}
	}

	exitCode := m.Run()
	os.Exit(exitCode)
}

func TestSubscriptionAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	ctx := t.Context()

	var theClient saclient.Client
	client, err := workflows.NewClient(&theClient)
	require.NoError(t, err)

	subscriptionAPI := workflows.NewSubscriptionOp(client)

	// ListPlans
	respListPlans, err := subscriptionAPI.ListPlans(ctx)
	require.NoError(t, err)
	require.NotNil(t, respListPlans)

	// Delete
	err = subscriptionAPI.Delete(ctx)
	require.NoError(t, err)

	// Create
	err = subscriptionAPI.Create(ctx, v1.CreateSubscriptionReq{PlanId: respListPlans.Plans[0].ID})
	require.NoError(t, err)

	// Read
	respRead, err := subscriptionAPI.Read(ctx)
	require.NoError(t, err)
	require.NotNil(t, respRead)
	assert.Equal(t, respListPlans.Plans[0].ID, respRead.CurrentPlan.Value.PlanId)
}
