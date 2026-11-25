// Copyright 2025- The sacloud/addon-api-go authors
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

package addon

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	ogen "github.com/ogen-go/ogen/validate"
	v1 "github.com/sacloud/sacloud-sdk-go/api/addon/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/httperror"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

type Error struct {
	msg string
	err error
}

type apiErrorResponse struct{ json.Marshaler }

func (e *Error) Unwrap() error { return e.err }
func (e *Error) Error() string {
	var buf strings.Builder

	buf.WriteString("addon")

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

func (e *apiErrorResponse) Error() string {
	if buf, err := e.Marshaler.MarshalJSON(); err != nil {
		o := fmt.Errorf("%#+v", e.Marshaler)
		return errors.Join(o, err).Error()
	} else {
		return string(buf)
	}
}

func NewError(msg string, err error) *Error { return &Error{msg: msg, err: err} }

func NewAPIError(method string, j json.Marshaler) *Error {
	return NewError(method, &apiErrorResponse{Marshaler: j})
}

func NewNotFoundError(method string) *Error {
	return NewError(method, saclient.NewError(http.StatusNotFound, "", nil))
}

type void struct{}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func ErrorFromDecodedResponse[
	T, E any,
](
	method string,
	yield func() (any, error),
) (
	res *T,
	err error,
) {
	var ret any

	ret, err = yield()

	if e, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](err); ok {
		err = NewError(method, httperror.NewError(method, e))
	} else if err != nil {
		err = NewError(method, err) // this is e.g. broken JSON
	} else if ng, is400 := ret.(*v1.ErrorResponse); is400 {
		err = NewAPIError(method, ng)
	} else if _, isErr := ret.(*E); isErr {
		err = NewNotFoundError(method)
	} else if r, isOk := ret.(*T); !isOk {
		var zero T
		panic(fmt.Sprintf("%s: unexpected response type %T (%T expected): %#+v", method, ret, zero, ret)) // UNLIKELY
	} else {
		res = r
	}

	return
}

func ErrorFromDecodedResponseE[
	T, E any,
](
	method string,
	yield func() (any, error),
) (
	err error,
) {
	_, err = ErrorFromDecodedResponse[T, E](method, yield)
	return
}
