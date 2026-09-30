// Copyright 2026- The sacloud/saclient-go Authors
// SPDX-License-Identifier: Apache-2.0

package saclient

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/suite"
)

type funcRoundTripper func(*http.Request) (*http.Response, error)

func (f funcRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) { return f(req) }

type timeoutError struct{}

func (timeoutError) Error() string   { return "timeout" }
func (timeoutError) Timeout() bool   { return true }
func (timeoutError) Temporary() bool { return true }

type RFC7523TestSuite struct{ suite.Suite }

func TestRFC7523(t *testing.T) { suite.Run(t, new(RFC7523TestSuite)) }

func (s *RFC7523TestSuite) TestInquireAccessTokenEOF() {
	require := s.Require()

	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(err)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	endpoint := "https://user:password@example.com/oauth2/token?secret=query#fragment" // #nosec G101 -- this is a test key and not a real secret
	cfg := config{
		"TokenEndpoint":         endpoint,
		"ServicePrincipalID":    "service-principal-id",
		"ServicePrincipalKeyID": "service-principal-key-id",
		"PrivateKey":            string(privateKeyPEM),
	}

	d := doer{
		client: &http.Client{
			Transport: funcRoundTripper(func(*http.Request) (*http.Response, error) {
				return nil, io.EOF
			}),
		},
	}

	_, err = d.inquireAccessToken(context.Background(), &cfg)
	require.Error(err)
	require.ErrorIs(err, io.EOF)

	var urlErr *url.Error
	require.ErrorAs(err, &urlErr)
	require.Equal("https://redacted:redacted@example.com/oauth2/token?secret=query#fragment", urlErr.URL)

	message := err.Error()
	require.Contains(message, "token endpoint request failed (EOF before response)")
	require.NotContains(message, "password")
}

func (s *RFC7523TestSuite) TestTokenEndpointRequestErrorClassification() {
	require := s.Require()
	tests := []struct {
		name  string
		cause error
		kind  string
	}{
		{name: "canceled", cause: context.Canceled, kind: "runtime context canceled"},
		{name: "deadline", cause: context.DeadlineExceeded, kind: "runtime deadline exceeded"},
		{name: "unexpected", cause: io.ErrUnexpectedEOF, kind: "unexpected EOF before response"},
		{name: "DNS", cause: &net.DNSError{Err: "no such host", Name: "example.com"}, kind: "DNS failure"},
		{name: "TLS", cause: tls.RecordHeaderError{Msg: "bad record header"}, kind: "TLS handshake failure"},
		{name: "timeout", cause: timeoutError{}, kind: "timeout"},
		{name: "transport", cause: errors.New("connection reset"), kind: "transport failure"},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			wrapped := fmt.Errorf("%w", tt.cause)
			err := classifiedRequestError(wrapped, "<<test>>")

			require.ErrorIs(err, tt.cause)
			require.Contains(err.Error(), "<<test>> ("+tt.kind+")")
		})
	}
}
