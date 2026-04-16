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
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/suite"
)

type HandleRequestTestSuite struct {
	suite.Suite
	client *Client
}

func TestHandleRetries(t *testing.T) { suite.Run(t, new(HandleRequestTestSuite)) }

//nolint:errcheck // this is only a test
func (s *HandleRequestTestSuite) SetupTest() {
	s.client = new(Client)
	// #nosec G104 -- this is only a test
	s.client.SetEnviron([]string{}) // intentional empty setup
}

//nolint:errcheck // this is only a test
func (s *HandleRequestTestSuite) TestNilResponse() {
	// :TRICK: create a httptest.Server that fails
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { /*:UNREACHABLE:*/ }))
	svr.Client().Transport = forceKillRoundTripper{}
	defer svr.Close()

	// #nosec G104 -- this is only a test
	s.client.SetWith(WithTestServer(svr))
	req, _ := http.NewRequest("GET", svr.URL, bytes.NewBuffer([]byte(nil)))
	actual, err := s.client.Do(req)
	s.Nil(actual)
	s.ErrorContains(err, "transport failure")
}

type forceKillRoundTripper struct{}

func (forceKillRoundTripper) RoundTrip(*http.Request) (*http.Response, error) {
	return nil, errors.New("transport failure")
}
