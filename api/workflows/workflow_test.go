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

func TestWorkflowAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	ctx := t.Context()

	var theClient saclient.Client
	client, err := workflows.NewClient(&theClient)
	require.NoError(t, err)

	workflowAPI := workflows.NewWorkflowOp(client)

	// CreateWorkflow
	respCreate, err := workflowAPI.Create(ctx, v1.CreateWorkflowReq{
		Name:    "test-workflow",
		Runbook: sampleRunbook,
		Publish: false,
		Logging: false,
	})
	require.NoError(t, err)
	require.NotNil(t, respCreate)
	assert.Equal(t, "test-workflow", respCreate.Name)

	defer func() {
		// DeleteWorkflow
		err := workflowAPI.Delete(ctx, respCreate.ID)
		require.NoError(t, err)
	}()

	// ReadWorkflow
	respRead, err := workflowAPI.Read(ctx, respCreate.ID)
	require.NoError(t, err)
	require.NotNil(t, respRead)
	assert.Equal(t, "test-workflow", respRead.Name)
	assert.Equal(t, v1.NewOptString(""), respRead.Description) // empty
	assert.False(t, respRead.Publish)
	assert.False(t, respRead.Logging)

	// ListWorkflows
	respList, err := workflowAPI.List(ctx, v1.ListWorkflowParams{})
	require.NoError(t, err)
	found := false
	for _, workflow := range respList.Workflows {
		if workflow.ID == respCreate.ID {
			found = true
			assert.Equal(t, "test-workflow", workflow.Name)
			assert.Equal(t, v1.NewOptString(""), workflow.Description) // empty
			assert.False(t, workflow.Publish)
			assert.False(t, workflow.Logging)
		}
	}
	assert.True(t, found, "Created Workflow not found in list")

	// ListWorkflowSuggest
	respListSuggest, err := workflowAPI.ListSuggest(ctx, v1.ListWorkflowSuggestParams{Name: "test-workflow"})
	require.NoError(t, err)
	found = false
	for _, workflow := range respListSuggest.Suggests {
		if workflow.Name == respCreate.Name {
			found = true
			assert.Equal(t, "test-workflow", workflow.Name)
		}
	}
	assert.True(t, found, "Created Workflow not found in suggest list")

	// UpdateWorkflow
	respUpdate, err := workflowAPI.Update(ctx, respCreate.ID, v1.UpdateWorkflowReq{
		Name:        v1.NewOptString("test-workflow-updated"),
		Description: v1.NewOptString("test workflow updated"),
		Publish:     v1.NewOptBool(true),
		Logging:     v1.NewOptBool(true),
	})
	require.NoError(t, err)
	require.NotNil(t, respUpdate)
	assert.Equal(t, "test-workflow-updated", respUpdate.Name)
	assert.Equal(t, v1.NewOptString("test workflow updated"), respUpdate.Description)
	assert.True(t, respUpdate.Publish)
	assert.True(t, respUpdate.Logging)
}

const sampleRunbook = `
meta:
  description: エラトステネスの篩
args:
  maxNumber:
    type: number
    description: 素数を求める最大の数
steps:
  setup:
    assign:
      sieve: ${array.fill(array.range(args.maxNumber), true)}
      primes: []
  initial:
    assign:
      _a: ${array.set(sieve, 0, false)}
      _b: ${array.set(sieve, 1, false)}
  loop:
    for:
      in: ${array.range(2, math.ceil(math.sqrt(args.maxNumber)))}
      as: index
      steps:
        if:
          switch:
            # falseだったら飛ばす
            - condition: ${sieve[index] == false}
              next: continue
            # trueだったら素数
            - condition: ${sieve[index] != false}
              steps:
                # 素数の倍数を篩にかける
                updateSieve:
                  for:
                    in: ${array.range(index * 2, args.maxNumber, index)}
                    as: n
                    steps:
                      set:
                        assign:
                          _a: ${array.set(sieve, n, false)}
        continue:
  printPrimes:
    for:
      in: ${array.range(2, args.maxNumber)}
      as: index
      steps:
        if:
          switch:
            - condition: ${sieve[index] == true}
              steps:
                push:
                  assign:
                    _a: ${array.push(primes, index)}
                log:
                  assign:
                    log: '${"素数: " + index}'
  done:
    return: ${primes}
`
