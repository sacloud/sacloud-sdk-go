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
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/nfs/builder"
	"github.com/sacloud/packages-go/validate"
)

type ApplyRequest struct {
	ID   types.ID // for update
	Zone string   `service:"-" validate:"required"`

	Name           string `validate:"required"`
	Description    string `validate:"min=0,max=512"`
	Tags           types.Tags
	IconID         types.ID
	SwitchID       types.ID       `validate:"required"`
	Plan           types.ID       `validate:"required,oneof=1 2"` // types.NFSPlans.HDD or types.NFSPlans.SSD
	Size           types.ENFSSize `validate:"required"`           // types.NFSPlans.HDD or types.NFSPlans.SSD
	IPAddresses    []string       `validate:"required,min=1,max=2,dive,ipv4"`
	NetworkMaskLen int            `validate:"required"`
	DefaultRoute   string         `validate:"omitempty,ipv4"`

	NoWait bool
}

func (req *ApplyRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *ApplyRequest) Builder(caller iaas.APICaller) *builder.Builder {
	return &builder.Builder{
		ID:             req.ID,
		Zone:           req.Zone,
		Name:           req.Name,
		Description:    req.Description,
		Tags:           req.Tags,
		IconID:         req.IconID,
		SwitchID:       req.SwitchID,
		Plan:           req.Plan,
		Size:           req.Size,
		IPAddresses:    req.IPAddresses,
		NetworkMaskLen: req.NetworkMaskLen,
		DefaultRoute:   req.DefaultRoute,
		Caller:         caller,
		NoWait:         req.NoWait,
	}
}
