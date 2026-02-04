// Copyright 2022-2025 The sacloud/api-client-go Authors
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

package saclient

import (
	"iter"
	"net/http"
	"net/http/httptest"
	"slices"
	"time"

	saht "github.com/sacloud/go-http"
	"go.uber.org/ratelimit"
)

type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

// A Middleware is much like http.RoundTripper, but it additionally takes a "pull" function.
// You might chain middlewares by pulling the next one.
type Middleware func(*http.Request, func() (Middleware, bool)) (*http.Response, error)

type doer struct {
	// configurable for tests (net/http/httptest)
	client      *http.Client
	root        string
	server      *httptest.Server
	rateLimiter ratelimit.Limiter
	middlewares []Middleware
}

var _ HttpRequestDoer = (*doer)(nil)

func (d *doer) Do(req *http.Request) (*http.Response, error) {
	pull, stop := iter.Pull(slices.Values(d.middlewares))
	defer stop()

	return pullThenCall(pull, req)
}

func newHttpRequestDoer(c *config) (HttpRequestDoer, error) {
	var d doer
	h := http.Client{
		Transport: http.DefaultTransport.(*http.Transport).Clone(),
	}

	if result := obtainFromConfig[time.Duration](c, "APIRequestTimeout"); result.isErr() {
		return nil, result.error()
	} else if v, ok := result.some(); ok {
		h.Timeout = v
	} else {
		// UNLIKELY: APIRequestTimeout has default value
		h.Timeout = 300 * time.Second
	}

	if result := obtainFromConfig[*httptest.Server](c, "MockServer"); result.isErr() {
		return nil, result.error()
	} else if svr, ok := result.some(); ok {
		d.client = svr.Client()
		d.root = svr.URL
		d.server = svr
	} else if result := obtainFromConfig[string](c, "APIRootURL"); result.isErr() {
		return nil, result.error()
	} else if apiRootURL, ok := result.some(); ok {
		d.client = &h
		d.root = apiRootURL
	} else {
		d.client = &h
	}
	// OK when root is absent

	if result := obtainFromConfig[int64](c, "APIRequestRateLimit"); result.isErr() {
		return nil, result.error()
	} else if v, ok := result.some(); ok {
		d.rateLimiter = ratelimit.New(int(v))
	} else {
		d.rateLimiter = ratelimit.NewUnlimited()
	}

	// basic middlewares
	// note that they are called in order
	middlewares := []Middleware{
		// upper layer vvvvv
		c.middlewareSetHeader(),
		d.middlewareAuthorization(c),
		d.tracer(c),
		gzipExpander,
		d.middlewareRateLimitter(),
		d.middlewareRequestCustomizers(c),
		d.middlewareHandleRetries(c),
		// lower layer ^^^^^
	}

	if result := obtainFromConfig[[]Middleware](c, "Middlewares"); result.isErr() {
		return nil, result.error()
	} else if m, ok := result.some(); ok {
		//nolint:gocritic // this is intentional
		d.middlewares = append(m, middlewares...) // prepend
	} else {
		d.middlewares = middlewares
	}

	return &d, nil
}

// Deprecated: only for compatibility.
func (d *doer) middlewareRequestCustomizers(c *config) Middleware {
	return func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		iter := func(yield saht.RequestCustomizer) error { return yield(req) }

		if result := obtainFromConfig[[]saht.RequestCustomizer](c, "RequestCustomizer"); result.isErr() {
			return nil, result.error()
		} else if customizers, ok := result.some(); !ok {
			// nothing to do here
		} else if err := findFirstError(slices.Values(customizers), iter); err != nil {
			return nil, err
		}
		return pullThenCall(pull, req)
	}
}

func pullThenCall(pull func() (Middleware, bool), req *http.Request) (*http.Response, error) {
	if cont, ok := pull(); !ok {
		return nil, NewErrorf("no next middleware to pull")
	} else {
		return cont(req, pull)
	}
}
