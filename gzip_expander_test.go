// Copyright 2022-2025 The sacloud/saclient-go Authors
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

package saclient

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

var testHttpClient = http.Client{
	Transport: &http.Transport{
		DisableCompression: true,
	},
}

// #nosec G704 -- This is only a test
func terminalMiddleware() Middleware {
	return func(req *http.Request, _ func() (Middleware, bool)) (*http.Response, error) {
		return testHttpClient.Do(req)
	}
}

func singlePull(next Middleware) func() (Middleware, bool) {
	called := false
	return func() (Middleware, bool) {
		if called {
			return nil, false
		}
		called = true
		return next, true
	}
}

func gzipBytes(t *testing.T, s string) []byte {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	if _, err := gw.Write([]byte(s)); err != nil {
		t.Fatalf("failed to write gzip: %v", err)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("failed to close gzip writer: %v", err)
	}
	return buf.Bytes()
}

func TestGzipExpander(t *testing.T) {
	type serverFn func() *httptest.Server

	tests := []struct {
		name              string
		makeServer        serverFn
		usePull           func() (Middleware, bool)
		wantBody          string
		wantEncoding      string
		wantCLHeaderEmpty bool
		wantContentLength int64
		wantErr           bool
	}{
		{
			name: "No gzip: passthrough",
			makeServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, "plain-body")
				}))
			},
			usePull:           singlePull(terminalMiddleware()),
			wantBody:          "plain-body",
			wantEncoding:      "",
			wantCLHeaderEmpty: false,
			wantContentLength: 10,
			wantErr:           false,
		},
		{
			name: "Gzip: decompressed and strip headers",
			makeServer: func() *httptest.Server {
				original := "Merry Christmas, efficiently delivered"
				gz := gzipBytes(t, original)
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Type", "text/plain")
					w.WriteHeader(http.StatusOK)
					_, _ = w.Write(gz)
				}))
			},
			usePull:           singlePull(terminalMiddleware()),
			wantBody:          "Merry Christmas, efficiently delivered",
			wantEncoding:      "",
			wantCLHeaderEmpty: true,
			wantContentLength: -1,
			wantErr:           false,
		},
		{
			name: "Empty gzip (io.EOF) with ContentLength=0: return as-is, no error",
			makeServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Encoding", "gzip")
					w.Header().Set("Content-Type", "text/plain")
					w.Header().Set("Content-Length", "0")
					w.WriteHeader(http.StatusOK)

					// no body
				}))
			},
			usePull:           singlePull(terminalMiddleware()),
			wantBody:          "",
			wantEncoding:      "gzip",
			wantCLHeaderEmpty: false,
			wantContentLength: 0,
			wantErr:           false,
		},
		{
			name: "Invalid gzip: returns error",
			makeServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Header().Set("Content-Type", "application/octet-stream")
					w.Header().Set("Content-Encoding", "gzip")
					w.WriteHeader(http.StatusOK)
					_, _ = io.WriteString(w, "not-a-valid-gzip-stream")
				}))
			},
			usePull: singlePull(terminalMiddleware()),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// setup
			testServer := tt.makeServer()
			defer testServer.Close()

			req, err := http.NewRequest(http.MethodGet, testServer.URL, nil)
			if err != nil {
				t.Fatalf("NewRequest error: %v", err)
			}

			resp, err := gzipExpander(req, tt.usePull)
			if tt.wantErr {
				if resp != nil && resp.Body != nil {
					_ = resp.Body.Close()
				}
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("gzipExpander returned error: %v", err)
			}
			defer func() { _ = resp.Body.Close() }()

			// tests
			if enc := resp.Header.Get("Content-Encoding"); enc != tt.wantEncoding {
				t.Fatalf("Content-Encoding mismatch: got %q, want %q", enc, tt.wantEncoding)
			}

			if tt.wantCLHeaderEmpty {
				if cl := resp.Header.Get("Content-Length"); cl != "" {
					t.Fatalf("expected Content-Length header removed, got %q", cl)
				}
			}

			if resp.ContentLength != tt.wantContentLength {
				t.Fatalf("ContentLength mismatch: got %d, want %d", resp.ContentLength, tt.wantContentLength)
			}

			bodyBytes, readErr := io.ReadAll(resp.Body)
			if readErr != nil {
				t.Fatalf("read body error: %v", readErr)
			}
			if string(bodyBytes) != tt.wantBody {
				t.Fatalf("body mismatch: got %q, want %q", string(bodyBytes), tt.wantBody)
			}
		})
	}
}
