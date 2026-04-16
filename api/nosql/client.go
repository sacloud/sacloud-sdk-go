// Copyright 2016-2025 The terraform-provider-sakura Authors
// SPDX-License-Identifier: Apache-2.0

package nosql

import (
	"fmt"
	"runtime"

	v1 "github.com/sacloud/nosql-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

const (
	// DefaultAPIRootURL デフォルトのAPIルートURL
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/tk1b/api/cloud/1.1"

	// ServiceKey SDKの種別を示すキー、プロファイルでのエンドポイント取得に利用する
	ServiceKey = "nosql"
)

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"nosql-api-go/%s (%s/%s; +https://github.com/sacloud/nosql-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

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

	augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
		saclient.WithBigInt(false), // 文字列を勝手に数値に変換しないようヘッダーで指定
	)
	if err != nil {
		return nil, err
	}
	return v1.NewClient(apiRootURL, v1.WithClient(augmented))
}
