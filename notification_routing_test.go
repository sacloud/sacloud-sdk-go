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

func TestNotificationRoutingService_List(t *testing.T) {
	expected := v1.PaginatedNotificationRoutingList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.NotificationRouting{TemplateNotificationRouting},
	}
	client := newTestClient(expected)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	routings, err := api.List(ctx, "12345", ref(20), ref(0))
	require.NoError(t, err)
	require.NotNil(t, routings)
	require.Equal(t, 1, len(routings))

	routing := routings[0]
	require.Equal(t, TemplateNotificationRouting.GetUID(), routing.GetUID())
	require.Equal(t, TemplateNotificationRouting.GetProjectID(), routing.GetProjectID())
	require.Equal(t, TemplateNotificationRouting.GetNotificationTargetUID(), routing.GetNotificationTargetUID())
	require.Equal(t, TemplateNotificationRouting.GetMatchLabels(), routing.GetMatchLabels())
	require.Equal(t, TemplateNotificationRouting.GetResendIntervalMinutes(), routing.GetResendIntervalMinutes())
}

func TestNotificationRoutingService_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	_, err := api.List(ctx, "12345", ref(20), ref(0))
	require.Error(t, err)
	require.ErrorContains(t, err, "request not authorized")
}

func TestNotificationRoutingService_Read(t *testing.T) {
	client := newTestClient(TemplateNotificationRouting)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "12345", uuid.New())
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateNotificationRouting.GetUID(), actual.GetUID())
	require.Equal(t, TemplateNotificationRouting.GetProjectID(), actual.GetProjectID())
	require.Equal(t, TemplateNotificationRouting.GetNotificationTargetUID(), actual.GetNotificationTargetUID())
	require.Equal(t, TemplateNotificationRouting.GetMatchLabels(), actual.GetMatchLabels())
	require.Equal(t, TemplateNotificationRouting.GetResendIntervalMinutes(), actual.GetResendIntervalMinutes())
}

func TestNotificationRoutingService_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No NotificationRouting matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	_, err := api.Read(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "No NotificationRouting matches the given query.")
}

func TestNotificationRoutingService_Create(t *testing.T) {
	nr := TemplateNotificationRouting
	client := newTestClient(nr, http.StatusCreated)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	createParams := NotificationRoutingCreateParams{
		NotificationTargetUID: uuid.New(),
		MatchLabels:           []v1.MatchLabelsItem{},
		ResendIntervalMinutes: ref(60),
	}
	actual, err := api.Create(ctx, "12345", createParams)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, nr.GetUID(), actual.GetUID())
	require.Equal(t, nr.GetProjectID(), actual.GetProjectID())
	require.Equal(t, nr.GetNotificationTargetUID(), actual.GetNotificationTargetUID())
	require.Equal(t, nr.GetMatchLabels(), actual.GetMatchLabels())
	require.Equal(t, nr.GetResendIntervalMinutes(), actual.GetResendIntervalMinutes())
}

func TestNotificationRoutingService_Create_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid create parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	createParams := NotificationRoutingCreateParams{}
	_, err := api.Create(ctx, "12345", createParams)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

func TestNotificationRoutingService_Update(t *testing.T) {
	nr := TemplateNotificationRouting
	client := newTestClient(nr)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	updateParams := NotificationRoutingUpdateParams{
		NotificationTargetUID: ref(uuid.New()),
		MatchLabels:           []v1.MatchLabelsItem{},
		ResendIntervalMinutes: ref(120),
	}
	updated, err := api.Update(ctx, "12345", uuid.New(), updateParams)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, nr.GetUID(), updated.GetUID())
	require.Equal(t, nr.GetProjectID(), updated.GetProjectID())
	require.Equal(t, nr.GetNotificationTargetUID(), updated.GetNotificationTargetUID())
	require.Equal(t, nr.GetMatchLabels(), updated.GetMatchLabels())
	require.Equal(t, nr.GetResendIntervalMinutes(), updated.GetResendIntervalMinutes())
}

func TestNotificationRoutingService_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	updateParams := NotificationRoutingUpdateParams{}
	updated, err := api.Update(ctx, "12345", uuid.New(), updateParams)
	require.Nil(t, updated)
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid update parameters.")
}

func TestNotificationRoutingService_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "12345", uuid.New())
	require.NoError(t, err)
}

func TestNotificationRoutingService_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid delete request.")
}

func TestNotificationRoutingService_Reorder(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	orders := []v1.NotificationRoutingOrder{
		{NotificationRoutingUID: uuid.New(), Order: 1},
		{NotificationRoutingUID: uuid.New(), Order: 2},
	}
	err := api.Reorder(ctx, "12345", orders)
	require.NoError(t, err)
}

func TestNotificationRoutingService_Reorder_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid reorder parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()

	orders := []v1.NotificationRoutingOrder{
		{NotificationRoutingUID: uuid.New(), Order: 1},
	}
	err := api.Reorder(ctx, "12345", orders)
	require.Error(t, err)
	require.ErrorContains(t, err, "Invalid reorder parameters.")
}

// Integration test for NotificationRouting
func TestNotificationRoutingIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)
	api := NewNotificationRoutingOp(client)
	ctx := context.Background()
	project := WithAlertProject(t, client, ctx)
	createdTarget := WithNotificationTarget(t, client, ctx, project.GetID())
	targetID := createdTarget.GetUID()
	id := fmt.Sprintf("%d", project.GetID())

	// Create routing
	createParams := NotificationRoutingCreateParams{
		NotificationTargetUID: targetID,
		MatchLabels: []v1.MatchLabelsItem{
			{Name: "severity", Value: "critical"},
		},
		ResendIntervalMinutes: ref(30),
	}
	created, err := api.Create(ctx, id, createParams)
	require.NoError(t, err)
	require.NotNil(t, created)
	rid := created.GetUID()

	// Cleanup routing
	t.Cleanup(func() {
		err := api.Delete(ctx, id, rid)
		require.NoError(t, err)
	})

	// Read
	read, err := api.Read(ctx, id, rid)
	require.NoError(t, err)
	require.NotNil(t, read)
	require.Equal(t, created.GetUID(), read.GetUID())
	require.Equal(t, created.GetProjectID(), read.GetProjectID())
	require.Equal(t, created.GetNotificationTargetUID(), read.GetNotificationTargetUID())
	require.Equal(t, created.GetMatchLabels(), read.GetMatchLabels())
	require.Equal(t, created.GetResendIntervalMinutes(), read.GetResendIntervalMinutes())

	// List
	list, err := api.List(ctx, id, nil, nil)
	require.NoError(t, err)
	require.NotNil(t, list)
	require.NotEmpty(t, list)

	// Update
	newInterval := 60
	updateParams := NotificationRoutingUpdateParams{
		ResendIntervalMinutes: &newInterval,
	}
	updated, err := api.Update(ctx, id, rid, updateParams)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, newInterval, updated.GetResendIntervalMinutes().Or(0))

	// Reorder
	orders := []v1.NotificationRoutingOrder{
		{NotificationRoutingUID: rid, Order: 1},
	}
	err = api.Reorder(ctx, id, orders)
	require.NoError(t, err)
}
