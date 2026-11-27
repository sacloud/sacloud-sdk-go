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

	v1 "github.com/sacloud/sacloud-sdk-go/api/simple-notification/apis/v1"
)

type GroupAPI interface {
	List(ctx context.Context) (*v1.ListCommonServiceItemsResponse, error)
	Create(ctx context.Context, request v1.PostCommonServiceItemRequest) (*v1.CreateCommonServiceItemCreated, error)
	Read(ctx context.Context, id string) (*v1.ReadCommonServiceItemOK, error)
	Update(ctx context.Context, id string, request v1.PutCommonServiceItemRequest) (*v1.UpdateCommonServiceItemOK, error)
	Delete(ctx context.Context, id string) error
	SendMessage(ctx context.Context, id string,
		request v1.SendNotificationRequest) (*v1.SendNotificationResponse, error)
}

var _ GroupAPI = (*GroupOp)(nil)

type GroupOp struct {
	client *v1.Client
}

func NewGroupOp(client *v1.Client) GroupAPI {
	return &GroupOp{client: client}
}

func (o *GroupOp) List(ctx context.Context) (*v1.ListCommonServiceItemsResponse, error) {
	ctx = setContextProviderClass(ctx, v1.CommonServiceItemProviderClassSaknoticegroup)
	return errorFromDecodedResponse("Group.List", func() (*v1.ListCommonServiceItemsResponse, error) { return o.client.ListCommonServiceItems(ctx) })
}

func (o *GroupOp) Create(ctx context.Context, request v1.PostCommonServiceItemRequest) (*v1.CreateCommonServiceItemCreated, error) {
	if request.CommonServiceItem.Tags == nil {
		request.CommonServiceItem.Tags = make([]string, 0)
	}
	request.CommonServiceItem.ServiceClass = v1.PostCommonServiceItemRequestCommonServiceItemServiceClassCloudSaknoticegroup2
	request.CommonServiceItem.Provider.Class = v1.PostCommonServiceItemRequestCommonServiceItemProviderClassSaknoticegroup
	request.CommonServiceItem.Provider.ServiceClass = v1.NewOptString("cloud/saknotice")
	request.CommonServiceItem.Settings.Type = v1.CommonServiceItemGroupSettingsPostCommonServiceItemRequestCommonServiceItemSettings

	return errorFromDecodedResponse("Group.Create", func() (*v1.CreateCommonServiceItemCreated, error) {
		return o.client.CreateCommonServiceItem(ctx, v1.NewOptPostCommonServiceItemRequest(request))
	})
}

func (o *GroupOp) Read(ctx context.Context, id string) (*v1.ReadCommonServiceItemOK, error) {
	return errorFromDecodedResponse("Group.Read", func() (*v1.ReadCommonServiceItemOK, error) {
		return o.client.ReadCommonServiceItem(ctx, v1.ReadCommonServiceItemParams{ID: id})
	})
}

func (o *GroupOp) Update(ctx context.Context, id string, request v1.PutCommonServiceItemRequest) (*v1.UpdateCommonServiceItemOK, error) {
	if request.CommonServiceItem.Tags == nil {
		request.CommonServiceItem.Tags = []string{}
	}
	request.CommonServiceItem.Settings.Value.Type = v1.CommonServiceItemGroupSettingsPutCommonServiceItemRequestCommonServiceItemSettings
	return errorFromDecodedResponse("Group.Update", func() (*v1.UpdateCommonServiceItemOK, error) {
		return o.client.UpdateCommonServiceItem(ctx, v1.NewOptPutCommonServiceItemRequest(request), v1.UpdateCommonServiceItemParams{ID: id})
	})
}

func (o *GroupOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE("Group.Delete", func() (*v1.DeleteCommonServiceItemOK, error) {
		return o.client.DeleteCommonServiceItem(ctx, v1.DeleteCommonServiceItemParams{ID: id})
	})
}

func (o *GroupOp) SendMessage(ctx context.Context, id string, request v1.SendNotificationRequest) (*v1.SendNotificationResponse, error) {
	return errorFromDecodedResponse("Group.SendMessage", func() (*v1.SendNotificationResponse, error) {
		return o.client.SendNotification(ctx, v1.NewOptSendNotificationRequest(request), v1.SendNotificationParams{ID: id})
	})
}
