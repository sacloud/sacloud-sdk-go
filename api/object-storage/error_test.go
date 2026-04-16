// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

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
			want: "object-storage: something failed: base error",
		},
		{
			name: "with msg only",
			err:  &Error{msg: "only msg"},
			want: "object-storage: only msg",
		},
		{
			name: "with err only",
			err:  &Error{err: baseErr},
			want: "object-storage: base error",
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
