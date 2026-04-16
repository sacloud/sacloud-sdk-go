// Copyright 2022-2025 The sacloud/api-client-go Authors
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

package client

import (
	"errors"
	"fmt"
	"net/http"
)

// APIError APIコール時のエラーを表す型。各APIクライアントでこれをラップして利用する
type APIError struct {
	Code    int
	Message string
	Err     error
}

func (e *APIError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("API Error %d - %s", e.Code, e.Message)
	} else {
		return fmt.Sprintf("API Error %d - %s: %s", e.Code, e.Message, e.Err.Error())
	}
}

func (e *APIError) Unwrap() error {
	return e.Err
}

func NewAPIError(code int, msg string, err error) *APIError {
	if len(msg) == 0 {
		msg = http.StatusText(code)
		if msg == "" { // client uses 0 for unknown error
			msg = "unknown error"
		}
	}
	return &APIError{
		Code:    code,
		Message: msg,
		Err:     err,
	}
}

func IsNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	var apiError *APIError
	if errors.As(err, &apiError) {
		return apiError.Code == http.StatusNotFound
	}

	return false
}
