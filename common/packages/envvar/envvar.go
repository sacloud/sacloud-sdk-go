// Copyright 2022-2025 The sacloud/packages-go Authors
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

package envvar

import (
	"os"
	"strconv"
	"strings"
)

// StringFromEnv 環境変数から指定のキーの値を読み取って返す。ゼロ値だった場合はdefaultValueを返す
func StringFromEnv(key, defaultValue string) string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	return v
}

// StringFromEnvMulti 指定の環境変数を順番に読み取り、最初に見つかった非ゼロ値を返す。すべてゼロ値だった場合はdefaultValueを返す
func StringFromEnvMulti(keys []string, defaultValue string) string {
	for _, key := range keys {
		v := os.Getenv(key)
		if v != "" {
			return v
		}
	}
	return defaultValue
}

// StringSliceFromEnv 環境変数から指定のキーの値を読み取って、スライスに変換して返す。
//
// 読み取った環境変数の値がゼロ値の場合はdefaultValueを返す。
// 以外の場合はカンマで区切り、半角スペースをトリムした上でスライスに格納して返す。
func StringSliceFromEnv(key string, defaultValue []string) []string {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	values := strings.Split(v, ",")
	for i := range values {
		values[i] = strings.Trim(values[i], " ")
	}
	return values
}

// StringSliceFromEnvMulti 指定の環境変数を順番に読み取り、スライスに変換して返す。
//
// 読み取った環境変数の値がゼロ値の場合はdefaultValueを返す。
// 以外の場合はカンマで区切り、半角スペースをトリムした上でスライスに格納して返す。
func StringSliceFromEnvMulti(keys []string, defaultValue []string) []string {
	for _, key := range keys {
		v := os.Getenv(key)
		if v != "" {
			values := strings.Split(v, ",")
			for i := range values {
				values[i] = strings.Trim(values[i], " ")
			}
			return values
		}
	}
	return defaultValue
}

// IntFromEnv 環境変数から指定のキーの値を読み取って返す。ゼロ値だった場合はdefaultValueを返す
func IntFromEnv(key string, defaultValue int) int {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultValue
	}
	return int(i)
}

// IntFromEnvMulti 指定の環境変数を順番に読み取り、最初に見つかった非ゼロ値を返す。すべてゼロ値だった場合はdefaultValueを返す
func IntFromEnvMulti(keys []string, defaultValue int) int {
	for _, key := range keys {
		v := os.Getenv(key)
		if v != "" {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return defaultValue
			}
			return int(i)
		}
	}
	return defaultValue
}

// Int64FromEnv 環境変数から指定のキーの値を読み取って返す。ゼロ値だった場合はdefaultValueを返す
func Int64FromEnv(key string, defaultValue int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return defaultValue
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return defaultValue
	}
	return i
}

// Int64FromEnvMulti 指定の環境変数を順番に読み取り、最初に見つかった非ゼロ値を返す。すべてゼロ値だった場合はdefaultValueを返す
func Int64FromEnvMulti(keys []string, defaultValue int64) int64 {
	for _, key := range keys {
		v := os.Getenv(key)
		if v != "" {
			i, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return defaultValue
			}
			return i
		}
	}
	return defaultValue
}
