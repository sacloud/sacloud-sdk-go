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
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type maybeUninit[T any] struct {
	value T
	set   bool
}

type storage struct {
	profileName         maybeUninit[string]
	privateKeyPath      maybeUninit[string]
	accessToken         maybeUninit[string]
	accessTokenSecret   maybeUninit[string]
	zone                maybeUninit[string]
	defaultZone         maybeUninit[string]
	zones               maybeUninit[[]string]
	retryMax            maybeUninit[int64]
	retryWaitMax        maybeUninit[int64]
	retryWaitMin        maybeUninit[int64]
	apiRootURL          maybeUninit[string]
	apiRequestTimeout   maybeUninit[int64]
	apiRequestRateLimit maybeUninit[int64]
	traceMode           maybeUninit[string]
}

// :INTERNAL: it is intentional that this is not a struct
// This is also not JSONable because it contains functions etc.
type config map[string]any

type parameter struct {
	profileOp *ProfileOp
	envp      storage
	argv      storage
	hcl       storage
}

func (p *parameter) setEnvironIter() func(string, string) error {
	return func(k, v string) error {
		if p == nil {
			return NewErrorf("nil parameter")
		} else {
			switch k {
			case "SACLOUD_PROFILE":
				return p.envp.profileName.Set(v)

			case "SACLOUD_PRIVATE_KEY_PATH":
				return p.envp.privateKeyPath.Set(v)

			case "SAKURACLOUD_ACCESS_TOKEN":
				return p.envp.accessToken.Set(v)

			case "SAKURACLOUD_ACCESS_TOKEN_SECRET":
				return p.envp.accessTokenSecret.Set(v)

			case "SAKURACLOUD_ZONE":
				return p.envp.zone.Set(v)

			case "SAKURACLOUD_ZONES":
				return p.envp.zones.Set(v)

			case "SAKURACLOUD_RETRY_MAX":
				return p.envp.retryMax.Set(v)

			case "SAKURACLOUD_RETRY_WAIT_MAX":
				return p.envp.retryWaitMax.Set(v)

			case "SAKURACLOUD_RETRY_WAIT_MIN":
				return p.envp.retryWaitMin.Set(v)

			case "SAKURACLOUD_API_ROOT_URL":
				return p.envp.apiRootURL.Set(v)

			case "SAKURACLOUD_API_REQUEST_TIMEOUT":
				return p.envp.apiRequestTimeout.Set(v)

			case "SAKURACLOUD_API_REQUEST_RATE_LIMIT":
				return p.envp.apiRequestRateLimit.Set(v)

			// These names are historical, cannot change at this point.
			case "USACLOUD_PROFILE":
				return p.envp.profileName.Set(v)

			case "SAKURACLOUD_RATE_LIMIT":
				return p.envp.apiRequestRateLimit.Set(v)

			case "SAKURACLOUD_TRACE":
				return p.envp.traceMode.Set(v)

			case "USACLOUD_TRACE":
				return p.envp.traceMode.Set(v)

			default:
				return nil
			}
		}
	}
}

func (p *parameter) setEnviron(env []string) error {
	if p == nil {
		return NewErrorf("nil parameter")
	} else {
		p.profileOp = NewProfileOp(env)
	}

	q := slices.Values(env)
	w := intoSeq2(q, func(i string) (string, string, bool) { return strings.Cut(i, "=") })
	e := transformSeq3(w, p.setEnvironIter())
	r := transformSeq3(e, func(k string, v error) error { return Wrapf(v, "in parsing %s", k) })
	t := valuesOfSeq2(r)
	y := slices.Collect(t)

	return errors.Join(y...)
}

func (p *parameter) setHCL(config TerraformProviderInterface) {
	if p == nil {
		return
	}

	p.hcl.profileName.from(config.LookupClientConfigProfileName)
	p.hcl.privateKeyPath.from(config.LookupClientConfigPrivateKeyPath)
	p.hcl.accessToken.from(config.LookupClientConfigAccessToken)
	p.hcl.accessTokenSecret.from(config.LookupClientConfigAccessTokenSecret)
	p.hcl.zone.from(config.LookupClientConfigZone)
	p.hcl.zones.from(config.LookupClientConfigZones)
	p.hcl.defaultZone.from(config.LookupClientConfigDefaultZone)
	p.hcl.retryMax.from(config.LookupClientConfigRetryMax)
	p.hcl.retryWaitMax.from(config.LookupClientConfigRetryWaitMax)
	p.hcl.retryWaitMin.from(config.LookupClientConfigRetryWaitMin)
	p.hcl.apiRootURL.from(config.LookupClientConfigAPIRootURL)
	p.hcl.apiRequestTimeout.from(config.LookupClientConfigAPIRRequestTimeout)
	p.hcl.apiRequestRateLimit.from(config.LookupClientConfigAPIRequestRateLimit)
	p.hcl.traceMode.from(config.LookupClientConfigTraceMode)
}

