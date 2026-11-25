// Copyright 2025- The sacloud/iam-api-go authors
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

package common

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
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

	buf.WriteString("iam")

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

type apiErrorResponse struct{ json.Marshaler }

func (e *apiErrorResponse) Error() string {
	if buf, err := e.MarshalJSON(); err != nil {
		o := fmt.Errorf("%#+v", e.Marshaler)
		return errors.Join(o, err).Error()
	} else {
		return string(buf)
	}
}

func newAPIErrorFromResponse[
	T any,
	E interface {
		json.Marshaler
		GetStatus() int
		GetDetail() string
	},
](
	method string,
	error E,
) (
	t *T,
	e *Error,
) {
	t = (*T)(nil)
	e = NewError(
		method,
		saclient.NewError(
			error.GetStatus(),
			error.GetDetail(),
			&apiErrorResponse{error},
		),
	)
	return
}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func ErrorFromDecodedResponse[T any](method string, yield func() (any, error)) (*T, error) {
	resp, err := yield()

	switch r := resp.(type) {
	case *T:
		return r, nil
	case *v1.Http400BadRequest:
		return newAPIErrorFromResponse[T](method, r)
	case *v1.Http401Unauthorized:
		return newAPIErrorFromResponse[T](method, r)
	case *v1.Http403Forbidden:
		return newAPIErrorFromResponse[T](method, r)
	case *v1.Http404NotFound:
		return newAPIErrorFromResponse[T](method, r)
	case *v1.Http409Conflict:
		return newAPIErrorFromResponse[T](method, r)
	case *v1.Http429TooManyRequests:
		return newAPIErrorFromResponse[T](method, r)
	case *v1.Http503ServiceUnavailable:
		return newAPIErrorFromResponse[T](method, r)
	default:
		if e, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](err); ok {
			return nil, NewError(method, httperror.NewError(method, e))
		} else if err != nil {
			return nil, NewError(method, err) // this is e.g. broken JSON
		} else {
			return nil, NewError(method, fmt.Errorf("unexpected type %#+v", resp))
		}
	}
}

func ErrorFromDecodedResponseE[T any](method string, yield func() (any, error)) (err error) {
	_, err = ErrorFromDecodedResponse[T](method, yield)
	return
}
