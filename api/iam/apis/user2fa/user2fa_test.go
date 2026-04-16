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

package user2fa_test

import (
	"math/rand/v2"
	"net/http"
	"testing"
	"time"

	. "github.com/sacloud/iam-api-go/apis/user2fa"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, User2FAAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	var user v1.User
	user.SetFake()
	user.SetID(rand.Int()) //nolint:gosec
	api := NewUser2FAOp(client, &user)
	return assert, api
}

var Time = time.UnixMicro(0).UTC()

func TestNewUser2FAOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestDeactivateOTP(t *testing.T) {
	assert, api := setup(t, &v1.CompatUsersUserIDDeactivateOtpPostNoContent{}, http.StatusNoContent)

	err := api.DeactivateOTP(t.Context())
	assert.NoError(err)
}

func TestDeactivateOTP_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	err := api.DeactivateOTP(t.Context())
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestListSecurityKeys(t *testing.T) {
	var expected v1.CompatUsersUserIDSecurityKeysGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.UserSecurityKey, 1))
	expected.Items[0].SetFake()
	expected.Items[0].SetID(123)
	expected.Items[0].SetRegisteredAt(Time)
	expected.Items[0].SetLastUsedAt(v1.NewNilDateTime(Time))
	assert, api := setup(t, &expected)

	actual, err := api.ListSecurityKeys(t.Context())
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestListSecurityKeys_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ListSecurityKeys(t.Context())
	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestGetSecurityKey(t *testing.T) {
	var expected v1.UserSecurityKey
	expected.SetFake()
	expected.SetID(456)
	expected.SetRegisteredAt(Time)
	expected.SetLastUsedAt(v1.NewNilDateTime(Time))
	assert, api := setup(t, &expected)

	actual, err := api.ReadSecurityKey(t.Context(), 123)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestGetSecurityKey_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ReadSecurityKey(t.Context(), 123)
	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestUpdateSecurityKey(t *testing.T) {
	var expected v1.UserSecurityKey
	name := testutil.RandomName("key", 32, testutil.CharSetAlphaNum)
	expected.SetFake()
	expected.SetName(name)
	expected.SetID(456)
	expected.SetRegisteredAt(Time)
	expected.SetLastUsedAt(v1.NewNilDateTime(Time))
	assert, api := setup(t, &expected)

	actual, err := api.UpdateSecurityKey(t.Context(), 123, name)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestUpdateSecurityKey_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	name := testutil.RandomName("key", 32, testutil.CharSetAlphaNum)
	actual, err := api.UpdateSecurityKey(t.Context(), 123, name)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestDeleteSecurityKey(t *testing.T) {
	assert, api := setup(t, &v1.CompatUsersUserIDSecurityKeysSecurityKeyIDDeleteNoContent{}, http.StatusNoContent)

	err := api.DeleteSecurityKey(t.Context(), 123)
	assert.NoError(err)
}

func TestDeleteSecurityKey_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	err := api.DeleteSecurityKey(t.Context(), 123)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestListTrustedDevices(t *testing.T) {
	var expected v1.CompatUsersUserIDTrustedDevicesGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.UserTrustedDevice, 1))
	expected.Items[0].SetFake()
	expected.Items[0].SetCreatedAt(Time)
	assert, api := setup(t, &expected)

	actual, err := api.ListTrustedDevices(t.Context())
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestListTrustedDevices_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ListTrustedDevices(t.Context())
	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestDeleteTrustedDevice(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	err := api.DeleteTrustedDevice(t.Context(), 123)
	assert.NoError(err)
}

func TestDeleteTrustedDevice_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	err := api.DeleteTrustedDevice(t.Context(), 123)
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

func TestClearTrustedDevices(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	err := api.ClearTrustedDevices(t.Context())
	assert.NoError(err)
}

func TestClearTrustedDevices_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	err := api.ClearTrustedDevices(t.Context())
	assert.Error(err)
	assert.Contains(err.Error(), expected)
}

// There is no TestIntegrated
// You can delete MFA but canot create one.
