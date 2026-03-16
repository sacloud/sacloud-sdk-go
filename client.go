// Copyright 2022-2025 The sacloud/dedicated-storage-api-go Authors
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

package dedicatedstorage

import (
	"fmt"
	"runtime"

	v1 "github.com/sacloud/dedicated-storage-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

const (
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/tk1b/api/cloud/1.0/"
	serviceKey        = "dedicated_storage"
)

var UserAgent = fmt.Sprintf(
	"dedicated-storage-api-go/%s (%s/%s; +https://github.com/sacloud/dedicated-storage-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

func NewClient(client *saclient.Client) (*v1.Client, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load endpoint configuration", err)
	}

	apiURL := DefaultAPIRootURL
	if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
		apiURL = ep
	}
	return NewClientWithAPIRootURL(client, apiURL)
}

func NewClientWithAPIRootURL(client *saclient.Client, apiRootURL string) (*v1.Client, error) {
	c, err := client.DupWith(saclient.WithUserAgent(UserAgent))
	if err != nil {
		return nil, err
	}
	return v1.NewClient(apiRootURL, v1.WithClient(c))
}
