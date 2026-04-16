// Copyright 2025- The sacloud/iam-api-go authors
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

// Package user provides the UserAPI that wraps the generated v1 client.
package user

import (
	"context"

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

// UserAPI is the interface for user operations.
type UserAPI interface {
	List(ctx context.Context, params ListParams) (*v1.CompatUsersGetOK, error)
	Create(ctx context.Context, params CreateParams) (*v1.User, error)
	Read(ctx context.Context, id int) (*v1.User, error)
	Update(ctx context.Context, id int, params UpdateParams) (*v1.User, error)
	Delete(ctx context.Context, id int) error

	RegisterEmail(ctx context.Context, userID int, email string) error
	UnregisterEmail(ctx context.Context, userID int) error
}

type userOp struct {
	client *v1.Client
}

func NewUserOp(client *v1.Client) UserAPI {
	return &userOp{client: client}
}

type ListParams struct {
	Page     *int
	PerPage  *int
	Ordering *v1.CompatUsersGetOrdering
}

func (u *userOp) List(ctx context.Context, params ListParams) (*v1.CompatUsersGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.CompatUsersGetOK]("User.List", func() (any, error) {
		return u.client.CompatUsersGet(ctx, v1.CompatUsersGetParams{
			Page:     common.IntoOpt[v1.OptInt](params.Page),
			PerPage:  common.IntoOpt[v1.OptInt](params.PerPage),
			Ordering: common.IntoOpt[v1.OptCompatUsersGetOrdering](params.Ordering),
		})
	})
}

type CreateParams struct {
	Name        string
	Password    string
	Code        string
	Description string
	Email       *string
}

func (u *userOp) Create(ctx context.Context, params CreateParams) (*v1.User, error) {
	return common.ErrorFromDecodedResponse[v1.User]("User.Create", func() (any, error) {
		return u.client.CompatUsersPost(ctx, &v1.CompatUsersPostReq{
			Name:        params.Name,
			Password:    params.Password,
			Code:        params.Code,
			Description: params.Description,
			Email:       common.IntoOpt[v1.OptString](params.Email),
		})
	})
}

func (u *userOp) Read(ctx context.Context, id int) (*v1.User, error) {
	return common.ErrorFromDecodedResponse[v1.User]("User.Read", func() (any, error) {
		return u.client.CompatUsersUserIDGet(ctx, v1.CompatUsersUserIDGetParams{UserID: id})
	})
}

type UpdateParams struct {
	Name        string
	Password    *string
	Description string
}

func (u *userOp) Update(ctx context.Context, id int, params UpdateParams) (*v1.User, error) {
	return common.ErrorFromDecodedResponse[v1.User]("User.Update", func() (any, error) {
		req := v1.CompatUsersUserIDPutReq{
			Name:        params.Name,
			Password:    common.IntoOpt[v1.OptString](params.Password),
			Description: params.Description,
		}
		p := v1.CompatUsersUserIDPutParams{
			UserID: id,
		}
		return u.client.CompatUsersUserIDPut(ctx, &req, p)
	})
}

func (u *userOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.CompatUsersUserIDDeleteNoContent]("User.Delete", func() (any, error) {
		return u.client.CompatUsersUserIDDelete(ctx, v1.CompatUsersUserIDDeleteParams{UserID: id})
	})

	return err
}

func (u *userOp) RegisterEmail(ctx context.Context, userID int, email string) error {
	_, err := common.ErrorFromDecodedResponse[v1.CompatUsersUserIDRegisterEmailPostNoContent]("User.RegisterEmail", func() (any, error) {
		req := v1.CompatUsersUserIDRegisterEmailPostReq{Email: email}
		p := v1.CompatUsersUserIDRegisterEmailPostParams{UserID: userID}
		return u.client.CompatUsersUserIDRegisterEmailPost(ctx, &req, p)
	})
	return err
}

func (u *userOp) UnregisterEmail(ctx context.Context, userID int) error {
	_, err := common.ErrorFromDecodedResponse[v1.CompatUsersUserIDUnregisterEmailPostNoContent]("User.UnregisterEmail", func() (any, error) {
		return u.client.CompatUsersUserIDUnregisterEmailPost(ctx, v1.CompatUsersUserIDUnregisterEmailPostParams{
			UserID: userID,
		})
	})
	return err
}
