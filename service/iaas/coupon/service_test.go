// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package coupon

import (
	"os"
	"testing"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/helper/api"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/saclient-go"
)

func testCaller() iaas.APICaller {
	return api.NewCallerWithOptions(&api.CallerOptions{
		Options: &client.Options{
			AccessToken:       os.Getenv("SAKURACLOUD_ACCESS_TOKEN"),
			AccessTokenSecret: os.Getenv("SAKURACLOUD_ACCESS_TOKEN_SECRET"),
			UserAgent:         "test-" + iaas.DefaultUserAgent,
			RetryMax:          20,
			Trace:             testutil.IsEnableTrace() || testutil.IsEnableHTTPTrace(),
		},
		TraceAPI: testutil.IsEnableTrace() || testutil.IsEnableHTTPTrace(),
	})
}

func testSaclient() saclient.ClientAPI {
	return &saclient.Client{}
}

func TestService_List(t *testing.T) {
	if !testutil.IsAccTest() {
		t.SkipNow()
	}

	svc := New(testCaller(), testSaclient())

	t.Run("List coupons", func(t *testing.T) {
		coupons, err := svc.List()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("Got %d coupons", len(coupons))
	})
}
