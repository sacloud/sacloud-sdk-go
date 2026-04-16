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
	"github.com/sacloud/iaas-api-go/types"
)

const (
	simpleNotificationDestinationAPIName     = "SimpleNotificationDestination"
	simpleNotificationDestinationAPIPathName = "commonserviceitem"
)

var simpleNotificationDestinationAPI = &dsl.Resource{
	Name:       simpleNotificationDestinationAPIName,
	PathName:   simpleNotificationDestinationAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	IsGlobal:   true,
	Operations: dsl.Operations{
		// find
		ops.FindCommonServiceItem(simpleNotificationDestinationAPIName, simpleNotificationDestinationNakedType, findParameter, simpleNotificationDestinationView),

		// create
		ops.CreateCommonServiceItem(simpleNotificationDestinationAPIName, simpleNotificationDestinationNakedType, simpleNotificationDestinationCreateParam, simpleNotificationDestinationView),

		// read
		ops.ReadCommonServiceItem(simpleNotificationDestinationAPIName, simpleNotificationDestinationNakedType, simpleNotificationDestinationView),

		// update
		ops.UpdateCommonServiceItem(simpleNotificationDestinationAPIName, simpleNotificationDestinationNakedType, simpleNotificationDestinationUpdateParam, simpleNotificationDestinationView),
		// updateSettings
		ops.UpdateCommonServiceItemSettings(simpleNotificationDestinationAPIName, simpleNotificationDestinationUpdateSettingsNakedType, simpleNotificationDestinationUpdateSettingsParam, simpleNotificationDestinationView),

		// delete
		ops.Delete(simpleNotificationDestinationAPIName),

		// status
		{
			ResourceName: simpleNotificationDestinationAPIName,
			Name:         "Status",
			PathFormat:   dsl.DefaultPathFormatWithID + "/simplenotification/status",
			Method:       http.MethodGet,
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
			},
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: meta.Static(naked.SimpleNotificationDestinationRunningStatus{}),
				Name: "SimpleNotificationDestination",
			}),
			Results: dsl.Results{
				{
					SourceField: "SimpleNotificationDestination",
					DestField:   simpleNotificationDestinationStatusView.Name,
					IsPlural:    false,
					Model:       simpleNotificationDestinationStatusView,
				},
			},
		},
	},
}

var (
	simpleNotificationDestinationNakedType               = meta.Static(naked.SimpleNotificationDestination{})
	simpleNotificationDestinationUpdateSettingsNakedType = meta.Static(naked.SimpleNotificationDestinationSettingsUpdate{})

	simpleNotificationDestinationView = &dsl.Model{
		Name:      simpleNotificationDestinationAPIName,
		NakedType: simpleNotificationDestinationNakedType,
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
			{
				Name: "Type",
				Type: meta.Static(types.ESimpleNotificationDestinationTypes("")),
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Type",
				},
			},
			{
				Name: "Disabled",
				Type: meta.TypeFlag,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Disabled",
				},
			},
			{
				Name: "Value",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Value",
				},
			},
			fields.SettingsHash(),
		},
	}

	simpleNotificationDestinationCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(simpleNotificationDestinationAPIName),
		NakedType: simpleNotificationDestinationNakedType,
		ConstFields: []*dsl.ConstFieldDesc{
			{
				Name: "Class",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "Provider.Class",
				},
				Value: `"saknoticedestination"`,
			},
			{
				Name:  "ServiceClass",
				Type:  meta.TypeString,
				Value: `"cloud/saknoticedestination/2"`,
			},
		},
		Fields: []*dsl.FieldDesc{
			// common fields
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),

			// settings
			{
				Name: "Type",
				Type: meta.Static(types.ESimpleNotificationDestinationTypes("")),
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Type",
				},
			},
			{
				Name: "Disabled",
				Type: meta.TypeFlag,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Disabled",
				},
			},
			{
				Name: "Value",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Value",
				},
			},
		},
	}

	simpleNotificationDestinationUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(simpleNotificationDestinationAPIName),
		NakedType: simpleNotificationDestinationNakedType,
		Fields: []*dsl.FieldDesc{
			// common fields
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),

			// settings
			{
				Name: "Disabled",
				Type: meta.TypeFlag,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Disabled",
				},
			},
			// settings hash
			fields.SettingsHash(),
		},
	}

	simpleNotificationDestinationUpdateSettingsParam = &dsl.Model{
		Name:      names.UpdateSettingsParameterName(simpleNotificationDestinationAPIName),
		NakedType: simpleNotificationDestinationNakedType,
		Fields: []*dsl.FieldDesc{
			// settings
			{
				Name: "Disabled",
				Type: meta.TypeFlag,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Disabled",
				},
			},
			// settings hash
			fields.SettingsHash(),
		},
	}

	simpleNotificationDestinationStatusView = &dsl.Model{
		Name:      "SimpleNotificationDestinationStatus",
		NakedType: meta.Static(naked.SimpleNotificationDestinationRunningStatus{}),
		Fields: []*dsl.FieldDesc{
			{
				Name: "Disabled",
				Type: meta.TypeFlag,
			},
			{
				Name: "ModifiedAt",
				Type: meta.TypeTime,
			},
		},
	}
)