func (p *parameter) flagSet() *flag.FlagSet {
	var fs *flag.FlagSet

	if p != nil {
		fs = flag.NewFlagSet("http-client-go", flag.PanicOnError)

		// :NOTE: these help messages are from usacloud's old --help output
		fs.Var(&p.argv.profileName, "profile", "the name of saved credentials")
		fs.Var(&p.argv.privateKeyPath, "private-key-path", "path to an RSA 2048 bit private key PEM format")
		fs.Var(&p.argv.accessToken, "token", "the API token used when calling SAKURA Cloud API")
		fs.Var(&p.argv.accessTokenSecret, "secret", "the API secret used when calling SAKURA Cloud API")
		fs.Var(&p.argv.zones, "zones", "permitted zone names")
		fs.Var(&p.argv.traceMode, "trace", "enable trace logs for API calling")

		// Not sure why but not everything can be specified from command line
		// for instance usacloud lacks --zone, in spiyte of having --zones.
	}

	return fs
}

func (p *parameter) populate(c *config) error {
	// This is the mother-of-all populate function.
	ret := make([]error, 0, 16) // <- 16 is the # of `append` calls below

	//nolint:gocritic
	if p == nil {
		return NewErrorf("nil parameter")
	} else if c == nil {
		return NewErrorf("nil config")
	} else if p.profileOp == nil {
		// Operator not initialized, means there was no call to SetEnviron()
		// This could be meddling, but we initialize it here for safety.
		ret = append(ret, p.setEnviron(os.Environ()))
	}

	*c = make(config)
	ret = append(ret, p.populateProfile(c))
	ret = append(ret, p.populatePrivateKeyPath(c))
	ret = append(ret, p.populateAccessToken(c))
	ret = append(ret, p.populateAccessTokenSecret(c))
	ret = append(ret, p.populateZone(c))
	ret = append(ret, p.populateZones(c)...)
	ret = append(ret, p.populateDefaultZone(c))
	ret = append(ret, p.populateRetryMax(c))
	ret = append(ret, p.populateRetryWaitMax(c))
	ret = append(ret, p.populateRetryWaitMin(c))
	ret = append(ret, p.populateAPIRootURL(c))
	ret = append(ret, p.populateAPIRequestTimeout(c))
	ret = append(ret, p.populateAPIRequestRateLimit(c))
	ret = append(ret, p.populateTraceMode(c))

	if v, err := c.get("AccessToken"); err == nil && v.isSet() {
		// Take that,
	} else if v, err := c.get("PrivateKeyPEMPath"); err == nil && v.isSet() {
		// Take that,
	} else {
		// This is fatal.  Stop here.
		ret = append(ret, NewErrorf("neither AccessToken nor PrivateKeyPEMPath is set"))
	}

	return errors.Join(ret...)
}

func (p *parameter) populateProfile(c *config) error {
	// We need to load profile.
	// The one from env var has the highest priority,
	// then the one from command-line flag,
	// and finally the one from Terraform provider is the lowest priority.
	// In case none of them are set, the "current" profile is used.
	var profileName maybeUninit[string]

	if p == nil {
		return NewErrorf("nil parameter")
	} else if v, ok := p.envp.profileName.Get(); ok {
		profileName.initialize(v)
	} else if v, ok := p.argv.profileName.Get(); ok {
		profileName.initialize(v)
	} else if v, ok := p.hcl.profileName.Get(); ok {
		profileName.initialize(v)
	} else if v, err := p.profileOp.GetCurrentName(); err == nil {
		profileName.initialize(v)
	}

	if v, ok := profileName.Get(); !ok {
		// None of above succeeded, and there is no "current" profile.
		// Maybe the user opted to not use profiles at all.
		// This is not an error, continue populating with empty profile.
		return nil
	} else if profile, err := p.profileOp.Read(v); err != nil {
		// Explicitly specified profile not found, this is surely an error.
		return err
	} else {
		return c.set("Profile", profile)
	}
}

//nolint:gocritic
func (p *parameter) populatePrivateKeyPath(c *config) error {
	if err := p.populateString(c, "PrivateKeyPEMPath"); err != nil {
		return err
	} else if path, err := c.get("PrivateKeyPEMPath"); err != nil {
		return err
	} else if v, ok := path.Get(); !ok {
		return nil // just not set
	} else if v, ok := v.(string); !ok {
		return NewErrorf("invalid type for PrivateKeyPEMPath in config: %T", v)
	} else if s, err := os.Stat(v); err != nil {
		return NewErrorf("private key file not found: %s", v)
	} else if !s.Mode().IsRegular() {
		return NewErrorf("private key not a file: %s", v)
	} else if s.Mode().Perm()&0o077 != 0 {
		return NewErrorf("private key file %s permission is too lax: %o", v, s.Mode().Perm())
	} else {
		return nil
	}
}

