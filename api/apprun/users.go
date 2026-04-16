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

package apprun

import (
	"context"
	"net/http"
)

type UserAPI interface {
	// Read ログイン中のユーザー情報を取得
	Read(ctx context.Context) (*http.Response, error)
	// Create さくらのAppRunにサインアップ
	Create(ctx context.Context) (*http.Response, error)
}

var _ UserAPI = (*userOp)(nil)

type userOp struct {
	client *Client
}

// NewUserOp ユーザー操作関連API
func NewUserOp(client *Client) UserAPI {
	return &userOp{client: client}
}

func (op *userOp) Read(ctx context.Context) (*http.Response, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.GetUser(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (op *userOp) Create(ctx context.Context) (*http.Response, error) {
	apiClient, err := op.client.apiClient()
	if err != nil {
		return nil, err
	}
	resp, err := apiClient.PostUser(ctx)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
