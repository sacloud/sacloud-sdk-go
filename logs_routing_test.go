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

func TestLogRoutingOp_List(t *testing.T) {
	expected := v1.PaginatedLogRoutingList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.LogRouting{TemplateLogRouting},
	}
	client := newTestClient(expected)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	routings, err := api.List(ctx, LogsRoutingsListParams{})
	require.NoError(t, err)
	require.NotNil(t, routings)
	require.Equal(t, 1, len(routings))
}

func TestLogRoutingOp_Read(t *testing.T) {
	client := newTestClient(TemplateWrappedLogRouting)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	res, err := api.Read(ctx, uuid.New())
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, TemplateWrappedLogRouting.GetID(), res.GetID())
	require.Equal(t, TemplateWrappedLogRouting.GetPublisher(), res.GetPublisher())
	require.Equal(t, TemplateWrappedLogRouting.GetPublisherCode(), res.GetPublisherCode())
	require.Equal(t, TemplateWrappedLogRouting.GetLogStorage(), res.GetLogStorage())
	require.Equal(t, TemplateWrappedLogRouting.GetLogStorageID(), res.GetLogStorageID())
	require.Equal(t, TemplateWrappedLogRouting.GetResourceID(), res.GetResourceID())
	require.Equal(t, TemplateWrappedLogRouting.GetVariant(), res.GetVariant())
}

func TestLogRoutingOp_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No LogRouting matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	routing, err := api.Read(ctx, uuid.New())
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "No LogRouting matches the given query.")
}

func TestLogRoutingOp_Create(t *testing.T) {
	client := newTestClient(TemplateWrappedLogRouting, http.StatusCreated)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	createReq := LogsRoutingCreateParams{
		PublisherCode: "appliance",
		Variant:       "variant",
		LogStorageID:  "12345",
	}
	res, err := api.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, TemplateWrappedLogRouting.GetID(), res.GetID())
	require.Equal(t, TemplateWrappedLogRouting.GetPublisher(), res.GetPublisher())
	require.Equal(t, TemplateWrappedLogRouting.GetPublisherCode(), res.GetPublisherCode())
	require.Equal(t, TemplateWrappedLogRouting.GetLogStorage(), res.GetLogStorage())
	require.Equal(t, TemplateWrappedLogRouting.GetLogStorageID(), res.GetLogStorageID())
	require.Equal(t, TemplateWrappedLogRouting.GetResourceID(), res.GetResourceID())
	require.Equal(t, TemplateWrappedLogRouting.GetVariant(), res.GetVariant())
}

func TestLogRoutingOp_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid request body.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	createReq := LogsRoutingCreateParams{
		LogStorageID: "0",
	}
	routing, err := api.Create(ctx, createReq)
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid request body.")
}

func TestLogRoutingOp_Update(t *testing.T) {
	client := newTestClient(TemplateWrappedLogRouting)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	updateReq := LogsRoutingUpdateParams{
		PublisherCode: ref("appliance"),
		Variant:       ref("variant"),
		LogStorageID:  ref("12345"),
	}
	res, err := api.Update(ctx, uuid.New(), updateReq)
	require.NoError(t, err)
	require.NotNil(t, res)
	require.Equal(t, TemplateWrappedLogRouting.GetID(), res.GetID())
	require.Equal(t, TemplateWrappedLogRouting.GetPublisher(), res.GetPublisher())
	require.Equal(t, TemplateWrappedLogRouting.GetPublisherCode(), res.GetPublisherCode())
	require.Equal(t, TemplateWrappedLogRouting.GetLogStorage(), res.GetLogStorage())
	require.Equal(t, TemplateWrappedLogRouting.GetLogStorageID(), res.GetLogStorageID())
	require.Equal(t, TemplateWrappedLogRouting.GetResourceID(), res.GetResourceID())
	require.Equal(t, TemplateWrappedLogRouting.GetVariant(), res.GetVariant())
}

func TestLogRoutingOp_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	routing, err := api.Update(ctx, uuid.New(), LogsRoutingUpdateParams{})
	require.Nil(t, routing)
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid update parameters.")
}

func TestLogRoutingOp_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, uuid.New())
	require.NoError(t, err)
}

func TestLogRoutingOp_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid delete request.")
}

func TestLogRoutingIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)
	api := NewLogRoutingOp(client)
	ctx := context.Background()

	// obtain a sane publisher object
	publisherOp := NewPublisherOp(client)
	publishers, err := publisherOp.List(ctx, nil, nil)
	require.NoError(t, err)
	require.NotEmpty(t, publishers)
	var pub *v1.Publisher
	var v *v1.PublisherVariant
	for _, p := range publishers {
		for _, q := range p.GetVariants() {
			if q.GetStorage() == v1.PublisherVariantStorageLogs {
				pub = &p
				v = &q
				break
			}
		}
	}
	require.NotNil(t, pub)
	require.NotNil(t, v)

	// and a storage
	storage := WithLogStorage(t, client, ctx)
	require.NotNil(t, storage)
	sid := fmt.Sprintf("%d", storage.GetID())

	// Create
	createReq := LogsRoutingCreateParams{
		PublisherCode: pub.GetCode(),
		Variant:       v.GetName(),
		LogStorageID:  sid,
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
	require.Equal(t, created.GetLogStorage(), read.GetLogStorage())
	require.Equal(t, created.GetLogStorageID(), read.GetLogStorageID())
	require.Equal(t, created.GetResourceID(), read.GetResourceID())
	require.Equal(t, created.GetVariant(), read.GetVariant())

	// List
	routings, err := api.List(ctx, LogsRoutingsListParams{})
	require.NoError(t, err)
	require.NotEmpty(t, routings)

	// Update
	updateReq := LogsRoutingUpdateParams{
		ResourceID: ref("12345"),
	}
	updated, err := api.Update(ctx, rid, updateReq)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, "12345", updated.GetResourceID().Or(^0))
}
