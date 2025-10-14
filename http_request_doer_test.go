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

	req, _ := http.NewRequest("GET", svr.URL, bytes.NewReader(j))
	actual, err := subject.Do(req)
	s.NoError(err)
	s.Equal(200, actual.StatusCode)

	actualBody, err := io.ReadAll(actual.Body)
	s.NoError(err)
	s.JSONEq(string(j), string(actualBody))
}
