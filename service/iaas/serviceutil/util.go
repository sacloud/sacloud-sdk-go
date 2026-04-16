// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package serviceutil

import (
	"fmt"
	"time"

	iaas "github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/mapconv"
	"github.com/sacloud/packages-go/size"
)

func HandleNotFoundError(err error, ignoreNotFoundError bool) error {
	if ignoreNotFoundError && iaas.IsNotFoundError(err) {
		return nil // ignore: 404 not found
	}
	return err
}

func MonitorCondition(start, end time.Time) (*iaas.MonitorCondition, error) {
	e := end
	if e.IsZero() {
		e = time.Now()
	}

	s := start
	if s.IsZero() {
		s = e.Add(-1 * time.Hour)
	}
	if !(s.Unix() <= e.Unix()) {
		return nil, fmt.Errorf("start(%s) or end(%s) is invalid", start.String(), end.String())
	}
	return &iaas.MonitorCondition{Start: s, End: e}, nil
}

func RequestConvertTo(source interface{}, dest interface{}) error {
	decoder := &mapconv.Decoder{
		Config: &mapconv.DecoderConfig{
			TagName: "service",
			FilterFuncs: map[string]mapconv.FilterFunc{
				"gb_to_mb": gbToMb,
			},
		},
	}
	return decoder.ConvertTo(source, dest)
}

func gbToMb(v interface{}) (interface{}, error) {
	s, ok := v.(int)
	if !ok {
		return nil, fmt.Errorf("invalid size value: %v", v)
	}
	return size.GiBToMiB(s), nil
}
