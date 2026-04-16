// Copyright 2022-2026 The object-storage-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package objectstorage

import (
	"context"
	"fmt"
	"net/url"
	"runtime"

	v2 "github.com/sacloud/object-storage-api-go/apis/v2"
	"github.com/sacloud/saclient-go"
)

const (
	DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/zone/is1a/api/objectstorage/1.0/"
	serviceKey        = "object_storage"
)

var NewUserAgent = fmt.Sprintf(
	"object-storage-api-go/%s (%s/%s; +https://github.com/sacloud/object-storage-api-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

type dummySecuritySource struct{}

func (ss dummySecuritySource) BasicAuth(ctx context.Context, operationName v2.OperationName) (v2.BasicAuth, error) {
	return v2.BasicAuth{Username: "", Password: "", Roles: nil}, nil
}

type FedClient struct {
	client *v2.Client
}

func NewFedClient(client saclient.ClientAPI) (*FedClient, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load message endpoint configuration", err)
	}

	endpoint := DefaultAPIRootURL
	if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
		endpoint = ep
	}

	return NewFedClientWithAPIRootURL(client, endpoint)
}

func NewFedClientWithAPIRootURL(client saclient.ClientAPI, apiRootURL string) (*FedClient, error) {
	dupable, ok := client.(saclient.ClientOptionAPI)
	if !ok {
		return nil, NewError("client does not implement saclient.ClientOptionAPI", nil)
	}
	augmented, err := dupable.DupWith(
		saclient.WithUserAgent(NewUserAgent),
		saclient.WithBigInt(false),
		saclient.WithForceAutomaticAuthentication(),
	)
	if err != nil {
		return nil, err
	}
	u, err := url.JoinPath(apiRootURL, "fed", "v1")
	if err != nil {
		return nil, err
	}
	c, err := v2.NewClient(u, &dummySecuritySource{}, v2.WithClient(augmented))
	if err != nil {
		return nil, err
	}
	return &FedClient{client: c}, nil
}

type SiteClient struct {
	client *v2.Client
}

func NewSiteClient(client saclient.ClientAPI, siteId string) (*SiteClient, error) {
	endpointConfig, err := client.EndpointConfig()
	if err != nil {
		return nil, NewError("unable to load message endpoint configuration", err)
	}

	endpoint := DefaultAPIRootURL
	if ep, ok := endpointConfig.Endpoints[serviceKey]; ok && ep != "" {
		endpoint = ep
	}

	return NewSiteClientWithAPIRootURL(client, endpoint, siteId)
}

func NewSiteClientWithAPIRootURL(client saclient.ClientAPI, apiRootURL string, siteId string) (*SiteClient, error) {
	dupable, ok := client.(saclient.ClientOptionAPI)
	if !ok {
		return nil, NewError("client does not implement saclient.ClientOptionAPI", nil)
	}
	argumented, err := dupable.DupWith(
		saclient.WithUserAgent(NewUserAgent),
		saclient.WithBigInt(false),
		saclient.WithForceAutomaticAuthentication(),
	)
	if err != nil {
		return nil, err
	}
	u, err := url.JoinPath(apiRootURL, siteId, "v2")
	if err != nil {
		return nil, err
	}
	c, err := v2.NewClient(u, &dummySecuritySource{}, v2.WithClient(argumented))
	if err != nil {
		return nil, err
	}
	return &SiteClient{client: c}, nil
}
