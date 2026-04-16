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
	"context"
	"fmt"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/serviceutil"
	vpcRouterBuilder "github.com/sacloud/iaas-service-go/vpcrouter/builder"
	"github.com/sacloud/packages-go/validate"
)

type UpdateRequest struct {
	Zone string   `service:"-" validate:"required"`
	ID   types.ID `service:"-" validate:"required"`

	Name        *string     `service:",omitempty" validate:"omitempty,min=1"`
	Description *string     `service:",omitempty" validate:"omitempty,max=512"`
	Tags        *types.Tags `service:",omitempty"`
	IconID      *types.ID   `service:",omitempty"`

	NICSetting            *PremiumNICSettingUpdate              `service:",omitempty,recursive"`
	AdditionalNICSettings *[]*AdditionalPremiumNICSettingUpdate `service:"-"` // Indexが同じものを手動でマージする
	RouterSetting         *RouterSettingUpdate                  `service:",omitempty,recursive"`
	NoWait                bool

	SettingsHash string
}

type PremiumNICSettingUpdate struct {
	IPAddresses      *[]string `service:",omitempty"`
	VirtualIPAddress *string   `service:",omitempty"`
	IPAliases        *[]string `service:",omitempty"`
}

type AdditionalPremiumNICSettingUpdate struct {
	SwitchID         *types.ID `service:",omitempty"`
	IPAddresses      *[]string `service:",omitempty"`
	VirtualIPAddress *string   `service:",omitempty"`
	NetworkMaskLen   *int      `service:",omitempty"`
	Index            int
}

type RouterSettingUpdate struct {
	InternetConnectionEnabled *bool                               `service:",omitempty"`
	StaticNAT                 *[]*iaas.VPCRouterStaticNAT         `service:",omitempty,recursive"`
	PortForwarding            *[]*iaas.VPCRouterPortForwarding    `service:",omitempty,recursive"`
	Firewall                  *[]*iaas.VPCRouterFirewall          `service:",omitempty,recursive"`
	DHCPServer                *[]*iaas.VPCRouterDHCPServer        `service:",omitempty,recursive"`
	DHCPStaticMapping         *[]*iaas.VPCRouterDHCPStaticMapping `service:",omitempty,recursive"`
	PPTPServer                *iaas.VPCRouterPPTPServer           `service:",omitempty,recursive"`
	L2TPIPsecServer           *iaas.VPCRouterL2TPIPsecServer      `service:",omitempty,recursive"`
	WireGuard                 *iaas.VPCRouterWireGuard            `service:",omitempty,recursive"`
	RemoteAccessUsers         *[]*iaas.VPCRouterRemoteAccessUser  `service:",omitempty,recursive"`
	SiteToSiteIPsecVPN        *iaas.VPCRouterSiteToSiteIPsecVPN   `service:",omitempty,recursive"`
	StaticRoute               *[]*iaas.VPCRouterStaticRoute       `service:",omitempty,recursive"`
	SyslogHost                *string                             `service:",omitempty"`
	ScheduledMaintenance      *iaas.VPCRouterScheduledMaintenance `service:",omitempty,recursive"`
	MonitoringSuite           *iaas.MonitoringSuite               `service:",omitempty,recursive"`
}

