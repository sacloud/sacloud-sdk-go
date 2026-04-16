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

package disk

import (
	"fmt"

	"github.com/sacloud/iaas-api-go/ostype"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/packages-go/validate"
)

// CreateRequest ディスク作成リクエスト
type CreateRequest struct {
	Zone string `service:"-" validate:"required"`

	Name                       string `validate:"required"`
	Description                string `validate:"min=0,max=512"`
	Tags                       types.Tags
	IconID                     types.ID
	DiskPlanID                 types.ID              `validate:"oneof=4 2"`
	Connection                 types.EDiskConnection `validate:"oneof=virtio ide"`
	EncryptionAlgorithm        types.EDiskEncryptionAlgorithm
	KMSKeyID                   types.ID
	DedicatedStorageContractID types.ID
	SourceDiskID               types.ID
	SourceArchiveID            types.ID
	ServerID                   types.ID
	SizeGB                     int `service:"SizeMB,filters=gb_to_mb"`
	DistantFrom                []types.ID
	OSType                     ostype.ArchiveOSType
	EditParameter              *EditParameter

	NoWait bool
}

func (req *CreateRequest) Validate() error {
	if req.OSType != ostype.Custom {
		if !req.SourceDiskID.IsEmpty() || !req.SourceArchiveID.IsEmpty() {
			return fmt.Errorf("SourceDiskID or SourceArchiveID must be empty if OSType has a value")
		}
	}
	return validate.New().Struct(req)
}

func (req *CreateRequest) ApplyRequest() *ApplyRequest {
	return &ApplyRequest{
		Zone:                       req.Zone,
		Name:                       req.Name,
		Description:                req.Description,
		Tags:                       req.Tags,
		IconID:                     req.IconID,
		DiskPlanID:                 req.DiskPlanID,
		Connection:                 req.Connection,
		EncryptionAlgorithm:        req.EncryptionAlgorithm,
		KMSKeyID:                   req.KMSKeyID,
		DedicatedStorageContractID: req.DedicatedStorageContractID,
		SourceDiskID:               req.SourceDiskID,
		SourceArchiveID:            req.SourceArchiveID,
		ServerID:                   req.ServerID,
		SizeGB:                     req.SizeGB,
		DistantFrom:                req.DistantFrom,
		OSType:                     req.OSType,
		EditParameter:              req.EditParameter,
		NoWait:                     req.NoWait,
	}
}
