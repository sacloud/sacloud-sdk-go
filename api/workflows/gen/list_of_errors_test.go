// Copyright 2026- The sacloud/workflows-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	path := filepath.Join(t.TempDir(), "schemas.go")
	source := `package v1
type ZError struct { IsOk bool; Message string }
func (*ZError) operationRes() {}
type AError struct { IsOk bool; Message string }
func (*AError) operationRes() {}
type NotResponse struct { IsOk bool; Message string }
type Success struct { IsOk bool }
func (*Success) operationRes() {}
`
	require.NoError(t, os.WriteFile(path, []byte(source), 0o600))
	types, err := list(path)
	require.NoError(t, err)
	require.Equal(t, []string{"AError", "ZError"}, types)
}

func TestListWithoutErrors(t *testing.T) {
	path := filepath.Join(t.TempDir(), "schemas.go")
	require.NoError(t, os.WriteFile(path, []byte("package v1\n"), 0o600))
	_, err := list(path)
	require.ErrorContains(t, err, "no error response types")
}

func TestStatusCodeUnknown(t *testing.T) {
	_, err := statusCode("OperationUnknown")
	require.ErrorContains(t, err, "unknown HTTP status in error type OperationUnknown")
}
