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

package project_test

import (
	"net/http"
	"testing"

	"github.com/sacloud/iam-api-go"
	"github.com/sacloud/iam-api-go/apis/folder"
	. "github.com/sacloud/iam-api-go/apis/project"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, ProjectAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewProjectOp(client)
	return assert, api
}

func TestNewProjectOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestList(t *testing.T) {
	var expected v1.ProjectsGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.Project, 2))
	expected.Items[0].SetFake()
	expected.Items[1].SetFake()
	assert, api := setup(t, &expected)

	params := ListParams{}
	actual, err := api.List(t.Context(), params)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestList_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	params := ListParams{}
	actual, err := api.List(t.Context(), params)
	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestCreate(t *testing.T) {
	var expected v1.Project
	expected.SetFake()
	assert, api := setup(t, &expected, http.StatusCreated)

	params := CreateParams{
		Code:        testutil.RandomName("proj", 12, testutil.CharSetAlphaNum),
		Name:        testutil.RandomName("project", 32, testutil.CharSetAlphaNum),
		Description: testutil.Random(64, testutil.CharSetAlphaNum),
	}
	actual, err := api.Create(t.Context(), params)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestCreate_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	params := CreateParams{
		Code:        testutil.RandomName("proj", 12, testutil.CharSetAlphaNum),
		Name:        testutil.RandomName("project", 32, testutil.CharSetAlphaNum),
		Description: testutil.Random(64, testutil.CharSetAlphaNum),
	}
	actual, err := api.Create(t.Context(), params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGet(t *testing.T) {
	var expected v1.Project
	expected.SetFake()
	assert, api := setup(t, &expected)

	projectID := 1
	actual, err := api.Read(t.Context(), projectID)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestGet_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	projectID := 1
	actual, err := api.Read(t.Context(), projectID)
	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestUpdate(t *testing.T) {
	var expected v1.Project
	name := testutil.RandomName("project", 32, testutil.CharSetAlphaNum)
	description := testutil.Random(64, testutil.CharSetAlphaNum)
	expected.SetFake()
	expected.SetName(name)
	expected.SetDescription(description)
	assert, api := setup(t, &expected)

	projectID := 1
	actual, err := api.Update(t.Context(), projectID, name, description)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestUpdate_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	projectID := 1
	name := testutil.RandomName("project", 32, testutil.CharSetAlphaNum)
	description := testutil.Random(64, testutil.CharSetAlphaNum)
	actual, err := api.Update(t.Context(), projectID, name, description)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, &v1.ProjectsProjectIDDeleteNoContent{}, http.StatusNoContent)

	projectID := 1
	err := api.Delete(t.Context(), projectID)
	assert.NoError(err)
}

func TestDelete_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	projectID := 1
	err := api.Delete(t.Context(), projectID)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestMove(t *testing.T) {
	assert, api := setup(t, &v1.MoveProjectsPostNoContent{}, http.StatusNoContent)

	ids := []int{1, 2, 3}
	parentFolderID := 1
	err := api.Move(t.Context(), ids, &parentFolderID)
	assert.NoError(err)
}

func TestMove_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	ids := []int{1, 2, 3}
	parentFolderID := 1
	err := api.Move(t.Context(), ids, &parentFolderID)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	op := NewProjectOp(client)

	// Create
	createParams := CreateParams{
		Code:        testutil.RandomName("proj", 12, testutil.CharSetAlphaNum),
		Name:        testutil.RandomName("project", 32, testutil.CharSetAlphaNum),
		Description: testutil.Random(64, testutil.CharSetAlphaNum),
	}
	created, err := op.Create(t.Context(), createParams)
	assert.NoError(err)
	assert.NotNil(created)

	defer func() {
		// Delete
		err = op.Delete(t.Context(), created.GetID())
		assert.NoError(err)
	}()

	// Read
	read, err := op.Read(t.Context(), created.GetID())
	assert.NoError(err)
	assert.NotNil(read)
	assert.Equal(created, read)

	// Update
	newName := testutil.RandomName("project-updated", 32, testutil.CharSetAlphaNum)
	newDescription := testutil.Random(64, testutil.CharSetAlphaNum)
	updated, err := op.Update(t.Context(), created.GetID(), newName, newDescription)
	assert.NoError(err)
	assert.NotNil(updated)
	assert.Equal(newName, updated.GetName())
	assert.Equal(newDescription, updated.GetDescription())

	fop := iam.NewFolderOp(client)
	folderName := testutil.RandomName("folder", 16, testutil.CharSetAlphaNum)
	folder, err := fop.Create(t.Context(), folder.CreateParams{Name: folderName})
	assert.NoError(err)
	assert.NotNil(folder)

	defer func() {
		// Delete Folder
		err = fop.Delete(t.Context(), folder.ID)
		assert.NoError(err)
	}()

	// Move
	err = op.Move(t.Context(), []int{created.GetID()}, &folder.ID)
	assert.NoError(err)

	defer func() {
		// Move back to root folder
		err = op.Move(t.Context(), []int{created.GetID()}, nil)
		assert.NoError(err)
	}()

	// List
	listParams := ListParams{}
	listed, err := op.List(t.Context(), listParams)
	assert.NoError(err)
	assert.NotNil(listed)
}
