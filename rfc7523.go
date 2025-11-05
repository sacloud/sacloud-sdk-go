// Copyright 2025- The sacloud/http-client-go Authors
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
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofrs/flock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sacloud/packages-go/testutil"
)

const skew = 30 * time.Second // :FIXME: should be configurable?
const retryDelay = 1024 * time.Millisecond

type tokenResponse struct {
	// This is non-standard but it _seems_ the server returns token_expired_at
	TokenExpiredAt time.Time `json:"token_expired_at,omitzero"`

	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope,omitempty"`
}

type cachedTokenResponse struct {
	Token    tokenResponse `json:"token"`
	CachedAt time.Time     `json:"cached_at"`
}

func (d *doer) inquireAccessToken(ctx context.Context, cfg *config) (*tokenResponse, error) {
	var err error
	var buf []byte
	var ret tokenResponse
	iat := time.Now().Unix()
	exp := iat + 600 // 10 min

	// these fields are checked beforehand
	aud := obtainFromConfig[string](cfg, "TokenEndpoint").unwrap()
	sub := obtainFromConfig[string](cfg, "ServicePrincipalID").unwrap()
	kid := obtainFromConfig[string](cfg, "ServicePrincipalKeyID").unwrap()

	if key := obtainFromConfig[string](cfg, "PrivateKey").asPtr(); key != nil {
		buf = []byte(*key)
	} else if path := obtainFromConfig[string](cfg, "PrivateKeyPEMPath").asPtr(); path != nil {
		//nolint:gosec // This `os.ReadFile` does not reveal any secret info
		if buf, err = os.ReadFile(*path); err != nil {
			return nil, err
		}
	} else {
		// UNLIKELY this is checked a priori
		return nil, NewErrorf("neither PrivateKeyPEMPath nor PrivateKey is present")
	}

	k, err := jwt.ParseRSAPrivateKeyFromPEM(buf)
	if err != nil {
		return nil, err
	}

	claims := jwt.MapClaims{
		"iss": sub,
		"sub": sub,
		"aud": aud,
		"iat": iat,
		"exp": exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	token.Header["kid"] = kid
	assertion, err := token.SignedString(k)
	if err != nil {
		return nil, err
	}

	form := make(url.Values)
	form.Set("grant_type", "urn:ietf:params:oauth:grant-type:jwt-bearer")
	form.Set("assertion", assertion)
	req, err := http.NewRequestWithContext(ctx, "POST", aud, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	} else {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	var resp *http.Response
	if d.server == nil {
		if d.client == nil {
			// should we error here...?
			d.client = http.DefaultClient
		}
		if resp, err = d.client.Do(req); err != nil {
			return nil, err
		} else {
			defer func() {
				_ = resp.Body.Close()
			}()
		}
	} else {
		// :BEWARE: in case of using httptest.Server, it is not a wise idea to
		// issue actual HTTP request to the token endpoint.  Instead we create
		// another httptest.Server here for one-shot mock and use it.
		svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Header().Set("Content-Type", "application/json")
			if j, err := json.Marshal(map[string]any{
				"access_token": testutil.Random(32, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"),
				"token_type":   "Bearer",
				"expires_in":   3600,
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else if _, err = w.Write(j); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}))
		defer svr.Close()

		if req.URL, err = url.Parse(svr.URL); err != nil {
			return nil, err
		} else if resp, err = svr.Client().Do(req); err != nil {
			return nil, err
		} else {
			defer func() {
				_ = resp.Body.Close()
			}()
		}
	}

	if resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return nil, NewError(resp.StatusCode, string(b), nil)

	} else if err := json.NewDecoder(resp.Body).Decode(&ret); err != nil {
		return nil, err

	} else {
		return &ret, nil
	}
}

func (d *doer) newTokenResponse(ctx context.Context, cfg *config) (*tokenResponse, error) {
	var profile *Profile
	var key, path *string
	var exists bool

	if result := obtainFromConfig[*Profile](cfg, "Profile"); result.isErr() {
		return nil, result.error()

	} else if profile, exists = result.Get(); !exists {
		// It is possible that Profile is absent.
		// Token request is still tried then, but no way to cache its response.
		return d.inquireAccessToken(ctx, cfg)

	} else if key = obtainFromConfig[string](cfg, "PrivateKey").asPtr(); key != nil {
		// ok

	} else if path = obtainFromConfig[string](cfg, "PrivateKeyPEMPath").asPtr(); path != nil {
		// ok

	} else {
		// UNLIKELY this is checked a priori
		return nil, NewErrorf("neither PrivateKeyPEMPath nor PrivateKey is present")
	}

	if path, err := profile.GetCacheFilePath(key, path); err != nil {
		// This is e.g. malformed PEM; worth propagating the error
		return nil, err

	} else if rel, err := filepath.Rel(profile.dir, path); err != nil {
		// (unlikely)
		return nil, err

	} else {
		return openFileAt(profile.dir, rel, os.O_RDWR|os.O_CREATE, func(fp *os.File) (*tokenResponse, error) {
			var lock *flock.Flock
			var cache cachedTokenResponse

			// At this point, because we open a file without O_EXCL, another process could be touching it.
			// We need to make sure that we don't read something incomplete.
			if lock = flock.New(path + ".lock"); lock == nil {
				// no lock; no caching
				return d.inquireAccessToken(ctx, cfg)

			} else if locked, err := lock.TryRLockContext(ctx, retryDelay); err != nil {
				// context deadline exceeded etc.; worth propagating the error
				return nil, err

			} else if !locked {
				return d.inquireAccessToken(ctx, cfg)

			} else {
				defer func() {
					_ = lock.Unlock()
				}()
			}

			dec := json.NewDecoder(fp)
			if err := dec.Decode(&cache); err != nil {
				// corrupt or nonexistent cache; ignore and request new token

			} else if !cache.isExpired() {
				// got it! return cached token
				return &cache.Token, nil
			}

			enc := json.NewEncoder(fp)
			enc.SetIndent("", "  ")
			if res, err := d.inquireAccessToken(ctx, cfg); err != nil {
				return nil, err

			} else if locked, err := lock.TryLockContext(ctx, retryDelay); err != nil {
				// lock promotion failed; no write and return
				return res, nil

			} else if !locked {
				// ditto
				return res, nil

			} else if err := fp.Truncate(0); err != nil {
				return nil, err

			} else if _, err := fp.Seek(0, 0); err != nil {
				return nil, err

			} else if err := enc.Encode(cachedTokenResponse{*res, time.Now()}); err != nil {
				// write failed; purge garbages
				_ = fp.Truncate(0)
				return nil, err

			} else {
				return res, nil
			}
		})
	}
}

func (c *cachedTokenResponse) isExpired() bool {
	expiresIn := time.Duration(c.Token.ExpiresIn) * time.Second
	expiresAt := c.CachedAt.Add(expiresIn)
	return time.Now().Add(-skew).After(expiresAt)
}
