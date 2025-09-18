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
	"net/url"
	"testing"

	"github.com/google/uuid"
	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestNotificationTargetService_List(t *testing.T) {
	expected := v1.PaginatedNotificationTargetList{
		IsOk:    v1.NewOptBool(true),
		Count:   1,
		From:    0,
		Results: []v1.NotificationTarget{TemplateNotificationTarget},
	}
	client := newTestClient(expected)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()
	params := NotificationTargetsListParams{
		Count: ref[int64](20),
		From:  ref[int64](0),
	}
	targets, err := api.List(ctx, "12345", params)
	require.NoError(t, err)
	require.NotNil(t, targets)
	require.Equal(t, 1, len(targets))

	target := targets[0]
	require.Equal(t, TemplateNotificationTarget.GetUID(), target.GetUID())
	require.Equal(t, TemplateNotificationTarget.GetProjectID(), target.GetProjectID())
	require.Equal(t, TemplateNotificationTarget.GetServiceType(), target.GetServiceType())
	require.Equal(t, TemplateNotificationTarget.GetURL(), target.GetURL())
	require.Equal(t, TemplateNotificationTarget.GetConfig(), target.GetConfig())
	require.Equal(t, TemplateNotificationTarget.GetDescription(), target.GetDescription())
}

func TestNotificationTargetService_List_403(t *testing.T) {
	expected := newErrorResponse(403, "request not authorized")
	client := newTestClient(expected, http.StatusForbidden)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()
	params := NotificationTargetsListParams{}
	_, err := api.List(ctx, "12345", params)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient permission")
}

func TestNotificationTargetService_Read(t *testing.T) {
	client := newTestClient(TemplateNotificationTarget)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()

	actual, err := api.Read(ctx, "12345", uuid.New())
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, TemplateNotificationTarget.GetUID(), actual.GetUID())
	require.Equal(t, TemplateNotificationTarget.GetProjectID(), actual.GetProjectID())
	require.Equal(t, TemplateNotificationTarget.GetServiceType(), actual.GetServiceType())
	require.Equal(t, TemplateNotificationTarget.GetURL(), actual.GetURL())
	require.Equal(t, TemplateNotificationTarget.GetConfig(), actual.GetConfig())
	require.Equal(t, TemplateNotificationTarget.GetDescription(), actual.GetDescription())
}

func TestNotificationTargetService_Read_404(t *testing.T) {
	expected := newErrorResponse(404, "No NotificationTarget matches the given query.")
	client := newTestClient(expected, http.StatusNotFound)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()

	_, err := api.Read(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "Not Found")
}

func TestNotificationTargetService_Create(t *testing.T) {
	nt := TemplateNotificationTarget
	client := newTestClient(nt, http.StatusCreated)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()

	url, _ := url.Parse("https://example.com/notify")
	createParams := NotificationTargetCreateParams{
		ServiceType: v1.NotificationTargetServiceTypeSAKURASIMPLENOTICE,
		URL:         *url,
	}
	actual, err := api.Create(ctx, "12345", createParams)
	require.NoError(t, err)
	require.NotNil(t, actual)
	require.Equal(t, nt.GetUID(), actual.GetUID())
	require.Equal(t, nt.GetProjectID(), actual.GetProjectID())
	require.Equal(t, nt.GetServiceType(), actual.GetServiceType())
	require.Equal(t, nt.GetURL(), actual.GetURL())
	require.Equal(t, nt.GetConfig(), actual.GetConfig())
	require.Equal(t, nt.GetDescription(), actual.GetDescription())
}

// Update
func TestNotificationTargetService_Update(t *testing.T) {
	nt := TemplateNotificationTarget
	client := newTestClient(nt)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()

	updateParams := NotificationTargetUpdateParams{
		ServiceType: func() *v1.PatchedNotificationTargetServiceType {
			v := v1.PatchedNotificationTargetServiceType(nt.GetServiceType())
			return &v
		}(),
		URL:         func() *string { v := nt.GetURL(); return &v }(),
		Description: func() *string { v := nt.GetDescription().Or(""); return &v }(),
	}
	updated, err := api.Update(ctx, "12345", uuid.New(), updateParams)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, nt.GetUID(), updated.GetUID())
	require.Equal(t, nt.GetProjectID(), updated.GetProjectID())
	require.Equal(t, nt.GetServiceType(), updated.GetServiceType())
	require.Equal(t, nt.GetURL(), updated.GetURL())
	require.Equal(t, nt.GetConfig(), updated.GetConfig())
	require.Equal(t, nt.GetDescription(), updated.GetDescription())
}

func TestNotificationTargetService_Update_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid update parameters.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()

	updateParams := NotificationTargetUpdateParams{}
	updated, err := api.Update(ctx, "12345", uuid.New(), updateParams)
	require.Nil(t, updated)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid")
}

// Delete
func TestNotificationTargetService_Delete(t *testing.T) {
	client := newTestClient(nil, http.StatusNoContent)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "12345", uuid.New())
	require.NoError(t, err)
}

func TestNotificationTargetService_Delete_400(t *testing.T) {
	expected := newErrorResponse(400, "Invalid delete request.")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewNotificationTargetOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "12345", uuid.New())
	require.Error(t, err)
	require.ErrorContains(t, err, "Bad Request")
}