func (p *parameter) populateAccessToken(c *config) error {
	return p.populateString(c, "AccessToken")
}

func (p *parameter) populateAccessTokenSecret(c *config) error {
	return p.populateString(c, "AccessTokenSecret")
}

func (p *parameter) populateZone(c *config) error {
	return p.populateString(c, "Zone")
}

func (p *parameter) populateDefaultZone(c *config) error {
	return p.populateString(c, "DefaultZone")
}

func (p *parameter) populateZones(c *config) []error {
	var ret []error
	var whence string
	var val maybeUninit[[]string]

	if p == nil {
		ret = append(ret, NewErrorf("nil parameter"))
	} else if c == nil {
		ret = append(ret, NewErrorf("nil config"))
	} else if v, ok := p.envp.zones.Get(); ok {
		val.initialize(v)
		whence = "environment variable"
	} else if v, ok := p.argv.zones.Get(); ok {
		val.initialize(v)
		whence = "command-line argument"
	} else if v, ok := p.hcl.zones.Get(); ok {
		val.initialize(v)
		whence = "terraform configuration"
	} else if wal, whence, err := obtainFromProfile[[]any](c, "Zones", "profile"); err != nil {
		ret = append(ret, err)
	} else if v, ok := wal.Get(); !ok {
		// just not set

	} else {
		w := []string{}
		for i, z := range v {
			if s, ok := z.(string); !ok {
				ret = append(ret, NewErrorf("nonstring zone %v in %s's #%d", z, whence, i))
			} else {
				w = append(w, s)
			}
		}
		val.initialize(w)
	}

	if v, ok := val.Get(); !ok {
		// just not set
	} else if len(v) == 0 {
		ret = append(ret, NewErrorf("empty Zones (from %s)", whence))
	} else if err := c.set("Zones", v); err != nil {
		ret = append(ret, err)
	}

	return ret
}

func (this *parameter) populateRetryMax(c *config) error {
	return this.populateUInt64(c, "RetryMax")
}

func (this *parameter) populateRetryWaitMax(c *config) error {
	return this.populateUInt64(c, "RetryWaitMax")
}

func (this *parameter) populateRetryWaitMin(c *config) error {
	return this.populateUInt64(c, "RetryWaitMin")
}

func (this *parameter) populateAPIRootURL(c *config) error {
	// :TODO: validate URL format?
	return this.populateString(c, "APIRootURL")
}

func (this *parameter) populateAPIRequestTimeout(c *config) error {
	return this.populateUInt64(c, "APIRequestTimeout")
}

func (this *parameter) populateAPIRequestRateLimit(c *config) error {
	return this.populateUInt64(c, "APIRequestRateLimit")
}

func (this *parameter) populateTraceMode(c *config) error {
	// TraceMode _seems_ like an enum.
	// Known values so far:
	//
	// - unset (no trace)
	// - "all" (trace everything)
	// - "error" (only after errors)
	// - "api" (???)
	// - "http" (???)
	return this.populateString(c, "TraceMode")
}

func (p *parameter) populateString(c *config, key string) error {
	if val, whence, err := prioritizedParameterValue[string](p, c, key); err != nil {
		return err
	} else if v, ok := val.Get(); !ok {
		return nil // just not set; leave blank
	} else if v == "" {
		return NewErrorf("empty %s (from %s)", key, whence)
	} else {
		return c.set(key, v)
	}
}

func (p *parameter) populateUInt64(c *config, key string) error {
	if val, whence, err := prioritizedParameterValue[int64](p, c, key); err != nil {
		return err
	} else if v, ok := val.Get(); !ok {
		return nil // just not set; leave blank
	} else if v < 0 {
		return NewErrorf("negative %s (from %s): %d", key, whence, v)
	} else {
		return c.set(key, v)
	}
}

func prioritizedParameterValue[
	T any,
](
	p *parameter,
	c *config,
	k string,
) (
	maybeUninit[T],
	string,
	error,
) {
	var val maybeUninit[T]
	var whence string

	if p == nil {
		return val, whence, NewErrorf("nil parameter")
	} else if c == nil {
		return val, whence, NewErrorf("nil config")
	} else if val, whence, err := obtainFromStorage[T](&p.envp, k, "environment variable"); val.isSet() || err != nil {
		return val, whence, err
	} else if val, whence, err := obtainFromStorage[T](&p.argv, k, "command-line argument"); val.isSet() || err != nil {
		return val, whence, err
	} else if val, whence, err := obtainFromStorage[T](&p.hcl, k, "terraform configuration"); val.isSet() || err != nil {
		return val, whence, err
	} else {
		return obtainFromProfile[T](c, k, "profile")
	}
}

