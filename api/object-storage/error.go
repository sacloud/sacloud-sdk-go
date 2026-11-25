// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"errors"
	"fmt"
	"strings"

	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

type Error struct {
	msg string
	err error
}

func (e *Error) Unwrap() error { return e.err }
func (e *Error) Error() string {
	var buf strings.Builder

	buf.WriteString("object-storage")

	if e.msg != "" {
		buf.WriteString(": ")
		buf.WriteString(e.msg)
	}

	if e.err != nil {
		buf.WriteString(": ")
		buf.WriteString(e.err.Error())
	}

	return buf.String()
}

func NewError(msg string, err error) *Error { return &Error{msg: msg, err: err} }
func NewAPIError(method string, code int, err error) *Error {
	return NewError(method, saclient.NewError(code, "", err))
}

func unknownError(method string, res any) *Error {
	var status int
	var message string

	switch res := res.(type) {
	case *v2.ErrorResponseStatusCode:
		status = res.StatusCode
		message = string(res.Response.Error.Message)

	// v2.ErrorResponse
	case interface{ GetError() v2.ErrorBody }:
		err := res.GetError()
		status = int(err.Code)
		message = string(err.Message)

	default:
		// :FIXME: unintentional info leakage could happen in the error message
		// but don't know what to do otherwise.
		message = fmt.Sprintf("unexpected error type: %#+v", res)
	}

	return NewAPIError(method, status, errors.New(message))
}
