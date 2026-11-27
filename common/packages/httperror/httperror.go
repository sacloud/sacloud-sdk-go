// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package httperror

import (
	"errors"
	"net/http/httputil"
	"reflect"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

// Error type that represents ogen's UnexpectedStatusCodeError.
type Error struct {
	*ogen.UnexpectedStatusCodeError
}

func (e *Error) Unwrap() error { return e.UnexpectedStatusCodeError }

func (e *Error) Error() string {
	var buf strings.Builder

	buf.WriteString(e.UnexpectedStatusCodeError.Error())

	payload := e.Payload
	if payload == nil {
		buf.WriteString(" (nil response payload)")
		return buf.String()
	}

	res, err := httputil.DumpResponse(payload, true)
	if err != nil {
		buf.WriteString(" (failed to dump response: ")
		buf.WriteString(err.Error())
		buf.WriteString(")")
		return buf.String()
	}

	buf.WriteString(": ")
	buf.Write(res)
	return buf.String()
}

func NewError(method string, unexpected *ogen.UnexpectedStatusCodeError) error {
	if unexpected == nil {
		return errors.New("unexpected nil status code error")
	}

	code := unexpected.StatusCode
	err := Error{UnexpectedStatusCodeError: unexpected}

	return saclient.NewError(code, method, &err)
}

func IsTypedNil[T any](x T) (ret bool) {
	defer func() {
		if r := recover(); r != nil {
			ret = false
		}
	}()
	ret = reflect.ValueOf(x).IsNil()
	return
}
