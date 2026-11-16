// Copyright 2026- The sacloud/saclient-go Authors
// SPDX-License-Identifier: Apache-2.0

package saclient

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewTokenResponseUsesMemoryCacheWhenFileCacheIsUnavailable(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		request := requests.Add(1)
		_, _ = fmt.Fprintf(w, `{"access_token":"token-%d","token_type":"Bearer","expires_in":3600}`, request)
	}))
	defer server.Close()

	blockedDir := filepath.Join(t.TempDir(), "not-a-directory")
	require.NoError(t, os.WriteFile(blockedDir, []byte("block"), 0o600))

	cfg := config{
		"TokenEndpoint":          server.URL,
		"ServicePrincipalID":     "principal",
		"ServicePrincipalKeyKID": "key-id",
		"PrivateKey":             newTestPrivateKeyPEM(t),
		"Profile":                &Profile{dir: blockedDir, Name: "profile"},
	}
	for i := range 2 {
		d := &doer{client: server.Client()}
		response, err := d.newTokenResponse(context.Background(), &cfg)
		require.NoError(t, err)
		require.Equal(t, "token-1", response.AccessToken, "call %d", i+1)
	}
	require.EqualValues(t, 1, requests.Load())
}

func TestNewTokenResponseUsesMemoryCacheWithoutProfile(t *testing.T) {
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		request := requests.Add(1)
		_, _ = fmt.Fprintf(w, `{"access_token":"token-%d","token_type":"Bearer","expires_in":3600}`, request)
	}))
	defer server.Close()

	cfg := config{
		"TokenEndpoint":          server.URL,
		"ServicePrincipalID":     "principal",
		"ServicePrincipalKeyKID": "key-id",
		"PrivateKey":             newTestPrivateKeyPEM(t),
	}
	for i := range 2 {
		d := &doer{client: server.Client()}
		response, err := d.newTokenResponse(context.Background(), &cfg)
		require.NoError(t, err)
		require.Equal(t, "token-1", response.AccessToken, "call %d", i+1)
	}
	require.EqualValues(t, 1, requests.Load())
}

func newTestPrivateKeyPEM(t *testing.T) string {
	t.Helper()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	return string(pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}))
}
