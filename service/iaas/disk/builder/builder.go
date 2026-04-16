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
	"errors"
	"fmt"
	"reflect"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/ostype"
	"github.com/sacloud/iaas-api-go/types"
	service "github.com/sacloud/iaas-service-go"
	"github.com/sacloud/packages-go/size"
)

// Builder ディスクの構築インターフェース
type Builder interface {
	Validate(ctx context.Context, zone string) error
	Build(ctx context.Context, zone string, serverID types.ID) (*BuildResult, error)
	Update(ctx context.Context, zone string) (*UpdateResult, error)
	DiskID() types.ID
	UpdateLevel(ctx context.Context, zone string, disk *iaas.Disk) service.UpdateLevel
	NoWaitFlag() bool
}

// BuildResult ディスク構築結果
type BuildResult struct {
	DiskID types.ID
}

// UpdateResult ディスク更新結果
type UpdateResult struct {
	Disk *iaas.Disk
}

// FromUnixBuilder Unix系パブリックアーカイブからディスクを作成するリクエスト
type FromUnixBuilder struct {
	OSType ostype.ArchiveOSType

	Name                       string
	SizeGB                     int
	DistantFrom                []types.ID
	PlanID                     types.ID
	Connection                 types.EDiskConnection
	EncryptionAlgorithm        types.EDiskEncryptionAlgorithm
	KMSKeyID                   types.ID
	DedicatedStorageContractID types.ID
	Description                string
	Tags                       types.Tags
	IconID                     types.ID

	EditParameter *UnixEditRequest

	Client *APIClient
	NoWait bool

	ID types.ID

	generatedNotes []*iaas.Note
}

// Validate 設定値の検証
func (d *FromUnixBuilder) Validate(ctx context.Context, zone string) error {
	if !d.OSType.IsSupportDiskEdit() {
		return fmt.Errorf("invalid OSType: %s", d.OSType.String())
	}
	if err := validateDiskPlan(ctx, d.Client, zone, d.PlanID, d.SizeGB); err != nil {
		return err
	}

	if d.EditParameter != nil {
		return d.EditParameter.Validate(ctx, d.Client)
	}
	return nil
}

// Build ディスクの構築
func (d *FromUnixBuilder) Build(ctx context.Context, zone string, serverID types.ID) (*BuildResult, error) {
	res, err := build(ctx, d.Client, zone, serverID, d.DistantFrom, d.KMSKeyID, d.DedicatedStorageContractID, d)
	if err != nil {
		return nil, err
	}
	d.ID = res.DiskID

	if d.EditParameter != nil {
		if d.EditParameter.IsNotesEphemeral {
			for _, note := range d.generatedNotes {
				if err := d.Client.Note.Delete(ctx, note.ID); err != nil {
					return nil, err
				}
			}
		}
	}
	return res, nil
}

// Update ディスクの更新
func (d *FromUnixBuilder) Update(ctx context.Context, zone string) (*UpdateResult, error) {
	return update(ctx, d.Client, zone, d)
}

// DiskID ディスクID取得
func (d *FromUnixBuilder) DiskID() types.ID {
	return d.ID
}

// UpdateLevel Update時にどのレベルの変更が必要か
func (d *FromUnixBuilder) UpdateLevel(ctx context.Context, zone string, disk *iaas.Disk) service.UpdateLevel {
	return updateLevel(disk, d.EditParameter != nil, d)
}

func (d *FromUnixBuilder) updateDiskParameter() *iaas.DiskUpdateRequest {
	return &iaas.DiskUpdateRequest{
		Name:        d.Name,
		Description: d.Description,
		Tags:        d.Tags,
		IconID:      d.IconID,
		Connection:  d.Connection,
	}
}

