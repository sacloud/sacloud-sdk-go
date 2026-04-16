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
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type UpdateSubnetRequest struct {
	Zone string   `service:"-" validate:"required"`
	ID   types.ID `service:"-" validate:"required"`

	SubnetID types.ID `validate:"required"`
	NextHop  string   `validate:"required,ipv4"`
}

func (req *UpdateSubnetRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateSubnetRequest) ToRequestParameter(current *iaas.Internet) (*iaas.InternetUpdateSubnetRequest, error) {
	r := &iaas.InternetUpdateSubnetRequest{}
	if err := serviceutil.RequestConvertTo(current, r); err != nil {
		return nil, err
	}
	if err := serviceutil.RequestConvertTo(req, r); err != nil {
		return nil, err
	}
	return r, nil
}
