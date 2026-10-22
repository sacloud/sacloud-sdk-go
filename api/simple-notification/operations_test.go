// Copyright 2026- The sacloud/simple-notification-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package simplenotification_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	simplenotification "github.com/sacloud/sacloud-sdk-go/api/simple-notification"
	v1 "github.com/sacloud/sacloud-sdk-go/api/simple-notification/apis/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServiceItemRequests(t *testing.T) {
	tests := []struct {
		name     string
		class    string
		settings string
		create   func(context.Context, *v1.Client, v1.PostCommonServiceItemRequest) error
		update   func(context.Context, *v1.Client, v1.PutCommonServiceItemRequest) error
	}{
		{
			name:     "destination",
			class:    "saknoticedestination",
			settings: `{"Type":"email","Value":"test@example.com"}`,
			create: func(ctx context.Context, client *v1.Client, request v1.PostCommonServiceItemRequest) error {
				_, err := simplenotification.NewDestinationOp(client).Create(ctx, request)
				return err
			},
			update: func(ctx context.Context, client *v1.Client, request v1.PutCommonServiceItemRequest) error {
				_, err := simplenotification.NewDestinationOp(client).Update(ctx, "123", request)
				return err
			},
		},
		{
			name:     "group",
			class:    "saknoticegroup",
			settings: `{"Destinations":["123700010001"]}`,
			create: func(ctx context.Context, client *v1.Client, request v1.PostCommonServiceItemRequest) error {
				_, err := simplenotification.NewGroupOp(client).Create(ctx, request)
				return err
			},
			update: func(ctx context.Context, client *v1.Client, request v1.PutCommonServiceItemRequest) error {
				_, err := simplenotification.NewGroupOp(client).Update(ctx, "123", request)
				return err
			},
		},
		{
			name:     "routing",
			class:    "saknoticerouting",
			settings: `{"MatchLabels":[{"Name":"environment","Value":"test"}],"SourceID":"789","TargetGroupID":"456","PriorityRank":1}`,
			create: func(ctx context.Context, client *v1.Client, request v1.PostCommonServiceItemRequest) error {
				_, err := simplenotification.NewRoutingOp(client).Create(ctx, request)
				return err
			},
			update: func(ctx context.Context, client *v1.Client, request v1.PutCommonServiceItemRequest) error {
				_, err := simplenotification.NewRoutingOp(client).Update(ctx, "123", request)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for _, operation := range []string{"create", "update", "update without settings"} {
				t.Run(operation, func(t *testing.T) {
					called := false
					server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						called = true
						var body struct {
							Item struct {
								Provider struct {
									Class        string `json:"Class"`
									ServiceClass string `json:"ServiceClass"`
								} `json:"Provider"`
								ServiceClass string          `json:"ServiceClass"`
								Settings     json.RawMessage `json:"Settings"`
								Tags         json.RawMessage `json:"Tags"`
							} `json:"CommonServiceItem"`
						}
						if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
							t.Error(err)
							w.WriteHeader(http.StatusBadRequest)
							return
						}
						assert.JSONEq(t, `[]`, string(body.Item.Tags))
						if operation == "update without settings" {
							assert.Empty(t, body.Item.Settings)
						} else {
							assert.JSONEq(t, tt.settings, string(body.Item.Settings))
						}
						w.Header().Set("Content-Type", "application/json")
						if operation == "create" {
							assert.Equal(t, http.MethodPost, r.Method)
							assert.Equal(t, "/commonserviceitem", r.URL.Path)
							assert.Equal(t, tt.class, body.Item.Provider.Class)
							assert.Equal(t, "cloud/"+tt.class+"/2", body.Item.ServiceClass)
							if tt.name == "group" {
								assert.Equal(t, "cloud/saknotice", body.Item.Provider.ServiceClass)
							}
							w.WriteHeader(http.StatusCreated)
						} else {
							assert.Equal(t, http.MethodPut, r.Method)
							assert.Equal(t, "/commonserviceitem/123", r.URL.Path)
						}
						_, _ = w.Write([]byte(`{"CommonServiceItem":{"ID":"123","Name":"test","Description":"","Settings":` +
							tt.settings + `,"CreatedAt":"2026-10-19T00:00:00Z","ModifiedAt":"2026-10-19T00:00:00Z",` +
							`"Provider":{"Class":"` + tt.class + `"},"Icon":null,"Tags":[]}}`))
					}))
					defer server.Close()
					client, err := v1.NewClient(server.URL, v1.WithClient(server.Client()))
					require.NoError(t, err)

					if operation == "create" {
						var settings v1.PostCommonServiceItemRequestCommonServiceItemSettings
						require.NoError(t, json.Unmarshal([]byte(tt.settings), &settings))
						settings.Type = "" // The facade selects the request variant.
						err = tt.create(t.Context(), client, v1.PostCommonServiceItemRequest{
							CommonServiceItem: v1.PostCommonServiceItemRequestCommonServiceItem{
								Name: "test", Tags: []string{}, Icon: v1.NilCommonServiceItemIcon{Null: true},
								Settings: settings,
							},
						})
					} else {
						var settings v1.OptPutCommonServiceItemRequestCommonServiceItemSettings
						if operation == "update" {
							require.NoError(t, json.Unmarshal([]byte(tt.settings), &settings.Value))
							settings.Set = true
							settings.Value.Type = ""
						}
						err = tt.update(t.Context(), client, v1.PutCommonServiceItemRequest{
							CommonServiceItem: v1.PutCommonServiceItemRequestCommonServiceItem{
								Name: "test", Tags: []string{}, Icon: v1.NilCommonServiceItemIcon{Null: true},
								Settings: settings,
							},
						})
					}
					require.NoError(t, err)
					require.True(t, called)
				})
			}
		})
	}
}

