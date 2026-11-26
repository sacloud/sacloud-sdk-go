// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite

import (
	"errors"
	"fmt"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
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

	buf.WriteString("networking-suite")

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

type apiError struct{ *v1.ApiErrorStatusCode }

func (e *apiError) Error() string {
	if e == nil {
		return "(nil error)"
	}

	if e.ApiErrorStatusCode == nil {
		return "(nil error)"
	}

	return new(e.GetResponse()).GetErrorMsg().Or("(empty message)")
}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse[T any](method string, yield func() (any, error)) (ret *T, err error) {
	res, err := yield()

	if e2, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](err); ok {
		err = NewError(method, httperror.NewError(method, e2))
	} else if e3, ok := res.(*v1.ApiErrorStatusCode); ok {
		err = NewAPIError(method, e3.GetStatusCode(), &apiError{e3})
	} else if err != nil {
		err = NewError(method, err) // this is e.g. broken JSON
	} else if ret, ok = res.(*T); !ok {
		err = NewError(method, fmt.Errorf("unexpected response type %T", res))
	}

	return
}

func errorFromDecodedResponseE[T any](method string, yield func() (any, error)) (err error) {
	_, err = errorFromDecodedResponse[T](method, yield)
	return
}
