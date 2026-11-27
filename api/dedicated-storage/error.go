// Copyright 2022-2025 The sacloud/dedicated-storage-api-go Authors
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

package dedicatedstorage

import (
	"errors"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/dedicated-storage/apis/v1"
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

	buf.WriteString("dedicated-storage")

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

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse[T any](method string, yield func() (*T, error)) (ret *T, err error) {
	ret, err = yield()

	if e2, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](err); ok {
		err = NewError(method, httperror.NewError(method, e2))
	} else if e3, ok := errors.AsType[*v1.ErrorStatusCode](err); ok {
		err = NewAPIError(method, e3.GetStatusCode(), e3)
	} else if err != nil {
		err = NewError(method, err) // this is e.g. broken JSON
	}
	return
}

func errorFromDecodedResponseE(method string, yield func() error) (err error) {
	type never struct{}
	var zero never
	_, err = errorFromDecodedResponse(method, func() (*never, error) { return &zero, yield() })
	return
}
