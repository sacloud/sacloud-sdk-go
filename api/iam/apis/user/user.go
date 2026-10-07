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

	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

// UserAPI is the interface for user operations.
type UserAPI interface {
	List(ctx context.Context, params ListParams) (*v1.ListUsersOK, error)
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
	Ordering *v1.ListUsersOrdering
}

func (u *userOp) List(ctx context.Context, params ListParams) (*v1.ListUsersOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListUsersOK]("User.List", func() (any, error) {
		return u.client.ListUsers(ctx, v1.ListUsersParams{
			Page:     into.Opt[v1.OptInt](params.Page),
			PerPage:  into.Opt[v1.OptInt](params.PerPage),
			Ordering: into.Opt[v1.OptListUsersOrdering](params.Ordering),
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
		return u.client.CreateUser(ctx, &v1.CreateUserReq{
			Name:        params.Name,
			Password:    params.Password,
			Code:        params.Code,
			Description: params.Description,
			Email:       into.Opt[v1.OptString](params.Email),
		})
	})
}

func (u *userOp) Read(ctx context.Context, id int) (*v1.User, error) {
	return common.ErrorFromDecodedResponse[v1.User]("User.Read", func() (any, error) {
		return u.client.ReadUser(ctx, v1.ReadUserParams{UserID: id})
	})
}

type UpdateParams struct {
	Name        string
	Password    *string
	Description string
}

func (u *userOp) Update(ctx context.Context, id int, params UpdateParams) (*v1.User, error) {
	return common.ErrorFromDecodedResponse[v1.User]("User.Update", func() (any, error) {
		req := v1.UpdateUserReq{
			Name:        params.Name,
			Password:    into.Opt[v1.OptString](params.Password),
			Description: params.Description,
		}
		p := v1.UpdateUserParams{
			UserID: id,
		}
		return u.client.UpdateUser(ctx, &req, p)
	})
}

func (u *userOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteUserNoContent]("User.Delete", func() (any, error) {
		return u.client.DeleteUser(ctx, v1.DeleteUserParams{UserID: id})
	})

	return err
}

func (u *userOp) RegisterEmail(ctx context.Context, userID int, email string) error {
	_, err := common.ErrorFromDecodedResponse[v1.RegisterEmailNoContent]("User.RegisterEmail", func() (any, error) {
		req := v1.RegisterEmailReq{Email: email}
		p := v1.RegisterEmailParams{UserID: userID}
		return u.client.RegisterEmail(ctx, &req, p)
	})
	return err
}

func (u *userOp) UnregisterEmail(ctx context.Context, userID int) error {
	_, err := common.ErrorFromDecodedResponse[v1.UnregisterEmailNoContent]("User.UnregisterEmail", func() (any, error) {
		return u.client.UnregisterEmail(ctx, v1.UnregisterEmailParams{
			UserID: userID,
		})
	})
	return err
}
