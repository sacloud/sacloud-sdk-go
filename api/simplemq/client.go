// Copyright 2025- The sacloud/simplemq-api-go authors
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

package simplemq

import (
	"context"
	"fmt"
	"runtime"

	"github.com/sacloud/saclient-go"
	"github.com/sacloud/simplemq-api-go/apis/v1/message"
	"github.com/sacloud/simplemq-api-go/apis/v1/queue"
)

const (
	// DefaultQueueAPIRootURL デフォルトのQueue APIルートURL
	DefaultQueueAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1a/api/cloud/1.1"
	// ServiceKeyQueue SDKの種別を示すキー、プロファイルでのエンドポイント取得に利用するもので、Queue APIのエンドポイント切り替えに利用する
	ServiceKeyQueue = "simple_mq_queue"

	// DefaultMessageAPIRootURL デフォルトのMessage APIルートURL
	DefaultMessageAPIRootURL = "https://simplemq.tk1b.api.sacloud.jp"
	// ServiceKeyMessage SDKの種別を示すキー、プロファイルでのエンドポイント取得に利用するもので、Message APIのエンドポイント切り替えに利用する
	ServiceKeyMessage = "simple_mq_message"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"simplemq-api-go/%s (%s/%s; +https://github.com/sacloud/simplemq-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

// SecuritySourceはOpenAPI定義で使用されている認証のための仕組み。saclient-goが処理するので、ogen用はダミーで誤魔化す
type DummySecuritySource struct {
	Token string
}

func (ss DummySecuritySource) ApiKeyAuth(ctx context.Context, operationName queue.OperationName) (queue.ApiKeyAuth, error) {
	return queue.ApiKeyAuth{Username: ss.Token}, nil
}

func NewQueueClient(client saclient.ClientAPI) (*queue.Client, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load queue endpoint configuration", err)
	}
	endpoint := DefaultQueueAPIRootURL
	if ep, ok := endpointConfig.Endpoints[ServiceKeyQueue]; ok && ep != "" {
		endpoint = ep
	}

	return NewQueueClientWithApiUrl(endpoint, client)
}

func NewQueueClientWithApiUrl(apiUrl string, client saclient.ClientAPI) (*queue.Client, error) {
	var ss DummySecuritySource

	if dupable, ok := client.(saclient.ClientOptionAPI); !ok {
		return nil, NewError("client does not implement saclient.ClientOptionAPI", nil)
	} else if augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
		// これはなにか:
		// DummySecuritySource.ApiKeyAuth()がBasic認証を生成
		// しかし実際の通信で必ずしもBasic認証が使われると限らない
		//　そのあたりをsaclient-go側で吸収させる設定が下記↓
		saclient.WithForceAutomaticAuthentication(),
	); err != nil {
		return nil, err
	} else {
		return queue.NewClient(apiUrl, ss, queue.WithClient(augmented))
	}
}

type ApiKeySecuritySource struct {
	Token string
}

func (ss ApiKeySecuritySource) ApiKeyAuth(ctx context.Context, operationName message.OperationName) (message.ApiKeyAuth, error) {
	return message.ApiKeyAuth{Token: ss.Token}, nil
}

func NewMessageClient(apiKey string, client saclient.ClientAPI) (*message.Client, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load message endpoint configuration", err)
	}
	endpoint := DefaultMessageAPIRootURL
	if ep, ok := endpointConfig.Endpoints[ServiceKeyMessage]; ok && ep != "" {
		endpoint = ep
	}

	return NewMessageClientWithApiUrl(endpoint, apiKey, client)
}

func NewMessageClientWithApiUrl(apiUrl, apiKey string, client saclient.ClientAPI) (*message.Client, error) {
	var ss ApiKeySecuritySource

	if dupable, ok := client.(saclient.ClientOptionAPI); !ok {
		return nil, NewError("client does not implement saclient.ClientOptionAPI", nil)
	} else if augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
		// これはなにか:
		// ApiKeySecuritySource.ApiKeyAuth()がBasic認証を生成
		// しかし実際の通信で必要なのはBearer認証
		//　そのあたりをsaclient-go側で吸収させる設定が下記↓
		saclient.WithBearerToken(apiKey),
	); err != nil {
		return nil, err
	} else {
		return message.NewClient(apiUrl, ss, message.WithClient(augmented))
	}
}
