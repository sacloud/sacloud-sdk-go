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

// Package group provides the GroupAPI that wraps the generated v1 client.
package group

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

// GroupAPI is the interface for group operations.
type GroupAPI interface {
	List(ctx context.Context, params ListParams) (*v1.ListGroupsOK, error)
	Create(ctx context.Context, name string, description string) (*v1.Group, error)
	Read(ctx context.Context, id int) (*v1.Group, error)
	Update(ctx context.Context, id int, name string, description string) (*v1.Group, error)
	Delete(ctx context.Context, id int) error

	ReadMemberships(ctx context.Context, groupID int) ([]v1.GroupMembershipsCompatUsersItem, error)
	UpdateMemberships(ctx context.Context, groupID int, userIDs []int) ([]v1.GroupMembershipsCompatUsersItem, error)
}

type groupOp struct {
	client *v1.Client
}

func NewGroupOp(client *v1.Client) GroupAPI { return &groupOp{client: client} }

type ListParams struct {
	Page     *int
	PerPage  *int
	Ordering *v1.ListGroupsOrdering
	User     *v1.User
}

func (g *groupOp) List(ctx context.Context, params ListParams) (*v1.ListGroupsOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListGroupsOK]("Group.List", func() (any, error) {
		var userID *int = nil
		if params.User != nil {
			uid := params.User.GetID()
			userID = &uid
		}

		return g.client.ListGroups(ctx, v1.ListGroupsParams{
			Page:         into.Opt[v1.OptInt](params.Page),
			PerPage:      into.Opt[v1.OptInt](params.PerPage),
			Ordering:     into.Opt[v1.OptListGroupsOrdering](params.Ordering),
			CompatUserID: into.Opt[v1.OptInt](userID),
		})
	})
}

func (g *groupOp) Create(ctx context.Context, name string, description string) (*v1.Group, error) {
	return common.ErrorFromDecodedResponse[v1.Group]("Group.Create", func() (any, error) {
		return g.client.CreateGroup(ctx, &v1.CreateGroupReq{
			Name:        name,
			Description: description,
		})
	})
}

func (g *groupOp) Read(ctx context.Context, id int) (*v1.Group, error) {
	return common.ErrorFromDecodedResponse[v1.Group]("Group.Read", func() (any, error) {
		return g.client.ReadGroup(ctx, v1.ReadGroupParams{GroupID: id})
	})
}

func (g *groupOp) Update(ctx context.Context, id int, name string, description string) (*v1.Group, error) {
	return common.ErrorFromDecodedResponse[v1.Group]("Group.Update", func() (any, error) {
		req := v1.UpdateGroupReq{
			Name:        name,
			Description: description,
		}
		p := v1.UpdateGroupParams{
			GroupID: id,
		}
		return g.client.UpdateGroup(ctx, &req, p)
	})
}

func (g *groupOp) Delete(ctx context.Context, id int) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteGroupNoContent]("Group.Delete", func() (any, error) {
		return g.client.DeleteGroup(ctx, v1.DeleteGroupParams{GroupID: id})
	})

	return err
}

func (g *groupOp) ReadMemberships(ctx context.Context, id int) ([]v1.GroupMembershipsCompatUsersItem, error) {
	if ret, err := common.ErrorFromDecodedResponse[v1.GroupMemberships]("Group.ReadMemberships", func() (any, error) {
		return g.client.ListMemberships(ctx, v1.ListMembershipsParams{GroupID: id})
	}); err != nil {
		return nil, err
	} else {
		return ret.CompatUsers, nil
	}
}

func (g *groupOp) UpdateMemberships(ctx context.Context, groupID int, userIDs []int) ([]v1.GroupMembershipsCompatUsersItem, error) {
	if ret, err := common.ErrorFromDecodedResponse[v1.GroupMemberships]("Group.UpdateMemberships", func() (any, error) {
		compatUsers := make([]v1.UpdateMembershipsReqCompatUsersItem, len(userIDs))
		for i, j := range userIDs {
			compatUsers[i].ID = j
		}
		req := v1.UpdateMembershipsReq{CompatUsers: compatUsers}
		p := v1.UpdateMembershipsParams{GroupID: groupID}
		return g.client.UpdateMemberships(ctx, &req, p)
	}); err != nil {
		return nil, err
	} else {
		return ret.CompatUsers, nil
	}
}
