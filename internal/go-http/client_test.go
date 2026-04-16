// Copyright 2021-2025 The sacloud/go-http authors
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

package http

import (
	"bytes"
	"compress/gzip"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type dummyHandler struct {
	called       []time.Time
	responseCode int
}

func (s *dummyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.responseCode == http.StatusMovedPermanently {
		w.Header().Set("Location", "/index.html")
	}
	w.WriteHeader(s.responseCode)
	switch s.responseCode {
	case http.StatusMultipleChoices,
		http.StatusMovedPermanently,
		http.StatusFound,
		http.StatusSeeOther,
		http.StatusNotModified,
		http.StatusUseProxy,
		http.StatusTemporaryRedirect,
		http.StatusPermanentRedirect:
		s.responseCode = http.StatusOK
	default:
		s.called = append(s.called, time.Now())
	}
}

func (s *dummyHandler) isRetried() bool {
	return len(s.called) > 1
}

func TestClient_Do_CheckRetryWithContext(t *testing.T) {
	client := &Client{RetryMax: 1, RetryWaitMin: 10 * time.Millisecond, RetryWaitMax: 10 * time.Millisecond}

	t.Run("context.Canceled", func(t *testing.T) {
		h := &dummyHandler{
			responseCode: http.StatusServiceUnavailable,
		}
		dummyServer := httptest.NewServer(h)
		defer dummyServer.Close()

		ctx, cancel := context.WithCancel(context.Background())
		// make ctx to Canceled
		cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dummyServer.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		client.Do(req) //nolint
		require.False(t, h.isRetried(), "don't retry when context was canceled")
	})

	t.Run("context.DeadlineExceeded", func(t *testing.T) {
		h := &dummyHandler{
			responseCode: http.StatusServiceUnavailable,
		}
		dummyServer := httptest.NewServer(h)
		defer dummyServer.Close()

		ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
		defer cancel()
		// make ctx to DeadlineExceeded
		time.Sleep(time.Millisecond)

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, dummyServer.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		client.Do(req) //nolint
		require.False(t, h.isRetried(), "don't retry when context exceeded deadline")
	})
}

func TestClient_Do_withGzip(t *testing.T) {
	var buf bytes.Buffer
	writer := gzip.NewWriter(&buf)
	io.WriteString(writer, "ok") //nolint //エラーは無視
	writer.Close()
	gzipped := buf.Bytes()

	dummyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusOK)
		w.Write(gzipped) //nolint
	}))
	defer dummyServer.Close()

	t.Run("with Gzip=true", func(t *testing.T) {
		client := &Client{
			Gzip:       true,
			HTTPClient: &http.Client{Transport: &http.Transport{DisableCompression: true}},
		}

		req, err := http.NewRequest(http.MethodGet, dummyServer.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		read, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		require.Equal(t, "ok", string(read))
	})

	t.Run("with Gzip=false", func(t *testing.T) {
		client := &Client{
			Gzip:       false,
			HTTPClient: &http.Client{Transport: &http.Transport{DisableCompression: true}},
		}

		req, err := http.NewRequest(http.MethodGet, dummyServer.URL, nil)
		if err != nil {
			t.Fatal(err)
		}

		res, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}

		read, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}

		// Client.Gzipが: true、かつhttp.TransportのDisableCompression: trueの場合はgzipされたデータがそのまま読めるはず
		require.EqualValues(t, gzipped, read)
	})
}

func TestClient_Do_RequestCustomizer(t *testing.T) {
	cases := []struct {
		name        string
		client      *Client
		queryString string
	}{
		{
			name:        "without request customizer",
			client:      &Client{},
			queryString: "",
		},
		{
			name: "with request customizer",
			client: &Client{
				RequestCustomizer: func(r *http.Request) error {
					r.URL.RawQuery = "foo=bar"
					return nil
				},
			},
			queryString: "foo=bar",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			queryString := ""
			dummyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
				queryString = req.URL.Query().Encode()
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("ok")) //nolint
			}))
			defer dummyServer.Close()

			req, err := http.NewRequest(http.MethodGet, dummyServer.URL, nil)
			if err != nil {
				t.Fatal(err)
			}

			_, err = tt.client.Do(req)
			if err != nil {
				t.Fatal(err)
			}
			require.Equal(t, tt.queryString, queryString)
		})
	}
}
