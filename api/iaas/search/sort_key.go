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

package search

import (
	"encoding/json"
	"fmt"
)

// SortOrder ソート順
type SortOrder int

const (
	// SortAsc 昇順(デフォルト)
	SortAsc SortOrder = iota
	// SortDesc 降順
	SortDesc
)

// SortKeys ソート順指定
type SortKeys []SortKey

// SortKey ソート順指定対象のフィールド名
type SortKey struct {
	Key   string
	Order SortOrder
}

// SortKeyAsc 昇順ソートキー
func SortKeyAsc(key string) SortKey {
	return SortKey{
		Key:   key,
		Order: SortAsc,
	}
}

// SortKeyDesc 降順ソートキー
func SortKeyDesc(key string) SortKey {
	return SortKey{
		Key:   key,
		Order: SortDesc,
	}
}

// MarshalJSON キーの文字列表現
func (k SortKey) MarshalJSON() ([]byte, error) {
	s := k.Key
	if k.Order == SortDesc {
		s = "-" + k.Key
	}
	return json.Marshal(s)
}

// UnmarshalJSON キーの文字列表現("Name"、降順は"-Name")からのデコード
//
// MarshalJSON が出力する文字列形式に加え、従来のオブジェクト形式({"Key":"Name","Order":1})も受け付ける。
func (k *SortKey) UnmarshalJSON(data []byte) error {
	// null の場合は値を変更しない(encoding/json の規約に従う)
	if string(data) == "null" {
		return nil
	}

	if len(data) > 0 && data[0] == '"' {
		// 文字列形式: 昇順は"Name"、降順は"-Name"
		var s string
		if err := json.Unmarshal(data, &s); err != nil {
			return err
		}
		if len(s) > 0 && s[0] == '-' {
			k.Key = s[1:]
			k.Order = SortDesc
		} else {
			k.Key = s
			k.Order = SortAsc
		}
		return nil
	}

	if len(data) > 0 && data[0] == '{' {
		// オブジェクト形式: {"Key":"Name","Order":1}
		type alias SortKey
		var a alias
		if err := json.Unmarshal(data, &a); err != nil {
			return err
		}
		*k = SortKey(a)
		return nil
	}

	return fmt.Errorf("search: cannot unmarshal %s into SortKey (must be a string like \"Name\" or \"-Name\", or an object like {\"Key\":\"Name\",\"Order\":1})", data)
}