func (d *FromUnixBuilder) createDiskParameter(ctx context.Context, client *APIClient, zone string, serverID types.ID) (*iaas.DiskCreateRequest, *iaas.DiskEditRequest, error) {
	archive, err := query.FindArchiveByOSType(ctx, client.Archive, zone, d.OSType)
	if err != nil {
		return nil, nil, err
	}

	createReq := &iaas.DiskCreateRequest{
		DiskPlanID:          d.PlanID,
		SizeMB:              d.SizeGB * size.GiB,
		Connection:          d.Connection,
		EncryptionAlgorithm: d.EncryptionAlgorithm,
		SourceArchiveID:     archive.ID,
		ServerID:            serverID,
		Name:                d.Name,
		Description:         d.Description,
		Tags:                d.Tags,
		IconID:              d.IconID,
	}

	var editReq *iaas.DiskEditRequest
	if d.EditParameter != nil {
		req, notes, err := d.EditParameter.prepareDiskEditParameter(ctx, client)
		if err != nil {
			return nil, nil, err
		}
		editReq = req
		if len(notes) > 0 {
			d.generatedNotes = notes
		}
	}

	return createReq, editReq, nil
}

func (d *FromUnixBuilder) NoWaitFlag() bool {
	return d.NoWait
}

// FromFixedArchiveBuilder ディスクの修正をサポートしないパブリックアーカイブからディスクを作成するリクエスト
type FromFixedArchiveBuilder struct {
	OSType ostype.ArchiveOSType

	Name                       string
	SizeGB                     int
	DistantFrom                []types.ID
	PlanID                     types.ID
	Connection                 types.EDiskConnection
	EncryptionAlgorithm        types.EDiskEncryptionAlgorithm
	KMSKeyID                   types.ID
	DedicatedStorageContractID types.ID
	Description                string
	Tags                       types.Tags
	IconID                     types.ID

	Client *APIClient
	NoWait bool

	ID types.ID
}

// Validate 設定値の検証
func (d *FromFixedArchiveBuilder) Validate(ctx context.Context, zone string) error {
	if d.OSType.IsSupportDiskEdit() {
		return fmt.Errorf("invalid OSType: %s", d.OSType.String())
	}
	if err := validateDiskPlan(ctx, d.Client, zone, d.PlanID, d.SizeGB); err != nil {
		return err
	}

	return nil
}

// Build ディスクの構築
func (d *FromFixedArchiveBuilder) Build(ctx context.Context, zone string, serverID types.ID) (*BuildResult, error) {
	res, err := build(ctx, d.Client, zone, serverID, d.DistantFrom, d.KMSKeyID, d.DedicatedStorageContractID, d)
	if err != nil {
		return nil, err
	}
	d.ID = res.DiskID
	return res, nil
}

// Update ディスクの更新
func (d *FromFixedArchiveBuilder) Update(ctx context.Context, zone string) (*UpdateResult, error) {
	return update(ctx, d.Client, zone, d)
}

// DiskID ディスクID取得
func (d *FromFixedArchiveBuilder) DiskID() types.ID {
	return d.ID
}

// UpdateLevel Update時にどのレベルの変更が必要か
func (d *FromFixedArchiveBuilder) UpdateLevel(ctx context.Context, zone string, disk *iaas.Disk) service.UpdateLevel {
	return updateLevel(disk, false, d)
}

func (d *FromFixedArchiveBuilder) updateDiskParameter() *iaas.DiskUpdateRequest {
	return &iaas.DiskUpdateRequest{
		Name:        d.Name,
		Description: d.Description,
		Tags:        d.Tags,
		IconID:      d.IconID,
		Connection:  d.Connection,
	}
}

func (d *FromFixedArchiveBuilder) createDiskParameter(ctx context.Context, client *APIClient, zone string, serverID types.ID) (*iaas.DiskCreateRequest, *iaas.DiskEditRequest, error) {
	archive, err := query.FindArchiveByOSType(ctx, client.Archive, zone, d.OSType)
	if err != nil {
		return nil, nil, err
	}

	createReq := &iaas.DiskCreateRequest{
		DiskPlanID:          d.PlanID,
		SizeMB:              d.SizeGB * size.GiB,
		Connection:          d.Connection,
		EncryptionAlgorithm: d.EncryptionAlgorithm,
		SourceArchiveID:     archive.ID,
		ServerID:            serverID,
		Name:                d.Name,
		Description:         d.Description,
		Tags:                d.Tags,
		IconID:              d.IconID,
	}
	return createReq, nil, nil
}

