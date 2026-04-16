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

package testutil

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/accessor"
	"github.com/sacloud/iaas-api-go/types"
)

// CRUDTestCase CRUD操作テストケース
type CRUDTestCase struct {
	// PreCheck テスト実行 or スキップを判定するためのFunc
	PreCheck func(TestT)

	// APICallerのセットアップ用Func、テストケースごとに1回呼ばれる
	SetupAPICallerFunc func() iaas.APICaller

	// Setup テスト前の準備(依存リソースの作成など)を行うためのFunc(省略可)
	Setup func(*CRUDTestContext, iaas.APICaller) error

	// Create Create操作のテスト用Func(省略可)
	Create *CRUDTestFunc

	// Read Read操作のテスト用Func(必須)
	Read *CRUDTestFunc

	// Updates Update操作のテスト用Func(省略可)
	Updates []*CRUDTestFunc

	// Shutdown Delete操作の前のシャットダウン(省略可)
	Shutdown func(*CRUDTestContext, iaas.APICaller) error

	// Delete Delete操作のテスト用Func(省略可)
	Delete *CRUDTestDeleteFunc

	// Cleanup APIで作成/変更したリソースなどのクリーンアップ用Func(省略化)
	Cleanup func(*CRUDTestContext, iaas.APICaller) error

	// Parallel t.Parallelを呼ぶかのフラグ
	Parallel bool

	// IgnoreStartupWait リソース作成後の起動待ちを行わない
	IgnoreStartupWait bool
}

// CRUDTestContext CRUD操作テストでのコンテキスト、一連のテスト中に共有される
type CRUDTestContext struct {
	// ID CRUDテスト対象リソースのID
	//
	// Create/Read/Updateの戻り値がidAccessorの場合に各操作の後で設定される
	ID types.ID

	// Values 一連のテスト中に共有したい値
	//
	// 依存リソースのIDの保持などで利用する
	Values map[string]interface{}

	// LastValue 最後の操作での戻り値
	LastValue interface{}

	ctx  context.Context
	once sync.Once
}

func (c *CRUDTestContext) initInnerContext() {
	c.once.Do(func() {
		c.ctx = context.TODO()
	})
}

// Deadline context.Context実装
func (c *CRUDTestContext) Deadline() (deadline time.Time, ok bool) {
	c.initInnerContext()
	return c.ctx.Deadline()
}

// Done context.Context実装
func (c *CRUDTestContext) Done() <-chan struct{} {
	c.initInnerContext()
	return c.ctx.Done()
}

// Err context.Context実装
func (c *CRUDTestContext) Err() error {
	c.initInnerContext()
	return c.ctx.Err()
}

// Value context.Context実装
func (c *CRUDTestContext) Value(key interface{}) interface{} {
	c.initInnerContext()
	return c.ctx.Value(key)
}

// CRUDTestFunc CRUD操作(DELETE以外)テストでのテスト用Func
type CRUDTestFunc struct {
	// Func API操作を行うFunc
	Func func(*CRUDTestContext, iaas.APICaller) (interface{}, error)

	// CheckFunc 任意のチェックを行うためのFunc、省略可能。
	CheckFunc func(TestT, *CRUDTestContext, interface{}) error

	// SkipExtractID Trueの場合Funcの戻り値からのID抽出を行わない
	SkipExtractID bool
}

// CRUDTestDeleteFunc CRUD操作テストのDeleteテスト用Func
type CRUDTestDeleteFunc struct {
	// Func API操作を行うFunc
	Func func(*CRUDTestContext, iaas.APICaller) error
}

// CRUDTestExpect CRUD操作(DELETE以外)テストでの期待値
type CRUDTestExpect struct {
	// ExpectValue CRUD操作実行後の期待値
	ExpectValue interface{}

	// IgnoreFields比較時に無視する項目
	IgnoreFields []string
}

// Prepare テスト対象値を受け取り、比較可能な状態に加工した対象値と期待値を返す
func (c *CRUDTestExpect) Prepare(actual interface{}) (interface{}, interface{}) {
	toMap := func(v interface{}) map[string]interface{} {
		data, err := json.Marshal(v)
		if err != nil {
			log.Fatalf("prepare is failed: json.Marshal returned error: %s", err)
		}
		var m map[string]interface{}
		if err := json.Unmarshal(data, &m); err != nil {
			log.Fatalf("prepare is failed: json.Unmarshal returned error: %s", err)
		}
		for _, key := range c.IgnoreFields {
			c.deleteByPath(m, key)
		}
		return m
	}

	return toMap(c.ExpectValue), toMap(actual)
}

