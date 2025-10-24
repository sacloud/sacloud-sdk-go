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

package client

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"strings"
)

func (d *doer) tracer(c *config) middleware {
	return func(req *http.Request, pull func() (middleware, bool)) (*http.Response, error) {
		var buf []byte

		if result := obtainFromConfig[string](c, "TraceMode"); result.isErr() {
			return nil, result.error()

		} else if mode, ok := result.some(); !ok {
			return pullThenCall(pull, req)

		} else {
			if req.Body == nil {
				// go below
			} else if buf, err := io.ReadAll(req.Body); err != nil {
				return nil, err

			} else {
				copied := bytes.NewBuffer(buf)
				req.Body = io.NopCloser(copied)
			}

			if res, err := pullThenCall(pull, req); err != nil {
				return res, err

			} else if mode == "error" && res.StatusCode < 300 {
				// why this is 300 rather than 400 -----> ^^^ <--- is not obvious to @shyouhei.
				// Just mimicing sacloud/go-http.
				return res, err

			} else {
				if req.Body != nil {
					// write back buffer
					copied := bytes.NewBuffer(buf)
					req.Body = io.NopCloser(copied)
				}

				return __traceDump(req, res)
			}
		}
	}
}

func __traceDump(req *http.Request, res *http.Response) (*http.Response, error) {
	if dump, err := httputil.DumpRequest(req, true); err != nil {
		return nil, err
	} else {
		log.Printf("[TRACE] \trequest: %s %s\n", req.Method, req.URL.String())
		log.Printf("==============================\n")
		for line := range strings.Lines(string(dump)) {
			log.Printf("%s", line)
		}
		log.Printf("==============================\n")
	}

	if dump, err := httputil.DumpResponse(res, true); err != nil {
		return nil, err
	} else {
		log.Printf("[TRACE] \tresponse: %s %s\n", req.Method, req.URL.String())
		log.Printf("==============================\n")
		for line := range strings.Lines(string(dump)) {
			log.Printf("%s", line)
		}
		log.Printf("==============================\n")
	}

	return res, nil
}
