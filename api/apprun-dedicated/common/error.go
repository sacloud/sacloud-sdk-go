// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package common

import (
	"errors"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/httperror"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
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

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func ErrorFromDecodedResponse[T any](method string, yield func() (T, error)) (resp T, err error) {
	resp, err = yield()

	if e, ok := errors.AsType[*v1.ErrorStatusCode](err); ok {
		err = NewAPIError(method, e.GetStatusCode(), e)
	} else if e, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](err); ok {
		err = NewError(method, httperror.NewError(method, e))
	} else if err != nil {
		err = NewError(method, err) // this is e.g. broken JSON
	}

	return
}
