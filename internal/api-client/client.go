// Copyright 2022-2025 The sacloud/api-client-go Authors
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

package client

import (
	"fmt"
	"net/http"
	"runtime"
	"strings"
	"sync"

	sacloudhttp "github.com/sacloud/go-http"
)

// DefaultUserAgent デフォルトのユーザーエージェント
var DefaultUserAgent = fmt.Sprintf(
	"api-client-go/v%s (%s/%s; +https://github.com/sacloud/api-client-go) %s",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
	sacloudhttp.DefaultUserAgent,
)

// Client APIクライアント
type Client struct {
	// APIRootURL APIのリクエスト先URLプレフィックス
	APIRootURL string

	initOnce sync.Once
	factory  *Factory
}

func (c *Client) ServerURL() string {
	v := c.APIRootURL
	if !strings.HasSuffix(v, "/") {
		v += "/"
	}
	return v
}

func (c *Client) NewHttpRequestDoer() HttpRequestDoer {
	return c.factory.NewHttpRequestDoer()
}

func (c *Client) Options() *Options {
	return c.factory.Options()
}

func (c *Client) init(params *ClientParams) error {
	var initError error
	c.initOnce.Do(func() {
		var opts []*Options
		// 1: Profile
		if !params.DisableProfile {
			o, err := OptionsFromProfile(params.Profile)
			if err != nil {
				initError = err
				return
			}
			opts = append(opts, o)
		}

		// 2: Env
		if !params.DisableEnv {
			opts = append(opts, OptionsFromEnv())
		}

		// 3: UserAgent
		opts = append(opts, &Options{
			UserAgent: params.UserAgent,
		})

		// 4: Options
		if params.Options != nil {
			opts = append(opts, params.Options)
		}

		// 5: フィールドのAPIキー
		opts = append(opts, &Options{
			AccessToken:       params.Token,
			AccessTokenSecret: params.Secret,
		})

		// 6: Custom HTTP client (for httptest)
		if params.HTTPClient != nil {
			opts = append(opts, &Options{
				HttpClient: params.HTTPClient,
			})
		}

		c.factory = NewFactory(opts...)
	})
	return initError
}

// SDKライブラリから設定するパラメータ。WithXXXを使って特定のパラメータだけ設定可能
type ClientParams struct {
	// APIのリクエスト先URLプレフィックス
	APIRootURL string
	// APIキー群
	Token  string
	Secret string
	// クライアントから送られるユーザエージェント
	UserAgent string
	// Options HTTPクライアント関連オプション
	Options *Options
	// Profile usacloud互換プロファイル名
	Profile string
	// usacloud互換プロファイルからの設定読み取りを無効化
	DisableProfile bool
	// 環境変数からの設定読み取りを無効化
	DisableEnv bool
	// カスタムHTTPクライアント (例: httptest.NewServer().Client())
	HTTPClient *http.Client
	// SDKライブラリから追加したいパラメータがあったら随時追加
}

type ClientParam func(*ClientParams)

func WithUserAgent(ua string) ClientParam {
	return func(params *ClientParams) {
		params.UserAgent = ua
	}
}

func WithApiKeys(accessToken string, secretToken string) ClientParam {
	return func(params *ClientParams) {
		params.Token = accessToken
		params.Secret = secretToken
	}
}

func WithProfile(profile string) ClientParam {
	return func(params *ClientParams) {
		params.Profile = profile
	}
}

func WithDisableProfile(disable bool) ClientParam {
	return func(params *ClientParams) {
		params.DisableProfile = disable
	}
}

func WithDisableEnv(disable bool) ClientParam {
	return func(params *ClientParams) {
		params.DisableEnv = disable
	}
}

func WithOptions(options *Options) ClientParam {
	return func(params *ClientParams) {
		params.Options = options
	}
}

// WithHTTPClient allows setting a custom http.Client (e.g., from httptest)
func WithHTTPClient(client *http.Client) ClientParam {
	return func(params *ClientParams) {
		params.HTTPClient = client
	}
}

// Clientを初期化してから返す。WithXXXを使って特定の設定を初期化可能
func NewClient(apiUrl string, params ...ClientParam) (*Client, error) {
	clientParams := &ClientParams{
		APIRootURL: apiUrl,
		UserAgent:  DefaultUserAgent,
	}

	for _, param := range params {
		param(clientParams)
	}

	return NewClientWithParams(clientParams)
}

// 設定したいパラメータが多い場合は直接ClientParamsを初期化して渡す
func NewClientWithParams(params *ClientParams) (*Client, error) {
	client := &Client{
		APIRootURL: params.APIRootURL,
	}
	if err := client.init(params); err != nil {
		return nil, err
	}

	return client, nil
}
