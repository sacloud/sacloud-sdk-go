// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/nosql/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/httperror"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

type Error struct {
	msg string
	err error
}

func (e *Error) Error() string {
	if e.msg != "" {
		if e.err != nil {
			return "nosql: " + e.msg + ": " + e.err.Error()
		} else {
			return "nosql: " + e.msg
		}
	} else {
		return "nosql: " + e.err.Error()
	}
}

func (e *Error) Unwrap() error {
	return e.err
}

func NewError(msg string, err error) *Error {
	return &Error{msg: msg, err: err}
}

func NewAPIError(method string, code int, err error) *Error {
	return &Error{msg: method, err: saclient.NewError(code, "", err)}
}

type errorish interface {
	GetStatus() *v1.OptString
	GetErrorMsg() *v1.OptString
}

type apiErrorResponse struct{ errorish }

func (e *apiErrorResponse) Error() string {
	switch {
	case e == nil:
		return "(nil error)"
	case e.errorish == nil:
		return "(nil error)"
	case httperror.IsTypedNil(e.errorish):
		return "(nil error)"
	default:
		opt := e.GetErrorMsg()
		if opt == nil {
			return "(no status code)"
		}
		return opt.Or("(empty error message)")
	}
}

func newErrorFromAPIError(err errorish) error {
	e := apiErrorResponse{errorish: err}

	if err == nil {
		return errors.Join(&e, errors.New("nil API error response"))
	} else if httperror.IsTypedNil(err) {
		return errors.Join(&e, errors.New("nil API error response"))
	} else if opt := err.GetStatus(); opt == nil {
		return errors.Join(&e, errors.New("no status code"))
	} else if s, ok := opt.Get(); !ok {
		return errors.Join(&e, errors.New("no status code"))
	} else if n, _, ok := strings.Cut(s, " "); !ok {
		return errors.Join(&e, fmt.Errorf("invalid HTTP status %q", s))
	} else if c, err := strconv.Atoi(n); err != nil {
		return errors.Join(&e, fmt.Errorf("invalid HTTP status %q: %w", s, err))
	} else {
		return saclient.NewError(c, "", &e)
	}
}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse2[T, U any](method string, yield func() (any, error)) (ret any, err error) {
	res, e1 := yield()

	if e2, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](e1); ok {
		err = NewError(method, httperror.NewError(method, e2))
	} else if e1 != nil {
		err = NewError(method, e1) // this is e.g. broken JSON
	} else if e3, ok := res.(errorish); ok {
		err = NewError(method, newErrorFromAPIError(e3))
	} else if ret, ok = res.(*T); ok {
		return
	} else if ret, ok = res.(*U); ok {
		return
	} else {
		err = NewError(method, fmt.Errorf("unknown API response type %v", res))
	}

	return
}

func errorFromDecodedResponse[T any](method string, yield func() (any, error)) (ret *T, err error) {
	type never struct{}

	if intermediate, err := errorFromDecodedResponse2[T, never](method, yield); err != nil {
		return nil, err
	} else if t, ok := intermediate.(*T); ok {
		return t, nil
	} else {
		return nil, NewError(method, fmt.Errorf("unexpected response type %v", intermediate))
	}
}

func errorFromDecodedResponseE[T any](method string, yield func() (any, error)) (err error) {
	_, err = errorFromDecodedResponse[T](method, yield)
	return
}
