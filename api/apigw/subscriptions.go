// Copyright 2025- The sacloud/apigw-api-go authors
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

package apigw

import (
	"context"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/sacloud/sacloud-sdk-go/api/apigw/apis/v1"
)

type SubscriptionAPI interface {
	ListPlans(ctx context.Context) ([]v1.Plan, error)
	List(ctx context.Context) ([]v1.Subscription, error)
	Create(ctx context.Context, id uuid.UUID, name string) error
	Read(ctx context.Context, id uuid.UUID) (*v1.SubscriptionDetailResponse, error)
	Update(ctx context.Context, id uuid.UUID, name string) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ SubscriptionAPI = (*subscriptionOp)(nil)

type subscriptionOp struct {
	client *v1.Client
}

func NewSubscriptionOp(client *v1.Client) SubscriptionAPI {
	return &subscriptionOp{client: client}
}

func (op *subscriptionOp) ListPlans(ctx context.Context) ([]v1.Plan, error) {
	res, err := errorFromDecodedResponse[v1.GetPlansOK]("Subscription.ListPlans", func() (any, error) { return op.client.GetPlans(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetPlans(), nil
}

func (op *subscriptionOp) List(ctx context.Context) ([]v1.Subscription, error) {
	res, err := errorFromDecodedResponse[v1.GetSubscriptionsOK]("Subscription.List", func() (any, error) { return op.client.GetSubscriptions(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetSubscriptions(), nil
}

func (op *subscriptionOp) Create(ctx context.Context, id uuid.UUID, name string) error {
	return errorFromDecodedResponseE[v1.SubscribeNoContent]("Subscription.Create", func() (any, error) {
		return op.client.Subscribe(ctx, &v1.SubscriptionCreate{PlanId: id, Name: name})
	})
}

func (op *subscriptionOp) Read(ctx context.Context, id uuid.UUID) (*v1.SubscriptionDetailResponse, error) {
	res, err := errorFromDecodedResponse[v1.GetSubscriptionByIdOK]("Subscription.Read", func() (any, error) {
		return op.client.GetSubscriptionById(ctx, v1.GetSubscriptionByIdParams{SubscriptionId: id})
	})
	if err != nil {
		return nil, err
	}

	s, ok := new(res.GetApigw()).GetSubscription().Get()
	if !ok {
		return nil, NewError("Subscription.Read", errors.New("response does not contain a subscription"))
	}

	return &s, nil
}

func (op *subscriptionOp) Update(ctx context.Context, id uuid.UUID, name string) error {
	return errorFromDecodedResponseE[v1.UpdateSubscriptionNoContent]("Subscription.Update", func() (any, error) {
		return op.client.UpdateSubscription(ctx, &v1.SubscriptionUpdate{Name: name}, v1.UpdateSubscriptionParams{SubscriptionId: id})
	})
}

func (op *subscriptionOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.UnsubscribeNoContent]("Subscription.Delete", func() (any, error) {
		return op.client.Unsubscribe(ctx, v1.UnsubscribeParams{SubscriptionId: id})
	})
}
