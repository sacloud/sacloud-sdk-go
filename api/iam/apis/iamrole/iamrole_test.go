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

package iamrole_test

import (
	"net/http"
	"testing"

	. "github.com/sacloud/iam-api-go/apis/iamrole"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, IAMRoleAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewIAMRoleOp(client)
	return assert, api
}

func TestNewIAMRoleOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestList(t *testing.T) {
	var expected v1.IamRolesGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.IamRole, 1))
	expected.Items[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.List(t.Context(), nil, nil)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestList_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := "forbidden"
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.List(t.Context(), nil, nil)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGet(t *testing.T) {
	var expected v1.IamRole
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), "123")
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestGet_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := "not found"
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.Read(t.Context(), "123")
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)

	op := NewIAMRoleOp(client)

	roles, err := op.List(t.Context(), nil, nil)
	assert.NoError(err)
	assert.NotEmpty(roles.Items)

	roleID := roles.Items[0].ID
	role, err := op.Read(t.Context(), roleID)
	assert.NoError(err)
	assert.Equal(roles.Items[0], *role)
}
