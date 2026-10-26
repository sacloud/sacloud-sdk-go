// Copyright 2026- The sacloud/sacloud-sdk-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package saclient_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/sacloud/sacloud-sdk-go/common/saclient"
	old "github.com/sacloud/sacloud-sdk-go/internal/api-client"
	"github.com/stretchr/testify/require"
)

// TestTracerShouldNotPanicOnTimeout reproduces
// https://github.com/sacloud/sacloud-sdk-go/issues/304
//
// When the "error" trace mode is set and the HTTP request times out,
// the underlying doer returns (nil, error).  The tracer middleware
// must return the error instead of panicking on a nil pointer
// dereference at res.StatusCode.
func TestTracerShouldNotPanicOnTimeout(t *testing.T) {
	// a server that never responds
	hanging := make(chan struct{})
	hangingServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-hanging
	}))
	defer hangingServer.Close()
	defer close(hanging)

	client := new(Client)
	require.NoError(t, client.CompatSettingsFromAPIClientParams(
		hangingServer.URL,
		old.WithDisableEnv(true),
		old.WithDisableProfile(true),
	))
	require.NoError(t, client.SetWith(WithTraceMode("error"), WithoutRetry(), WithDefaultTimeout(1*time.Millisecond)))

	req, err := http.NewRequest(http.MethodGet, hangingServer.URL+"/", nil)
	require.NoError(t, err)

	resp, err := client.Do(req)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close()
	}
	// The expected behavior: the error is returned, not a panic.
	require.Nil(t, resp)
	require.Error(t, err)
}
