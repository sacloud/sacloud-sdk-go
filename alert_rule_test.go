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

	monitoringsuite "github.com/sacloud/monitoring-suite-api-go"

	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestAlertRuleOp_List(t *testing.T) {
	expected := v1.PaginatedAlertRuleList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.AlertRule{TemplateAlertRule},
	}
	client := newTestClient(expected)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	params := v1.AlertsProjectsRulesListParams{
		Count: v1.NewOptInt(32),
		From:  v1.NewOptInt(0),
	}
	rules, err := api.List(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, rules)
	require.Equal(t, 1, len(rules))
	rule := rules[0]
	require.Equal(t, TemplateAlertRule.GetID(), rule.GetID())
	require.Equal(t, TemplateAlertRule.GetProjectID(), rule.GetProjectID())
	require.Equal(t, TemplateAlertRule.GetName(), rule.GetName())
	require.Equal(t, TemplateAlertRule.GetQuery(), rule.GetQuery())
}

func TestAlertRuleOp_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	params := v1.AlertsProjectsRulesListParams{}
	_, err := api.List(ctx, params)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permission")
}

func TestAlertRuleOp_Read(t *testing.T) {
	client := newTestClient(TemplateAlertRule)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	actual, err := api.Read(ctx, "12345", "56789")
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateAlertRule.GetID(), actual.GetID())
	require.Equal(t, TemplateAlertRule.GetProjectID(), actual.GetProjectID())
	require.Equal(t, TemplateAlertRule.GetName(), actual.GetName())
	require.Equal(t, TemplateAlertRule.GetQuery(), actual.GetQuery())
}

func TestAlertRuleOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No AlertRule matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	_, err := api.Read(ctx, "12345", "99999")
	require.Error(t, err)
	require.ErrorContains(t, err, "Not Found")
}

func TestAlertRuleOp_Create(t *testing.T) {
	client := newTestClient(TemplateAlertRule, http.StatusCreated)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	rule := &TemplateAlertRule
	actual, err := api.Create(ctx, "12345", rule)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, rule.GetName(), actual.GetName())
	require.Equal(t, rule.GetQuery(), actual.GetQuery())
}

func TestAlertRuleOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	actual, err := api.Create(ctx, "12345", &v1.AlertRule{})
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestAlertRuleOp_Update(t *testing.T) {
	client := newTestClient(TemplateAlertRule)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	rule := &TemplateAlertRule
	actual, err := api.Update(ctx, "12345", "56789", rule)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, rule.GetName(), actual.GetName())
	require.Equal(t, rule.GetQuery(), actual.GetQuery())
}

func TestAlertRuleOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	actual, err := api.Update(ctx, "12345", "56789", &v1.AlertRule{})
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

// ...existing code...
func TestAlertRuleOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	err := api.Delete(ctx, "12345", "56789")
	require.NoError(t, err)
}

func TestAlertRuleOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	err := api.Delete(ctx, "12345", "56789")
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestAlertRuleOp_ListHistories(t *testing.T) {
	expected := v1.PaginatedHistoryList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Total:   1,
		Results: []v1.History{TemplateHistory},
	}
	client := newTestClient(expected)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	params := v1.AlertsProjectsRulesHistoriesListParams{
		Count: v1.NewOptInt(10),
		From:  v1.NewOptInt(0),
	}
	histories, err := api.ListHistories(ctx, params)
	require.NoError(t, err)
	require.Equal(t, 1, len(histories))
	require.Equal(t, TemplateHistory.GetID(), histories[0].GetID())
	require.Equal(t, TemplateHistory.GetRuleID(), histories[0].GetRuleID())
}

func TestAlertRuleOp_ListHistories_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	params := v1.AlertsProjectsRulesHistoriesListParams{}
	_, err := api.ListHistories(ctx, params)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permission")
}

func TestAlertRuleOp_ReadHistory(t *testing.T) {
	client := newTestClient(TemplateHistory)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	actual, err := api.ReadHistory(ctx, "123", "456", "789")
	require.NoError(t, err)
	require.Equal(t, TemplateHistory.GetID(), actual.GetID())
	require.Equal(t, TemplateHistory.GetRuleID(), actual.GetRuleID())
}

func TestAlertRuleOp_ReadHistory_404(t *testing.T) {
	expected := newErrorResponse(404, "No History matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := monitoringsuite.NewAlertRuleOp(client)
	ctx := context.Background()
	_, err := api.ReadHistory(ctx, "123", "456", "789")
	require.Error(t, err)
	require.ErrorContains(t, err, "Not Found")
}
