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
	"errors"
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
	res, err := op.client.CreateQueue(ctx, &req)
	if err != nil {
		return nil, NewAPIError("Queue.Create", 0, err)
	}

	switch r := res.(type) {
	case *queue.CreateQueueCreated:
		return &r.CommonServiceItem, nil
	case *queue.CreateQueueUnauthorized:
		return nil, NewAPIError("Queue.Create", 401, errors.New(r.ErrorMsg.Value))
	case *queue.CreateQueueBadRequest:
		return nil, NewAPIError("Queue.Create", 400, errors.New(r.ErrorMsg.Value))
	case *queue.CreateQueueConflict:
		return nil, NewAPIError("Queue.Create", 409, errors.New(r.ErrorMsg.Value))
	case *queue.CreateQueueInternalServerError:
		return nil, NewAPIError("Queue.Create", 500, errors.New(r.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Queue.Create", 0, nil)
	}
}

func (op *queueOp) List(ctx context.Context) ([]queue.CommonServiceItem, error) {
	res, err := op.client.ListQueues(ctx)
	if err != nil {
		return nil, NewAPIError("Queue.List", 0, err)
	}

	switch r := res.(type) {
	case *queue.ListQueuesOK:
		return r.CommonServiceItems, nil
	case *queue.ListQueuesUnauthorized:
		return nil, NewAPIError("Queue.List", 401, errors.New(r.ErrorMsg.Value))
	case *queue.ListQueuesBadRequest:
		return nil, NewAPIError("Queue.List", 400, errors.New(r.ErrorMsg.Value))
	case *queue.ListQueuesInternalServerError:
		return nil, NewAPIError("Queue.List", 500, errors.New(r.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Queue.List", 0, nil)
	}
}

func (op *queueOp) Read(ctx context.Context, id string) (*queue.CommonServiceItem, error) {
	res, err := op.client.ReadQueue(ctx, queue.ReadQueueParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Queue.Read", 0, err)
	}

	switch r := res.(type) {
	case *queue.ReadQueueOK:
		return &r.CommonServiceItem, nil
	case *queue.ReadQueueUnauthorized:
		return nil, NewAPIError("Queue.Read", 401, errors.New(r.ErrorMsg.Value))
	case *queue.ReadQueueBadRequest:
		return nil, NewAPIError("Queue.Read", 400, errors.New(r.ErrorMsg.Value))
	case *queue.ReadQueueNotFound:
		return nil, NewAPIError("Queue.Read", 404, errors.New(r.ErrorMsg.Value))
	case *queue.ReadQueueInternalServerError:
		return nil, NewAPIError("Queue.Read", 500, errors.New(r.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Queue.Read", 0, nil)
	}
}

func (op *queueOp) Config(ctx context.Context, id string, req queue.ConfigQueueRequest) (*queue.CommonServiceItem, error) {
	res, err := op.client.UpdateQueue(ctx, &req, queue.UpdateQueueParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Queue.Config", 0, err)
	}

	switch r := res.(type) {
	case *queue.UpdateQueueOK:
		return &r.CommonServiceItem, nil
	case *queue.UpdateQueueUnauthorized:
		return nil, NewAPIError("Queue.Config", 401, errors.New(r.ErrorMsg.Value))
	case *queue.UpdateQueueBadRequest:
		return nil, NewAPIError("Queue.Config", 400, errors.New(r.ErrorMsg.Value))
	case *queue.UpdateQueueNotFound:
		return nil, NewAPIError("Queue.Config", 404, errors.New(r.ErrorMsg.Value))
	case *queue.UpdateQueueInternalServerError:
		return nil, NewAPIError("Queue.Config", 500, errors.New(r.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Queue.Config", 0, nil)
	}
}

func (op *queueOp) Delete(ctx context.Context, id string) error {
	res, err := op.client.DeleteQueue(ctx, queue.DeleteQueueParams{
		ID: id,
	})
	if err != nil {
		return NewError("Delete", err)
	}

	switch r := res.(type) {
	case *queue.DeleteQueueOK:
		return nil
	case *queue.DeleteQueueUnauthorized:
		return NewAPIError("Queue.Delete", 401, errors.New(r.ErrorMsg.Value))
	case *queue.DeleteQueueBadRequest:
		return NewAPIError("Queue.Delete", 400, errors.New(r.ErrorMsg.Value))
	case *queue.DeleteQueueNotFound:
		return NewAPIError("Queue.Delete", 404, errors.New(r.ErrorMsg.Value))
	case *queue.DeleteQueueInternalServerError:
		return NewAPIError("Queue.Delete", 500, errors.New(r.ErrorMsg.Value))
	default:
		return NewAPIError("Queue.Delete", 0, nil)
	}
}

func (op *queueOp) CountMessages(ctx context.Context, id string) (int, error) {
	res, err := op.client.ReadMessageCount(ctx, queue.ReadMessageCountParams{ID: id})
	if err != nil {
		return 0, NewError("CountMessages", err)
	}

	switch r := res.(type) {
	case *queue.ReadMessageCountOK:
		return r.SimpleMQ.GetCount(), nil
	case *queue.ReadMessageCountUnauthorized:
		return 0, NewAPIError("Queue.CountMessages", 401, errors.New(r.ErrorMsg.Value))
	case *queue.ReadMessageCountBadRequest:
		return 0, NewAPIError("Queue.CountMessages", 400, errors.New(r.ErrorMsg.Value))
	case *queue.ReadMessageCountNotFound:
		return 0, NewAPIError("Queue.CountMessages", 404, errors.New(r.ErrorMsg.Value))
	case *queue.ReadMessageCountInternalServerError:
		return 0, NewAPIError("Queue.CountMessages", 500, errors.New(r.ErrorMsg.Value))
	default:
		return 0, NewAPIError("Queue.CountMessages", 0, nil)
	}
}

func (op *queueOp) RotateAPIKey(ctx context.Context, id string) (string, error) {
	res, err := op.client.RotateAPIKey(ctx, queue.RotateAPIKeyParams{ID: id})
	if err != nil {
		return "", NewAPIError("Queue.RotateAPIKey", 0, err)
	}

	switch r := res.(type) {
	case *queue.RotateAPIKeyOK:
		return r.SimpleMQ.GetApikey(), nil
	case *queue.RotateAPIKeyUnauthorized:
		return "", NewAPIError("Queue.RotateAPIKey", 401, errors.New(r.ErrorMsg.Value))
	case *queue.RotateAPIKeyBadRequest:
		return "", NewAPIError("Queue.RotateAPIKey", 400, errors.New(r.ErrorMsg.Value))
	case *queue.RotateAPIKeyNotFound:
		return "", NewAPIError("Queue.RotateAPIKey", 404, errors.New(r.ErrorMsg.Value))
	case *queue.RotateAPIKeyInternalServerError:
		return "", NewAPIError("Queue.RotateAPIKey", 500, errors.New(r.ErrorMsg.Value))
	default:
		return "", NewAPIError("Queue.RotateAPIKey", 0, nil)
	}
}

func (op *queueOp) ClearMessages(ctx context.Context, id string) error {
	res, err := op.client.ClearQueue(ctx, queue.ClearQueueParams{ID: id})
	if err != nil {
		return NewAPIError("Queue.ClearMessages", 0, err)
	}

	switch r := res.(type) {
	case *queue.ClearQueueOK:
		return nil
	case *queue.ClearQueueUnauthorized:
		return NewAPIError("Queue.ClearMessages", 401, errors.New(r.ErrorMsg.Value))
	case *queue.ClearQueueBadRequest:
		return NewAPIError("Queue.ClearMessages", 400, errors.New(r.ErrorMsg.Value))
	case *queue.ClearQueueNotFound:
		return NewAPIError("Queue.ClearMessages", 404, errors.New(r.ErrorMsg.Value))
	case *queue.ClearQueueInternalServerError:
		return NewAPIError("Queue.ClearMessages", 500, errors.New(r.ErrorMsg.Value))
	default:
		return NewAPIError("Queue.ClearMessages", 0, nil)
	}
}
