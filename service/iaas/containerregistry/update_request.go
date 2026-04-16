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
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/containerregistry/builder"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	ID types.ID `service:"-" validate:"required"`

	Name          *string                              `service:",omitempty" validate:"omitempty,min=1"`
	Description   *string                              `service:",omitempty" validate:"omitempty,max=512"`
	Tags          *types.Tags                          `service:",omitempty"`
	IconID        *types.ID                            `service:",omitempty"`
	AccessLevel   *types.EContainerRegistryAccessLevel `service:",omitempty"`
	VirtualDomain *string                              `service:",omitempty"`
	Users         *[]*builder.User                     `service:",omitempty"`
	SettingsHash  string
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) ApplyRequest(ctx context.Context, caller iaas.APICaller) (*ApplyRequest, error) {
	client := iaas.NewContainerRegistryOp(caller)
	current, err := client.Read(ctx, req.ID)
	if err != nil {
		return nil, err
	}
	users, err := client.ListUsers(ctx, req.ID) // NOTE: ユーザーが登録されていなくても200が返る
	if err != nil {
		return nil, err
	}

	applyRequest := &ApplyRequest{
		ID:             req.ID,
		Name:           current.Name,
		Description:    current.Description,
		Tags:           current.Tags,
		IconID:         current.IconID,
		AccessLevel:    current.AccessLevel,
		VirtualDomain:  current.VirtualDomain,
		SubDomainLabel: current.SubDomainLabel,
		SettingsHash:   current.SettingsHash,
	}
	if users != nil {
		for _, user := range users.Users {
			applyRequest.Users = append(applyRequest.Users, &builder.User{
				UserName:   user.UserName,
				Password:   "", // パスワードは参照できないため常に空
				Permission: user.Permission,
			})
		}
	}

	if err := serviceutil.RequestConvertTo(req, applyRequest); err != nil {
		return nil, err
	}
	return applyRequest, nil
}
