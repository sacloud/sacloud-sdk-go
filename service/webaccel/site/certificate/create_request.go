// Copyright 2022-2025 The sacloud/webaccel-service-go authors
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

package certificate

import (
	"github.com/sacloud/webaccel-api-go"
)

type CreateRequest struct {
	SiteId string `service:"-" validate:"required"`

	CertificateChain string `validate:"required"`
	Key              string `validate:"required"`
}

func (req *CreateRequest) ToRequestParameter() *webaccel.CreateOrUpdateCertificateRequest {
	return &webaccel.CreateOrUpdateCertificateRequest{
		CertificateChain: req.CertificateChain,
		Key:              req.Key,
	}
}
