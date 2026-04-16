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
	"time"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/sacloud/simplemq-api-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/message"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var messageClient saclient.Client

func TestMessageAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	ctx := context.Background()

	qClient, err := simplemq.NewQueueClient(&queueClient)
	require.NoError(t, err)
	queueOp := simplemq.NewQueueOp(qClient)

	resCreate, err := queueOp.Create(ctx, queue.CreateQueueRequest{
		CommonServiceItem: queue.CreateQueueRequestCommonServiceItem{
			Name: "SDK-Test-Queue",
		},
	})
	require.NoError(t, err)
	defer func() {
		_ = queueOp.Delete(ctx, simplemq.GetQueueID(resCreate))
	}()
	queueName := simplemq.GetQueueName(resCreate)
	apiKey, err := queueOp.RotateAPIKey(ctx, simplemq.GetQueueID(resCreate))
	assert.NoError(t, err)

	client, err := simplemq.NewMessageClient(apiKey, &messageClient)
	require.NoError(t, err)
	messageOp := simplemq.NewMessageOp(client, queueName)

	// SendMessage
	resSend, err := messageOp.Send(ctx, "HelloSimpleMQ")
	assert.NoError(t, err)
	assert.Equal(t, message.MessageContent("HelloSimpleMQ"), resSend.Content)

	// ReceiveMessage
	resReceive, err := messageOp.Receive(ctx)
	assert.NoError(t, err)
	assert.NotEmpty(t, resReceive)
	assert.Equal(t, message.MessageContent("HelloSimpleMQ"), resReceive[0].Content)
	messageID := string(resReceive[0].ID)

	// ExtendMessageTimeout
	resExtend, err := messageOp.ExtendTimeout(ctx, messageID)
	assert.NoError(t, err)
	timeBefore := time.UnixMilli(resReceive[0].VisibilityTimeoutAt)
	timeExtended := time.UnixMilli(resExtend.VisibilityTimeoutAt)
	assert.True(t, timeExtended.After(timeBefore))

	// DeleteMessage
	err = messageOp.Delete(ctx, messageID)
	require.NoError(t, err)
}
