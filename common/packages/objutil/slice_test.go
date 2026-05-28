// Copyright 2022-2025 The sacloud/packages-go Authors
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

package objutil

import (
	"reflect"
	"testing"
)

func TestToSlice(t *testing.T) {
	tests := []struct {
		name string
		args interface{}
		want []interface{}
	}{
		{
			name: "string",
			args: "a",
			want: []interface{}{"a"},
		},
		{
			name: "string-slice",
			args: []string{"a", "b"},
			want: []interface{}{"a", "b"},
		},
		{
			name: "int",
			args: 1,
			want: []interface{}{1},
		},
		{
			name: "int-slice",
			args: []int{1, 2},
			want: []interface{}{1, 2},
		},
		{
			name: "struct",
			args: dummy{value: "1"},
			want: []interface{}{dummy{value: "1"}},
		},
		{
			name: "struct-slice",
			args: []dummy{{value: "1"}, {value: "2"}},
			want: []interface{}{dummy{value: "1"}, dummy{value: "2"}},
		},
		{
			name: "pointer",
			args: &dummy{value: "1"},
			want: []interface{}{&dummy{value: "1"}},
		},
		{
			name: "pointer-slice",
			args: []*dummy{{value: "1"}, {value: "2"}},
			want: []interface{}{&dummy{value: "1"}, &dummy{value: "2"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ToSlice(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

type dummy struct {
	value string
}
