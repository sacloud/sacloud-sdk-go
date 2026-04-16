// Copyright 2025- The sacloud/secretmanager-api-go authors
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

package secretmanager

import (
	"errors"

	ogen "github.com/ogen-go/ogen/validate"
	"github.com/sacloud/saclient-go"
)

type Error struct {
	msg string
	err error
}

func (e *Error) Error() string {
	if e.msg != "" {
		if e.err != nil {
			return "secretmanager: " + e.msg + ": " + e.err.Error()
		} else {
			return "secretmanager: " + e.msg
		}
	} else {
		return "secretmanager: " + e.err.Error()
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

// secretmanagerのOpenAPI定義でエラーケースが定義されていないので、現状はogenのエラーから状態を取り出す
// NewAPIErrorは他のクライアントとインターフェイスを揃えるために維持し、別で生成関数を用意
func createAPIError(method string, err error) error {
	var unexpected *ogen.UnexpectedStatusCodeError
	if errors.As(err, &unexpected) {
		return NewAPIError(method, unexpected.StatusCode, err)
	}

	return NewAPIError(method, 0, err)
}
