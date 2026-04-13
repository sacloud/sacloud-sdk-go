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

package scim_test

import (
	"net/http"
	"testing"

	"github.com/google/uuid"
	. "github.com/sacloud/iam-api-go/apis/scim"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, ScimAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewScimOp(client)
	return assert, api
}

func TestNewScimOp(t *testing.T) {
	assert := require.New(t)
	client := iam_test.NewTestClient(nil, 200)
	api := NewScimOp(client)
	assert.NotNil(api)
}

func intPtr(i int) *int {
	return &i
}

func TestList(t *testing.T) {
	var expected v1.ScimConfigurationsGetOK
	expected.SetFake()
	// Itemsフィールドは空のスライスになるように調整
	expected.Items = []v1.ScimConfigurationBase{}
	assert, api := setup(t, &expected, 200)

	params := ListParams{
		Page:    intPtr(1),
		PerPage: intPtr(10),
	}
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

	params := ListParams{
		Page:    intPtr(1),
		PerPage: intPtr(10),
	}
	actual, err := api.List(t.Context(), params)
	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestCreate(t *testing.T) {
	var expected v1.ScimConfiguration
	expected.SetFake()
	assert, api := setup(t, &expected, 201)

	params := CreateParams{
		Name: testutil.RandomName("scim", 32, testutil.CharSetAlphaNum),
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
		Name: testutil.RandomName("scim", 32, testutil.CharSetAlphaNum),
	}
	actual, err := api.Create(t.Context(), params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGet(t *testing.T) {
	var expected v1.ScimConfigurationBase
	expected.SetFake()
	assert, api := setup(t, &expected, 200)

	id := uuid.New().String()
	actual, err := api.Read(t.Context(), id)
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

	id := uuid.New().String()
	actual, err := api.Read(t.Context(), id)
	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestUpdate(t *testing.T) {
	var expected v1.ScimConfigurationBase
	expected.SetFake()
	assert, api := setup(t, &expected, 200)

	id := uuid.New().String()
	params := UpdateParams{
		Name: testutil.RandomName("scim", 32, testutil.CharSetAlphaNum),
	}
	actual, err := api.Update(t.Context(), id, params)
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

	id := uuid.New().String()
	params := UpdateParams{
		Name: testutil.RandomName("scim", 32, testutil.CharSetAlphaNum),
	}
	actual, err := api.Update(t.Context(), id, params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, nil, 204)

	id := uuid.New().String()
	err := api.Delete(t.Context(), id)
	assert.NoError(err)
}

func TestDelete_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	id := uuid.New().String()
	err := api.Delete(t.Context(), id)
	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestRegenerateToken(t *testing.T) {
	var expected v1.ScimConfigurationsIDRegenerateTokenPostOK
	expected.SetFake()
	assert, api := setup(t, &expected, 200)

	id := uuid.New().String()
	actual, err := api.RegenerateToken(t.Context(), id)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestRegenerateToken_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	id := uuid.New().String()
	actual, err := api.RegenerateToken(t.Context(), id)
	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	api := NewScimOp(client)

	// Create
	name := testutil.RandomName("scim", 32, testutil.CharSetAlphaNum)
	createParams := CreateParams{
		Name: name,
	}
	created, err := api.Create(t.Context(), createParams)
	assert.NoError(err)
	assert.NotNil(created)
	assert.Equal(name, created.GetName())

	defer func() {
		// Delete
		err = api.Delete(t.Context(), created.GetID().String())
		assert.NoError(err)
	}()

	// Read
	read, err := api.Read(t.Context(), created.GetID().String())
	assert.NoError(err)
	assert.NotNil(read)
	assert.Equal(created.GetID(), read.GetID())
	assert.Equal(created.GetName(), read.GetName())
	assert.Equal(created.GetBaseURL(), read.GetBaseURL())
	assert.Equal(created.GetCreatedAt(), read.GetCreatedAt())
	assert.Equal(created.GetUpdatedAt(), read.GetUpdatedAt())

	// Update
	newName := testutil.RandomName("scim", 32, testutil.CharSetAlphaNum)
	updateParams := UpdateParams{
		Name: newName,
	}
	updated, err := api.Update(t.Context(), created.GetID().String(), updateParams)
	assert.NoError(err)
	assert.NotNil(updated)
	assert.Equal(newName, updated.GetName())

	// List
	list, err := api.List(t.Context(), ListParams{})
	assert.NoError(err)
	assert.NotNil(list)
	assert.GreaterOrEqual(len(list.GetItems()), 1)

	// RegenerateToken
	token, err := api.RegenerateToken(t.Context(), created.GetID().String())
	assert.NoError(err)
	assert.NotNil(token)
	assert.NotEmpty(token.GetSecretToken())
}
