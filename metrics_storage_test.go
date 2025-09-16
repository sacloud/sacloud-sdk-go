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
		Description: "Created metrics tank",
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

	createReq := v1.MetricsStorageCreate{
		Name:        "",
		Description: "",
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

	key, err := api.ReadKey(ctx, "12345", "3")
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

	key, err := api.ReadKey(ctx, "12345", "3")
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_UpdateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, "12345", "4", nil)
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

	key, err := api.UpdateKey(ctx, "12345", "4", nil)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_DeleteKey(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", "5")
	require.NoError(t, err)
}

func TestMetricsStorageOp_DeleteKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, "12345", "5")
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}
