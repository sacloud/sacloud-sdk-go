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

type CreateRequest struct {
	Target             string `validate:"required"`
	Description        string `validate:"min=0,max=512"`
	Tags               types.Tags
	IconID             types.ID
	MaxCheckAttempts   int                            `mapconv:"Settings.SimpleMonitor.MaxCheckAttempts" validate:"min=1,max=10"`
	RetryInterval      int                            `mapconv:"Settings.SimpleMonitor.RetryInterval" validate:"min=10,max=3600"`
	DelayLoop          int                            `mapconv:"Settings.SimpleMonitor.DelayLoop" validate:"min=60,max=3600"`
	Enabled            types.StringFlag               `mapconv:"Settings.SimpleMonitor.Enabled"`
	HealthCheck        *iaas.SimpleMonitorHealthCheck `mapconv:"Settings.SimpleMonitor.HealthCheck,recursive"`
	NotifyEmailEnabled types.StringFlag               `mapconv:"Settings.SimpleMonitor.NotifyEmail.Enabled"`
	NotifyEmailHTML    types.StringFlag               `mapconv:"Settings.SimpleMonitor.NotifyEmail.HTML"`
	NotifySlackEnabled types.StringFlag               `mapconv:"Settings.SimpleMonitor.NotifySlack.Enabled"`
	SlackWebhooksURL   string                         `mapconv:"Settings.SimpleMonitor.NotifySlack.IncomingWebhooksURL"`
	NotifyInterval     int                            `mapconv:"Settings.SimpleMonitor.NotifyInterval" validate:"min=3600,max=259200"`
	Timeout            int
	MonitoringSuiteLog *iaas.MonitoringSuiteLog
}

func (req *CreateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *CreateRequest) ToRequestParameter() (*iaas.SimpleMonitorCreateRequest, error) {
	params := &iaas.SimpleMonitorCreateRequest{}
	if err := serviceutil.RequestConvertTo(req, params); err != nil {
		return nil, err
	}
	return params, nil
}
