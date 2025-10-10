// Copyright 2025- The sacloud/http-client-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package client

import (
	"errors"
	"fmt"
	"strings"
)

// Error represents an error within the HTTP client.
type Error struct {
	code int
	msg  string
	err  error
}

// :TODO: shim for golang 1.26+
func asType[T error](e error) (T, bool) {
	var t T
	ok := errors.As(e, &t)
	return t, ok
}

func compose(code int, msg string, err error) Error { return Error{code: code, msg: msg, err: err} }
func (e *Error) decompose() (int, string, error)    { return e.code, e.msg, e.err }
func (e *Error) codeIs(code int) bool               { return e.code == code }

func NewError(code int, msg string, err error) error {
	ret := compose(code, msg, err)

	return &ret
}

func NewErrorf(format string, args ...any) error {
	ret := compose(0, fmt.Sprintf(format, args...), nil)

	return &ret
}

func Wrapf(err error, format string, args ...any) error {
	if err == nil {
		return nil

	} else {
		ret := compose(0, fmt.Sprintf(format, args...), err)

		return &ret
	}
}

// Implements the error interface.
// Returns a stringized representation of the error.
func (e *Error) Error() string {
	// This error message format mimics the one used in sacloud/api-client-go.

	if e == nil {
		return ""
	}

	var buf strings.Builder
	code, msg, err := e.decompose()

	buf.WriteString("API Error")

	if code != 0 {
		buf.WriteString(fmt.Sprintf(" %d", code))
	}

	if msg != "" {
		buf.WriteString(" - ")
		buf.WriteString(msg)
	}

	if err != nil {
		buf.WriteString(": ")
		buf.WriteString(err.Error())
	}

	return buf.String()
}

// Unwrap returns the underlying error.
func (e *Error) Unwrap() error {
	if e == nil {
		return nil

	} else {
		_, _, err := e.decompose()

		return err
	}
}

// Returns whether the given error is an Error with a 404 status code.
// Provided here for compatibility with sacloud/api-client-go.
func IsNotFoundError(err error) bool {
	if e, ok := asType[*Error](err); ok {
		return e.codeIs(404)

	} else {
		return false
	}
}
