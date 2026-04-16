// Copyright 2025- The sacloud/apigw-api-go authors
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

package apigw

import (
	"fmt"
	"runtime"

	v1 "github.com/sacloud/apigw-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/api/apigw/1.0/"

	// ServiceKey SDKの種別を示すキー、プロファイルでのエンドポイント取得に利用する
	ServiceKey = "apigw"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"apigw-api-go/%s (%s/%s; +https://github.com/sacloud/apigw-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load endpoint configuration", err)
	}
	endpoint := DefaultAPIRootURL
	if ep, ok := endpointConfig.Endpoints[ServiceKey]; ok && ep != "" {
		endpoint = ep
	}
	return NewClientWithAPIRootURL(client, endpoint)
}

func NewClientWithAPIRootURL(client saclient.ClientAPI, apiRootURL string) (*v1.Client, error) {
	dupable, ok := client.(saclient.ClientOptionAPI)
	if !ok {
		return nil, NewError("client does not implement saclient.ClientOptionAPI", nil)
	}

	augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
	)
	if err != nil {
		return nil, err
	}
	return v1.NewClient(apiRootURL, v1.WithClient(augmented))
}
