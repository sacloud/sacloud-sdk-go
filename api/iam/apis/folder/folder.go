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

	v1 "github.com/sacloud/iam-api-go/apis/v1"
	"github.com/sacloud/iam-api-go/common"
)

type FolderAPI interface {
	List(ctx context.Context, params ListParams) (*v1.FoldersGetOK, error)

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

func (f *folderOp) List(ctx context.Context, params ListParams) (*v1.FoldersGetOK, error) {
	return common.ErrorFromDecodedResponse[v1.FoldersGetOK]("Folder.List", func() (any, error) {
		return f.client.FoldersGet(ctx, v1.FoldersGetParams{
			Page:       common.IntoOpt[v1.OptInt](params.Page),
			PerPage:    common.IntoOpt[v1.OptInt](params.PerPage),
			FolderName: common.IntoOpt[v1.OptString](params.Name),
			ParentID:   common.IntoOpt[v1.OptInt](params.ParentID),
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
		return f.client.FoldersPost(ctx, &v1.FoldersPostReq{
			Name:        params.Name,
			Description: common.IntoOpt[v1.OptString](params.Description),
			ParentID:    common.IntoOpt[v1.OptNilInt](params.ParentID),
		})
	})
}

func (f *folderOp) Read(ctx context.Context, id int) (*v1.Folder, error) {
	return common.ErrorFromDecodedResponse[v1.Folder]("Folder.Read", func() (any, error) {
		return f.client.FoldersFolderIDGet(ctx, v1.FoldersFolderIDGetParams{FolderID: id})
	})
}

func (f *folderOp) Update(ctx context.Context, id int, name string, description *string) (*v1.Folder, error) {
	return common.ErrorFromDecodedResponse[v1.Folder]("Folder.Update", func() (any, error) {
		params := v1.FoldersFolderIDPutParams{
			FolderID: id,
		}
		request := v1.FoldersFolderIDPutReq{
			Name:        name,
			Description: common.IntoOpt[v1.OptString](description),
		}
		return f.client.FoldersFolderIDPut(ctx, &request, params)
	})
}

func (f *folderOp) Delete(ctx context.Context, folderID int) error {
	_, err := common.ErrorFromDecodedResponse[v1.FoldersFolderIDDeleteNoContent]("Folder.Delete", func() (any, error) {
		return f.client.FoldersFolderIDDelete(ctx, v1.FoldersFolderIDDeleteParams{FolderID: folderID})
	})

	return err
}

func (f *folderOp) Move(ctx context.Context, ids []int, parent *int) error {
	_, err := common.ErrorFromDecodedResponse[v1.MoveFoldersPostNoContent]("Folder.Move", func() (any, error) {
		return f.client.MoveFoldersPost(ctx, &v1.MoveFolders{
			FolderIds: ids,
			ParentID:  common.IntoNullable[v1.NilInt](parent),
		})
	})

	return err
}
