// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite

import (
	"fmt"
	"net/http"
	"runtime"

	client "github.com/sacloud/api-client-go"
	saht "github.com/sacloud/go-http"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1a/api/monitoring/1.0/"
)

var (
	// UserAgent APIリクエスト時のユーザーエージェント
	UserAgent = fmt.Sprintf(
		"monitoring-suite-api-go/%s (%s/%s; +https://github.com/sacloud/monitoring-suite-api-go) %s",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
		client.DefaultUserAgent,
	)

	RequestCustomizers = []saht.RequestCustomizer{
		func(req *http.Request) error {
			req.Header.Set("X-Sakura-Bigint-As-Int", "0")
			return nil
		},
	}
)

func NewClient(params ...client.ClientParam) (*v1.Client, error) {
	return NewClientWithApiUrl(DefaultAPIRootURL, params...)
}

func NewClientWithApiUrl(apiUrl string, params ...client.ClientParam) (*v1.Client, error) {
	return NewClientWithApiUrlAndClient(apiUrl, nil, params...)
}

func NewClientWithApiUrlAndClient(apiUrl string, apiClient *http.Client, params ...client.ClientParam) (*v1.Client, error) {
	var cli client.ClientParam
	if apiClient == nil {
		cli = func(i *client.ClientParams) {}
	} else {
		cli = client.WithHTTPClient(apiClient)
	}
	ua := client.WithUserAgent(UserAgent)
	opts := client.WithOptions(&client.Options{RequestCustomizers: RequestCustomizers})
	c, err := client.NewClient(apiUrl, append(params, ua, cli, opts)...)
	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	d, err := v1.NewClient(c.ServerURL(), v1.WithClient(c.NewHttpRequestDoer()))
	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	return d, nil
}
