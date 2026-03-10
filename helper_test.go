// Copyright 2025- The sacloud/cloudhsm-api-go Authors
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

package cloudhsm_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/sacloud/cloudhsm-api-go"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

var theClient saclient.Client

type ErrorResponse struct {
	Message string `json:"error_msg"`
	IsOk    bool   `json:"is_ok"`
}

func newErrorResponse(message string) ErrorResponse {
	return ErrorResponse{
		Message: message,
		IsOk:    false,
	}
}

func newTestClient(v any, s ...int) *v1.Client {
	s = append(s, http.StatusOK)
	j, e := json.Marshal(v)
	if e != nil {
		panic(e)
	}

	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		st := s[0]

		w.WriteHeader(st)
		if st == http.StatusNoContent {
			return
		}
		if _, e = w.Write(j); e != nil {
			panic(e)
		}
	})
	sv := httptest.NewServer(h)
	api, e := theClient.DupWith(saclient.WithTestServer(sv))
	if e != nil {
		panic(e)
	}
	c, e := NewClientWithApiUrl(sv.URL, api)
	if e != nil {
		panic(e)
	}
	return c
}

func newIntegratedClient(t *testing.T) *v1.Client {
	testutil.PreCheckEnvsFunc(
		"SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET",
	)(t)

	apiUrl := DefaultAPIRootURL
	if root, ok := os.LookupEnv("SAKURA_LOCAL_ENDPOINT_CLOUDHSM"); ok {
		apiUrl = root
	}
	api, err := theClient.DupWith(saclient.WithTraceMode("error"))
	require.NoError(t, err)
	ret, err := NewClientWithApiUrl(apiUrl, api)
	require.NoError(t, err)
	return ret
}

func ref[T any](v T) *T {
	return &v
}

var TemplateDateTime = func() v1.DateTime {
	var ret v1.DateTime
	ret.SetFake()
	return ret
}

var TemplateTags = []string{"tag1", "tag2"}

var TemplateLicense = func() v1.CloudHSMSoftwareLicense {
	var ret v1.CloudHSMSoftwareLicense
	ret.SetFake()
	ret.SetTags(TemplateTags)

	return ret
}()

var TemplateCreateLicense = func() v1.CreateCloudHSMSoftwareLicense {
	var ret v1.CreateCloudHSMSoftwareLicense
	ret.SetFake()
	ret.SetTags(TemplateTags)

	return ret
}()

var TemplateWrappedCreateLicense = func() v1.WrappedCreateCloudHSMSoftwareLicense {
	var ret v1.WrappedCreateCloudHSMSoftwareLicense
	ret.SetLicense(v1.NewOptCreateCloudHSMSoftwareLicense(TemplateCreateLicense))

	return ret
}()

var TemplateWrappedLicense = func() v1.WrappedCloudHSMSoftwareLicense {
	var ret v1.WrappedCloudHSMSoftwareLicense
	ret.SetLicense(v1.NewOptCloudHSMSoftwareLicense(TemplateLicense))

	return ret
}()

var TemplateCloudHSM = func() v1.CloudHSM {
	var ret v1.CloudHSM
	ret.SetFake()
	ret.SetTags(TemplateTags)
	ret.SetAvailability(v1.AvailabilityEnumAvailable)

	return ret
}()

var TemplateCreateCloudHSM = func() v1.CreateCloudHSM {
	var ret v1.CreateCloudHSM
	ret.SetFake()
	ret.SetTags(TemplateTags)

	return ret
}()

var TemplateWrappedCreateCloudHSM = func() v1.WrappedCreateCloudHSM {
	var ret v1.WrappedCreateCloudHSM
	ret.SetCloudHSM(TemplateCreateCloudHSM)

	return ret
}()

var TemplateWrappedCloudHSM = func() v1.WrappedCloudHSM {
	var ret v1.WrappedCloudHSM
	ret.SetCloudHSM(TemplateCloudHSM)

	return ret
}()

var TemplateCloudHSMPeer = func() v1.CloudHSMPeer {
	var ret v1.CloudHSMPeer
	ret.SetFake()

	return ret
}()

var TemplateCreateCloudHSMPeer = func() v1.CreateCloudHSMPeer {
	var ret v1.CreateCloudHSMPeer
	ret.SetFake()

	return ret
}()

var TemplateWrappedCreateCloudHSMPeer = func() v1.WrappedCreateCloudHSMPeer {
	var ret v1.WrappedCreateCloudHSMPeer
	ret.SetPeer(TemplateCreateCloudHSMPeer)

	return ret
}()

var TemplateCloudHSMClient = func() v1.CloudHSMClient {
	var ret v1.CloudHSMClient
	ret.SetFake()
	ret.SetAvailability(v1.AvailabilityEnumAvailable)
	return ret
}()

var TemplateWrappedCloudHSMClient = func() v1.WrappedCloudHSMClient {
	var ret v1.WrappedCloudHSMClient
	ret.SetClient(TemplateCloudHSMClient)

	return ret
}

var TemplateCreateCloudHSMClient = func() v1.CreateCloudHSMClient {
	var ret v1.CreateCloudHSMClient
	ret.SetFake()
	ret.SetAvailability(v1.AvailabilityEnumAvailable)

	return ret
}()

var TemplateWrappedCreateCloudHSMClient = func() v1.WrappedCreateCloudHSMClient {
	var ret v1.WrappedCreateCloudHSMClient
	ret.SetClient(TemplateCreateCloudHSMClient)

	return ret
}()
