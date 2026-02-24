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
	"os"
	"testing"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	simplenotification "github.com/sacloud/simple-notification-api-go"
	v1 "github.com/sacloud/simple-notification-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func groupAPISetup(t *testing.T) (ctx context.Context, api simplenotification.GroupAPI) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_DESTINATION_TEST_ID")(t)

	ctx = t.Context()

	var saClient saclient.Client

	client, err := simplenotification.NewClient(&saClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	api = simplenotification.NewGroupOp(client)

	return ctx, api
}

func TestGroupOp(t *testing.T) {
	ctx, groupAPI := groupAPISetup(t)

	id := "" // set your pre-created group ID here
	groupname := "test-group-1"
	description := "test-group-description"
	tags := []string{"test"}
	destinationID := os.Getenv("SAKURA_DESTINATION_TEST_ID")

	result := t.Run("Create", func(t *testing.T) {
		request := v1.PostCommonServiceItemRequest{
			CommonServiceItem: v1.PostCommonServiceItemRequestCommonServiceItem{
				Name:        groupname,
				Description: description,
				Tags:        tags,
				Settings: v1.CommonServiceItemSettings{
					Type: v1.GroupSettingsCommonServiceItemSettings,
					GroupSettings: v1.GroupSettings{
						Destinations: []string{destinationID},
					},
				},
			},
		}
		resp, err := groupAPI.Create(ctx, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("GroupOp.Create response: %+v", resp)
		id = resp.CommonServiceItem.ID
	})
	defer func() {
		if id != "" {
			err := groupAPI.Delete(ctx, id)
			if err != nil {
				t.Fatalf("unexpected error on cleanup: %v", err)
			}
			t.Log("GroupOp.Delete succeeded on cleanup")
		}
	}()

	if !result {
		t.Fatal("skipping rest of tests due to Create failure")
	}

	t.Run("List", func(t *testing.T) {
		resp, err := groupAPI.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("GroupOp.List response: %+v", resp)
	})
	t.Run("Read", func(t *testing.T) {
		resp, err := groupAPI.Read(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("GroupOp.Read response: %+v", resp)
	})
	t.Run("Update", func(t *testing.T) {
		groupnameUpdate := "updated-group"
		descriptionUpdate := "updated-description"
		tagsUpdate := []string{"updated"}
		request := v1.PutCommonServiceItemRequest{
			CommonServiceItem: v1.PutCommonServiceItemRequestCommonServiceItem{
				Name:        groupnameUpdate,
				Description: descriptionUpdate,
				Tags:        tagsUpdate,
				Settings: v1.OptCommonServiceItemSettings{
					Set: true,
					Value: v1.CommonServiceItemSettings{
						Type: v1.GroupSettingsCommonServiceItemSettings,
						GroupSettings: v1.GroupSettings{
							Destinations: []string{destinationID},
						},
					},
				},
			},
		}
		resp, err := groupAPI.Update(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		require.Equal(t, groupnameUpdate, resp.CommonServiceItem.Name, "group name should be updated")
		require.Equal(t, descriptionUpdate, resp.CommonServiceItem.Description, "description should be updated")
		require.Equal(t, tagsUpdate, resp.CommonServiceItem.Tags, "tags should be updated")
		t.Logf("GroupOp.Update response: %+v", resp)
	})
	t.Run("UpdatewithoutSetting", func(t *testing.T) {
		groupnameUpdatewithoutSetting := "updated-group-withoutSetting"
		descriptionUpdatewithoutSetting := "updated-description-withoutSetting"
		tagsUpdatewithoutSetting := []string{"updated-withoutSetting"}

		request := v1.PutCommonServiceItemRequest{
			CommonServiceItem: v1.PutCommonServiceItemRequestCommonServiceItem{
				Name:        groupnameUpdatewithoutSetting,
				Description: descriptionUpdatewithoutSetting,
				Tags:        tagsUpdatewithoutSetting,
				Settings: v1.OptCommonServiceItemSettings{
					Set: false,
				},
			},
		}
		resp, err := groupAPI.Update(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		require.Equal(t, groupnameUpdatewithoutSetting, resp.CommonServiceItem.Name, "group name should be updated")
		require.Equal(t, descriptionUpdatewithoutSetting, resp.CommonServiceItem.Description, "description should be updated")
		require.Equal(t, tagsUpdatewithoutSetting, resp.CommonServiceItem.Tags, "tags should be updated")
		t.Logf("GroupOp.UpdatewithoutSetting response: %+v", resp)
	})
	t.Run("SendMessage", func(t *testing.T) {
		request := v1.SendNotificationMessageRequest{
			Message: "test message from GroupOp.SendMessage",
		}

		resp, err := groupAPI.SendMessage(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("GroupOp.SendMessage response: %+v", resp)
	})
	t.Run("Delete", func(t *testing.T) {
		err := groupAPI.Delete(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		t.Log("GroupOp.Delete succeeded")
		id = "" // prevent double delete in defer
	})
}
