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

package autoscale

import (
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	ID types.ID `service:"-" validate:"required"`

	Name        *string     `service:",omitempty" validate:"omitempty,min=1"`
	Description *string     `service:",omitempty" validate:"omitempty,max=512"`
	Tags        *types.Tags `service:",omitempty"`
	IconID      *types.ID   `service:",omitempty"`

	Zones                  *[]string                     `service:",omitempty" validate:"omitempty,required"`
	Config                 *string                       `service:",omitempty" validate:"omitempty,required"`
	TriggerType            *string                       `service:",omitempty" validate:"omitempty,required,oneof=cpu router"`
	CPUThresholdScaling    *UpdateCPUThresholdScaling    `service:"-" validate:"omitempty"`
	RouterThresholdScaling *UpdateRouterThresholdScaling `service:"-" validate:"omitempty"`
	ScheduleScaling        *[]*UpdateScheduleScaling     `service:"-" validate:"omitempty,dive"`
	Disabled               *bool                         `service:",omitempty"`

	SettingsHash string
}

type UpdateCPUThresholdScaling struct {
	ServerPrefix *string `service:",omitempty" validate:"omitempty,required"`
	Up           *int    `service:",omitempty" validate:"omitempty,required"`
	Down         *int    `service:",omitempty" validate:"omitempty,required"`
}

type UpdateRouterThresholdScaling struct {
	RouterPrefix *string `service:",omitempty" validate:"omitempty,required"`
	Direction    *string `service:",omitempty" validate:"omitempty,required,oneof=in out"`
	Mbps         *int    `service:",omitempty" validate:"omitempty,required"`
}

type UpdateScheduleScaling struct {
	Action    types.EAutoScaleAction `validate:"required,oneof=up down"`
	Hour      int                    `validate:"gte=0,lte=23"`
	Minute    int                    `validate:"oneof=0 15 30 45"`
	DayOfWeek []types.EDayOfTheWeek  `validate:"gt=0,unique,dive,required,oneof=sun mon tue wed thu fri sat"`
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) ToRequestParameter(current *iaas.AutoScale) (*iaas.AutoScaleUpdateRequest, error) {
	r := &iaas.AutoScaleUpdateRequest{}
	if err := serviceutil.RequestConvertTo(current, r); err != nil {
		return nil, err
	}
	if err := serviceutil.RequestConvertTo(req, r); err != nil {
		return nil, err
	}
	if req.CPUThresholdScaling != nil {
		if r.CPUThresholdScaling == nil {
			r.CPUThresholdScaling = &iaas.AutoScaleCPUThresholdScaling{}
		}
		if req.CPUThresholdScaling.ServerPrefix != nil {
			r.CPUThresholdScaling.ServerPrefix = *req.CPUThresholdScaling.ServerPrefix
		}
		if req.CPUThresholdScaling.Up != nil {
			r.CPUThresholdScaling.Up = *req.CPUThresholdScaling.Up
		}
		if req.CPUThresholdScaling.Down != nil {
			r.CPUThresholdScaling.Down = *req.CPUThresholdScaling.Down
		}
	}
	if req.RouterThresholdScaling != nil {
		if r.RouterThresholdScaling == nil {
			r.RouterThresholdScaling = &iaas.AutoScaleRouterThresholdScaling{}
		}

		if req.RouterThresholdScaling.RouterPrefix != nil {
			r.RouterThresholdScaling.RouterPrefix = *req.RouterThresholdScaling.RouterPrefix
		}
		if req.RouterThresholdScaling.Direction != nil {
			r.RouterThresholdScaling.Direction = *req.RouterThresholdScaling.Direction
		}
		if req.RouterThresholdScaling.Mbps != nil {
			r.RouterThresholdScaling.Mbps = *req.RouterThresholdScaling.Mbps
		}
	}
	if req.ScheduleScaling != nil {
		r.ScheduleScaling = []*iaas.AutoScaleScheduleScaling{}
		for _, ss := range *req.ScheduleScaling {
			r.ScheduleScaling = append(r.ScheduleScaling, &iaas.AutoScaleScheduleScaling{
				Action:    ss.Action,
				Hour:      ss.Hour,
				Minute:    ss.Minute,
				DayOfWeek: ss.DayOfWeek,
			})
		}
	}

	return r, nil
}
