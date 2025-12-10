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

func (d *doer) middlewareAuthorization(c *config) Middleware {
	return func(req *http.Request, pull func() (Middleware, bool)) (*http.Response, error) {
		var mode string

		if req.Header.Get("Authorization") != "" {
			// already set, skip
			return pullThenCall(pull, req)
		} else if result := obtainFromConfig[string](c, "AuthPreference"); result.isErr() {
			return nil, result.error()
		} else if pref, ok := result.some(); ok {
			mode = pref
		} else if result := obtainFromConfig[string](c, "AccessToken"); result.isSome() {
			mode = "basic"
		} else if result := obtainFromConfig[string](c, "PrivateKeyPEMPath"); result.isSome() {
			mode = "bearer"
		} else if result := obtainFromConfig[string](c, "PrivateKey"); result.isSome() {
			mode = "bearer"
		} else {
			// no auth info found ... ?
			// let's just let it go without auth
			// and let server return 401.
			return pullThenCall(pull, req)
		}

		switch mode {
		default:
			panic("unknown authPreference: " + mode)

		case "basic":
			if result := obtainFromConfig[string](c, "AccessToken"); result.isErr() {
				return nil, result.error()
			} else if user, ok := result.some(); !ok || user == "" {
				return nil, NewErrorf("missing AccessToken")
			} else if result := obtainFromConfig[string](c, "AccessTokenSecret"); result.isErr() {
				return nil, result.error()
			} else if pass, ok := result.some(); !ok || pass == "" {
				return nil, NewErrorf("missing AccessTokenSecret")
			} else {
				req.SetBasicAuth(user, pass)
				return pullThenCall(pull, req)
			}

		case "bearer":
			if result := obtainFromConfig[string](c, "TokenEndpoint"); result.isErr() {
				return nil, result.error()
			} else if !result.isSome() {
				// UNLIKELY it has default value
				return nil, NewErrorf("TokenEndpoint is absent")
			} else if result := obtainFromConfig[string](c, "ServicePrincipalID"); result.isErr() {
				return nil, result.error()
			} else if !result.isSome() {
				return nil, NewErrorf("ServicePrincipalID is absent")
			} else if result := obtainFromConfig[string](c, "ServicePrincipalKeyID"); result.isErr() {
				return nil, result.error()
			} else if !result.isSome() {
				return nil, NewErrorf("ServicePrincipalKeyID is absent")
			} else if !obtainFromConfig[string](c, "PrivateKeyPEMPath").isSome() && !obtainFromConfig[string](c, "PrivateKey").isSome() {
				return nil, NewErrorf("neither PrivateKeyPEMPath nor PrivateKey is present")
			} else if token, err := d.newTokenResponse(req.Context(), c); err != nil {
				return nil, err
			} else {
				req.Header.Set("Authorization", token.HTTPAuthorizationHeader())
				return pullThenCall(pull, req)
			}
		}
	}
}
