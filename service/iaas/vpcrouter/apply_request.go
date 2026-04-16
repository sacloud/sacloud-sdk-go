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

package vpcrouter

import (
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/setup"
	"github.com/sacloud/iaas-service-go/vpcrouter/builder"
	"github.com/sacloud/packages-go/validate"
)

// ApplyRequest Applyサービスへのパラメータ
type ApplyRequest struct {
	Zone string   `service:"-" validate:"required"`
	ID   types.ID `service:"-"`

	Name        string `validate:"required"`
	Description string `validate:"min=0,max=512"`
	Tags        types.Tags
	IconID      types.ID

	PlanID                types.ID `validate:"required"`
	Version               int
	NICSetting            builder.NICSettingHolder             // StandardNICSetting または PremiumNICSetting を指定する
	AdditionalNICSettings []builder.AdditionalNICSettingHolder // AdditionalStandardNICSetting または AdditionalPremiumNICSetting を指定する
	RouterSetting         *RouterSetting
	NoWait                bool
	BootAfterCreate       bool
}

func (req *ApplyRequest) Validate() error {
	return validate.New().Struct(req)
}

// RouterSetting VPCルータの設定
type RouterSetting struct {
	VRID                      int
	InternetConnectionEnabled types.StringFlag
	StaticNAT                 []*iaas.VPCRouterStaticNAT
	PortForwarding            []*iaas.VPCRouterPortForwarding
	Firewall                  []*iaas.VPCRouterFirewall
	DHCPServer                []*iaas.VPCRouterDHCPServer
	DHCPStaticMapping         []*iaas.VPCRouterDHCPStaticMapping
	DNSForwarding             *iaas.VPCRouterDNSForwarding
	PPTPServer                *iaas.VPCRouterPPTPServer
	L2TPIPsecServer           *iaas.VPCRouterL2TPIPsecServer
	WireGuard                 *iaas.VPCRouterWireGuard
	RemoteAccessUsers         []*iaas.VPCRouterRemoteAccessUser
	SiteToSiteIPsecVPN        *iaas.VPCRouterSiteToSiteIPsecVPN
	StaticRoute               []*iaas.VPCRouterStaticRoute
	SyslogHost                string
	ScheduledMaintenance      *iaas.VPCRouterScheduledMaintenance
	MonitoringSuite           *iaas.MonitoringSuite
}

func (req *ApplyRequest) Builder(caller iaas.APICaller) *builder.Builder {
	return &builder.Builder{
		ID:   req.ID,
		Zone: req.Zone,

		Name:                  req.Name,
		Description:           req.Description,
		Tags:                  req.Tags,
		IconID:                req.IconID,
		PlanID:                req.PlanID,
		Version:               req.Version,
		NICSetting:            req.NICSetting,
		AdditionalNICSettings: req.AdditionalNICSettings,
		RouterSetting:         req.routerSetting(),
		NoWait:                req.NoWait,
		Client:                iaas.NewVPCRouterOp(caller),
		SetupOptions: &setup.Options{
			BootAfterBuild: req.BootAfterCreate,
		},
	}
}

func (req *ApplyRequest) routerSetting() *builder.RouterSetting {
	if req.RouterSetting == nil {
		return nil
	}

	return &builder.RouterSetting{
		VRID:                      req.RouterSetting.VRID,
		InternetConnectionEnabled: req.RouterSetting.InternetConnectionEnabled,
		StaticNAT:                 req.RouterSetting.StaticNAT,
		PortForwarding:            req.RouterSetting.PortForwarding,
		Firewall:                  req.RouterSetting.Firewall,
		DHCPServer:                req.RouterSetting.DHCPServer,
		DHCPStaticMapping:         req.RouterSetting.DHCPStaticMapping,
		DNSForwarding:             req.RouterSetting.DNSForwarding,
		PPTPServer:                req.RouterSetting.PPTPServer,
		L2TPIPsecServer:           req.RouterSetting.L2TPIPsecServer,
		WireGuard:                 req.RouterSetting.WireGuard,
		RemoteAccessUsers:         req.RouterSetting.RemoteAccessUsers,
		SiteToSiteIPsecVPN:        req.RouterSetting.SiteToSiteIPsecVPN,
		StaticRoute:               req.RouterSetting.StaticRoute,
		SyslogHost:                req.RouterSetting.SyslogHost,
		ScheduledMaintenance:      req.RouterSetting.ScheduledMaintenance,
		MonitoringSuite:           req.RouterSetting.MonitoringSuite,
	}
}
