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

	"github.com/google/uuid"
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

	routings, err := api.List(ctx, MetricsRoutingsListParams{})
	require.NoError(t, err)
	require.NotNil(t, routings)
	require.Equal(t, 1, len(routings))
}

func TestMetricsRoutingOp_Read(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsRouting)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	res, err := api.Read(ctx, uuid.New())
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
	expected := newErrorResponse(404, "No MetricsRouting matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	routing, err := api.Read(ctx, uuid.New())
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "internal server error")
}

func TestMetricsRoutingOp_Create(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsRouting, http.StatusCreated)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	createReq := MetricsRoutingCreateParams{
		PublisherCode:    "appliance",
		Variant:          "variant",
		MetricsStorageID: "12355",
	}
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

	createReq := MetricsRoutingCreateParams{}
	routing, err := api.Create(ctx, createReq)
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

func TestMetricsRoutingOp_Update(t *testing.T) {
	client := newTestClient(TemplateWrappedMetricsRouting)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	updateReq := MetricsRoutingUpdateParams{
		PublisherCode:    ref("appliance"),
		Variant:          ref("variant"),
		MetricsStorageID: ref("12355"),
	}
	res, err := api.Update(ctx, uuid.New(), updateReq)
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
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	updateReq := MetricsRoutingUpdateParams{}
	routing, err := api.Update(ctx, uuid.New(), updateReq)
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

func TestMetricsRoutingOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, uuid.New())
	require.NoError(t, err)
}

func TestMetricsRoutingOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "not eligible for deletion")
}

func TestMetricsRoutingIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)
	api := NewMetricsRoutingOp(client)
	ctx := context.Background()

	// we need to obtain a sane publisher object
	publisherOp := NewPublisherOp(client)
	publishers, err := publisherOp.List(ctx, nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, publishers)
	var pub *v1.Publisher
	var v *v1.PublisherVariant
	for _, p := range publishers {
		for _, q := range p.GetVariants() {
			if q.GetStorage() == v1.PublisherVariantStorageMetrics {
				pub = &p
				v = &q
				break
			}
		}
	}
	require.NotNil(t, pub)
	require.NotNil(t, v)

	// aaand a storage
	storage := WithMetricsStorage(t, client, ctx)
	require.NotNil(t, storage)
	sid := fmt.Sprintf("%d", storage.GetID())

	// Create
	createReq := MetricsRoutingCreateParams{
		PublisherCode:    pub.GetCode(),
		Variant:          v.GetName(),
		MetricsStorageID: sid,
	}
	created, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, created)
	rid := created.GetUID()

	// Delete
	t.Cleanup(func() {
		err := api.Delete(ctx, rid)
		require.NoError(t, err)
	})

	// Read
	read, err := api.Read(ctx, rid)
	require.NoError(t, err)
	require.NotNil(t, read)
	require.Equal(t, created.GetID(), read.GetID())
	require.Equal(t, created.GetPublisher(), read.GetPublisher())
	require.Equal(t, created.GetPublisherCode(), read.GetPublisherCode())
	require.Equal(t, created.GetMetricsStorage(), read.GetMetricsStorage())
	require.Equal(t, created.GetMetricsStorageID(), read.GetMetricsStorageID())
	require.Equal(t, created.GetResourceID(), read.GetResourceID())
	require.Equal(t, created.GetVariant(), read.GetVariant())

	// List
	routings, err := api.List(ctx, MetricsRoutingsListParams{})
	require.NoError(t, err)
	require.NotEmpty(t, routings)

	// Update
	updateReq := MetricsRoutingUpdateParams{
		ResourceID: ref("12345"),
	}
	updated, err := api.Update(ctx, rid, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, "12345", updated.GetResourceID().Or(^0))
}
