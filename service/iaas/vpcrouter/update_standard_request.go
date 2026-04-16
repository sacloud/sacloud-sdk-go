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

type UpdateStandardRequest struct {
	Zone string   `service:"-" validate:"required"`
	ID   types.ID `service:"-" validate:"required"`

	Name        *string     `service:",omitempty" validate:"omitempty,min=1"`
	Description *string     `service:",omitempty" validate:"omitempty,max=512"`
	Tags        *types.Tags `service:",omitempty"`
	IconID      *types.ID   `service:",omitempty"`

	AdditionalNICSettings *[]*AdditionalStandardNICSettingUpdate `service:"-"` // Indexが同じものを手動でマージする
	RouterSetting         *RouterSettingUpdate                   `service:",omitempty,recursive"`
	NoWait                bool

	SettingsHash string
}

type AdditionalStandardNICSettingUpdate struct {
	SwitchID       *types.ID `service:",omitempty"`
	IPAddress      *string   `service:",omitempty"`
	NetworkMaskLen *int      `service:",omitempty"`
	Index          int       `service:",omitempty"`
}

func (req *UpdateStandardRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *UpdateStandardRequest) ApplyRequest(ctx context.Context, caller iaas.APICaller) (*ApplyRequest, error) {
	current, err := iaas.NewVPCRouterOp(caller).Read(ctx, req.Zone, req.ID)
	if err != nil {
		return nil, err
	}

	if current.PlanID != types.VPCRouterPlans.Standard {
		return nil, fmt.Errorf("target is not a standard plan: Zone=%s ID=%s", req.Zone, req.ID)
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

		ip := ""
		if len(setting.IPAddress) > 0 {
			ip = setting.IPAddress[0]
		}
		additionalNICs = append(additionalNICs, &vpcRouterBuilder.AdditionalStandardNICSetting{
			SwitchID:       nic.SwitchID,
			IPAddress:      ip,
			NetworkMaskLen: setting.NetworkMaskLen,
			Index:          setting.Index,
		})
	}

	applyRequest := &ApplyRequest{
		Zone:                  req.Zone,
		ID:                    req.ID,
		Name:                  current.Name,
		Description:           current.Description,
		Tags:                  current.Tags,
		IconID:                current.IconID,
		PlanID:                current.PlanID,
		NICSetting:            &vpcRouterBuilder.StandardNICSetting{},
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

func (req *UpdateStandardRequest) mergeAdditionalNICSettings(applyRequest *ApplyRequest) error {
	if req.AdditionalNICSettings != nil {
		var newAdditionalNICs []vpcRouterBuilder.AdditionalNICSettingHolder
		for _, reqNIC := range *req.AdditionalNICSettings {
			var nic vpcRouterBuilder.AdditionalNICSettingHolder
			for _, n := range applyRequest.AdditionalNICSettings {
				if reqNIC.Index == n.(*vpcRouterBuilder.AdditionalStandardNICSetting).Index {
					nic = n
					break
				}
			}
			if nic == nil {
				nic = &vpcRouterBuilder.AdditionalStandardNICSetting{}
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
