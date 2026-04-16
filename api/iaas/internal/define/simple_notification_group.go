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
	simpleNotificationGroupAPIName     = "SimpleNotificationGroup"
	simpleNotificationGroupAPIPathName = "commonserviceitem"
)

var simpleNotificationGroupAPI = &dsl.Resource{
	Name:       simpleNotificationGroupAPIName,
	PathName:   simpleNotificationGroupAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	IsGlobal:   true,
	Operations: dsl.Operations{
		// find
		ops.FindCommonServiceItem(simpleNotificationGroupAPIName, simpleNotificationGroupNakedType, findParameter, simpleNotificationGroupView),

		// create
		ops.CreateCommonServiceItem(simpleNotificationGroupAPIName, simpleNotificationGroupNakedType, simpleNotificationGroupCreateParam, simpleNotificationGroupView),

		// read
		ops.ReadCommonServiceItem(simpleNotificationGroupAPIName, simpleNotificationGroupNakedType, simpleNotificationGroupView),

		// update
		ops.UpdateCommonServiceItem(simpleNotificationGroupAPIName, simpleNotificationGroupNakedType, simpleNotificationGroupUpdateParam, simpleNotificationGroupView),
		// updateSettings
		ops.UpdateCommonServiceItemSettings(simpleNotificationGroupAPIName, simpleNotificationGroupUpdateSettingsNakedType, simpleNotificationGroupUpdateSettingsParam, simpleNotificationGroupView),

		// delete
		ops.Delete(simpleNotificationGroupAPIName),

		// Post Message
		{
			ResourceName: simpleNotificationGroupAPIName,
			Name:         "PostMessage",
			PathFormat:   dsl.IDAndSuffixPathFormat("simplenotification/message"),
			Method:       http.MethodPost,
			RequestEnvelope: dsl.RequestEnvelope(
				&dsl.EnvelopePayloadDesc{
					Type: meta.TypeString,
					Name: "Message",
				},
			),
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
				&dsl.Argument{Name: "message", Type: meta.TypeString, MapConvTag: "Message"},
			},
		},

		// history
		{
			ResourceName: simpleNotificationGroupAPIName,
			Name:         "History",
			PathFormat:   dsl.DefaultPathFormat + "/simplenotification/history",
			Method:       http.MethodGet,
			Arguments:    dsl.Arguments{},
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: meta.Static(naked.SimpleNotificationHistories{}),
				Name: "NotificationHistories",
			}),
			Results: dsl.Results{
				{
					SourceField: "NotificationHistories",
					DestField:   simpleNotificationHistory.Name,
					Model:       simpleNotificationHistory,
				},
			},
		},
	},
}

var (
	simpleNotificationGroupNakedType               = meta.Static(naked.SimpleNotificationGroup{})
	simpleNotificationGroupUpdateSettingsNakedType = meta.Static(naked.SimpleNotificationGroupSettingsUpdate{})

	simpleNotificationGroupView = &dsl.Model{
		Name:      simpleNotificationGroupAPIName,
		NakedType: simpleNotificationGroupNakedType,
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
				Name: "Destinations",
				Type: meta.TypeStringSlice,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Destinations",
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
				Name: "Sources",
				Type: meta.TypeStringSlice,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Sources",
				},
			},
			fields.SettingsHash(),
		},
	}

	simpleNotificationGroupCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(simpleNotificationGroupAPIName),
		NakedType: simpleNotificationGroupNakedType,
		ConstFields: []*dsl.ConstFieldDesc{
			{
				Name: "Class",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "Provider.Class",
				},
				Value: `"saknoticegroup"`,
			},
			{
				Name:  "ServiceClass",
				Type:  meta.TypeString,
				Value: `"cloud/saknoticegroup/2"`,
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
				Name: "Destinations",
				Type: meta.TypeStringSlice,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Destinations",
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
				Name: "Sources",
				Type: meta.TypeStringSlice,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Sources",
				},
			},
		},
	}

	simpleNotificationGroupUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(simpleNotificationGroupAPIName),
		NakedType: simpleNotificationGroupNakedType,
		Fields: []*dsl.FieldDesc{
			// common fields
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),

			// settings
			{
				Name: "Destinations",
				Type: meta.TypeStringSlice,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Destinations",
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
				Name: "Sources",
				Type: meta.TypeStringSlice,
				Tags: &dsl.FieldTags{
					MapConv: "Settings.Sources",
				},
			},

			// settings hash
			fields.SettingsHash(),
		},
	}

	simpleNotificationGroupUpdateSettingsParam = &dsl.Model{
		Name:      names.UpdateSettingsParameterName(simpleNotificationGroupAPIName),
		NakedType: simpleNotificationGroupNakedType,
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

	simpleNotificationHistory = &dsl.Model{
		Name:      "SimpleNotificationHistories",
		NakedType: meta.Static(naked.SimpleNotificationHistories{}),
		Fields: []*dsl.FieldDesc{
			{
				Name: "NotificationHistories",
				Type: &dsl.Model{
					Name:      "SimpleNotificationHistory",
					NakedType: meta.Static(naked.SimpleNotificationHistory{}),
					IsArray:   true,
					Fields: []*dsl.FieldDesc{
						fields.Def("RequestID", meta.TypeString),
						fields.Def("SourceID", meta.TypeString),
						fields.Def("ReceivedAt", meta.TypeTime),
						{
							Name: "Message",
							Type: &dsl.Model{
								Name:      "SimpleNotificationHistoryMessage",
								NakedType: meta.Static(naked.SimpleNotificationHistoryMessage{}),
								Fields: []*dsl.FieldDesc{
									fields.Def("Body", meta.TypeString),
									fields.Def("Color", meta.TypeString),
									fields.Def("ColorCode", meta.TypeString),
									fields.Def("IconURL", meta.TypeString),
									fields.Def("ImageURL", meta.TypeString),
									fields.Def("Title", meta.TypeString),
								},
							},
						},
						{
							Name: "Statuses",
							Type: &dsl.Model{
								Name:      "SimpleNotificationHistoryStatus",
								NakedType: meta.Static(naked.SimpleNotificationHistoryStatus{}),
								IsArray:   true,
								Fields: []*dsl.FieldDesc{
									fields.Def("ID", meta.TypeString),
									fields.Def("Status", meta.TypeInt),
									fields.Def("ErrorInfo", meta.TypeString),
									fields.Def("NotificationRequestID", meta.TypeString),
									fields.Def("GroupID", meta.TypeString),
									fields.Def("CreatedAt", meta.TypeTime),
									fields.Def("UpdatedAt", meta.TypeTime),
								},
							},
						},
					},
				},
			},
		},
	}
)
