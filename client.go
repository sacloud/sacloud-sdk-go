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
	"net/url"
	"path"

	"github.com/sacloud/saclient-go"
	v1 "github.com/sacloud/service-endpoint-gateway-api-go/apis/v1"
)

const (
	defaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1a/api/cloud/1.1/"
	serviceKey        = "service_endpoint_gateway"

	// used zone info from saclient
	defaultEndpoint = "https://secure.sakura.ad.jp/cloud/zone/"
)

func NewClient(clientAPI saclient.ClientAPI) (*v1.Client, error) {
	endpointConfig, err := clientAPI.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load service endpoint configuration", err)
	}
	endpoint := defaultAPIRootURL

	if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
		endpoint = ep
	}

	if endpointConfig.Zone != "" {
		u, err := url.Parse(defaultEndpoint)
		if err != nil {
			return nil, NewError("invalid endpoint URL", err)
		}
		u.Path = path.Join(u.Path, endpointConfig.Zone, "api", "cloud", "1.1")
		endpoint = u.String()
		// Ensure the endpoint URL ends with a slash
		endpoint += "/"
	}

	return NewClientWithAPIRootURL(clientAPI, endpoint)
}

// NewClientWithAPIRootURL creates a new service-endpoint-gateway API client with a custom API root URL
func NewClientWithAPIRootURL(clientAPI saclient.ClientAPI, apiRootURL string) (*v1.Client, error) {
	// cast the clientAPI to *saclient.Client to access DupWith method
	client, ok := clientAPI.(*saclient.Client)
	if !ok {
		return nil, NewError("invalid client type", nil)
	}

	newcl, err := client.DupWith(saclient.WithBigInt(false),
		saclient.WithMiddleware(modifyMiddleware()))
	if err != nil {
		return nil, err
	}

	return v1.NewClient(apiRootURL, v1.WithClient(newcl))
}
