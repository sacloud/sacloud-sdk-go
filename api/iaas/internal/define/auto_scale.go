// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package define

import (
	"net/http"

	"github.com/sacloud/iaas-api-go/internal/define/names"
	"github.com/sacloud/iaas-api-go/internal/define/ops"
	"github.com/sacloud/iaas-api-go/internal/dsl"
	"github.com/sacloud/iaas-api-go/internal/dsl/meta"
	"github.com/sacloud/iaas-api-go/naked"
)

const (
	autoScaleAPIName     = "AutoScale"
	autoScaleAPIPathName = "commonserviceitem"
)

var autoScaleAPI = &dsl.Resource{
	Name:       autoScaleAPIName,
	PathName:   autoScaleAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	IsGlobal:   true,
	Operations: dsl.Operations{
		// find
		ops.FindCommonServiceItem(autoScaleAPIName, autoScaleNakedType, findParameter, autoScaleView),

		// create
		ops.CreateCommonServiceItem(autoScaleAPIName, autoScaleNakedType, autoScaleCreateParam, autoScaleView),

		// read
		ops.ReadCommonServiceItem(autoScaleAPIName, autoScaleNakedType, autoScaleView),

		// update
		ops.UpdateCommonServiceItem(autoScaleAPIName, autoScaleNakedType, autoScaleUpdateParam, autoScaleView),
		// updateSettings
		ops.UpdateCommonServiceItemSettings(autoScaleAPIName, autoScaleUpdateSettingsNakedType, autoScaleUpdateSettingsParam, autoScaleView),

		// delete
		ops.Delete(autoScaleAPIName),

		// status
		{
			ResourceName: autoScaleAPIName,
			Name:         "Status",
			PathFormat:   dsl.DefaultPathFormatWithID + "/autoscale/status",
			Method:       http.MethodGet,
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
			},
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: meta.Static(naked.AutoScaleRunningStatus{}),
				Name: "AutoScale",
			}),
			Results: dsl.Results{
				{
					SourceField: "AutoScale",
					DestField:   autoScaleStatusView.Name,
					IsPlural:    false,
					Model:       autoScaleStatusView,
				},
			},
		},

		// scale up
		{
			ResourceName: autoScaleAPIName,
			Name:         "ScaleUp",
			Method:       http.MethodPut,
			PathFormat:   dsl.DefaultPathFormatWithID + "/autoscale/up",
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
			},
		},
		// scale down
		{
			ResourceName: autoScaleAPIName,
			Name:         "ScaleDown",
			Method:       http.MethodPut,
			PathFormat:   dsl.DefaultPathFormatWithID + "/autoscale/down",
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
			},
		},
	},
}

var (
	autoScaleNakedType               = meta.Static(naked.AutoScale{})
	autoScaleUpdateSettingsNakedType = meta.Static(naked.AutoScaleSettingsUpdate{})

	autoScaleView = &dsl.Model{
		Name:      autoScaleAPIName,
		NakedType: autoScaleNakedType,
		Fields: []*dsl.FieldDesc{
			fields.ID(),
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.Availability(),
			fields.IconID(),
			fields.CreatedAt(),
			fields.ModifiedAt(),

			// settings
			fields.AutoScaleDisabled(),
			fields.AutoScaleZones(),
			fields.AutoScaleConfig(),
			fields.AutoScaleTriggerType(),
			fields.AutoScaleCPUThreshold(),
			fields.AutoScaleRouterThreshold(),
			fields.AutoScaleScheduleScaling(),
			fields.SettingsHash(),

			// status
			fields.AutoScaleAPIKeyID(),
		},
	}

	autoScaleCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(autoScaleAPIName),
		NakedType: autoScaleNakedType,
		ConstFields: []*dsl.ConstFieldDesc{
			{
				Name: "Class",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "Provider.Class",
				},
				Value: `"autoscale"`,
			},
			{
				Name:  "ServiceClass",
				Type:  meta.TypeString,
				Value: `"cloud/autoscale/1"`,
			},
		},
		Fields: []*dsl.FieldDesc{
			// common fields
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),

			// settings
			fields.AutoScaleDisabled(),
			fields.AutoScaleZones(),
			fields.AutoScaleConfig(),
			fields.AutoScaleTriggerType(),
			fields.AutoScaleCPUThreshold(),
			fields.AutoScaleRouterThreshold(),
			fields.AutoScaleScheduleScaling(),

			// status
			fields.AutoScaleAPIKeyID(),
		},
	}

	autoScaleUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(autoScaleAPIName),
		NakedType: autoScaleNakedType,
		Fields: []*dsl.FieldDesc{
			// common fields
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),

			// settings
			fields.AutoScaleDisabled(),
			fields.AutoScaleZones(),
			fields.AutoScaleConfig(),
			fields.AutoScaleTriggerType(),
			fields.AutoScaleCPUThreshold(),
			fields.AutoScaleRouterThreshold(),
			fields.AutoScaleScheduleScaling(),
			// settings hash
			fields.SettingsHash(),
		},
	}

	autoScaleUpdateSettingsParam = &dsl.Model{
		Name:      names.UpdateSettingsParameterName(autoScaleAPIName),
		NakedType: autoScaleNakedType,
		Fields: []*dsl.FieldDesc{
			// settings
			fields.AutoScaleDisabled(),
			fields.AutoScaleZones(),
			fields.AutoScaleConfig(),
			fields.AutoScaleTriggerType(),
			fields.AutoScaleCPUThreshold(),
			fields.AutoScaleRouterThreshold(),
			fields.AutoScaleScheduleScaling(),
			// settings hash
			fields.SettingsHash(),
		},
	}

	autoScaleStatusView = &dsl.Model{
		Name:      "AutoScaleStatus",
		NakedType: meta.Static(naked.AutoScaleRunningStatus{}),
		Fields: []*dsl.FieldDesc{
			{
				Name: "LatestLogs",
				Type: meta.TypeStringSlice,
			},
			{
				Name: "ResourcesText",
				Type: meta.TypeString,
			},
		},
	}
)
