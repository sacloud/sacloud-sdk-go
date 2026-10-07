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

package nfs

import (
	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/search"
	"github.com/sacloud/sacloud-sdk-go/common/packages/objutil"
	"github.com/sacloud/sacloud-sdk-go/common/packages/validate"
	"github.com/sacloud/sacloud-sdk-go/service/iaas/serviceutil"
)

type FindRequest struct {
	Zone string `service:"-" validate:"required"`

	Names []string `service:"-"`
	Tags  []string `service:"-"`

	Sort  search.SortKeys
	Count int
	From  int
}

func (req *FindRequest) Validate() error {
	return validate.New().Struct(req)
}

func (req *FindRequest) ToRequestParameter() (*iaas.FindCondition, error) {
	condition := &iaas.FindCondition{
		Filter: map[search.FilterKey]any{},
	}
	if err := serviceutil.RequestConvertTo(req, condition); err != nil {
		return nil, err
	}

	if !objutil.IsEmpty(req.Names) {
		condition.Filter[search.Key("Name")] = search.AndEqual(req.Names...)
	}
	if !objutil.IsEmpty(req.Tags) {
		condition.Filter[search.Key("Tags.Name")] = search.TagsAndEqual(req.Tags...)
	}
	return condition, nil
}
