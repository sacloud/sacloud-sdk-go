// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol

import (
	"context"
	"fmt"
	"runtime"

	"github.com/sacloud/saclient-go"
	v1 "github.com/sacloud/security-control-api-go/apis/v1"
)

const (
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/tk1b/api/securitycontrol/1.0/"
	serviceKey        = "security_control"
)

var UserAgent = fmt.Sprintf(
	"security-control-api-go/%s (%s/%s; +https://github.com/sacloud/security-control-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

type dummySecuritySource struct{}

func (ss dummySecuritySource) BasicAuth(ctx context.Context, operationName v1.OperationName) (v1.BasicAuth, error) {
	return v1.BasicAuth{Username: "", Password: "", Roles: nil}, nil
}
func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load message endpoint configuration", err)
	}

	endpoint := DefaultAPIRootURL
	if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
		endpoint = ep
	}

	return NewClientWithAPIRootURL(client, endpoint)
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
		return v1.NewClient(apiRootURL, dummySecuritySource{}, v1.WithClient(augmented))
	}
}
