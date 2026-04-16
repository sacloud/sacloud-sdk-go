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

package services_test

import (
	"reflect"
	"testing"

	"github.com/sacloud/services"
	"github.com/sacloud/services/dummy"
)

func TestOperations_Categories(t *testing.T) {
	tests := []struct {
		name string
		o    services.Operations
		want []string
	}{
		{
			name: "dummy",
			o:    dummy.New().Operations(),
			want: []string{"Basic", "Category1", "Category2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.o.Categories(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Categories() = %v, want %v", got, tt.want)
			}
		})
	}
}
