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

package dispatcher_test

import (
	"fmt"

	"github.com/sacloud/services"
	"github.com/sacloud/services/dispatcher"
	"github.com/sacloud/services/dummy"
)

func Example() {
	// サービスの登録
	dispatcher.Register("my-platform", services.Services{dummy.New()})

	// サービスの呼び出し(プラットフォーム名 + 対象リソース名 + 操作, パラメータを渡す)
	arguments := []string{"my-platform", "dummy", "read"}
	parameter := map[string]interface{}{"Param1": "example"}

	result, err := dispatcher.Dispatch(arguments, parameter)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.(*dummy.ReadResult).Dummy)
	// output: result
}
