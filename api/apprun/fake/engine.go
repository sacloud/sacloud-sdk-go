// Copyright 2021-2024 The sacloud/apprun-api-go authors
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

package fake

import (
	"sync"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

// Engine Fakeサーバであつかうダミーデータを表す
//
// Serverに渡した後はDataStore内のデータを外部から操作しないこと
type Engine struct {
	User         *User
	Applications []*v1.Application
	Versions     []*v1.Version

	// MapのkeyにApplicationのIDを利用する
	Traffics map[string][]*v1.Traffic

	// MapのkeyにApplicationのIDを利用する
	appVersionRelations map[string][]*appVersionRelation

	// MapのkeyにApplicationのIDを利用する
	appPacketFilterRelations map[string]*v1.HandlerGetPacketFilter

	// GeneratedID 採番済みの最終ID
	//
	// DataStoreの各フィールドの値との整合性は確認されないため利用者側が管理する必要がある
	GeneratedVersionID int

	mu sync.RWMutex
}

type appVersionRelation struct {
	application *v1.Application
	version     *v1.Version
}

func NewEngine() *Engine {
	return &Engine{
		appVersionRelations:      make(map[string][]*appVersionRelation),
		appPacketFilterRelations: make(map[string]*v1.HandlerGetPacketFilter),
		Traffics:                 make(map[string][]*v1.Traffic),
	}
}

func (engine *Engine) lock() func() {
	engine.mu.Lock()
	return engine.mu.Unlock
}

func (engine *Engine) rLock() func() {
	engine.mu.RLock()
	return engine.mu.RUnlock
}

// nextId GeneratedVersionID+1したものを返す
//
// ロックは行わないため呼び出し側で適切に制御すること
func (engine *Engine) nextVersionId() int {
	engine.GeneratedVersionID++
	id := engine.GeneratedVersionID
	return id
}
