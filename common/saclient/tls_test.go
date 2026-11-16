// Copyright 2026- The sacloud/saclient-go Authors
// SPDX-License-Identifier: Apache-2.0

package saclient_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	. "github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/suite"
)

// TLSTestSuite tests trusting a custom root CA, which is how a client reaches
// an endpoint served with a privately signed certificate (a development mock
// server or a staging endpoint).
type TLSTestSuite struct {
	suite.Suite
	client *Client
}

func TestTLS(t *testing.T) { suite.Run(t, new(TLSTestSuite)) }

func (s *TLSTestSuite) SetupTest() {
	s.client = new(Client)
}

var tlsDummyEnviron = []string{
	"SAKURACLOUD_ACCESS_TOKEN=foo",
	"SAKURACLOUD_ACCESS_TOKEN_SECRET=bar",
}

// newTLSServer starts a TLS test server answering 200 to GET /, and returns it
// along with its certificate in PEM format.
func (s *TLSTestSuite) newTLSServer() (*httptest.Server, []byte) {
	svr := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	s.T().Cleanup(svr.Close)

	cert := svr.Certificate()
	s.Require().NotNil(cert)
	pemBytes := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: cert.Raw})
	s.Require().NotEmpty(pemBytes)
	return svr, pemBytes
}

// doGet issues GET svr.URL with the given client and returns the status code.
func (s *TLSTestSuite) doGet(subject ClientAPI, svr *httptest.Server) (int, error) {
	req, err := http.NewRequest("GET", svr.URL, nil)
	s.Require().NoError(err)
	resp, err := subject.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	s.Require().NoError(err)
	return resp.StatusCode, nil
}

func (s *TLSTestSuite) TestWithCACertPEM() {
	svr, pemBytes := s.newTLSServer()

	// Without the CA the self-signed certificate is rejected.
	subject, err := s.client.DupWith()
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(tlsDummyEnviron))
	s.Require().NoError(subject.Populate())
	_, err = s.doGet(subject, svr)
	s.Error(err)

	// With the CA PEM the request succeeds.
	subject, err = s.client.DupWith(WithCACertPEM(pemBytes))
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(tlsDummyEnviron))
	s.Require().NoError(subject.Populate())
	status, err := s.doGet(subject, svr)
	s.NoError(err)
	s.Equal(http.StatusOK, status)
}

func (s *TLSTestSuite) TestWithCACertPath() {
	svr, pemBytes := s.newTLSServer()

	path := filepath.Join(s.T().TempDir(), "ca.pem")
	s.Require().NoError(os.WriteFile(path, pemBytes, 0o600))

	subject, err := s.client.DupWith(WithCACertPath(path))
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(tlsDummyEnviron))
	s.Require().NoError(subject.Populate())

	status, err := s.doGet(subject, svr)
	s.NoError(err)
	s.Equal(http.StatusOK, status)
}

func (s *TLSTestSuite) TestCACertPathEnv() {
	svr, pemBytes := s.newTLSServer()

	path := filepath.Join(s.T().TempDir(), "ca.pem")
	s.Require().NoError(os.WriteFile(path, pemBytes, 0o600))

	subject, err := s.client.DupWith()
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(append([]string{"SAKURA_CA_CERT_PATH=" + path}, tlsDummyEnviron...)))
	s.Require().NoError(subject.Populate())

	status, err := s.doGet(subject, svr)
	s.NoError(err)
	s.Equal(http.StatusOK, status)
}

// TestCACertPEMEnv proves the SAKURA_CA_CERT environment variable accepts the
// PEM content itself (like SAKURA_PRIVATE_KEY does), not just a path.
func (s *TLSTestSuite) TestCACertPEMEnv() {
	svr, pemBytes := s.newTLSServer()

	subject, err := s.client.DupWith()
	s.Require().NoError(err)
	// The PEM spans multiple lines; intoEnvmap splits at the first '=' so the
	// whole PEM becomes the value.
	environ := append([]string{"SAKURA_CA_CERT=" + string(pemBytes)}, tlsDummyEnviron...)
	s.Require().NoError(subject.SetEnviron(environ))
	s.Require().NoError(subject.Populate())

	status, err := s.doGet(subject, svr)
	s.NoError(err)
	s.Equal(http.StatusOK, status)
}

// TestWithCACertPEMChain exercises the realistic CA setup end to end: a
// self-signed CA issues a server certificate, the server presents that leaf
// (with the CA in its chain), and a client trusting only the CA's PEM — as it
// would when a development mock server emits its own CA — completes a real
// HTTPS request. Without the CA the handshake fails.
func (s *TLSTestSuite) TestWithCACertPEMChain() {
	// A self-signed CA.
	caKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	s.Require().NoError(err)
	caTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(1),
		Subject:               pkix.Name{Organization: []string{"saclient-test"}, CommonName: "test CA"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign,
		BasicConstraintsValid: true,
		IsCA:                  true,
	}
	caDER, err := x509.CreateCertificate(rand.Reader, caTmpl, caTmpl, &caKey.PublicKey, caKey)
	s.Require().NoError(err)
	caCert, err := x509.ParseCertificate(caDER)
	s.Require().NoError(err)
	caPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: caDER})

	// A server certificate for the httptest listener, signed by that CA.
	leafKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	s.Require().NoError(err)
	leafTmpl := &x509.Certificate{
		SerialNumber:          big.NewInt(2),
		Subject:               pkix.Name{Organization: []string{"saclient-test"}, CommonName: "localhost"},
		NotBefore:             time.Now().Add(-time.Hour),
		NotAfter:              time.Now().Add(24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{"localhost"},
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1"), net.ParseIP("::1")},
	}
	leafDER, err := x509.CreateCertificate(rand.Reader, leafTmpl, caCert, &leafKey.PublicKey, caKey)
	s.Require().NoError(err)

	svr := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	svr.TLS = &tls.Config{Certificates: []tls.Certificate{{
		Certificate: [][]byte{leafDER, caDER}, // leaf followed by its issuer, as real servers serve
		PrivateKey:  leafKey,
	}}}
	svr.StartTLS()
	s.T().Cleanup(svr.Close)

	// Trusting the CA PEM makes the HTTPS request succeed.
	subject, err := s.client.DupWith(WithCACertPEM(caPEM))
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(tlsDummyEnviron))
	s.Require().NoError(subject.Populate())
	status, err := s.doGet(subject, svr)
	s.NoError(err)
	s.Equal(http.StatusOK, status)

	// Without the CA the handshake is rejected.
	bare, err := s.client.DupWith()
	s.Require().NoError(err)
	s.Require().NoError(bare.SetEnviron(tlsDummyEnviron))
	s.Require().NoError(bare.Populate())
	_, err = s.doGet(bare, svr)
	s.Error(err)
}

func (s *TLSTestSuite) TestCACertInvalidPEM() {
	svr, _ := s.newTLSServer()

	subject, err := s.client.DupWith(WithCACertPEM([]byte("not a certificate")))
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(tlsDummyEnviron))
	err = subject.Populate()
	s.ErrorContains(err, "no certificate found in CA certificate PEM")
	_ = svr
}

func (s *TLSTestSuite) TestCACertFileNotFound() {
	_, _ = s.newTLSServer()

	subject, err := s.client.DupWith(WithCACertPath(filepath.Join(s.T().TempDir(), "missing.pem")))
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(tlsDummyEnviron))
	err = subject.Populate()
	s.ErrorContains(err, "cannot read CA certificate file")
}
