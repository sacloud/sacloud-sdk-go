// Copyright 2025-2026 The sacloud/eventbus-api-go authors
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

package eventbus_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	. "github.com/sacloud/eventbus-api-go"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	assert := require.New(t)

	var theClient saclient.Client
	client, err := NewClient(&theClient)
	assert.NoError(err)
	assert.NotNil(client)
}

func TestNewClient_WithCustomEndpoint(t *testing.T) {
	assert := require.New(t)

	tracker := newMockRequestTracker()
	defer tracker.Close()

	var theClient saclient.Client
	err := theClient.SetEnviron([]string{"SAKURA_ENDPOINTS_EVENTBUS=" + tracker.URL()})
	assert.NoError(err)

	client, err := NewClient(&theClient)
	assert.NoError(err)
	assert.NotNil(client)

	op := NewProcessConfigurationOp(client)
	_, _ = op.List(t.Context())

	requests := tracker.Requests()
	assert.Len(requests, 1)
}

type mockRequestTracker struct {
	mu       sync.Mutex
	requests []*http.Request
	server   *httptest.Server
}

func (m *mockRequestTracker) handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m.mu.Lock()
		m.requests = append(m.requests, r)
		m.mu.Unlock()

		w.WriteHeader(http.StatusNoContent)
	}
}

func newMockRequestTracker() *mockRequestTracker {
	tracker := &mockRequestTracker{}
	tracker.server = httptest.NewServer(tracker.handler())
	return tracker
}

func (m *mockRequestTracker) Close() {
	if m.server != nil {
		m.server.Close()
	}
}

func (m *mockRequestTracker) URL() string {
	return m.server.URL
}

func (m *mockRequestTracker) Requests() []*http.Request {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.requests
}
