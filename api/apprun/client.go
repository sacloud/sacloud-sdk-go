// Copyright 2021-2024 The sacloud/apprun-api-go authors
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

package apprun

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	client "github.com/sacloud/api-client-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api"
	serviceKey        = "apprun"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"apprun-api-go/%s (%s/%s; +https://github.com/sacloud/apprun-api-go) %s",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
	client.DefaultUserAgent,
)

// Client APIクライアント
type Client struct {
	// Profile usacloud互換プロファイル名
	//
	// Saclientフィールドが指定されている場合は無視される
	Profile string

	// Token APIキー: トークン
	//
	// Saclientフィールドが指定されている場合は無視される
	Token string
	// Token APIキー: シークレット
	//
	// Saclientフィールドが指定されている場合は無視される
	Secret string //nolint:gosec

	// APIRootURL APIのリクエスト先URLプレフィックス、省略可能
	//
	// Saclientフィールドが指定されている場合は無視される
	APIRootURL string

	// Options HTTPクライアント関連オプション
	//
	// Saclientフィールドが指定されている場合は無視される
	Options *client.Options

	// DisableProfile usacloud互換プロファイルからの設定読み取りを無効化
	//
	// Saclientフィールドが指定されている場合は無視される
	DisableProfile bool
	// DisableEnv 環境変数からの設定読み取りを無効化
	//
	// Saclientフィールドが指定されている場合は無視される
	DisableEnv bool

	// Saclient APIクライアント
	//
	// 従来のapi-cient-goを置き換えるもので、通常は*saclient.Clientを呼び出し側で組み立ててから渡すことを想定している。
	// 互換性維持のためにこの値が空の場合はClientの残りのフィールドから組み立てられる。
	Saclient saclient.ClientAPI

	initOnce sync.Once
}

func (c *Client) serverURL() string {
	v := DefaultAPIRootURL
	if c.APIRootURL != "" {
		v = c.APIRootURL
	}

	if !strings.HasSuffix(v, "/") {
		v += "/"
	}
	return v
}

func (c *Client) init() error {
	var initError error
	c.initOnce.Do(func() {
		var opts []*client.Options
		// 1: Profile
		if !c.DisableProfile {
			o, err := client.OptionsFromProfile(c.Profile)
			if err != nil {
				initError = err
				return
			}
			opts = append(opts, o)
		}

		// 2: Env
		if !c.DisableEnv {
			opts = append(opts, client.OptionsFromEnv())
		}

		// 3: UserAgent
		opts = append(opts, &client.Options{
			UserAgent: UserAgent,
		})

		// 4: Options
		if c.Options != nil {
			opts = append(opts, c.Options)
		}

		// 5: フィールドのAPIキー
		opts = append(opts, &client.Options{
			AccessToken:       c.Token,
			AccessTokenSecret: c.Secret,
		})

		if c.Saclient == nil {
			c.Saclient = saclient.NewFactory(opts...)

			endpointConfig, err := c.Saclient.EndpointConfig()
			if err != nil {
				initError = err
				return
			}
			// エンドポイント設定にapprunのエンドポイントがあれば上書きする
			if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
				c.APIRootURL = ep
			}
		}
	})
	return initError
}

func (c *Client) apiClient() (*v1.ClientWithResponses, error) {
	if err := c.init(); err != nil {
		return nil, err
	}

	return &v1.ClientWithResponses{
		ClientInterface: &v1.Client{
			Server: c.serverURL(),
			Client: c.Saclient,
		},
	}, nil
}