// deleteByPath removes a nested field specified by dot-separated path (e.g. "aaa.bbb")
func (c *CRUDTestExpect) deleteByPath(m map[string]interface{}, path string) {
	keys := strings.Split(path, ".")
	last := len(keys) - 1

	current := m
	for i, k := range keys {
		v, ok := current[k]
		if !ok {
			// key not found → nothing to do
			return
		}

		if i == last {
			// 最終キーなら削除
			delete(current, k)
			return
		}

		// 次に進むための map かどうか確認
		next, ok := v.(map[string]interface{})
		if !ok {
			// 途中で map でなくなったら削除不能
			return
		}
		current = next
	}
}

// RunCRUD 任意の条件でCRUD操作をテストする
func RunCRUD(t TestT, testCase *CRUDTestCase) {
	if testCase.SetupAPICallerFunc == nil {
		t.Fatal("CRUDTestCase.SetupAPICallerFunc is required")
	}

	if testCase.Parallel {
		t.Parallel()
	}

	if testCase.PreCheck != nil {
		testCase.PreCheck(t)
	}

	testContext := &CRUDTestContext{
		Values: make(map[string]interface{}),
	}
	defer func() {
		// Cleanup
		if testCase.Cleanup != nil {
			if err := testCase.Cleanup(testContext, testCase.SetupAPICallerFunc()); err != nil {
				t.Logf("Cleanup is failed: ", err)
			}
		}
		if err := recover(); err != nil {
			t.Logf("Unexpected error is occurred: %v, Trace: %s", err, string(debug.Stack()))
		}
	}()

	if testCase.Setup != nil {
		if err := testCase.Setup(testContext, testCase.SetupAPICallerFunc()); err != nil {
			t.Error("Setup is failed: ", err)
			return
		}
	}

	testFunc := func(f *CRUDTestFunc) error {
		actual, err := f.Func(testContext, testCase.SetupAPICallerFunc())
		if err != nil {
			return err
		}
		testContext.LastValue = actual

		if actual != nil && f.CheckFunc != nil {
			if err := f.CheckFunc(t, testContext, actual); err != nil {
				return err
			}
		}

		// extract ID from result of f.Func()
		if actual != nil && !f.SkipExtractID {
			if idHolder, ok := actual.(accessor.ID); ok {
				testContext.ID = idHolder.GetID()
			}
		}

		return nil
	}

	// Create
	if testCase.Create != nil {
		if err := testFunc(testCase.Create); err != nil {
			t.Error("Create is failed: ", err)
			return
		}

		if !testCase.IgnoreStartupWait && testCase.Read != nil && testContext.LastValue != nil {
			waiter := iaas.WaiterForApplianceUp(func() (interface{}, error) {
				return testCase.Read.Func(testContext, testCase.SetupAPICallerFunc())
			}, 100)
			if _, err := waiter.WaitForState(context.TODO()); err != nil {
				t.Error("WaitForUp is failed: ", err)
				return
			}
		}
	}

	// Read
	if testCase.Read != nil {
		if err := testFunc(testCase.Read); err != nil {
			t.Fatal("Read is failed: ", err)
		}
	}

	// Updates
	for _, updFunc := range testCase.Updates {
		if err := testFunc(updFunc); err != nil {
			t.Error("Update is failed: ", err)
			return
		}
	}

	// Shutdown
	if testCase.Shutdown != nil {
		if testCase.Read == nil {
			t.Log("CRUDTestCase.Shutdown is set, but CRUDTestCase.Read is nil. Shutdown is skipped")
		} else {
			v, err := testCase.Read.Func(testContext, testCase.SetupAPICallerFunc())
			if err != nil {
				t.Error("Shutdown is failed: ", err)
				return
			}
			if v, ok := v.(accessor.InstanceStatus); ok && v.GetInstanceStatus().IsUp() {
				if err := testCase.Shutdown(testContext, testCase.SetupAPICallerFunc()); err != nil {
					t.Error("Shutdown is failed: ", err)
					return
				}

				waiter := iaas.WaiterForDown(func() (interface{}, error) {
					return testCase.Read.Func(testContext, testCase.SetupAPICallerFunc())
				})
				if _, err := waiter.WaitForState(context.TODO()); err != nil {
					t.Error("WaitForDown is failed: ", err)
					return
				}
			}
		}
	}

	// Delete
	if testCase.Delete != nil {
		if err := testCase.Delete.Func(testContext, testCase.SetupAPICallerFunc()); err != nil {
			t.Error("Delete is failed: ", err)
			return
		}
		if testCase.Read != nil {
			// check not exists
			_, err := testCase.Read.Func(testContext, testCase.SetupAPICallerFunc())
			if err == nil {
				t.Error("Resource still exists: ", testContext.ID)
				return
			}
			if e, ok := err.(iaas.APIError); ok {
				if e.ResponseCode() != http.StatusNotFound {
					t.Error("Reading after delete is failed: ", e)
					return
				}
			}
		}
	}
}
