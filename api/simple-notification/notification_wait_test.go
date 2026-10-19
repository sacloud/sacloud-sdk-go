// Copyright 2026- The sacloud/simple-notification-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package simplenotification_test

import (
	"context"
	"errors"
	"testing"
	"time"

	v1 "github.com/sacloud/sacloud-sdk-go/api/simple-notification/apis/v1"
	"github.com/stretchr/testify/require"
)

type notificationHistoryStub struct {
	list func(context.Context) (*v1.ListNotificationHistoriesResponse, error)
}

func (s notificationHistoryStub) List(ctx context.Context) (*v1.ListNotificationHistoriesResponse, error) {
	return s.list(ctx)
}

func (s notificationHistoryStub) Read(context.Context, string) (*v1.ReadNotificationHistoryResponse, error) {
	return nil, errors.New("unexpected history read")
}

func TestWaitForNotification(t *testing.T) {
	status := func(value v1.NotificationStatusStatus) v1.NotificationStatus {
		return v1.NotificationStatus{GroupID: "group", DestinationID: "destination", Status: value}
	}
	item := func(message string, statuses ...v1.NotificationStatus) v1.NotificationHistory {
		return v1.NotificationHistory{Message: v1.NotificationMessage{Body: message}, Statuses: statuses}
	}
	tests := []struct {
		name      string
		responses [][]v1.NotificationHistory
		wantError string
	}{
		{
			name: "wait for matching history and sending to finish",
			responses: [][]v1.NotificationHistory{
				{},
				{item("unrelated", status(v1.NotificationStatusStatus2))},
				{item("message", v1.NotificationStatus{GroupID: "other", DestinationID: "destination", Status: v1.NotificationStatusStatus2})},
				{item("message", v1.NotificationStatus{GroupID: "group", DestinationID: "other", Status: v1.NotificationStatusStatus2})},
				{item("message")},
				{item("message", status(v1.NotificationStatusStatus0))},
				{item("message", status(v1.NotificationStatusStatus1))},
				{item("message", status(v1.NotificationStatusStatus2))},
			},
		},
		{
			name: "all matching statuses must finish",
			responses: [][]v1.NotificationHistory{
				{item("message", status(v1.NotificationStatusStatus2), status(v1.NotificationStatusStatus1))},
				{item("message", status(v1.NotificationStatusStatus2), status(v1.NotificationStatusStatus2))},
			},
		},
		{
			name: "failure still waits for pending statuses",
			responses: [][]v1.NotificationHistory{
				{item("message", status(v1.NotificationStatusStatus9), status(v1.NotificationStatusStatus0))},
				{item("message", status(v1.NotificationStatusStatus9), status(v1.NotificationStatusStatus2))},
			},
			wantError: "status 9",
		},
		{
			name: "skipped is not success",
			responses: [][]v1.NotificationHistory{
				{item("message", status(v1.NotificationStatusStatus3))},
			},
			wantError: "status 3",
		},
		{
			name: "failure includes reason",
			responses: [][]v1.NotificationHistory{
				{item("message", v1.NotificationStatus{
					GroupID: "group", DestinationID: "destination",
					Status: v1.NotificationStatusStatus9, ErrorInfo: "notification group not found",
				})},
			},
			wantError: "notification group not found",
		},
		{
			name: "unknown status is an error",
			responses: [][]v1.NotificationHistory{
				{item("message", status(v1.NotificationStatusStatus(99)))},
			},
			wantError: "unexpected notification status",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			calls := 0
			history := notificationHistoryStub{list: func(context.Context) (*v1.ListNotificationHistoriesResponse, error) {
				require.Less(t, calls, len(tt.responses), "unexpected additional poll")
				res := &v1.ListNotificationHistoriesResponse{NotificationHistories: tt.responses[calls]}
				calls++
				return res, nil
			}}
			ctx, cancel := context.WithTimeout(t.Context(), time.Second)
			defer cancel()
			err := waitForNotification(ctx, history, "group", "destination", "message", time.Millisecond)
			if tt.wantError == "" {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.wantError)
			}
			require.Equal(t, len(tt.responses), calls)
		})
	}
}

func TestWaitForNotificationHistoryError(t *testing.T) {
	want := errors.New("history unavailable")
	history := notificationHistoryStub{list: func(context.Context) (*v1.ListNotificationHistoriesResponse, error) {
		return nil, want
	}}
	err := waitForNotification(t.Context(), history, "group", "destination", "message", time.Millisecond)
	require.ErrorIs(t, err, want)
}

func TestWaitForNotificationTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(t.Context(), time.Millisecond)
	defer cancel()
	history := notificationHistoryStub{list: func(context.Context) (*v1.ListNotificationHistoriesResponse, error) {
		return &v1.ListNotificationHistoriesResponse{}, nil
	}}
	err := waitForNotification(ctx, history, "group", "destination", "message", time.Hour)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
