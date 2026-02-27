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
	"compress/gzip"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

func (d *doer) middlewareHandleRetries(c *config) Middleware {
	return func(r *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		if cont, ok := pull(); ok {
			// This _has_ to be the last middleware
			return nil, NewErrorf("broken middleware chain: got %v", cont)
		} else if client, err := d.retryableClient(c); err != nil {
			return nil, err
		} else if req, err := retryablehttp.FromRequest(r); err != nil {
			return nil, err
		} else {
			return client.Do(req)
		}
	}
}

func (d *doer) retryableClient(c *config) (*retryablehttp.Client, error) {
	var ret retryablehttp.Client

	ret.HTTPClient = d.client

	v, ok, err := obtainFromConfig[int64](c, "RetryWaitMin").decompose()

	if err != nil {
		return nil, err
	}

	if ok {
		ret.RetryWaitMin = time.Duration(v) * time.Second
	}

	v, ok, err = obtainFromConfig[int64](c, "RetryWaitMax").decompose()

	if err != nil {
		return nil, err
	}

	if ok {
		ret.RetryWaitMax = time.Duration(v) * time.Second
	}

	v, ok, err = obtainFromConfig[int64](c, "RetryMax").decompose()

	if err != nil {
		return nil, err
	}

	if ok {
		ret.RetryMax = int(v)
	}

	f, ok, err := obtainFromConfig[retryablehttp.CheckRetry](c, "CheckRetryFunc").decompose()

	if err != nil {
		return nil, err
	}

	if ok {
		ret.CheckRetry = f
	}

	if c.hasSome("TraceMode") {
		superCheckRetry := ret.CheckRetry

		ret.CheckRetry = func(ctx context.Context, res *http.Response, e error) (ret bool, err error) {
			ret, err = superCheckRetry(ctx, res, e)

			// superCheckRetry thinks it's worth retrying.
			// Dump the previous attempt at this point,
			// since it would be lost after retrying.
			// (Request payload had already been lost here
			// which we cannot do anything about)
			if ret == true {
				if res.Header.Get("Content-Encoding") == "gzip" {
					// make it readable
					res.Header.Del("Content-Encoding")
					res.Header.Del("Content-Length") // unknown length
					res.ContentLength = -1           // ditto
					if body, err := gzip.NewReader(res.Body); err == nil {
						res.Body = body
					} else {
						// :UNLIKELY: this is e.g. empty error message
						res.Body = io.NopCloser(strings.NewReader("(redacted)"))
					}
				}
				_, err2 := dumpTracePair(res.Request, res)
				err = errors.Join(err, err2)
			}

			return
		}
	}

	// This is callled right before the retryablehttp client gives up.
	// without it the `res` would be lost
	ret.ErrorHandler = func(res *http.Response, err error, n int) (*http.Response, error) {
		req := res.Request
		msg := fmt.Sprintf("%s %s giving up after %d attempt(s)", req.Method, req.URL, n)

		if err != nil {
			err = fmt.Errorf("%s: %w", msg, err)
		} else {
			err = fmt.Errorf("%s", msg)
		}

		return res, err
	}

	ret.Backoff = retryablehttp.DefaultBackoff

	return &ret, nil
}
