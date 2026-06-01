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

package note

import (
	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
	"github.com/sacloud/sacloud-sdk-go/common/packages/validate"
	"github.com/sacloud/sacloud-sdk-go/service/iaas/serviceutil"
)

type CreateRequest struct {
	Name    string `validate:"required"`
	Tags    types.Tags
	IconID  types.ID
	Class   string `validate:"required,oneof=shell yaml_cloud_config"`
	Content string `validate:"required"`
}

func (req *CreateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *CreateRequest) ToRequestParameter() (*iaas.NoteCreateRequest, error) {
	params := &iaas.NoteCreateRequest{}
	if err := serviceutil.RequestConvertTo(req, params); err != nil {
		return nil, err
	}
	return params, nil
}
