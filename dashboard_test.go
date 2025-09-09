// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite_test

import (
	"context"
	"net/http"
	"testing"

	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func newTestDashboardClient(resp interface{}, status ...int) *v1.Client {
	return newTestClient(resp, status...)
}

func TestDashboardOp_List(t *testing.T) {
	expected := v1.PaginatedDashboardProjectList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.DashboardProject{TemplateDashboardProject},
	}
	client := newTestDashboardClient(expected)
	api := NewDashboardOp(client)
	ctx := context.Background()

	projects, err := api.List(ctx, 0, 0)
	require.NoError(t, err)
	require.NotNil(t, projects)
	require.Equal(t, 1, len(projects))
}

func TestDashboardOp_Read(t *testing.T) {
	client := newTestDashboardClient(TemplateWrappedDashboardProject)
	api := NewDashboardOp(client)
	ctx := context.Background()

	res, err := api.Read(ctx, "12345")
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, TemplateWrappedDashboardProject.GetID(), res.GetID())
	require.Equal(t, TemplateWrappedDashboardProject.GetName(), res.GetName())
	require.Equal(t, TemplateWrappedDashboardProject.GetDescription(), res.GetDescription())
	require.Equal(t, TemplateWrappedDashboardProject.GetIsSystem(), res.GetIsSystem())
	require.Equal(t, TemplateWrappedDashboardProject.GetTags(), res.GetTags())
	require.Equal(t, TemplateWrappedDashboardProject.GetAccountID(), res.GetAccountID())
	require.Equal(t, TemplateWrappedDashboardProject.GetResourceID(), res.GetResourceID())
}

func TestDashboardOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No DashboardProject matches the given query.")
	client := newTestDashboardClient(expected, http.StatusNotFound)
	api := NewDashboardOp(client)
	ctx := context.Background()

	project, err := api.Read(ctx, "99999")
	require.Nil(t, project)
	require.Error(t, err)
	require.ErrorContains(t, err, "Not Found")
}

func TestDashboardOp_Create(t *testing.T) {
	client := newTestDashboardClient(TemplateDashboardProject, http.StatusCreated)
	api := NewDashboardOp(client)
	ctx := context.Background()

	createReq := v1.DashboardProjectCreate{}
	res, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestDashboardOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestDashboardClient(expected, http.StatusBadRequest)
	api := NewDashboardOp(client)
	ctx := context.Background()

	createReq := v1.DashboardProjectCreate{}
	project, err := api.Create(ctx, createReq)
	require.Nil(t, project)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

func TestDashboardOp_Update(t *testing.T) {
	client := newTestDashboardClient(TemplateWrappedDashboardProject)
	api := NewDashboardOp(client)
	ctx := context.Background()

	updateReq := TemplateDashboardProject
	res, err := api.Update(ctx, "12345", &updateReq)
	require.NoError(t, err)
	require.NotNil(t, res)
}

func TestDashboardOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestDashboardClient(expected, http.StatusBadRequest)
	api := NewDashboardOp(client)
	ctx := context.Background()

	updateReq := v1.DashboardProject{}
	project, err := api.Update(ctx, "0", &updateReq)
	require.Nil(t, project)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

func TestDashboardOp_Delete(t *testing.T) {
	client := newTestDashboardClient(nil, http.StatusNoContent)
	api := NewDashboardOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "12345")
	require.NoError(t, err)
}

func TestDashboardOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestDashboardClient(expected, http.StatusBadRequest)
	api := NewDashboardOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "0")
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}
