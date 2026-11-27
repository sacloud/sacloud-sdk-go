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

package workflows

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/workflows/apis/v1"
)

type SubscriptionAPI interface {
	ListPlans(ctx context.Context) (*v1.ListPlansOK, error)
	Read(ctx context.Context) (*v1.ReadSubscriptionOK, error)
	Create(ctx context.Context, request v1.CreateSubscriptionReq) error
	Delete(ctx context.Context) error
}

var _ SubscriptionAPI = (*subscriptionOp)(nil)

type subscriptionOp struct {
	client *v1.Client
}

func NewSubscriptionOp(client *v1.Client) SubscriptionAPI {
	return &subscriptionOp{client: client}
}

func (op *subscriptionOp) ListPlans(ctx context.Context) (*v1.ListPlansOK, error) {
	return errorFromDecodedResponse[v1.ListPlansOK]("Subscription.ListPlans", func() (any, error) { return op.client.ListPlans(ctx) })
}

func (op *subscriptionOp) Read(ctx context.Context) (*v1.ReadSubscriptionOK, error) {
	return errorFromDecodedResponse[v1.ReadSubscriptionOK]("Subscription.Read", func() (any, error) { return op.client.ReadSubscription(ctx) })
}

func (op *subscriptionOp) Create(ctx context.Context, req v1.CreateSubscriptionReq) error {
	return errorFromDecodedResponseE[v1.CreateSubscriptionNoContent]("Subscription.Create", func() (any, error) { return op.client.CreateSubscription(ctx, &req) })
}

func (op *subscriptionOp) Delete(ctx context.Context) error {
	return errorFromDecodedResponseE[v1.DeleteSubscriptionNoContent]("Subscription.Delete", func() (any, error) { return op.client.DeleteSubscription(ctx) })
}
