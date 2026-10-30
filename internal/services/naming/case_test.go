// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package naming

import "testing"

func TestCaseConversions(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
		fn   func(string) string
	}{
		{name: "snake acronym", in: "HTTPServer", want: "http_server", fn: SnakeCase},
		{name: "kebab digits", in: "HTTP20xOK", want: "http-20x-ok", fn: KebabCase},
		{name: "pascal", in: "ipv6_enable", want: "Ipv6Enable", fn: PascalCase},
		{name: "camel", in: "http_server", want: "httpServer", fn: CamelCase},
		{name: "unicode", in: "ÉclairName", want: "éclair_name", fn: SnakeCase},
		{name: "separators", in: "some words-and_names", want: "some-words-and-names", fn: KebabCase},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fn(tt.in); got != tt.want {
				t.Errorf("%s(%q) = %q, want %q", tt.name, tt.in, got, tt.want)
			}
		})
	}
}

func TestFirstRuneToLower(t *testing.T) {
	if got := FirstRuneToLower("Éclair"); got != "éclair" {
		t.Errorf("FirstRuneToLower() = %q, want %q", got, "éclair")
	}
}
