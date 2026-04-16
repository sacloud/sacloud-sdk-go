// Copyright 2022-2025 The sacloud/services Authors
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

package services

import (
	"strings"

	"github.com/sacloud/services/naming"
)

// Info サービスについての情報
type Info struct {
	// サービスの名前(ケバブケース)
	Name string

	// サービスの説明
	Description string

	// 親サービスの名称リスト
	//
	// 例えばservice1サービスがservice1/service2/service3という階層構造の場合、
	// ["service1", "service2"]となる。
	// この場合、service3サービスの各操作ではService1IdとService2Idというフィールドが存在し、必須パラメータであることが期待される。
	// (ここで指定した名前をアッパーキャメルケースにしたものがフィールド名になる)
	ParentServices []string

	// カテゴリー名(ケバブケース)
	CategoryName string
}

// FullName 親サービス名まで含めたサービス名称を返す
func (info *Info) FullName() string {
	elems := info.ParentServices
	return strings.Join(append(elems, info.Name), "/")
}

func (info *Info) Category() string {
	return naming.ToUpperCamelCase(info.CategoryName)
}
