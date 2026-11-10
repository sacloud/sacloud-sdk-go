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
	"context"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gofrs/flock"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
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

	isFake bool `json:"-"`
}

type cachedTokenResponse struct {
	Token    tokenResponse `json:"token"`
	CachedAt time.Time     `json:"cached_at"`
}

type inMemoryCache struct {
	sync.Mutex
	response option[cachedTokenResponse]
}

var caches sync.Map

func (d *doer) inquireAccessToken(ctx context.Context, cfg *config) (*cachedTokenResponse, error) {
	if tok, err := d.inquireAccessTokenInternal(ctx, cfg); err != nil {
		return nil, err
	} else {
		// This return value can include fake tokens if using a test server.
		// Should be checked at caller site.
		return &cachedTokenResponse{*tok, time.Now()}, nil
	}
}

func (d *doer) inquireAccessTokenInternal(ctx context.Context, cfg *config) (*tokenResponse, error) {
	var err error
	var buf []byte
	var ret tokenResponse
	iat := time.Now().Unix()
	exp := iat + 600 // 10 min
	nbf := iat - 60  // 1 min

	// these fields are checked beforehand
	aud := obtainFromConfig[string](cfg, "TokenEndpoint").unwrap()
	sub := obtainFromConfig[string](cfg, "ServicePrincipalID").unwrap()
	kid := obtainFromConfig[string](cfg, "ServicePrincipalKeyKID").unwrap()

	if key := obtainFromConfig[string](cfg, "PrivateKey").asPtr(); key != nil {
		buf = []byte(*key)
	} else if path := obtainFromConfig[string](cfg, "PrivateKeyPEMPath").asPtr(); path != nil {
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
		"nbf": nbf,
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
		// a temporary fake token response to be used later.
		return &tokenResponse{
			AccessToken: testutil.Random(32, testutil.CharSetAlphaNum),
			TokenType:   "Bearer",
			ExpiresIn:   3600,
			isFake:      true,
		}, nil
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
	if k, err := cacheKey(cfg); err != nil {
		return nil, err
	} else {
		return cacheFor(k).Transaction(func() (*cachedTokenResponse, error) {
			return d.newTokenResponseInternal(ctx, cfg)
		})
	}
}

func (d *doer) newTokenResponseInternal(ctx context.Context, cfg *config) (*cachedTokenResponse, error) {
	var profile *Profile
	var key, path *string
	var exists bool
	var err error

	profile, exists, err = obtainFromConfig[*Profile](cfg, "Profile").decompose()

	if err != nil {
		return nil, err
	}

	if !exists {
		// It is possible that Profile is absent.
		// Token request is still tried then, kept in memory only.
		return d.inquireAccessToken(ctx, cfg)
	}

	if key = obtainFromConfig[string](cfg, "PrivateKey").asPtr(); key != nil {
		// ok

	} else if path = obtainFromConfig[string](cfg, "PrivateKeyPEMPath").asPtr(); path != nil {
		// ok

	} else {
		// UNLIKELY this is checked a priori
		return nil, NewErrorf("neither PrivateKeyPEMPath nor PrivateKey is present")
	}

	if path, err := profile.GetCacheFilePath(path, key); err != nil {
		// This is e.g. malformed PEM; worth propagating the error
		return nil, err
	} else if rel, err := filepath.Rel(profile.dir, path); err != nil {
		// (unlikely)
		return nil, err
	} else {
		yielded := false
		res, err := openFileAt(profile.dir, rel, os.O_RDWR|os.O_CREATE, func(fp *os.File) (*cachedTokenResponse, error) {
			var lock *flock.Flock
			var cache cachedTokenResponse
			yielded = true

			// At this point, because we open a file without O_EXCL, another process could be touching it.
			// We need to make sure that we don't read something incomplete.
			if lock = flock.New(path + ".lock"); lock == nil {
				// no file lock; use the in-process cache
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
				return &cache, nil
			}

			enc := json.NewEncoder(fp)
			enc.SetIndent("", "  ")
			locked, eerr := lock.TryLockContext(ctx, retryDelay)
			if res, err := d.inquireAccessToken(ctx, cfg); err != nil {
				return nil, err
			} else if res.Token.isFake {
				// This response is not for caching.
				return res, nil
			} else if eerr != nil {
				// lock promotion failed; no write and return
				return res, nil
			} else if !locked {
				// ditto
				return res, nil
			} else if err := fp.Truncate(0); err != nil {
				// write failed; keep the response in memory
				return res, nil
			} else if _, err := fp.Seek(0, 0); err != nil {
				// ditto
				return res, nil
			} else if err := enc.Encode(res); err != nil {
				// write failed; purge garbages
				_ = fp.Truncate(0)
				return res, nil
			} else {
				return res, nil
			}
		})

		if yielded {
			return res, err
		} else {
			// premature failure of openFileAt for some reason;
			// go on without filesystem caches
			return d.inquireAccessToken(ctx, cfg)
		}
	}
}

func (r *tokenResponse) HTTPAuthorizationHeader() string {
	return fmt.Sprintf("%s %s", r.TokenType, r.AccessToken)
}

func (c *cachedTokenResponse) isExpired() bool {
	expiresIn := time.Duration(c.Token.ExpiresIn) * time.Second
	expiresAt := c.CachedAt.Add(expiresIn)
	return time.Now().Add(skew).After(expiresAt)
}

func cacheKey(cfg *config) (string, error) {
	var keyPEM []byte

	if cfg == nil {
		return "", NewErrorf("nil config passed to tokenResponseCacheKey")
	}

	aud, ok, err := obtainFromConfig[string](cfg, "TokenEndpoint").decompose()

	if err != nil {
		return "", err
	}

	if !ok {
		return "", NewErrorf("TokenEndpoint not found in config")
	}

	sub, ok, err := obtainFromConfig[string](cfg, "ServicePrincipalID").decompose()

	if err != nil {
		return "", err
	}

	if !ok {
		return "", NewErrorf("ServicePrincipalID not found in config")
	}

	kid, ok, err := obtainFromConfig[string](cfg, "ServicePrincipalKeyKID").decompose()

	if err != nil {
		return "", err
	}

	if !ok {
		return "", NewErrorf("ServicePrincipalKeyKID not found in config")
	}

	if key := obtainFromConfig[string](cfg, "PrivateKey").asPtr(); key != nil {
		keyPEM = []byte(*key)
	} else if path := obtainFromConfig[string](cfg, "PrivateKeyPEMPath").asPtr(); path != nil {
		var err error
		keyPEM, err = os.ReadFile(*path)
		if err != nil {
			return "", err
		}
	} else {
		return "", NewErrorf("neither PrivateKeyPEMPath nor PrivateKey is present")
	}

	key, err := jwt.ParseRSAPrivateKeyFromPEM(keyPEM)
	if err != nil {
		return "", err
	}

	publicKey, err := x509.MarshalPKIXPublicKey(&key.PublicKey)
	if err != nil {
		return "", err
	}

	identity, err := json.Marshal([]string{aud, sub, kid, string(publicKey)})
	if err != nil {
		return "", err
	}

	sum := sha256.Sum256(identity)
	return hex.EncodeToString(sum[:]), nil
}

func cacheFor(key string) *inMemoryCache {
	allocated := new(inMemoryCache)
	cache, _ := caches.LoadOrStore(key, allocated)
	return cache.(*inMemoryCache)
}

func (c *inMemoryCache) Transaction(yield func() (*cachedTokenResponse, error)) (res *tokenResponse, err error) {
	if c == nil {
		err = NewErrorf("in-memory cache is nil")
		return
	}

	c.Lock()
	defer c.Unlock()

	if cached, ok := c.response.Get(); ok && !cached.isExpired() {
		return &cached.Token, nil
	}

	obtained, err := yield()

	if err != nil {
		return
	}

	if obtained == nil {
		return
	}

	res = &obtained.Token

	if res.isFake {
		return
	}

	c.response = optionSome(*obtained)

	return
}
