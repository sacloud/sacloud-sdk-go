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

	"github.com/sacloud/simplemq-api-go/apis/v1/message"
)

type MessageAPI interface {
	Send(_ context.Context, content string) (*message.NewMessage, error)
	Receive(_ context.Context) ([]message.Message, error)
	ExtendTimeout(_ context.Context, messageID string) (*message.Message, error)
	Delete(_ context.Context, messageID string) error
}

var _ MessageAPI = (*messageOp)(nil)

type messageOp struct {
	queueName message.QueueName
	client    *message.Client
}

func NewMessageOp(client *message.Client, queueName string) MessageAPI {
	return &messageOp{
		queueName: message.QueueName(queueName),
		client:    client,
	}
}

func (op *messageOp) Send(ctx context.Context, content string) (*message.NewMessage, error) {
	res, err := op.client.SendMessage(ctx,
		&message.SendRequest{
			Content: message.MessageContent(content),
		},
		message.SendMessageParams{
			QueueName: op.queueName,
		})
	if err != nil {
		return nil, NewAPIError("Message.Send", 0, err)
	}

	switch r := res.(type) {
	case *message.SendMessageOK:
		return &r.Message, nil
	case *message.SendMessageUnauthorized:
		return nil, NewAPIError("Message.Send", 401, errors.New(r.Message.Value))
	case *message.SendMessageBadRequest:
		return nil, NewAPIError("Message.Send", 400, errors.New(r.Message.Value))
	case *message.SendMessageTooManyRequests:
		return nil, NewAPIError("Message.Send", 429, errors.New(r.Message.Value))
	case *message.SendMessageInternalServerError:
		return nil, NewAPIError("Message.Send", 500, errors.New(r.Message.Value))
	default:
		return nil, NewAPIError("Message.Send", 0, nil)
	}
}

func (op *messageOp) Receive(ctx context.Context) ([]message.Message, error) {
	res, err := op.client.ReceiveMessage(ctx,
		message.ReceiveMessageParams{
			QueueName: op.queueName,
		})
	if err != nil {
		return nil, NewAPIError("Message.Receive", 0, err)
	}

	switch r := res.(type) {
	case *message.ReceiveMessageOK:
		return r.Messages, nil
	case *message.ReceiveMessageUnauthorized:
		return nil, NewAPIError("Message.Receive", 401, errors.New(r.Message.Value))
	case *message.ReceiveMessageBadRequest:
		return nil, NewAPIError("Message.Receive", 400, errors.New(r.Message.Value))
	case *message.ReceiveMessageTooManyRequests:
		return nil, NewAPIError("Message.Receive", 429, errors.New(r.Message.Value))
	case *message.ReceiveMessageInternalServerError:
		return nil, NewAPIError("Message.Receive", 500, errors.New(r.Message.Value))
	default:
		return nil, NewAPIError("Message.Receive", 0, nil)
	}
}

func (op *messageOp) ExtendTimeout(ctx context.Context, messageID string) (*message.Message, error) {
	res, err := op.client.ExtendMessageTimeout(ctx,
		message.ExtendMessageTimeoutParams{
			QueueName: op.queueName,
			MessageId: message.MessageId(messageID),
		})
	if err != nil {
		return nil, NewAPIError("Message.ExtendTimeout", 0, err)
	}

	switch r := res.(type) {
	case *message.ExtendMessageTimeoutOK:
		return &r.Message, nil
	case *message.ExtendMessageTimeoutUnauthorized:
		return nil, NewAPIError("Message.ExtendTimeout", 401, errors.New(r.Message.Value))
	case *message.ExtendMessageTimeoutBadRequest:
		return nil, NewAPIError("Message.ExtendTimeout", 400, errors.New(r.Message.Value))
	case *message.ExtendMessageTimeoutNotFound:
		return nil, NewAPIError("Message.ExtendTimeout", 404, errors.New(r.Message.Value))
	case *message.ExtendMessageTimeoutTooManyRequests:
		return nil, NewAPIError("Message.ExtendTimeout", 429, errors.New(r.Message.Value))
	case *message.ExtendMessageTimeoutInternalServerError:
		return nil, NewAPIError("Message.ExtendTimeout", 500, errors.New(r.Message.Value))
	default:
		return nil, NewAPIError("Message.ExtendTimeout", 0, nil)
	}
}

func (op *messageOp) Delete(ctx context.Context, messageID string) error {
	res, err := op.client.DeleteMessage(ctx,
		message.DeleteMessageParams{
			QueueName: op.queueName,
			MessageId: message.MessageId(messageID),
		})
	if err != nil {
		return NewAPIError("Message.Delete", 0, err)
	}

	switch r := res.(type) {
	case *message.DeleteMessageOK:
		return nil
	case *message.DeleteMessageUnauthorized:
		return NewAPIError("Message.Delete", 401, errors.New(r.Message.Value))
	case *message.DeleteMessageBadRequest:
		return NewAPIError("Message.Delete", 400, errors.New(r.Message.Value))
	case *message.DeleteMessageNotFound:
		return NewAPIError("Message.Delete", 404, errors.New(r.Message.Value))
	case *message.DeleteMessageTooManyRequests:
		return NewAPIError("Message.Delete", 429, errors.New(r.Message.Value))
	case *message.DeleteMessageInternalServerError:
		return NewAPIError("Message.Delete", 500, errors.New(r.Message.Value))
	default:
		return NewAPIError("Message.Delete", 0, nil)
	}
}
