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

package builder

import (
	"context"
	"errors"
	"reflect"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/query"
	"github.com/sacloud/iaas-api-go/helper/wait"
	"github.com/sacloud/iaas-api-go/types"
)

type Builder struct {
	ID   types.ID
	Zone string

	Name        string
	Description string
	Tags        types.Tags
	IconID      types.ID
	SwitchID    types.ID

	Plan types.ID
	Size types.ENFSSize

	IPAddresses    []string
	NetworkMaskLen int
	DefaultRoute   string

	Caller iaas.APICaller
	NoWait bool
}

func (b *Builder) Build(ctx context.Context) (*iaas.NFS, error) {
	if b.ID.IsEmpty() {
		return b.create(ctx)
	}
	return b.update(ctx)
}

func (b *Builder) findPlanID(ctx context.Context) (types.ID, error) {
	return query.FindNFSPlanID(ctx, iaas.NewNoteOp(b.Caller), b.Plan, b.Size)
}

func (b *Builder) create(ctx context.Context) (*iaas.NFS, error) {
	client := iaas.NewNFSOp(b.Caller)
	planID, err := b.findPlanID(ctx)
	if err != nil {
		return nil, err
	}
	created, err := client.Create(ctx, b.Zone, &iaas.NFSCreateRequest{
		SwitchID:       b.SwitchID,
		PlanID:         planID,
		IPAddresses:    b.IPAddresses,
		NetworkMaskLen: b.NetworkMaskLen,
		DefaultRoute:   b.DefaultRoute,
		Name:           b.Name,
		Description:    b.Description,
		Tags:           b.Tags,
		IconID:         b.IconID,
	})
	if err != nil {
		return nil, err
	}
	if b.NoWait {
		return created, nil
	}

	return wait.UntilNFSIsUp(ctx, client, b.Zone, created.ID)
}

func (b *Builder) validateForUpdate(ctx context.Context, current *iaas.NFS) error {
	planID, err := b.findPlanID(ctx)
	if err != nil {
		return err
	}
	if current.SwitchID != b.SwitchID {
		return errors.New("SwitchID cannot be changed")
	}
	if current.PlanID != planID {
		return errors.New("Plan/Size cannot be changed")
	}
	if !reflect.DeepEqual(current.IPAddresses, b.IPAddresses) {
		return errors.New("IPAddresses cannot be changed")
	}
	if current.NetworkMaskLen != b.NetworkMaskLen {
		return errors.New("NetworkMaskLen cannot be changed")
	}
	if current.DefaultRoute != b.DefaultRoute {
		return errors.New("DefaultRoute cannot be changed")
	}
	return nil
}

func (b *Builder) update(ctx context.Context) (*iaas.NFS, error) {
	client := iaas.NewNFSOp(b.Caller)
	current, err := client.Read(ctx, b.Zone, b.ID)
	if err != nil {
		return nil, err
	}
	if err := b.validateForUpdate(ctx, current); err != nil {
		return nil, err
	}

	return client.Update(ctx, b.Zone, b.ID, &iaas.NFSUpdateRequest{
		Name:        b.Name,
		Description: b.Description,
		Tags:        b.Tags,
		IconID:      b.IconID,
	})
}
