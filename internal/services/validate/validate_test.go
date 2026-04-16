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

package validate

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/go-multierror"
	"github.com/stretchr/testify/require"
)

type Foo struct {
	Required string `validate:"required"`
}

func TestValidator_Struct(t *testing.T) {
	err := (&Validator{}).StructWithMultiError(&Foo{})
	require.Error(t, err)
	require.EqualError(t, err.Errors[0], "Required: required")
}

type Bar struct {
	Values []string `validate:"omitempty,dive,my-value"`
	Value  string   `validate:"omitempty,my-value"`
}

func TestValidator_RegisterCollectionValidator(t *testing.T) {
	validator := &Validator{}
	allowedValues := []string{"allowed1", "allowed2"}
	validator.RegisterCollectionValidator("my-value", allowedValues)

	cases := []struct {
		name   string
		target *Bar
		want   *multierror.Error
	}{
		{
			name: "no error",
			target: &Bar{
				Values: []string{"allowed1"},
				Value:  "allowed2",
			},
			want: nil,
		},
		{
			name: "with error from singular alias",
			target: &Bar{
				Value: "invalid",
			},
			want: &multierror.Error{Errors: []error{
				fmt.Errorf("Value: oneof=%s", strings.Join(allowedValues, " ")),
			}},
		},
		{
			name: "with error from plural alias",
			target: &Bar{
				Values: []string{"allowed1", "invalid", "allowed2"},
			},
			want: &multierror.Error{Errors: []error{
				fmt.Errorf("Values[1]: oneof=%s", strings.Join(allowedValues, " ")),
			}},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := validator.StructWithMultiError(tc.target)
			require.Equal(t, tc.want, got)
		})
	}
}
