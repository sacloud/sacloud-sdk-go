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
	sshKeyAPIName     = "SSHKey"
	sshKeyAPIPathName = "sshkey"
)

var sshKeyAPI = &dsl.Resource{
	Name:       sshKeyAPIName,
	PathName:   sshKeyAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	IsGlobal:   true,
	Operations: dsl.Operations{
		// find
		ops.Find(sshKeyAPIName, sshKeyNakedType, findParameter, sshKeyView),

		// create
		ops.Create(sshKeyAPIName, sshKeyNakedType, sshKeyCreateParam, sshKeyView),

		// read
		ops.Read(sshKeyAPIName, sshKeyNakedType, sshKeyView),

		// update
		ops.Update(sshKeyAPIName, sshKeyNakedType, sshKeyUpdateParam, sshKeyView),

		// delete
		ops.Delete(sshKeyAPIName),
	},
}

var (
	sshKeyNakedType = meta.Static(naked.SSHKey{})

	sshKeyFields = []*dsl.FieldDesc{
		fields.ID(),
		fields.Name(),
		fields.Description(),
		fields.CreatedAt(),
		fields.PublicKey(),
		fields.Fingerprint(),
	}

	sshKeyView = &dsl.Model{
		Name:      sshKeyAPIName,
		NakedType: sshKeyNakedType,
		Fields:    sshKeyFields,
	}

	sshKeyCreateParam = &dsl.Model{
		Name:      names.CreateParameterName(sshKeyAPIName),
		NakedType: sshKeyNakedType,
		Fields: []*dsl.FieldDesc{
			fields.Name(),
			fields.Description(),
			fields.PublicKey(),
		},
	}

	sshKeyUpdateParam = &dsl.Model{
		Name:      names.UpdateParameterName(sshKeyAPIName),
		NakedType: sshKeyNakedType,
		Fields: []*dsl.FieldDesc{
			fields.Name(),
			fields.Description(),
		},
	}
)
