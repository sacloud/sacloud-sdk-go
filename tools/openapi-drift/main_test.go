// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunSeparatesVersionAndStructuralDifferences(t *testing.T) {
	tests := []struct {
		name       string
		version    string
		paths      string
		wantCounts string
		wantRow    string
	}{
		{
			name:       "identical",
			version:    "1.0.0",
			paths:      "{}",
			wantCounts: "version differences: 0; APIs with structural differences: 0.",
			wantRow:    "| Test API | 1.0.0 | 1.0.0 | :white_check_mark: | 0 |",
		},
		{
			name:       "version only",
			version:    "1.1.0",
			paths:      "{}",
			wantCounts: "version differences: 1; APIs with structural differences: 0.",
			wantRow:    "| Test API | 1.0.0 | 1.1.0 | :x: | 0 |",
		},
		{
			name:       "structure only",
			version:    "1.0.0",
			paths:      `{"/new": {}}`,
			wantCounts: "version differences: 0; APIs with structural differences: 1.",
			wantRow:    "| Test API | 1.0.0 | 1.0.0 | :x: | 1 |",
		},
		{
			name:       "version and structure",
			version:    "1.1.0",
			paths:      `{"/new": {}}`,
			wantCounts: "version differences: 1; APIs with structural differences: 1.",
			wantRow:    "| Test API | 1.0.0 | 1.1.0 | :x: | 1 |",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			t.Chdir(dir)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(`{"openapi":"3.0.0","info":{"version":"` + tt.version + `"},"paths":` + tt.paths + `}`))
			}))
			defer server.Close()

			config := manifest{
				Sources:  []source{{ID: "test", URL: server.URL}},
				Packages: []apiSpec{{Name: "Test API", Local: "local.json", Source: "test"}},
			}
			encoded, err := json.Marshal(config)
			if err != nil {
				t.Fatal(err)
			}
			if err := os.MkdirAll(filepath.Dir(manifestPath), 0o755); err != nil {
				t.Fatal(err)
			}
			for filename, content := range map[string][]byte{
				manifestPath: encoded,
				"local.json": []byte(`{"openapi":"3.0.0","info":{"version":"1.0.0"},"paths":{}}`),
				"summary.md": nil,
			} {
				if err := os.WriteFile(filename, content, 0o644); err != nil {
					t.Fatal(err)
				}
			}
			t.Setenv("GITHUB_STEP_SUMMARY", filepath.Join(dir, "summary.md"))

			if err := run(); err != nil {
				t.Fatalf("run() error = %v; want success even when differences are found", err)
			}

			summary, err := os.ReadFile("summary.md")
			if err != nil {
				t.Fatal(err)
			}
			for _, want := range []string{tt.wantCounts, tt.wantRow} {
				if !strings.Contains(string(summary), want) {
					t.Errorf("summary missing %q:\n%s", want, summary)
				}
			}
		})
	}
}

func TestRunFailsWhenManifestCannotBeRead(t *testing.T) {
	t.Chdir(t.TempDir())

	if err := run(); err == nil || !strings.Contains(err.Error(), "read manifest:") {
		t.Fatalf("run() error = %v; want manifest read error", err)
	}
}
