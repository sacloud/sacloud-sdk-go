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

package loadbalancer

import (
	"context"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	Zone string   `service:"-" validate:"required"`
	ID   types.ID `service:"-" validate:"required"`

	Name               *string                              `service:",omitempty" validate:"omitempty,min=1"`
	Description        *string                              `service:",omitempty" validate:"omitempty,max=512"`
	Tags               *types.Tags                          `service:",omitempty"`
	IconID             *types.ID                            `service:",omitempty"`
	VirtualIPAddresses *iaas.LoadBalancerVirtualIPAddresses `service:",omitempty"`

	SettingsHash string
	NoWait       bool
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) ApplyRequest(ctx context.Context, caller iaas.APICaller) (*ApplyRequest, error) {
	client := iaas.NewLoadBalancerOp(caller)
	current, err := client.Read(ctx, req.Zone, req.ID)
	if err != nil {
		return nil, err
	}
	if current.Availability != types.Availabilities.Available {
		return nil, fmt.Errorf("target has invalid Availability: Zone=%s ID=%s Availability=%v", req.Zone, req.ID.String(), current.Availability)
	}

	applyRequest := &ApplyRequest{
		ID:                 req.ID,
		Zone:               req.Zone,
		Name:               current.Name,
		Description:        current.Description,
		Tags:               current.Tags,
		IconID:             current.IconID,
		SwitchID:           current.SwitchID,
		PlanID:             current.PlanID,
		VRID:               current.VRID,
		IPAddresses:        current.IPAddresses,
		NetworkMaskLen:     current.NetworkMaskLen,
		DefaultRoute:       current.DefaultRoute,
		VirtualIPAddresses: current.VirtualIPAddresses,
	}
	if err := serviceutil.RequestConvertTo(req, applyRequest); err != nil {
		return nil, err
	}
	return applyRequest, nil
}
