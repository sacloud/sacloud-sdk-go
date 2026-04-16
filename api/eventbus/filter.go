// Copyright 2025-2026 The sacloud/eventbus-api-go authors
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

package eventbus

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	v1 "github.com/sacloud/eventbus-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

func injectFilterMiddleware(apiRootURL string) (saclient.Middleware, error) {
	u, err := url.Parse(apiRootURL)
	if err != nil {
		return nil, err
	}
	listAPIPath := u.JoinPath("/commonserviceitem").Path

	return func(req *http.Request, pull func() (saclient.Middleware, bool)) (*http.Response, error) {
		injectFilterToRequest(req, listAPIPath)
		cont, ok := pull()
		if !ok {
			return nil, errors.New("middleware not found error")
		}
		return cont(req, pull)
	}, nil
}

func injectFilterToRequest(req *http.Request, listAPIPath string) {
	// NOTE: OpenAPIで表現できないクエリの書き込みを行う
	// 同じエンドポイントに3種類のProvider.Classでフィルタしたいため、生成コードの書き換えでなくclient middlewareで対応
	// `GET /commonserviceitem?{"Filter":{"Provider.Class":"eventbusschedule"}}`
	// `GET /commonserviceitem?{"Filter":{"Provider.Class":"eventbustrigger"}}`
	// `GET /commonserviceitem?{"Filter":{"Provider.Class":"eventbusprocessconfiguration"}}`.
	if req.Method == http.MethodGet && req.URL.Path == listAPIPath {
		pc := getFilterProviderClass(req.Context())
		req.URL.RawQuery = fmt.Sprintf(`{"Filter":{"Provider.Class":"%s"}}`, pc)
	}
}

type ctxKeyFilterProviderClass struct{}

func setFilterProviderClass(ctx context.Context, v v1.ProviderClass) context.Context {
	return context.WithValue(ctx, ctxKeyFilterProviderClass{}, v)
}

func getFilterProviderClass(ctx context.Context) v1.ProviderClass {
	return ctx.Value(ctxKeyFilterProviderClass{}).(v1.ProviderClass)
}
