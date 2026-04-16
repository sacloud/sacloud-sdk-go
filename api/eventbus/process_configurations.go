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
	"encoding/json"
	"errors"

	v1 "github.com/sacloud/eventbus-api-go/apis/v1"
)

type ProcessConfigurationAPI interface {
	List(ctx context.Context) ([]v1.CommonServiceItem, error)
	Read(ctx context.Context, id string) (*v1.CommonServiceItem, error)
	Create(ctx context.Context, request v1.CreateCommonServiceItemRequest) (*v1.CommonServiceItem, error)
	Update(ctx context.Context, id string, request v1.UpdateCommonServiceItemRequest) (*v1.CommonServiceItem, error)
	UpdateSecret(ctx context.Context, id string, secret v1.SetSecretRequest) error
	Delete(ctx context.Context, id string) error
}

var _ ProcessConfigurationAPI = (*processConfigurationOp)(nil)

type processConfigurationOp struct {
	client *v1.Client
}

func NewProcessConfigurationOp(client *v1.Client) ProcessConfigurationAPI {
	return &processConfigurationOp{client: client}
}

func (op *processConfigurationOp) List(ctx context.Context) ([]v1.CommonServiceItem, error) {
	ctx = setFilterProviderClass(ctx, v1.ProviderClassEventbusprocessconfiguration)
	res, err := op.client.GetCommonServiceItems(ctx)
	if err != nil {
		return nil, NewAPIError("ProcessConfiguration.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetCommonServiceItemsOK:
		return p.CommonServiceItems, nil
	case *v1.GetCommonServiceItemsUnauthorized:
		return nil, NewAPIError("ProcessConfiguration.List", 401, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemsBadRequest:
		return nil, NewAPIError("ProcessConfiguration.List", 400, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemsInternalServerError:
		return nil, NewAPIError("ProcessConfiguration.List", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("ProcessConfiguration.List", 0, nil)
	}
}

func (op *processConfigurationOp) Read(ctx context.Context, id string) (*v1.CommonServiceItem, error) {
	res, err := op.client.GetCommonServiceItem(ctx, v1.GetCommonServiceItemParams{ID: id})
	if err != nil {
		return nil, NewAPIError("ProcessConfiguration.Read", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetCommonServiceItemOK:
		return &p.CommonServiceItem, nil
	case *v1.GetCommonServiceItemUnauthorized:
		return nil, NewAPIError("ProcessConfiguration.Read", 401, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemBadRequest:
		return nil, NewAPIError("ProcessConfiguration.Read", 400, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemNotFound:
		return nil, NewAPIError("ProcessConfiguration.Read", 404, errors.New(p.ErrorMsg.Value))
	case *v1.GetCommonServiceItemInternalServerError:
		return nil, NewAPIError("ProcessConfiguration.Read", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("CommonServiceItem.Read", 0, nil)
	}
}

func (op *processConfigurationOp) Create(ctx context.Context, request v1.CreateCommonServiceItemRequest) (*v1.CommonServiceItem, error) {
	if !request.CommonServiceItem.Settings.IsProcessConfigurationSettings() {
		return nil, NewError("invalid settings as ProcessConfiguration", nil)
	}
	request.CommonServiceItem.Provider = v1.Provider{Class: v1.ProviderClassEventbusprocessconfiguration}
	res, err := op.client.CreateCommonServiceItem(ctx, &request)
	if err != nil {
		return nil, NewAPIError("ProcessConfiguration.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.CreateCommonServiceItemCreated:
		return &p.CommonServiceItem, nil
	case *v1.CreateCommonServiceItemBadRequest:
		return nil, NewAPIError("ProcessConfiguration.Create", 400, errors.New(p.ErrorMsg.Value))
	case *v1.CreateCommonServiceItemUnauthorized:
		return nil, NewAPIError("ProcessConfiguration.Create", 401, errors.New(p.ErrorMsg.Value))
	case *v1.CreateCommonServiceItemConflict:
		return nil, NewAPIError("ProcessConfiguration.Create", 409, errors.New(p.ErrorMsg.Value))
	case *v1.CreateCommonServiceItemInternalServerError:
		return nil, NewAPIError("ProcessConfiguration.Create", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("ProcessConfiguration.Create", 0, nil)
	}
}

func (op *processConfigurationOp) Update(ctx context.Context, id string, request v1.UpdateCommonServiceItemRequest) (*v1.CommonServiceItem, error) {
	if settings := request.CommonServiceItem.Settings; settings.IsSet() && !settings.Value.IsProcessConfigurationSettings() {
		return nil, NewError("invalid settings as ProcessConfiguration", nil)
	}
	request.CommonServiceItem.Provider = v1.NewOptProvider(v1.Provider{Class: v1.ProviderClassEventbusprocessconfiguration})
	res, err := op.client.UpdateCommonServiceItem(ctx, &request, v1.UpdateCommonServiceItemParams{ID: id})
	if err != nil {
		return nil, NewAPIError("ProcessConfiguration.Update", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpdateCommonServiceItemOK:
		return &p.CommonServiceItem, nil
	case *v1.UpdateCommonServiceItemBadRequest:
		return nil, NewAPIError("ProcessConfiguration.Update", 400, errors.New(p.ErrorMsg.Value))
	case *v1.UpdateCommonServiceItemUnauthorized:
		return nil, NewAPIError("ProcessConfiguration.Update", 401, errors.New(p.ErrorMsg.Value))
	case *v1.UpdateCommonServiceItemNotFound:
		return nil, NewAPIError("ProcessConfiguration.Update", 404, errors.New(p.ErrorMsg.Value))
	case *v1.UpdateCommonServiceItemInternalServerError:
		return nil, NewAPIError("ProcessConfiguration.Update", 500, errors.New(p.ErrorMsg.Value))
	default:
		return nil, NewAPIError("ProcessConfiguration.Update", 0, nil)
	}
}

func (op *processConfigurationOp) UpdateSecret(ctx context.Context, id string, secret v1.SetSecretRequest) error {
	res, err := op.client.SetProcessConfigurationSecret(ctx, &secret, v1.SetProcessConfigurationSecretParams{ID: id})
	if err != nil {
		return NewAPIError("ProcessConfiguration.UpdateSecret", 0, err)
	}

	switch p := res.(type) {
	case *v1.SetProcessConfigurationSecretOK:
		return nil
	case *v1.SetProcessConfigurationSecretBadRequest:
		return NewAPIError("ProcessConfiguration.UpdateSecret", 400, errors.New(p.ErrorMsg.Value))
	case *v1.SetProcessConfigurationSecretUnauthorized:
		return NewAPIError("ProcessConfiguration.UpdateSecret", 401, errors.New(p.ErrorMsg.Value))
	case *v1.SetProcessConfigurationSecretNotFound:
		return NewAPIError("ProcessConfiguration.UpdateSecret", 404, errors.New(p.ErrorMsg.Value))
	case *v1.SetProcessConfigurationSecretInternalServerError:
		return NewAPIError("ProcessConfiguration.UpdateSecret", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("ProcessConfiguration.UpdateSecret", 0, nil)
	}
}

func (op *processConfigurationOp) Delete(ctx context.Context, id string) error {
	res, err := op.client.DeleteCommonServiceItem(ctx, v1.DeleteCommonServiceItemParams{ID: id})
	if err != nil {
		return NewAPIError("ProcessConfiguration.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.DeleteCommonServiceItemOK:
		return nil
	case *v1.DeleteCommonServiceItemUnauthorized:
		return NewAPIError("ProcessConfiguration.Delete", 401, errors.New(p.ErrorMsg.Value))
	case *v1.DeleteCommonServiceItemBadRequest:
		return NewAPIError("ProcessConfiguration.Delete", 400, errors.New(p.ErrorMsg.Value))
	case *v1.DeleteCommonServiceItemNotFound:
		return NewAPIError("ProcessConfiguration.Delete", 404, errors.New(p.ErrorMsg.Value))
	case *v1.DeleteCommonServiceItemInternalServerError:
		return NewAPIError("ProcessConfiguration.Delete", 500, errors.New(p.ErrorMsg.Value))
	default:
		return NewAPIError("ProcessConfiguration.Delete", 0, nil)
	}
}

func CreateSimpleNotificationSettings(groupId string, message string) v1.Settings {
	param, _ := json.Marshal(map[string]string{"group_id": groupId, "message": message})
	return v1.NewProcessConfigurationSettingsSettings(v1.ProcessConfigurationSettings{
		Destination: v1.ProcessConfigurationSettingsDestinationSimplenotification,
		Parameters:  string(param),
	})
}

func CreateSimpleMqSettings(queueName string, content string) v1.Settings {
	param, _ := json.Marshal(map[string]string{"queue_name": queueName, "content": content})
	return v1.NewProcessConfigurationSettingsSettings(v1.ProcessConfigurationSettings{
		Destination: v1.ProcessConfigurationSettingsDestinationSimplemq,
		Parameters:  string(param),
	})
}

type AutoScaleAction string

const (
	AutoScaleActionUp   AutoScaleAction = "scale_up"
	AutoScaleActionDown AutoScaleAction = "scale_down"
)

func CreateAutoScaleSettings(action AutoScaleAction, resourceID string) v1.Settings {
	param, _ := json.Marshal(map[string]string{"action": string(action), "resource_id": resourceID})
	return v1.NewProcessConfigurationSettingsSettings(v1.ProcessConfigurationSettings{
		Destination: v1.ProcessConfigurationSettingsDestinationAutoscale,
		Parameters:  string(param),
	})
}
