// Copyright 2026- The sacloud/simple-notification-api-go Authors
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

package simplenotification

import (
	"github.com/sacloud/saclient-go"
	v1 "github.com/sacloud/simple-notification-api-go/apis/v1"
)

const (
	defaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1a/api/cloud/1.1/"
	serviceKey        = "simple_notification"
)

// NewClient creates a new simple-notification API client with default settings
func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load message endpoint configuration", err)
	}

	endpoint := defaultAPIRootURL
	if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
		endpoint = ep
	}

	return NewClientWithAPIRootURL(client, endpoint)
}

// NewClientWithAPIRootURL creates a new simple-notification API client with a custom API root URL
func NewClientWithAPIRootURL(client saclient.ClientAPI, apiRootURL string) (*v1.Client, error) {
	clientOption, ok := client.(saclient.ClientOptionAPI)
	if !ok {
		return nil, NewError("client requires saclient.ClientOptionAPI interface", nil)
	}

	newcl, err := clientOption.DupWith(saclient.WithBigInt(false), saclient.WithMiddleware(modifiyMiddleware()))
	if err != nil {
		return nil, err
	}

	return v1.NewClient(apiRootURL, v1.WithClient(newcl))
}
