// Copyright 2025- The sacloud/iam-api-go Authors
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

package folder

import (
	"context"

	v1 "github.com/sacloud/sacloud-sdk-go/api/iam/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/api/iam/common"
	"github.com/sacloud/sacloud-sdk-go/common/packages/into"
)

type FolderAPI interface {
	List(ctx context.Context, params ListParams) (*v1.ListFoldersOK, error)

	Create(ctx context.Context, params CreateParams) (*v1.Folder, error)
	Read(ctx context.Context, id int) (*v1.Folder, error)
	Update(ctx context.Context, id int, name string, description *string) (*v1.Folder, error)
	Delete(ctx context.Context, id int) error

	Move(ctx context.Context, ids []int, parent *int) error
}

type folderOp struct {
	client *v1.Client
}

var _ FolderAPI = (*folderOp)(nil)

func NewFolderOp(client *v1.Client) FolderAPI { return &folderOp{client} }

type ListParams struct {
	Page     *int
	PerPage  *int
	Name     *string
	ParentID *int
}

func (f *folderOp) List(ctx context.Context, params ListParams) (*v1.ListFoldersOK, error) {
	return common.ErrorFromDecodedResponse[v1.ListFoldersOK]("Folder.List", func() (any, error) {
		return f.client.ListFolders(ctx, v1.ListFoldersParams{
			Page:       into.Opt[v1.OptInt](params.Page),
			PerPage:    into.Opt[v1.OptInt](params.PerPage),
			FolderName: into.Opt[v1.OptString](params.Name),
			ParentID:   into.Opt[v1.OptInt](params.ParentID),
		})
	})
}

type CreateParams struct {
	Name        string
	Description *string
	ParentID    *int
}

func (f *folderOp) Create(ctx context.Context, params CreateParams) (*v1.Folder, error) {
	return common.ErrorFromDecodedResponse[v1.Folder]("Folder.Create", func() (any, error) {
		return f.client.CreateFolder(ctx, &v1.CreateFolderReq{
			Name:        params.Name,
			Description: into.Opt[v1.OptString](params.Description),
			ParentID:    into.Opt[v1.OptNilInt](params.ParentID),
		})
	})
}

func (f *folderOp) Read(ctx context.Context, id int) (*v1.Folder, error) {
	return common.ErrorFromDecodedResponse[v1.Folder]("Folder.Read", func() (any, error) {
		return f.client.ReadFolder(ctx, v1.ReadFolderParams{FolderID: id})
	})
}

func (f *folderOp) Update(ctx context.Context, id int, name string, description *string) (*v1.Folder, error) {
	return common.ErrorFromDecodedResponse[v1.Folder]("Folder.Update", func() (any, error) {
		params := v1.UpdateFolderParams{
			FolderID: id,
		}
		request := v1.UpdateFolderReq{
			Name:        name,
			Description: into.Opt[v1.OptString](description),
		}
		return f.client.UpdateFolder(ctx, &request, params)
	})
}

func (f *folderOp) Delete(ctx context.Context, folderID int) error {
	_, err := common.ErrorFromDecodedResponse[v1.DeleteFolderNoContent]("Folder.Delete", func() (any, error) {
		return f.client.DeleteFolder(ctx, v1.DeleteFolderParams{FolderID: folderID})
	})

	return err
}

func (f *folderOp) Move(ctx context.Context, ids []int, parent *int) error {
	_, err := common.ErrorFromDecodedResponse[v1.MoveFoldersNoContent]("Folder.Move", func() (any, error) {
		return f.client.MoveFolders(ctx, &v1.MoveFolders{
			FolderIds: ids,
			ParentID:  into.Nil[v1.NilInt](parent),
		})
	})

	return err
}
