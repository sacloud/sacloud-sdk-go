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
	licenseAPIName     = "License"
	licenseAPIPathName = "license"
)

var licenseAPI = &dsl.Resource{
	Name:       licenseAPIName,
	PathName:   licenseAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	IsGlobal:   true,
	Operations: dsl.Operations{
		// find
		ops.Find(licenseAPIName, licenseNakedType, findParameter, licenseView),

		// create
		ops.Create(licenseAPIName, licenseNakedType, licenseCreateParam, licenseView),

		// read
		ops.Read(licenseAPIName, licenseNakedType, licenseView),

		// update
		ops.Update(licenseAPIName, licenseNakedType, licenseUpdateParam, licenseView),

		// delete
		ops.Delete(licenseAPIName),
	},
}

var (
	licenseNakedType = meta.Static(naked.License{})

	licenseView = &dsl.Model{
		Name:      licenseAPIName,
		NakedType: licenseNakedType,
		Fields: []*dsl.FieldDesc{
			fields.ID(),
			fields.Name(),
			fields.LicenseInfoID(),
			fields.LicenseInfoName(),
			fields.CreatedAt(),
			fields.ModifiedAt(),
		},
	}

	licenseCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(licenseAPIName),
		NakedType: licenseNakedType,
		Fields: []*dsl.FieldDesc{
			fields.Name(),
			fields.LicenseInfoID(),
		},
	}

	licenseUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(licenseAPIName),
		NakedType: licenseNakedType,
		Fields: []*dsl.FieldDesc{
			fields.Name(),
		},
	}
)
