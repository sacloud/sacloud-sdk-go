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

import "net/http"

func (c *config) middlewareSetHeader() Middleware {
	return func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		if req.Header.Get("User-Agent") == "" {
			ua := obtainFromConfig[string](c, "UserAgent").unwrap()
			req.Header.Set("User-Agent", ua)
		}

		if req.Body != nil && req.Header.Get("Content-Type") == "" {
			req.Header.Set("Content-Type", "application/json")
		}

		if req.Header.Get("Accept-Encoding") == "" {
			// :TODO: do we want to _disable_ gzip?
			req.Header.Add("Accept-Encoding", "gzip")
		}

		if req.Header.Get("X-Requested-With") == "" {
			req.Header.Add("X-Requested-With", "XMLHttpRequest")
		}

		if req.Header.Get("X-Sakura-Bigint-As-Int") == "" {
			req.Header.Add("X-Sakura-Bigint-As-Int", "1") // default
		}

		return pullThenCall(pull, req)
	}
}
