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

func TestLogsStorageOp_List(t *testing.T) {
	expected := v1.PaginatedLogStorageList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.LogStorage{TemplateLogStorage},
	}
	client := newTestClient(expected)
	api := NewLogsStorageOp(client)
	ctx := context.Background()
	params := LogsStoragesListParams{
		IsSystem:             ref(false),
		BucketClassification: ref(v1.LogsStoragesListBucketClassificationShared),
		Status:               ref(v1.LogsStoragesListStatusAssigned),
	}
	tables, err := api.List(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, tables)
	require.Equal(t, 1, len(tables))

	table := tables[0]
	require.Equal(t, TemplateLogStorage.GetName(), table.GetName())
	require.Equal(t, TemplateLogStorage.GetDescription(), table.GetDescription())
	require.Equal(t, TemplateLogStorage.GetIsSystem(), table.GetIsSystem())
	require.Equal(t, TemplateLogStorage.GetAccountID(), table.GetAccountID())
	require.Equal(t, TemplateLogStorage.GetResourceID(), table.GetResourceID())
	require.Equal(t, TemplateLogStorage.GetEndpoints(), table.GetEndpoints())
	require.Equal(t, TemplateLogStorage.GetUsage(), table.GetUsage())
}

func TestLogsStorageOp_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewLogsStorageOp(client)
	ctx := context.Background()
	params := LogsStoragesListParams{}
	_, err := api.List(ctx, params)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permission")
}

func TestLogsStorageOp_Read(t *testing.T) {
	client := newTestClient(TemplateWrappedLogStorage)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "12345")
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedLogStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedLogStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedLogStorage.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateWrappedLogStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedLogStorage.GetResourceID(), actual.GetResourceID())
	require.Equal(t, TemplateLogStorageEndpoints, actual.GetEndpoints())
	require.Equal(t, TemplateLogStorageUsage, actual.GetUsage())
	for i, e := range TemplateWrappedLogStorage.GetTags() {
		a := actual.GetTags()[i]
		require.Equal(t, e, a)
	}
}

func TestLogsStorageOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No LogStorage matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	_, err := api.Read(ctx, "12345")
	require.Error(t, err)
	require.ErrorContains(t, err, "not found")
}

func TestLogsStorageOp_Create(t *testing.T) {
	client := newTestClient(TemplateLogStorage, http.StatusCreated)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	createReq := LogStorageCreateParams{
		Name:        "created-table",
		Description: ref("Created log table"),
		IsSystem:    false,
	}
	actual, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateLogStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateLogStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateLogStorage.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateLogStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateLogStorage.GetResourceID(), actual.GetResourceID())
	require.Equal(t, TemplateLogStorage.GetEndpoints(), actual.GetEndpoints())
	require.Equal(t, TemplateLogStorage.GetUsage(), actual.GetUsage())
	require.Equal(t, TemplateLogStorage.GetTags(), actual.GetTags())
}

func TestLogsStorageOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	createReq := LogStorageCreateParams{
		Name:        "",
		Description: nil,
		IsSystem:    false,
	}
	actual, err := api.Create(ctx, createReq)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid parameter")
}

func TestLogsStorageOp_Update(t *testing.T) {
	client := newTestClient(TemplateWrappedLogStorage)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	updateReq := LogStorageUpdateParams{Name: ref("new name")}
	actual, err := api.Update(ctx, "54321", updateReq)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedLogStorage.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedLogStorage.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedLogStorage.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateWrappedLogStorage.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedLogStorage.GetResourceID(), actual.GetResourceID())
	require.Equal(t, TemplateLogStorageEndpoints, actual.GetEndpoints())
	require.Equal(t, TemplateLogStorageUsage, actual.GetUsage())
}

func TestLogsStorageOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	updateReq := LogStorageUpdateParams{}
	actual, err := api.Update(ctx, "0", updateReq)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid parameter")
}

func TestLogsStorageOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "54321")
	require.NoError(t, err)
}

func TestLogsStorageOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "0")
	require.Error(t, err)
	require.ErrorContains(t, err, " not eligible for deletion")
}

func TestLogsStorageOp_SetExpire(t *testing.T) {
	client := newTestClient(TemplateWrappedLogStorage)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	result, err := api.SetExpire(ctx, "12345", 365)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, TemplateWrappedLogStorage.GetName(), result.GetName())
}

func TestLogsStorageOp_SetExpire_400(t *testing.T) {
	expected := newErrorResponse(400, "invalid parameter, days must be between 1 and 730")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	result, err := api.SetExpire(ctx, "12345", 999)
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}
func TestLogsStorageOp_StatsDaily(t *testing.T) {
	expected := v1.LogStorageDailyUsageBody{
		Usages: []v1.LogStorageDailyUsage{TemplateLogStorageDailyUsage},
	}
	client := newTestClient(expected)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := api.ReadDailyStats(ctx, "12345", &startDate, &endDate)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(result))
}

func TestLogsStorageOp_StatsDaily_400(t *testing.T) {
	expected := newErrorResponse(400, "invalid parameter")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 0, 0, 0, 0, time.UTC)

	result, err := api.ReadDailyStats(ctx, "invalid", &startDate, &endDate)
	require.Nil(t, result)
	require.Error(t, err)
}

