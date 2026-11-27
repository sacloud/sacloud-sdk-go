// Copyright 2025- The sacloud/nosql-api-go authors
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

package apigw

//go:generate go run ./gen/list_of_errors.go

import (
	"errors"
	"fmt"

	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/apigw/apis/v1"
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
			return "apigw: " + e.msg + ": " + e.err.Error()
		} else {
			return "apigw: " + e.msg
		}
	} else {
		return "apigw: " + e.err.Error()
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

type apiErrorResponse struct{ *v1.ErrorSchema }

func (e *apiErrorResponse) Error() string {
	if e == nil {
		return fmt.Sprintf("nil API error response %T", e)
	} else if e.ErrorSchema == nil {
		return "nil ErrorSchema"
	} else if msg, ok := e.GetMessage().Get(); !ok {
		return fmt.Sprintf("API error response missing message: %#+v", e)
	} else if msg == "" {
		return fmt.Sprintf("API error response has empty message: %#+v", e)
	} else {
		return msg
	}
}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse[T any](method string, yield func() (any, error)) (ret *T, err error) {
	res, e1 := yield()

	if e2, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](e1); ok {
		err = NewError(method, httperror.NewError(method, e2))
	} else if e1 != nil {
		err = NewError(method, e1) // this is e.g. broken JSON
	} else if schema, ok := asErrorSchema(res); ok {
		err = NewError(method, &apiErrorResponse{ErrorSchema: schema})
	} else if ret, ok = res.(*T); !ok {
		err = NewError(method, fmt.Errorf("unknown API response type %#+v", res))
	}

	return
}

func errorFromDecodedResponseE[T any](method string, yield func() (any, error)) (err error) {
	_, err = errorFromDecodedResponse[T](method, yield)
	return
}
