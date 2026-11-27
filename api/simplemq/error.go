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

	"errors"

	"github.com/go-faster/jx"

	ogen "github.com/ogen-go/ogen/validate"
	message "github.com/sacloud/sacloud-sdk-go/api/simplemq/apis/v1/message"
	queue "github.com/sacloud/sacloud-sdk-go/api/simplemq/apis/v1/queue"
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
	if p == nil || httperror.IsTypedNil(p) {
		return fmt.Errorf("nil API error response %T", p)
	}
	if err, _, ok := asQueueError(p); !ok {
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
	if p == nil || httperror.IsTypedNil(p) {
		return fmt.Errorf("nil API error response %T", p)
	}
	if err, _, ok := asMessageError(p); !ok {
		return fmt.Errorf("unknown API response type %#+v", p)
	} else if err == nil {
		return fmt.Errorf("nil API error response %T", p)
	} else {
		var e jx.Encoder
		err.Encode(&e)
		return fmt.Errorf("unknown API error response: %s", e.String())
	}
}

type apiError[E interface{ Encode(*jx.Encoder) }] struct{ err E }

func (e *apiError[E]) Error() string {
	if e == nil || httperror.IsTypedNil(e.err) {
		return fmt.Sprintf("nil API error response %T", e.err)
	}
	var enc jx.Encoder
	e.err.Encode(&enc)
	return enc.String()
}

// UnexpectedStatusCodeError -> IsNotFoundError conversion
func errorFromDecodedResponse[T any](method string, yield func() (any, error)) (ret *T, err error) {
	res, e := yield()

	if e1, ok := errors.AsType[*ogen.UnexpectedStatusCodeError](e); ok {
		err = NewError(method, httperror.NewError(method, e1))
	} else if e != nil {
		err = NewError(method, e) // this is e.g. broken JSON
	} else if e2, code, ok := asQueueError(res); ok {
		err = NewAPIError(method, code, &apiError[*queue.Error]{err: e2})
	} else if e3, code, ok := asMessageError(res); ok {
		err = NewAPIError(method, code, &apiError[*message.Error]{err: e3})
	} else if ret, ok = res.(*T); !ok {
		err = fmt.Errorf("unexpected response type %T", res)
	}
	return
}

func errorFromDecodedResponseE[T any](method string, yield func() (any, error)) (err error) {
	_, err = errorFromDecodedResponse[T](method, yield)
	return
}
