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
	v1 "github.com/sacloud/sacloud-sdk-go/api/apigw/apis/v1"
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
	res, err := errorFromDecodedResponse[v1.GetGroupsOK]("Group.List", func() (any, error) { return op.client.GetGroups(ctx) })
	if err != nil {
		return nil, err
	}

	return new(res.GetApigw()).GetGroups(), nil
}

func (op *groupOp) Create(ctx context.Context, request *v1.Group) (*v1.Group, error) {
	res, err := errorFromDecodedResponse[v1.AddGroupCreated]("Group.Create", func() (any, error) { return op.client.AddGroup(ctx, request) })
	if err != nil {
		return nil, err
	}

	g, ok := new(res.GetApigw()).GetGroup().Get()
	if !ok {
		return nil, NewError("Group.Create", errors.New("response does not contain a group"))
	}

	return &g, nil
}

func (op *groupOp) Read(ctx context.Context, id uuid.UUID) (*v1.Group, error) {
	res, err := errorFromDecodedResponse[v1.GetGroupOK]("Group.Read", func() (any, error) { return op.client.GetGroup(ctx, v1.GetGroupParams{GroupId: id}) })
	if err != nil {
		return nil, err
	}

	g, ok := new(res.GetApigw()).GetGroup().Get()
	if !ok {
		return nil, NewError("Group.Read", errors.New("response does not contain a group"))
	}

	return &g, nil
}

func (op *groupOp) Update(ctx context.Context, request *v1.Group, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.UpdateGroupNoContent]("Group.Update", func() (any, error) {
		return op.client.UpdateGroup(ctx, request, v1.UpdateGroupParams{GroupId: id})
	})
}

func (op *groupOp) Delete(ctx context.Context, id uuid.UUID) error {
	return errorFromDecodedResponseE[v1.DeleteGroupNoContent]("Group.Delete", func() (any, error) {
		return op.client.DeleteGroup(ctx, v1.DeleteGroupParams{GroupId: id})
	})
}
