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
	"github.com/sacloud/packages-go/validate"
)

type CreateRequest struct {
	Name        string `validate:"required"`
	Description string `validate:"min=0,max=512"`
	Tags        types.Tags
	IconID      types.ID

	Zones                  []string                      `validate:"required"`
	Config                 string                        `validate:"required"`
	TriggerType            types.EAutoScaleTriggerType   `validate:"omitempty,oneof=cpu router schedule"`
	CPUThresholdScaling    *CreateCPUThresholdScaling    `validate:"omitempty"`
	RouterThresholdScaling *CreateRouterThresholdScaling `validate:"omitempty"`
	ScheduleScaling        []*CreateScheduleScaling      `validate:"omitempty,dive"`
	Disabled               bool

	APIKeyID string `validate:"required"`
}

type CreateCPUThresholdScaling struct {
	ServerPrefix string `validate:"required"`
	Up           int    `validate:"required"`
	Down         int    `validate:"required"`
}

type CreateRouterThresholdScaling struct {
	RouterPrefix string `validate:"required"`
	Direction    string `validate:"required,oneof=in out"`
	Mbps         int    `validate:"required"`
}

type CreateScheduleScaling struct {
	Action    types.EAutoScaleAction `validate:"required,oneof=up down"`
	Hour      int                    `validate:"gte=0,lte=23"`
	Minute    int                    `validate:"oneof=0 15 30 45"`
	DayOfWeek []types.EDayOfTheWeek  `validate:"gt=0,unique,dive,required,oneof=sun mon tue wed thu fri sat"`
}

func (req *CreateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *CreateRequest) ToRequestParameter() (*iaas.AutoScaleCreateRequest, error) {
	createReq := &iaas.AutoScaleCreateRequest{
		Name:        req.Name,
		Description: req.Description,
		Tags:        req.Tags,
		IconID:      req.IconID,
		Zones:       req.Zones,
		Config:      req.Config,
		TriggerType: req.TriggerType,
		Disabled:    req.Disabled,

		APIKeyID: req.APIKeyID,
	}
	if req.CPUThresholdScaling != nil {
		createReq.CPUThresholdScaling = &iaas.AutoScaleCPUThresholdScaling{
			ServerPrefix: req.CPUThresholdScaling.ServerPrefix,
			Up:           req.CPUThresholdScaling.Up,
			Down:         req.CPUThresholdScaling.Down,
		}
	}
	if req.RouterThresholdScaling != nil {
		createReq.RouterThresholdScaling = &iaas.AutoScaleRouterThresholdScaling{
			RouterPrefix: req.RouterThresholdScaling.RouterPrefix,
			Direction:    req.RouterThresholdScaling.Direction,
			Mbps:         req.RouterThresholdScaling.Mbps,
		}
	}
	for _, ss := range req.ScheduleScaling {
		createReq.ScheduleScaling = append(createReq.ScheduleScaling, &iaas.AutoScaleScheduleScaling{
			Action:    ss.Action,
			Hour:      ss.Hour,
			Minute:    ss.Minute,
			DayOfWeek: ss.DayOfWeek,
		})
	}

	return createReq, nil
}
