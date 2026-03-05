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
	"github.com/stretchr/testify/require"
)

func TestTracesStorageOp_List(t *testing.T) {
	expected := v1.PaginatedTraceStorageList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.TraceStorage{TemplateTraceStorage},
	}
	client := newTestClient(expected)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	tanks, err := api.List(ctx, TracesStorageListParams{Count: nil, From: nil})
	require.NoError(t, err)
	require.NotNil(t, tanks)
	require.Equal(t, 1, len(tanks))

	tank := tanks[0]
	require.Equal(t, TemplateTraceStorage.GetName(), tank.GetName())
	require.Equal(t, TemplateTraceStorage.GetDescription(), tank.GetDescription())
	require.Equal(t, TemplateTraceStorage.GetAccountID(), tank.GetAccountID())
	require.Equal(t, TemplateTraceStorage.GetResourceID(), tank.GetResourceID())
	// Add more field checks as appropriate for TraceStorage
}

func TestTracesStorageOp_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	tanks, err := api.List(ctx, TracesStorageListParams{})
	require.Nil(t, tanks)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestTracesStorageOp_Read(t *testing.T) {
	client := newTestClient(TemplateWrappedTraceStorage)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "12345")
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedTraceStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedTraceStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedTraceStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedTraceStorage.GetResourceID(), actual.GetResourceID())
	// Add more field checks as appropriate for TraceStorage
}

func TestTracesStorageOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No TraceStorage matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "99999")
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "not found")
}

func TestTracesStorageOp_Create(t *testing.T) {
	client := newTestClient(TemplateTraceStorage, http.StatusCreated)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	createReq := TracesStorageCreateParams{
		Name:        "created-tank",
		Description: ref("Created traces tank"),
	}
	actual, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateTraceStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateTraceStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateTraceStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateTraceStorage.GetResourceID(), actual.GetResourceID())
	// Add more field checks as appropriate for TraceStorage
}

func TestTracesStorageOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	createReq := TracesStorageCreateParams{
		Name:        "",
		Description: nil,
	}
	actual, err := api.Create(ctx, createReq)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid request")
}

func TestTracesStorageOp_Update(t *testing.T) {
	client := newTestClient(TemplateWrappedTraceStorage)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	updatedName := "updated-tank"
	actual, err := api.Update(ctx, "54321", TracesStorageUpdateParams{&updatedName, nil})
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedTraceStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedTraceStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedTraceStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedTraceStorage.GetResourceID(), actual.GetResourceID())
}

func TestTracesStorageOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	actual, err := api.Update(ctx, "54321", TracesStorageUpdateParams{})
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid parameter")
}

func TestTracesStorageOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "54321")
	require.NoError(t, err)
}

func TestTracesStorageOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "0")
	require.Error(t, err)
	require.ErrorContains(t, err, "not eligible for deletion")
}

func TestTracesStorageOp_SetExpire(t *testing.T) {
	client := newTestClient(TemplateWrappedTraceStorage)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	result, err := api.SetExpire(ctx, "12345", 365)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, TemplateWrappedTraceStorage.GetName(), result.GetName())
}

func TestTracesStorageOp_SetExpire_400(t *testing.T) {
	expected := newErrorResponse(400, "invalid parameter, days must be between 1 and 730")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	result, err := api.SetExpire(ctx, "12345", 999)
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

func TestTracesStorageOp_StatsDaily(t *testing.T) {
	expected := v1.LogStorageDailyUsageBody{
		Usages: []v1.LogStorageDailyUsage{TemplateLogStorageDailyUsage},
	}
	client := newTestClient(expected)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := api.ReadDailyStats(ctx, "12345", &startDate, &endDate)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(result))
}

func TestTracesStorageOp_StatsDaily_400(t *testing.T) {
	expected := newErrorResponse(400, "invalid parameter")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := api.ReadDailyStats(ctx, "invalid", &startDate, &endDate)
	require.Nil(t, result)
	require.Error(t, err)
}

func TestTracesStorageOp_StatsMonthly(t *testing.T) {
	expected := v1.LogStorageMonthlyUsageBody{
		Usages: []v1.LogStorageMonthlyUsage{TemplateLogStorageMonthlyUsage},
	}
	client := newTestClient(expected)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	result, err := api.ReadMonthlyStats(ctx, "12345", 2025)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(result))
}

func TestTracesStorageOp_StatsMonthly_400(t *testing.T) {
	expected := newErrorResponse(400, "invalid parameter, year must be between 1970 and 2100")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	result, err := api.ReadMonthlyStats(ctx, "99999", 2200)
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid parameter")
}

// --- Access Key API tests ---

