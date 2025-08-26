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
	expected := v1.PaginatedMetricsTankList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.MetricsTank{TemplateMetricsTank},
	}
	client := newTestClient(expected)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	tanks, err := api.List(ctx, 1, 0)
	require.NoError(t, err)
	require.NotNil(t, tanks)
	require.Equal(t, 1, len(tanks))

	tank := tanks[0]
	require.Equal(t, TemplateMetricsTank.GetName(), tank.GetName())
	require.Equal(t, TemplateMetricsTank.GetDescription(), tank.GetDescription())
	require.Equal(t, TemplateMetricsTank.GetIsSystem(), tank.GetIsSystem())
	require.Equal(t, TemplateMetricsTank.GetAccountID(), tank.GetAccountID())
	require.Equal(t, TemplateMetricsTank.GetResourceID(), tank.GetResourceID())
	require.Equal(t, TemplateMetricsTank.GetEndpoints(), tank.GetEndpoints())
	require.Equal(t, TemplateMetricsTank.GetUsage(), tank.GetUsage())
}

func TestMetricsStorageOp_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	tanks, err := api.List(ctx, -1, -1)
	require.Nil(t, tanks)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_Read(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsTank)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, 12345)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedMetricsTank.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedMetricsTank.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedMetricsTank.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateWrappedMetricsTank.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedMetricsTank.GetResourceID(), actual.GetResourceID())
	require.Equal(t, (&TemplateWrappedMetricsTank.Endpoints).GetAddress(), (&actual.Endpoints).GetAddress())
	require.Equal(t, (&TemplateWrappedMetricsTank.Usage).GetMetricsRoutings(), (&actual.Usage).GetMetricsRoutings())
	require.Equal(t, (&TemplateWrappedMetricsTank.Usage).GetAlertRules(), (&actual.Usage).GetAlertRules())
	require.Equal(t, (&TemplateWrappedMetricsTank.Usage).GetLogRecordingRules(), (&actual.Usage).GetLogRecordingRules())
	require.Equal(t, TemplateWrappedMetricsTank.GetTags(), actual.GetTags())
}

func TestMetricsStorageOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No MetricsTank matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, 99999)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "not found")
}

func TestMetricsStorageOp_Create(t *testing.T) {
	client := newTestClient(TemplateMetricsTank, http.StatusCreated)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	createReq := v1.MetricsTankCreate{
		Name:        "created-tank",
		Description: "Created metrics tank",
		IsSystem:    false,
	}
	actual, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateMetricsTank.GetName(), actual.GetName())
	require.Equal(t, TemplateMetricsTank.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateMetricsTank.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateMetricsTank.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateMetricsTank.GetResourceID(), actual.GetResourceID())
	require.Equal(t, (&TemplateMetricsTank.Endpoints).GetAddress(), (&actual.Endpoints).GetAddress())
	require.Equal(t, (&TemplateMetricsTank.Usage).GetMetricsRoutings(), (&actual.Usage).GetMetricsRoutings())
	require.Equal(t, (&TemplateMetricsTank.Usage).GetAlertRules(), (&actual.Usage).GetAlertRules())
	require.Equal(t, (&TemplateMetricsTank.Usage).GetLogRecordingRules(), (&actual.Usage).GetLogRecordingRules())
	require.Equal(t, TemplateMetricsTank.GetTags(), actual.GetTags())
}

func TestMetricsStorageOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	createReq := v1.MetricsTankCreate{
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
	client := newTestClient(TemplateWrappedMetricsTank)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Update(ctx, 54321, &TemplateMetricsTank)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateWrappedMetricsTank.GetName(), actual.GetName())
	require.Equal(t, TemplateWrappedMetricsTank.GetDescription(), actual.GetDescription())
	require.Equal(t, TemplateWrappedMetricsTank.GetIsSystem(), actual.GetIsSystem())
	require.Equal(t, TemplateWrappedMetricsTank.GetAccountID(), actual.GetAccountID())
	require.Equal(t, TemplateWrappedMetricsTank.GetResourceID(), actual.GetResourceID())
	require.Equal(t, (&TemplateWrappedMetricsTank.Endpoints).GetAddress(), (&actual.Endpoints).GetAddress())
	require.Equal(t, (&TemplateWrappedMetricsTank.Usage).GetMetricsRoutings(), (&actual.Usage).GetMetricsRoutings())
	require.Equal(t, (&TemplateWrappedMetricsTank.Usage).GetAlertRules(), (&actual.Usage).GetAlertRules())
	require.Equal(t, (&TemplateWrappedMetricsTank.Usage).GetLogRecordingRules(), (&actual.Usage).GetLogRecordingRules())
}

func TestMetricsStorageOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	actual, err := api.Update(ctx, 0, &TemplateMetricsTank)
	require.Nil(t, actual)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

func TestMetricsStorageOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, 54321)
	require.NoError(t, err)
}

func TestMetricsStorageOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, 0)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}

// --- Access Key API tests ---

func TestMetricsStorageOp_ListKeys(t *testing.T) {
	expected := v1.PaginatedMetricsTankAccessKeyList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.MetricsTankAccessKey{TemplateAccessKey},
	}
	client := newTestClient(expected)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	keys, err := api.ListKeys(ctx, 12345, 1, 0)
	require.NoError(t, err)
	require.NotNil(t, keys)
	require.Equal(t, 1, len(keys))
	require.Equal(t, TemplateAccessKey.GetID(), keys[0].GetID())
	require.Contains(t, keys[0].GetDescription().Value, TemplateAccessKey.GetDescription().Value)
}

func TestMetricsStorageOp_ListKeys_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	keys, err := api.ListKeys(ctx, 12345, 1, 0)
	require.Nil(t, keys)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_CreateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey, http.StatusCreated)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.CreateKey(ctx, 12345, &TemplateAccessKey)
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

	key, err := api.CreateKey(ctx, 12345, &TemplateAccessKey)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_ReadKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.ReadKey(ctx, 12345, 3)
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

	key, err := api.ReadKey(ctx, 12345, 3)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_UpdateKey(t *testing.T) {
	client := newTestClient(TemplateWrappedAccessKey)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	key, err := api.UpdateKey(ctx, 12345, 4, &TemplateAccessKey)
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

	key, err := api.UpdateKey(ctx, 12345, 4, &TemplateAccessKey)
	require.Nil(t, key)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}

func TestMetricsStorageOp_DeleteKey(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, 12345, 5)
	require.NoError(t, err)
}

func TestMetricsStorageOp_DeleteKey_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewMetricsStorageOp(client)
	ctx := context.Background()

	err := api.DeleteKey(ctx, 12345, 5)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permissions")
}
