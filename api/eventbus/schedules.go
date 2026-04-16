// Copyright 2025- The sacloud/eventbus-api-go authors
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

package eventbus

import (
	"context"
	"errors"

	v1 "github.com/sacloud/eventbus-api-go/apis/v1"
)

type ScheduleAPI interface {
	List(ctx context.Context) ([]v1.CommonServiceItem, error)
	Read(ctx context.Context, id string) (*v1.CommonServiceItem, error)
	Create(ctx context.Context, request v1.CreateCommonServiceItemRequest) (*v1.CommonServiceItem, error)
	Update(ctx context.Context, id string, request v1.UpdateCommonServiceItemRequest) (*v1.CommonServiceItem, error)
	Delete(ctx context.Context, id string) error
}

var _ ScheduleAPI = (*scheduleOp)(nil)

type scheduleOp struct {
	client *v1.Client
}

func NewScheduleOp(client *v1.Client) ScheduleAPI {
	return &scheduleOp{client: client}
}

func (op *scheduleOp) List(ctx context.Context) ([]v1.CommonServiceItem, error) {
	ctx = setFilterProviderClass(ctx, v1.ProviderClassEventbusschedule)
	res, err := op.client.GetCommonServiceItems(ctx)
	if err != nil {
		return nil, NewAPIError("Schedule.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetCommonServiceItemsOK:
		return p.CommonServiceItems, nil
	case *v1.GetCommonServiceItemsUnauthorized:
		return nil, NewAPIError("Schedule.List", 401, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemsBadRequest:
		return nil, NewAPIError("Schedule.List", 400, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemsInternalServerError:
		return nil, NewAPIError("Schedule.List", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Schedule.List", 0, nil)
	}
}

func (op *scheduleOp) Read(ctx context.Context, id string) (*v1.CommonServiceItem, error) {
	res, err := op.client.GetCommonServiceItem(ctx, v1.GetCommonServiceItemParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Schedule.Read", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetCommonServiceItemOK:
		return &p.CommonServiceItem, nil
	case *v1.GetCommonServiceItemUnauthorized:
		return nil, NewAPIError("Schedule.Read", 401, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemBadRequest:
		return nil, NewAPIError("Schedule.Read", 400, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemNotFound:
		return nil, NewAPIError("Schedule.Read", 404, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemInternalServerError:
		return nil, NewAPIError("Schedule.Read", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Schedule.Read", 0, nil)
	}
}

func (op *scheduleOp) Create(ctx context.Context, request v1.CreateCommonServiceItemRequest) (*v1.CommonServiceItem, error) {
	if !request.CommonServiceItem.Settings.IsScheduleSettings() {
		return nil, NewError("invalid settings as Schedule", nil)
	}
	request.CommonServiceItem.Provider = v1.Provider{Class: v1.ProviderClassEventbusschedule}
	res, err := op.client.CreateCommonServiceItem(ctx, &request)
	if err != nil {
		return nil, NewAPIError("Schedule.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.CreateCommonServiceItemCreated:
		return &p.CommonServiceItem, nil
	case *v1.CreateCommonServiceItemBadRequest:
		return nil, NewAPIError("Schedule.Create", 400, errors.New(p.ErrorMsg.Value))
	case *v1.CreateCommonServiceItemUnauthorized:
		return nil, NewAPIError("Schedule.Create", 401, errors.New(p.ErrorMsg.Value))
	case *v1.CreateCommonServiceItemConflict:
		return nil, NewAPIError("Schedule.Create", 409, errors.New(p.ErrorMsg.Value))
	case *v1.CreateCommonServiceItemInternalServerError:
		return nil, NewAPIError("Schedule.Create", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Schedule.Create", 0, nil)
	}
}

func (op *scheduleOp) Update(ctx context.Context, id string, request v1.UpdateCommonServiceItemRequest) (*v1.CommonServiceItem, error) {
	if settings := request.CommonServiceItem.Settings; settings.IsSet() && !settings.Value.IsScheduleSettings() {
		return nil, NewError("invalid settings as Schedule", nil)
	}
	request.CommonServiceItem.Provider = v1.NewOptProvider(v1.Provider{Class: v1.ProviderClassEventbusschedule})
	res, err := op.client.UpdateCommonServiceItem(ctx, &request, v1.UpdateCommonServiceItemParams{ID: id})
	if err != nil {
		return nil, NewAPIError("Schedule.Update", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpdateCommonServiceItemOK:
		return &p.CommonServiceItem, nil
	case *v1.UpdateCommonServiceItemBadRequest:
		return nil, NewAPIError("Schedule.Update", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UpdateCommonServiceItemUnauthorized:
		return nil, NewAPIError("Schedule.Update", 401, errors.New(p.ErrorMsg.Value))
	case *v1.UpdateCommonServiceItemNotFound:
		return nil, NewAPIError("Schedule.Update", 404, errors.New(p.ErrorMsg.Value))
	case *v1.UpdateCommonServiceItemInternalServerError:
		return nil, NewAPIError("Schedule.Update", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("Schedule.Update", 0, nil)
	}
}

func (op *scheduleOp) Delete(ctx context.Context, id string) error {
	res, err := op.client.DeleteCommonServiceItem(ctx, v1.DeleteCommonServiceItemParams{ID: id})
	if err != nil {
		return NewAPIError("Schedule.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.DeleteCommonServiceItemOK:
		return nil
	case *v1.DeleteCommonServiceItemUnauthorized:
		return NewAPIError("Schedule.Delete", 401, errors.New(p.ErrorMsg.Value))
	case *v1.DeleteCommonServiceItemBadRequest:
		return NewAPIError("Schedule.Delete", 400, errors.New(p.ErrorMsg.Value))
	case *v1.DeleteCommonServiceItemNotFound:
		return NewAPIError("Schedule.Delete", 404, errors.New(p.ErrorMsg.Value))
	case *v1.DeleteCommonServiceItemInternalServerError:
		return NewAPIError("Schedule.Delete", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("Schedule.Delete", 0, nil)
	}
}
