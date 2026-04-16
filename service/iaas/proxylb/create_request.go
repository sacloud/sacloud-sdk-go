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

package proxylb

import (
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/packages-go/validate"
)

type CreateRequest struct {
	Name                 string `validate:"required"`
	Description          string `validate:"min=0,max=512"`
	Tags                 types.Tags
	IconID               types.ID
	Plan                 types.EProxyLBPlan `validate:"required,oneof=100 500 1000 5000 10000 50000 100000 400000"`
	HealthCheck          *iaas.ProxyLBHealthCheck
	SorryServer          *iaas.ProxyLBSorryServer
	BindPorts            []*iaas.ProxyLBBindPort
	Servers              []*iaas.ProxyLBServer
	Rules                []*iaas.ProxyLBRule
	LetsEncrypt          *iaas.ProxyLBACMESetting
	StickySession        *iaas.ProxyLBStickySession
	Gzip                 *iaas.ProxyLBGzip
	BackendHttpKeepAlive *iaas.ProxyLBBackendHttpKeepAlive
	ProxyProtocol        *iaas.ProxyLBProxyProtocol
	Syslog               *iaas.ProxyLBSyslog
	Timeout              *iaas.ProxyLBTimeout
	UseVIPFailover       bool
	Region               types.EProxyLBRegion
	MonitoringSuiteLog   *iaas.MonitoringSuiteLog
}

func (req *CreateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *CreateRequest) ToRequestParameter() (*iaas.ProxyLBCreateRequest, error) {
	params := &iaas.ProxyLBCreateRequest{}
	if err := serviceutil.RequestConvertTo(req, params); err != nil {
		return nil, err
	}
	return params, nil
}
