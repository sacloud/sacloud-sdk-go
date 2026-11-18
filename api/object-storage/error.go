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
	var detail interface {
		GetCode() v2.OptErrorCode
		GetMessage() v2.OptErrorMessage
	}
	var ok bool

	switch res := res.(type) {
	case *v2.ErrorDefaultStatusCode:
		status = res.GetStatusCode()
		if err, ok := new(res.GetResponse()).GetError().Get(); ok {
			if msg, ok := err.GetMessage().Get(); ok {
				message = string(msg)
			}
		}
		return NewAPIError(method, status, errors.New(message))

	case *v2.Error400:
		d, s := res.GetError().Get()
		detail, ok = &d, s
	case *v2.Error401:
		d, s := res.GetError().Get()
		detail, ok = &d, s
	case *v2.Error403:
		d, s := res.GetError().Get()
		detail, ok = &d, s
	case *v2.Error404:
		d, s := res.GetError().Get()
		detail, ok = &d, s
	case *v2.Error405:
		d, s := res.GetError().Get()
		detail, ok = &d, s
	case *v2.Error409:
		d, s := res.GetError().Get()
		detail, ok = &d, s
	case *v2.ErrorDefault:
		d, s := res.GetError().Get()
		detail, ok = &d, s

	default:
		// :FIXME: unintentional info leakage could happen in the error message
		// but don't know what to do otherwise.
		return NewAPIError(method, 0, fmt.Errorf("unexpected error type: %#+v", res))
	}

	if !ok {
		return NewAPIError(method, 0, fmt.Errorf("got null: %#+v", res))
	}

	if code, ok := detail.GetCode().Get(); ok {
		status = int(code)
	}

	if msg, ok := detail.GetMessage().Get(); ok {
		message = string(msg)
	}

	return NewAPIError(method, status, errors.New(message))
}
