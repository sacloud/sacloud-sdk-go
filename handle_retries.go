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
	"net/http"
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

	if result := obtainFromConfig[int64](c, "RetryWaitMin"); result.isErr() {
		return nil, result.error()
	} else if v, ok := result.some(); ok {
		ret.RetryWaitMin = time.Duration(v) * time.Second
	}

	if result := obtainFromConfig[int64](c, "RetryWaitMax"); result.isErr() {
		return nil, result.error()
	} else if v, ok := result.some(); ok {
		ret.RetryWaitMax = time.Duration(v) * time.Second
	}

	if result := obtainFromConfig[int64](c, "RetryMax"); result.isErr() {
		return nil, result.error()
	} else if v, ok := result.some(); ok {
		ret.RetryMax = int(v)
	}

	if result := obtainFromConfig[retryablehttp.CheckRetry](c, "CheckRetryFunc"); result.isErr() {
		return nil, result.error()
	} else if v, ok := result.some(); ok {
		ret.CheckRetry = v
	}

	ret.Backoff = retryablehttp.DefaultBackoff

	return &ret, nil
}
