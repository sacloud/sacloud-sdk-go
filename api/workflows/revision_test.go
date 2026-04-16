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

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/sacloud/workflows-api-go"
	v1 "github.com/sacloud/workflows-api-go/apis/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRevisionAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	ctx := t.Context()

	var theClient saclient.Client
	client, err := workflows.NewClient(&theClient)
	require.NoError(t, err)

	workflowAPI := workflows.NewWorkflowOp(client)

	// setup
	workflow, err := workflowAPI.Create(ctx, v1.CreateWorkflowReq{
		Name:    "test-workflow",
		Runbook: sampleRunbook,
		Publish: false,
		Logging: false,
	})
	require.NoError(t, err)
	require.NotNil(t, workflow)
	defer func() {
		err := workflowAPI.Delete(ctx, workflow.ID)
		require.NoError(t, err)
	}()

	revisionAPI := workflows.NewRevisionOp(client)

	// Create
	respCreate, err := revisionAPI.Create(ctx, workflow.ID, v1.CreateWorkflowRevisionReq{
		Runbook:       sampleRunbook,
		RevisionAlias: v1.NewOptString("v1"),
	})
	require.NoError(t, err)
	require.NotNil(t, respCreate)
	assert.Equal(t, workflow.ID, respCreate.WorkflowId)
	assert.Equal(t, v1.NewOptString("v1"), respCreate.RevisionAlias)
	assert.Equal(t, sampleRunbook, respCreate.Runbook)

	// Read
	respRead, err := revisionAPI.Read(ctx, workflow.ID, respCreate.RevisionId)
	require.NoError(t, err)
	require.NotNil(t, respRead)
	assert.Equal(t, workflow.ID, respRead.WorkflowId)
	assert.Equal(t, v1.NewOptString("v1"), respRead.RevisionAlias)

	// List
	respList, err := revisionAPI.List(ctx, v1.ListWorkflowRevisionsParams{
		ID: workflow.ID,
	})
	require.NoError(t, err)
	found := false
	for _, revision := range respList.Revisions {
		if revision.RevisionId == respCreate.RevisionId {
			found = true
			assert.Equal(t, workflow.ID, revision.WorkflowId)
			assert.Equal(t, v1.NewOptString("v1"), revision.RevisionAlias)
			assert.Equal(t, respCreate.CreatedAt, revision.CreatedAt)
		}
	}
	assert.True(t, found, "Created Revision not found in list")

	// UpdateAlias
	respUpdate, err := revisionAPI.UpdateAlias(ctx, workflow.ID, respCreate.RevisionId, v1.UpdateWorkflowRevisionAliasReq{
		RevisionAlias: "v2",
	})
	require.NoError(t, err)
	require.NotNil(t, respUpdate)
	assert.Equal(t, workflow.ID, respUpdate.WorkflowId)
	assert.Equal(t, v1.NewOptString("v2"), respUpdate.RevisionAlias)
	assert.Equal(t, respCreate.CreatedAt, respUpdate.CreatedAt)

	// DeleteAlias
	err = revisionAPI.DeleteAlias(ctx, workflow.ID, respCreate.RevisionId)
	require.NoError(t, err)
}
