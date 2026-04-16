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

package acl

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

	t.Run("Create", func(t *testing.T) {
		created, err := svc.Create(&CreateRequest{
			SiteId: siteId,
			ACL:    "deny 192.0.2.5/25\ndeny 198.51.100.0\nallow all",
		})

		require.NoError(t, err)
		require.Equal(t, "deny 192.0.2.5/25\ndeny 198.51.100.0\nallow all", created.ACL)
	})

	t.Run("Read", func(t *testing.T) {
		read, err := svc.Read(&ReadRequest{SiteId: siteId})

		require.NoError(t, err)
		require.Equal(t, "deny 192.0.2.5/25\ndeny 198.51.100.0\nallow all", read.ACL)
	})

	t.Run("Update", func(t *testing.T) {
		updated, err := svc.Update(&UpdateRequest{
			SiteId: siteId,
			ACL:    "allow 192.0.2.5/25\nallow 198.51.100.0\ndeny all",
		})

		require.NoError(t, err)
		require.Equal(t, "allow 192.0.2.5/25\nallow 198.51.100.0\ndeny all", updated.ACL)
	})

	t.Run("Delete", func(t *testing.T) {
		err := svc.Delete(&DeleteRequest{SiteId: siteId})
		require.NoError(t, err)

		_, err = svc.Read(&ReadRequest{SiteId: siteId})
		require.Error(t, err) // 404 error
	})
}
