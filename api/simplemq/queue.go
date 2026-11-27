// Copyright 2025- The sacloud/simplemq-api-go authors
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

package simplemq

import (
	"context"
	"strconv"

	"github.com/sacloud/sacloud-sdk-go/api/simplemq/apis/v1/queue"
)

type QueueAPI interface {
	List(context.Context) ([]queue.CommonServiceItem, error)
	Read(_ context.Context, id string) (*queue.CommonServiceItem, error)
	Create(context.Context, queue.CreateQueueRequest) (*queue.CommonServiceItem, error)
	Config(_ context.Context, id string, req queue.ConfigQueueRequest) (*queue.CommonServiceItem, error)
	Delete(_ context.Context, id string) error

	CountMessages(_ context.Context, id string) (int, error)
	RotateAPIKey(_ context.Context, id string) (string, error)
	ClearMessages(_ context.Context, id string) error
}

var _ QueueAPI = (*queueOp)(nil)

type queueOp struct {
	client *queue.Client
}

func NewQueueOp(client *queue.Client) QueueAPI {
	return &queueOp{client: client}
}

func GetQueueID(q *queue.CommonServiceItem) string {
	if q == nil {
		return ""
	}
	if q.ID.IsString() {
		return q.ID.String
	}
	return strconv.Itoa(q.ID.Int)
}

func GetQueueName(q *queue.CommonServiceItem) string {
	if q == nil {
		return ""
	}
	return q.Status.GetQueueName()
}

func (op *queueOp) Create(ctx context.Context, req queue.CreateQueueRequest) (*queue.CommonServiceItem, error) {
	req.CommonServiceItem.Provider.Class = queue.CreateQueueRequestCommonServiceItemProviderClassSimplemq
	res, err := errorFromDecodedResponse[queue.CreateQueueCreated]("Queue.Create", func() (any, error) { return op.client.CreateQueue(ctx, &req) })
	if err != nil {
		return nil, err
	}

	return new(res.GetCommonServiceItem()), nil
}

func (op *queueOp) List(ctx context.Context) ([]queue.CommonServiceItem, error) {
	res, err := errorFromDecodedResponse[queue.ListQueuesOK]("Queue.List", func() (any, error) { return op.client.ListQueues(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetCommonServiceItems(), nil
}

func (op *queueOp) Read(ctx context.Context, id string) (*queue.CommonServiceItem, error) {
	res, err := errorFromDecodedResponse[queue.ReadQueueOK]("Queue.Read", func() (any, error) { return op.client.ReadQueue(ctx, queue.ReadQueueParams{ID: id}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetCommonServiceItem()), nil
}

func (op *queueOp) Config(ctx context.Context, id string, req queue.ConfigQueueRequest) (*queue.CommonServiceItem, error) {
	res, err := errorFromDecodedResponse[queue.UpdateQueueOK]("Queue.Config", func() (any, error) { return op.client.UpdateQueue(ctx, &req, queue.UpdateQueueParams{ID: id}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetCommonServiceItem()), nil
}

func (op *queueOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[queue.DeleteQueueOK]("Queue.Delete", func() (any, error) { return op.client.DeleteQueue(ctx, queue.DeleteQueueParams{ID: id}) })
}

func (op *queueOp) CountMessages(ctx context.Context, id string) (int, error) {
	res, err := errorFromDecodedResponse[queue.ReadMessageCountOK]("Queue.CountMessages", func() (any, error) { return op.client.ReadMessageCount(ctx, queue.ReadMessageCountParams{ID: id}) })
	if err != nil {
		return 0, err
	}

	return new(res.GetSimpleMQ()).GetCount(), nil
}

func (op *queueOp) RotateAPIKey(ctx context.Context, id string) (string, error) {
	res, err := errorFromDecodedResponse[queue.RotateAPIKeyOK]("Queue.RotateAPIKey", func() (any, error) { return op.client.RotateAPIKey(ctx, queue.RotateAPIKeyParams{ID: id}) })
	if err != nil {
		return "", err
	}

	return new(res.GetSimpleMQ()).GetApikey(), nil
}

func (op *queueOp) ClearMessages(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[queue.ClearQueueOK]("Queue.ClearMessages", func() (any, error) { return op.client.ClearQueue(ctx, queue.ClearQueueParams{ID: id}) })
}
