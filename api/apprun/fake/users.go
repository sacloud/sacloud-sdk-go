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

var (
	UserId   = 111
	UserName = "user_name"
)

// API上でUserの定義は存在しないが、ユーザーのサインアップを管理するために定義する
type User struct {
	Id   int
	Name string
}

func (engine *Engine) GetUser() error {
	defer engine.rLock()()

	if u := engine.User; u == nil {
		return newError(
			ErrorTypeNotFound, "user", nil,
			"さくらのAppRunにユーザーが存在しません。")
	}

	return nil
}

func (engine *Engine) CreateUser() error {
	defer engine.lock()()

	if u := engine.User; u != nil {
		return newError(
			ErrorTypeConflict, "user", u.Id,
			"さくらのAppRunにユーザーが既に存在します。")
	}

	engine.User = &User{
		Id:   UserId,
		Name: UserName,
	}
	return nil
}
