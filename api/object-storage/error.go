// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"errors"
	"fmt"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	v2 "github.com/sacloud/sacloud-sdk-go/api/object-storage/apis/v2"
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

type errorish interface{ GetError() v2.ErrorBody }

type apierror struct{ *v2.ErrorBody }

func (e *apierror) Error() string {
	if e == nil {
		return "(nil error)"
	}

	if e.ErrorBody == nil {
		return "(nil error)"
	}

	return string(e.GetMessage())
}

func newErrorFromAPIError(err errorish) error {
	if err == nil {
		return errors.New("nil API error response")
	}
	if httperror.IsTypedNil(err) {
		return errors.New("nil API error response")
	}

	b := err.GetError()
	e := apierror{ErrorBody: &b}
	c := int(b.GetCode())
	return saclient.NewError(c, "", &e)
}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse[T any](method string, yield func() (any, error)) (ret *T, err error) {
	res, e1 := yield()

	if e2, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](e1); ok {
		err = NewError(method, httperror.NewError(method, e2))
	} else if e1 != nil {
		err = NewError(method, e1) // this is e.g. broken JSON
	} else if e3, ok := res.(*v2.ErrorResponseStatusCode); ok {
		status := e3.GetStatusCode()
		message := string(new(new(e3.GetResponse()).GetError()).GetMessage())
		err = NewAPIError(method, status, errors.New(message))
	} else if e3, ok := res.(errorish); ok {
		err = NewError(method, newErrorFromAPIError(e3))
	} else if ret, ok = res.(*T); !ok {
		err = NewError(method, fmt.Errorf("unknown API response type %v", res))
	}

	return
}

func errorFromDecodedResponseE[T any](method string, yield func() (any, error)) (err error) {
	_, err = errorFromDecodedResponse[T](method, yield)
	return
}
