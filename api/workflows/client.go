// Copyright 2025- The sacloud/workflows-api-go Authors
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

package workflows

import (
	"context"
	"fmt"
	"runtime"

	"github.com/sacloud/saclient-go"
	v1 "github.com/sacloud/workflows-api-go/apis/v1"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/tk1b/api/workflow/1.0/"
	// serviceKey APIエンドポイントのキー
	serviceKey = "workflows"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"workflows-api-go/%s (%s/%s; +https://github.com/sacloud/workflows-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

// voidSecuritySource is a placeholder to satisfy the SecuritySource interface.
// saclientにて処理するためここにはロジック不要だが何か渡さないといけないので空の構造体を用意する
type voidSecuritySource struct{}

func (voidSecuritySource) ApiKeyAuth(context.Context, v1.OperationName) (v1.ApiKeyAuth, error) {
	return v1.ApiKeyAuth{}, nil
}

// NewClient creates a new workflows API client with default settings
func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load endpoint configuration", err)
	}
	apiURL := DefaultAPIRootURL
	if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
		apiURL = ep
	}
	return NewClientWithAPIRootURL(client, apiURL)
}

// NewClientWithAPIRootURL creates a new workflows API client with a custom API root URL
func NewClientWithAPIRootURL(client saclient.ClientAPI, apiRootURL string) (*v1.Client, error) {
	if dupable, ok := client.(saclient.ClientOptionAPI); !ok {
		return nil, NewError("client does not implement saclient.ClientOptionAPI", nil)
	} else if augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
		saclient.WithForceAutomaticAuthentication(),
	); err != nil {
		return nil, err
	} else {
		return v1.NewClient(apiRootURL, voidSecuritySource{}, v1.WithClient(augmented))
	}
}
