// Copyright 2022-2025 The sacloud/simplemq-api-go Authors
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

package simplemq_test

import (
	"context"
	"testing"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var queueClient saclient.Client

func TestQueueAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	ctx := context.Background()

	client, err := simplemq.NewQueueClient(&queueClient)
	require.NoError(t, err)
	queueOp := simplemq.NewQueueOp(client)

	// CreateQueue
	resCreate, err := queueOp.Create(ctx, queue.CreateQueueRequest{
		CommonServiceItem: queue.CreateQueueRequestCommonServiceItem{
			Name:        "SDK-Test-Queue",
			Description: queue.NewOptString("SDK-Test-Queueの概要"),
		},
	})
	require.NoError(t, err)
	assert.Equal(t, "SDK-Test-Queue", resCreate.Status.QueueName)

	queueID := simplemq.GetQueueID(resCreate)

	// ListQueues
	resList, err := queueOp.List(ctx)
	assert.NoError(t, err)
	found := false
	for _, q := range resList {
		if queueID == simplemq.GetQueueID(&q) {
			found = true
			assert.Equal(t, "SDK-Test-Queue", resCreate.Status.QueueName)
		}
	}
	assert.True(t, found, "Created queue not found in list")

	// ConfigQueue
	resConfig, err := queueOp.Config(ctx, queueID, queue.ConfigQueueRequest{
		CommonServiceItem: queue.ConfigQueueRequestCommonServiceItem{
			Description: queue.NewOptString("SDK-Test-Queueの概要を変更"),
			Settings: queue.Settings{
				VisibilityTimeoutSeconds: 99,
				ExpireSeconds:            resCreate.Settings.ExpireSeconds,
			},
			Tags: []string{"tag1", "tag2"},
			Icon: queue.NewOptNilConfigQueueRequestCommonServiceItemIcon(queue.ConfigQueueRequestCommonServiceItemIcon{
				ID: queue.NewOptConfigQueueRequestCommonServiceItemIconID(queue.NewStringConfigQueueRequestCommonServiceItemIconID("112901627751")),
			}),
		},
	})
	assert.NoError(t, err)
	assert.Equal(t, queue.VisibilityTimeoutSeconds(99), resConfig.Settings.VisibilityTimeoutSeconds)

	// ReadQueue
	resRead, err := queueOp.Read(ctx, queueID)
	assert.NoError(t, err)
	assert.Equal(t, "SDK-Test-Queueの概要を変更", resRead.Description.Value)

	// GetMessageCount
	resMessageCount, err := queueOp.CountMessages(ctx, queueID)
	assert.NoError(t, err)
	assert.Equal(t, 0, resMessageCount)

	// RotateAPIKey
	_, err = queueOp.RotateAPIKey(ctx, queueID)
	assert.NoError(t, err)

	// ClearQueue
	err = queueOp.ClearMessages(ctx, queueID)
	assert.NoError(t, err)

	// DeleteQueue
	err = queueOp.Delete(ctx, queueID)
	require.NoError(t, err)
}