func TestRenamedOperations(t *testing.T) {
	const itemResponse = `{"CommonServiceItem":{"ID":"123","Name":"test","Description":"","Settings":{"Type":"email","Value":"test@example.com"},` +
		`"CreatedAt":"2026-10-19T00:00:00Z","ModifiedAt":"2026-10-19T00:00:00Z","Provider":{"Class":"saknoticedestination"},"Icon":null,"Tags":[]}}`

	tests := []struct {
		name     string
		method   string
		path     string
		status   int
		body     string
		response string
		invoke   func(*testing.T, context.Context, *v1.Client) error
	}{
		{
			name: "destination read", method: http.MethodGet,
			path: "/commonserviceitem/123", status: http.StatusOK, response: itemResponse,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				res, err := simplenotification.NewDestinationOp(client).Read(ctx, "123")
				if err == nil {
					require.Equal(t, "123", res.CommonServiceItem.ID)
				}
				return err
			},
		},
		{
			name: "group read", method: http.MethodGet,
			path: "/commonserviceitem/123", status: http.StatusOK, response: itemResponse,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				res, err := simplenotification.NewGroupOp(client).Read(ctx, "123")
				if err == nil {
					require.Equal(t, "123", res.CommonServiceItem.ID)
				}
				return err
			},
		},
		{
			name: "routing read", method: http.MethodGet,
			path: "/commonserviceitem/123", status: http.StatusOK, response: itemResponse,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				res, err := simplenotification.NewRoutingOp(client).Read(ctx, "123")
				if err == nil {
					require.Equal(t, "123", res.CommonServiceItem.ID)
				}
				return err
			},
		},
		{
			name: "destination status", method: http.MethodGet,
			path: "/commonserviceitem/123/simplenotification/status", status: http.StatusOK,
			response: `{"NotificationStatus":{"IsValid":true,"ModifiedAt":"2026-10-19T00:00:00Z"}}`,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				res, err := simplenotification.NewDestinationOp(client).ReadStatus(ctx, "123")
				if err == nil {
					require.True(t, res.NotificationStatus.IsValid)
				}
				return err
			},
		},
		{
			name: "send message", method: http.MethodPost,
			path: "/commonserviceitem/123/simplenotification/message", status: http.StatusAccepted,
			body: `{"Message":"test"}`, response: `{"is_ok":true}`,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				res, err := simplenotification.NewGroupOp(client).SendMessage(ctx, "123", v1.SendNotificationRequest{Message: "test"})
				if err == nil {
					require.True(t, res.IsOk)
				}
				return err
			},
		},
		{
			name: "history list", method: http.MethodGet,
			path: "/commonserviceitem/simplenotification/history", status: http.StatusOK,
			response: `{"NotificationHistories":[]}`,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				_, err := simplenotification.NewHistoryOp(client).List(ctx)
				return err
			},
		},
		{
			name: "history read", method: http.MethodGet,
			path: "/commonserviceitem/simplenotification/history/123", status: http.StatusNotFound,
			response: `{"is_fatal":false,"serial":"1","status":"404 Not Found","error_code":"not_found","error_msg":"not found"}`,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				_, err := simplenotification.NewHistoryOp(client).Read(ctx, "123")
				require.Error(t, err)
				var apiErr *v1.ErrorStatusCode
				require.ErrorAs(t, err, &apiErr)
				require.Equal(t, http.StatusNotFound, apiErr.StatusCode)
				return nil
			},
		},
		{
			name: "routing reorder", method: http.MethodPut,
			path: "/commonserviceitem/simplenotification/routing/reorder", status: http.StatusAccepted,
			body: `{"Orders":[{"RoutingID":"123","PriorityRank":1}]}`, response: `{"is_ok":true}`,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				_, err := simplenotification.NewRoutingOp(client).Reorder(ctx, v1.ReorderNotificationRoutingRequest{
					Orders: []v1.ReorderNotificationRoutingRequestOrdersItem{{RoutingID: "123", PriorityRank: 1}},
				})
				return err
			},
		},
		{
			name: "routing sources", method: http.MethodGet,
			path: "/commonserviceitem/simplenotification/sources", status: http.StatusOK,
			response: `{"Sources":[]}`,
			invoke: func(t *testing.T, ctx context.Context, client *v1.Client) error {
				_, err := simplenotification.NewRoutingOp(client).ListSource(ctx)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				called = true
				assert.Equal(t, tt.method, r.Method)
				assert.Equal(t, tt.path, r.URL.Path)
				if tt.body != "" {
					var body json.RawMessage
					assert.NoError(t, json.NewDecoder(r.Body).Decode(&body))
					assert.JSONEq(t, tt.body, string(body))
				}
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(tt.response))
			}))
			defer server.Close()
			client, err := v1.NewClient(server.URL, v1.WithClient(server.Client()))
			require.NoError(t, err)
			require.NoError(t, tt.invoke(t, t.Context(), client))
			require.True(t, called)
		})
	}
}
