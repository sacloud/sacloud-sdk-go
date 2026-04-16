// Copyright 2025- The sacloud/workflows-api-go Authors
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

package workflows

import (
	"errors"
	"testing"

	"github.com/sacloud/saclient-go"
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
			want: "workflows: something failed: base error",
		},
		{
			name: "with msg only",
			err:  &Error{msg: "only msg"},
			want: "workflows: only msg",
		},
		{
			name: "with err only",
			err:  &Error{err: baseErr},
			want: "workflows: base error",
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
	baseErr := errors.New("base error")

	err := NewAPIError("msg", 404, baseErr)
	assert.Equal("msg", err.msg)
	assert.True(saclient.IsNotFoundError(err))

	err2 := NewAPIError("msg", 503, nil)
	assert.Equal("msg", err2.msg)
	assert.False(saclient.IsNotFoundError(err2))
}
