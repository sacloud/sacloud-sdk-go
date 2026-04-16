// Copyright 2022-2025 The sacloud/services Authors
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

package dummy

import (
	"github.com/sacloud/services"
	"github.com/sacloud/services/meta"
	"github.com/sacloud/services/validate"
)

var _ services.Service = (*Service)(nil)

type Service struct{}

func (s *Service) Info() *services.Info {
	return &services.Info{
		Name:           "dummy",
		Description:    "Description for Dummy service",
		ParentServices: nil,
	}
}

func (s *Service) Operations() services.Operations {
	return services.Operations{
		{
			Name:          "find",
			OperationType: services.OperationTypeList,
			CategoryName:  "basic",
		},
		{
			Name:          "read",
			OperationType: services.OperationTypeRead,
			CategoryName:  "basic",
		},
		{
			Name:          "error-read",
			OperationType: services.OperationTypeRead,
			CategoryName:  "category1",
		},
		{
			Name:          "echo",
			OperationType: services.OperationTypeRead,
			CategoryName:  "category1",
		},
		{
			Name:          "action",
			OperationType: services.OperationTypeAction,
			CategoryName:  "category2",
		},
		{
			Name:          "error-action",
			OperationType: services.OperationTypeAction,
			CategoryName:  "category2",
		},
	}
}

var optionDefs = []*meta.Option{
	{Key: "option1", Values: []string{"o1", "o2"}},
}

func (s *Service) Config() *services.Config {
	return &services.Config{
		OptionDefs: optionDefs,
	}
}

func (s *Service) Validate(p interface{}) error {
	return validate.New(s).Struct(p)
}

func New() *Service {
	return &Service{}
}
