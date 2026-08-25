// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite

import (
	"fmt"
	"runtime"
	"strings"

	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

const (
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1c/api/cloud/1.1/"
	DefaultEndpoint   = "https://secure.sakura.ad.jp/cloud/zone/"
	DefaultZone       = "is1c"
	serviceKey        = "networking_suite"
)

var UserAgent = fmt.Sprintf(
	"networking-suite-api-go/%s (%s/%s; +https://github.com/sacloud/sacloud-sdk-go/api/networking-suite)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
	cfg, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("NewClient", err)
	}

	apiUrl := DefaultAPIRootURL
	// 他のクライアントと挙動をあわせるため、エンドポイントの設定がある場合はそちらを最優先とする
	if ep, ok := cfg.Endpoints[serviceKey]; ok && ep != "" {
		apiUrl = ep
	} else if cfg.Zone != "" {
		const path = "api/cloud/1.1/"
		apiUrl = fmt.Sprintf("%s/%s/%s", strings.TrimSuffix(DefaultEndpoint, "/"), strings.TrimPrefix(strings.TrimSuffix(cfg.Zone, "/"), "/"), path)
	}
	return NewClientWithAPIRootURL(client, apiUrl)
}

func NewClientWithAPIRootURL(client saclient.ClientAPI, apiRootURL string) (*v1.Client, error) {
	if dupable, ok := client.(saclient.ClientOptionAPI); !ok {
		return nil, NewError("client does not implement saclient.ClientOptionAPI", nil)
	} else if augmented, err := dupable.DupWith(
		saclient.WithUserAgent(UserAgent),
		saclient.WithForceAutomaticAuthentication(),
	); err != nil {
		return nil, err
	} else {
		// Use lower retry settings until stable API release
		_ = augmented.SetEnviron([]string{
			"SAKURA_RETRY_MAX=3",
			"SAKURA_RETRY_WAIT_MAX=7",
			"SAKURA_RETRY_WAIT_MIN=3",
		})
		err = augmented.Populate()
		if err != nil {
			return nil, NewError("failed to populate client", err)
		}
		return v1.NewClient(apiRootURL, v1.WithClient(augmented))
	}
}
