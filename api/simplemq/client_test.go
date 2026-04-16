// Copyright 2025- The sacloud/simplemq-api-go authors
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

package simplemq_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/sacloud/saclient-go"
	"github.com/sacloud/simplemq-api-go"
	"github.com/stretchr/testify/require"
)

func TestNewQueueClient(t *testing.T) {
	assert := require.New(t)

	var theClient saclient.Client
	client, err := simplemq.NewQueueClient(&theClient)
	assert.NoError(err)
	assert.NotNil(client)
}

func TestNewQueueClient_WithCustomEndpoint(t *testing.T) {
	assert := require.New(t)

	tracker := newMockRequestTracker()
	defer tracker.Close()

	var theClient saclient.Client
	err := theClient.SetEnviron([]string{"SAKURA_ENDPOINTS_SIMPLE_MQ_QUEUE=" + tracker.URL()})
	assert.NoError(err)

	client, err := simplemq.NewQueueClient(&theClient)
	assert.NoError(err)
	assert.NotNil(client)

	queueAPI := simplemq.NewQueueOp(client)
	_, _ = queueAPI.List(t.Context())

	requests := tracker.Requests()
	assert.Len(requests, 1)
}

// TestNewQueueClient_FilterQueryParam OpenAPIで表現できず生成コードへのpatchで対応しているURLクエリパラメータが正しく付与されているか
func TestNewQueueClient_FilterQueryParam(t *testing.T) {
	assert := require.New(t)

	tracker := newMockRequestTracker()
	defer tracker.Close()

	var theClient saclient.Client
	err := theClient.SetEnviron([]string{"SAKURA_ENDPOINTS_SIMPLE_MQ_QUEUE=" + tracker.URL()})
	assert.NoError(err)

	client, err := simplemq.NewQueueClient(&theClient)
	assert.NoError(err)
	assert.NotNil(client)

	queueAPI := simplemq.NewQueueOp(client)
	_, _ = queueAPI.List(t.Context())

	requests := tracker.Requests()
	assert.Len(requests, 1)
	assert.Equal("/commonserviceitem", requests[0].URL.Path)
	expectedQuery := url.QueryEscape(`{"Filter":{"Provider.Class":"simplemq"}}`)
	assert.Equal(expectedQuery, requests[0].URL.RawQuery)
}

func TestNewMessageClient(t *testing.T) {
	assert := require.New(t)

	var theClient saclient.Client
	client, err := simplemq.NewMessageClient("test-api-key", &theClient)
	assert.NoError(err)
	assert.NotNil(client)
}

func TestNewMessageClient_WithCustomEndpoint(t *testing.T) {
	assert := require.New(t)

	tracker := newMockRequestTracker()
	defer tracker.Close()

	var theClient saclient.Client
	err := theClient.SetEnviron([]string{"SAKURA_ENDPOINTS_SIMPLE_MQ_MESSAGE=" + tracker.URL()})
	assert.NoError(err)

	client, err := simplemq.NewMessageClient("test-api-key", &theClient)
	assert.NoError(err)
	assert.NotNil(client)

	messageAPI := simplemq.NewMessageOp(client, "test-queue")
	_, _ = messageAPI.Send(t.Context(), "test")

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

		switch {
		case r.URL.Path == "/commonserviceitem":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck,gosec
				"CommonServiceItems": []interface{}{},
			})
		case r.URL.Path == "/v1/queues/test-queue/messages" && r.Method == "POST":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck,gosec
				"result": "success",
				"message": map[string]interface{}{
					"id":         "0193b878-1b25-7775-87f5-9c698206a7e7",
					"content":    "test message",
					"created_at": 1704067200000,
					"updated_at": 1704067200000,
					"expires_at": 1704070800000,
				},
			})
		default:
			w.WriteHeader(http.StatusNoContent)
		}
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
