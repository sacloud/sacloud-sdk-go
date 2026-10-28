// Copyright 2026- The sacloud/saclient-go Authors
// SPDX-License-Identifier: Apache-2.0

//go:build unix

package saclient

import (
	"os"
	"syscall"
)

func sanePrivateKeyPermissions(s os.FileInfo) error {
	// Here, "me" is the UID of the current user, not the effective UID of the process.
	// This distinction is important when the process is running with elevated privileges (e.g., setuid)
	// For instance when under `sudo usacloud`, "me" here is that of the original user, not root.
	me := syscall.Getuid()
	stat, ok := s.Sys().(*syscall.Stat_t)
	perm := s.Mode().Perm()

	switch {
	case !ok:
		return NewErrorf("failed to get file stat")

	case perm&0o077 == 0:
		return nil

	case int(stat.Uid) == me:
		return NewErrorf("private key file %s permission is too lax: %o", s.Name(), perm)

	case perm&0o037 != 0:
		return NewErrorf("private key file %s permission is too lax: %o", s.Name(), perm)

	default:
		return nil
	}
}
