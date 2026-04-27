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

package iaas

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBuildURL(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		param   map[string]interface{}
		want    string
		wantErr bool
	}{
		{
			name: "basic substitution with multiple keys",
			tmpl: "{{.rootURL}}/{{.zone}}/{{.pathSuffix}}/{{.pathName}}",
			param: map[string]interface{}{
				"rootURL":    "https://api.example.com",
				"zone":       "is1a",
				"pathSuffix": "api/cloud/1.1",
				"pathName":   "server",
			},
			want:    "https://api.example.com/is1a/api/cloud/1.1/server",
			wantErr: false,
		},
		{
			name: "substitution with ID",
			tmpl: "{{.rootURL}}/{{.zone}}/{{.pathSuffix}}/{{.pathName}}/{{.id}}",
			param: map[string]interface{}{
				"rootURL":    "https://api.example.com",
				"zone":       "is1b",
				"pathSuffix": "api/cloud/1.1",
				"pathName":   "archive",
				"id":         12345,
			},
			want:    "https://api.example.com/is1b/api/cloud/1.1/archive/12345",
			wantErr: false,
		},
		{
			name: "numeric values are stringified",
			tmpl: "{{.rootURL}}/{{.accountID}}/{{.year}}/{{.month}}",
			param: map[string]interface{}{
				"rootURL":   "https://api.example.com",
				"accountID": 111111111111,
				"year":      2024,
				"month":     12,
			},
			want:    "https://api.example.com/111111111111/2024/12",
			wantErr: false,
		},
		{
			name: "empty template",
			tmpl: "",
			param: map[string]interface{}{
				"rootURL": "https://api.example.com",
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "undefined variable in template",
			tmpl: "{{.rootURL}}/{{.zone}}/{{.undefinedKey}}",
			param: map[string]interface{}{
				"rootURL": "https://api.example.com",
				"zone":    "is1a",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "same placeholder appears multiple times",
			tmpl: "{{.zone}}/{{.zone}}/{{.zone}}",
			param: map[string]interface{}{
				"zone": "is1a",
			},
			want:    "is1a/is1a/is1a",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := buildURL(tt.tmpl, tt.param)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tt.want, got)
		})
	}
}
