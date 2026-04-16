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

package nfs

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	Zone string   `validate:"required"`
	ID   types.ID `validate:"required"`

	Name        *string     `service:",omitempty" validate:"omitempty,min=1"`
	Description *string     `service:",omitempty" validate:"omitempty,max=512"`
	Tags        *types.Tags `service:",omitempty"`
	IconID      *types.ID   `service:",omitempty"`
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) ApplyRequest(ctx context.Context, caller iaas.APICaller) (*ApplyRequest, error) {
	client := iaas.NewNFSOp(caller)
	current, err := client.Read(ctx, req.Zone, req.ID)
	if err != nil {
		return nil, err
	}

	// find plan
	plan, err := query.GetNFSPlanInfo(ctx, iaas.NewNoteOp(caller), current.PlanID)
	if err != nil {
		return nil, err
	}

	applyRequest := &ApplyRequest{
		ID:             current.ID,
		Zone:           req.Zone,
		Name:           current.Name,
		Description:    current.Description,
		Tags:           current.Tags,
		IconID:         current.IconID,
		SwitchID:       current.SwitchID,
		Plan:           plan.DiskPlanID,
		Size:           plan.Size,
		IPAddresses:    current.IPAddresses,
		NetworkMaskLen: current.NetworkMaskLen,
		DefaultRoute:   current.DefaultRoute,
	}
	if err := serviceutil.RequestConvertTo(req, applyRequest); err != nil {
		return nil, err
	}
	return applyRequest, nil
}
