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
	params := v1.LogsStoragesListParams{
		AccountID:            v1.NewOptString("account-id-12345"),
		BucketClassification: v1.NewOptLogsStoragesListBucketClassification(v1.LogsStoragesListBucketClassificationSeparated),
		Count:                v1.NewOptInt(20),
		From:                 v1.NewOptInt(0),
		IsSystem:             v1.NewOptBool(false),
		ResourceID:           v1.NewOptInt(12345),
		Status:               v1.NewOptLogsStoragesListStatus(v1.LogsStoragesListStatusAssigned),
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
	params := v1.LogsStoragesListParams{}
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
	require.ErrorContains(t, err, "Not Found")
}

func TestLogsStorageOp_Create(t *testing.T) {
	client := newTestClient(TemplateLogStorage, http.StatusCreated)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	createReq := v1.LogStorageCreate{
		Name:        "created-table",
		Description: "Created log table",
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

	createReq := v1.LogStorageCreate{
		Name:        "",
		Description: "",
		IsSystem:    false,
	}
	actual, err := api.Create(ctx, createReq)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestLogsStorageOp_Update(t *testing.T) {
	client := newTestClient(TemplateWrappedLogStorage)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Update(ctx, "54321", &TemplateLogStorage)
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

	actual, err := api.Update(ctx, "0", &TemplateLogStorage)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
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
	require.ErrorContains(t, err, "Bad Request")
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

	keys, err := api.ListKeys(ctx, "12345", 1, 0)
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

	keys, err := api.ListKeys(ctx, "12345", 1, 0)
	require.Nil(t, keys)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogsStorageOp_CreateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey, http.StatusCreated)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.CreateKey(ctx, "12345", &TemplateLogStorageAccessKey)
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

	key, err := api.CreateKey(ctx, "12345", &TemplateLogStorageAccessKey)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogsStorageOp_ReadKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.ReadKey(ctx, "12345", "3")
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

	key, err := api.ReadKey(ctx, "12345", "3")
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogsStorageOp_UpdateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, "12345", "4", &TemplateLogStorageAccessKey)
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

	key, err := api.UpdateKey(ctx, "12345", "4", &TemplateLogStorageAccessKey)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestLogsStorageOp_DeleteKey(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", "5")
	require.NoError(t, err)
}

func TestLogsStorageOp_DeleteKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewLogsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", "5")
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}
