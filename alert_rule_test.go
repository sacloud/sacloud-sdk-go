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
	"fmt"
	"net/http"
	"testing"

	"github.com/google/uuid"
	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
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
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	rules, err := api.List(ctx, "12345", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, rules)
	require.Equal(t, 1, len(rules))
	rule := rules[0]
	require.Equal(t, TemplateAlertRule.GetUID(), rule.GetUID())
	require.Equal(t, TemplateAlertRule.GetProjectID(), rule.GetProjectID())
	require.Equal(t, TemplateAlertRule.GetName(), rule.GetName())
	require.Equal(t, TemplateAlertRule.GetQuery(), rule.GetQuery())
}

func TestAlertRuleOp_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	_, err := api.List(ctx, "12345", nil, nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "request not authorized")
}

func TestAlertRuleOp_Read(t *testing.T) {
	client := newTestClient(TemplateAlertRule)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	actual, err := api.Read(ctx, "12345", uuid.New())
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateAlertRule.GetUID(), actual.GetUID())
	require.Equal(t, TemplateAlertRule.GetProjectID(), actual.GetProjectID())
	require.Equal(t, TemplateAlertRule.GetName(), actual.GetName())
	require.Equal(t, TemplateAlertRule.GetQuery(), actual.GetQuery())
}

func TestAlertRuleOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No AlertRule matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	_, err := api.Read(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "No AlertRule matches the given query.")
}

func TestAlertRuleOp_Create(t *testing.T) {
	client := newTestClient(TemplateAlertRule, http.StatusCreated)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	actual, err := api.Create(ctx, "12345", AlertRuleCreateParams{
		MetricsStorageID: "56789",
		Query:            "q",
	})
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateAlertRule.GetName(), actual.GetName())
	require.Equal(t, TemplateAlertRule.GetQuery(), actual.GetQuery())
}

func TestAlertRuleOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	actual, err := api.Create(ctx, "12345", AlertRuleCreateParams{MetricsStorageID: "56789"})
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid request body.")
}

func TestAlertRuleOp_Update(t *testing.T) {
	client := newTestClient(TemplateAlertRule)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	name := "rule"
	actual, err := api.Update(ctx, "12345", uuid.New(), AlertRuleUpdateParams{
		Name: &name,
	})
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateAlertRule.GetName(), actual.GetName())
	require.Equal(t, TemplateAlertRule.GetQuery(), actual.GetQuery())
}

func TestAlertRuleOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	actual, err := api.Update(ctx, "12345", uuid.New(), AlertRuleUpdateParams{})
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid update parameters.")
}

func TestAlertRuleOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	err := api.Delete(ctx, "12345", uuid.New())
	require.NoError(t, err)
}

func TestAlertRuleOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	err := api.Delete(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid delete request.")
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
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	params := AlertRuleListHistoriesParams{
		Count: nil,
		From:  nil,
	}
	histories, err := api.ListHistories(ctx, "12345", uuid.New(), params)
	require.NoError(t, err)
	require.Equal(t, 1, len(histories))
	require.Equal(t, TemplateHistory.GetUID(), histories[0].GetUID())
	require.Equal(t, TemplateHistory.GetRuleUID(), histories[0].GetRuleUID())
}

func TestAlertRuleOp_ListHistories_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	params := AlertRuleListHistoriesParams{
		Count: nil,
		From:  nil,
	}
	_, err := api.ListHistories(ctx, "12345", uuid.New(), params)
	require.Error(t, err)
	require.ErrorContains(t, err, "request not authorized")
}

func TestAlertRuleOp_ReadHistory(t *testing.T) {
	client := newTestClient(TemplateHistory)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	actual, err := api.ReadHistory(ctx, "123", uuid.New(), uuid.New())
	require.NoError(t, err)
	require.Equal(t, TemplateHistory.GetUID(), actual.GetUID())
	require.Equal(t, TemplateHistory.GetRuleUID(), actual.GetRuleUID())
}

func TestAlertRuleOp_ReadHistory_404(t *testing.T) {
	expected := newErrorResponse(404, "No History matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewAlertRuleOp(client)
	ctx := t.Context()
	_, err := api.ReadHistory(ctx, "123", uuid.New(), uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "No History matches the given query.")
}

func TestAlertRuleIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)

	api := NewAlertRuleOp(client)
	ctx := t.Context()

	project := WithAlertProject(t, client, ctx)
	projectId := fmt.Sprintf("%d", project.GetID())

	storage := WithMetricsStorage(t, client, ctx)
	storageId := fmt.Sprintf("%d", storage.GetID())

	// Create
	name := testutil.RandomName("test-alert-rule-", 16, testutil.CharSetAlphaNum)
	created, err := api.Create(ctx, projectId, AlertRuleCreateParams{
		MetricsStorageID:         storageId,
		Name:                     &name,
		Query:                    "count_values",
		EnabledWarning:           ref(true),
		ThresholdWarning:         ref("==0"),
		ThresholdDurationWarning: ref[int64](32768),
	})
	require.NoError(t, err)
	require.NotNil(t, created)
	uid := created.GetUID()

	// Delete
	t.Cleanup(func() {
		err = api.Delete(ctx, projectId, uid)
		require.NoError(t, err)
	})

	// Read
	read, err := api.Read(ctx, projectId, uid)
	require.NoError(t, err)
	require.NotNil(t, read)
	require.Equal(t, uid, read.GetUID())
	require.Equal(t, created.GetName(), read.GetName())
	require.Equal(t, created.GetQuery(), read.GetQuery())

	// Update
	rename := testutil.Random(16, testutil.CharSetAlphaNum)
	updated, err := api.Update(ctx, projectId, uid, AlertRuleUpdateParams{
		// :TODO: this `EnabledWarning` is tentatively mandatory
		// but subject to change in the future.
		EnabledWarning: ref(true),
		Name:           &rename,
	})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, rename, updated.GetName().Or("failure"))

	// List rules
	rules, err := api.List(ctx, projectId, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, rules)

	// List histories
	histories, err := api.ListHistories(ctx, projectId, uid, AlertRuleListHistoriesParams{})
	require.NoError(t, err)
	require.NotNil(t, histories)
}
