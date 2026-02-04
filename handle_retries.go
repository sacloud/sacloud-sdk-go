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

	ret.Backoff = retryablehttp.DefaultBackoff

	return &ret, nil
}
