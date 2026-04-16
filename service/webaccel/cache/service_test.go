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

package cache

import (
	"os"
	"strings"
	"testing"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/webaccel-api-go"
	"github.com/stretchr/testify/require"
)

var caller = &webaccel.Client{Options: &client.Options{UserAgent: "webaccel-service-go/v" + webaccel.Version}}

// TestService_CRUD_plus_L CRUD+L
func TestService_CRUD_plus_L(t *testing.T) {
	if !testutil.IsAccTest() {
		t.Skip("environment variables required: TESTACC")
	}
	testutil.PreCheckEnvsFunc(
		"SAKURACLOUD_ACCESS_TOKEN",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET",
		"SAKURACLOUD_WEBACCEL_DOMAIN",
		"SAKURACLOUD_WEBACCEL_URLS",
	)(t)
	svc := New(caller)

	t.Run("Purge with invalid parameter returns error", func(t *testing.T) {
		err := svc.Purge(&PurgeRequest{})
		require.Error(t, err)
	})

	t.Run("Purge by URLs", func(t *testing.T) {
		urls := strings.Split(os.Getenv("SAKURACLOUD_WEBACCEL_URLS"), ",")
		err := svc.Purge(&PurgeRequest{URL: urls})
		require.NoError(t, err)
	})

	t.Run("Purge by Domain", func(t *testing.T) {
		domain := os.Getenv("SAKURACLOUD_WEBACCEL_DOMAIN")
		err := svc.Purge(&PurgeRequest{Domain: domain})
		require.NoError(t, err)
	})
}
