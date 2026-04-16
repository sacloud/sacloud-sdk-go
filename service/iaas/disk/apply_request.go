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
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/ostype"
	"github.com/sacloud/iaas-api-go/types"
	disk "github.com/sacloud/iaas-service-go/disk/builder"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type ApplyRequest struct {
	Zone string   `service:"-" validate:"required"`
	ID   types.ID `service:"-"` // TODO Builderを更新対応させる(変更できない値の場合はエラーにするとか)

	Name                       string `validate:"required"`
	Description                string `validate:"min=0,max=512"`
	Tags                       types.Tags
	IconID                     types.ID
	DiskPlanID                 types.ID
	Connection                 types.EDiskConnection
	EncryptionAlgorithm        types.EDiskEncryptionAlgorithm
	KMSKeyID                   types.ID
	DedicatedStorageContractID types.ID
	SourceDiskID               types.ID
	SourceArchiveID            types.ID
	ServerID                   types.ID
	SizeGB                     int
	DistantFrom                []types.ID

	OSType        ostype.ArchiveOSType
	EditParameter *EditParameter

	NoWait bool
}

// EditParameter ディスクの修正用パラメータ
type EditParameter struct {
	HostName string
	Password string

	DisablePWAuth       bool
	EnableDHCP          bool
	ChangePartitionUUID bool

	IPAddress      string
	NetworkMaskLen int
	DefaultRoute   string

	SSHKeys   []string
	SSHKeyIDs []types.ID

	// IsSSHKeysEphemeral trueの場合、SSHキーを生成する場合に生成したSSHキーリソースをサーバ作成後に削除する
	IsSSHKeysEphemeral bool

	IsNotesEphemeral bool
	NoteContents     []string
	Notes            []*iaas.DiskEditNote
}

func (req *ApplyRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *ApplyRequest) Builder(caller iaas.APICaller) (disk.Builder, error) {
	var editParameter *disk.EditRequest

	if req.EditParameter != nil {
		editParameter = &disk.EditRequest{}
		if err := serviceutil.RequestConvertTo(req.EditParameter, editParameter); err != nil {
			return nil, err
		}
	}

	director := &disk.Director{
		OSType:                     req.OSType,
		Name:                       req.Name,
		SizeGB:                     req.SizeGB,
		DistantFrom:                req.DistantFrom,
		PlanID:                     req.DiskPlanID,
		Connection:                 req.Connection,
		EncryptionAlgorithm:        req.EncryptionAlgorithm,
		KMSKeyID:                   req.KMSKeyID,
		DedicatedStorageContractID: req.DedicatedStorageContractID,
		Description:                req.Description,
		Tags:                       req.Tags,
		IconID:                     req.IconID,
		DiskID:                     req.ID,
		SourceDiskID:               req.SourceDiskID,
		SourceArchiveID:            req.SourceArchiveID,
		EditParameter:              editParameter,
		NoWait:                     req.NoWait,
		Client:                     disk.NewBuildersAPIClient(caller),
	}
	return director.Builder(), nil
}
