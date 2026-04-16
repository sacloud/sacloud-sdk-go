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
	"errors"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	diskService "github.com/sacloud/iaas-service-go/disk"
	diskBuilder "github.com/sacloud/iaas-service-go/disk/builder"
	server "github.com/sacloud/iaas-service-go/server/builder"
	"github.com/sacloud/packages-go/validate"
)

type ApplyRequest struct {
	Zone string `validate:"required"`
	ID   types.ID

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

	ForceShutdown bool
}

func (req *ApplyRequest) Validate() error {
	if err := validate.New().Struct(req); err != nil {
		return err
	}
	// nic
	for i, nic := range req.NetworkInterfaces {
		if err := nic.Validate(); err != nil {
			return err
		}
		if i != 0 && nic.Upstream == "shared" {
			return errors.New("upstream=shared is not supported for additional NICs")
		}
	}
	return nil
}

func (req *ApplyRequest) nicSetting() server.NICSettingHolder {
	if len(req.NetworkInterfaces) == 0 {
		return nil
	}
	return req.NetworkInterfaces[0].NICSettingHolder()
}

func (req *ApplyRequest) additionalNICSetting() []server.AdditionalNICSettingHolder {
	var results []server.AdditionalNICSettingHolder
	for i, s := range req.NetworkInterfaces {
		if i == 0 {
			continue
		}
		results = append(results, s.AdditionalNICSettingHolder())
	}
	return results
}

func (req *ApplyRequest) Builder(caller iaas.APICaller) (*server.Builder, error) {
	var diskBuilders []diskBuilder.Builder
	for _, d := range req.Disks {
		b, err := d.Builder(caller)
		if err != nil {
			return nil, err
		}
		diskBuilders = append(diskBuilders, b)
	}

	return &server.Builder{
		Name:            req.Name,
		CPU:             req.CPU,
		MemoryGB:        req.MemoryGB,
		GPU:             req.GPU,
		GPUModel:        req.GPUModel,
		CPUModel:        req.CPUModel,
		Commitment:      req.Commitment,
		Generation:      req.Generation,
		ConfidentialVM:  req.ConfidentialVM,
		InterfaceDriver: req.InterfaceDriver,
		Description:     req.Description,
		IconID:          req.IconID,
		Tags:            req.Tags,
		BootAfterCreate: req.BootAfterCreate,
		CDROMID:         req.CDROMID,
		PrivateHostID:   req.PrivateHostID,
		NIC:             req.nicSetting(),
		AdditionalNICs:  req.additionalNICSetting(),
		DiskBuilders:    diskBuilders,
		Client:          server.NewBuildersAPIClient(caller),
		ServerID:        req.ID,
		ForceShutdown:   req.ForceShutdown,
		NoWait:          req.NoWait,
	}, nil
}
