// Copyright 2021-2026 The sacloud/apprun-api-go authors
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

package apprun

import (
	stderrs "errors"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
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

type modelErrorGetter interface {
	GetModelDefaultError() (v1.ModelDefaultError, bool)
	GetModelCloudctrlError() (v1.ModelCloudctrlError, bool)
}

func apiErrorFromModel(method string, code int, response modelErrorGetter) *Error {
	if v, ok := response.GetModelDefaultError(); ok {
		if v.Error.Message != "" {
			return NewAPIError(method, code, stderrs.New(v.Error.Message))
		}
	}
	if v, ok := response.GetModelCloudctrlError(); ok {
		if v.ErrorMsg != "" {
			return NewAPIError(method, code, stderrs.New(v.ErrorMsg))
		}
	}
	return NewAPIError(method, code, stderrs.New("unknown error"))
}
