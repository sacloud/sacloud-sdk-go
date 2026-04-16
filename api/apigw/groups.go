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
	"errors"

	"github.com/google/uuid"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
)

type GroupAPI interface {
	List(ctx context.Context) ([]v1.Group, error)
	Create(ctx context.Context, request *v1.Group) (*v1.Group, error)
	Read(ctx context.Context, id uuid.UUID) (*v1.Group, error)
	Update(ctx context.Context, request *v1.Group, id uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
}

var _ GroupAPI = (*groupOp)(nil)

type groupOp struct {
	client *v1.Client
}

func NewGroupOp(client *v1.Client) GroupAPI {
	return &groupOp{client: client}
}

func (op *groupOp) List(ctx context.Context) ([]v1.Group, error) {
	res, err := op.client.GetGroups(ctx)
	if err != nil {
		return nil, NewAPIError("Group.List", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetGroupsOK:
		return p.Apigw.Groups, nil
	case *v1.GetGroupsBadRequest:
		return nil, NewAPIError("Group.List", 400, errors.New(p.Message.Value))
	case *v1.GetGroupsInternalServerError:
		return nil, NewAPIError("Group.List", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Group.List", 0, nil)
}

func (op *groupOp) Create(ctx context.Context, request *v1.Group) (*v1.Group, error) {
	res, err := op.client.AddGroup(ctx, request)
	if err != nil {
		return nil, NewAPIError("Group.Create", 0, err)
	}

	switch p := res.(type) {
	case *v1.AddGroupCreated:
		return &p.Apigw.Group.Value, nil
	case *v1.AddGroupBadRequest:
		return nil, NewAPIError("Group.Create", 400, errors.New(p.Message.Value))
	case *v1.AddGroupConflict:
		return nil, NewAPIError("Group.Create", 409, errors.New(p.Message.Value))
	case *v1.AddGroupInternalServerError:
		return nil, NewAPIError("Group.Create", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Group.Create", 0, nil)
}

func (op *groupOp) Read(ctx context.Context, id uuid.UUID) (*v1.Group, error) {
	res, err := op.client.GetGroup(ctx, v1.GetGroupParams{GroupId: id})
	if err != nil {
		return nil, NewAPIError("Group.Read", 0, err)
	}

	switch p := res.(type) {
	case *v1.GetGroupOK:
		return &p.Apigw.Group.Value, nil
	case *v1.GetGroupBadRequest:
		return nil, NewAPIError("Group.Read", 400, errors.New(p.Message.Value))
	case *v1.GetGroupNotFound:
		return nil, NewAPIError("Group.Read", 404, errors.New(p.Message.Value))
	case *v1.GetGroupInternalServerError:
		return nil, NewAPIError("Group.Read", 500, errors.New(p.Message.Value))
	}

	return nil, NewAPIError("Group.Read", 0, nil)
}

func (op *groupOp) Update(ctx context.Context, request *v1.Group, id uuid.UUID) error {
	res, err := op.client.UpdateGroup(ctx, request, v1.UpdateGroupParams{GroupId: id})
	if err != nil {
		return NewAPIError("Group.Update", 0, err)
	}

	switch p := res.(type) {
	case *v1.UpdateGroupNoContent:
		return nil
	case *v1.UpdateGroupBadRequest:
		return NewAPIError("Group.Update", 400, errors.New(p.Message.Value))
	case *v1.UpdateGroupNotFound:
		return NewAPIError("Group.Update", 404, errors.New(p.Message.Value))
	case *v1.UpdateGroupInternalServerError:
		return NewAPIError("Group.Update", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Group.Update", 0, nil)
}

func (op *groupOp) Delete(ctx context.Context, id uuid.UUID) error {
	res, err := op.client.DeleteGroup(ctx, v1.DeleteGroupParams{GroupId: id})
	if err != nil {
		return NewAPIError("Group.Delete", 0, err)
	}

	switch p := res.(type) {
	case *v1.DeleteGroupNoContent:
		return nil
	case *v1.DeleteGroupBadRequest:
		return NewAPIError("Group.Delete", 400, errors.New(p.Message.Value))
	case *v1.DeleteGroupUnauthorized:
		return NewAPIError("Group.Delete", 401, errors.New(p.Message.Value))
	case *v1.DeleteGroupNotFound:
		return NewAPIError("Group.Delete", 404, errors.New(p.Message.Value))
	case *v1.DeleteGroupInternalServerError:
		return NewAPIError("Group.Delete", 500, errors.New(p.Message.Value))
	}

	return NewAPIError("Group.Delete", 0, nil)
}