func obtainFromStorage[
	T any,
](
	s *storage,
	k string,
	msg ...string,
) (
	maybeUninit[T],
	string,
	error,
) {
	var val maybeUninit[T]
	whence := append(msg, "storage")[0]

	if s == nil {
		return val, whence, NewErrorf("nil %s", whence)
	} else if v, ok := s.get(k); !ok {
		return val, whence, nil
	} else if t, ok := v.(T); !ok {
		return val, whence, NewErrorf("invalid type for %s in %s: %T", k, whence, v)
	} else {
		return maybeUninit[T]{t, true}, whence, nil
	}
}

func obtainFromProfile[
	T any,
](
	c *config,
	k string,
	msg ...string,
) (
	maybeUninit[T],
	string,
	error,
) {
	var val maybeUninit[T]
	whence := fmt.Sprintf("%s %s", append(msg, "profile")[0], k)

	if c == nil {
		return val, whence, NewErrorf("nil config")
	} else if v, err := c.get("Profile"); err != nil {
		return val, whence, err
	} else if w, ok := v.Get(); !ok {
		// profile not set; ok unspecified
		return val, whence, nil
	} else if p, ok := w.(*Profile); !ok {
		return val, whence, NewErrorf("invalid profile in config: %v", v)
	} else if v, ok := p.Get(k); !ok {
		// profile does not have this key; ok unspecified
		return val, whence, nil
	} else if w, ok := v.(T); !ok {
		return val, whence, NewErrorf("invalid type for %s in %s: %T", k, whence, v)
	} else {
		return maybeUninit[T]{w, true}, whence, nil
	}
}

func (m *maybeUninit[T]) from(src func() (T, bool)) {
	if m != nil {
		value, set := src()
		*m = maybeUninit[T]{value, set}
	}
}

func (m *maybeUninit[T]) initialize(v T) {
	m.from(func() (T, bool) {
		return v, true
	})
}

func (m *maybeUninit[T]) isSet() bool { return m.set }

func (m *maybeUninit[T]) String() string {
	// This is used by flag package
	if m == nil {
		panic("nil dereference")
	}
	if v, ok := m.Get(); ok {
		return fmt.Sprintf("%v", v)
	} else {
		return ""
	}
}

func (m *maybeUninit[T]) Get() (T, bool) {
	var zero T

	if m == nil {
		return zero, false
	} else {
		return m.value, m.set
	}
}

func (m *maybeUninit[T]) Set(s string) error {
	// This is used by flag package

	switch m := any(m).(type) {
	case *maybeUninit[string]:
		m.initialize(s)

	case *maybeUninit[int64]:
		if v, err := strconv.ParseInt(s, 0, 64); err != nil {
			return err
		} else {
			m.initialize(v)
		}

	case *maybeUninit[[]string]:
		// This behaviour mimics usacloud's --zones= option
		r := strings.NewReader(s)
		c := csv.NewReader(r)
		if v, err := c.Read(); err != nil {
			return err
		} else {
			m.initialize(v)
		}

	default:
		// Should be unreachable because of the constraint,
		// but good practice to guard anyway.
		panic("unsupported type")
	}
	return nil
}

func (c *config) set(k string, v any) error {
	if c == nil {
		return NewErrorf("nil config")
	} else {
		(*c)[k] = v
		return nil
	}
}

func (c *config) get(k string) (maybeUninit[any], error) {
	var ret maybeUninit[any]

	if c == nil {
		return ret, NewErrorf("nil config")
	} else if v, ok := (*c)[k]; !ok {
		return ret, nil
	} else {
		ret.initialize(v)
		return ret, nil
	}
}

func (s *storage) get(k string) (any, bool) {
	switch k {
	case "profileName":
		return s.profileName.Get()

	case "PrivateKeyPEMPath":
		return s.privateKeyPath.Get()

	case "AccessToken":
		return s.accessToken.Get()

	case "AccessTokenSecret":
		return s.accessTokenSecret.Get()

	case "Zone":
		return s.zone.Get()

	case "DefaultZone":
		return s.defaultZone.Get()

	case "Zones":
		return s.zones.Get()

	case "RetryMax":
		return s.retryMax.Get()

	case "RetryWaitMax":
		return s.retryWaitMax.Get()

	case "RetryWaitMin":
		return s.retryWaitMin.Get()

	case "APIRootURL":
		return s.apiRootURL.Get()

	case "APIRequestTimeout":
		return s.apiRequestTimeout.Get()

	case "APIRequestRateLimit":
		return s.apiRequestRateLimit.Get()

	case "TraceMode":
		return s.traceMode.Get()

	default:
		panic("unknown key: " + k)
	}
}

var _ flag.Value = (*maybeUninit[string])(nil)
var _ flag.Value = (*maybeUninit[int64])(nil)
var _ flag.Value = (*maybeUninit[[]string])(nil)
