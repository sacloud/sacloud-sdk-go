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

package sso_test

import (
	"encoding/pem"
	"net/http"
	"testing"

	"github.com/sacloud/iam-api-go/apis/sso"
	. "github.com/sacloud/iam-api-go/apis/sso"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, SSOAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewSSOOp(client)
	return assert, api
}

func TestList(t *testing.T) {
	var expected v1.SSOProfilesGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.SSOProfile, 1))
	expected.Items[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.List(t.Context(), nil, nil)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestList_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail("forbidden")
	assert, api := setup(t, &res, res.Status)

	actual, err := api.List(t.Context(), nil, nil)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), "forbidden")
}

func TestCreate(t *testing.T) {
	var expected v1.SSOProfile
	expected.SetFake()
	var req v1.SSOProfilesPostReq
	req.SetFake()
	assert, api := setup(t, &expected, http.StatusCreated)

	actual, err := api.Create(t.Context(), req)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestCreate_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail("bad request")
	var req v1.SSOProfilesPostReq
	req.SetFake()
	assert, api := setup(t, &res, res.Status)

	actual, err := api.Create(t.Context(), req)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), "bad request")
}

func TestGet(t *testing.T) {
	var expected v1.SSOProfile
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), 123)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestGet_Fail(t *testing.T) {
	var res v1.Http404NotFound
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail("not found")
	assert, api := setup(t, &res, res.Status)

	actual, err := api.Read(t.Context(), 123)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), "not found")
}

func TestUpdate(t *testing.T) {
	var expected v1.SSOProfile
	expected.SetFake()
	var req v1.SSOProfilesSSOProfileIDPutReq
	req.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Update(t.Context(), 123, req)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestUpdate_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail("forbidden")
	var req v1.SSOProfilesSSOProfileIDPutReq
	req.SetFake()
	assert, api := setup(t, &res, res.Status)

	actual, err := api.Update(t.Context(), 123, req)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), "forbidden")
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, &v1.SSOProfilesSSOProfileIDDeleteNoContent{}, http.StatusNoContent)

	err := api.Delete(t.Context(), 123)
	assert.NoError(err)
}

func TestDelete_Fail(t *testing.T) {
	var res v1.Http404NotFound
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail("not found")
	assert, api := setup(t, &res, res.Status)

	err := api.Delete(t.Context(), 123)
	assert.Error(err)
	assert.Contains(err.Error(), "not found")
}

func TestLink(t *testing.T) {
	var expected v1.SSOProfile
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Link(t.Context(), 123)
	assert.NoError(err)
	assert.Equal(&expected, actual)
}

func TestLink_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail("forbidden")
	assert, api := setup(t, &res, res.Status)

	_, err := api.Link(t.Context(), 123)
	assert.Error(err)
	assert.Contains(err.Error(), "forbidden")
}

func TestUnlinkProfile(t *testing.T) {
	var expected v1.SSOProfile
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Unlink(t.Context(), 123)
	assert.NoError(err)
	assert.Equal(&expected, actual)
}

func TestUnlinkProfile_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail("forbidden")
	assert, api := setup(t, &res, res.Status)

	_, err := api.Unlink(t.Context(), 123)
	assert.Error(err)
	assert.Contains(err.Error(), "forbidden")
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	api := NewSSOOp(client)

	// Create
	var createParam v1.SSOProfilesPostReq
	createParam.IdpLoginURL = "https://example.com/sso/login"
	createParam.IdpLogoutURL = "https://example.com/sso/logout"
	createParam.IdpEntityID = "https://example.com/sso/issuer"
	createParam.IdpCertificate = string(cert(assert))
	createParam.Name = testutil.RandomName("sso-", 32, testutil.CharSetAlphaNum)
	created, err := api.Create(t.Context(), createParam)
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

	// Update
	var updateParam sso.UpdateParams
	updateParam.IdpLoginURL = read.GetIdpLoginURL()
	updateParam.IdpLogoutURL = read.GetIdpLogoutURL()
	updateParam.IdpEntityID = read.GetIdpEntityID()
	updateParam.IdpCertificate = read.GetIdpCertificate()
	updateParam.Name = read.GetName()
	updateParam.Description = testutil.Random(64, testutil.CharSetAlphaNum)
	updated, err := api.Update(t.Context(), created.GetID(), updateParam)
	assert.NoError(err)
	assert.NotNil(updated)

	// Link
	linked, err := api.Link(t.Context(), created.GetID())
	assert.NoError(err)
	assert.NotNil(linked)

	// Unlink
	unlinked, err := api.Unlink(t.Context(), created.GetID())
	assert.NoError(err)
	assert.NotNil(unlinked)
}

func cert(assert *require.Assertions) []byte {
	return pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: []byte(nil)})
}
