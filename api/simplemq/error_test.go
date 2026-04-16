// Copyright 2022-2025 The sacloud/simplemq-api-go Authors
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

import (
	"errors"
	"testing"

	"github.com/sacloud/saclient-go"
)

func TestError_Error(t *testing.T) {
	baseErr := errors.New("base error")

	tests := []struct {
		name string
		err  *Error
		want string
	}{
		{
			name: "with msg and err",
			err:  &Error{msg: "something failed", err: baseErr},
			want: "simplemq: something failed: base error",
		},
		{
			name: "with msg only",
			err:  &Error{msg: "only msg"},
			want: "simplemq: only msg",
		},
		{
			name: "with err only",
			err:  &Error{err: baseErr},
			want: "simplemq: base error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestNewError(t *testing.T) {
	baseErr := errors.New("base error")

	err := NewError("msg", baseErr)
	if err.msg != "msg" || err.err != baseErr {
		t.Errorf("NewError() did not set fields correctly")
	}

	err2 := NewError("msg only", nil)
	if err2.msg != "msg only" || err2.err != nil {
		t.Errorf("NewError() with nil err did not set fields correctly")
	}
}

func TestNewAPIError(t *testing.T) {
	baseErr := errors.New("base error")

	err := NewAPIError("msg", 404, baseErr)
	if !saclient.IsNotFoundError(err) {
		t.Errorf("IsNotFoundError is true for NewAPIError with 404")
	}

	err2 := NewAPIError("msg", 503, nil)
	if saclient.IsNotFoundError(err2) {
		t.Errorf("IsNotFoundError is false for NewAPIError with 503")
	}
}
