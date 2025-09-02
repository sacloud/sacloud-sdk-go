// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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
	"encoding/json"
	"reflect"

	"github.com/go-faster/errors"
)

func Unwrap[T json.Unmarshaler, U json.Marshaler](dst T, src U) (T, error) {
	// :FIXME: I know it's super inefficient, but works anyways...
	// Contributions welcomed.
	if bytes, err := src.MarshalJSON(); err != nil {
		return dst, errors.Wrapf(err, "failed to marshal source")
	} else if err := dst.UnmarshalJSON(bytes); err != nil {
		return dst, errors.Wrapf(err, "failed to unmarshal to destination: %s", reflect.TypeOf(dst).String())
	} else {
		return dst, nil
	}
}
