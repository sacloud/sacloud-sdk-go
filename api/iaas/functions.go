// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package iaas

import (
	"fmt"
	"strings"
)

func buildURL(pathFormat string, param map[string]interface{}) (string, error) {
	var replPairs []string
	for k, v := range param {
		replPairs = append(replPairs, fmt.Sprintf("{{.%s}}", k), fmt.Sprint(v))
	}
	replacer := strings.NewReplacer(replPairs...)
	result := replacer.Replace(pathFormat)

	if strings.Contains(result, "{{") {
		return "", fmt.Errorf("undefined variable in URL template: %s", result)
	}
	return result, nil
}
