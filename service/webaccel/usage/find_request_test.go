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

package usage

import (
	"testing"

	"github.com/sacloud/services/helper"
	"github.com/stretchr/testify/require"
)

func TestFindRequest_validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *FindRequest
		wantErr bool
	}{
		{
			name:    "empty",
			req:     &FindRequest{},
			wantErr: false,
		},
		{
			name:    "returns error with only Year",
			req:     &FindRequest{Year: 1},
			wantErr: true,
		},
		{
			name:    "returns error with only Month",
			req:     &FindRequest{Month: 1},
			wantErr: true,
		},
		{
			name:    "valid",
			req:     &FindRequest{Year: 2020, Month: 1},
			wantErr: false,
		},
		{
			name:    "returns error with invalid month",
			req:     &FindRequest{Year: 2020, Month: -1},
			wantErr: true,
		},
		{
			name:    "returns error with multiple error",
			req:     &FindRequest{Year: 0, Month: -1},
			wantErr: true,
		},
	}

	svc := New(nil)
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := helper.ValidateStruct(svc, tt.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
