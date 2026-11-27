// Copyright 2025- The sacloud/apigw-api-go authors
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

package apigw

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	v1 "github.com/sacloud/sacloud-sdk-go/api/apigw/apis/v1"
)

type UserAPI interface {
	List(ctx context.Context) ([]v1.User, error)
	Create(ctx context.Context, request *v1.UserDetail) (*v1.UserDetail, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.UserDetail, error)
	Update(ctx context.Context, request *v1.UserDetail, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ UserAPI = (*userOp)(nil)

type userOp struct {
	client *v1.Client
}

func NewUserOp(client *v1.Client) UserAPI {
	return &userOp{client: client}
}

func (op *userOp) List(ctx context.Context) ([]v1.User, error) {
	res, err := errorFromDecodedResponse[v1.GetUsersOK]("User.List", func() (any, error) { return op.client.GetUsers(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetUsers(), nil
}

func (op *userOp) Create(ctx context.Context, request *v1.UserDetail) (*v1.UserDetail, error) {
	res, err := errorFromDecodedResponse[v1.AddUserCreated]("User.Create", func() (any, error) { return op.client.AddUser(ctx, request) })
	if err != nil {
		return nil, err
	}

	u, ok := new(res.GetApigw()).GetUser().Get()
	if !ok {
		return nil, NewError("User.Create", errors.New("response does not contain a user"))
	}

	return &u, nil
}

func (op *userOp) Read(ctx context.Context, id uuid.UUID) (*v1.UserDetail, error) {
	res, err := errorFromDecodedResponse[v1.GetUserOK]("User.Read", func() (any, error) { return op.client.GetUser(ctx, v1.GetUserParams{UserId: id}) })
	if err != nil {
		return nil, err
	}

	u, ok := new(res.GetApigw()).GetUser().Get()
	if !ok {
		return nil, NewError("User.Read", errors.New("response does not contain a user"))
	}

	return &u, nil
}

func (op *userOp) Update(ctx context.Context, request *v1.UserDetail, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.UpdateUserNoContent]("User.Update", func() (any, error) {
		return op.client.UpdateUser(ctx, request, v1.UpdateUserParams{UserId: id})
	})
}

func (op *userOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.DeleteUserNoContent]("User.Delete", func() (any, error) {
		return op.client.DeleteUser(ctx, v1.DeleteUserParams{UserId: id})
	})
}

type UserExtraAPI interface {
	ListGroup(ctx context.Context) ([]v1.UserGroupDetail, error)
	UpdateGroup(ctx context.Context, groupIdOrName string, isAssigned bool) error
	ReadAuth(ctx context.Context) (*v1.UserAuthentication, error)
	UpdateAuth(ctx context.Context, request v1.UserAuthentication) error
}

var _ UserExtraAPI = (*userExtraOp)(nil)

type userExtraOp struct {
	client *v1.Client
	userId uuid.UUID
}

func NewUserExtraOp(client *v1.Client, userId uuid.UUID) UserExtraAPI {
	return &userExtraOp{client: client, userId: userId}
}

func (op *userExtraOp) ListGroup(ctx context.Context) ([]v1.UserGroupDetail, error) {
	res, err := errorFromDecodedResponse[v1.GetUserGroupOK]("UserExtra.ListGroup", func() (any, error) { return op.client.GetUserGroup(ctx, v1.GetUserGroupParams{UserId: op.userId}) })
	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetGroups(), nil
}

func (op *userExtraOp) UpdateGroup(ctx context.Context, groupIdOrName string, isAssigned bool) error {
	var req []byte
	idOrName, err := uuid.Parse(groupIdOrName)
	if err != nil {
		temp := []struct {
			IsAssigned bool   `json:"isAssigned"`
			Name       string `json:"name"`
		}{{
			isAssigned,
			groupIdOrName,
		}}
		req, _ = json.Marshal(temp)
	} else {
		temp := []struct {
			IsAssigned bool      `json:"isAssigned"`
			Id         uuid.UUID `json:"id"`
		}{{
			isAssigned,
			idOrName,
		}}
		req, _ = json.Marshal(temp)
	}

	return errorFromDecodedResponseE[v1.UpdateUserGroupNoContent]("UserExtra.UpdateGroup", func() (any, error) {
		return op.client.UpdateUserGroup(ctx, req, v1.UpdateUserGroupParams{UserId: op.userId})
	})
}

func (op *userExtraOp) ReadAuth(ctx context.Context) (*v1.UserAuthentication, error) {
	res, err := errorFromDecodedResponse[v1.GetUserAuthenticationOK]("UserExtra.ReadAuth", func() (any, error) {
		return op.client.GetUserAuthentication(ctx, v1.GetUserAuthenticationParams{UserId: op.userId})
	})
	if err != nil {
		return nil, err
	}

	auth, ok := new(res.GetApigw()).GetUserAuthentication().Get()
	if !ok {
		return nil, NewError("UserExtra.ReadAuth", errors.New("response does not contain a user authentication"))
	}

	return &auth, nil
}

func (op *userExtraOp) UpdateAuth(ctx context.Context, request v1.UserAuthentication) error {
	return errorFromDecodedResponseE[v1.UpsertUserAuthenticationNoContent]("UserExtra.UpdateAuth", func() (any, error) {
		return op.client.UpsertUserAuthentication(ctx, v1.NewOptUserAuthentication(request), v1.UpsertUserAuthenticationParams{UserId: op.userId})
	})
}
