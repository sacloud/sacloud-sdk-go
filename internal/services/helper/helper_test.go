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
	"testing"

	"github.com/sacloud/services/dummy"
	"github.com/stretchr/testify/require"
)

func TestParameterMeta(t *testing.T) {
	svc := dummy.New()
	_, err := ServiceMeta(svc)
	require.NoError(t, err)

	PrintServiceMeta(svc)
}

func TestNewParameter(t *testing.T) {
	parameter, err := NewParameter(dummy.New(), "Find")
	require.NoError(t, err)
	require.NotNil(t, parameter)

	// 手動でnew & Initializeした結果と同等なはず
	want := &dummy.FindRequest{}
	want.Initialize()

	require.EqualValues(t, want, parameter)
}

func TestValidateParameter(t *testing.T) {
	tests := []struct {
		name      string
		parameter interface{}
		wantErr   bool
	}{
		{
			name: "required",
			parameter: &dummy.FindRequest{
				Field1: "",
				Field2: "",
			},
			wantErr: true,
		},
		{
			name: "oneof",
			parameter: &dummy.FindRequest{
				Field1: "dummy",
				Field2: "o3",
			},
			wantErr: false,
		},
		{
			name: "invalid options",
			parameter: &dummy.FindRequest{
				Field1: "dummy",
				Field2: "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		service := dummy.New()
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateStruct(service, tt.parameter); (err != nil) != tt.wantErr {
				t.Errorf("ValidateParameter() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
