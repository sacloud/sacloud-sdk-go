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
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
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
	res, err := op.client.GetPlans(ctx)
	if err != nil {
		return nil, NewAPIError("Subscription.ListPlans", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetPlansOK:
		return p.Apigw.Plans, nil
	case *v1.ErrorSchema:
		return nil, NewAPIError("Subscription.ListPlans", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Subscription.ListPlans", 0, nil)
}

func (op *subscriptionOp) List(ctx context.Context) ([]v1.Subscription, error) {
	res, err := op.client.GetSubscriptions(ctx)
	if err != nil {
		return nil, NewAPIError("Subscription.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetSubscriptionsOK:
		return p.Apigw.Subscriptions, nil
	case *v1.GetSubscriptionsUnauthorized:
		return nil, NewAPIError("Subscription.List", 401, errors.New(p.Message.Value))
	case *v1.GetSubscriptionsInternalServerError:
		return nil, NewAPIError("Subscription.List", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Subscription.List", 0, nil)
}

func (op *subscriptionOp) Create(ctx context.Context, id uuid.UUID, name string) error {
	res, err := op.client.Subscribe(ctx, &v1.SubscriptionCreate{PlanId: id, Name: name})
	if err != nil {
		return NewAPIError("Subscription.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.SubscribeNoContent:
		return nil
	case *v1.SubscribeBadRequest:
		return NewAPIError("Subscription.Create", 400, errors.New(p.Message.Value))
	case *v1.SubscribeUnauthorized:
		return NewAPIError("Subscription.Create", 401, errors.New(p.Message.Value))
	case *v1.SubscribeInternalServerError:
		return NewAPIError("Subscription.Create", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Subscription.Create", 0, nil)
}

func (op *subscriptionOp) Read(ctx context.Context, id uuid.UUID) (*v1.SubscriptionDetailResponse, error) {
	res, err := op.client.GetSubscriptionById(ctx, v1.GetSubscriptionByIdParams{SubscriptionId: id})
	if err != nil {
		return nil, NewAPIError("Subscription.Read", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetSubscriptionByIdOK:
		return &p.Apigw.Subscription.Value, nil
	case *v1.GetSubscriptionByIdNotFound:
		return nil, NewAPIError("Subscription.Read", 404, errors.New(p.Message.Value))
	case *v1.GetSubscriptionByIdInternalServerError:
		return nil, NewAPIError("Subscription.Read", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Subscription.Read", 0, nil)
}

func (op *subscriptionOp) Update(ctx context.Context, id uuid.UUID, name string) error {
	res, err := op.client.UpdateSubscription(ctx, &v1.SubscriptionUpdate{Name: name},
		v1.UpdateSubscriptionParams{SubscriptionId: id})
	if err != nil {
		return NewAPIError("Subscription.Update", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpdateSubscriptionNoContent:
		return nil
	case *v1.UpdateSubscriptionBadRequest:
		return NewAPIError("Subscription.Update", 400, errors.New(p.Message.Value))
	case *v1.UpdateSubscriptionUnauthorized:
		return NewAPIError("Subscription.Update", 401, errors.New(p.Message.Value))
	case *v1.UpdateSubscriptionInternalServerError:
		return NewAPIError("Subscription.Update", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Subscription.Update", 0, nil)
}

func (op *subscriptionOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.Unsubscribe(ctx, v1.UnsubscribeParams{SubscriptionId: id})
	if err != nil {
		return NewAPIError("Subscription.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.UnsubscribeNoContent:
		return nil
	case *v1.UnsubscribeBadRequest:
		return NewAPIError("Subscription.Delete", 400, errors.New(p.Message.Value))
	case *v1.UnsubscribeNotFound:
		return NewAPIError("Subscription.Delete", 404, errors.New(p.Message.Value))
	case *v1.UnsubscribeInternalServerError:
		return NewAPIError("Subscription.Delete", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Subscription.Delete", 0, nil)
}
