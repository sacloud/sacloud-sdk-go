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

package mobilegateway

import (
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/mobilegateway/builder"
	"github.com/sacloud/iaas-service-go/serviceutil"
	"github.com/sacloud/iaas-service-go/setup"
	"github.com/sacloud/packages-go/validate"
)

type ApplyRequest struct {
	Zone string `service:"-" validate:"required"`

	ID                              types.ID `service:"-"`
	Name                            string   `validate:"required"`
	Description                     string   `validate:"min=0,max=512"`
	Tags                            types.Tags
	IconID                          types.ID
	PrivateInterface                *PrivateInterfaceSetting `validate:"omitempty"`
	StaticRoutes                    []*iaas.MobileGatewayStaticRoute
	SIMRoutes                       []*SIMRouteSetting
	InternetConnectionEnabled       bool
	InterDeviceCommunicationEnabled bool
	DNS                             *DNSSetting
	SIMs                            []*SIMSetting
	TrafficConfig                   *TrafficConfig

	SettingsHash    string
	BootAfterCreate bool
	NoWait          bool
}

// PrivateInterfaceSetting represents API parameter/response structure
type PrivateInterfaceSetting struct {
	SwitchID       types.ID `service:",omitempty"`
	IPAddress      string   `service:",omitempty" validate:"required,ipv4"`
	NetworkMaskLen int      `service:",omitempty"`
}

// SIMRouteSetting represents API parameter/response structure
type SIMRouteSetting struct {
	SIMID  types.ID
	Prefix string `validate:"required"`
}

// SIMSetting represents API parameter/response structure
type SIMSetting struct {
	SIMID     types.ID
	IPAddress string `validate:"required,ipv4"`
}

type DNSSetting struct {
	DNS1 string `service:",omitempty" validate:"required_with=DNS2,omitempty,ipv4"`
	DNS2 string `service:",omitempty" validate:"required_with=DNS1,omitempty,ipv4"`
}

type TrafficConfig struct {
	TrafficQuotaInMB       int    `service:",omitempty"`
	BandWidthLimitInKbps   int    `service:",omitempty"`
	EmailNotifyEnabled     bool   `service:",omitempty"`
	SlackNotifyEnabled     bool   `service:",omitempty"`
	SlackNotifyWebhooksURL string `service:",omitempty"`
	AutoTrafficShaping     bool   `service:",omitempty"`
}

func (req *ApplyRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *ApplyRequest) Builder(caller iaas.APICaller) (*builder.Builder, error) {
	var privateInterface *builder.PrivateInterfaceSetting
	if req.PrivateInterface != nil {
		privateInterface = &builder.PrivateInterfaceSetting{
			SwitchID:       req.PrivateInterface.SwitchID,
			IPAddress:      req.PrivateInterface.IPAddress,
			NetworkMaskLen: req.PrivateInterface.NetworkMaskLen,
		}
	}

	var simRoutes []*builder.SIMRouteSetting
	for _, sr := range req.SIMRoutes {
		simRoutes = append(simRoutes, &builder.SIMRouteSetting{
			SIMID:  sr.SIMID,
			Prefix: sr.Prefix,
		})
	}

	var sims []*builder.SIMSetting
	for _, s := range req.SIMs {
		sims = append(sims, &builder.SIMSetting{
			SIMID:     s.SIMID,
			IPAddress: s.IPAddress,
		})
	}

	var dns *iaas.MobileGatewayDNSSetting
	if req.DNS != nil {
		dns = &iaas.MobileGatewayDNSSetting{
			DNS1: req.DNS.DNS1,
			DNS2: req.DNS.DNS2,
		}
	}

	var trafficConfig *iaas.MobileGatewayTrafficControl
	if req.TrafficConfig != nil {
		trafficConfig = &iaas.MobileGatewayTrafficControl{}
		if err := serviceutil.RequestConvertTo(req.TrafficConfig, trafficConfig); err != nil {
			return nil, err
		}
	}

	return &builder.Builder{
		ID:                              req.ID,
		Zone:                            req.Zone,
		Name:                            req.Name,
		Description:                     req.Description,
		Tags:                            req.Tags,
		IconID:                          req.IconID,
		PrivateInterface:                privateInterface,
		StaticRoutes:                    req.StaticRoutes,
		SIMRoutes:                       simRoutes,
		InternetConnectionEnabled:       req.InternetConnectionEnabled,
		InterDeviceCommunicationEnabled: req.InterDeviceCommunicationEnabled,
		DNS:                             dns,
		SIMs:                            sims,
		TrafficConfig:                   trafficConfig,
		SettingsHash:                    req.SettingsHash,
		NoWait:                          req.NoWait,
		SetupOptions:                    &setup.Options{BootAfterBuild: req.BootAfterCreate},
		Client:                          builder.NewAPIClient(caller),
	}, nil
}
