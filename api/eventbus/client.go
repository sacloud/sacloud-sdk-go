// Copyright 2025-2026 The sacloud/eventbus-api-go authors
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

package eventbus

import (
	"context"
	"fmt"
	"runtime"

	v1 "github.com/sacloud/eventbus-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1a/api/cloud/1.1"

	// ServiceKey SDKの種別を示すキー、プロファイルでのエンドポイント取得に利用する
	ServiceKey = "eventbus"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"eventbus-api-go/%s (%s/%s; +https://github.com/sacloud/eventbus-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

// DummySecuritySource SecuritySourceはOpenAPI定義で使用されている認証のための仕組み。saclient-goが処理するので、ogen用はダミーで誤魔化す
type DummySecuritySource struct {
	Token string
}

func (ss DummySecuritySource) ApiKeyAuth(ctx context.Context, operationName v1.OperationName) (v1.ApiKeyAuth, error) {
	return v1.ApiKeyAuth{Username: ss.Token}, nil
}

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
	middleware, err := injectFilterMiddleware(apiRootURL)
	if err != nil {
		return nil, NewError("failed to load middleware", nil)
	}

	augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
		saclient.WithForceAutomaticAuthentication(),
		saclient.WithBigInt(false),          // 文字列を勝手に数値に変換しないようヘッダーで指定
		saclient.WithMiddleware(middleware), // TODO: filterがOpenAPI定義で表現できるようになったら不要となる。その後に削除する。
	)
	if err != nil {
		return nil, err
	}
	return v1.NewClient(apiRootURL, DummySecuritySource{Token: "eventbus-client"}, v1.WithClient(augmented))
}
