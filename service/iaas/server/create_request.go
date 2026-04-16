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

package server

import (
	"github.com/sacloud/iaas-api-go/types"
	diskService "github.com/sacloud/iaas-service-go/disk"
	"github.com/sacloud/packages-go/validate"
)

type CreateRequest struct {
	Zone string `service:"-" validate:"required"`

	Name            string `validate:"required"`
	Description     string `validate:"min=0,max=512"`
	Tags            types.Tags
	IconID          types.ID
	CPU             int
	MemoryGB        int
	GPU             int
	GPUModel        string
	CPUModel        string
	Commitment      types.ECommitment
	Generation      types.EPlanGeneration
	ConfidentialVM  bool
	InterfaceDriver types.EInterfaceDriver

	BootAfterCreate bool
	CDROMID         types.ID
	PrivateHostID   types.ID

	NetworkInterfaces []*NetworkInterface
	Disks             []*diskService.ApplyRequest
	NoWait            bool
}

func (req *CreateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *CreateRequest) ApplyRequest() *ApplyRequest {
	return &ApplyRequest{
		Zone:              req.Zone,
		Name:              req.Name,
		Description:       req.Description,
		Tags:              req.Tags,
		IconID:            req.IconID,
		CPU:               req.CPU,
		MemoryGB:          req.MemoryGB,
		GPU:               req.GPU,
		GPUModel:          req.GPUModel,
		CPUModel:          req.CPUModel,
		Commitment:        req.Commitment,
		Generation:        req.Generation,
		ConfidentialVM:    req.ConfidentialVM,
		InterfaceDriver:   req.InterfaceDriver,
		BootAfterCreate:   req.BootAfterCreate,
		CDROMID:           req.CDROMID,
		PrivateHostID:     req.PrivateHostID,
		NetworkInterfaces: req.NetworkInterfaces,
		Disks:             req.Disks,
		NoWait:            req.NoWait,
	}
}
