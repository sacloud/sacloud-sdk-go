// Copyright 2025- The sacloud/monitoring-suite-api-go authors
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

package monitoringsuite

import (
	"fmt"
	"io"

	"github.com/go-faster/errors"
	ogen "github.com/ogen-go/ogen/validate"
	"github.com/sacloud/saclient-go"
)

type Error struct {
	msg string
	err error
}

func (e *Error) Error() string {
	if e.msg == "" {
		return "monitoringsuite: " + e.err.Error()
	}
	if e.err != nil {
		return "monitoringsuite: " + e.msg + ": " + e.err.Error()
	}
	return "monitoringsuite: " + e.msg
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

// convert UnexpectedStatusCodeError -> IsNotFoundError
func errorFromDecodedResponse[T any](method string, yield func() (T, error)) (res T, err error) {
	res, err = yield()
	if err == nil {
		return // no error
	}
	e, ok := errors.Into[*ogen.UnexpectedStatusCodeError](err)
	if !ok {
		// expected error
		err = NewAPIError(method, 0, err)
		return
	}
	msg, ee := io.ReadAll(e.Payload.Body)
	if ee != nil {
		// payload read error, maybe empty
		err = errors.Join(e, ee)
		return
	}
	// expect error payload to be at least human readable
	err = NewAPIError(method, e.StatusCode, fmt.Errorf("%s", string(msg)))

	return
}

func errorFromDecodedResponse1(method string, yield func() error) (err error) {
	_, err = errorFromDecodedResponse(method, func() ([]any, error) {
		return nil, yield()
	})
	return
}
