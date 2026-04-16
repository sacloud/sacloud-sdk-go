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
	"fmt"

	"github.com/sacloud/services/naming"
)

// SupportedOperation サービスが提供する操作のメタデータ
type SupportedOperation struct {
	// 操作(メソッド)の名前(ケバブケース)
	Name string

	// 操作についての説明
	Description string

	// 操作種別、種別によりIdが必要/不要が決定される
	OperationType OperationType

	// この操作が属するカテゴリーの名称(省略可、指定する場合はケバブケース)
	CategoryName string
}

func (op *SupportedOperation) EqualsByName(name string) bool {
	return naming.ToKebabCase(name) == naming.ToKebabCase(op.Name)
}

func (op *SupportedOperation) FuncName() string {
	return naming.ToUpperCamelCase(op.Name)
}

// WithContextFuncName Nameを持つFuncに対応する、context.Contextを受け取るFuncの名前を返す
func (op *SupportedOperation) WithContextFuncName() string {
	return fmt.Sprintf("%sWithContext", naming.ToUpperCamelCase(op.Name))
}

func (op *SupportedOperation) Category() string {
	return naming.ToUpperCamelCase(op.CategoryName)
}
