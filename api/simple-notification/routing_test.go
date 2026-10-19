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

	simplenotification "github.com/sacloud/sacloud-sdk-go/api/simple-notification"
	v1 "github.com/sacloud/sacloud-sdk-go/api/simple-notification/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/require"
)

func routingAPISetup(t *testing.T) (ctx context.Context, api simplenotification.RoutingAPI) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	ctx = t.Context()

	var saClient saclient.Client

	client, err := simplenotification.NewClient(&saClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	api = simplenotification.NewRoutingOp(client)

	return ctx, api
}

func TestRoutingOp(t *testing.T) {
	ctx, routingAPI := routingAPISetup(t)
	testutil.PreCheckEnvsFunc("SAKURA_ROUTING_TEST_LABEL", "SAKURA_ROUTING_TEST_LABELVAL", "SAKURA_ROUTING_TEST_SOURCEID", "SAKURA_ROUTING_TEST_TARGETGROUPID")(t)
	id := "" // set created routing ID
	routingName := "test-routing-1"
	description := "test-routing-description"
	tags := []string{"test"}
	testLabel := os.Getenv("SAKURA_ROUTING_TEST_LABEL")
	testLabelValue := os.Getenv("SAKURA_ROUTING_TEST_LABELVAL")
	testSourceID := os.Getenv("SAKURA_ROUTING_TEST_SOURCEID")
	testTargetGroupID := os.Getenv("SAKURA_ROUTING_TEST_TARGETGROUPID")
	setting := v1.CommonServiceItemRoutingSettings{
		MatchLabels: []v1.CommonServiceItemRoutingSettingsMatchLabelsItem{
			{Name: testLabel, Value: testLabelValue},
		},
		SourceID:      testSourceID,
		TargetGroupID: testTargetGroupID,
		// Although it is necessary to set PriorityRank,
		// PriorityRank is determined from the routing information set on the service side.
		// If you set PriorityRank, use Reorder API to change the order after creation.
		PriorityRank: 1,
	}
	result := t.Run("Create", func(t *testing.T) {
		request := v1.PostCommonServiceItemRequest{
			CommonServiceItem: v1.PostCommonServiceItemRequestCommonServiceItem{
				Name:        routingName,
				Description: description,
				Tags:        tags,
				Icon: v1.NilCommonServiceItemIcon{
					Null: false,
					Value: v1.CommonServiceItemIcon{
						ID: "112901627732", // debian icon ID
					},
				},
				Provider: v1.PostCommonServiceItemRequestCommonServiceItemProvider{
					Class: v1.PostCommonServiceItemRequestCommonServiceItemProviderClassSaknoticerouting,
				},
				Settings: v1.PostCommonServiceItemRequestCommonServiceItemSettings{
					Type:                             v1.CommonServiceItemRoutingSettingsPostCommonServiceItemRequestCommonServiceItemSettings,
					CommonServiceItemRoutingSettings: setting,
				},
			},
		}
		resp, err := routingAPI.Create(ctx, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("RoutingOp.Create response: %+v", resp)
		id = resp.CommonServiceItem.ID
	})
	defer func() {
		if id != "" {
			err := routingAPI.Delete(ctx, id)
			if err != nil {
				t.Fatalf("unexpected error on cleanup: %v", err)
			}
			t.Log("RoutingOp.Delete succeeded on cleanup")
		}
	}()

	if !result {
		t.Fatal("skipping rest of tests due to Create failure")
	}

	t.Run("List", func(t *testing.T) {
		resp, err := routingAPI.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("RoutingOp.List response: %+v", resp)
	})
	t.Run("Read", func(t *testing.T) {
		resp, err := routingAPI.Read(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("RoutingOp.Read response: %+v", resp)
	})
	t.Run("Update", func(t *testing.T) {
		routingNameUpdate := "updated-routing"
		descriptionUpdate := "updated-description"
		tagsUpdate := []string{"updated"}
		updateSetting := v1.CommonServiceItemRoutingSettings{
			MatchLabels: []v1.CommonServiceItemRoutingSettingsMatchLabelsItem{
				{Name: testLabel, Value: "sdk-test-updated"},
			},
			SourceID:      testSourceID,
			TargetGroupID: testTargetGroupID,
			// Although it is necessary to set PriorityRank,
			// PriorityRank can't be updated, so not set here (if want to update, use Reorder API)
			PriorityRank: 1,
		}
		request := v1.PutCommonServiceItemRequest{
			CommonServiceItem: v1.PutCommonServiceItemRequestCommonServiceItem{
				Name:        routingNameUpdate,
				Description: descriptionUpdate,
				Tags:        tagsUpdate,
				Icon: v1.NilCommonServiceItemIcon{
					Null: true,
				},
				Settings: v1.OptPutCommonServiceItemRequestCommonServiceItemSettings{
					Set: true,
					Value: v1.PutCommonServiceItemRequestCommonServiceItemSettings{
						Type:                             v1.CommonServiceItemRoutingSettingsPutCommonServiceItemRequestCommonServiceItemSettings,
						CommonServiceItemRoutingSettings: updateSetting,
					},
				},
			},
		}
		resp, err := routingAPI.Update(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		require.Equal(t, routingNameUpdate, resp.CommonServiceItem.Name, "routing name should be updated")
		require.Equal(t, descriptionUpdate, resp.CommonServiceItem.Description, "description should be updated")
		require.Equal(t, tagsUpdate, resp.CommonServiceItem.Tags, "tags should be updated")
		require.True(t, resp.CommonServiceItem.Settings.IsCommonServiceItemRoutingSettings())
		actualSetting := resp.CommonServiceItem.Settings.CommonServiceItemRoutingSettings
		require.Equal(t, updateSetting.MatchLabels, actualSetting.MatchLabels)
		require.Equal(t, updateSetting.SourceID, actualSetting.SourceID)
		require.Equal(t, updateSetting.TargetGroupID, actualSetting.TargetGroupID)
		read, err := routingAPI.Read(ctx, id)
		require.NoError(t, err)
		require.Equal(t, resp.CommonServiceItem.Settings, read.CommonServiceItem.Settings)
		t.Logf("RoutingOp.Update response: %+v", resp)
	})
	t.Run("UpdateWithoutSetting", func(t *testing.T) {
		routingNameUpdateWithoutSetting := "updated-routing-withoutSetting"
		descriptionUpdateWithoutSetting := "updated-description-withoutSetting"
		tagsUpdateWithoutSetting := []string{"updated-withoutSetting"}
		before, err := routingAPI.Read(ctx, id)
		require.NoError(t, err)
		request := v1.PutCommonServiceItemRequest{
			CommonServiceItem: v1.PutCommonServiceItemRequestCommonServiceItem{
				Name:        routingNameUpdateWithoutSetting,
				Description: descriptionUpdateWithoutSetting,
				Tags:        tagsUpdateWithoutSetting,
				Icon: v1.NilCommonServiceItemIcon{
					Null: false,
					Value: v1.CommonServiceItemIcon{
						ID: "0", // erase icon
					},
				},
				Settings: v1.OptPutCommonServiceItemRequestCommonServiceItemSettings{
					Set: false,
				},
			},
		}
		resp, err := routingAPI.Update(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		require.Equal(t, routingNameUpdateWithoutSetting, resp.CommonServiceItem.Name, "routing name should be updated")
		require.Equal(t, descriptionUpdateWithoutSetting, resp.CommonServiceItem.Description, "description should be updated")
		require.Equal(t, tagsUpdateWithoutSetting, resp.CommonServiceItem.Tags, "tags should be updated")
		require.Equal(t, before.CommonServiceItem.Settings, resp.CommonServiceItem.Settings)
		read, err := routingAPI.Read(ctx, id)
		require.NoError(t, err)
		require.Equal(t, before.CommonServiceItem.Settings, read.CommonServiceItem.Settings)
		t.Logf("RoutingOp.UpdateWithoutSetting response: %+v", resp)
	})
	t.Run("ClearTags", func(t *testing.T) {
		before, err := routingAPI.Read(ctx, id)
		require.NoError(t, err)
		require.NotEmpty(t, before.CommonServiceItem.Tags)
		resp, err := routingAPI.Update(ctx, id, v1.PutCommonServiceItemRequest{
			CommonServiceItem: v1.PutCommonServiceItemRequestCommonServiceItem{
				Name: before.CommonServiceItem.Name, Description: before.CommonServiceItem.Description,
				Tags: []string{}, Icon: before.CommonServiceItem.Icon,
			},
		})
		require.NoError(t, err)
		require.Empty(t, resp.CommonServiceItem.Tags)
		require.Equal(t, before.CommonServiceItem.Settings, resp.CommonServiceItem.Settings)
		read, err := routingAPI.Read(ctx, id)
		require.NoError(t, err)
		require.Empty(t, read.CommonServiceItem.Tags)
		require.Equal(t, before.CommonServiceItem.Settings, read.CommonServiceItem.Settings)
	})
	t.Run("Reorder", func(t *testing.T) {
		secondCreateRequest := v1.PostCommonServiceItemRequest{
			CommonServiceItem: v1.PostCommonServiceItemRequestCommonServiceItem{
				Name:        "test-routing-2",
				Description: "test-routing-description-2",
				Tags:        tags,
				Provider: v1.PostCommonServiceItemRequestCommonServiceItemProvider{
					Class: v1.PostCommonServiceItemRequestCommonServiceItemProviderClassSaknoticerouting,
				},
				Settings: v1.PostCommonServiceItemRequestCommonServiceItemSettings{
					Type:                             v1.CommonServiceItemRoutingSettingsPostCommonServiceItemRequestCommonServiceItemSettings,
					CommonServiceItemRoutingSettings: setting,
				},
			},
		}
		createResp, err := routingAPI.Create(ctx, secondCreateRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if createResp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("RoutingOp.Create (for Reorder) response: %+v", createResp)
		defer func() {
			err := routingAPI.Delete(ctx, createResp.CommonServiceItem.ID)
			if err != nil {
				t.Fatalf("unexpected error on cleanup: %v", err)
			}
			t.Log("RoutingOp.Delete (for Reorder) succeeded on cleanup")
		}()
		firstPriorityRank := 3
		secondPriorityRank := 4
		reorderRequest := v1.ReorderNotificationRoutingRequest{
			Orders: []v1.ReorderNotificationRoutingRequestOrdersItem{
				{
					RoutingID:    id,
					PriorityRank: firstPriorityRank,
				},
				{
					RoutingID:    createResp.CommonServiceItem.ID,
					PriorityRank: secondPriorityRank,
				},
			},
		}
		resp, err := routingAPI.Reorder(ctx, reorderRequest)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("RoutingOp.Reorder response: %+v", resp)

		listResp, err := routingAPI.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if listResp == nil {
			t.Fatal("expected response but got nil")
		}
		var firstRouting, secondRouting *v1.CommonServiceItem
		for _, item := range listResp.CommonServiceItems {
			if item.ID == id {
				firstRouting = &item
			}
			if item.ID == createResp.CommonServiceItem.ID {
				secondRouting = &item
			}
		}
		if firstRouting == nil || secondRouting == nil {
			t.Fatal("could not find created routings in list response")
		}
		require.Equal(t, firstPriorityRank, firstRouting.Settings.CommonServiceItemRoutingSettings.PriorityRank, "first routing's priority rank should be updated")
		require.Equal(t, secondPriorityRank, secondRouting.Settings.CommonServiceItemRoutingSettings.PriorityRank, "second routing's priority rank should be updated")
	})
	t.Run("ListSource", func(t *testing.T) {
		resp, err := routingAPI.ListSource(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("RoutingOp.ListSource response: %+v", resp)
	})

	t.Run("Delete", func(t *testing.T) {
		err := routingAPI.Delete(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		t.Log("RoutingOp.Delete succeeded")
		id = "" // prevent double delete in defer
	})
}
