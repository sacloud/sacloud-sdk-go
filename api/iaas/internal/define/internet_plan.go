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
	"github.com/sacloud/iaas-api-go/internal/define/ops"
	"github.com/sacloud/iaas-api-go/internal/dsl"
	"github.com/sacloud/iaas-api-go/internal/dsl/meta"
	"github.com/sacloud/iaas-api-go/naked"
)

const (
	internetPlanAPIName     = "InternetPlan"
	internetPlanAPIPathName = "product/internet"
)

var internetPlanAPI = &dsl.Resource{
	Name:       internetPlanAPIName,
	PathName:   internetPlanAPIPathName,
	PathSuffix: dsl.CloudAPISuffix,
	Operations: dsl.Operations{
		ops.Find(internetPlanAPIName, internetPlanNakedType, findParameter, internetPlanView),
		ops.Read(internetPlanAPIName, internetPlanNakedType, internetPlanView),
	},
}

var (
	internetPlanNakedType = meta.Static(naked.InternetPlan{})
	internetPlanView      = &dsl.Model{
		Name:      internetPlanAPIName,
		NakedType: internetPlanNakedType,
		Fields: []*dsl.FieldDesc{
			fields.ID(),
			fields.Name(),
			fields.Def("BandWidthMbps", meta.TypeInt),
			fields.Availability(),
		},
	}
)
