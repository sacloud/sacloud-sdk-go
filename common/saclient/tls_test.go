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

package saclient_test

import (
	"encoding/pem"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

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

func (s *TLSTestSuite) TestWithCACertFile() {
	svr, pemBytes := s.newTLSServer()

	path := filepath.Join(s.T().TempDir(), "ca.pem")
	s.Require().NoError(os.WriteFile(path, pemBytes, 0o600))

	subject, err := s.client.DupWith(WithCACertFile(path))
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(tlsDummyEnviron))
	s.Require().NoError(subject.Populate())

	status, err := s.doGet(subject, svr)
	s.NoError(err)
	s.Equal(http.StatusOK, status)
}

func (s *TLSTestSuite) TestCACertEnv() {
	svr, pemBytes := s.newTLSServer()

	path := filepath.Join(s.T().TempDir(), "ca.pem")
	s.Require().NoError(os.WriteFile(path, pemBytes, 0o600))

	subject, err := s.client.DupWith()
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(append([]string{"SAKURA_CACERT=" + path}, tlsDummyEnviron...)))
	s.Require().NoError(subject.Populate())

	status, err := s.doGet(subject, svr)
	s.NoError(err)
	s.Equal(http.StatusOK, status)
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

	subject, err := s.client.DupWith(WithCACertFile(filepath.Join(s.T().TempDir(), "missing.pem")))
	s.Require().NoError(err)
	s.Require().NoError(subject.SetEnviron(tlsDummyEnviron))
	err = subject.Populate()
	s.ErrorContains(err, "cannot read CA certificate file")
}
