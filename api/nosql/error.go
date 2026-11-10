// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"encoding/json"
	"errors"
	"fmt"

	v1 "github.com/sacloud/sacloud-sdk-go/api/nosql/apis/v1"
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

type apiErrorResponse struct{ json.Marshaler }

func (e *apiErrorResponse) Error() string {
	if buf, err := e.MarshalJSON(); err != nil {
		return errors.Join(err, fmt.Errorf("%#+v", e)).Error()
	} else {
		return string(buf)
	}
}

func unknownError(res any) error {
	switch v := res.(type) {
	case v1.BadRequestResponse:
		return new(apiErrorResponse{&v})

	case v1.ConflictErrorResponse:
		return new(apiErrorResponse{&v})

	case v1.NotFoundResponse:
		return new(apiErrorResponse{&v})

	case v1.ServerErrorResponse:
		return new(apiErrorResponse{&v})

	case v1.UnauthorizedResponse:
		return new(apiErrorResponse{&v})

	default:
		return fmt.Errorf("unknown error response: %#+v", v)
	}
}
