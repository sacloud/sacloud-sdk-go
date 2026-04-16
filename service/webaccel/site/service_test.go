// Copyright 2022-2025 The sacloud/webaccel-service-go authors
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

package site

import (
	"os"
	"testing"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/webaccel-api-go"
	"github.com/stretchr/testify/require"
)

var caller = &webaccel.Client{Options: &client.Options{UserAgent: "webaccel-service-go/v" + webaccel.Version}}

// TestService_CRUD_plus_L CRUD+L
//
// Note: 実行時にサイトが1件以上登録済みであること
func TestService_CRUD_plus_L(t *testing.T) {
	if !testutil.IsAccTest() {
		t.Skip("environment variables required: TESTACC")
	}
	testutil.PreCheckEnvsFunc(
		"SAKURACLOUD_ACCESS_TOKEN",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET",
		"SAKURACLOUD_WEBACCEL_SITE_ID",
	)(t)
	svc := New(caller)
	siteId := os.Getenv("SAKURACLOUD_WEBACCEL_SITE_ID")

	t.Run("List", func(t *testing.T) {
		found, err := svc.Find(&FindRequest{})
		require.NoError(t, err)
		require.NotEmpty(t, found)
	})

	t.Run("Read", func(t *testing.T) {
		read, err := svc.Read(&ReadRequest{Id: siteId})
		require.NoError(t, err)
		require.NotEmpty(t, read)
		require.Equal(t, siteId, read.ID)
	})

	t.Run("Update", func(t *testing.T) {
		newName := testutil.RandomName("webaccel-service-go-test-", 8, testutil.CharSetAlpha)
		updated, err := svc.Update(&UpdateRequest{
			Id: siteId,
			UpdateSiteRequest: webaccel.UpdateSiteRequest{
				Name: newName,
			},
		})
		require.NoError(t, err)
		require.NotEmpty(t, updated)
		require.Equal(t, newName, updated.Name)
	})
}
