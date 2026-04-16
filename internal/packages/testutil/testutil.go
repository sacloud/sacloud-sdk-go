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

package testutil

import (
	"math/rand"
	"os"
	"time"
)

func init() {
	rand.Seed(time.Now().UTC().UnixNano())
}

const (
	// CharSetAlphaNum アフファベット(小文字)+数値
	CharSetAlphaNum = "abcdefghijklmnopqrstuvwxyz012346789"

	// CharSetAlpha アフファベット(小文字)
	CharSetAlpha = "abcdefghijklmnopqrstuvwxyz"

	// CharSetNumber 数値
	CharSetNumber = "012346789"
)

// Random ランダムな文字列を生成して返す
func Random(strlen int, charSet string) string {
	return RandomName("", strlen, charSet)
}

// RandomName ランダムな文字列を生成して返す
func RandomName(prefix string, strlen int, charSet string) string {
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = charSet[rand.Intn(len(charSet))] //nolint: gosec
	}
	return prefix + string(result)
}

// IsAccTest TESTACC環境変数が指定されているか
func IsAccTest() bool {
	return os.Getenv("TESTACC") != ""
}

// TestT テストのライフサイクルを管理するためのインターフェース.
//
// 通常は*testing.Tを実装として利用する
type TestT interface {
	Log(args ...interface{})
	Logf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	FailNow()
	Fatal(args ...interface{})
	Skip(args ...interface{})
	Skipf(format string, args ...interface{})
	Name() string
	Parallel()
}

// PreCheckEnvsFunc 指定の環境変数が指定されていなかった場合にテストをスキップするためのFuncを返す
func PreCheckEnvsFunc(envs ...string) func(TestT) {
	return func(t TestT) {
		for _, env := range envs {
			v := os.Getenv(env)
			if v == "" {
				t.Skipf("environment variable %q is not set. skip", env)
			}
		}
	}
}
