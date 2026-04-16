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
	"errors"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/internet/builder"
	"github.com/sacloud/packages-go/validate"
)

type CreateRequest struct {
	Zone string `service:"-" validate:"required"`

	Name           string `validate:"required"`
	Description    string `validate:"min=0,max=512"`
	Tags           types.Tags
	IconID         types.ID
	NetworkMaskLen int
	BandWidthMbps  int
	EnableIPv6     bool
	NoWait         bool
	NotFoundRetry  int // スイッチ+ルータは作成直後だと404を返すことがあることへの対応でリトライする際のリトライ上限回数、省略時はDefaultNotFoundRetry
}

func (req *CreateRequest) Validate() error {
	if err := validate.New().Struct(req); err != nil {
		return err
	}
	if req.EnableIPv6 && req.NoWait {
		return errors.New("NoWait=true is not supported when EnableIPv6=true")
	}
	return nil
}

func (req *CreateRequest) Builder(caller iaas.APICaller) *builder.Builder {
	return &builder.Builder{
		Name:           req.Name,
		Description:    req.Description,
		Tags:           req.Tags,
		IconID:         req.IconID,
		NetworkMaskLen: req.NetworkMaskLen,
		BandWidthMbps:  req.BandWidthMbps,
		EnableIPv6:     req.EnableIPv6,
		NotFoundRetry:  req.NotFoundRetry,
		NoWait:         req.NoWait,
		Client:         builder.NewAPIClient(caller),
	}
}
