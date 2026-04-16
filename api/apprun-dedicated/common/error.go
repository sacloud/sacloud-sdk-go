// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"strings"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

type Error struct {
	msg string
	err error
}

func (e *Error) Unwrap() error { return e.err }
func (e *Error) Error() string {
	var buf strings.Builder

	buf.WriteString("apprun-dedicated")

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

// Wrapper to reduce copy & paste
func ErrorFromDecodedResponseE(m string, y func() error) (e error) {
	type t struct{}
	_, e = ErrorFromDecodedResponse[*t](m, func() (*t, error) { return nil, y() })
	return
}

// saclient.IsNotFoundError() に対応させる
func ErrorFromDecodedResponse[T any](method string, yield func() (T, error)) (resp T, err error) {
	resp, err = yield()

	if e, ok := errors.Into[*v1.ErrorStatusCode](err); ok {
		err = NewAPIError(method, e.GetStatusCode(), e)
	} else if e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err); ok {
		err = NewAPIError(method, e.StatusCode, err)
	} else if err != nil {
		err = NewError(method, err) // :TODO: is this necessary...?
	}

	return
}
