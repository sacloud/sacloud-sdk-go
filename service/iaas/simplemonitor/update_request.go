// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package simplemonitor

import (
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	ID types.ID `service:"-" validate:"required"`

	Description        *string                        `service:",omitempty" validate:"omitempty,max=512"`
	Tags               *types.Tags                    `service:",omitempty"`
	IconID             *types.ID                      `service:",omitempty"`
	MaxCheckAttempts   *int                           `service:",omitempty"`
	RetryInterval      *int                           `service:",omitempty"`
	DelayLoop          *int                           `service:",omitempty"`
	Enabled            *types.StringFlag              `service:",omitempty"`
	HealthCheck        *iaas.SimpleMonitorHealthCheck `service:",omitempty"`
	NotifyEmailEnabled *types.StringFlag              `service:",omitempty"`
	NotifyEmailHTML    *types.StringFlag              `service:",omitempty"`
	NotifySlackEnabled *types.StringFlag              `service:",omitempty"`
	SlackWebhooksURL   *string                        `service:",omitempty"`
	NotifyInterval     *int                           `service:",omitempty"`
	Timeout            *int                           `service:",omitempty"`
	MonitoringSuiteLog *iaas.MonitoringSuiteLog       `service:",omitempty"`
	SettingsHash       string
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) ToRequestParameter(current *iaas.SimpleMonitor) (*iaas.SimpleMonitorUpdateRequest, error) {
	r := &iaas.SimpleMonitorUpdateRequest{}
	if err := serviceutil.RequestConvertTo(current, r); err != nil {
		return nil, err
	}
	if err := serviceutil.RequestConvertTo(req, r); err != nil {
		return nil, err
	}
	return r, nil
}