func TestTracesStorageOp_ListKeys(t *testing.T) {
	expected := v1.PaginatedTraceStorageAccessKeyList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.TraceStorageAccessKey{TemplateTraceStorageAccessKey},
	}
	client := newTestClient(expected)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	keys, err := api.ListKeys(ctx, "12345", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, keys)
	require.Equal(t, 1, len(keys))
	require.Equal(t, TemplateTraceStorageAccessKey.GetUID(), keys[0].GetUID())
	require.Contains(t, keys[0].GetDescription().Value, TemplateTraceStorageAccessKey.GetDescription().Value)
}

func TestTracesStorageOp_ListKeys_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	keys, err := api.ListKeys(ctx, "12345", nil, nil)
	require.Nil(t, keys)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestTracesStorageOp_CreateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedTraceStorageAccessKey, http.StatusCreated)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	key, err := api.CreateKey(ctx, "12345", ref("new key"))
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedTraceStorageAccessKey.GetUID(), key.GetUID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedTraceStorageAccessKey.GetDescription().Value)
}

func TestTracesStorageOp_CreateKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	key, err := api.CreateKey(ctx, "12345", nil)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestTracesStorageOp_ReadKey(t *testing.T) {
	client := newTestClient(TemplateWrappedTraceStorageAccessKey)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	key, err := api.ReadKey(ctx, "12345", uuid.New())
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedTraceStorageAccessKey.GetUID(), key.GetUID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedTraceStorageAccessKey.GetDescription().Value)
}

func TestTracesStorageOp_ReadKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	key, err := api.ReadKey(ctx, "12345", uuid.New())
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestTracesStorageOp_UpdateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedTraceStorageAccessKey)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, "12345", uuid.New(), ref("updated key"))
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedTraceStorageAccessKey.GetUID(), key.GetUID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedTraceStorageAccessKey.GetDescription().Value)
}

func TestTracesStorageOp_UpdateKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, "12345", uuid.New(), nil)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestTracesStorageOp_DeleteKey(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", uuid.New())
	require.NoError(t, err)
}

func TestTracesStorageOp_DeleteKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewTracesStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestTracesStorageIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)
	ctx := context.Background()
	api := NewTracesStorageOp(client)
	created := WithTraceStorage(t, client, ctx)
	tid := fmt.Sprintf("%d", created.GetID())

	// Read
	read, err := api.Read(ctx, tid)
	require.NoError(t, err)
	require.NotNil(t, read)
	require.Equal(t, created.GetName(), read.GetName())

	// Update
	updatedName := "integration-test-trace-storage-updated"
	updateReq := TracesStorageUpdateParams{Name: &updatedName}
	updated, err := api.Update(ctx, tid, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, updatedName, updated.GetName().Or("failure"))

	// List
	list, err := api.List(ctx, TracesStorageListParams{})
	require.NoError(t, err)
	require.NotNil(t, list)
	require.NotEmpty(t, list)

	// Create Key
	keyDesc := "integration test key"
	key, err := api.CreateKey(ctx, tid, &keyDesc)
	require.NoError(t, err)
	require.NotNil(t, key)
	keyID := key.GetUID()

	// Delete Key
	t.Cleanup(func() {
		err = api.DeleteKey(ctx, tid, keyID)
		require.NoError(t, err)
	})

	// List Keys
	keys, err := api.ListKeys(ctx, tid, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, keys)
	require.NotEmpty(t, keys)

	// Read Key
	keyRead, err := api.ReadKey(ctx, tid, keyID)
	require.NoError(t, err)
	require.NotNil(t, keyRead)
	require.Equal(t, keyID, keyRead.GetUID())

	// Update Key
	updatedKeyDesc := "integration test key updated"
	keyUpdated, err := api.UpdateKey(ctx, tid, keyID, &updatedKeyDesc)
	require.NoError(t, err)
	require.NotNil(t, keyUpdated)
	require.Equal(t, updatedKeyDesc, keyUpdated.GetDescription().Value)

	// SetExpire
	expireDays := 365
	expired, err := api.SetExpire(ctx, tid, expireDays)
	require.NoError(t, err)
	require.NotNil(t, expired)
	require.Equal(t, created.GetID(), expired.GetID())
	require.Equal(t, expireDays, expired.GetRetentionPeriodDays())

	// ReadDailyStats
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	endDate := now
	dailyStats, err := api.ReadDailyStats(ctx, tid, &startDate, &endDate)
	require.NoError(t, err)
	require.NotNil(t, dailyStats)
	// Note: May be empty for newly created resources, which is acceptable

	// ReadMonthlyStats
	currentYear := now.Year()
	monthlyStats, err := api.ReadMonthlyStats(ctx, tid, currentYear)
	require.NoError(t, err)
	require.NotNil(t, monthlyStats)
	// Note: May be empty for newly created resources, which is acceptable
}
