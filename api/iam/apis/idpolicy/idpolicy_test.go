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

package idpolicy_test

import (
	"net/http"
	"testing"

	. "github.com/sacloud/iam-api-go/apis/idpolicy"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, IDPolicyAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewIDPolicyOp(client)
	return assert, api
}

func TestNewIdPolicyOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestGetOrganizationIdPolicy(t *testing.T) {
	var expected v1.OrganizationIDPolicyGetOK
	expected.SetFake()
	expected.SetBindings(make([]v1.IdPolicy, 1))
	expected.Bindings[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.ReadOrganizationIdPolicy(t.Context())
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetBindings(), actual)
}

func TestGetOrganizationIdPolicy_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	detail := testutil.Random(32, testutil.CharSetAlphaNum)
	res.SetDetail(detail)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ReadOrganizationIdPolicy(t.Context())
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), detail)
}

func TestUpdateOrganizationIdPolicy(t *testing.T) {
	var expected v1.OrganizationIDPolicyGetOK
	expected.SetFake()
	expected.SetBindings(make([]v1.IdPolicy, 1))
	expected.Bindings[0].SetFake()
	assert, api := setup(t, &expected)

	bindings := make([]v1.IdPolicy, 32)
	for i := range bindings {
		bindings[i].SetFake()
	}
	actual, err := api.UpdateOrganizationIdPolicy(t.Context(), bindings)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetBindings(), actual)
}

func TestUpdateOrganizationIdPolicy_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	detail := testutil.Random(32, testutil.CharSetAlphaNum)
	res.SetDetail(detail)
	assert, api := setup(t, &res, res.Status)

	bindings := make([]v1.IdPolicy, 32)
	for i := range bindings {
		bindings[i].SetFake()
	}
	actual, err := api.UpdateOrganizationIdPolicy(t.Context(), bindings)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), detail)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	api := NewIDPolicyOp(client)

	// Organization ID Policy
	orgPolicy, err := api.ReadOrganizationIdPolicy(t.Context())
	assert.NoError(err)
	assert.NotNil(orgPolicy)
	if len(orgPolicy) == 0 {
		return
	}

	updatedOrgPolicy, err := api.UpdateOrganizationIdPolicy(t.Context(), orgPolicy)
	assert.NoError(err)
	assert.NotNil(updatedOrgPolicy)
}
