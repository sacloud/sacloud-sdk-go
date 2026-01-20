// Copyright 2025- The sacloud/saclient-go Authors
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

package saclient

import (
	"strings"
	"testing"
)

// TestNormalizeEndpoints tests the normalizeEndpoints function
func TestNormalizeEndpoints(t *testing.T) {
	tests := []struct {
		name      string
		endpoints map[string]string
		expected  map[string]string
	}{
		{
			name:      "nil input",
			endpoints: nil,
			expected:  nil,
		},
		{
			name:      "empty map",
			endpoints: map[string]string{},
			expected:  map[string]string{},
		},
		{
			name: "all lowercase keys",
			endpoints: map[string]string{
				"iaas": "https://iaas.example.com",
				"iam":  "https://iam.example.com",
			},
			expected: map[string]string{
				"iaas": "https://iaas.example.com",
				"iam":  "https://iam.example.com",
			},
		},
		{
			name: "mixed case keys",
			endpoints: map[string]string{
				"IaaS":          "https://iaas.example.com",
				"IAM":           "https://iam.example.com",
				"ObjectStorage": "https://storage.example.com",
			},
			expected: map[string]string{
				"iaas":          "https://iaas.example.com",
				"iam":           "https://iam.example.com",
				"objectstorage": "https://storage.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeEndpoints(tt.endpoints)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("expected %v, got nil", tt.expected)
				} else {
					for k, expectedValue := range tt.expected {
						if actualValue, ok := result[k]; !ok {
							t.Errorf("expected key %s not found", k)
						} else if actualValue != expectedValue {
							t.Errorf("key %s: expected %s, got %s", k, expectedValue, actualValue)
						}
					}
					// Check that all keys are lowercase
					for k := range result {
						if k != strings.ToLower(k) {
							t.Errorf("key %s is not lowercase", k)
						}
					}
				}
			}
		})
	}
}

// TestNormalizeEndpointsFromAny tests the normalizeEndpointsFromAny function
func TestNormalizeEndpointsFromAny(t *testing.T) {
	tests := []struct {
		name      string
		endpoints map[string]any
		expected  map[string]string
	}{
		{
			name:      "nil input",
			endpoints: nil,
			expected:  nil,
		},
		{
			name:      "empty map",
			endpoints: map[string]any{},
			expected:  map[string]string{},
		},
		{
			name: "all string values with lowercase keys",
			endpoints: map[string]any{
				"iaas": "https://iaas.example.com",
				"iam":  "https://iam.example.com",
			},
			expected: map[string]string{
				"iaas": "https://iaas.example.com",
				"iam":  "https://iam.example.com",
			},
		},
		{
			name: "mixed case keys with string values",
			endpoints: map[string]any{
				"IaaS":          "https://iaas.example.com",
				"IAM":           "https://iam.example.com",
				"ObjectStorage": "https://storage.example.com",
			},
			expected: map[string]string{
				"iaas":          "https://iaas.example.com",
				"iam":           "https://iam.example.com",
				"objectstorage": "https://storage.example.com",
			},
		},
		{
			name: "mixed value types (only strings included)",
			endpoints: map[string]any{
				"iaas":  "https://iaas.example.com",  // string: included
				"iam":   12345,                       // int: excluded
				"other": map[string]string{},         // map: excluded
				"space": "https://space.example.com", // string: included
			},
			expected: map[string]string{
				"iaas":  "https://iaas.example.com",
				"space": "https://space.example.com",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeEndpointsFromAny(tt.endpoints)
			if tt.expected == nil {
				if result != nil {
					t.Errorf("expected nil, got %v", result)
				}
			} else {
				if result == nil {
					t.Errorf("expected %v, got nil", tt.expected)
				} else {
					// Check that all expected keys exist
					for k, expectedValue := range tt.expected {
						if actualValue, ok := result[k]; !ok {
							t.Errorf("expected key %s not found", k)
						} else if actualValue != expectedValue {
							t.Errorf("key %s: expected %s, got %s", k, expectedValue, actualValue)
						}
					}
					// Check that all keys are lowercase
					for k := range result {
						if k != strings.ToLower(k) {
							t.Errorf("key %s is not lowercase", k)
						}
					}
					// Check that no extra keys exist
					if len(result) != len(tt.expected) {
						t.Errorf("expected %d keys, got %d keys", len(tt.expected), len(result))
					}
				}
			}
		})
	}
}
