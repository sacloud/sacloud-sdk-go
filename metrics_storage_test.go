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
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/google/uuid"
	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func TestMetricsStorageOp_List(t *testing.T) {
	expected := v1.PaginatedMetricsStorageList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.MetricsStorage{TemplateMetricsStorage},
	}
	client := newTestClient(expected)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	tanks, err := api.List(ctx, MetricsStorageListParams{Count: nil, From: nil})
	require.NoError(t, err)
	require.NotNil(t, tanks)
	require.Equal(t, 1, len(tanks))

	tank := tanks[0]
	require.Equal(t, TemplateMetricsStorage.GetName(), tank.GetName())
	require.Equal(t, TemplateMetricsStorage.GetDescription(), tank.GetDescription())
	require.Equal(t, TemplateMetricsStorage.GetIsSystem(), tank.GetIsSystem())
	require.Equal(t, TemplateMetricsStorage.GetAccountID(), tank.GetAccountID())
	require.Equal(t, TemplateMetricsStorage.GetResourceID(), tank.GetResourceID())
	require.Equal(t, TemplateMetricsStorage.GetEndpoints(), tank.GetEndpoints())
	require.Equal(t, TemplateMetricsStorage.GetUsage(), tank.GetUsage())
}

func TestMetricsStorageOp_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	tanks, err := api.List(ctx, MetricsStorageListParams{})
	require.Nil(t, tanks)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_Read(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsStorage)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "12345")
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedMetricsStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedMetricsStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedMetricsStorage.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateWrappedMetricsStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedMetricsStorage.GetResourceID(), actual.GetResourceID())
	require.Equal(t, (&TemplateWrappedMetricsStorage.Endpoints).GetAddress(), (&actual.Endpoints).GetAddress())
	require.Equal(t, (&TemplateWrappedMetricsStorage.Usage).GetMetricsRoutings(), (&actual.Usage).GetMetricsRoutings())
	require.Equal(t, (&TemplateWrappedMetricsStorage.Usage).GetAlertRules(), (&actual.Usage).GetAlertRules())
	require.Equal(t, (&TemplateWrappedMetricsStorage.Usage).GetLogMeasureRules(), (&actual.Usage).GetLogMeasureRules())
	require.Equal(t, TemplateWrappedMetricsStorage.GetTags(), actual.GetTags())
}

func TestMetricsStorageOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No MetricsStorage matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "99999")
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "not found")
}

func TestMetricsStorageOp_Create(t *testing.T) {
	client := newTestClient(TemplateMetricsStorage, http.StatusCreated)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	createReq := MetricsStorageCreateParams{
		Name:        "created-tank",
		Description: ref("Created metrics tank"),
		IsSystem:    false,
	}
	actual, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateMetricsStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateMetricsStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateMetricsStorage.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateMetricsStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateMetricsStorage.GetResourceID(), actual.GetResourceID())
	require.Equal(t, (&TemplateMetricsStorage.Endpoints).GetAddress(), (&actual.Endpoints).GetAddress())
	require.Equal(t, (&TemplateMetricsStorage.Usage).GetMetricsRoutings(), (&actual.Usage).GetMetricsRoutings())
	require.Equal(t, (&TemplateMetricsStorage.Usage).GetAlertRules(), (&actual.Usage).GetAlertRules())
	require.Equal(t, (&TemplateMetricsStorage.Usage).GetLogMeasureRules(), (&actual.Usage).GetLogMeasureRules())
	require.Equal(t, TemplateMetricsStorage.GetTags(), actual.GetTags())
}

func TestMetricsStorageOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	createReq := MetricsStorageCreateParams{
		Name:        "",
		Description: nil,
		IsSystem:    false,
	}
	actual, err := api.Create(ctx, createReq)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestMetricsStorageOp_Update(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsStorage)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	updatedName := "updated-tank"
	actual, err := api.Update(ctx, "54321", MetricsStorageUpdateParams{&updatedName, nil})
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedMetricsStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedMetricsStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedMetricsStorage.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateWrappedMetricsStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedMetricsStorage.GetResourceID(), actual.GetResourceID())
	require.Equal(t, (&TemplateWrappedMetricsStorage.Endpoints).GetAddress(), (&actual.Endpoints).GetAddress())
	require.Equal(t, (&TemplateWrappedMetricsStorage.Usage).GetMetricsRoutings(), (&actual.Usage).GetMetricsRoutings())
	require.Equal(t, (&TemplateWrappedMetricsStorage.Usage).GetAlertRules(), (&actual.Usage).GetAlertRules())
	require.Equal(t, (&TemplateWrappedMetricsStorage.Usage).GetLogMeasureRules(), (&actual.Usage).GetLogMeasureRules())
}

func TestMetricsStorageOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Update(ctx, "54321", MetricsStorageUpdateParams{})
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestMetricsStorageOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "54321")
	require.NoError(t, err)
}

func TestMetricsStorageOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "0")
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestMetricsStorageOp_StatsDaily(t *testing.T) {
	expected := v1.MetricsStorageDailyUsageBody{
		Usages: []v1.MetricsStorageDailyUsage{TemplateMetricsStorageDailyUsage},
	}
	client := newTestClient(expected)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := api.StatsDaily(ctx, "12345", &startDate, &endDate)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(result))
}

