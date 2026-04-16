// Copyright 2021-2024 The sacloud/apprun-api-go authors
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

package v1

import (
	"fmt"
	"net/http"
	"reflect"
)

func eCoalesce(errs ...interface{}) error {
	for _, e := range errs {
		if !(e == nil || reflect.ValueOf(e).IsNil()) {
			return toError(e)
		}
	}
	return nil
}

// toError errorまたはerrorResponserの実装からerrorを返す
// 上記以外またはnilを渡した場合はpanicするため呼び出し側でチェックする必要がある
func toError(v interface{}) error {
	switch err := v.(type) {
	case error:
		return err
	default:
		msg := fmt.Sprintf("invalid arg: %#+v", v)
		panic(msg)
	}
}

var osStatusCodes = map[int]bool{
	http.StatusOK:        true,
	http.StatusCreated:   true,
	http.StatusAccepted:  true,
	http.StatusNoContent: true,
}

func isOKStatus(httpStatusCode int) bool {
	_, ok := osStatusCodes[httpStatusCode]
	return ok
}
