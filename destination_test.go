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
	"github.com/stretchr/testify/require"

	v1 "github.com/sacloud/simple-notification-api-go/apis/v1"
)

func destinationAPISetup(t *testing.T) (ctx context.Context, api simplenotification.DestinationAPI) {

	ctx = t.Context()
	var saClient saclient.Client

	client, err := simplenotification.NewClient(&saClient)
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	api = simplenotification.NewDestinationOp(client)

	return ctx, api
}

func TestDestinationOp(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_DESTINATION_TEST_MAILADDRESS")(t)
	ctx, destinationAPI := destinationAPISetup(t)
	id := ""
	destName := "test-destination-1"
	description := "test-destination-description"
	tags := []string{"test"}
	mailAddress := os.Getenv("SAKURA_DESTINATION_TEST_MAILADDRESS")
	setting := v1.CommonServiceItemDestinationSettings{
		Type:  v1.CommonServiceItemDestinationSettingsTypeEmail,
		Value: mailAddress,
	}
	result := t.Run("Create", func(t *testing.T) {
		request := v1.PostCommonServiceItemRequest{
			CommonServiceItem: v1.PostCommonServiceItemRequestCommonServiceItem{
				Name:        destName,
				Description: description,
				Tags:        tags,
				Icon: v1.NilCommonServiceItemIcon{
					Null: false,
					Value: v1.CommonServiceItemIcon{
						ID: "112901627732", //Debian icon ID
					},
				},
				Settings: v1.PostCommonServiceItemRequestCommonServiceItemSettings{
					CommonServiceItemDestinationSettings: setting,
				},
			},
		}
		resp, err := destinationAPI.Create(ctx, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("DestinationOp.Create response: %+v", resp)
		id = resp.CommonServiceItem.ID
	})
	defer func() {
		if id != "" {
			err := destinationAPI.Delete(ctx, id)
			if err != nil {
				t.Errorf("unexpected error on cleanup: %v", err)
			}
			t.Log("DestinationOp.Delete succeeded on cleanup")
		}
	}()

	if !result {
		t.Fatal("skipping rest of tests due to Create failure")
	}

	t.Run("List", func(t *testing.T) {
		resp, err := destinationAPI.List(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("DestinationOp.List response: %+v", resp)
	})
	t.Run("Read", func(t *testing.T) {
		resp, err := destinationAPI.Read(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("DestinationOp.Read response: %+v", resp)
	})
	t.Run("Update", func(t *testing.T) {
		destNameUpdate := "updated-destination"
		descriptionUpdate := "updated-description"
		tagsUpdate := []string{"updated"}

		request := v1.PutCommonServiceItemRequest{
			CommonServiceItem: v1.PutCommonServiceItemRequestCommonServiceItem{
				Name:        destNameUpdate,
				Description: descriptionUpdate,
				Tags:        tagsUpdate,
				Settings: v1.OptPutCommonServiceItemRequestCommonServiceItemSettings{
					Set: true,
					Value: v1.PutCommonServiceItemRequestCommonServiceItemSettings{
						Type: v1.CommonServiceItemDestinationSettingsPutCommonServiceItemRequestCommonServiceItemSettings,
						CommonServiceItemDestinationSettings: v1.CommonServiceItemDestinationSettings{
							Type:  v1.CommonServiceItemDestinationSettingsTypeEmail,
							Value: mailAddress,
						},
					},
				},
			},
		}
		resp, err := destinationAPI.Update(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		require.Equal(t, destNameUpdate, resp.CommonServiceItem.Name, "destination name should be updated")
		require.Equal(t, descriptionUpdate, resp.CommonServiceItem.Description, "description should be updated")
		require.Equal(t, tagsUpdate, resp.CommonServiceItem.Tags, "tags should be updated")
		t.Logf("DestinationOp.Update response: %+v", resp)
	})
	t.Run("UpdatewithoutSetting", func(t *testing.T) {
		destNameUpdateWithoutSetting := "updated-destination-withoutSetting"
		descriptionUpdateWithoutSetting := "updated-description-withoutSetting"
		tagsUpdateWithoutSetting := []string{"updated-withoutSetting"}

		request := v1.PutCommonServiceItemRequest{
			CommonServiceItem: v1.PutCommonServiceItemRequestCommonServiceItem{
				Name:        destNameUpdateWithoutSetting,
				Description: descriptionUpdateWithoutSetting,
				Tags:        tagsUpdateWithoutSetting,
				Settings: v1.OptPutCommonServiceItemRequestCommonServiceItemSettings{
					Set: false,
				},
			},
		}
		resp, err := destinationAPI.Update(ctx, id, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		require.Equal(t, destNameUpdateWithoutSetting, resp.CommonServiceItem.Name, "destination name should be updated")
		require.Equal(t, descriptionUpdateWithoutSetting, resp.CommonServiceItem.Description, "description should be updated")
		require.Equal(t, tagsUpdateWithoutSetting, resp.CommonServiceItem.Tags, "tags should be updated")
		t.Logf("DestinationOp.UpdatewithoutSetting response: %+v", resp)
	})
	t.Run("ReadStatus", func(t *testing.T) {
		resp, err := destinationAPI.ReadStatus(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if resp == nil {
			t.Fatal("expected response but got nil")
		}
		t.Logf("DestinationOp.GetStatus response: %+v", resp)
	})
	t.Run("Delete", func(t *testing.T) {
		err := destinationAPI.Delete(ctx, id)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		t.Log("DestinationOp.Delete succeeded")
		id = "" // prevent double delete in defer
	})
}
