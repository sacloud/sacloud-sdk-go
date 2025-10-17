// Copyright 2025- The sacloud/http-client-go Authors
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

package client_test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	. "github.com/sacloud/http-client-go"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/suite"
)

type HttpRequestDoerTestSuite struct {
	suite.Suite
	XDG_CONFIG_HOME *string
	client          *Client
}

func TestHttpRequestDoer(t *testing.T) { suite.Run(t, new(HttpRequestDoerTestSuite)) }

//nolint:errcheck,gosec
func (s *HttpRequestDoerTestSuite) SetupSuite() {
	if current, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		s.XDG_CONFIG_HOME = &current
		os.Setenv("XDG_CONFIG_HOME", s.T().TempDir())
	}
}

//nolint:errcheck,gosec
func (s *HttpRequestDoerTestSuite) TearDownSuite() {
	if s.XDG_CONFIG_HOME != nil {
		os.Setenv("XDG_CONFIG_HOME", *s.XDG_CONFIG_HOME)
	} else {
		os.Unsetenv("XDG_CONFIG_HOME")
	}
}

func (s *HttpRequestDoerTestSuite) SetupTest() {
	s.client = new(Client)
}

func (s *HttpRequestDoerTestSuite) TestSmoke() {
	str := testutil.Random(256, "qwertyuiop")
	j, _ := json.Marshal(map[string]string{"result": str})

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, error := io.ReadAll(r.Body)
		s.NoError(error)
		s.JSONEq(string(j), string(body))
		w.WriteHeader(200)
		_, err := w.Write(j)
		s.NoError(err)
	}))
	defer svr.Close()

	subject, err := s.client.DupWith(WithTestServer(svr))
	s.NoError(err)
	_ = subject.SetEnviron([]string{
		"SAKURACLOUD_ACCESS_TOKEN=foo",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET=bar",
	})

	req, _ := http.NewRequest("GET", svr.URL, bytes.NewReader(j))
	actual, err := subject.Do(req)
	s.NoError(err)
	s.Equal(200, actual.StatusCode)

	actualBody, err := io.ReadAll(actual.Body)
	s.NoError(err)
	s.JSONEq(string(j), string(actualBody))
}

func (s *HttpRequestDoerTestSuite) TestBearer() {
	var requests []*http.Request
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests = append(requests, r)

		body, e := io.ReadAll(r.Body)
		s.NoError(e)
		w.WriteHeader(200)
		s.NoError(e)
		_, e = w.Write(body)
		s.NoError(e)
	}))
	defer svr.Close()

	subject, err := s.client.DupWith(WithTestServer(svr))
	s.NoError(err)
	_ = subject.SetEnviron([]string{
		"SAKURACLOUD_SERVICE_PRINCIPAL_ID=113702516320",
		`SAKURACLOUD_PRIVATE_KEY=
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC/AfcvlUlhcPpDD/1HqWBWvGQZxdER+fy6jbm1BlhVT156hjZi
UUwetUMrVGuy+bYE50j+qJB2VYKIhUTIUYCJg/AlruszlmydV0dOWPpSsLMXA5XU
GhoijaZY9l8vsbGN3n0QJ313GvFQQ+GrP1PzmRbpK686weAwtCx+PXYPQwIDAQAB
AoGAEvG69nk0AfoWmDgpwsXFzFR7CSNZjRLiQg50cMPkVvG8SSKumim+Bv2rX8zL
scCakPnvf3JwgYwRmkC9hbCvssfQK2o0Zzc6zPa560TxXYK5rADTfMXqeLnF6nFZ
sKLlE5vxyv2XD6zDcc1K2q25ARYMeWOGQ2WfuMYexBd36EECQQD0va3JquOaPQI7
2yRXNumv2fRwYohnJxOymu4vKZp11R0gTGljsv7y8I+mcVDJnJy27t9a7tUSLS4F
G1FMId0LAkEAx8t39aRzchpUoJYl9KmigFQ5AS6qAmDqdGIOBFQ5hf6HErukbRBd
2q+tNXAKF62ecXR3dlaS54CpSXkQVxlJqQJBANJD1/hIEk0kFzQ3nSw06GaFmcWo
UcpVv02WYAYy9xo/I0vpei4GzZUI6lG0TxU3sUhVR53HTVXVbRFEG/+NpGsCQQCi
qPilOJn0z5MOmq+UHXd7WxZ96+vlu9mlnx8iTx/2A18c1T/su2Jt5JDz7J+K34Mb
g2KvKZS4fXtVoga3opLhAkAtR4iVtxGi3NxOw0XrTXClzJD1e357/MrSDQ09gdRG
sP9Knwr9WVBtRYPRFjC3YccLTwoQnjVcF1qJN6ybMvnS
-----END RSA PRIVATE KEY-----
		`,
	})

	req, _ := http.NewRequest("GET", svr.URL, strings.NewReader("{}"))
	actual, err := subject.Do(req)
	s.NoError(err)
	s.Equal(200, actual.StatusCode)

	actualBody, err := io.ReadAll(actual.Body)
	s.NoError(err)
	s.JSONEq("{}", string(actualBody))

	s.Len(requests, 1)
	s.Equal("Bearer ", requests[0].Header.Get("Authorization")[:7])
}
