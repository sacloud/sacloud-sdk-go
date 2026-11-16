// Copyright 2025- The sacloud/simplemq-api-go authors
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

package simplemq

//go:generate go run ./gen/list_of_errors.go

import (
	"fmt"

	"github.com/go-faster/jx"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

type Error struct {
	msg string
	err error
}

func (e *Error) Error() string {
	if e.msg != "" {
		if e.err != nil {
			return "simplemq: " + e.msg + ": " + e.err.Error()
		} else {
			return "simplemq: " + e.msg
		}
	} else {
		return "simplemq: " + e.err.Error()
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

func unknownQueueError(p any) error {
	if err, ok := asQueueError(p); !ok {
		return fmt.Errorf("unknown API response type %#+v", p)
	} else if err == nil {
		return fmt.Errorf("nil API error response %T", p)
	} else {
		var e jx.Encoder
		err.Encode(&e)
		return fmt.Errorf("unknown API error response: %s", e.String())
	}
}

func unknownMessageError(p any) error {
	if err, ok := asMessageError(p); !ok {
		return fmt.Errorf("unknown API response type %#+v", p)
	} else if err == nil {
		return fmt.Errorf("nil API error response %T", p)
	} else {
		var e jx.Encoder
		err.Encode(&e)
		return fmt.Errorf("unknown API error response: %s", e.String())
	}
}
