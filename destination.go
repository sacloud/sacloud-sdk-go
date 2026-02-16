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

type DestinationAPI interface {
	List(ctx context.Context) (*v1.ListCommonServiceItemsResponse, error)
	Create(ctx context.Context, request v1.PostCommonServiceItemRequest) (*v1.CreateCommonServiceItemCreated, error)
	Read(ctx context.Context, id string) (*v1.GetCommonServiceItemOK, error)
	Update(ctx context.Context, id string, request v1.PutCommonServiceItemRequest) (*v1.UpdateCommonServiceItemOK, error)
	Delete(ctx context.Context, id string) error
	ReadStatus(ctx context.Context, id string) (*v1.GetCommonServiceItemStatusResponse, error)
}

var _ DestinationAPI = (*DestinationOp)(nil)

type DestinationOp struct {
	client *v1.Client
}

func NewDestinationOp(client *v1.Client) DestinationAPI {
	return &DestinationOp{client: client}
}

func (o *DestinationOp) List(ctx context.Context) (*v1.ListCommonServiceItemsResponse, error) {
	const methodName = "Destination.List"
	ctx = setContextProviderClass(ctx, v1.CommonServiceItemProviderClassSaknoticedestination)
	res, err := o.client.ListCommonServiceItems(ctx)
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, NewError(methodName, err)
	}
	return res, nil
}

func (o *DestinationOp) Create(ctx context.Context, request v1.PostCommonServiceItemRequest) (*v1.CreateCommonServiceItemCreated, error) {
	const methodName = "Destination.Create"
	request.CommonServiceItem.ServiceClass = v1.PostCommonServiceItemRequestCommonServiceItemServiceClassCloudSaknoticedestination2
	request.CommonServiceItem.Provider.Class = v1.PostCommonServiceItemRequestCommonServiceItemProviderClassSaknoticedestination
	request.CommonServiceItem.Settings.Type = v1.CommonServiceItemDestinationSettingsPostCommonServiceItemRequestCommonServiceItemSettings
	res, err := o.client.CreateCommonServiceItem(ctx, v1.NewOptPostCommonServiceItemRequest(request))
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, NewError(methodName, err)
	}
	return res, nil
}

func (o *DestinationOp) Read(ctx context.Context, id string) (*v1.GetCommonServiceItemOK, error) {
	const methodName = "Destination.Read"
	res, err := o.client.GetCommonServiceItem(ctx, v1.GetCommonServiceItemParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, NewError(methodName, err)
	}
	return res, nil
}

func (o *DestinationOp) Update(ctx context.Context, id string, request v1.PutCommonServiceItemRequest) (*v1.UpdateCommonServiceItemOK, error) {
	const methodName = "Destination.Update"
	request.CommonServiceItem.Settings.Value.Type = v1.CommonServiceItemDestinationSettingsPutCommonServiceItemRequestCommonServiceItemSettings
	res, err := o.client.UpdateCommonServiceItem(ctx, v1.NewOptPutCommonServiceItemRequest(request), v1.UpdateCommonServiceItemParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, NewError(methodName, err)
	}
	return res, nil
}

func (o *DestinationOp) Delete(ctx context.Context, id string) error {
	const methodName = "Destination.Delete"
	_, err := o.client.DeleteCommonServiceItem(ctx, v1.DeleteCommonServiceItemParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return NewAPIError(methodName, e.StatusCode, err)
		}
		return NewError(methodName, err)
	}
	return nil
}

func (o *DestinationOp) ReadStatus(ctx context.Context, id string) (*v1.GetCommonServiceItemStatusResponse, error) {
	const methodName = "Destination.ReadStatus"
	res, err := o.client.GetCommonServiceItemStatus(ctx, v1.GetCommonServiceItemStatusParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, NewError(methodName, err)
	}
	return res, nil
}
