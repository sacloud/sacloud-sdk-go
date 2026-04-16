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

package auth_test

import (
	"net/http"
	"testing"

	. "github.com/sacloud/iam-api-go/apis/auth"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, AuthAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewAuthOp(client)
	return assert, api
}

func TestNewAuthOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestGetPasswordPolicy(t *testing.T) {
	var expected v1.PasswordPolicy
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.ReadPasswordPolicy(t.Context())
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestGetPasswordPolicy_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ReadPasswordPolicy(t.Context())
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestPutPasswordPolicy(t *testing.T) {
	var expected v1.PasswordPolicy
	expected.SetFake()
	assert, api := setup(t, &expected)
	var req v1.PasswordPolicy
	req.SetFake()

	actual, err := api.UpdatePasswordPolicy(t.Context(), req)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestPutPasswordPolicy_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)
	var req v1.PasswordPolicy
	req.SetFake()

	actual, err := api.UpdatePasswordPolicy(t.Context(), req)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGetAuthConditions(t *testing.T) {
	var expected v1.AuthConditions
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.ReadAuthConditions(t.Context())
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestGetAuthConditions_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ReadAuthConditions(t.Context())
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestPutAuthConditions(t *testing.T) {
	var expected v1.AuthConditions
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.UpdateAuthConditions(t.Context(), &expected)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestPutAuthConditions_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	var req v1.AuthConditions
	req.SetFake()
	actual, err := api.UpdateAuthConditions(t.Context(), &req)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestReadAuthContext(t *testing.T) {
	var expected v1.GetAuthContextOK
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.ReadAuthContext(t.Context())
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestReadAuthContext_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ReadAuthContext(t.Context())
	assert.Error(err)
	assert.Nil(actual)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	op := NewAuthOp(client)

	pp, err := op.ReadPasswordPolicy(t.Context())
	assert.NoError(err)
	assert.NotNil(pp)

	_, err = op.UpdatePasswordPolicy(t.Context(), *pp)
	assert.NoError(err)

	ac, err := op.ReadAuthConditions(t.Context())
	assert.NoError(err)
	assert.NotNil(ac)

	_, err = op.UpdateAuthConditions(t.Context(), ac)
	assert.NoError(err)

	authContext, err := op.ReadAuthContext(t.Context())
	assert.NoError(err)
	assert.NotNil(authContext)
}
