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
	"github.com/sacloud/iaas-api-go/internal/define/names"
	"github.com/sacloud/iaas-api-go/internal/define/ops"
	"github.com/sacloud/iaas-api-go/internal/dsl"
	"github.com/sacloud/iaas-api-go/internal/dsl/meta"
	"github.com/sacloud/iaas-api-go/naked"
)

const (
	iconAPIName     = "Icon"
	iconAPIPathName = "icon"
)

// iconAPI アイコンAPI
//
// Note: libsacloudでは画像データ取得(GET /icon/:id?Size=[small|medium|large])はサポートしない。
var iconAPI = &dsl.Resource{
	Name:       iconAPIName,
	PathName:   iconAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	IsGlobal:   true,
	Operations: dsl.Operations{
		// find
		ops.Find(iconAPIName, iconNakedType, findParameter, iconView),

		// create
		ops.Create(iconAPIName, iconNakedType, iconCreateParam, iconView),

		// read
		ops.Read(iconAPIName, iconNakedType, iconView),

		// update
		ops.Update(iconAPIName, iconNakedType, iconUpdateParam, iconView),

		// delete
		ops.Delete(iconAPIName),
	},
}

var (
	iconNakedType = meta.Static(naked.Icon{})

	iconView = &dsl.Model{
		Name:      iconAPIName,
		NakedType: iconNakedType,
		Fields: []*dsl.FieldDesc{
			fields.ID(),
			fields.Name(),
			fields.Tags(),
			fields.Availability(),
			fields.Scope(),
			fields.IconURL(),
			fields.CreatedAt(),
			fields.ModifiedAt(),
		},
	}
	iconCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(iconAPIName),
		NakedType: iconNakedType,
		Fields: []*dsl.FieldDesc{
			fields.Name(),
			fields.Tags(),
			fields.IconImage(),
		},
	}

	iconUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(iconAPIName),
		NakedType: iconNakedType,
		Fields: []*dsl.FieldDesc{
			fields.Name(),
			fields.Tags(),
		},
	}
)
