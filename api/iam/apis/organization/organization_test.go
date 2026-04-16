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

package organization_test

import (
	"net/http"
	"testing"

	"github.com/sacloud/iam-api-go"
	. "github.com/sacloud/iam-api-go/apis/organization"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, OrganizationAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewOrganizationOp(client)
	return assert, api
}

func TestNewOrganizationOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestGet(t *testing.T) {
	var expected v1.Organization
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context())
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

	actual, err := api.Read(t.Context())
	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestPut(t *testing.T) {
	var expected v1.Organization
	name := testutil.RandomName("org", 32, testutil.CharSetAlphaNum)
	expected.SetFake()
	expected.SetName(name)
	assert, api := setup(t, &expected)

	actual, err := api.Update(t.Context(), name)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestPut_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.Update(t.Context(), "updated")
	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestGetServicePolicy(t *testing.T) {
	var expected v1.OrganizationServicePolicyGetOK
	expected.SetFake()
	expected.SetRules(make([]v1.RuleResponse, 1))
	expected.Rules[0].SetFake()
	assert, api := setup(t, &expected)

	params := GetServicePolicyParams{IsActive: new(bool)}
	actual, err := api.ReadServicePolicy(t.Context(), params)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.Rules, actual)
}

func TestGetServicePolicy_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	params := GetServicePolicyParams{IsActive: new(bool)}
	actual, err := api.ReadServicePolicy(t.Context(), params)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestPutServicePolicy(t *testing.T) {
	var expected v1.OrganizationServicePolicyPutOK
	expected.SetFake()
	expected.SetRules(make([]v1.RuleResponse, 1))
	expected.Rules[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.UpdateServicePolicy(t.Context(), []v1.Rule{})
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.Rules, actual)
}

func TestPutServicePolicy_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.UpdateServicePolicy(t.Context(), []v1.Rule{})
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	api := NewOrganizationOp(client)

	actual, err := api.Read(t.Context())
	assert.NoError(err)
	assert.NotNil(actual)

	name := actual.Name
	defer func() {
		_, err = api.Update(t.Context(), name)
		assert.NoError(err)
	}()

	actual, err = api.Update(t.Context(), testutil.RandomName("test::", 128, testutil.CharSetAlphaNum))
	assert.NoError(err)
	assert.NotNil(actual)

	policyop := iam.NewServicePolicyOp(client)
	before, err := policyop.IsEnabled(t.Context())
	assert.NoError(err)
	if before {
		defer func() {
			err = policyop.Enable(t.Context())
			assert.NoError(err)
		}()
	} else {
		defer func() {
			err = policyop.Disable(t.Context())
			assert.NoError(err)
		}()
	}

	err = policyop.Enable(t.Context())
	assert.NoError(err)

	actualPolicies, err := api.ReadServicePolicy(t.Context(), GetServicePolicyParams{})
	assert.NoError(err)
	assert.NotNil(actualPolicies)

	rules := make([]v1.Rule, 0, len(actualPolicies))
	for _, r := range actualPolicies {
		rules = append(rules, into(&r))
	}
	actualPolicies, err = api.UpdateServicePolicy(t.Context(), rules)
	assert.NoError(err)
	assert.NotNil(actualPolicies)
}

func into(from *v1.RuleResponse) (ret v1.Rule) {
	if val, ok := from.GetCode().Get(); ok {
		ret.SetCode(v1.NewOptString(val))
	}
	if val, ok := from.GetSpec().Get(); ok {
		if len(val.GetContents()) > 0 {
			ret.SetSpec(v1.NewOptRuleSpec(val))
		}
	}
	if val, ok := from.GetDryRunSpec().Get(); ok {
		if len(val.GetContents()) > 0 {
			ret.SetDryRunSpec(v1.NewOptRuleSpec(val))
		}
	}

	if (!ret.GetSpec().IsSet()) && (!ret.GetDryRunSpec().IsSet()) {
		// fall back to empty spec
		var spec v1.RuleSpec
		spec.SetFake()
		spec.SetContents(make([]v1.RuleContent, 1))
		spec.Contents[0].SetFake()
		spec.Contents[0].SetAllowAll(v1.NewOptBool(true))
		spec.Contents[0].SetDenyAll(v1.NewOptBool(false))
		ret.SetSpec(v1.NewOptRuleSpec(spec))
	}

	if val, ok := from.GetIsActive().Get(); ok {
		ret.SetIsActive(v1.NewOptBool(val))
	}
	if val, ok := from.GetIsDryRun().Get(); ok {
		ret.SetIsDryRun(v1.NewOptBool(val))
	}
	return ret
}