func (d *FromFixedArchiveBuilder) NoWaitFlag() bool {
	return d.NoWait
}

// FromDiskOrArchiveBuilder ディスクorアーカイブからディスクを作成するリクエスト
//
// ディスクの修正が可能かは実行時にさくらのクラウドAPI側にて判定される
type FromDiskOrArchiveBuilder struct {
	SourceDiskID    types.ID
	SourceArchiveID types.ID

	Name                       string
	SizeGB                     int
	DistantFrom                []types.ID
	PlanID                     types.ID
	Connection                 types.EDiskConnection
	EncryptionAlgorithm        types.EDiskEncryptionAlgorithm
	KMSKeyID                   types.ID
	DedicatedStorageContractID types.ID
	Description                string
	Tags                       types.Tags
	IconID                     types.ID

	EditParameter *UnixEditRequest

	Client *APIClient

	ID     types.ID
	NoWait bool

	generatedNotes []*iaas.Note
}

// Validate 設定値の検証
func (d *FromDiskOrArchiveBuilder) Validate(ctx context.Context, zone string) error {
	if d.SourceArchiveID.IsEmpty() && d.SourceDiskID.IsEmpty() {
		return errors.New("SourceArchiveID or SourceDiskID is required")
	}
	if err := validateDiskPlan(ctx, d.Client, zone, d.PlanID, d.SizeGB); err != nil {
		return err
	}

	if !d.SourceArchiveID.IsEmpty() {
		if _, err := d.Client.Archive.Read(ctx, zone, d.SourceArchiveID); err != nil {
			return err
		}
	}
	if !d.SourceDiskID.IsEmpty() {
		if _, err := d.Client.Disk.Read(ctx, zone, d.SourceDiskID); err != nil {
			return err
		}
	}

	return nil
}

// Build ディスクの構築
func (d *FromDiskOrArchiveBuilder) Build(ctx context.Context, zone string, serverID types.ID) (*BuildResult, error) {
	res, err := build(ctx, d.Client, zone, serverID, d.DistantFrom, d.KMSKeyID, d.DedicatedStorageContractID, d)
	if err != nil {
		return nil, err
	}
	d.ID = res.DiskID
	if d.EditParameter != nil {
		if d.EditParameter.IsNotesEphemeral {
			for _, note := range d.generatedNotes {
				if err := d.Client.Note.Delete(ctx, note.ID); err != nil {
					return nil, err
				}
			}
		}
	}
	return res, nil
}

// Update ディスクの更新
func (d *FromDiskOrArchiveBuilder) Update(ctx context.Context, zone string) (*UpdateResult, error) {
	return update(ctx, d.Client, zone, d)
}

// DiskID ディスクID取得
func (d *FromDiskOrArchiveBuilder) DiskID() types.ID {
	return d.ID
}

// UpdateLevel Update時にどのレベルの変更が必要か
func (d *FromDiskOrArchiveBuilder) UpdateLevel(ctx context.Context, zone string, disk *iaas.Disk) service.UpdateLevel {
	return updateLevel(disk, d.EditParameter != nil, d)
}

func (d *FromDiskOrArchiveBuilder) updateDiskParameter() *iaas.DiskUpdateRequest {
	return &iaas.DiskUpdateRequest{
		Name:        d.Name,
		Description: d.Description,
		Tags:        d.Tags,
		IconID:      d.IconID,
		Connection:  d.Connection,
	}
}

