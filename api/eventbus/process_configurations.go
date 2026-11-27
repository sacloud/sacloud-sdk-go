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

	v1 "github.com/sacloud/sacloud-sdk-go/api/eventbus/apis/v1"
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
	res, err := errorFromDecodedResponse[v1.ListCommonServiceItemsOK]("ProcessConfiguration.List", func() (any, error) { return op.client.ListCommonServiceItems(ctx) })
	if err != nil {
		return nil, err
	}

	return res.GetCommonServiceItems(), nil
}

func (op *processConfigurationOp) Read(ctx context.Context, id string) (*v1.CommonServiceItem, error) {
	res, err := errorFromDecodedResponse[v1.ReadCommonServiceItemOK]("ProcessConfiguration.Read", func() (any, error) {
		return op.client.ReadCommonServiceItem(ctx, v1.ReadCommonServiceItemParams{ID: id})
	})
	if err != nil {
		return nil, err
	}
	return new(res.GetCommonServiceItem()), nil
}

func (op *processConfigurationOp) Create(ctx context.Context, request v1.CreateCommonServiceItemRequest) (*v1.CommonServiceItem, error) {
	if !request.CommonServiceItem.Settings.IsProcessConfigurationSettings() {
		return nil, NewError("invalid settings as ProcessConfiguration", nil)
	}
	request.CommonServiceItem.Provider = v1.Provider{Class: v1.ProviderClassEventbusprocessconfiguration}

	res, err := errorFromDecodedResponse[v1.CreateCommonServiceItemCreated]("ProcessConfiguration.Create", func() (any, error) { return op.client.CreateCommonServiceItem(ctx, &request) })
	if err != nil {
		return nil, err
	}

	return new(res.GetCommonServiceItem()), nil
}

func (op *processConfigurationOp) Update(ctx context.Context, id string, request v1.UpdateCommonServiceItemRequest) (*v1.CommonServiceItem, error) {
	if settings := request.CommonServiceItem.Settings; settings.IsSet() && !settings.Value.IsProcessConfigurationSettings() {
		return nil, NewError("invalid settings as ProcessConfiguration", nil)
	}
	request.CommonServiceItem.Provider = v1.NewOptProvider(v1.Provider{Class: v1.ProviderClassEventbusprocessconfiguration})
	res, err := errorFromDecodedResponse[v1.UpdateCommonServiceItemOK]("ProcessConfiguration.Update", func() (any, error) {
		return op.client.UpdateCommonServiceItem(ctx, &request, v1.UpdateCommonServiceItemParams{ID: id})
	})
	if err != nil {
		return nil, err
	}

	return new(res.GetCommonServiceItem()), nil
}

func (op *processConfigurationOp) UpdateSecret(ctx context.Context, id string, secret v1.SetSecretRequest) error {
	return errorFromDecodedResponseE[v1.SetProcessConfigurationSecretOK]("ProcessConfiguration.UpdateSecret", func() (any, error) {
		return op.client.SetProcessConfigurationSecret(ctx, &secret, v1.SetProcessConfigurationSecretParams{ID: id})
	})
}

func (op *processConfigurationOp) Delete(ctx context.Context, id string) error {
	return errorFromDecodedResponseE[v1.DeleteCommonServiceItemOK]("ProcessConfiguration.Delete", func() (any, error) {
		return op.client.DeleteCommonServiceItem(ctx, v1.DeleteCommonServiceItemParams{ID: id})
	})
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
