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

package servicepolicy_test

import (
	"net/http"
	"testing"

	. "github.com/sacloud/iam-api-go/apis/servicepolicy"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, ServicePolicyAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewServicePolicyOp(client)
	return assert, api
}

func TestNewServicePolicyOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestEnable(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	err := api.Enable(t.Context())
	assert.NoError(err)
}

func TestEnable_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	err := api.Enable(t.Context())
	assert.Error(err)
	assert.False(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestDisable(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	err := api.Disable(t.Context())
	assert.NoError(err)
}

func TestDisable_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	err := api.Disable(t.Context())
	assert.Error(err)
	assert.False(saclient.IsNotFoundError(err))
	assert.Contains(err.Error(), expected)
}

func TestGetStatus(t *testing.T) {
	var expected v1.ServicePolicyStatusGetOK
	expected.SetFake()
	expected.SetEnabled(true)
	assert, api := setup(t, &expected)

	actual, err := api.IsEnabled(t.Context())
	assert.NoError(err)
	assert.True(actual)
	assert.Equal(expected.Enabled, actual)
}

func TestGetStatus_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.IsEnabled(t.Context())
	assert.Error(err)
	assert.False(actual)
	assert.Contains(err.Error(), expected)
}

func TestGetRuleTemplates(t *testing.T) {
	var expected v1.ServicePolicyRuleTemplatesGetOK
	expected.SetFake()
	expected.SetItems(make([]v1.RuleTemplate, 1))
	expected.Items[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.ListRuleTemplates(t.Context(), ListRuleTemplatesParams{})
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(&expected, actual)
}

func TestGetRuleTemplates_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := testutil.Random(128, testutil.CharSetAlphaNum)
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ListRuleTemplates(t.Context(), ListRuleTemplatesParams{})
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	op := NewServicePolicyOp(client)

	ok, err := op.IsEnabled(t.Context())
	assert.NoError(err)

	if ok {
		defer func() {
			ok, err := op.IsEnabled(t.Context())
			assert.NoError(err)
			if !ok {
				err = op.Enable(t.Context())
				assert.NoError(err)
			}
		}()
	} else {
		defer func() {
			ok, err := op.IsEnabled(t.Context())
			assert.NoError(err)
			if ok {
				err = op.Disable(t.Context())
				assert.NoError(err)
			}
		}()
	}

	err = op.Enable(t.Context())
	assert.NoError(err)

	ok, err = op.IsEnabled(t.Context())
	assert.NoError(err)
	assert.True(ok)

	err = op.Disable(t.Context())
	assert.NoError(err)

	ok, err = op.IsEnabled(t.Context())
	assert.NoError(err)
	assert.False(ok)
}
