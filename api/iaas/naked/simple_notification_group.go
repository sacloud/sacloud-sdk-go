// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package naked

import (
	"time"

	"github.com/sacloud/iaas-api-go/types"
)

// SimpleNotificationGroup シンプル通知グループ
type SimpleNotificationGroup struct {
	ID           types.ID                         `json:",omitempty" yaml:"id,omitempty" structs:",omitempty"`
	Name         string                           `json:",omitempty" yaml:"name,omitempty" structs:",omitempty"`
	Description  string                           `yaml:"description"`
	Tags         types.Tags                       `yaml:"tags"`
	Icon         *Icon                            `json:",omitempty" yaml:"icon,omitempty" structs:",omitempty"`
	CreatedAt    *time.Time                       `json:",omitempty" yaml:"created_at,omitempty" structs:",omitempty"`
	ModifiedAt   *time.Time                       `json:",omitempty" yaml:"modified_at,omitempty" structs:",omitempty"`
	Availability types.EAvailability              `json:",omitempty" yaml:"availability,omitempty" structs:",omitempty"`
	Provider     *Provider                        `json:",omitempty" yaml:"provider,omitempty" structs:",omitempty"`
	Settings     *SimpleNotificationGroupSettings `json:",omitempty" yaml:"settings,omitempty" structs:",omitempty"`
	SettingsHash string                           `json:",omitempty" yaml:"settings_hash,omitempty" structs:",omitempty"`
	Status       *SimpleNotificationGroupStatus   `json:",omitempty" yaml:"status" structs:",omitempty"`
	ServiceClass string                           `json:",omitempty" yaml:"service_class,omitempty" structs:",omitempty"`
}

// SimpleNotificationGroupSettingsUpdate シンプル通知更新パラメータ
type SimpleNotificationGroupSettingsUpdate struct {
	Settings     *SimpleNotificationGroupSettings `json:",omitempty" yaml:"settings,omitempty" structs:",omitempty"`
	SettingsHash string                           `json:",omitempty" yaml:"settings_hash,omitempty" structs:",omitempty"`
}

// SimpleNotificationGroupSettings セッティング
type SimpleNotificationGroupSettings struct {
	Destinations []string
	Disabled     bool
	Sources      []string
}

// SimpleNotificationGroupStatus ステータス
type SimpleNotificationGroupStatus struct {
	Disabled     bool
	ErrorMessage string `json:",omitempty" yaml:",omitempty"`
	IsValid      bool
}

type SimpleNotificationMessageRequest struct {
	Message string
}

type SimpleNotificationHistory struct {
	RequestID  string    `json:"request_id" yaml:"request_id"`
	SourceID   string    `json:"source_id" yaml:"source_id"`
	ReceivedAt time.Time `json:"received_at" yaml:"received_at"`
	Message    *SimpleNotificationHistoryMessage
	Statuses   []*SimpleNotificationHistoryStatus
}

type SimpleNotificationHistoryMessage struct {
	Body      string `json:"body" yaml:"body"`
	Color     string `json:"color" yaml:"color"`
	ColorCode string `json:"color_code" yaml:"color_code"`
	IconURL   string `json:"icon_url" yaml:"icon_url"`
	ImageURL  string `json:"image_url" yaml:"image_url"`
	Title     string `json:"title" yaml:"title"`
}

type SimpleNotificationHistoryStatus struct {
	ID                    string    `json:"id" yaml:"id"`
	Status                int       `json:"status" yaml:"status"`
	ErrorInfo             string    `json:"error_info" yaml:"error_info"`
	NotificationRequestID string    `json:"notification_request_id" yaml:"notification_request_id"`
	GroupID               string    `json:"group_id" yaml:"group_id"`
	CreatedAt             time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" yaml:"updated_at"`
}

type SimpleNotificationHistories struct {
	NotificationHistories []*SimpleNotificationHistory `json:",omitempty"`
}
