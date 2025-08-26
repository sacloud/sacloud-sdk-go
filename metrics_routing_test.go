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

func TestMetricsRoutingOp_List(t *testing.T) {
	expected := v1.PaginatedMetricsRoutingList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.MetricsRouting{TemplateMetricsRouting},
	}
	client := newTestClient(expected)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	routings, err := api.List(ctx, v1.MetricsRoutingsListParams{})
	require.NoError(t, err)
	require.NotNil(t, routings)
	require.Equal(t, 1, len(routings))
}

func TestMetricsRoutingOp_Read(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsRouting)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	res, err := api.Read(ctx, TemplateWrappedMetricsRouting.GetID())
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, TemplateWrappedMetricsRouting.GetID(), res.GetID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetPublisher(), res.GetPublisher())
	require.Equal(t, TemplateWrappedMetricsRouting.GetPublisherCode(), res.GetPublisherCode())
	require.Equal(t, TemplateWrappedMetricsRouting.GetMetricsStorage(), res.GetMetricsStorage())
	require.Equal(t, TemplateWrappedMetricsRouting.GetMetricsStorageID(), res.GetMetricsStorageID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetResourceID(), res.GetResourceID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetVariant(), res.GetVariant())
}

func TestMetricsRoutingOp_Read_404(t *testing.T) {
	expected := ErrorResponse{
		Code:    "not_found",
		Message: "No MetricsRouting matches the given query.",
		IsOk:    false,
		Status:  404,
	}
	client := newTestClient(expected, http.StatusNotFound)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	routing, err := api.Read(ctx, 99999)
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "Not Found")
}

func TestMetricsRoutingOp_Create(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsRouting, http.StatusCreated)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	createReq := TemplateMetricsRouting
	res, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, TemplateWrappedMetricsRouting.GetID(), res.GetID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetPublisher(), res.GetPublisher())
	require.Equal(t, TemplateWrappedMetricsRouting.GetPublisherCode(), res.GetPublisherCode())
	require.Equal(t, TemplateWrappedMetricsRouting.GetMetricsStorage(), res.GetMetricsStorage())
	require.Equal(t, TemplateWrappedMetricsRouting.GetMetricsStorageID(), res.GetMetricsStorageID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetResourceID(), res.GetResourceID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetVariant(), res.GetVariant())
}

func TestMetricsRoutingOp_Create_400(t *testing.T) {
	expected := ErrorResponse{
		Code:    "bad_request",
		Message: "Invalid request body.",
		IsOk:    false,
		Status:  400,
	}
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	createReq := v1.MetricsRouting{}
	routing, err := api.Create(ctx, createReq)
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

func TestMetricsRoutingOp_Update(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsRouting)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	updateReq := TemplateMetricsRouting
	res, err := api.Update(ctx, TemplateWrappedMetricsRouting.GetID(), &updateReq)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, TemplateWrappedMetricsRouting.GetID(), res.GetID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetPublisher(), res.GetPublisher())
	require.Equal(t, TemplateWrappedMetricsRouting.GetPublisherCode(), res.GetPublisherCode())
	require.Equal(t, TemplateWrappedMetricsRouting.GetMetricsStorage(), res.GetMetricsStorage())
	require.Equal(t, TemplateWrappedMetricsRouting.GetMetricsStorageID(), res.GetMetricsStorageID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetResourceID(), res.GetResourceID())
	require.Equal(t, TemplateWrappedMetricsRouting.GetVariant(), res.GetVariant())
}

func TestMetricsRoutingOp_Update_400(t *testing.T) {
	expected := ErrorResponse{
		Code:    "bad_request",
		Message: "Invalid update parameters.",
		IsOk:    false,
		Status:  400,
	}
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	updateReq := v1.MetricsRouting{}
	routing, err := api.Update(ctx, 0, &updateReq)
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

func TestMetricsRoutingOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, TemplateMetricsRouting.GetID())
	require.NoError(t, err)
}

func TestMetricsRoutingOp_Delete_400(t *testing.T) {
	expected := ErrorResponse{
		Code:    "bad_request",
		Message: "Invalid delete request.",
		IsOk:    false,
		Status:  400,
	}
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, 0)
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}
