// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package workernode_test

import (
	"net/http"
	"testing"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/workernode"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v interface{ Encode(*jx.Encoder) }, s ...int) (assert *require.Assertions, api WorkerNodeAPI) {
	assert = require.New(t)
	c, e := apprun_test.NewTestClient(v, s...)
	assert.NoError(e)
	api = NewWorkerNodeOp(c, v1.ClusterID(uuid.New()), v1.AutoScalingGroupID(uuid.New()))

	return
}

func TestList(t *testing.T) {
	next := v1.WorkerNodeID(uuid.New())
	var expected v1.ListWorkerNodesResponse
	expected.SetFake()
	expected.NextCursor.SetTo(next)
	expected.SetWorkerNodes(make([]v1.ReadWorkerNodeSummary, 3))
	for i := 0; i < len(expected.GetWorkerNodes()); i++ {
		var node v1.ReadWorkerNodeSummary
		node.SetFake()
		expected.WorkerNodes[i] = node
	}
	assert, api := setup(t, &expected)

	actual, cursor, err := api.List(t.Context(), 10, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.NotNil(cursor)
	assert.Len(actual, 3)
	assert.Equal(expected.GetWorkerNodes()[0].GetWorkerNodeID(), actual[0].WorkerNodeID)
}

func TestList_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, cursor, err := api.List(t.Context(), 0, nil)

	assert.Error(err)
	assert.Nil(actual)
	assert.Nil(cursor)
	assert.False(saclient.IsNotFoundError(err))
}

func TestRead(t *testing.T) {
	var expected v1.GetWorkerNodeResponse
	var node v1.ReadWorkerNodeDetail
	node.SetFake()
	id := node.GetWorkerNodeID()
	expected.SetWorkerNode(node)

	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), id)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(id, actual.WorkerNodeID)
}

func TestRead_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.WorkerNodeID(uuid.New())
	actual, err := api.Read(t.Context(), id)

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestUpdate(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	id := v1.WorkerNodeID(uuid.New())
	err := api.Update(t.Context(), id, true)

	assert.NoError(err)
}

func TestUpdate_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.WorkerNodeID(uuid.New())
	err := api.Update(t.Context(), id, false)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}
