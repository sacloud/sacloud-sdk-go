// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

//go:generate go run ./gen/list_of_errors.go

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	v1 "github.com/sacloud/sacloud-sdk-go/api/security-control/apis/v1"
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

func responseError(method string, response any) *Error {
	if r, ok := response.(*v1.UnexpectedErrorStatusCode); ok {
		if r == nil {
			return NewAPIError(method, 0, fmt.Errorf("nil API error response %T", response))
		}
		return responseDetail(method, r.StatusCode, r.Response.Detail.Value, r.Response.Detail.Set)
	}
	err, ok := asErrorResponse(response)
	if !ok {
		return NewAPIError(method, 0, fmt.Errorf("unknown API response type %T", response))
	}
	switch r := err.(type) {
	case *v1.BadRequest:
		if r != nil {
			return responseDetail(method, responseStatus(int(r.Status.Value), r.Status.Set, http.StatusBadRequest), string(r.Detail.Value), r.Detail.Set)
		}
	case *v1.Conflict:
		if r != nil {
			return responseDetail(method, responseStatus(int(r.Status.Value), r.Status.Set, http.StatusConflict), string(r.Detail.Value), r.Detail.Set)
		}
	case *v1.Forbidden:
		if r != nil {
			return responseDetail(method, responseStatus(int(r.Status.Value), r.Status.Set, http.StatusForbidden), string(r.Detail.Value), r.Detail.Set)
		}
	case *v1.NotFound:
		if r != nil {
			return responseDetail(method, responseStatus(int(r.Status.Value), r.Status.Set, http.StatusNotFound), string(r.Detail.Value), r.Detail.Set)
		}
	case *v1.Unauthorized:
		if r != nil {
			return responseDetail(method, responseStatus(int(r.Status.Value), r.Status.Set, http.StatusUnauthorized), string(r.Detail.Value), r.Detail.Set)
		}
	case *v1.TooManyRequestsHeaders:
		if r != nil {
			return responseDetail(method, responseStatus(int(r.Response.Status.Value), r.Response.Status.Set, http.StatusTooManyRequests), string(r.Response.Detail.Value), r.Response.Detail.Set)
		}
	}
	return NewAPIError(method, 0, fmt.Errorf("nil API error response %T", response))
}

func responseStatus(status int, set bool, fallback int) int {
	if set {
		return status
	}
	return fallback
}

func responseDetail(method string, code int, detail string, set bool) *Error {
	switch {
	case !set:
		return NewAPIError(method, code, errors.New("API error response missing detail"))
	case detail == "":
		return NewAPIError(method, code, errors.New("API error response has empty detail"))
	default:
		return NewAPIError(method, code, errors.New(detail))
	}
}
