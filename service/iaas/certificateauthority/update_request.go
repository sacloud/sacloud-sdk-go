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

package certificateauthority

import (
	"context"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/certificateauthority/builder"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	ID types.ID `service:"-" validate:"required"`

	Name        *string     `service:",omitempty" validate:"omitempty,min=1"`
	Description *string     `service:",omitempty" validate:"omitempty,max=512"`
	Tags        *types.Tags `service:",omitempty"`
	IconID      *types.ID   `service:",omitempty"`

	Clients []*builder.ClientCert `service:",omitempty"` // Note: API的に証明書の削除はできないため、指定した以上の証明書が存在する可能性がある
	Servers []*builder.ServerCert `service:",omitempty"` // Note: API的に証明書の削除はできないため、指定した以上の証明書が存在する可能性がある

	PollingTimeout  time.Duration // 証明書発行待ちのタイムアウト
	PollingInterval time.Duration // 証明書発行待ちのポーリング間隔
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) ApplyRequest(ctx context.Context, caller iaas.APICaller) (*ApplyRequest, error) {
	client := iaas.NewCertificateAuthorityOp(caller)
	current, err := client.Read(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	applyRequest := &ApplyRequest{
		ID:          req.ID,
		Name:        current.Name,
		Description: current.Description,
		Tags:        current.Tags,
		IconID:      current.IconID,

		Clients:         req.Clients,
		Servers:         req.Servers,
		PollingTimeout:  req.PollingTimeout,
		PollingInterval: req.PollingInterval,
	}

	if err := serviceutil.RequestConvertTo(req, applyRequest); err != nil {
		return nil, err
	}
	return applyRequest, nil
}
