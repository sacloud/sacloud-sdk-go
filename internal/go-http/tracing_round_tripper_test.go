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
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTracingRoundTripper_OutputOnlyError(t *testing.T) {
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "ok")
	})
	server := httptest.NewServer(h)

	type fields struct {
		Transport       http.RoundTripper
		OutputOnlyError bool
	}
	tests := []struct {
		name         string
		fields       fields
		wantTraceLog bool
	}{
		{
			name:         "all",
			fields:       fields{},
			wantTraceLog: true,
		},
		{
			name: "only error",
			fields: fields{
				OutputOnlyError: true,
			},
			wantTraceLog: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &TracingRoundTripper{
				Transport:       tt.fields.Transport,
				OutputOnlyError: tt.fields.OutputOnlyError,
			}
			client := &http.Client{Transport: r}
			req, err := http.NewRequest("POST", server.URL, bytes.NewBufferString("request body"))
			if err != nil {
				t.Fatal(err)
			}

			logs := bytes.NewBufferString("")
			log.SetOutput(logs)
			log.SetFlags(0)

			_, err = client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			require.Equal(t, tt.wantTraceLog, strings.Contains(logs.String(), "POST / HTTP/1.1"), "error log not found: request")
			require.Equal(t, tt.wantTraceLog, strings.Contains(logs.String(), "request body"), "error log not found: request body")
			require.Equal(t, tt.wantTraceLog, strings.Contains(logs.String(), "HTTP/1.1 200 OK"), "error log not found: response")
		})
	}
}
