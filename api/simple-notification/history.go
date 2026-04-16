// Copyright 2026- The sacloud/simple-notification-api-go Authors
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

package simplenotification

import (
	"context"
	"errors"

	v1 "github.com/sacloud/simple-notification-api-go/apis/v1"
)

type HistoryAPI interface {
	List(ctx context.Context) (*v1.ListSimpleNotificationHistoriesResponse, error)
	Read(ctx context.Context, id string) (*v1.GetSimpleNotificationHistoryResponse, error)
}

var _ HistoryAPI = (*HistoryOp)(nil)

type HistoryOp struct {
	client *v1.Client
}

func NewHistoryOp(client *v1.Client) HistoryAPI {
	return &HistoryOp{client: client}
}

func (o *HistoryOp) List(ctx context.Context) (*v1.ListSimpleNotificationHistoriesResponse, error) {
	const methodName = "History.List"
	res, err := o.client.ListNotificationHistories(ctx)
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, NewError(methodName, err)
	}
	return res, nil
}

func (o *HistoryOp) Read(ctx context.Context, id string) (*v1.GetSimpleNotificationHistoryResponse, error) {
	const methodName = "History.Read"
	res, err := o.client.GetNotificationHistory(ctx, v1.GetNotificationHistoryParams{RequestID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, NewError(methodName, err)
	}
	return res, nil
}
