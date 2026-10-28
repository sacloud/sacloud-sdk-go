// Copyright 2026- The sacloud/saclient-go Authors
// SPDX-License-Identifier: Apache-2.0

//go:build unix

package saclient

import (
	"os"
	"syscall"
	"testing"
)

type fakeFileInfo struct {
	os.FileInfo
	mode os.FileMode
	sys  any
}

func (i fakeFileInfo) Mode() os.FileMode { return i.mode }
func (i fakeFileInfo) Sys() any          { return i.sys }

func TestPrivateKeyPermissionsTooLax(t *testing.T) {
	base, err := os.Stat(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		mode       os.FileMode
		ownerUID   int
		wantTooLax bool
	}{
		{name: "group-readable by real user", mode: 0o640, ownerUID: os.Getuid(), wantTooLax: true},
		{name: "group-readable by another user", mode: 0o640, ownerUID: os.Getuid() + 1},
		{name: "group-writable by another user", mode: 0o660, ownerUID: os.Getuid() + 1, wantTooLax: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := fakeFileInfo{
				FileInfo: base,
				mode:     tt.mode,
				sys: &syscall.Stat_t{
					Uid: uint32(tt.ownerUID), // #nosec G115 -- no overflow here
				},
			}
			if err := sanePrivateKeyPermissions(info); err == nil {
				// expected
			} else if tt.wantTooLax {
				// expected
			} else {
				t.Errorf("sanePrivateKeyPermissions() = %v, want %v", err, !tt.wantTooLax)
			}
		})
	}
}
