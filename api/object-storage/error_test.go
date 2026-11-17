// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
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

func TestUnknownError(t *testing.T) {
	tests := []struct {
		response any
		code     int
	}{
		{response: &v2.Error400{}, code: 400},
		{response: &v2.Error401{}, code: 401},
		{response: &v2.Error403{}, code: 403},
		{response: &v2.Error404{}, code: 404},
		{response: &v2.Error405{}, code: 405},
		{response: &v2.Error409{}, code: 409},
		{response: &v2.ErrorDefault{}, code: 500},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%T", tt.response), func(t *testing.T) {
			body := fmt.Sprintf(`{"error":{"code":%d,"message":"API request failed"}}`, tt.code)
			require.NoError(t, json.Unmarshal([]byte(body), tt.response))

			err := unknownError("test", tt.response)
			require.Equal(t, fmt.Sprintf("object-storage: test: API Error %d: API request failed", tt.code), err.Error())
			require.Equal(t, tt.code == 404, saclient.IsNotFoundError(err))
			require.Equal(t, tt.code == 409, saclient.IsConflictError(err))
		})
	}
}

func TestUnknownErrorOptionalFields(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "missing code", body: `{"error":{"message":"unauthorized"}}`, want: "API Error: unauthorized"},
		{name: "missing message", body: `{"error":{"code":401}}`, want: "API Error 401: "},
		{name: "empty detail", body: `{"error":{}}`, want: "API Error: "},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var response v2.Error401
			require.NoError(t, json.Unmarshal([]byte(tt.body), &response))
			require.Equal(t, "object-storage: test: "+tt.want, unknownError("test", &response).Error())
		})
	}
}

func TestUnknownErrorMissingDetail(t *testing.T) {
	require.Contains(t, unknownError("test", &v2.Error401{}).Error(), "object-storage: test: API Error: got null:")
}

func TestUnknownErrorDefaultStatusCode(t *testing.T) {
	response := &v2.ErrorDefaultStatusCode{StatusCode: 503}
	require.NoError(t, json.Unmarshal([]byte(`{"error":{"code":500,"message":"unavailable"}}`), &response.Response))
	require.Equal(t, "object-storage: test: API Error 503: unavailable", unknownError("test", response).Error())

	response.Response = v2.ErrorDefault{}
	require.Equal(t, "object-storage: test: API Error 503: ", unknownError("test", response).Error())
}

func TestUnknownErrorUnexpectedResponse(t *testing.T) {
	require.Contains(t, unknownError("test", struct{}{}).Error(), "object-storage: test: API Error: unexpected error type:")
}
