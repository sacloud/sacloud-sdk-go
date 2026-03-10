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

package cloudhsm

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1b/api/cloud/1.1/"

	// DefaultEndpoint デフォルトのエンドポイントURL
	DefaultEndpoint = "https://secure.sakura.ad.jp/cloud/zone/"

	// DefaultZone クラウドHSM設置ゾーンのデフォルト
	DefaultZone = "is1b"
)

var (
	// UserAgent APIリクエスト時のユーザーエージェント
	UserAgent = fmt.Sprintf(
		"cloudhsm-api-go/%s (%s/%s; +https://github.com/sacloud/cloudhsm-api-go)",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
	)

	// エンドポイントサービスキー
	ServiceKey = "cloudhsm"
)

type EmptySecuritySource struct{}

func (this EmptySecuritySource) BasicAuth(ctx context.Context, operationName v1.OperationName) (v1.BasicAuth, error) {
	return v1.BasicAuth{}, nil
}

func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
	const path = "api/cloud/1.1/"
	zone := DefaultZone
	endpoint := DefaultEndpoint

	cfg, err := client.EndpointConfig()

	if err != nil {
		return nil, NewError("NewClient", err)
	}

	if ep, ok := cfg.Endpoints[ServiceKey]; ok && ep != "" {
		endpoint = ep
	}

	if cfg.Zone != "" {
		zone = cfg.Zone
	}

	apiUrl := fmt.Sprintf(
		"%s/%s/%s",
		strings.TrimSuffix(endpoint, "/"),
		strings.TrimPrefix(strings.TrimSuffix(zone, "/"), "/"),
		path,
	)

	return NewClientWithApiUrl(apiUrl, client)
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
		// これはなにか:
		// EmptySecuritySource.BasicAuth()がBasic認証を生成
		// しかし実際の通信で必ずしもBasic認証が使われると限らない
		//　そのあたりをsaclient-go側で吸収させる設定が下記↓
		saclient.WithForceAutomaticAuthentication(),
	)

	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	d, err := v1.NewClient(apiUrl, EmptySecuritySource{}, v1.WithClient(augmented))
	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	return d, nil
}
