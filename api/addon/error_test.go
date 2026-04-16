// Copyright 2025- The sacloud/addon-api-go Authors
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

package addon

import (
	"errors"
	"testing"

	v1 "github.com/sacloud/addon-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestError_Error(t *testing.T) {
	baseErr := errors.New("base error")

	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "with msg and err",
			err:  &Error{msg: "something failed", err: baseErr},
			want: "addon: something failed: base error",
		},
		{
			name: "with msg only",
			err:  &Error{msg: "only msg"},
			want: "addon: only msg",
		},
		{
			name: "with err only",
			err:  &Error{err: baseErr},
			want: "addon: base error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.New(t).Equal(tt.want, tt.err.Error())
		})
	}
}

func TestNewError(t *testing.T) {
	assert := require.New(t)
	baseErr := errors.New("base error")

	err := NewError("msg", baseErr)
	assert.Equal("msg", err.msg)
	assert.Equal(baseErr, err.err)

	err2 := NewError("msg only", nil)
	assert.Equal("msg only", err2.msg)
	assert.Nil(err2.err)
}

func TestNewAPIError(t *testing.T) {
	assert := require.New(t)

	inner := v1.ErrorResponse{
		Errors: v1.NewOptNilErrorInfoArray([]v1.ErrorInfo{
			{
				Code:    v1.NewOptNilString("InvalidParameter"),
				Message: v1.NewOptNilString("The parameter is invalid."),
			},
		}),
	}
	err := NewAPIError("TestMethod", &inner)

	assert.Equal("TestMethod", err.msg)
	assert.Contains(err.Error(), "The parameter is invalid.")
}
