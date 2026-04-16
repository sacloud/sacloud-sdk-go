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
	"testing"
	"time"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/sacloud/workflows-api-go"
	v1 "github.com/sacloud/workflows-api-go/apis/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecutionAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	ctx := t.Context()

	var theClient saclient.Client
	client, err := workflows.NewClient(&theClient)
	require.NoError(t, err)

	// setup
	workflowAPI := workflows.NewWorkflowOp(client)
	workflow, err := workflowAPI.Create(ctx, v1.CreateWorkflowReq{
		Name:    "test-workflow",
		Runbook: sampleRunbook,
		Publish: true, // workflow executable only when published
		Logging: false,
	})
	require.NoError(t, err)
	require.NotNil(t, workflow)
	defer func() {
		time.Sleep(10 * time.Second) // wait for execution to be done/cancelled
		err := workflowAPI.Delete(ctx, workflow.ID)
		require.NoError(t, err)
	}()

	executionAPI := workflows.NewExecutionOp(client)

	// Create
	respCreate, err := executionAPI.Create(ctx, workflow.ID, v1.NewOptCreateExecutionReq(v1.CreateExecutionReq{
		Args: v1.NewOptString(`{"maxNumber": 100000}`), // run longer so that it can be cancelled later
	}))
	require.NoError(t, err)
	require.NotNil(t, respCreate)
	assert.Equal(t, workflow.ID, respCreate.Workflow.ID)
	assert.Equal(t, `{"maxNumber": 100000}`, respCreate.Args)

	// Read
	respRead, err := executionAPI.Read(ctx, workflow.ID, respCreate.ExecutionId)
	require.NoError(t, err)
	require.NotNil(t, respRead)
	assert.Equal(t, workflow.ID, respRead.Workflow.ID)
	assert.Equal(t, respCreate.CreatedAt, respRead.CreatedAt)

	// Cancel
	respCancel, err := executionAPI.Cancel(ctx, workflow.ID, respCreate.ExecutionId)
	require.NoError(t, err)
	require.NotNil(t, respCancel)
	assert.Equal(t, workflow.ID, respCreate.Workflow.ID)

	// List
	respList, err := executionAPI.List(ctx, v1.ListExecutionParams{
		ID: workflow.ID,
	})
	require.NoError(t, err)
	require.NotNil(t, respList)
	found := false
	for _, execution := range respList.Executions {
		if execution.ExecutionId == respCreate.ExecutionId {
			found = true
			assert.Equal(t, workflow.ID, execution.Workflow.ID)
			assert.Equal(t, respCreate.CreatedAt, execution.CreatedAt)
		}
	}
	assert.True(t, found, "Created Execution not found in list")

	// ListHistory
	respListHistory, err := executionAPI.ListHistory(ctx, v1.ListExecutionHistoryParams{
		ID:          workflow.ID,
		ExecutionId: respCreate.ExecutionId,
	})
	require.NoError(t, err)
	require.NotNil(t, respListHistory)

	// Delete
	// NOTE: Cancelが完了するまで409で失敗するので、完了するまでリトライする。
	require.EventuallyWithT(t, func(c *assert.CollectT) {
		err = executionAPI.Delete(ctx, workflow.ID, respCreate.ExecutionId)
		require.NoError(c, err)
	}, 30*time.Second, 5*time.Second, "Delete execution failed")
}
