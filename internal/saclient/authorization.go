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
		}

		pref, ok, err := obtainFromConfig[string](c, "AuthPreference").decompose()

		if err != nil {
			return nil, err
		}

		switch {
		case ok:
			mode = pref
		case c.hasSome("PrivateKeyPEMPath"):
			mode = "bearer"
		case c.hasSome("PrivateKey"):
			mode = "bearer"
		case c.hasSome("AccessToken"):
			mode = "basic"
		default:
			// no auth info found ... ?
			// let's just let it go without auth
			// and let server return 401.
			return pullThenCall(pull, req)
		}

		switch mode {
		default:
			panic("unknown authPreference: " + mode)

		case "basic":
			user, ok, err := obtainFromConfig[string](c, "AccessToken").decompose()

			if err != nil {
				return nil, err
			}

			if !ok || user == "" {
				return nil, NewErrorf("missing AccessToken")
			}

			pass, ok, err := obtainFromConfig[string](c, "AccessTokenSecret").decompose()

			if err != nil {
				return nil, err
			}

			if !ok || pass == "" {
				return nil, NewErrorf("missing AccessTokenSecret")
			}

			req.SetBasicAuth(user, pass)
			return pullThenCall(pull, req)

		case "bearer":
			_, ok, err := obtainFromConfig[string](c, "TokenEndpoint").decompose()

			if err != nil {
				return nil, err
			}

			if !ok {
				// UNLIKELY it has default value
				return nil, NewErrorf("TokenEndpoint is absent")
			}

			_, ok, err = obtainFromConfig[string](c, "ServicePrincipalID").decompose()

			if err != nil {
				return nil, err
			}

			if !ok {
				return nil, NewErrorf("ServicePrincipalID is absent")
			}

			_, ok, err = obtainFromConfig[string](c, "ServicePrincipalKeyID").decompose()

			if err != nil {
				return nil, err
			}

			if !ok {
				return nil, NewErrorf("ServicePrincipalKeyID is absent")
			}

			_, path, err := obtainFromConfig[string](c, "PrivateKeyPEMPath").decompose()

			if err != nil {
				return nil, err
			}

			_, key, err := obtainFromConfig[string](c, "PrivateKey").decompose()

			if err != nil {
				return nil, err
			}

			if !key && !path {
				return nil, NewErrorf("neither PrivateKeyPEMPath nor PrivateKey is present")
			}

			token, err := d.newTokenResponse(req.Context(), c)

			if err != nil {
				return nil, err
			}

			req.Header.Set("Authorization", token.HTTPAuthorizationHeader())
			return pullThenCall(pull, req)
		}
	}
}
