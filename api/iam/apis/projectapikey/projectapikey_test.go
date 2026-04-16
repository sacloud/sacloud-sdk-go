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

package projectapikey_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/sacloud/iam-api-go"
	"github.com/sacloud/iam-api-go/apis/projectapikey"
	. "github.com/sacloud/iam-api-go/apis/projectapikey"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, ProjectAPIKeyAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewProjectAPIKeyOp(client)
	return assert, api
}

var Time = time.UnixMicro(0).UTC()

func TestNewProjectAPIKeyOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestList(t *testing.T) {
	var expected v1.CompatAPIKeysGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.ProjectApiKey, 1))
	expected.Items[0].SetFake()
	expected.Items[0].SetCreatedAt(v1.NewOptString(Time.String()))
	expected.Items[0].SetUpdatedAt(v1.NewOptString(Time.String()))
	expected.Items[0].SetIamRoles(make([]string, 1))
	expected.Items[0].IamRoles[0] = "role1"
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

	actual, err := api.List(t.Context(), ListParams{})
	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestCreate(t *testing.T) {
	var expected v1.ProjectApiKeyWithSecret
	expected.SetFake()
	expected.SetIamRoles([]string{"foo", "bar"})
	expected.SetCreatedAt(v1.NewOptString(Time.String()))
	expected.SetUpdatedAt(v1.NewOptString(Time.String()))
	assert, api := setup(t, &expected, http.StatusCreated)

	params := CreateParams{
		ProjectID:   123,
		Name:        testutil.RandomName("key", 32, testutil.CharSetAlphaNum),
		Description: testutil.Random(64, testutil.CharSetAlphaNum),
		IamRoles:    []string{"foo", "bar"},
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
		ProjectID:   123,
		Name:        testutil.RandomName("key", 32, testutil.CharSetAlphaNum),
		Description: testutil.Random(64, testutil.CharSetAlphaNum),
		IamRoles:    []string{"foo", "bar"},
	}
	actual, err := api.Create(t.Context(), params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGet(t *testing.T) {
	var expected v1.ProjectApiKey
	expected.SetFake()
	expected.SetCreatedAt(v1.NewOptString(Time.String()))
	expected.SetUpdatedAt(v1.NewOptString(Time.String()))
	expected.SetIamRoles([]string{"role1", "role2"})
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
	assert.Contains(err.Error(), expected)
}

func TestUpdate(t *testing.T) {
	var expected v1.ProjectApiKey
	expected.SetFake()
	expected.SetCreatedAt(v1.NewOptString(Time.String()))
	expected.SetUpdatedAt(v1.NewOptString(Time.String()))
	expected.SetIamRoles([]string{"role1", "role2"})
	assert, api := setup(t, &expected)

	params := UpdateParams{
		Name:        "foo",
		Description: "bar",
		IamRoles:    []string{"role1", "role2"},
	}
	actual, err := api.Update(t.Context(), 123, params)
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

	name := testutil.RandomName("key", 32, testutil.CharSetAlphaNum)
	description := testutil.Random(64, testutil.CharSetAlphaNum)
	params := UpdateParams{
		Name:        name,
		Description: description,
		IamRoles:    []string{"role1", "role2"},
	}
	actual, err := api.Update(t.Context(), 123, params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, &v1.CompatAPIKeysApikeyIDDeleteNoContent{}, http.StatusNoContent)

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

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	api := NewProjectAPIKeyOp(client)
	svc := iam.NewServicePrincipalOp(client)

	// an api key must be created by a service principal of the project
	myself := iam_test.Myself()
	p, err := svc.Read(t.Context(), myself.PrincipalID)
	assert.NoError(err)

	// Create
	created, err := api.Create(t.Context(), projectapikey.CreateParams{
		ProjectID:   p.GetProjectID(),
		Name:        testutil.RandomName("key-", 32, testutil.CharSetAlphaNum),
		Description: testutil.Random(64, testutil.CharSetAlphaNum),
		IamRoles:    []string{"resource-viewer"},
	})
	assert.NoError(err)
	assert.NotNil(created)

	// Delete
	defer func() {
		err = api.Delete(t.Context(), created.GetID())
		assert.NoError(err)
	}()

	// Read
	read, err := api.Read(t.Context(), created.GetID())
	assert.NoError(err)
	assert.NotNil(read)

	assert.Equal(created.GetID(), read.GetID())
	assert.Equal(created.GetProjectID(), read.GetProjectID())
	assert.Equal(created.GetName(), read.GetName())
	assert.Equal(created.GetDescription(), read.GetDescription())
	assert.Equal(created.GetAccessToken(), read.GetAccessToken())
	assert.Equal(created.GetServerResourceID(), read.GetServerResourceID())
	assert.Equal(created.GetIamRoles(), read.GetIamRoles())
	assert.Equal(created.GetZoneID(), read.GetZoneID())

	// List
	listed, err := api.List(t.Context(), ListParams{})
	assert.NoError(err)
	assert.NotNil(listed)
	assert.NotEmpty(listed)

	// Update
	updated, err := api.Update(t.Context(), created.GetID(), projectapikey.UpdateParams{
		Name:        read.GetName(),
		Description: read.GetDescription(),
		IamRoles:    []string{"resource-viewer"},
	})
	assert.NoError(err)
	assert.NotNil(updated)
}
