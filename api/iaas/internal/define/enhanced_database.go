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
	enhancedDatabaseAPIName     = "EnhancedDB"
	enhancedDatabaseAPIPathName = "commonserviceitem"
)

var enhancedDatabaseAPI = &dsl.Resource{
	Name:       enhancedDatabaseAPIName,
	PathName:   enhancedDatabaseAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	IsGlobal:   true,
	Operations: dsl.Operations{
		// find
		ops.FindCommonServiceItem(enhancedDatabaseAPIName, enhancedDatabaseNakedType, findParameter, enhancedDatabaseView),

		// create
		ops.CreateCommonServiceItem(enhancedDatabaseAPIName, enhancedDatabaseNakedType, enhancedDatabaseCreateParam, enhancedDatabaseView),

		// read
		ops.ReadCommonServiceItem(enhancedDatabaseAPIName, enhancedDatabaseNakedType, enhancedDatabaseView),

		// update
		ops.UpdateCommonServiceItem(enhancedDatabaseAPIName, enhancedDatabaseNakedType, enhancedDatabaseUpdateParam, enhancedDatabaseView),

		// delete
		ops.Delete(enhancedDatabaseAPIName),

		// Set Password
		{
			ResourceName: enhancedDatabaseAPIName,
			Name:         "SetPassword",
			PathFormat:   dsl.IDAndSuffixPathFormat("enhanceddb/set-password"),
			Method:       http.MethodPut,
			RequestEnvelope: dsl.RequestEnvelope(&dsl.EnvelopePayloadDesc{
				Type: enhancedDatabaseUserNakedType,
				Name: "CommonServiceItem",
			}),
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
				dsl.MappableArgument("param", enhancedDatabaseSetPasswordParam, "CommonServiceItem"),
			},
		},

		// Get Config
		{
			ResourceName: enhancedDatabaseAPIName,
			Name:         "GetConfig",
			PathFormat:   dsl.IDAndSuffixPathFormat("enhanceddb/config"),
			Method:       http.MethodGet,
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
			},
			ResponseEnvelope: dsl.ResponseEnvelope(&dsl.EnvelopePayloadDesc{
				Type: meta.Static(naked.EnhancedDBConfig{}),
				Name: "EnhancedDB",
			}),
			Results: dsl.Results{
				{
					SourceField: "EnhancedDB",
					DestField:   enhancedDatabaseConfig.Name,
					IsPlural:    false,
					Model:       enhancedDatabaseConfig,
				},
			},
		},

		// Set Config
		{
			ResourceName: enhancedDatabaseAPIName,
			Name:         "SetConfig",
			PathFormat:   dsl.IDAndSuffixPathFormat("enhanceddb/config"),
			Method:       http.MethodPut,
			RequestEnvelope: dsl.RequestEnvelope(&dsl.EnvelopePayloadDesc{
				Type: enhancedDatabaseConfigNakedType,
				Name: "CommonServiceItem",
			}),
			Arguments: dsl.Arguments{
				dsl.ArgumentID,
				dsl.MappableArgument("param", enhancedDatabaseSetConfigParam, "CommonServiceItem"),
			},
		},
	},
}

var (
	enhancedDatabaseNakedType       = meta.Static(naked.EnhancedDB{})
	enhancedDatabaseUserNakedType   = meta.Static(naked.EnhancedDBPasswordSettings{})
	enhancedDatabaseConfigNakedType = meta.Static(naked.EnhancedDBConfigSettings{})

	enhancedDatabaseView = &dsl.Model{
		Name:      enhancedDatabaseAPIName,
		NakedType: enhancedDatabaseNakedType,
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
			fields.SettingsHash(),

			// status
			fields.EnhancedDBDatabaseName(),
			fields.EnhancedDBDatabaseType(),
			fields.EnhancedDBDatabaseRegion(),
			fields.EnhancedDBDatabaseHostName(),
			fields.EnhancedDBDatabasePort(),
		},
	}

	enhancedDatabaseCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(enhancedDatabaseAPIName),
		NakedType: enhancedDatabaseNakedType,
		ConstFields: []*dsl.ConstFieldDesc{
			{
				Name: "Class",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "Provider.Class",
				},
				Value: `"enhanceddb"`,
			},
			{
				Name: "MaxConnections",
				Type: meta.TypeInt,
				Tags: &dsl.FieldTags{
					MapConv: "Config.MaxConnections",
				},
				Value: `50`,
			},
			{
				Name: "AllowedNetworks",
				Type: meta.TypeStringSlice,
				Tags: &dsl.FieldTags{
					MapConv: "Config.AllowedNetworks",
				},
				Value: `[]string{}`,
			},
		},
		Fields: []*dsl.FieldDesc{
			// common fields
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),

			// settings
			fields.EnhancedDBDatabaseName(),
			fields.EnhancedDBDatabaseType(),
			fields.EnhancedDBDatabaseRegion(),
		},
	}

	enhancedDatabaseUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(enhancedDatabaseAPIName),
		NakedType: enhancedDatabaseNakedType,
		Fields: []*dsl.FieldDesc{
			// common fields
			fields.Name(),
			fields.Description(),
			fields.Tags(),
			fields.IconID(),

			// settings hash
			fields.SettingsHash(),
		},
	}

	enhancedDatabaseSetPasswordParam = &dsl.Model{
		Name:      enhancedDatabaseAPIName + "SetPasswordRequest",
		NakedType: meta.Static(naked.EnhancedDBPasswordSetting{}),
		Fields: []*dsl.FieldDesc{
			{
				Name: "Password",
				Type: meta.TypeString,
				Tags: &dsl.FieldTags{
					MapConv: "EnhancedDB.Password",
				},
			},
		},
	}

	enhancedDatabaseConfig = &dsl.Model{
		Name:      enhancedDatabaseAPIName + "Config",
		NakedType: meta.Static(naked.EnhancedDBConfig{}),
		Fields: []*dsl.FieldDesc{
			{
				Name: "MaxConnections",
				Type: meta.TypeInt,
			},
			{
				Name: "AllowedNetworks",
				Type: meta.TypeStringSlice,
			},
		},
	}

	enhancedDatabaseSetConfigParam = &dsl.Model{
		Name:      enhancedDatabaseAPIName + "SetConfigRequest",
		NakedType: meta.Static(naked.EnhancedDBConfigSettings{}),
		Fields: []*dsl.FieldDesc{
			{
				Name: "AllowedNetworks",
				Type: meta.TypeStringSlice,
				Tags: &dsl.FieldTags{
					MapConv: "EnhancedDB.AllowedNetworks",
				},
			},
		},
		ConstFields: []*dsl.ConstFieldDesc{
			{
				Name: "MaxConnections",
				Type: meta.TypeInt,
				Tags: &dsl.FieldTags{
					MapConv: "EnhancedDB.MaxConnections",
				},
				Value: `50`,
			},
		},
	}
)
