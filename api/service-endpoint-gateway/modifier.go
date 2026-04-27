// Copyright 2026- The sacloud/service-endpoint-gateway-api-go Authors
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

package seg

import (
	"encoding/json"
	"errors"
	"net/http"
	"path"

	"github.com/sacloud/saclient-go"
	v1 "github.com/sacloud/service-endpoint-gateway-api-go/apis/v1"
)

func modifyMiddleware() saclient.Middleware {
	return func(req *http.Request, pull func() (saclient.Middleware, bool)) (*http.Response, error) {
		if err := requestModifier(req); err != nil {
			return nil, err
		}

		cont, ok := pull()
		if !ok {
			return nil, errors.New("middleware chain exhausted")
		}

		resp, err := cont(req, pull)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}
}
func requestModifier(req *http.Request) error {
	// Only modify GET requests to the appliance list endpoint
	listpath := path.Clean(req.URL.Path)
	if path.Base(listpath) == "appliance" &&
		req.Method == http.MethodGet {
		if err := setJSONOnlyQuery(req); err != nil {
			return err
		}
	}
	return nil
}

type filterQuery struct {
	Filter map[string]string `json:"Filter"`
}

func setJSONOnlyQuery(req *http.Request) error {
	q := filterQuery{
		Filter: map[string]string{
			"Class": string(v1.ModelsApplianceApplianceClassServiceendpointgateway),
		},
	}

	b, err := json.Marshal(q)
	if err != nil {
		return err
	}

	req.URL.RawQuery = string(b)
	return nil
}
