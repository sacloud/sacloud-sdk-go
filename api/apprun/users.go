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
	"errors"
	"net/http"

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
	const methodName = "Users.Read"
	res, err := op.client.ReadUser(ctx)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerReadUser:
		return result, nil
	case *v1.ReadUserUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.ReadUserForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.ReadUserNotFound:
		return nil, apiErrorFromModel(methodName, http.StatusNotFound, result)
	case *v1.ReadUserInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}

func (op *userOp) Create(ctx context.Context) (*v1.HandlerCreateUser, error) {
	const methodName = "Users.Create"
	res, err := op.client.CreateUser(ctx)
	if err != nil {
		return nil, NewAPIError(methodName, 0, err)
	}
	switch result := res.(type) {
	case *v1.HandlerCreateUser:
		return result, nil
	case *v1.CreateUserUnauthorized:
		return nil, apiErrorFromModel(methodName, http.StatusUnauthorized, result)
	case *v1.CreateUserForbidden:
		return nil, apiErrorFromModel(methodName, http.StatusForbidden, result)
	case *v1.CreateUserConflict:
		return nil, apiErrorFromModel(methodName, http.StatusConflict, result)
	case *v1.CreateUserInternalServerError:
		return nil, apiErrorFromModel(methodName, http.StatusInternalServerError, result)
	default:
		return nil, NewAPIError(methodName, 0, errors.New("unknown error"))
	}
}