func TestLogsStorageOp_StatsMonthly(t *testing.T) {
	expected := v1.LogStorageMonthlyUsageBody{
		Usages: []v1.LogStorageMonthlyUsage{TemplateLogStorageMonthlyUsage},
	}
	client := newTestClient(expected)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	result, err := api.ReadMonthlyStats(ctx, "12345", 2025)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, 1, len(result))
}

func TestLogsStorageOp_StatsMonthly_400(t *testing.T) {
	expected := newErrorResponse(400, "invalid parameter, year must be between 1970 and 2100")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	result, err := api.ReadMonthlyStats(ctx, "99999", 2200)
	require.Nil(t, result)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid parameter")
}

// --- Access Key API tests ---

func TestLogsStorageOp_ListKeys(t *testing.T) {
	expected := v1.PaginatedLogStorageAccessKeyList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.LogStorageAccessKey{TemplateLogStorageAccessKey},
	}
	client := newTestClient(expected)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	keys, err := api.ListKeys(ctx, "12345", nil, nil)
	require.NoError(t, err)
	require.NotNil(t, keys)
	require.Equal(t, 1, len(keys))
	require.Equal(t, TemplateLogStorageAccessKey.GetID(), keys[0].GetID())
	require.Contains(t, keys[0].GetDescription().Value, TemplateLogStorageAccessKey.GetDescription().Value)
}

func TestLogsStorageOp_ListKeys_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	keys, err := api.ListKeys(ctx, "12345", nil, nil)
	require.Nil(t, keys)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogsStorageOp_CreateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey, http.StatusCreated)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.CreateKey(ctx, "12345", ref("new key"))
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedAccessKey.GetID(), key.GetID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedAccessKey.GetDescription().Value)
}

func TestLogsStorageOp_CreateKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.CreateKey(ctx, "12345", nil)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogsStorageOp_ReadKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.ReadKey(ctx, "12345", uuid.New())
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedAccessKey.GetID(), key.GetID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedAccessKey.GetDescription().Value)
}

func TestLogsStorageOp_ReadKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.ReadKey(ctx, "12345", uuid.New())
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogsStorageOp_UpdateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, "12345", uuid.New(), ref("updated key"))
	require.NoError(t, err)
	require.NotNil(t, key)
	require.Equal(t, TemplateWrappedAccessKey.GetID(), key.GetID())
	require.Contains(t, key.GetDescription().Value, TemplateWrappedAccessKey.GetDescription().Value)
}

func TestLogsStorageOp_UpdateKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, "12345", uuid.New(), nil)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogsStorageOp_DeleteKey(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", uuid.New())
	require.NoError(t, err)
}

func TestLogsStorageOp_DeleteKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogStorageIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)
	api := NewLogsStorageOp(client)
	ctx := context.Background()
	tmp := WithLogStorage(t, client, ctx)
	lid := fmt.Sprintf("%d", tmp.GetID())

	// Read
	read, err := api.Read(ctx, lid)
	require.NoError(t, err)
	require.NotNil(t, read)
	require.Equal(t, tmp.GetID(), read.GetID())
	require.Equal(t, tmp.GetName(), read.GetName())

	// List
	params := LogsStoragesListParams{
		IsSystem:             ref(false),
		BucketClassification: ref(v1.LogsStoragesListBucketClassificationShared),
		Status:               ref(v1.LogsStoragesListStatusAssigned),
	}
	list, err := api.List(ctx, params)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.NotEmpty(t, list)

	// Update
	desc := "updated-integrated-test-storage"
	req := LogStorageUpdateParams{Name: &desc}
	updated, err := api.Update(ctx, lid, req)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, desc, updated.GetName().Or("failure"))

	// CreateKey
	ckey, err := api.CreateKey(ctx, lid, ref("integrated-test-key"))
	require.NoError(t, err)
	require.NotNil(t, ckey)
	kid := ckey.GetUID()

	// DeleteKey
	t.Cleanup(func() {
		err = api.DeleteKey(ctx, lid, kid)
		require.NoError(t, err)
	})

	// ListKeys
	keys, err := api.ListKeys(ctx, lid, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, keys)
	require.NotEmpty(t, keys)

	// ReadKey
	rkey, err := api.ReadKey(ctx, lid, kid)
	require.NoError(t, err)
	require.NotNil(t, rkey)
	require.Equal(t, ckey.GetID(), rkey.GetID())

	// UpdateKey
	desc = "updated-integrated-test-key"
	ukey, err := api.UpdateKey(ctx, lid, kid, &desc)
	require.NoError(t, err)
	require.NotNil(t, ukey)
	require.Equal(t, desc, ukey.GetDescription().Or("failure"))

	// SetExpire
	expireDays := 365
	expired, err := api.SetExpire(ctx, lid, expireDays)
	require.NoError(t, err)
	require.NotNil(t, expired)
	require.Equal(t, tmp.GetID(), expired.GetID())
	require.Equal(t, expireDays, expired.GetExpireDay())

	// ReadDailyStats
	now := time.Now()
	startDate := now.AddDate(0, 0, -30)
	endDate := now
	dailyStats, err := api.ReadDailyStats(ctx, lid, &startDate, &endDate)
	require.NoError(t, err)
	require.NotNil(t, dailyStats)
	// Note: May be empty for newly created resources, which is acceptable

	// ReadMonthlyStats
	currentYear := now.Year()
	monthlyStats, err := api.ReadMonthlyStats(ctx, lid, currentYear)
	require.NoError(t, err)
	require.NotNil(t, monthlyStats)
	// Note: May be empty for newly created resources, which is acceptable
}
