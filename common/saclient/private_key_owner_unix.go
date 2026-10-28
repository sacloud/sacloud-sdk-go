// Copyright 2026- The sacloud/saclient-go Authors
// SPDX-License-Identifier: Apache-2.0

//go:build unix

package saclient

import (
	"os"
	"syscall"
)

func sanePrivateKeyPermissions(s os.FileInfo, v string) error {
	// Here, "me" is the UID of the current user, not the effective UID of the process.
	// This distinction is important when the process is running with elevated privileges (e.g., setuid)
	// When invoked via sudo, the real UID is normally the target user's UID (root by default).
	me := syscall.Getuid()
	stat, ok := s.Sys().(*syscall.Stat_t)
	perm := s.Mode().Perm()

	switch {
	case !ok:
		return NewErrorf("failed to get file stat")

	case perm&0o077 == 0:
		return nil

	case int(stat.Uid) == me:
		return NewErrorf("private key file %s permission is too lax: %o", v, perm)

	case perm&0o037 != 0:
		return NewErrorf("private key file %s permission is too lax: %o", v, perm)

	default:
		return nil
	}
}
