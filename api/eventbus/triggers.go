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

	v1 "github.com/sacloud/sacloud-sdk-go/api/eventbus/apis/v1"
)

type TriggerAPI interface {
	List(ctx context.Context) ([]v1.CommonServiceItem, error)
	Read(ctx context.Context, id string) (*v1.CommonServiceItem, error)
	Create(ctx context.Context, request v1.CreateCommonServiceItemRequest) (*v1.CommonServiceItem, error)
	Update(ctx context.Context, id string, request v1.UpdateCommonServiceItemRequest) (*v1.CommonServiceItem, error)
	Delete(ctx context.Context, id string) error
}

var _ TriggerAPI = (*triggerOp)(nil)

type triggerOp struct {
	client *v1.Client
}

func NewTriggerOp(client *v1.Client) TriggerAPI {
	return &triggerOp{client: client}
}

func (op *triggerOp) List(ctx context.Context) ([]v1.CommonServiceItem, error) {
	ctx = setFilterProviderClass(ctx, v1.ProviderClassEventbustrigger)
	res, err := errorFromDecodedResponse[v1.ListCommonServiceItemsOK]("Trigger.List", func() (any, error) { return op.client.ListCommonServiceItems(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetCommonServiceItems(), nil
}

func (op *triggerOp) Read(ctx context.Context, id string) (*v1.CommonServiceItem, error) {
	res, err := errorFromDecodedResponse[v1.ReadCommonServiceItemOK]("Trigger.Read", func() (any, error) {
		return op.client.ReadCommonServiceItem(ctx, v1.ReadCommonServiceItemParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetCommonServiceItem()), nil
}

func (op *triggerOp) Create(ctx context.Context, request v1.CreateCommonServiceItemRequest) (*v1.CommonServiceItem, error) {
	if !request.CommonServiceItem.Settings.IsTriggerSettings() {
		return nil, NewError("invalid settings as Trigger", nil)
	}
	request.CommonServiceItem.Provider = v1.Provider{Class: v1.ProviderClassEventbustrigger}

	res, err := errorFromDecodedResponse[v1.CreateCommonServiceItemCreated]("Trigger.Create", func() (any, error) { return op.client.CreateCommonServiceItem(ctx, &request) })
	if err != nil {
		return nil, err
	}

	return new(res.GetCommonServiceItem()), nil
}

func (op *triggerOp) Update(ctx context.Context, id string, request v1.UpdateCommonServiceItemRequest) (*v1.CommonServiceItem, error) {
	if settings := request.CommonServiceItem.Settings; settings.IsSet() && !settings.Value.IsTriggerSettings() {
		return nil, NewError("invalid settings as Trigger", nil)
	}
	request.CommonServiceItem.Provider = v1.NewOptProvider(v1.Provider{Class: v1.ProviderClassEventbustrigger})

	res, err := errorFromDecodedResponse[v1.UpdateCommonServiceItemOK]("Trigger.Update", func() (any, error) {
		return op.client.UpdateCommonServiceItem(ctx, &request, v1.UpdateCommonServiceItemParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetCommonServiceItem()), nil
}

func (op *triggerOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[v1.DeleteCommonServiceItemOK]("Trigger.Delete", func() (any, error) {
		return op.client.DeleteCommonServiceItem(ctx, v1.DeleteCommonServiceItemParams{ID: id})
	})
}
