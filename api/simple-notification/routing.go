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

type RoutingAPI interface {
	List(ctx context.Context) (*v1.ListCommonServiceItemsResponse, error)
	Create(ctx context.Context, request v1.PostCommonServiceItemRequest) (*v1.CreateCommonServiceItemCreated, error)
	Read(ctx context.Context, id string) (*v1.ReadCommonServiceItemOK, error)
	Update(ctx context.Context, id string, request v1.PutCommonServiceItemRequest) (*v1.UpdateCommonServiceItemOK, error)
	Delete(ctx context.Context, id string) error
	Reorder(ctx context.Context, request v1.ReorderNotificationRoutingRequest) (*v1.ReorderNotificationRoutingAccepted, error)
	ListSource(ctx context.Context) (*v1.ListNotificationSourcesResponse, error)
}

var _ RoutingAPI = (*RoutingOp)(nil)

type RoutingOp struct {
	client *v1.Client
}

func NewRoutingOp(client *v1.Client) RoutingAPI {
	return &RoutingOp{client: client}
}

func (o *RoutingOp) List(ctx context.Context) (*v1.ListCommonServiceItemsResponse, error) {
	ctx = setContextProviderClass(ctx, v1.CommonServiceItemProviderClassSaknoticerouting)

	return errorFromDecodedResponse("Routing.List", func() (*v1.ListCommonServiceItemsResponse, error) { return o.client.ListCommonServiceItems(ctx) })
}

func (o *RoutingOp) Create(ctx context.Context, request v1.PostCommonServiceItemRequest) (*v1.CreateCommonServiceItemCreated, error) {
	if request.CommonServiceItem.Tags == nil {
		request.CommonServiceItem.Tags = make([]string, 0)
	}
	request.CommonServiceItem.ServiceClass = v1.PostCommonServiceItemRequestCommonServiceItemServiceClassCloudSaknoticerouting2
	request.CommonServiceItem.Provider.Class = v1.PostCommonServiceItemRequestCommonServiceItemProviderClassSaknoticerouting
	request.CommonServiceItem.Settings.Type = v1.CommonServiceItemRoutingSettingsPostCommonServiceItemRequestCommonServiceItemSettings

	return errorFromDecodedResponse("Routing.Create", func() (*v1.CreateCommonServiceItemCreated, error) {
		return o.client.CreateCommonServiceItem(ctx, v1.NewOptPostCommonServiceItemRequest(request))
	})
}

func (o *RoutingOp) Read(ctx context.Context, id string) (*v1.ReadCommonServiceItemOK, error) {
	return errorFromDecodedResponse("Routing.Read", func() (*v1.ReadCommonServiceItemOK, error) {
		return o.client.ReadCommonServiceItem(ctx, v1.ReadCommonServiceItemParams{ID: id})
	})
}

func (o *RoutingOp) Update(ctx context.Context, id string, request v1.PutCommonServiceItemRequest) (*v1.UpdateCommonServiceItemOK, error) {
	if request.CommonServiceItem.Tags == nil {
		request.CommonServiceItem.Tags = []string{}
	}
	request.CommonServiceItem.Settings.Value.Type = v1.CommonServiceItemRoutingSettingsPutCommonServiceItemRequestCommonServiceItemSettings

	return errorFromDecodedResponse("Routing.Update", func() (*v1.UpdateCommonServiceItemOK, error) {
		return o.client.UpdateCommonServiceItem(ctx, v1.NewOptPutCommonServiceItemRequest(request), v1.UpdateCommonServiceItemParams{ID: id})
	})
}

func (o *RoutingOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE("Routing.Delete", func() (*v1.DeleteCommonServiceItemOK, error) {
		return o.client.DeleteCommonServiceItem(ctx, v1.DeleteCommonServiceItemParams{ID: id})
	})
}

func (o *RoutingOp) Reorder(ctx context.Context, request v1.ReorderNotificationRoutingRequest) (*v1.ReorderNotificationRoutingAccepted, error) {
	return errorFromDecodedResponse("Routing.Reorder", func() (*v1.ReorderNotificationRoutingAccepted, error) {
		return o.client.ReorderNotificationRouting(ctx, v1.NewOptReorderNotificationRoutingRequest(request))
	})
}

func (o *RoutingOp) ListSource(ctx context.Context) (*v1.ListNotificationSourcesResponse, error) {
	return errorFromDecodedResponse("Routing.ListSource", func() (*v1.ListNotificationSourcesResponse, error) { return o.client.ListNotificationSources(ctx) })
}
