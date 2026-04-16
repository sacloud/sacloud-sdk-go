// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	. "github.com/sacloud/nosql-api-go"
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
	err := theClient.SetEnviron([]string{"SAKURA_ENDPOINTS_NOSQL=" + tracker.URL()})
	assert.NoError(err)

	client, err := NewClient(&theClient)
	assert.NoError(err)
	assert.NotNil(client)

	op := NewDatabaseOp(client)
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
