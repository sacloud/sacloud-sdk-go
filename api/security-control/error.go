// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

//go:generate go run ./gen/list_of_errors.go

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/security-control/apis/v1"
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

	buf.WriteString("security-control")

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

func intoDetail[T ~string, U interface{ Or(T) T }](u U) (ret string) {
	ret = string(u.Or("(empty detail)"))
	if ret == "" {
		ret = "(empty detail)"
	}
	return
}

func newErrorFromAPIError(err any) error {
	status := 0
	detail := "(nil response)"

	switch r := err.(type) {
	case *v1.BadRequest:
		status = http.StatusBadRequest
		if r != nil {
			detail = intoDetail(r.GetDetail())
		}
	case *v1.Conflict:
		status = http.StatusConflict
		if r != nil {
			detail = intoDetail(r.GetDetail())
		}
	case *v1.Forbidden:
		status = http.StatusForbidden
		if r != nil {
			detail = intoDetail(r.GetDetail())
		}
	case *v1.NotFound:
		status = http.StatusNotFound
		if r != nil {
			detail = intoDetail(r.GetDetail())
		}
	case *v1.Unauthorized:
		status = http.StatusUnauthorized
		if r != nil {
			detail = intoDetail(r.GetDetail())
		}
	case *v1.TooManyRequestsHeaders:
		status = http.StatusTooManyRequests
		if r != nil {
			res := r.GetResponse()
			detail = intoDetail(res.GetDetail())
		}
	case *v1.UnexpectedErrorStatusCode:
		if r != nil {
			status = r.GetStatusCode()
			res := r.GetResponse()
			detail = intoDetail(res.GetDetail())
		}
	default:
		status = 0
		switch {
		case r == nil:
			detail = "(nil response)"
		case httperror.IsTypedNil(r):
			detail = "(nil response)"
		default:
			// should we panic ...?
			detail = fmt.Sprintf("unknown API error type %T", err)
		}
	}

	return saclient.NewError(status, "", errors.New(detail))
}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse[T any](method string, yield func() (any, error)) (ret *T, err error) {
	res, e1 := yield()

	if e2, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](e1); ok {
		err = NewError(method, httperror.NewError(method, e2))
	} else if e1 != nil {
		err = NewError(method, e1) // this is e.g. broken JSON
	} else if e3, ok := asErrorResponse(res); ok {
		err = NewError(method, newErrorFromAPIError(e3))
	} else if e4, ok := (res).(*v1.UnexpectedErrorStatusCode); ok {
		err = NewError(method, newErrorFromAPIError(e4))
	} else if ret, ok = res.(*T); !ok {
		err = NewError(method, fmt.Errorf("unknown API response type %T", res))
	}

	return
}

func errorFromDecodedResponseE[T any](method string, yield func() (any, error)) (err error) {
	_, err = errorFromDecodedResponse[T](method, yield)
	return
}
