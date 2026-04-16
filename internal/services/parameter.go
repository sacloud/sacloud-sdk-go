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

// ParameterInitializer パラメータの初期化時に特別な処理を行いたい場合に実装すべきインターフェース
type ParameterInitializer interface {
	// Initialize パラメータの初期カスタマイズ
	Initialize()
}

// ParameterValidator パラメータが独自のバリデーションを実装する場合に実装すべきインターフェース
type ParameterValidator interface {
	Validate() error
}

// ParameterExampleValuer パラメータが値の例示をサポートする場合に実装するインターフェース
//
// ここで返される値はJSONまたはYAMLで出力されることがあるため適切なタグを付与しておくこと
type ParameterExampleValuer interface {
	Examples() interface{}
}
