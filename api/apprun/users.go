// Copyright 2021-2026 The sacloud/apprun-api-go authors
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

package apprun

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun/apis/v1"
)

type UserAPI interface {
	// Read ログイン中のユーザー情報を取得
	Read(ctx context.Context) (*v1.HandlerReadUser, error)
	// Create さくらのAppRunにサインアップ
	Create(ctx context.Context) (*v1.HandlerCreateUser, error)
}

var _ UserAPI = (*userOp)(nil)

type userOp struct {
	client *v1.Client
}

// NewUserOp ユーザー操作関連API
func NewUserOp(client *v1.Client) UserAPI {
	return &userOp{client: client}
}

func (op *userOp) Read(ctx context.Context) (*v1.HandlerReadUser, error) {
	return errorFromDecodedResponse[v1.HandlerReadUser]("Users.Read", func() (any, error) {
		return op.client.ReadUser(ctx)
	})
}

func (op *userOp) Create(ctx context.Context) (*v1.HandlerCreateUser, error) {
	return errorFromDecodedResponse[v1.HandlerCreateUser]("Users.Create", func() (any, error) {
		return op.client.CreateUser(ctx)
	})
}
