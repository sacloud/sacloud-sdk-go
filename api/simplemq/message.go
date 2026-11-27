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

	"github.com/sacloud/sacloud-sdk-go/api/simplemq/apis/v1/message"
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
	res, err := errorFromDecodedResponse[message.SendMessageOK]("Message.Send", func() (any, error) {
		return op.client.SendMessage(ctx, &message.SendRequest{Content: message.MessageContent(content)}, message.SendMessageParams{QueueName: op.queueName})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetMessage()), nil
}

func (op *messageOp) Receive(ctx context.Context) ([]message.Message, error) {
	res, err := errorFromDecodedResponse[message.ReceiveMessageOK]("Message.Receive", func() (any, error) {
		return op.client.ReceiveMessage(ctx, message.ReceiveMessageParams{QueueName: op.queueName})
	})
	if err != nil {
		return nil, err
	}

	return res.GetMessages(), nil
}

func (op *messageOp) ExtendTimeout(ctx context.Context, messageID string) (*message.Message, error) {
	res, err := errorFromDecodedResponse[message.ExtendMessageTimeoutOK]("Message.ExtendTimeout", func() (any, error) {
		return op.client.ExtendMessageTimeout(ctx, message.ExtendMessageTimeoutParams{QueueName: op.queueName, MessageId: message.MessageId(messageID)})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetMessage()), nil
}

func (op *messageOp) Delete(ctx context.Context, messageID string) error {
	return errorFromDecodedResponseE[message.DeleteMessageOK]("Message.Delete", func() (any, error) {
		return op.client.DeleteMessage(ctx, message.DeleteMessageParams{QueueName: op.queueName, MessageId: message.MessageId(messageID)})
	})
}
