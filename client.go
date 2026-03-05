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
	"runtime"

	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1a/api/monitoring/1.0/"

	// エンドポイントサービスキー
	ServiceKey = "monitoring_suite"
)

var (
	// UserAgent APIリクエスト時のユーザーエージェント
	UserAgent = fmt.Sprintf(
		"monitoring-suite-api-go/%s (%s/%s; +https://github.com/sacloud/monitoring-suite-api-go)",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
	)
)

func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
	endpoint := DefaultAPIRootURL

	if endpointConfig, err := client.EndpointConfig(); err != nil {
		return nil, NewError("unable to load endpoint configuration", err)
	} else if ep, ok := endpointConfig.Endpoints[ServiceKey]; ok && ep != "" {
		endpoint = ep
	}

	return NewClientWithApiUrl(endpoint, client)
}

func NewClientWithApiUrl(apiUrl string, client saclient.ClientAPI) (*v1.Client, error) {
	dupable, ok := client.(*saclient.Client)
	if !ok {
		return nil, NewError("NewClientWithApiUrl", fmt.Errorf("client must be *saclient.Client"))
	}

	augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
		saclient.WithRootURL(apiUrl),
		saclient.WithBigInt(true),
	)

	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	d, err := v1.NewClient(apiUrl, v1.WithClient(augmented))
	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	return d, nil
}
