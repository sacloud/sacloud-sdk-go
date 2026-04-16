// Copyright 2022-2025 The sacloud/webaccel-service-go authors
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

package site

import (
	"github.com/sacloud/services"
	"github.com/sacloud/webaccel-api-go"
)

// Service provides a high-level API of for Service
type Service struct {
	client *webaccel.Client
}

var _ services.Service = (*Service)(nil)

// New returns new site instance of Service
func New(client *webaccel.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Info() *services.Info {
	return &services.Info{
		Name: "site",
	}
}

func (s *Service) Operations() services.Operations {
	return []services.SupportedOperation{
		{Name: "find", OperationType: services.OperationTypeList},
		{Name: "read", OperationType: services.OperationTypeRead},
		{Name: "update", OperationType: services.OperationTypeUpdate},
	}
}

func (s *Service) Config() *services.Config {
	return &services.Config{}
}
