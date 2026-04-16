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

package user_test

import (
	"net/http"
	"testing"

	. "github.com/sacloud/iam-api-go/apis/user"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, UserAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewUserOp(client)
	return assert, api
}

func TestNewUserOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestList(t *testing.T) {
	var expected v1.CompatUsersGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.User, 2))
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
	var expected v1.User
	expected.SetFake()
	assert, api := setup(t, &expected, http.StatusCreated)

	params := CreateParams{
		Name:     testutil.RandomName("user", 32, testutil.CharSetAlphaNum),
		Password: testutil.Random(16, testutil.CharSetAlphaNum),
		Email:    nil,
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
		Name:     testutil.RandomName("user", 32, testutil.CharSetAlphaNum),
		Password: testutil.Random(16, testutil.CharSetAlphaNum),
		Email:    nil,
	}
	actual, err := api.Create(t.Context(), params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGet(t *testing.T) {
	var expected v1.User
	expected.SetFake()
	assert, api := setup(t, &expected)

	userID := 1
	actual, err := api.Read(t.Context(), userID)
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

	userID := 1
	actual, err := api.Read(t.Context(), userID)
	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestUpdate(t *testing.T) {
	var expected v1.User
	name := testutil.RandomName("user", 32, testutil.CharSetAlphaNum)
	description := testutil.Random(64, testutil.CharSetAlphaNum)
	password := testutil.Random(16, testutil.CharSetAlphaNum)
	expected.SetFake()
	expected.SetName(name)
	expected.SetDescription(description)
	assert, api := setup(t, &expected)

	userID := 1
	params := UpdateParams{
		Name:        name,
		Password:    &password,
		Description: description,
	}
	actual, err := api.Update(t.Context(), userID, params)
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

	userID := 1
	name := testutil.RandomName("user", 32, testutil.CharSetAlphaNum)
	description := testutil.Random(64, testutil.CharSetAlphaNum)
	password := testutil.Random(16, testutil.CharSetAlphaNum)
	params := UpdateParams{
		Name:        name,
		Password:    &password,
		Description: description,
	}
	actual, err := api.Update(t.Context(), userID, params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, &v1.CompatUsersUserIDDeleteNoContent{}, http.StatusNoContent)

	userID := 1
	err := api.Delete(t.Context(), userID)
	assert.NoError(err)
}

func TestDelete_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	userID := 1
	err := api.Delete(t.Context(), userID)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestRegisterEmail(t *testing.T) {
	assert, api := setup(t, &v1.CompatUsersUserIDRegisterEmailPostNoContent{}, http.StatusNoContent)

	userID := 1
	email := testutil.RandomName("name-", 12, testutil.CharSetAlphaNum) + "@example.com"
	err := api.RegisterEmail(t.Context(), userID, email)
	assert.NoError(err)
}

func TestRegisterEmail_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	userID := 1
	email := testutil.RandomName("name-", 12, testutil.CharSetAlphaNum) + "@example.com"
	err := api.RegisterEmail(t.Context(), userID, email)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestUnregisterEmail(t *testing.T) {
	assert, api := setup(t, &v1.CompatUsersUserIDUnregisterEmailPostNoContent{}, http.StatusNoContent)

	userID := 1
	err := api.UnregisterEmail(t.Context(), userID)
	assert.NoError(err)
}

func TestUnregisterEmail_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	userID := 1
	err := api.UnregisterEmail(t.Context(), userID)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	api := NewUserOp(client)

	// Create
	name := testutil.RandomName("user", 32, testutil.CharSetAlphaNum)
	password := testutil.Random(64, testutil.CharSetAlphaNum)
	code := testutil.RandomName("c", 31, testutil.CharSetAlphaNum)
	createParams := CreateParams{
		Name:     name,
		Password: password,
		Code:     code,
		Email:    nil,
	}
	created, err := api.Create(t.Context(), createParams)
	assert.NoError(err)
	assert.NotNil(created)
	assert.Equal(name, created.GetName())

	defer func() {
		// Delete
		err = api.Delete(t.Context(), created.GetID())
		assert.NoError(err)
	}()

	// Read
	read, err := api.Read(t.Context(), created.GetID())
	assert.NoError(err)
	assert.NotNil(read)
	assert.Equal(created, read)

	// Update
	newDescription := testutil.Random(64, testutil.CharSetAlphaNum)
	updateParams := UpdateParams{
		Name:        name,
		Password:    &password,
		Description: newDescription,
	}
	updated, err := api.Update(t.Context(), created.GetID(), updateParams)
	assert.NoError(err)
	assert.NotNil(updated)
	assert.Equal(newDescription, updated.GetDescription())

	// RegisterEmail
	email := testutil.RandomName("name-", 12, testutil.CharSetAlphaNum) + "@example.com"
	err = api.RegisterEmail(t.Context(), created.GetID(), email)
	assert.NoError(err)

	// UnregisterEmail
	err = api.UnregisterEmail(t.Context(), created.GetID())
	assert.NoError(err)
}
