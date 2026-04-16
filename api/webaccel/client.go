// Copyright 2021-2022 The webaccel-api-go authors
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

package webaccel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/saclient-go"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1a/api/webaccel/1.0/"
	// serviceKey エンドポイントURLの取得に使用するサービスキー
	serviceKey = "webaccel"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"webaccel-api-go/v%s (%s/%s; +https://github.com/sacloud/webaccel-api-go) %s",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
	client.DefaultUserAgent,
)

// Client APIクライアント
type Client struct {
	// Profile usacloud互換プロファイル名
	Profile string

	// APIRootURL APIのリクエスト先URLプレフィックス、省略可能
	APIRootURL string

	// AccessToken APIキー:トークン
	// Optionsでの指定より優先される
	AccessToken string
	// AccessTokenSecret APIキー:シークレット
	// Optionsでの指定より優先される
	AccessTokenSecret string

	// Options HTTPクライアント関連オプション
	Options *client.Options

	// DisableProfile usacloud互換プロファイルからの設定読み取りを無効化
	DisableProfile bool

	// DisableEnv 環境変数からの設定読み取りを無効化
	DisableEnv bool

	Saclient saclient.ClientAPI

	initOnce sync.Once

	//onceで実行するInitializeのエラーを保持する
	initError error
}

func (c *Client) RootURL() string {
	// 初期化処理を実行して、エンドポイントURLを取得する
	// エラーはDo()でチェックするため、ここではエラーは返さない
	_ = c.init()

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
	c.initOnce.Do(func() {
		var opts []*client.Options
		// 1: Profile
		if !c.DisableProfile {
			o, err := client.OptionsFromProfile(c.Profile)
			if err != nil {
				c.initError = err
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
			AccessToken:       c.AccessToken,
			AccessTokenSecret: c.AccessTokenSecret,
		})

		if c.Saclient == nil {
			c.Saclient = saclient.NewFactory(opts...)
		}
		// エンドポイントURLの取得
		endpointConfig, err := c.Saclient.EndpointConfig()
		if err != nil {
			c.initError = err
		}
		if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
			c.APIRootURL = ep
		}
	})
	return c.initError
}

// Do APIコール実施
func (c *Client) Do(ctx context.Context, method, uri string, body interface{}) ([]byte, error) {
	if err := c.init(); err != nil {
		return nil, err
	}

	req, err := c.newRequest(ctx, method, uri, body)
	if err != nil {
		return nil, err
	}

	// API call
	resp, err := c.Saclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if !c.isOkStatus(resp.StatusCode) {
		errResponse := &APIErrorResponse{}
		err := json.Unmarshal(data, errResponse)
		if err != nil {
			return nil, fmt.Errorf("error in response: %s", string(data))
		}
		return nil, NewAPIError(req.Method, req.URL, resp.StatusCode, errResponse)
	}

	return data, nil
}

func (c *Client) newRequest(ctx context.Context, method, uri string, body interface{}) (*http.Request, error) {
	// setup url and body
	var url = uri
	var bodyReader io.ReadSeeker
	if body != nil {
		var bodyJSON []byte
		bodyJSON, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		if method == "GET" {
			url = fmt.Sprintf("%s?%s", url, bytes.NewBuffer(bodyJSON))
		} else {
			bodyReader = bytes.NewReader(bodyJSON)
		}
	}
	return http.NewRequestWithContext(ctx, method, url, bodyReader)
}

func (c *Client) isOkStatus(code int) bool {
	codes := map[int]bool{
		http.StatusOK:        true,
		http.StatusCreated:   true,
		http.StatusAccepted:  true,
		http.StatusNoContent: true,
	}
	_, ok := codes[code]
	return ok
}
