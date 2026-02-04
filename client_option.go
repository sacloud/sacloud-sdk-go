// Copyright 2025- The sacloud/saclient-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package saclient

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"slices"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

type clientOption func(*Client) error

type ClientOptionAPI interface {
	// destructive application of argument options
	SetWith(...clientOption) error

	// duplicative application of argument options;
	// use it when the receiver is already populated (or unknown)
	DupWith(...clientOption) (ClientAPI, error)
}

//nolint:gocritic
func (c *Client) SetWith(opts ...clientOption) error {
	if c == nil {
		return NewErrorf("nil client")
	} else if c.once.Done() {
		return NewErrorf("client already populated; cannot change settings")
	} else {
		q := slices.Values(opts)
		w := transformSeq(q, func(opt clientOption) (error, bool) { return opt(c), true })
		e := slices.Collect(w)
		return errors.Join(e...)
	}
}

func (c *Client) DupWith(opts ...clientOption) (ClientAPI, error) {
	ret := c.Dup().(*Client)
	err := ret.SetWith(opts...)
	return ret, err
}

func WithRootURL(url string) clientOption {
	return func(c *Client) error {
		c.params.dynamic.apiRootURL.initialize(url)
		return nil
	}
}

func WithTestServer(svr *httptest.Server) clientOption {
	return func(c *Client) error {
		c.params.dynamic.mockServer.initialize(svr)
		return nil
	}
}

// this is not strictly necessary because you can set it via env/flag/HCL,
// but can be handy on occasions.
func WithTraceMode(mode string) clientOption {
	return func(c *Client) error {
		c.params.dynamic.traceMode.initialize(mode)
		return nil
	}
}

func WithFavouringRFC7617() clientOption {
	return func(c *Client) error {
		c.params.dynamic.authPreference.initialize("basic")
		return nil
	}
}

func WithFavouringRFC7523() clientOption {
	return func(c *Client) error {
		c.params.dynamic.authPreference.initialize("bearer")
		return nil
	}
}

// Did you know...? These days the "Authorization:" headers are for authentication
// https://datatracker.ietf.org/doc/html/rfc7235#section-4.2
// ... which sounds quite confusing to me honestly though.

func WithFavouringBasicAuthentication() clientOption  { return WithFavouringRFC7617() } // alias
func WithFavouringBearerAuthentication() clientOption { return WithFavouringRFC7523() } // alias

func WithMiddleware(m ...Middleware) clientOption {
	return func(c *Client) error {
		// This option is cumulative, must merge
		if cur, ok := c.params.dynamic.middlewares.Get(); ok {
			m = append(m, cur...) // later ones have higher priority
		}
		c.params.dynamic.middlewares.initialize(m)
		return nil
	}
}

func WithBearerToken(bearer string) clientOption {
	return WithMiddleware(func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		req.Header.Set("Authorization", "Bearer "+bearer)

		return pullThenCall(pull, req)
	})
}

func WithBasicAuth(user, pass string) clientOption {
	return WithMiddleware(func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		req.SetBasicAuth(user, pass)

		return pullThenCall(pull, req)
	})
}

func WithBasicAuth1(userinfo string) clientOption {
	return WithMiddleware(func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		req.Header.Set("Authorization", "Basic "+userinfo)

		return pullThenCall(pull, req)
	})
}

func WithForceAutomaticAuthentication() clientOption {
	return WithMiddleware(func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		req.Header.Del("Authorization")

		return pullThenCall(pull, req)
	})
}

func WithBigInt(needed bool) clientOption {
	return WithMiddleware(func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		var val string
		if needed {
			val = "1"
		} else {
			val = "0"
		}
		req.Header.Set("X-Sakura-Bigint-As-Int", val)

		return pullThenCall(pull, req)
	})
}

func WithUserAgent(ua string) clientOption {
	return WithMiddleware(func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		req.Header.Set("User-Agent", ua)

		return pullThenCall(pull, req)
	})
}

func WithCheckRetryFunc(f retryablehttp.CheckRetry) clientOption {
	return func(c *Client) error {
		c.params.dynamic.checkRetryFunc.initialize(f)
		return nil
	}
}

// disables retries at all
func WithoutRetry() clientOption {
	return WithCheckRetryFunc(disableRetry)
}

// WithDefaultTimeout sets default timeout for requests.
// For each individual request, timeout can also be set using context.
//
// ```golang
//
//	var client saclient.Client
//
//	// This is setting default
//	client.SetEnviron([]string{"SAKURA_API_REQUEST_TIMEOUT=30"})
//
//	// Alternatively, set it via option
//	client.SetWith(WithDefaultTimeout(30 * time.Second))
//
//	// Below sets timeout only for this request
//	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//	defer cancel()
//	req, _ := http.NewRequestWithContext(ctx, "GET", "https://example.com", nil)
//	resp, err := client.Do(req)
//
// ```
func WithDefaultTimeout(t time.Duration) clientOption {
	return func(c *Client) error {
		c.params.dynamic.apiRequestTimeout.initialize(t)
		return nil
	}
}

var disableRetry = func(context.Context, *http.Response, error) (bool, error) { return false, nil }
var _ ClientOptionAPI = (*Client)(nil)