func (d *FromDiskOrArchiveBuilder) createDiskParameter(ctx context.Context, client *APIClient, zone string, serverID types.ID) (*iaas.DiskCreateRequest, *iaas.DiskEditRequest, error) {
	createReq := &iaas.DiskCreateRequest{
		DiskPlanID:          d.PlanID,
		SizeMB:              d.SizeGB * size.GiB,
		Connection:          d.Connection,
		EncryptionAlgorithm: d.EncryptionAlgorithm,
		SourceArchiveID:     d.SourceArchiveID,
		SourceDiskID:        d.SourceDiskID,
		ServerID:            serverID,
		Name:                d.Name,
		Description:         d.Description,
		Tags:                d.Tags,
		IconID:              d.IconID,
	}

	var editReq *iaas.DiskEditRequest
	if d.EditParameter != nil {
		req, notes, err := d.EditParameter.prepareDiskEditParameter(ctx, client)
		if err != nil {
			return nil, nil, err
		}
		editReq = req
		if len(notes) > 0 {
			d.generatedNotes = notes
		}
	}

	return createReq, editReq, nil
}

func (d *FromDiskOrArchiveBuilder) NoWaitFlag() bool {
	return d.NoWait
}

// BlankBuilder ブランクディスクを作成する場合のリクエスト
type BlankBuilder struct {
	Name                       string
	SizeGB                     int
	DistantFrom                []types.ID
	PlanID                     types.ID
	Connection                 types.EDiskConnection
	EncryptionAlgorithm        types.EDiskEncryptionAlgorithm
	KMSKeyID                   types.ID
	DedicatedStorageContractID types.ID
	Description                string
	Tags                       types.Tags
	IconID                     types.ID

	Client *APIClient
	NoWait bool
	ID     types.ID
}

// Validate 設定値の検証
func (d *BlankBuilder) Validate(ctx context.Context, zone string) error {
	if err := validateDiskPlan(ctx, d.Client, zone, d.PlanID, d.SizeGB); err != nil {
		return err
	}
	return nil
}

// Build ディスクの構築
func (d *BlankBuilder) Build(ctx context.Context, zone string, serverID types.ID) (*BuildResult, error) {
	res, err := build(ctx, d.Client, zone, serverID, d.DistantFrom, d.KMSKeyID, d.DedicatedStorageContractID, d)
	if err != nil {
		return nil, err
	}
	d.ID = res.DiskID
	return res, err
}

// Update ディスクの更新
func (d *BlankBuilder) Update(ctx context.Context, zone string) (*UpdateResult, error) {
	return update(ctx, d.Client, zone, d)
}

// DiskID ディスクID取得
func (d *BlankBuilder) DiskID() types.ID {
	return d.ID
}

// UpdateLevel Update時にどのレベルの変更が必要か
func (d *BlankBuilder) UpdateLevel(ctx context.Context, zone string, disk *iaas.Disk) service.UpdateLevel {
	return updateLevel(disk, false, d)
}

func (d *BlankBuilder) updateDiskParameter() *iaas.DiskUpdateRequest {
	return &iaas.DiskUpdateRequest{
		Name:        d.Name,
		Description: d.Description,
		Tags:        d.Tags,
		IconID:      d.IconID,
		Connection:  d.Connection,
	}
}

func (d *BlankBuilder) createDiskParameter(ctx context.Context, client *APIClient, zone string, serverID types.ID) (*iaas.DiskCreateRequest, *iaas.DiskEditRequest, error) {
	createReq := &iaas.DiskCreateRequest{
		DiskPlanID:          d.PlanID,
		SizeMB:              d.SizeGB * size.GiB,
		Connection:          d.Connection,
		EncryptionAlgorithm: d.EncryptionAlgorithm,
		ServerID:            serverID,
		Name:                d.Name,
		Description:         d.Description,
		Tags:                d.Tags,
		IconID:              d.IconID,
	}
	return createReq, nil, nil
}

func (d *BlankBuilder) NoWaitFlag() bool {
	return d.NoWait
}

// ConnectedDiskBuilder 既存ディスクを接続する場合のリクエスト
type ConnectedDiskBuilder struct {
	ID            types.ID
	EditParameter *UnixEditRequest

	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID
	Connection  types.EDiskConnection

	NoWait bool
	Client *APIClient
}

// Validate 設定値の検証
func (d *ConnectedDiskBuilder) Validate(ctx context.Context, zone string) error {
	if d.ID.IsEmpty() {
		return errors.New("DiskID is required")
	}

	if _, err := d.Client.Disk.Read(ctx, zone, d.ID); err != nil {
		return err
	}

	return nil
}

