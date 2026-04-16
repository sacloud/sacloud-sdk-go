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

package iampolicy_test

import (
	"net/http"
	"testing"

	. "github.com/sacloud/iam-api-go/apis/iampolicy"
	v1 "github.com/sacloud/iam-api-go/apis/v1"
	iam_test "github.com/sacloud/iam-api-go/testutil"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v any, s ...int) (*require.Assertions, IAMPolicyAPI) {
	assert := require.New(t)
	client := iam_test.NewTestClient(v, s...)
	api := NewIAMPolicyOp(client)
	return assert, api
}

func TestNewIamPolicyOp(t *testing.T) {
	assert, api := setup(t, make(map[string]any), http.StatusAccepted)
	assert.NotNil(api)
}

func TestGetOrganizationPolicy(t *testing.T) {
	var expected v1.OrganizationIamPolicyGetOK
	expected.SetFake()
	expected.SetBindings(make([]v1.IamPolicy, 1))
	expected.Bindings[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.ReadOrganizationPolicy(t.Context())
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetBindings(), actual)
}

func TestUpdateOrganizationPolicy(t *testing.T) {
	var expected v1.OrganizationIamPolicyPutOK
	expected.SetFake()
	expected.SetBindings(make([]v1.IamPolicy, 1))
	expected.Bindings[0].SetFake()
	bindings := []v1.IamPolicy{}
	assert, api := setup(t, &expected)

	actual, err := api.UpdateOrganizationPolicy(t.Context(), bindings)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetBindings(), actual)
}

func TestGetProjectPolicy(t *testing.T) {
	var expected v1.ProjectsProjectIDIamPolicyGetOK
	expected.SetFake()
	expected.SetBindings(make([]v1.IamPolicy, 1))
	expected.Bindings[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.ReadProjectPolicy(t.Context(), 123)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetBindings(), actual)
}

func TestUpdateProjectPolicy(t *testing.T) {
	var expected v1.ProjectsProjectIDIamPolicyPutOK
	expected.SetFake()
	expected.SetBindings(make([]v1.IamPolicy, 1))
	expected.Bindings[0].SetFake()
	bindings := []v1.IamPolicy{}
	assert, api := setup(t, &expected)

	actual, err := api.UpdateProjectPolicy(t.Context(), 123, bindings)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetBindings(), actual)
}

func TestGetFolderPolicy(t *testing.T) {
	var expected v1.FoldersFolderIDIamPolicyGetOK
	expected.SetFake()
	expected.SetBindings(make([]v1.IamPolicy, 1))
	expected.Bindings[0].SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.ReadFolderPolicy(t.Context(), 123)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetBindings(), actual)
}

func TestUpdateFolderPolicy(t *testing.T) {
	var expected v1.FoldersFolderIDIamPolicyPutOK
	expected.SetFake()
	expected.SetBindings(make([]v1.IamPolicy, 1))
	expected.Bindings[0].SetFake()
	bindings := []v1.IamPolicy{}
	assert, api := setup(t, &expected)

	actual, err := api.UpdateFolderPolicy(t.Context(), 123, bindings)
	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetBindings(), actual)
}

func TestGetOrganizationPolicy_Fail(t *testing.T) {
	var res v1.Http403Forbidden
	expected := "organization policy forbidden"
	res.SetFake()
	res.SetStatus(http.StatusForbidden)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	actual, err := api.ReadOrganizationPolicy(t.Context())
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestUpdateOrganizationPolicy_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := "organization policy bad request"
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	bindings := []v1.IamPolicy{}
	actual, err := api.UpdateOrganizationPolicy(t.Context(), bindings)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGetProjectPolicy_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := "project policy not found"
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	projectID := 1
	actual, err := api.ReadProjectPolicy(t.Context(), projectID)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestUpdateProjectPolicy_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := "project policy bad request"
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	projectID := 1
	bindings := []v1.IamPolicy{}
	actual, err := api.UpdateProjectPolicy(t.Context(), projectID, bindings)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestGetFolderPolicy_Fail(t *testing.T) {
	var res v1.Http404NotFound
	expected := "folder policy not found"
	res.SetFake()
	res.SetStatus(http.StatusNotFound)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	folderID := 1
	actual, err := api.ReadFolderPolicy(t.Context(), folderID)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestUpdateFolderPolicy_Fail(t *testing.T) {
	var res v1.Http400BadRequest
	expected := "folder policy bad request"
	res.SetFake()
	res.SetStatus(http.StatusBadRequest)
	res.SetDetail(expected)
	assert, api := setup(t, &res, res.Status)

	folderID := 1
	bindings := []v1.IamPolicy{}
	actual, err := api.UpdateFolderPolicy(t.Context(), folderID, bindings)
	assert.Error(err)
	assert.Nil(actual)
	assert.Contains(err.Error(), expected)
}

func TestIntegrated(t *testing.T) {
	assert, client := iam_test.IntegratedClient(t)
	api := NewIAMPolicyOp(client)

	// Organization Policy
	orgPolicy, err := api.ReadOrganizationPolicy(t.Context())
	assert.NoError(err)
	assert.NotNil(orgPolicy)
	assert.NotEmpty(orgPolicy) // being able to call this API means at least there is some policy whatsoever

	updatedOrgPolicy, err := api.UpdateOrganizationPolicy(t.Context(), orgPolicy)
	assert.NoError(err)
	assert.NotNil(updatedOrgPolicy)

	// Project Policy
	project, dproject := iam_test.NewProject(t, client)
	defer dproject()

	projectID := project.GetID()
	principal, dprincipal := iam_test.NewPrincipal(t, client, projectID)
	defer dprincipal()

	projectPolicy, err := api.ReadProjectPolicy(t.Context(), projectID)
	assert.NoError(err)
	assert.NotNil(projectPolicy)
	assert.NotEmpty(projectPolicy) // the newborn project is "owned" by the caller

	var newRole v1.IamPolicyRole
	var newPolicy v1.IamPolicy
	newRole.SetFake()
	newRole.SetType(v1.NewOptIamPolicyRoleType(v1.IamPolicyRoleTypePreset))
	newRole.SetID(v1.NewOptString("resource-viewer"))
	newPolicy.SetFake()
	newPolicy.SetRole(v1.NewOptIamPolicyRole(newRole))
	newPolicy.SetPrincipals(make([]v1.Principal, 1))
	newPolicy.Principals[0].SetFake()
	newPolicy.Principals[0].SetID(v1.NewOptInt(principal.GetID()))
	newPolicy.Principals[0].SetType(v1.NewOptString("service-principal"))

	updatedProjectPolicy, err := api.UpdateProjectPolicy(t.Context(), projectID, []v1.IamPolicy{newPolicy})
	assert.NoError(err)
	assert.NotNil(updatedProjectPolicy)

	// Folder Policy
	folder, dfolder := iam_test.NewFolder(t, client)
	defer dfolder()

	folderID := folder.GetID()
	folderPolicy, err := api.ReadFolderPolicy(t.Context(), folderID)
	assert.NoError(err)
	assert.NotNil(folderPolicy)
	assert.Empty(folderPolicy)

	updatedFolderPolicy, err := api.UpdateFolderPolicy(t.Context(), folderID, []v1.IamPolicy{newPolicy})
	assert.NoError(err)
	assert.NotNil(updatedFolderPolicy)
}
