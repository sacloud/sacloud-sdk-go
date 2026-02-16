// Copyright 2025- The sacloud/kms-api-go authors
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

package kms

import (
	"context"
	"fmt"
	"runtime"

	v1 "github.com/sacloud/kms-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

// DefaultAPIRootURL デフォルトのAPIルートURL
const (
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/tk1a/api/cloud/1.1"
	ServiceKey        = "kms"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"kms-api-go/%s (%s/%s; +https://github.com/sacloud/kms-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

// SecuritySourceはOpenAPI定義で使用されている認証のための仕組み。ogen用はダミーで誤魔化す
type dummySecuritySource struct {
	Username string
	Password string
}

func (ss dummySecuritySource) BasicAuth(ctx context.Context, operationName v1.OperationName) (v1.BasicAuth, error) {
	return v1.BasicAuth{Username: ss.Username, Password: ss.Password, Roles: nil}, nil
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
	argumented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),

		saclient.WithForceAutomaticAuthentication(),
	)
	if err != nil {
		return nil, err
	}
	c, err := v1.NewClient(apiRootURL, dummySecuritySource{}, v1.WithClient(argumented))
	if err != nil {
		return nil, err
	}
	return c, nil
}