// Build ディスクの構築
func (d *ConnectedDiskBuilder) Build(ctx context.Context, zone string, serverID types.ID) (*BuildResult, error) {
	res := &BuildResult{
		DiskID: d.ID,
	}
	if !serverID.IsEmpty() {
		if err := d.Client.Disk.ConnectToServer(ctx, zone, d.ID, serverID); err != nil {
			return nil, err
		}
	}

	if d.EditParameter != nil {
		req, _, err := d.EditParameter.prepareDiskEditParameter(ctx, d.Client)
		if err != nil {
			return nil, err
		}
		if err := d.Client.Disk.Config(ctx, zone, d.ID, req); err != nil {
			return nil, err
		}
		waiter := iaas.WaiterForReady(func() (interface{}, error) {
			return d.Client.Disk.Read(ctx, zone, d.ID)
		})
		if _, err := waiter.WaitForState(ctx); err != nil {
			return nil, err
		}
	}
	return res, nil
}

// Update ディスクの更新
func (d *ConnectedDiskBuilder) Update(ctx context.Context, zone string) (*UpdateResult, error) {
	disk, err := d.Client.Disk.Update(ctx, zone, d.ID, d.updateDiskParameter())
	if err != nil {
		return nil, err
	}

	if d.EditParameter != nil {
		req, _, err := d.EditParameter.prepareDiskEditParameter(ctx, d.Client)
		if err != nil {
			return nil, err
		}
		if err := d.Client.Disk.Config(ctx, zone, d.ID, req); err != nil {
			return nil, err
		}
		waiter := iaas.WaiterForReady(func() (interface{}, error) {
			return d.Client.Disk.Read(ctx, zone, d.ID)
		})
		if _, err := waiter.WaitForState(ctx); err != nil {
			return nil, err
		}
	}

	return &UpdateResult{Disk: disk}, nil
}

// DiskID ディスクID取得
func (d *ConnectedDiskBuilder) DiskID() types.ID {
	return d.ID
}

// UpdateLevel Update時にどのレベルの変更が必要か
func (d *ConnectedDiskBuilder) UpdateLevel(ctx context.Context, zone string, disk *iaas.Disk) service.UpdateLevel {
	return updateLevel(disk, d.EditParameter != nil, d)
}

func (d *ConnectedDiskBuilder) updateDiskParameter() *iaas.DiskUpdateRequest {
	return &iaas.DiskUpdateRequest{
		Name:        d.Name,
		Description: d.Description,
		Tags:        d.Tags,
		IconID:      d.IconID,
		Connection:  d.Connection,
	}
}

func (d *ConnectedDiskBuilder) createDiskParameter(
	_ context.Context,
	_ *APIClient,
	_ string,
	_ types.ID,
) (*iaas.DiskCreateRequest, *iaas.DiskEditRequest, error) {
	// noop
	return nil, nil, nil
}

func (d *ConnectedDiskBuilder) NoWaitFlag() bool {
	return d.NoWait
}

type diskBuilder interface {
	createDiskParameter(
		ctx context.Context,
		client *APIClient,
		zone string,
		serverID types.ID,
	) (*iaas.DiskCreateRequest, *iaas.DiskEditRequest, error)
	updateDiskParameter() *iaas.DiskUpdateRequest
	DiskID() types.ID
	NoWaitFlag() bool
}

