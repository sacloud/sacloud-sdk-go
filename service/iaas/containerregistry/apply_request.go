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

package containerregistry

import (
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/containerregistry/builder"
	"github.com/sacloud/packages-go/validate"
)

type ApplyRequest struct {
	ID types.ID `service:"-"`

	Name           string `validate:"required"`
	Description    string `validate:"min=0,max=512"`
	Tags           types.Tags
	IconID         types.ID
	AccessLevel    types.EContainerRegistryAccessLevel
	VirtualDomain  string
	SubDomainLabel string
	Users          []*builder.User

	SettingsHash string
}

func (req *ApplyRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *ApplyRequest) Builder(caller iaas.APICaller) (*builder.Builder, error) {
	return &builder.Builder{
		ID:             req.ID,
		Name:           req.Name,
		Description:    req.Description,
		Tags:           req.Tags,
		IconID:         req.IconID,
		AccessLevel:    req.AccessLevel,
		VirtualDomain:  req.VirtualDomain,
		SubDomainLabel: req.SubDomainLabel,
		Users:          req.Users,
		SettingsHash:   req.SettingsHash,
		Client:         iaas.NewContainerRegistryOp(caller),
	}, nil
}
