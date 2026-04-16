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

type GroupAPI interface {
	List(ctx context.Context) (*v1.ListCommonServiceItemsResponse, error)
	Create(ctx context.Context, request v1.PostCommonServiceItemRequest) (*v1.CreateCommonServiceItemCreated, error)
	Read(ctx context.Context, id string) (*v1.GetCommonServiceItemOK, error)
	Update(ctx context.Context, id string, request v1.PutCommonServiceItemRequest) (*v1.UpdateCommonServiceItemOK, error)
	Delete(ctx context.Context, id string) error
	SendMessage(ctx context.Context, id string,
		request v1.SendNotificationMessageRequest) (*v1.SendNotificationMessageResponse, error)
}

var _ GroupAPI = (*GroupOp)(nil)

type GroupOp struct {
	client *v1.Client
}

func NewGroupOp(client *v1.Client) GroupAPI {
	return &GroupOp{client: client}
}

func (o *GroupOp) List(ctx context.Context) (*v1.ListCommonServiceItemsResponse, error) {
	const methodName = "Group.List"

	ctx = setContextProviderClass(ctx, v1.CommonServiceItemProviderClassSaknoticegroup)
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

func (o *GroupOp) Create(ctx context.Context, request v1.PostCommonServiceItemRequest) (*v1.CreateCommonServiceItemCreated, error) {
	const methodName = "Group.Create"
	request.CommonServiceItem.ServiceClass = v1.CommonServiceItemServiceClassCloudSaknoticegroup2
	request.CommonServiceItem.Provider.Class = v1.CommonServiceItemProviderClassSaknoticegroup
	request.CommonServiceItem.Provider.ServiceClass = v1.OptString{Value: "cloud/saknotice", Set: true}
	request.CommonServiceItem.Settings.Type = v1.GroupSettingsCommonServiceItemSettings

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

func (o *GroupOp) Read(ctx context.Context, id string) (*v1.GetCommonServiceItemOK, error) {
	const methodName = "Group.Read"
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

func (o *GroupOp) Update(ctx context.Context, id string, request v1.PutCommonServiceItemRequest) (*v1.UpdateCommonServiceItemOK, error) {
	const methodName = "Group.Update"
	request.CommonServiceItem.Settings.Value.Type = v1.GroupSettingsCommonServiceItemSettings
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

func (o *GroupOp) Delete(ctx context.Context, id string) error {
	const methodName = "Group.Delete"
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

func (o *GroupOp) SendMessage(ctx context.Context, id string, request v1.SendNotificationMessageRequest) (*v1.SendNotificationMessageResponse, error) {
	const methodName = "Group.SendMessage"
	res, err := o.client.SendNotificationMessage(ctx, v1.NewOptSendNotificationMessageRequest(request), v1.SendNotificationMessageParams{ID: id})
	if err != nil {
		var e *v1.ErrorStatusCode
		if errors.As(err, &e) {
			return nil, NewAPIError(methodName, e.StatusCode, err)
		}
		return nil, NewError(methodName, err)
	}
	return res, nil
}
