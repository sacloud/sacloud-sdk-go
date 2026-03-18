// Copyright 2026- The sacloud/simple-notification-api-go Authors
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

package simplenotification_test

import (
	"context"
	"testing"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	simplenotification "github.com/sacloud/simple-notification-api-go"
)

func historyAPISetup(t *testing.T) (ctx context.Context, api simplenotification.HistoryAPI) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	ctx = t.Context()

	var saClient saclient.Client

	client, err := simplenotification.NewClient(&saClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	api = simplenotification.NewHistoryOp(client)

	return ctx, api
}

func TestHistoryOp(t *testing.T) {
	ctx, historyAPI := historyAPISetup(t)
	id := ""
	t.Run("List", func(t *testing.T) {
		resp, err := historyAPI.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("HistoryOp.List response: %+v", resp)
		if len(resp.NotificationHistories) == 0 {
			t.Skip("expected history 0. skipping Read test")
		} else {
			id = resp.NotificationHistories[0].RequestID
		}
	})
	t.Run("Read", func(t *testing.T) {
		if id == "" {
			t.Skip("no history found. skipping")
		}
		resp, err := historyAPI.Read(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("HistoryOp.Read response: %+v", resp)
	})
}
