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

package folder_test

import (
	"net/http"
	"testing"

	. "github.com/sacloud/iam-api-go/apis/folder"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, FolderAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewFolderOp(client)
	return assert, api
}

func TestNewFolderOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestCreate(t *testing.T) {
	var expected v1.Folder
	expected.SetFake()
	assert, api := setup(t, &expected, http.StatusCreated)

	actual, err := api.Create(t.Context(), CreateParams{Name: expected.GetName()})
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetName(), actual.GetName())
}

func TestCreate_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.Create(t.Context(), CreateParams{Name: "fake"})
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestList(t *testing.T) {
	var expected v1.FoldersGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.Folder, 1))
	expected.Items[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.List(t.Context(), ListParams{})
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

	params := ListParams{Page: new(int)}
	actual, err := api.List(t.Context(), params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGet(t *testing.T) {
	var expected v1.Folder
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), 123)
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

	actual, err := api.Read(t.Context(), 123)
	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestUpdate(t *testing.T) {
	var expected v1.Folder
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Update(t.Context(), 123, "name", nil)
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

	actual, err := api.Update(t.Context(), 123, "name", nil)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	err := api.Delete(t.Context(), 123)
	assert.NoError(err)
}

func TestDelete_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	err := api.Delete(t.Context(), 123)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestMove(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	err := api.Move(t.Context(), []int{123}, nil)
	assert.NoError(err)
}

func TestMove_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	err := api.Move(t.Context(), []int{123}, nil)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	op := NewFolderOp(client)

	name1 := testutil.RandomName("folder", 64, testutil.CharSetAlphaNum)
	folder1, err := op.Create(t.Context(), CreateParams{Name: name1})
	assert.NoError(err)
	assert.NotNil(folder1)
	assert.Equal(name1, folder1.GetName())
	defer func() {
		err := op.Delete(t.Context(), folder1.GetID())
		assert.NoError(err)
	}()

	name2 := testutil.RandomName("folder", 64, testutil.CharSetAlphaNum)
	folder2, err := op.Create(t.Context(), CreateParams{Name: name2})
	assert.NoError(err)
	assert.NotNil(folder2)
	assert.Equal(name2, folder2.GetName())
	defer func() {
		err := op.Delete(t.Context(), folder2.GetID())
		assert.NoError(err)
	}()

	list, err := op.List(t.Context(), ListParams{})
	assert.NoError(err)
	assert.NotEmpty(list.Items)

	readFolder, err := op.Read(t.Context(), folder1.GetID())
	assert.NoError(err)
	assert.NotNil(readFolder)
	assert.Equal(folder1, readFolder)

	updatedName := testutil.RandomName("folder-updated", 64, testutil.CharSetAlphaNum)
	updatedFolder, err := op.Update(t.Context(), folder1.GetID(), updatedName, nil)
	assert.NoError(err)
	assert.NotNil(updatedFolder)
	assert.Equal(updatedName, updatedFolder.GetName())

	defer func() {
		err = op.Move(t.Context(), []int{folder1.GetID()}, nil)
		assert.NoError(err)
	}()
	err = op.Move(t.Context(), []int{folder1.GetID()}, saclient.Ptr(folder2.GetID()))
	assert.NoError(err)
}
