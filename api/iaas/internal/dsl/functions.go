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

package dsl

import (
	"fmt"
	"strings"

	"github.com/huandu/xstrings"
)

func uniqStrings(ss []string) []string {
	seen := make(map[string]struct{}, len(ss))
	i := 0
	for _, v := range ss {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		ss[i] = v
		i++
	}
	return ss[:i]
}

func wrapByDoubleQuote(targets ...string) []string {
	var ss []string
	for _, s := range targets {
		ss = append(ss, fmt.Sprintf(`"%s"`, s))
	}
	return ss
}

func toSnakeCaseName(name string) string {
	return strings.ReplaceAll(normalizeResourceName(xstrings.ToSnakeCase(name)), "-", "_")
}

func toLower(name string) string {
	return strings.ReplaceAll(normalizeResourceName(xstrings.ToSnakeCase(name)), "_", "")
}

var normalizationWords = map[string]string{
	"Ip":    "IP",
	"i_pv":  "ipv",
	"i_pv_": "ipv",
	"i-pv-": "ipv",
	"Cpu":   "CPU",
	"Ssd":   "SSD",
	"Hdd":   "HDD",
}

func normalizeResourceName(name string) string {
	n := name
	for k, v := range normalizationWords {
		if strings.Contains(name, k) {
			n = strings.ReplaceAll(name, k, v)
			break
		}
	}
	return n
}

func firstRuneToLower(name string) string {
	return xstrings.FirstRuneToLower(name)
}
