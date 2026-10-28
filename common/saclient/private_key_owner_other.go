// Copyright 2026- The sacloud/saclient-go Authors
// SPDX-License-Identifier: Apache-2.0

//go:build !unix

package saclient

import "os"

func sanePrivateKeyPermissions(os.FileInfo, string) error {
	// There is no concept like world writable directory on Windows.
	// Permission checks make no sense here.
	return nil
}