func TestMetricsStorageOp_StatsDaily_400(t *testing.T) {
	expected := newErrorResponse(400, "invalid parameter")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := api.StatsDaily(ctx, "invalid", &startDate, &endDate)
	require.Nil(t, result)
	require.Error(t, err)
}

func TestMetricsStorageOp_StatsMonthly(t *testing.T) {
	expected := v1.MetricsStorageMonthlyUsageBody{
		Usages: []v1.MetricsStorageMonthlyUsage{TemplateMetricsStorageMonthlyUsage},
	}
	client := newTestClient(expected)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	result, err := api.StatsMonthly(ctx, "12345", 2025)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(result))
}

func TestMetricsStorageOp_StatsMonthly_400(t *testing.T) {
	expected := newErrorResponse(400, "invalid parameter, year must be between 1970 and 2100")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	result, err := api.StatsMonthly(ctx, "99999", 2200)
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid parameter")
}

// --- Access Key API tests ---

func TestMetricsStorageOp_ListKeys(t *testing.T) {
	expected := v1.PaginatedMetricsStorageAccessKeyList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.MetricsStorageAccessKey{TemplateMetricsStorageAccessKey},
	}
	client := newTestClient(expected)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	keys, err := api.ListKeys(ctx, "12345", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, keys)
	require.Equal(t, 1, len(keys))
	require.Equal(t, TemplateMetricsStorageAccessKey.GetID(), keys[0].GetID())
	require.Contains(t, keys[0].GetDescription().Value, TemplateMetricsStorageAccessKey.GetDescription().Value)
}

func TestMetricsStorageOp_ListKeys_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	keys, err := api.ListKeys(ctx, "12345", nil, nil)
	require.Nil(t, keys)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_CreateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey, http.StatusCreated)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.CreateKey(ctx, "12345", nil)
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedAccessKey.GetID(), key.GetID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedAccessKey.GetDescription().Value)
}

func TestMetricsStorageOp_CreateKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.CreateKey(ctx, "12345", nil)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_ReadKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.ReadKey(ctx, "12345", uuid.New())
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedAccessKey.GetID(), key.GetID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedAccessKey.GetDescription().Value)
}

func TestMetricsStorageOp_ReadKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.ReadKey(ctx, "12345", uuid.New())
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_UpdateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, "12345", uuid.New(), nil)
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedAccessKey.GetID(), key.GetID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedAccessKey.GetDescription().Value)
}

func TestMetricsStorageOp_UpdateKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, "12345", uuid.New(), nil)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_DeleteKey(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", uuid.New())
	require.NoError(t, err)
}

func TestMetricsStorageOp_DeleteKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()
	tmp := WithMetricsStorage(t, client, ctx)
	sid := fmt.Sprintf("%d", tmp.GetID())

	// Read
	read, err := api.Read(ctx, sid)
	require.NoError(t, err)
	require.NotNil(t, read)
	require.Equal(t, tmp.GetID(), read.GetID())
	require.Equal(t, tmp.GetName(), read.GetName())

	// Update
	updateReq := testutil.Random(128, testutil.CharSetAlphaNum)
	updated, err := api.Update(ctx, sid, MetricsStorageUpdateParams{nil, &updateReq})
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, updateReq, updated.GetDescription().Or("failure"))

	// Create Key
	createdKey, err := api.CreateKey(ctx, sid, nil)
	require.NoError(t, err)
	require.NotNil(t, createdKey)
	require.NotZero(t, createdKey.GetUID())
	require.Empty(t, createdKey.GetDescription().Or("failure"))
	require.NotEmpty(t, createdKey.GetSecret())
	kid := createdKey.GetUID()

	// Read Key
	readKey, err := api.ReadKey(ctx, sid, kid)
	require.NoError(t, err)
	require.NotNil(t, readKey)
	require.Equal(t, createdKey.GetUID(), readKey.GetUID())

	// Update Key
	updatedDesc := testutil.Random(128, testutil.CharSetAlphaNum)
	updatedKey, err := api.UpdateKey(ctx, sid, kid, &updatedDesc)
	require.NoError(t, err)
	require.NotNil(t, updatedKey)
	require.Equal(t, createdKey.GetUID(), updatedKey.GetUID())
	require.Equal(t, updatedDesc, updatedKey.GetDescription().Or("failure"))

	// Delete Key
	err = api.DeleteKey(ctx, sid, kid)
	require.NoError(t, err)

	// List
	tanks, err := api.List(ctx, MetricsStorageListParams{Count: nil, From: nil})
	require.NoError(t, err)
	require.NotNil(t, tanks)
	require.NotEmpty(t, tanks)

	// StatsDaily
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	endDate := now
	dailyStats, err := api.StatsDaily(ctx, sid, &startDate, &endDate)
	require.NoError(t, err)
	require.NotNil(t, dailyStats)
	// Note: May be empty for newly created resources, which is acceptable

	// StatsMonthly
	currentYear := now.Year()
	monthlyStats, err := api.StatsMonthly(ctx, sid, currentYear)
	require.NoError(t, err)
	require.NotNil(t, monthlyStats)
	// Note: May be empty for newly created resources, which is acceptable
}
