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

package group_test

import (
	"iter"
	"net/http"
	"slices"
	"testing"

	. "github.com/sacloud/iam-api-go/apis/group"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, GroupAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewGroupOp(client)
	return assert, api
}

func TestNewGroupOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestList(t *testing.T) {
	var expected v1.GroupsGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.Group, 2))
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
	var expected v1.Group
	expected.SetFake()
	assert, api := setup(t, &expected, http.StatusCreated)

	actual, err := api.Create(t.Context(), "name", "description")
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

	actual, err := api.Create(t.Context(), "name", "description")
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGet(t *testing.T) {
	var expected v1.Group
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
	assert.Contains(err.Error(), expected)
}

func TestUpdate(t *testing.T) {
	var expected v1.Group
	name := testutil.RandomName("group", 32, testutil.CharSetAlphaNum)
	description := testutil.Random(64, testutil.CharSetAlphaNum)
	expected.SetFake()
	expected.SetName(name)
	expected.SetDescription(description)
	assert, api := setup(t, &expected)

	actual, err := api.Update(t.Context(), 123, "name", "description")
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

	actual, err := api.Update(t.Context(), 123, "name", "description")
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, &v1.GroupsGroupIDDeleteNoContent{}, http.StatusNoContent)

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

func TestGetMemberships(t *testing.T) {
	var mock v1.GroupMemberships
	expected := make([]v1.GroupMembershipsCompatUsersItem, 2)
	expected[0].SetFake()
	expected[1].SetFake()
	mock.SetFake()
	mock.SetCompatUsers(expected)
	assert, api := setup(t, &mock)

	actual, err := api.ReadMemberships(t.Context(), 123)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected, actual)
}

func TestGetMemberships_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ReadMemberships(t.Context(), 123)
	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestUpdateMemberships(t *testing.T) {
	expected := make([]v1.GroupMembershipsCompatUsersItem, 2)
	expected[0].SetFake()
	expected[1].SetFake()
	var mock v1.GroupMemberships
	mock.SetFake()
	mock.SetCompatUsers(expected)
	assert, api := setup(t, &mock)

	actual, err := api.UpdateMemberships(t.Context(), 123, []int{1, 2})
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected, actual)
}

func TestUpdateMemberships_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.UpdateMemberships(t.Context(), 123, []int{1, 2})
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	op := NewGroupOp(client)

	// Create
	groupName := testutil.RandomName("group", 32, testutil.CharSetAlphaNum)
	groupDescription := testutil.Random(64, testutil.CharSetAlphaNum)
	created, err := op.Create(t.Context(), groupName, groupDescription)
	assert.NoError(err)
	assert.NotNil(created)

	id := created.GetID()

	defer func() {
		err = op.Delete(t.Context(), id)
		assert.NoError(err)
	}()

	// List
	groups, err := op.List(t.Context(), ListParams{})
	assert.NoError(err)
	assert.NotNil(groups)
	assert.NotEmpty(groups.GetItems())

	// Read
	read, err := op.Read(t.Context(), id)
	assert.NoError(err)
	assert.Equal(created, read)

	// Update
	newGroupName := testutil.RandomName("group-updated", 32, testutil.CharSetAlphaNum)
	newGroupDescription := testutil.Random(64, testutil.CharSetAlphaNum)
	updated, err := op.Update(t.Context(), id, newGroupName, newGroupDescription)
	assert.NoError(err)
	assert.NotNil(updated)
	assert.Equal(newGroupName, updated.GetName())
	assert.Equal(newGroupDescription, updated.GetDescription())

	user, deleter := iam_test.NewUser(t, client)
	defer deleter()

	// ReadMemberships
	memberships, err := op.ReadMemberships(t.Context(), id)
	assert.NoError(err)
	assert.NotNil(memberships)

	// UpdateMemberships
	ids := mapSeq(slices.Values(memberships), func(m v1.GroupMembershipsCompatUsersItem) int { return m.GetID() })
	updatedMemberships, err := op.UpdateMemberships(t.Context(), id, append(slices.Collect(ids), user.GetID()))
	assert.NoError(err)
	assert.NotNil(updatedMemberships)
}

func mapSeq[T any, U any](s iter.Seq[T], f func(T) U) iter.Seq[U] {
	return func(y func(U) bool) {
		for t := range s {
			if !y(f(t)) {
				return
			}
		}
	}
}
