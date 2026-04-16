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

package helper

import (
	"context"
	"testing"

	"github.com/sacloud/services"
	"github.com/stretchr/testify/require"
)

func Test_testServiceImpl(t *testing.T) {
	tests := []struct {
		name         string
		svc          services.Service
		wantError    bool
		wantWarnings bool
	}{
		{
			name: "minimum valid service",
			svc: &FakeService{
				info:   &services.Info{Name: "fake", Description: "desc"},
				config: &services.Config{},
				operations: services.Operations{
					{Name: "find", OperationType: services.OperationTypeList},
				},
			},
			wantError: false,
		},
		{
			name:      "Info() and Config() required",
			svc:       &FakeService{},
			wantError: true,
		},
		{
			name: "Operations() required",
			svc: &FakeService{
				info:   &services.Info{Name: "fake", Description: "desc"},
				config: &services.Config{},
			},
			wantError: true,
		},
		{
			name: "Operations().Name required",
			svc: &FakeService{
				info:   &services.Info{Name: "fake", Description: "desc"},
				config: &services.Config{},
				operations: services.Operations{
					{Name: ""},
				},
			},
			wantError: true,
		},
		{
			name: "Operations().OperationType required",
			svc: &FakeService{
				info:   &services.Info{Name: "fake", Description: "desc"},
				config: &services.Config{},
				operations: services.Operations{
					{Name: "fake"},
				},
			},
			wantError: true,
		},
		{
			name: "Operations().Name is duplicated",
			svc: &FakeService{
				info:   &services.Info{Name: "fake", Description: "desc"},
				config: &services.Config{},
				operations: services.Operations{
					{Name: "fake", OperationType: services.OperationTypeList},
					{Name: "fake", OperationType: services.OperationTypeList},
				},
			},
			wantError: true,
		},
		{
			name: "func xxxWithContext required",
			svc: &FakeService{
				info:   &services.Info{Name: "fake", Description: "desc"},
				config: &services.Config{},
				operations: services.Operations{
					{Name: "dummy", OperationType: services.OperationTypeList},
				},
			},
			wantError: true,
		},
		{
			name: "response+error required",
			svc: &FakeService{
				info:   &services.Info{Name: "fake", Description: "desc"},
				config: &services.Config{},
				operations: services.Operations{
					{Name: "invalid-return-values", OperationType: services.OperationTypeList},
				},
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, warn := testServiceImpl(tt.svc)
			require.Equal(t, tt.wantError, err != nil, "got unexpected errors: %s", err)
			require.Equal(t, tt.wantWarnings, warn != nil, "got unexpected warnings: %s", warn)
		})
	}
}

var _ services.Service = (*FakeService)(nil)

type FakeService struct {
	info       *services.Info
	operations services.Operations
	config     *services.Config
}

func (s *FakeService) Info() *services.Info {
	return s.info
}

func (s *FakeService) Operations() services.Operations {
	return s.operations
}

func (s *FakeService) Config() *services.Config {
	return s.config
}

func (s *FakeService) Find(req *FakeRequest) ([]interface{}, error) {
	return s.FindWithContext(context.Background(), req)
}

func (s *FakeService) FindWithContext(ctx context.Context, req *FakeRequest) ([]interface{}, error) {
	return nil, nil
}

type FakeRequest struct {
}

func (s *FakeService) Dummy(req *FakeRequest) ([]interface{}, error) {
	return nil, nil
}

func (s *FakeService) InvalidReturnValues(req *FakeRequest) []interface{} {
	return s.InvalidReturnValuesWithContext(context.Background(), req)
}

func (s *FakeService) InvalidReturnValuesWithContext(ctx context.Context, req *FakeRequest) []interface{} {
	return nil
}
