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
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
)

// APIClient builderが利用するAPIクライアント群
type APIClient struct {
	Archive  ArchiveFinder
	Disk     CreateDiskHandler
	DiskPlan PlanReader
	Note     NoteHandler
	SSHKey   SSHKeyHandler
}

// ArchiveFinder アーカイブ検索のためのインターフェース
type ArchiveFinder interface {
	Find(ctx context.Context, zone string, conditions *iaas.FindCondition) (*iaas.ArchiveFindResult, error)
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Archive, error)
}

// CreateDiskHandler ディスク操作のためのインターフェース
type CreateDiskHandler interface {
	Create(ctx context.Context, zone string, createParam *iaas.DiskCreateRequest, distantFrom []types.ID, kmsKeyID types.ID) (*iaas.Disk, error)
	CreateOnDedicatedStorage(ctx context.Context, zone string, createParam *iaas.DiskCreateRequest, distantFrom []types.ID, kmeKeyID types.ID, dedicatedStorageContractID types.ID) (*iaas.Disk, error)
	CreateWithConfig(
		ctx context.Context,
		zone string,
		createParam *iaas.DiskCreateRequest,
		editParam *iaas.DiskEditRequest,
		bootAtAvailable bool,
		distantFrom []types.ID,
		kmsKeyID types.ID,
	) (*iaas.Disk, error)
	CreateOnDedicatedStorageWithConfig(
		ctx context.Context,
		zone string,
		createParam *iaas.DiskCreateRequest,
		editParam *iaas.DiskEditRequest,
		bootAtAvailable bool,
		distantFrom []types.ID,
		kmeKeyID types.ID,
		dedicatedStorageContractID types.ID,
	) (*iaas.Disk, error)
	Update(ctx context.Context, zone string, id types.ID, updateParam *iaas.DiskUpdateRequest) (*iaas.Disk, error)
	Config(ctx context.Context, zone string, id types.ID, editParam *iaas.DiskEditRequest) error
	Read(ctx context.Context, zone string, id types.ID) (*iaas.Disk, error)
	ConnectToServer(ctx context.Context, zone string, id types.ID, serverID types.ID) error
}

// PlanReader ディスクプラン取得のためのインターフェース
type PlanReader interface {
	Read(ctx context.Context, zone string, id types.ID) (*iaas.DiskPlan, error)
}

// NoteHandler スタートアップスクリプト参照のためのインターフェース
type NoteHandler interface {
	Read(ctx context.Context, id types.ID) (*iaas.Note, error)
	Create(ctx context.Context, param *iaas.NoteCreateRequest) (*iaas.Note, error)
	Delete(ctx context.Context, id types.ID) error
}

// SSHKeyHandler SSHKey参照のためのインターフェース
type SSHKeyHandler interface {
	Read(ctx context.Context, id types.ID) (*iaas.SSHKey, error)
	Delete(ctx context.Context, id types.ID) error
}

// NewBuildersAPIClient APIクライアントの作成
func NewBuildersAPIClient(caller iaas.APICaller) *APIClient {
	return &APIClient{
		Archive:  iaas.NewArchiveOp(caller),
		Disk:     iaas.NewDiskOp(caller),
		DiskPlan: iaas.NewDiskPlanOp(caller),
		Note:     iaas.NewNoteOp(caller),
		SSHKey:   iaas.NewSSHKeyOp(caller),
	}
}
