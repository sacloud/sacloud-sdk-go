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

func unknownError(p any) error {
	if schema, ok := asErrorSchema(p); !ok {
		return fmt.Errorf("unknown API response type %#+v", p)
	} else if schema == nil {
		return fmt.Errorf("nil API error response %T", p)
	} else if msg, ok := schema.GetMessage().Get(); !ok {
		return fmt.Errorf("API error response missing message: %#+v", p)
	} else if msg == "" {
		return fmt.Errorf("API error response has empty message: %#+v", p)
	} else {
		return errors.New(msg)
	}
}