func build(ctx context.Context, client *APIClient, zone string, serverID types.ID, distantFrom []types.ID, kmsKeyID, dedicatedStorageContractID types.ID, builder diskBuilder) (*BuildResult, error) {
	var err error

	diskReq, editReq, err := builder.createDiskParameter(ctx, client, zone, serverID)
	if err != nil {
		return nil, err
	}
	if diskReq == nil {
		return nil, fmt.Errorf("disk create request is nil")
	}
	diskReq.ServerID = serverID

	var disk *iaas.Disk

	if editReq == nil {
		if dedicatedStorageContractID.IsEmpty() {
			disk, err = client.Disk.Create(ctx, zone, diskReq, distantFrom, kmsKeyID)
		} else {
			disk, err = client.Disk.CreateOnDedicatedStorage(ctx, zone, diskReq, distantFrom, kmsKeyID, dedicatedStorageContractID)
		}
	} else {
		if dedicatedStorageContractID.IsEmpty() {
			disk, err = client.Disk.CreateWithConfig(ctx, zone, diskReq, editReq, false, distantFrom, kmsKeyID)
		} else {
			disk, err = client.Disk.CreateOnDedicatedStorageWithConfig(ctx, zone, diskReq, editReq, false, distantFrom, kmsKeyID, dedicatedStorageContractID)
		}
	}
	if err != nil {
		if disk != nil {
			return &BuildResult{DiskID: disk.ID}, err
		}
		return nil, err
	}

	if builder.NoWaitFlag() {
		return &BuildResult{DiskID: disk.ID}, nil
	}

	waiter := iaas.WaiterForReady(func() (interface{}, error) {
		return client.Disk.Read(ctx, zone, disk.ID)
	})
	lastState, err := waiter.WaitForState(ctx)
	if err != nil {
		if lastState != nil {
			return &BuildResult{DiskID: lastState.(*iaas.Disk).ID}, err
		}
		return nil, err
	}
	disk = lastState.(*iaas.Disk)

	return &BuildResult{DiskID: disk.ID}, nil
}

func update(ctx context.Context, client *APIClient, zone string, builder diskBuilder) (*UpdateResult, error) {
	var err error

	diskID := builder.DiskID()
	if diskID.IsEmpty() {
		return nil, fmt.Errorf("disk id required")
	}

	diskReq, editReq, err := builder.createDiskParameter(ctx, client, zone, types.ID(0))
	if err != nil {
		return nil, err
	}
	if diskReq == nil {
		return nil, fmt.Errorf("disk update request is nil")
	}

	disk, err := client.Disk.Update(ctx, zone, diskID, &iaas.DiskUpdateRequest{
		Name:        diskReq.Name,
		Description: diskReq.Description,
		Tags:        diskReq.Tags,
		IconID:      diskReq.IconID,
		Connection:  diskReq.Connection,
	})
	if err != nil {
		return nil, err
	}

	if editReq != nil {
		if err := client.Disk.Config(ctx, zone, disk.ID, editReq); err != nil {
			return nil, err
		}
	}

	if builder.NoWaitFlag() {
		return &UpdateResult{Disk: disk}, nil
	}

	waiter := iaas.WaiterForReady(func() (interface{}, error) {
		return client.Disk.Read(ctx, zone, disk.ID)
	})
	lastState, err := waiter.WaitForState(ctx)
	if err != nil {
		return nil, err
	}
	disk = lastState.(*iaas.Disk)

	return &UpdateResult{Disk: disk}, nil
}

func validateDiskPlan(ctx context.Context, client *APIClient, zone string, diskPlanID types.ID, sizeGB int) error {
	plan, err := client.DiskPlan.Read(ctx, zone, diskPlanID)
	if err != nil {
		return err
	}
	found := false
	for _, size := range plan.Size {
		if size.Availability.IsAvailable() && size.GetSizeGB() == sizeGB {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("disk plan[%s:%dGB] is not found", plan.Name, sizeGB)
	}
	return nil
}

func updateLevel(disk *iaas.Disk, hasEditReq bool, b diskBuilder) service.UpdateLevel {
	if disk.ID != b.DiskID() || hasEditReq {
		return service.UpdateLevelNeedShutdown
	}

	current := &iaas.DiskUpdateRequest{
		Name:        disk.Name,
		Description: disk.Description,
		Tags:        disk.Tags,
		IconID:      disk.IconID,
		Connection:  disk.Connection,
	}
	desired := b.updateDiskParameter()
	if desired == nil {
		return service.UpdateLevelNone
	}
	if reflect.DeepEqual(current, desired) {
		if current.Connection != desired.Connection {
			return service.UpdateLevelNeedShutdown
		}
		return service.UpdateLevelSimple
	}
	return service.UpdateLevelNone
}
