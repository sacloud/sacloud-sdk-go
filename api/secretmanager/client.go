// Copyright 2025- The sacloud/secretmanager-api-go authors
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

package secretmanager

import (
	"context"
	"fmt"
	"runtime"

	"github.com/sacloud/saclient-go"
	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/tk1a/api/cloud/1.1"
	serviceKey        = "secretmanager"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"secretmanager-api-go/%s (%s/%s; +https://github.com/sacloud/secretmanager-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

// SecuritySourceはOpenAPI定義で使用されている認証のための仕組み。api-client-goが処理するので、ogen用はダミーで誤魔化す
type DummySecuritySource struct {
	Username string
	Password string
}

func (ss DummySecuritySource) BasicAuth(ctx context.Context, operationName v1.OperationName) (v1.BasicAuth, error) {
	return v1.BasicAuth{Username: ss.Username, Password: ss.Password, Roles: nil}, nil
}

func NewClient(sa saclient.ClientAPI) (*v1.Client, error) {
	endpointConfig, err := sa.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load endpoint configuration", err)
	}

	resetter := sa.(saclient.ClientOptionAPI)

	client, err := resetter.DupWith(
		saclient.WithUserAgent(UserAgent),
		// これはなにか:
		// DummySecuritySource.BasicAuth()がBasic認証を生成
		// しかし実際の通信で必ずしもBasic認証が使われると限らない
		//　そのあたりをsaclient-go側で吸収させる設定が下記↓
		saclient.WithForceAutomaticAuthentication(),
	)
	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	apiUrl := DefaultAPIRootURL
	if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
		apiUrl = ep
	}

	v1Client, err := v1.NewClient(apiUrl, DummySecuritySource{}, v1.WithClient(client))
	if err != nil {
		return nil, NewError("NewClientWithApiUrl", err)
	}

	return v1Client, nil
}
