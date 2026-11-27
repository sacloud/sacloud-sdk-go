// Copyright 2025- The sacloud/workflows-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package workflows

//go:generate go run ./gen/list_of_errors.go

import (
	"errors"
	"fmt"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
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

	buf.WriteString("workflows")

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

type errorish interface{ GetMessage() string }

type apierror struct{ errorish }

func (e *apierror) Error() string {
	if e == nil {
		return "(nil error)"
	} else if e.errorish == nil {
		return "(nil errorish)"
	} else if httperror.IsTypedNil(e.errorish) {
		return "(nil errorish)"
	} else if msg := e.GetMessage(); msg != "" {
		return msg
	} else {
		return "(empty error message)"
	}
}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse[T any](method string, yield func() (any, error)) (ret *T, err error) {
	res, e1 := yield()

	if e2, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](e1); ok {
		err = NewError(method, httperror.NewError(method, e2))
	} else if e1 != nil {
		err = NewError(method, e1) // this is e.g. broken JSON
	} else if e2, code, ok := asErrorResponse(res); ok {
		err = NewAPIError(method, code, &apierror{errorish: e2})
	} else if ret, ok = res.(*T); !ok {
		err = NewError(method, fmt.Errorf("unknown API response type %T", res))
	}

	return
}

func errorFromDecodedResponseE[T any](method string, yield func() (any, error)) (err error) {
	_, err = errorFromDecodedResponse[T](method, yield)
	return
}
