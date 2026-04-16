// Copyright 2025- The sacloud/addon-api-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package addon_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-faster/jx"
	"github.com/sacloud/addon-api-go"
	v1 "github.com/sacloud/addon-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

var theClient saclient.Client

type encodable interface {
	Encode(*jx.Encoder)
}

func newTestClient(v encodable, s ...int) *v1.Client {
	s = append(s, http.StatusOK)

	sv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		st := s[0]

		w.WriteHeader(st)
		if st == http.StatusNoContent {
			return
		}
		if v == nil {
			return
		}
		e := new(jx.Encoder)
		v.Encode(e)
		_, _ = w.Write(e.Bytes())
	}))

	sa, e := theClient.DupWith(saclient.WithTestServer(sv))
	if e != nil {
		panic(e)
	}

	c, e := addon.NewClientWithAPIRootURL(sa, sv.URL)
	if e != nil {
		panic(e)
	}

	if e := sa.Populate(); e != nil {
		panic(e)
	}

	return c
}

func IntegraetdClient(t *testing.T) (assert *require.Assertions, client *v1.Client) {
	var err error
	testutil.PreCheckEnvsFunc("TESTACC")(t)

	assert = require.New(t)
	client, err = addon.NewClient(&theClient)

	assert.NoError(err)

	return
}

var MockResourceGroupResource v1.ResourceGroupResource = func() (mock v1.ResourceGroupResource) {
	mock.SetFake()
	return
}()

var MockListResourcesResponse v1.ListResourcesResponse = func() (mock v1.ListResourcesResponse) {
	mock.SetFake()
	mock.SetResources(v1.NewOptNilResourceGroupResourceArray([]v1.ResourceGroupResource{
		MockResourceGroupResource,
	}))
	return
}()

var MockPostDeploymentResponse v1.PostDeploymentResponse = func() (mock v1.PostDeploymentResponse) {
	mock.SetFake()
	return
}()

var MockErrorInfo v1.ErrorInfo = func() (mock v1.ErrorInfo) {
	mock.SetFake()
	return
}()

var MockErrorResponse v1.ErrorResponse = func() (mock v1.ErrorResponse) {
	mock.SetFake()
	mock.SetErrors(v1.NewOptNilErrorInfoArray([]v1.ErrorInfo{
		MockErrorInfo,
	}))
	return
}()

var MockResourceResponse v1.GetResourceResponse = func() (mock v1.GetResourceResponse) {
	var raw jx.Raw = []byte(`{"test": "data"}`)
	mock.SetFake()
	mock.SetData(raw)
	return
}()

var MockDeploymentStatusProperties v1.DeploymentStatusProperties = func() (mock v1.DeploymentStatusProperties) {
	mock.SetFake()

	// Don't rely on local time zone
	mock.SetTimestamp(v1.NewOptDateTime(time.Time{}))
	return
}()

var MockDeploymentStatus v1.DeploymentStatus = func() (mock v1.DeploymentStatus) {
	mock.SetFake()
	mock.SetID(v1.NewOptNilString(testutil.RandomName("test-id-", 32, testutil.CharSetAlphaNum)))
	mock.SetName(v1.NewOptNilString(testutil.RandomName("test-name-", 32, testutil.CharSetAlphaNum)))
	mock.SetType(v1.NewOptNilString(testutil.RandomName("test-type-", 32, testutil.CharSetAlphaNum)))
	mock.SetProperties(v1.NewOptDeploymentStatusProperties(MockDeploymentStatusProperties))
	return
}()

var MockFrontDoorOrigin v1.FrontDoorOrigin = func() (mock v1.FrontDoorOrigin) {
	mock.SetFake()
	return
}()

var MockVulnerabilityResponseBody v1.VulnerabilityResponseBody = func() (mock v1.VulnerabilityResponseBody) {
	mock.SetFake()
	return
}()
