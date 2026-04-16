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

package internet

import (
	"context"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/internet/builder"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	Zone string   `service:"-" validate:"required"`
	ID   types.ID `service:"-" validate:"required"`

	Name          *string     `service:",omitempty" validate:"omitempty,min=1"`
	Description   *string     `service:",omitempty" validate:"omitempty,max=512"`
	Tags          *types.Tags `service:",omitempty"`
	IconID        *types.ID   `service:",omitempty"`
	BandWidthMbps *int        `service:",omitempty"`
	EnableIPv6    *bool       `service:",omitempty"`
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) Builder(ctx context.Context, caller iaas.APICaller) (*builder.Builder, error) {
	current, err := iaas.NewInternetOp(caller).Read(ctx, req.Zone, req.ID)
	if err != nil {
		return nil, err
	}

	builder := &builder.Builder{
		Name:           current.Name,
		Description:    current.Description,
		Tags:           current.Tags,
		IconID:         current.IconID,
		NetworkMaskLen: current.NetworkMaskLen,
		BandWidthMbps:  current.BandWidthMbps,
		EnableIPv6:     len(current.Switch.IPv6Nets) > 0,
		Client:         builder.NewAPIClient(caller),
	}

	if err := serviceutil.RequestConvertTo(req, builder); err != nil {
		return nil, err
	}
	return builder, nil
}