func (req *UpdateRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateRequest) ApplyRequest(ctx context.Context, caller iaas.APICaller) (*ApplyRequest, error) {
	current, err := iaas.NewVPCRouterOp(caller).Read(ctx, req.Zone, req.ID)
	if err != nil {
		return nil, err
	}

	if current.PlanID == types.VPCRouterPlans.Standard {
		return nil, fmt.Errorf("target is not a premium or higher plan: Zone=%s ID=%s", req.Zone, req.ID)
	}
	if current.Availability != types.Availabilities.Available {
		return nil, fmt.Errorf("target has invalid Availability: Zone=%s ID=%s Availability=%v", req.Zone, req.ID.String(), current.Availability)
	}

	var additionalNICs []vpcRouterBuilder.AdditionalNICSettingHolder
	for _, nic := range current.Interfaces {
		if nic.Index == 0 {
			continue
		}
		var setting *iaas.VPCRouterInterfaceSetting
		for _, s := range current.Settings.Interfaces {
			if s.Index == nic.Index {
				setting = s
				break
			}
		}
		if setting == nil {
			continue
		}

		additionalNICs = append(additionalNICs, &vpcRouterBuilder.AdditionalPremiumNICSetting{
			SwitchID:         nic.SwitchID,
			IPAddresses:      setting.IPAddress,
			VirtualIPAddress: setting.VirtualIPAddress,
			NetworkMaskLen:   setting.NetworkMaskLen,
			Index:            setting.Index,
		})
	}

	var nicSetting *vpcRouterBuilder.PremiumNICSetting
	for _, s := range current.Settings.Interfaces {
		if s.Index == 0 {
			nicSetting = &vpcRouterBuilder.PremiumNICSetting{
				SwitchID:         current.Interfaces[0].SwitchID,
				IPAddresses:      s.IPAddress,
				VirtualIPAddress: s.VirtualIPAddress,
				IPAliases:        s.IPAliases,
			}
			break
		}
	}

	applyRequest := &ApplyRequest{
		Zone:                  req.Zone,
		ID:                    req.ID,
		Name:                  current.Name,
		Description:           current.Description,
		Tags:                  current.Tags,
		IconID:                current.IconID,
		PlanID:                current.PlanID,
		NICSetting:            nicSetting,
		AdditionalNICSettings: additionalNICs,
		RouterSetting: &RouterSetting{
			VRID:                      current.Settings.VRID,
			InternetConnectionEnabled: current.Settings.InternetConnectionEnabled,
			StaticNAT:                 current.Settings.StaticNAT,
			PortForwarding:            current.Settings.PortForwarding,
			Firewall:                  current.Settings.Firewall,
			DHCPServer:                current.Settings.DHCPServer,
			DHCPStaticMapping:         current.Settings.DHCPStaticMapping,
			DNSForwarding:             current.Settings.DNSForwarding,
			PPTPServer:                current.Settings.PPTPServer,
			L2TPIPsecServer:           current.Settings.L2TPIPsecServer,
			RemoteAccessUsers:         current.Settings.RemoteAccessUsers,
			SiteToSiteIPsecVPN:        current.Settings.SiteToSiteIPsecVPN,
			StaticRoute:               current.Settings.StaticRoute,
			SyslogHost:                current.Settings.SyslogHost,
			MonitoringSuite:           current.Settings.MonitoringSuite,
		},
		NoWait: false,
	}

	if err := serviceutil.RequestConvertTo(req, applyRequest); err != nil {
		return nil, err
	}

	// NOTE: AdditionalNICSettingsは配列のインデックスではなく
	//       要素中のIndexフィールドを元にマージする必要があるためここで個別実装する
	if err := req.mergeAdditionalNICSettings(applyRequest); err != nil {
		return nil, err
	}

	return applyRequest, nil
}

func (req *UpdateRequest) mergeAdditionalNICSettings(applyRequest *ApplyRequest) error {
	if req.AdditionalNICSettings != nil {
		var newAdditionalNICs []vpcRouterBuilder.AdditionalNICSettingHolder
		for _, reqNIC := range *req.AdditionalNICSettings {
			var nic vpcRouterBuilder.AdditionalNICSettingHolder
			for _, n := range applyRequest.AdditionalNICSettings {
				if reqNIC.Index == n.(*vpcRouterBuilder.AdditionalPremiumNICSetting).Index {
					nic = n
					break
				}
			}
			if nic == nil {
				nic = &vpcRouterBuilder.AdditionalPremiumNICSetting{}
			}
			if err := serviceutil.RequestConvertTo(reqNIC, nic); err != nil {
				return err
			}
			newAdditionalNICs = append(newAdditionalNICs, nic)
		}
		applyRequest.AdditionalNICSettings = newAdditionalNICs
	}
	return nil
}
