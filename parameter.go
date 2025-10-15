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
	"net/http/httptest"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
)

type option[T any] struct {
	some T

	// This is intentionally flipped from Rust's Option<T>
	// to make the zero value mean "not set".
	set bool
}

// This is just Result<Option<T>, E>
type resultOption[T any] struct {
	option[T]
	err error
}

type storage struct {
	profileName         option[string]
	privateKeyPath      option[string]
	accessToken         option[string]
	accessTokenSecret   option[string]
	zone                option[string]
	defaultZone         option[string]
	zones               option[[]string]
	retryMax            option[int64]
	retryWaitMax        option[int64]
	retryWaitMin        option[int64]
	apiRootURL          option[string]
	apiRequestTimeout   option[int64]
	apiRequestRateLimit option[int64]
	traceMode           option[string]
	mockServer          option[*httptest.Server]
	userAgent           option[string]
}

// :INTERNAL: it is intentional that this is not a struct
// This is also not JSONable because it contains functions etc.
type config map[string]any

type parameter struct {
	profileOp *ProfileOp
	envp      storage
	argv      storage
	hcl       storage
	dynamic   storage
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
	p.hcl.apiRequestTimeout.from(config.LookupClientConfigAPIRequestTimeout)
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
	ret := make([]error, 0, 18) // <- 18 is the # of `append` calls below

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
	ret = append(ret, p.populateMockServer(c))
	ret = append(ret, p.populateUserAgent(c))

	if result := obtainFromConfig[string](c, "AccessToken"); result.isSome() {
		// Take that,

	} else if result := obtainFromConfig[string](c, "PrivateKeyPEMPath"); result.isSome() {
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
	var profileName option[string]

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

	} else if result := obtainFromConfig[string](c, "PrivateKeyPEMPath"); result.isErr() {
		return result.error()

	} else if v, ok := result.some(); !ok {
		return nil // just not set

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
	var val option[[]string]

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

	} else if whence, result := obtainFromProfile[[]any](c, "Zones", "profile"); result.isErr() {
		ret = append(ret, result.error())

	} else if v, ok := result.some(); !ok {
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

func (p *parameter) populateMockServer(c *config) error {
	if _, result := prioritizedParameterValue[*httptest.Server](p, c, "MockServer"); result.isErr() {
		return result.error()

	} else if v, ok := result.some(); !ok {
		return nil // just not set; leave blank

	} else if v == nil {
		return nil // avoid SEGV

	} else {
		c.set("MockServer", v)
		return nil
	}
}

func (p *parameter) populateUserAgent(c *config) error {
	return p.populateString(c, "UserAgent")
}

func (p *parameter) populateString(c *config, key string) error {
	if whence, result := prioritizedParameterValue[string](p, c, key); result.isErr() {
		return result.error()

	} else if v, ok := result.some(); !ok {
		return nil // just not set; leave blank

	} else if v == "" {
		return NewErrorf("empty %s (from %s)", key, whence)

	} else {
		return c.set(key, v)
	}
}

func (p *parameter) populateUInt64(c *config, key string) error {
	if whence, result := prioritizedParameterValue[int64](p, c, key); result.isErr() {
		return result.error()

	} else if v, ok := result.some(); !ok {
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
	string,
	resultOption[T],
) {
	var whence string

	if p == nil {
		return whence, resultOptionErr[T](NewErrorf("nil parameter"))

	} else if c == nil {
		return whence, resultOptionErr[T](NewErrorf("nil config"))

	} else if whence, result := obtainFromStorage[T](&p.envp, k, "environment variable"); result.isSome() {
		return whence, result

	} else if whence, result := obtainFromStorage[T](&p.argv, k, "command-line argument"); result.isSome() {
		return whence, result

	} else if whence, result := obtainFromStorage[T](&p.hcl, k, "terraform configuration"); result.isSome() {
		return whence, result

	} else if whence, result := obtainFromProfile[T](c, k, "profile"); result.isSome() {
		return whence, result

	} else if whence, result := obtainFromStorage[T](&p.dynamic, k, "on-the-fly"); result.isSome() {
		return whence, result

	} else {
		return obtainFromStorage[T](&defaults, k, "defaults")
	}
}

func obtainFromStorage[
	T any,
](
	s *storage,
	k string,
	whence string,
) (
	string,
	resultOption[T],
) {
	if s == nil {
		return whence, resultOptionErr[T](NewErrorf("nil %s", whence))

	} else if v, ok := s.get(k); !ok {
		return whence, resultOptionNone[T]()

	} else if t, ok := v.(T); !ok {
		return whence, resultOptionErr[T](NewErrorf("invalid type for %s in %s: %T", k, whence, v))

	} else {
		return whence, resultOptionSome(t)
	}
}

func obtainFromProfile[
	T any,
](
	c *config,
	k string,
	msg string,
) (
	string,
	resultOption[T],
) {
	whence := fmt.Sprintf("%s %s", msg, k)

	if c == nil {
		return whence, resultOptionErr[T](NewErrorf("nil config"))

	} else if result := obtainFromConfig[*Profile](c, "Profile"); result.isErr() {
		return whence, resultOptionErr[T](result.error())

	} else if p, ok := result.some(); !ok {
		// profile not set; ok unspecified
		return whence, resultOptionNone[T]()

	} else if v, ok := p.Get(k); !ok {
		// profile does not have this key; ok unspecified
		return whence, resultOptionNone[T]()

	} else if w, ok := v.(T); !ok {
		return whence, resultOptionErr[T](NewErrorf("invalid type for %s in %s: %T", k, whence, v))

	} else {
		return whence, resultOptionSome(w)
	}
}

func obtainFromConfig[T any](c *config, k string) resultOption[T] {
	if c == nil {
		return resultOptionErr[T](NewErrorf("nil config"))

	} else if v, ok := (*c)[k]; !ok {
		return resultOptionNone[T]()

	} else if w, ok := v.(T); !ok {
		return resultOptionErr[T](NewErrorf("invalid type for %s in config: %T", k, v))

	} else {
		return resultOptionSome(w)
	}
}

func (m *option[T]) from(src func() (T, bool)) {
	if m != nil {
		value, set := src()
		*m = option[T]{value, set}
	}
}

func (m *option[T]) initialize(v T) {
	m.from(func() (T, bool) {
		return v, true
	})
}

func (m *option[T]) isSome() bool { return m.set }

func (m *option[T]) String() string {
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

func (m *option[T]) Get() (T, bool) {
	var zero T

	if m == nil {
		return zero, false

	} else {
		return m.some, m.set
	}
}

func (m *option[T]) Set(s string) error {
	// This is used by flag package

	switch m := any(m).(type) {
	case *option[string]:
		m.initialize(s)

	case *option[int64]:
		if v, err := strconv.ParseInt(s, 0, 64); err != nil {
			return err

		} else {
			m.initialize(v)
		}

	case *option[[]string]:
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

func (r *resultOption[T]) isErr() bool     { return r.err != nil }
func (r *resultOption[T]) error() error    { return r.err }
func (r *resultOption[T]) ok() *option[T]  { return &r.option }
func (r *resultOption[T]) some() (T, bool) { return r.ok().Get() }
func (r *resultOption[T]) unwrap_or(zero T) T {
	if ret, ok := r.some(); ok {
		return ret
	} else {
		return zero
	}
}

func resultOptionErr[T any](err error) resultOption[T] { return resultOption[T]{err: err} }
func resultOptionNone[T any]() resultOption[T]         { return resultOption[T]{} }
func resultOptionSome[T any](some T) resultOption[T] {
	return resultOption[T]{option: option[T]{some: some, set: true}}
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

	case "MockServer":
		return s.mockServer.Get()

	case "UserAgent":
		return s.userAgent.Get()

	default:
		panic("unknown key: " + k)
	}
}

var _ flag.Value = (*option[string])(nil)
var _ flag.Value = (*option[int64])(nil)
var _ flag.Value = (*option[[]string])(nil)

// values copied from: sacloud/api-client-go/options.go:defaultOption
var defaults = storage{
	// absent keys are "not defaults"
	profileName:         option[string]{set: true, some: "default"},
	retryMax:            option[int64]{set: true, some: 10},
	retryWaitMax:        option[int64]{set: true, some: 64},
	retryWaitMin:        option[int64]{set: true, some: 1},
	apiRequestTimeout:   option[int64]{set: true, some: 300},
	apiRequestRateLimit: option[int64]{set: true, some: 5},
	userAgent: option[string]{set: true, some: fmt.Sprintf(
		// :INTENTIONAL: keeping "api-client-go" here for backward compatibility
		"api-client-go/v%s (%s/%s; +https://github.com/sacloud/http-client-go)",
		Version,
		runtime.GOOS,
		runtime.GOARCH,
	)},
}
