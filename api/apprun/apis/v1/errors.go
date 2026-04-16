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
	"encoding/json"
	"fmt"
	"strings"
)

var (
	_ error = (*ModelDefaultError)(nil)
	_ error = (*ModelCloudctrlError)(nil)
	_ error = (*ModelAppRunError)(nil)
)

func (e ModelDefaultError) Error() string {
	var in string
	if len(e.Detail.Errors) == 0 {
		in = "(empty)"
	} else {
		var errorStrings []string
		for _, errItem := range e.Detail.Errors {
			domain := ""
			reason := ""
			message := ""
			locationType := ""
			location := ""
			if errItem.Domain != nil {
				domain = *errItem.Domain
			}
			if errItem.Reason != nil {
				reason = *errItem.Reason
			}
			if errItem.Message != nil {
				message = *errItem.Message
			}
			if errItem.LocationType != nil {
				locationType = string(*errItem.LocationType)
			}
			if errItem.Location != nil {
				location = *errItem.Location
			}
			errorStrings = append(errorStrings, fmt.Sprintf("{domain: %s, reason: %s, message: %s, location_type: %s, location: %s}", domain, reason, message, locationType, location))
		}

		in = strings.Join(errorStrings, ", ")
	}

	return fmt.Sprintf("code: %d, message: %s, inner_error: %s", e.Detail.Code, e.Detail.Message, in)
}

func (e ModelCloudctrlError) Error() string {
	return fmt.Sprintf("CloudctrlError: %s (code: %s, fatal: %v, serial: %s, status: %s)",
		e.ErrorMsg, e.ErrorCode, e.IsFatal, e.Serial, e.Status)
}

func (e ModelAppRunError) Error() string {
	var defaultErr ModelDefaultError
	if err := json.Unmarshal(e.union, &defaultErr); err == nil && defaultErr.Detail.Message != "" {
		return defaultErr.Error()
	}

	var cloudctrlErr ModelCloudctrlError
	if err := json.Unmarshal(e.union, &cloudctrlErr); err == nil && cloudctrlErr.ErrorMsg != "" {
		return cloudctrlErr.Error()
	}

	return string(e.union)
}
