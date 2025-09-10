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

	"github.com/google/uuid"
	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestAlertProjectOp_List(t *testing.T) {
	expected := v1.PaginatedAlertProjectList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.AlertProject{TemplateAlertProject},
	}
	client := newTestClient(expected)
	api := NewAlertProjectOp(client)
	ctx := context.Background()
	projects, err := api.List(ctx, 32, 0)
	require.NoError(t, err)
	require.NotNil(t, projects)
	require.Equal(t, 1, len(projects))

	project := projects[0]
	require.Equal(t, TemplateAlertProject.GetName(), project.GetName())
	require.Equal(t, TemplateAlertProject.GetDescription(), project.GetDescription())
	require.Equal(t, TemplateAlertProject.GetIsSystem(), project.GetIsSystem())
	require.Equal(t, TemplateAlertProject.GetAccountID(), project.GetAccountID())
	require.Equal(t, TemplateAlertProject.GetResourceID(), project.GetResourceID())
	require.Equal(t, TemplateAlertProject.GetTags(), project.GetTags())
}

func TestAlertProjectOp_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewAlertProjectOp(client)
	ctx := context.Background()
	_, err := api.List(ctx, 0, 0)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permission")
}

func TestAlertProjectOp_Read(t *testing.T) {
	client := newTestClient(TemplateWrappedAlertProject)
	api := NewAlertProjectOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "12345")
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedAlertProject.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedAlertProject.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedAlertProject.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateWrappedAlertProject.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedAlertProject.GetResourceID(), actual.GetResourceID())
	require.Equal(t, TemplateWrappedAlertProject.GetTags(), actual.GetTags())
}

func TestAlertProjectOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No AlertProject matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewAlertProjectOp(client)
	ctx := context.Background()

	_, err := api.Read(ctx, "12345")
	require.Error(t, err)
	require.ErrorContains(t, err, "Not Found")
}

func TestAlertProjectOp_Create(t *testing.T) {
	client := newTestClient(TemplateAlertProject, http.StatusCreated)
	api := NewAlertProjectOp(client)
	ctx := context.Background()

	createReq := v1.AlertProjectCreate{
		Name:        "created-alert-project",
		Description: "Created alert project",
		IsSystem:    v1.NewOptBool(false),
	}
	actual, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateAlertProject.GetName(), actual.GetName())
	require.Equal(t, TemplateAlertProject.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateAlertProject.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateAlertProject.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateAlertProject.GetResourceID(), actual.GetResourceID())
	require.Equal(t, TemplateAlertProject.GetTags(), actual.GetTags())
}

func TestAlertProjectOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewAlertProjectOp(client)
	ctx := context.Background()

	createReq := v1.AlertProjectCreate{
		Name:        "",
		Description: "",
		IsSystem:    v1.NewOptBool(false),
	}
	actual, err := api.Create(ctx, createReq)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestAlertProjectOp_Update(t *testing.T) {
	client := newTestClient(TemplateWrappedAlertProject)
	api := NewAlertProjectOp(client)
	ctx := context.Background()

	actual, err := api.Update(ctx, "54321", &TemplateAlertProject)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedAlertProject.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedAlertProject.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedAlertProject.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateWrappedAlertProject.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedAlertProject.GetResourceID(), actual.GetResourceID())
	require.Equal(t, TemplateWrappedAlertProject.GetTags(), actual.GetTags())
}

func TestAlertProjectOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewAlertProjectOp(client)
	ctx := context.Background()

	actual, err := api.Update(ctx, "0", &TemplateAlertProject)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestAlertProjectOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewAlertProjectOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "54321")
	require.NoError(t, err)
}

func TestAlertProjectOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewAlertProjectOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "0")
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestAlertProjectOp_ListHistories(t *testing.T) {
	expected := v1.PaginatedHistoryList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.History{TemplateHistory},
	}
	client := newTestClient(expected)
	api := NewAlertProjectOp(client)
	ctx := context.Background()
	params := v1.AlertsProjectsHistoriesListParams{
		Count: v1.NewOptInt(32),
		From:  v1.NewOptInt(0),
	}
	histories, err := api.ListHistories(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, histories)
	require.Equal(t, 1, len(histories))
	require.Equal(t, TemplateHistory.GetUID(), histories[0].GetUID())
	require.Equal(t, TemplateHistory.GetProjectID(), histories[0].GetProjectID())
}

func TestAlertProjectOp_ListHistories_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewAlertProjectOp(client)
	ctx := context.Background()
	params := v1.AlertsProjectsHistoriesListParams{}
	_, err := api.ListHistories(ctx, params)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permission")
}

func TestAlertProjectOp_ReadHistory(t *testing.T) {
	client := newTestClient(TemplateHistory)
	api := NewAlertProjectOp(client)
	ctx := context.Background()
	history, err := api.ReadHistory(ctx, "12345", uuid.New())
	require.NoError(t, err)
	require.NotNil(t, history)
	require.Equal(t, TemplateHistory.GetUID(), history.GetUID())
	require.Equal(t, TemplateHistory.GetProjectID(), history.GetProjectID())
}

func TestAlertProjectOp_ReadHistory_404(t *testing.T) {
	expected := newErrorResponse(404, "No History matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewAlertProjectOp(client)
	ctx := context.Background()
	_, err := api.ReadHistory(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "not found")
}
