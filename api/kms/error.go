// Copyright 2025- The sacloud/kms-api-go authors
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

package kms

import (
	"errors"

	ogen "github.com/ogen-go/ogen/validate"
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
			return "kms: " + e.msg + ": " + e.err.Error()
		} else {
			return "kms: " + e.msg
		}
	} else {
		return "kms: " + e.err.Error()
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

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse[T any](method string, yield func() (*T, error)) (ret *T, err error) {
	ret, err = yield()

	if e, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](err); ok {
		err = NewError(method, httperror.NewError(method, e))
	} else if err != nil {
		err = NewError(method, err) // this is e.g. broken JSON
	}

	return
}
