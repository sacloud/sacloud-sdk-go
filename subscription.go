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
	"errors"
	"net/http"

	v1 "github.com/sacloud/workflows-api-go/apis/v1"
)

type SubscriptionAPI interface {
	ListPlans(ctx context.Context) (*v1.ListPlansOK, error)
	Read(ctx context.Context) (*v1.GetSubscriptionOK, error)
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
	const methodName = "Subscription.ListPlans"

	res, err := op.client.ListPlans(ctx)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.ListPlansOK:
		return r, nil
	case *v1.ListPlansBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.ListPlansUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.ListPlansForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.ListPlansNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.ListPlansInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subscriptionOp) Read(ctx context.Context) (*v1.GetSubscriptionOK, error) {
	const methodName = "Subscription.Read"

	res, err := op.client.GetSubscription(ctx)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.GetSubscriptionOK:
		return r, nil
	case *v1.GetSubscriptionBadRequest:
		return nil, NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.GetSubscriptionUnauthorized:
		return nil, NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.GetSubscriptionForbidden:
		return nil, NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.GetSubscriptionNotFound:
		return nil, NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.GetSubscriptionInternalServerError:
		return nil, NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subscriptionOp) Create(ctx context.Context, req v1.CreateSubscriptionReq) error {
	const methodName = "Subscription.Create"

	res, err := op.client.CreateSubscription(ctx, &req)
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.CreateSubscriptionNoContent:
		return nil
	case *v1.CreateSubscriptionBadRequest:
		return NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.CreateSubscriptionUnauthorized:
		return NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.CreateSubscriptionForbidden:
		return NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.CreateSubscriptionNotFound:
		return NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.CreateSubscriptionInternalServerError:
		return NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *subscriptionOp) Delete(ctx context.Context) error {
	const methodName = "Subscription.Delete"

	res, err := op.client.DeleteSubscription(ctx)
	if err != nil {
		return NewAPIError(methodName, 0, err)
	}

	switch r := res.(type) {
	case *v1.DeleteSubscriptionNoContent:
		return nil
	case *v1.DeleteSubscriptionBadRequest:
		return NewAPIError(methodName, http.StatusBadRequest, errors.New(r.Message))
	case *v1.DeleteSubscriptionUnauthorized:
		return NewAPIError(methodName, http.StatusUnauthorized, errors.New(r.Message))
	case *v1.DeleteSubscriptionForbidden:
		return NewAPIError(methodName, http.StatusForbidden, errors.New(r.Message))
	case *v1.DeleteSubscriptionNotFound:
		return NewAPIError(methodName, http.StatusNotFound, errors.New(r.Message))
	case *v1.DeleteSubscriptionInternalServerError:
		return NewAPIError(methodName, http.StatusInternalServerError, errors.New(r.Message))
	default:
		return NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
