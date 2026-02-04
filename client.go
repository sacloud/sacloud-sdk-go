// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package apprun_dedicated

import (
	"context"
	"fmt"
	"runtime"

	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/common"
	"github.com/sacloud/saclient-go"
)

// DefaultAPIRootURL デフォルトのAPIルートURL
const DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/api/apprun-dedicated/1.0/"

// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
	"apprun-dedicated-api-go/%s (%s/%s; +https://github.com/sacloud/apprun-dedicated-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

// voidSecuritySource is a placeholder to satisfy the SecuritySource interface.
// saclientにて処理するためここにはロジック不要だが何か渡さないといけないので空の構造体を用意する
type voidSecuritySource struct{}

func (voidSecuritySource) BasicAuth(context.Context, v1.OperationName) (v1.BasicAuth, error) {
	return v1.BasicAuth{}, nil
}

func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
	return NewClientWithAPIRootURL(client, DefaultAPIRootURL)
}

func NewClientWithAPIRootURL(client saclient.ClientAPI, apiRootURL string) (*v1.Client, error) {
	dupable, ok := client.(saclient.ClientOptionAPI)

	if !ok {
		return nil, common.NewError("client does not implement saclient.ClientOptionAPI", nil)
	}

	augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
		// これはなにか:
		// voidSecuritySource.BasicAuth()がBasic認証を生成
		// しかし実際の通信で必ずしもBasic認証が使われると限らない
		//　そのあたりをsaclient-go側で吸収させる設定が下記↓
		saclient.WithForceAutomaticAuthentication(),
	)

	if err != nil {
		return nil, err
	}

	return v1.NewClient(apiRootURL, voidSecuritySource{}, v1.WithClient(augmented))
}
