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

package mapconv

import (
	"testing"

	"github.com/fatih/structs"
	"github.com/stretchr/testify/require"
)

func TestMap_Set(t *testing.T) {
	expects := []struct {
		caseName string
		source   map[string]any
		dest     map[string]any
	}{
		{
			caseName: "minimum",
			source: map[string]any{
				"test": "test",
			},
			dest: map[string]any{
				"test": "test",
			},
		},
		{
			caseName: "nested",
			source: map[string]any{
				"test.A":    "A",
				"test.B":    "B",
				"test.C.C1": "C1",
				"test.C.C2": "C2",
				"outer":     "outer",
				"int":       1,
				"float":     1.1,
			},
			dest: map[string]any{
				"test": map[string]any{
					"A": "A",
					"B": "B",
					"C": map[string]any{
						"C1": "C1",
						"C2": "C2",
					},
				},
				"outer": "outer",
				"int":   1,
				"float": 1.1,
			},
		},
		{
			caseName: "slice",
			source: map[string]any{
				"slice.slice.value": []any{"value4", "value5"},
			},
			dest: map[string]any{
				"slice": map[string]any{
					"slice": map[string]any{
						"value": []any{"value4", "value5"},
					},
				},
			},
		},
		{
			caseName: "expanded slice",
			source: map[string]any{
				"[]slice.value": []any{"value1", "value2"},
			},
			dest: map[string]any{
				"slice": []map[string]any{
					{"value": "value1"},
					{"value": "value2"},
				},
			},
		},
		{
			caseName: "expanded nested slice",
			source: map[string]any{
				"[]slice.slice.value": []any{"value4", "value5"},
			},
			dest: map[string]any{
				"slice": []map[string]any{
					{
						"slice": map[string]any{
							"value": "value4",
						},
					},
					{
						"slice": map[string]any{
							"value": "value5",
						},
					},
				},
			},
		},
		{
			caseName: "expanded nested slice with middle slice",
			source: map[string]any{
				"slice.[]slice.value": []any{"value4", "value5"},
			},
			dest: map[string]any{
				"slice": map[string]any{
					"slice": []map[string]any{
						{"value": "value4"},
						{"value": "value5"},
					},
				},
			},
		},
		{
			caseName: "expanded nested slice with last slice",
			source: map[string]any{
				"slice.slice.[]value": []any{"value4", "value5"},
			},
			dest: map[string]any{
				"slice": map[string]any{
					"slice": map[string]any{
						"value": []any{"value4", "value5"},
					},
				},
			},
		},
		{
			caseName: "expanded deep nested slice",
			source: map[string]any{
				"[]slice.[]slice.value": []any{"value4", "value5"},
			},
			dest: map[string]any{
				"slice": []map[string]any{
					{
						"slice": []map[string]any{
							{"value": "value4"},
						},
					},
					{
						"slice": []map[string]any{
							{"value": "value5"},
						},
					},
				},
			},
		},
	}

	for _, expect := range expects {
		t.Run(expect.caseName, func(t *testing.T) {
			m := Map(make(map[string]any))
			for k, v := range expect.source {
				m.Set(k, v)
			}
			require.Equal(t, expect.dest, m.Map())
		})
	}
}

func TestMap_Get(t *testing.T) {
	expects := []struct {
		caseName  string
		keyValues map[string]any
		source    map[string]any
		err       error
	}{
		{
			caseName: "minimum",
			keyValues: map[string]any{
				"test": "test",
			},
			source: map[string]any{
				"test": "test",
			},
		},
		{
			caseName: "nested",
			keyValues: map[string]any{
				"test.A":    "A",
				"test.B":    "B",
				"test.C.C1": "C1",
				"test.C.C2": "C2",
				"outer":     "outer",
				"int":       1,
				"float":     1.1,
			},
			source: map[string]any{
				"test": map[string]any{
					"A": "A",
					"B": "B",
					"C": map[string]any{
						"C1": "C1",
						"C2": "C2",
					},
				},
				"outer": "outer",
				"int":   1,
				"float": 1.1,
			},
		},
		{
			caseName: "slice",
			keyValues: map[string]any{
				"slice.value": []any{"value1", "value2"},
			},
			source: map[string]any{
				"slice": []map[string]any{
					{"value": "value1"},
					{"value": "value2"},
				},
			},
		},
		{
			caseName: "nested slice",
			keyValues: map[string]any{
				"slice.slice.value": []any{"value4", "value5"},
			},
			source: map[string]any{
				"slice": []map[string]any{
					{"value": "value1"},
					{"value": "value2"},
					{
						"value": "value3",
						"slice": []map[string]any{
							{"value": "value4"},
							{"value": "value5"},
						},
					},
				},
			},
		},
		{
			caseName: "with invalid key",
			keyValues: map[string]any{
				"test.A.B": nil,
			},
			source: map[string]any{
				"test": map[string]any{
					"A": "test",
				},
			},
		},
	}

	for _, expect := range expects {
		t.Run(expect.caseName, func(t *testing.T) {
			m := Map(expect.source)
			for k, v := range expect.keyValues {
				value, err := m.Get(k)
				require.Equal(t, expect.err, err)
				if err == nil {
					require.Equal(t, v, value)
				}
			}
		})
	}
}

type A struct {
	Foo *B `structs:",omitempty"`
}

type B struct {
	Bar *C `structs:",omitempty"`
}

type C struct {
	Baz any `structs:",omitempty"`
}

func TestMap_GetWithEmptyStruct(t *testing.T) {
	cases := []struct {
		in     any
		expect any
	}{
		{in: &A{Foo: &B{Bar: &C{Baz: "FooBarBaz"}}}, expect: "FooBarBaz"},
		{in: &A{Foo: &B{Bar: &C{}}}, expect: nil},
		{in: &A{Foo: &B{}}, expect: nil},
	}
	for _, tc := range cases {
		m := Map(structs.Map(tc.in))
		got, err := m.Get("Foo.Bar.Baz")
		if err != nil {
			t.Fatal(err)
		}
		if got != tc.expect {
			t.Fatalf("got unexpected value: expected: %v actual:%v", tc.expect, got)
		}
	}
}
