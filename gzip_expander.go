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
	"errors"
	"io"
	"net/http"
)

var gzipExpander Middleware = func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
	if resp, err := pullThenCall(pull, req); err != nil {
		return resp, err

	} else if resp.Header.Get("Content-Encoding") != "gzip" {
		return resp, nil

	} else if body, err := gzip.NewReader(resp.Body); err != nil {
		if errors.Is(err, io.EOF) && resp.ContentLength == 0 {
			return resp, nil
		}
		return resp, err

	} else {
		resp.Body = body
		resp.Header.Del("Content-Encoding")
		resp.Header.Del("Content-Length") // unknown length
		resp.ContentLength = -1           // ditto
		return resp, nil
	}
}
