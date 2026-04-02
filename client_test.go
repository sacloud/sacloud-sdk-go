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

//go:debug rsa1024min=0
package saclient_test

import (
	"context"
	"crypto/rand"
	"encoding/pem"
	"flag"
	"fmt"
	"io"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	old "github.com/sacloud/api-client-go"
	saht "github.com/sacloud/go-http"
	. "github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/suite"
)

type providerModel struct {
	Profile             types.String `tfsdk:"profile"`
	AccessToken         types.String `tfsdk:"token"`
	AccessTokenSecret   types.String `tfsdk:"secret"`
	Zone                types.String `tfsdk:"zone"`
	Zones               types.List   `tfsdk:"zones"`
	DefaultZone         types.String `tfsdk:"default_zone"`
	APIRootURL          types.String `tfsdk:"api_root_url"`
	RetryMax            types.Int64  `tfsdk:"retry_max"`
	RetryWaitMax        types.Int64  `tfsdk:"retry_wait_max"`
	RetryWaitMin        types.Int64  `tfsdk:"retry_wait_min"`
	APIRequestTimeout   types.Int64  `tfsdk:"api_request_timeout"`
	APIRequestRateLimit types.Int64  `tfsdk:"api_request_rate_limit"`
	TraceMode           types.String `tfsdk:"trace"`
}

var _ TerraformProviderInterface = (*providerModel)(nil)

func (p *providerModel) LookupClientConfigProfileName() (string, bool) {
	return p.Profile.ValueString(), !p.Profile.IsNull() && !p.Profile.IsUnknown()
}

func (p *providerModel) LookupClientConfigPrivateKeyPath() (string, bool) {
	// Not supported in this test model
	return "", false
}

func (p *providerModel) LookupClientConfigServicePrincipalID() (string, bool) {
	// Not supported in this test model
	return "", false
}

func (p *providerModel) LookupClientConfigServicePrincipalKeyID() (string, bool) {
	// Not supported in this test model
	return "", false
}

func (p *providerModel) LookupClientConfigAccessToken() (string, bool) {
	return p.AccessToken.ValueString(), !p.AccessToken.IsNull() && !p.AccessToken.IsUnknown()
}

func (p *providerModel) LookupClientConfigAccessTokenSecret() (string, bool) {
	return p.AccessTokenSecret.ValueString(), !p.AccessTokenSecret.IsNull() && !p.AccessTokenSecret.IsUnknown()
}

func (p *providerModel) LookupClientConfigZone() (string, bool) {
	return p.Zone.ValueString(), !p.Zone.IsNull() && !p.Zone.IsUnknown()
}

func (p *providerModel) LookupClientConfigDefaultZone() (string, bool) {
	return p.DefaultZone.ValueString(), !p.DefaultZone.IsNull() && !p.DefaultZone.IsUnknown()
}

func (p *providerModel) LookupClientConfigZones() ([]string, bool) {
	if p.Zones.IsNull() || p.Zones.IsUnknown() {
		return nil, false
	} else {
		vals := p.Zones.Elements()
		result := make([]string, 0, len(vals))
		for _, v := range vals {
			if str, ok := v.(types.String); ok {
				if !str.IsNull() && !str.IsUnknown() {
					result = append(result, str.ValueString())
				}
			}
		}
		return result, true
	}
}

func (p *providerModel) LookupClientConfigRetryMax() (int64, bool) {
	return p.RetryMax.ValueInt64(), !p.RetryMax.IsNull() && !p.RetryMax.IsUnknown()
}

func (p *providerModel) LookupClientConfigRetryWaitMax() (int64, bool) {
	return p.RetryWaitMax.ValueInt64(), !p.RetryWaitMax.IsNull() && !p.RetryWaitMax.IsUnknown()
}

func (p *providerModel) LookupClientConfigRetryWaitMin() (int64, bool) {
	return p.RetryWaitMin.ValueInt64(), !p.RetryWaitMin.IsNull() && !p.RetryWaitMin.IsUnknown()
}

func (p *providerModel) LookupClientConfigAPIRootURL() (string, bool) {
	return p.APIRootURL.ValueString(), !p.APIRootURL.IsNull() && !p.APIRootURL.IsUnknown()
}

func (p *providerModel) LookupClientConfigAPIRequestTimeout() (int64, bool) {
	return p.APIRequestTimeout.ValueInt64(), !p.APIRequestTimeout.IsNull() && !p.APIRequestTimeout.IsUnknown()
}

func (p *providerModel) LookupClientConfigAPIRequestRateLimit() (int64, bool) {
	return p.APIRequestRateLimit.ValueInt64(), !p.APIRequestRateLimit.IsNull() && !p.APIRequestRateLimit.IsUnknown()
}

func (p *providerModel) LookupClientConfigTraceMode() (string, bool) {
	return p.TraceMode.ValueString(), !p.TraceMode.IsNull() && !p.TraceMode.IsUnknown()
}

// :FIXME: this does not cover any part of the implementation.
// Because this is nothing more than a copy & paste of the tested code itself.
var ua string = fmt.Sprintf(
	"api-client-go/v%s (%s/%s; +https://github.com/sacloud/saclient-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

type ClientTestSuite struct {
	suite.Suite
	XDG_CONFIG_HOME         *string
	SAKURACLOUD_PROFILE_DIR *string
	subject                 *Client
}

func TestClientTestSuite(t *testing.T) { suite.Run(t, new(ClientTestSuite)) }

//nolint:errcheck,gosec
func (s *ClientTestSuite) SetupSuite() {
	// Note that `s.T().TempDir()` is removed every time after a _test_, not afrer a suite.
	if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		s.XDG_CONFIG_HOME = &dir
	}
	if dir, ok := os.LookupEnv("SAKURACLOUD_PROFILE_DIR"); ok {
		s.SAKURACLOUD_PROFILE_DIR = &dir
	}
	dir, _ := os.MkdirTemp(os.TempDir(), "profile_test")
	os.Setenv("XDG_CONFIG_HOME", dir)
	os.Unsetenv("SAKURACLOUD_PROFILE_DIR")

	// create sample profiles
	os.MkdirAll(dir+"/usacloud/usacloud", 0o700)

	os.WriteFile(dir+"/usacloud/current", []byte("usacloud"), 0o600)
	os.WriteFile(dir+"/usacloud/usacloud/config.json",
		[]byte(fmt.Sprintf(`{
			"RetryMax":7,
			"TraceMode": "",
			"Zone":"usacloud",
			"PrivateKeyPEMPath":"%s/usacloud/usacloud/usamin.pem"
		}`, filepath.ToSlash(dir))),
		0o600,
	)

	garbage := make([]byte, 1024)
	rand.Read(garbage)
	os.MkdirAll(dir+"/usacloud/broken", 0o700)
	os.WriteFile(dir+"/usacloud/broken/config.json", garbage, 0o600)

	fp, _ := os.OpenFile(dir+"/usacloud/usacloud/usamin.pem", os.O_WRONLY|os.O_CREATE, 0o600)
	defer fp.Close()
	io.WriteString(fp, `
-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEAuvRcHR1FLlbL0hV/
oMcgGwCP4NS6HayccLsSXrBmXnVvG0YSntAUz9p+lMcOLm5/Ao+M7nJdWntHxgn4
REkS1wIDAQABAkAESOFrkWYqf7bAI9n+91FXDRY/EuEJGRGky8TKAsT12TUN7v/F
0G96JeBUUsH7ZHKMqyOui9SGypnR+6baR5kRAiEA94jnP7ZhRS0p8r0ScmjNs4Fl
qrcqI3CWC9LiIzVC3T0CIQDBWRf4m1vXt5fr77j85Zj28mNutiwXYQLJUKV5X7xZ
owIgCDcL7apg2gngrYSm2xMtWHq/5AWGKXzwDd5m0OJQoMUCIQCdbqERIddHr8s5
JonXClBiC32hIR6HrsspBsymJqjjxwIhAIyHd91UNdLPGxidqNp7c65sh1/BimtT
euIJBGkmzNop
-----END PRIVATE KEY-----
`)

	fp, _ = os.OpenFile(dir+"/another.pem", os.O_WRONLY|os.O_CREATE, 0o600)
	defer fp.Close()
	pem.Encode(fp, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: []byte("dummy"),
	})

	os.MkdirAll(dir+"/usacloud/withNull", 0o700)
	os.WriteFile(dir+"/usacloud/withNull/config.json",
		[]byte(`{
			"Zone": null,
			"Zones": null
		}`),
		0o600,
	)
}

//nolint:errcheck,gosec
func (s *ClientTestSuite) TearDownSuite() {
	if s.XDG_CONFIG_HOME != nil {
		os.Setenv("XDG_CONFIG_HOME", *s.XDG_CONFIG_HOME)
	} else {
		os.Unsetenv("XDG_CONFIG_HOME")
	}
	if s.SAKURACLOUD_PROFILE_DIR != nil {
		os.Setenv("SAKURACLOUD_PROFILE_DIR", *s.SAKURACLOUD_PROFILE_DIR)
	} else {
		os.Unsetenv("SAKURACLOUD_PROFILE_DIR")
	}
}

func (s *ClientTestSuite) SetupTest() {
	s.subject = new(Client)

	// AD HOC: easy test
	_ = s.subject.CompatSettingsFromAPIClientParams(
		"",
		old.WithDisableEnv(true),
		old.WithDisableProfile(true),
	)
}

// #nosec G101 -- This is only a test
func (s *ClientTestSuite) TestCLI() {
	e := s.subject.FlagSet(flag.PanicOnError).Parse([]string{
		"--secret=bar",
		"--token=foo",
		"--trace",
		"--zones=foo,\", bar\"",
	})
	s.NoError(e)
	e = s.subject.Populate()
	s.NoError(e)
	s.Equal(map[string]any{
		"AccessToken":         "foo",
		"AccessTokenSecret":   "bar",
		"APIRequestRateLimit": int64(5),
		"APIRequestTimeout":   time.Duration(300) * time.Second,
		"AuthPreference":      "basic",
		"RetryMax":            int64(10),
		"RetryWaitMax":        int64(64),
		"RetryWaitMin":        int64(1),
		"TokenEndpoint":       "https://secure.sakura.ad.jp/cloud/api/iam/1.0/service-principals/oauth2/token",
		"TraceMode":           "all",
		"UserAgent":           ua,
		"Zones": []string{
			"foo",
			", bar",
		},
	}, s.subject.JSON())
}

// #nosec G101 -- This is only a test
func (s *ClientTestSuite) TestEnviron() {
	s.Run("SAKURACLOUD_", func() {
		subject := s.subject.Dup().(*Client)
		e := subject.SetEnviron([]string{
			"SAKURACLOUD_ACCESS_TOKEN_SECRET=bar",
			"SAKURACLOUD_ACCESS_TOKEN=foo",
			"SAKURACLOUD_API_REQUEST_RATE_LIMIT=20",
			"SAKURACLOUD_API_REQUEST_TIMEOUT=30",
			"SAKURACLOUD_API_ROOT_URL=https://api.example.com",
			"SAKURACLOUD_PRIVATE_KEY_PATH=" + os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
			"SAKURACLOUD_PRIVATE_KEY=dummy-private-key",
			"SAKURACLOUD_RETRY_MAX=", // <= empty value
			"SAKURACLOUD_RETRY_WAIT_MAX=7",
			"SAKURACLOUD_RETRY_WAIT_MIN=5",
			"SAKURACLOUD_TOKEN_ENDPOINT=https://example.com/oauth2/token",
			"SAKURACLOUD_TRACE=error",
			"SAKURACLOUD_ZONE=foo",
			"SAKURACLOUD_ZONES=foo,\", bar\"",
			"XDG_CONFIG_HOME=" + os.Getenv("XDG_CONFIG_HOME"),
		})
		s.NoError(e)
		e = subject.Populate()
		s.NoError(e)
		s.Equal(map[string]any{
			"AccessToken":         "foo",
			"AccessTokenSecret":   "bar",
			"APIRequestRateLimit": int64(20),
			"APIRequestTimeout":   time.Duration(30) * time.Second,
			"APIRootURL":          "https://api.example.com",
			"AuthPreference":      "bearer",
			"PrivateKey":          "dummy-private-key",
			"PrivateKeyPEMPath":   os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
			"RetryMax":            int64(10), // <= default value instead of zero
			"RetryWaitMax":        int64(7),
			"RetryWaitMin":        int64(5),
			"TokenEndpoint":       "https://example.com/oauth2/token",
			"TraceMode":           "error",
			"UserAgent":           ua,
			"Zone":                "foo",
			"Zones": []string{
				"foo",
				", bar",
			},
		}, subject.JSON())
	})

	s.Run("SAKURA_", func() {
		subject := s.subject.Dup().(*Client)
		e := subject.SetEnviron([]string{
			"SAKURA_ACCESS_TOKEN_SECRET=bar",
			"SAKURA_ACCESS_TOKEN=foo",
			"SAKURA_RATE_LIMIT=20",
			"SAKURA_API_REQUEST_TIMEOUT=30",
			"SAKURA_API_ROOT_URL=https://api.example.com",
			"SAKURA_PRIVATE_KEY_PATH=" + os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
			"SAKURA_PRIVATE_KEY=dummy-private-key",
			"SAKURA_RETRY_MAX=", // <= empty value
			"SAKURA_RETRY_WAIT_MAX=7",
			"SAKURA_RETRY_WAIT_MIN=5",
			"SAKURA_TOKEN_ENDPOINT=https://example.com/oauth2/token",
			"SAKURA_TRACE=error",
			"SAKURA_ZONE=foo",
			"SAKURA_ZONES=foo,\", bar\"",
			"XDG_CONFIG_HOME=" + os.Getenv("XDG_CONFIG_HOME"),
		})
		s.NoError(e)
		e = subject.Populate()
		s.NoError(e)
		s.Equal(map[string]any{
			"AccessToken":         "foo",
			"AccessTokenSecret":   "bar",
			"APIRequestRateLimit": int64(20),
			"APIRequestTimeout":   time.Duration(30) * time.Second,
			"APIRootURL":          "https://api.example.com",
			"AuthPreference":      "bearer",
			"PrivateKey":          "dummy-private-key",
			"PrivateKeyPEMPath":   os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
			"RetryMax":            int64(10), // <= default value instead of zero
			"RetryWaitMax":        int64(7),
			"RetryWaitMin":        int64(5),
			"TokenEndpoint":       "https://example.com/oauth2/token",
			"TraceMode":           "error",
			"UserAgent":           ua,
			"Zone":                "foo",
			"Zones": []string{
				"foo",
				", bar",
			},
		}, subject.JSON())
	})

	s.Run("both", func() {
		subject := s.subject.Dup().(*Client)
		e := subject.SetEnviron([]string{
			"SAKURA_ACCESS_TOKEN_SECRET=bar",
			"SAKURA_ACCESS_TOKEN=foo",
			"SAKURA_RATE_LIMIT=20",
			"SAKURA_API_REQUEST_TIMEOUT=30",
			"SAKURA_API_ROOT_URL=https://api.example.com",
			"SAKURA_PRIVATE_KEY_PATH=" + os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
			"SAKURA_PRIVATE_KEY=dummy-private-key",
			"SAKURA_RETRY_MAX=", // <= empty value
			"SAKURA_RETRY_WAIT_MAX=7",
			"SAKURA_RETRY_WAIT_MIN=5",
			"SAKURA_TOKEN_ENDPOINT=https://example.com/oauth2/token",
			"SAKURA_TRACE=error",
			"SAKURA_ZONE=foo",
			"SAKURA_ZONES=foo,\", bar\"",
			"SAKURACLOUD_ACCESS_TOKEN_SECRET=baz",
			"SAKURACLOUD_ACCESS_TOKEN=quux",
			"SAKURACLOUD_API_REQUEST_RATE_LIMIT=64",
			"SAKURACLOUD_API_REQUEST_TIMEOUT=128",
			"SAKURACLOUD_API_ROOT_URL=https://sample.example.com",
			"SAKURACLOUD_PRIVATE_KEY_PATH=" + os.Getenv("XDG_CONFIG_HOME") + "/another.pem",
			"SAKURACLOUD_PRIVATE_KEY=nonexistent-key",
			"SAKURACLOUD_RETRY_MAX=1024",
			"SAKURACLOUD_RETRY_WAIT_MAX=512",
			"SAKURACLOUD_RETRY_WAIT_MIN=32",
			"SAKURACLOUD_TOKEN_ENDPOINT=https://another.example.com/oauth2/token",
			"SAKURACLOUD_TRACE=all",
			"SAKURACLOUD_ZONE=baz",
			"SAKURACLOUD_ZONES=baz,\", quux\"",
			"XDG_CONFIG_HOME=" + os.Getenv("XDG_CONFIG_HOME"),
		})
		s.NoError(e)
		e = subject.Populate()
		s.NoError(e)
		s.Equal(map[string]any{
			"AccessToken":         "foo",
			"AccessTokenSecret":   "bar",
			"APIRequestRateLimit": int64(20),
			"APIRequestTimeout":   time.Duration(30) * time.Second,
			"APIRootURL":          "https://api.example.com",
			"AuthPreference":      "bearer",
			"PrivateKey":          "dummy-private-key",
			"PrivateKeyPEMPath":   os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
			"RetryMax":            int64(1024), // <= SAKURACLOUD_ wins here
			"RetryWaitMax":        int64(7),
			"RetryWaitMin":        int64(5),
			"TokenEndpoint":       "https://example.com/oauth2/token",
			"TraceMode":           "error",
			"UserAgent":           ua,
			"Zone":                "foo",
			"Zones": []string{
				"foo",
				", bar",
			},
		}, subject.JSON())
	})
}

// #nosec G101 -- This is only a test
func (s *ClientTestSuite) TestTerraform() {
	_ = s.subject.SettingsFromTerraformProvider(&providerModel{
		AccessToken:         types.StringValue("foo"),
		AccessTokenSecret:   types.StringValue("bar"),
		APIRequestRateLimit: types.Int64Value(20),
		APIRequestTimeout:   types.Int64Value(30),
		APIRootURL:          types.StringValue("https://api.example.com"),
		DefaultZone:         types.StringValue("foo"),
		Profile:             types.StringValue("usacloud"),
		RetryMax:            types.Int64Value(3),
		RetryWaitMax:        types.Int64Value(7),
		RetryWaitMin:        types.Int64Value(5),
		TraceMode:           types.StringValue("error"),
		Zone:                types.StringValue("foo"),
		Zones: types.ListValueMust(types.StringType, []attr.Value{
			types.StringValue("foo"),
			types.StringValue(", bar"),
		}),
	})
	e := s.subject.Populate()
	s.NoError(e)
	s.Equal(map[string]any{
		"AccessToken":         "foo",
		"AccessTokenSecret":   "bar",
		"APIRequestRateLimit": int64(20),
		"APIRequestTimeout":   time.Duration(30) * time.Second,
		"APIRootURL":          "https://api.example.com",
		"AuthPreference":      "basic",
		"DefaultZone":         "foo",
		"RetryMax":            int64(3),
		"RetryWaitMax":        int64(7),
		"RetryWaitMin":        int64(5),
		"TokenEndpoint":       "https://secure.sakura.ad.jp/cloud/api/iam/1.0/service-principals/oauth2/token",
		"TraceMode":           "error",
		"UserAgent":           ua,
		"Zone":                "foo",
		"Zones": []string{
			"foo",
			", bar",
		},
	}, s.subject.JSON())
}

func (s *ClientTestSuite) TestNoProfile() {
	current := os.Getenv("XDG_CONFIG_HOME")
	_ = os.Setenv("XDG_CONFIG_HOME", s.T().TempDir())
	defer func() { _ = os.Setenv("XDG_CONFIG_HOME", current) }()

	_ = os.Setenv("SAKURACLOUD_PROFILE_DIR", s.T().TempDir())
	defer func() { _ = os.Unsetenv("SAKURACLOUD_PROFILE_DIR") }()

	e := s.subject.FlagSet(flag.PanicOnError).Parse([]string{
		"--secret=bar",
		"--token=foo",
	})
	s.NoError(e)
	e = s.subject.Populate()
	s.NoError(e)
	s.Nil(s.subject.Profile())
	s.Equal("bar", s.subject.JSON()["AccessTokenSecret"])
}

func (s *ClientTestSuite) TestDynamic() {
	api, err := s.subject.DupWith(WithTraceMode("all"))
	s.NoError(err)
	err = api.Populate()
	s.NoError(err)
	subject := api.(*Client)
	j := subject.JSON()
	s.Equal("all", j["TraceMode"])

	api, err = s.subject.DupWith(WithZone("is1c"))
	s.NoError(err)
	err = api.Populate()
	s.NoError(err)
	cfg, err := api.EndpointConfig()
	s.NoError(err)
	s.Equal("is1c", cfg.Zone)
}

// #nosec G101 -- This is only a test
func (s *ClientTestSuite) TestDynamicUsinfgClientParams() {
	err := s.subject.CompatSettingsFromAPIClientParams(
		"https://api.example.com",
		old.WithDisableEnv(true),
		old.WithDisableProfile(true),
		old.WithApiKeys("foo", "bar"),
		old.WithUserAgent(ua),
		old.WithOptions(&old.Options{
			HttpRequestTimeout:   300,
			HttpRequestRateLimit: 100,
			RetryMax:             30,
			RetryWaitMax:         70,
			RetryWaitMin:         50,
		}),
	)
	s.NoError(err)
	err = s.subject.Populate()
	s.NoError(err)
	s.Equal(map[string]any{
		"AccessToken":         "foo",
		"AccessTokenSecret":   "bar",
		"APIRequestRateLimit": int64(100),
		"APIRequestTimeout":   time.Duration(300) * time.Second,
		"APIRootURL":          "https://api.example.com",
		"RetryMax":            int64(30),
		"RetryWaitMax":        int64(70),
		"RetryWaitMin":        int64(50),
		"TokenEndpoint":       "https://secure.sakura.ad.jp/cloud/api/iam/1.0/service-principals/oauth2/token",
		"UserAgent":           ua,
	}, s.subject.JSON())
}

// #nosec G101 -- This is only a test
func (s *ClientTestSuite) TestDynamicUsinfgClientOptions() {
	err := s.subject.CompatSettingsFromAPIClientOptions(
		&old.Options{AccessToken: "foo"},
		&old.Options{AccessTokenSecret: "bar"},
		&old.Options{Gzip: true},
		&old.Options{HttpRequestTimeout: 300},
		&old.Options{HttpRequestRateLimit: 100},
		&old.Options{RetryMax: 30},
		&old.Options{RetryWaitMax: 70},
		&old.Options{RetryWaitMin: 50},
		&old.Options{UserAgent: ua},
		&old.Options{Trace: true},
		&old.Options{TraceOnlyError: true},
		&old.Options{RequestCustomizers: []saht.RequestCustomizer{
			func(req *http.Request) error {
				req.Header.Set("X-Custom-Header", "custom-value")
				return nil
			},
		}},
		&old.Options{CheckRetryFunc: func(ctx context.Context, resp *http.Response, err error) (bool, error) {
			if resp != nil && resp.StatusCode == 502 {
				return true, nil
			}
			return false, nil
		}},
		&old.Options{CheckRetryStatusCodes: []int{502}},
	)
	s.NoError(err)
	err = s.subject.Populate()
	s.NoError(err)
	s.Equal(map[string]any{
		"AccessToken":         "foo",
		"AccessTokenSecret":   "bar",
		"APIRequestRateLimit": int64(100),
		"APIRequestTimeout":   time.Duration(300) * time.Second,
		"RetryMax":            int64(30),
		"RetryWaitMax":        int64(70),
		"RetryWaitMin":        int64(50),
		"TokenEndpoint":       "https://secure.sakura.ad.jp/cloud/api/iam/1.0/service-principals/oauth2/token",
		"TraceMode":           "error",
		"UserAgent":           ua,
	}, s.subject.JSON())
}

func (s *ClientTestSuite) TestProfileName() {
	s.Run("Found sane", func() {
		subject := s.subject.Dup()
		_ = subject.CompatSettingsFromAPIClientParams("", old.WithDisableProfile(false))
		dir, name := subject.ProfileName()
		s.NotNil(name)
		s.NotNil(dir)
		s.Equal(filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "usacloud"), *dir)
		s.Equal("usacloud", *name)
	})

	s.Run("Found broken", func() {
		subject := s.subject.Dup()
		_ = subject.CompatSettingsFromAPIClientParams("", old.WithDisableProfile(false))
		err := subject.FlagSet(flag.PanicOnError).Parse([]string{"--profile=broken"})
		s.NoError(err)
		err = subject.Populate()
		s.ErrorContains(err, "failed to parse")
		dir, name := subject.ProfileName()
		s.Equal("broken", *name)
		s.Equal(filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "usacloud"), *dir)
	})

	s.Run("Specified, but not found", func() {
		subject := s.subject.Dup()
		_ = subject.CompatSettingsFromAPIClientParams("", old.WithDisableProfile(false))
		err := subject.FlagSet(flag.PanicOnError).Parse([]string{"--profile=nonexistent"})
		s.NoError(err)
		err = subject.Populate()
		s.ErrorContains(err, "failed to open")
		dir, name := subject.ProfileName()
		s.Equal("nonexistent", *name)
		s.Equal(filepath.Join(os.Getenv("XDG_CONFIG_HOME"), "usacloud"), *dir)
	})

	s.Run("unspecified", func() {
		expected := s.T().TempDir()
		current := os.Getenv("XDG_CONFIG_HOME")
		_ = os.Setenv("XDG_CONFIG_HOME", s.T().TempDir())
		defer func() { _ = os.Setenv("XDG_CONFIG_HOME", current) }()

		_ = os.Setenv("SAKURACLOUD_PROFILE_DIR", expected)
		defer func() { _ = os.Unsetenv("SAKURACLOUD_PROFILE_DIR") }()

		var subject Client
		err := subject.Populate()
		s.NoError(err)
		dir, name := subject.ProfileName()
		s.Equal(expected, *dir)
		s.Nil(name)
	})
}

func (s *ClientTestSuite) TestProfileWithNullValue() {
	var subject *Client = s.subject.Dup().(*Client)
	e := subject.CompatSettingsFromAPIClientParams("", old.WithDisableProfile(false))
	s.NoError(e)

	e = subject.FlagSet(flag.PanicOnError).Parse([]string{"--profile=withNull"})
	s.NoError(e)

	e = subject.Populate()
	s.NoError(e)

	actual := subject.JSON()
	s.NotContains(actual, "Zone")
	s.NotContains(actual, "Zones")
}

func (s *ClientTestSuite) TestPrecedence() {
	// Setup
	subject := s.subject.Dup().(*Client)

	if e := subject.FlagSet(flag.PanicOnError).Parse([]string{
		"--secret=argv@secret",
		"--trace", // <= TraceMode = "all"
	}); !s.NoError(e) {
		return
	}

	if e := subject.SetEnviron([]string{
		"SAKURACLOUD_ACCESS_TOKEN=envp@token",
		"SAKURACLOUD_ACCESS_TOKEN_SECRET=envp@secret",
		"SAKURACLOUD_ZONE=envp@zone",
		"SAKURACLOUD_TRACE=envp@trace",
		"XDG_CONFIG_HOME=" + os.Getenv("XDG_CONFIG_HOME"), // to load profile
	}); !s.NoError(e) {
		return
	}

	if e := subject.SettingsFromTerraformProvider(&providerModel{
		AccessToken:       types.StringValue("hcl@token"),
		AccessTokenSecret: types.StringValue("hcl@secret"),
		TraceMode:         types.StringValue("hcl@trace"),
	}); !s.NoError(e) {
		return
	}

	if e := subject.CompatSettingsFromAPIClientOptions(&old.Options{
		TraceOnlyError: true, // <= TraceMode = "error"
	}); !s.NoError(e) {
		return
	}

	if e := subject.CompatSettingsFromAPIClientParams("", old.WithDisableProfile(false)); !s.NoError(e) {
		return
	}

	if e := subject.Populate(); !s.NoError(e) {
		return
	}

	// Test

	// argv:exists, envp:exists, hcl:exists, profile:missing, dynamic: exists, default:missing
	// => dynamic wins
	s.Equal("error", subject.JSON()["TraceMode"])

	// argv:exists, envp:exists, hcl:exists, profile:missing, dynamic: missing, default:missing
	// => argv wins
	s.Equal("argv@secret", subject.JSON()["AccessTokenSecret"])

	// argv:missing, envp:exists, hcl:exists, profile:missing, dynamic: missing, default:missing
	// => hcl wins
	s.Equal("hcl@token", subject.JSON()["AccessToken"])

	// argv:missing, envp:exists, hcl:missing, profile:exists, dynamic: missing, default:missing
	// => envp wins
	s.Equal("envp@zone", subject.JSON()["Zone"])

	// argv:missing, envp:missing, hcl:missing, profile:exists, dynamic: missing, default:exists
	// => profile wins
	s.Equal(int64(7), subject.JSON()["RetryMax"])

	// argv:missing, envp:missing, hcl:missing, profile:missing, dynamic: missing, default:exists
	// => default wins
	s.Equal(int64(5), subject.JSON()["APIRequestRateLimit"])

	// argv:missing, envp:missing, hcl:missing, profile:missing, dynamic: missing, default:missing
	// => no key
	s.NotContains(subject.JSON(), "SomeNonexistentKey")
}

func (s *ClientTestSuite) authPreferenceMockServer() (
	hdr http.Header,
	req *http.Request,
	svr *httptest.Server,
) {
	var err error
	hdr = make(http.Header)
	svr = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		maps.Copy(hdr, r.Header)
		w.WriteHeader(http.StatusAccepted)
	}))
	req, err = http.NewRequestWithContext(s.T().Context(), "GET", svr.URL, nil)
	s.NoError(err)
	return
}

func (s *ClientTestSuite) TestAuthPreference() {
	// Setup
	e := s.subject.FlagSet(flag.PanicOnError).Parse([]string{
		"--secret=bar",
		"--token=foo",
		"--private-key-path=" + os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
		"--service-principal-id=9995069170671",                 // (nonexistent)
		"--service-principal-key-id=EoSzAQwRly2XTKib5EfCSj7jl", // (nonexistent)
	})
	s.NoError(e)
	e = s.subject.Populate()
	s.NoError(e)

	s.Run("explicit AuthPreference=basic", func() {
		hdr, req, svr := s.authPreferenceMockServer()
		subject, e := s.subject.DupWith(WithTestServer(svr), WithFavouringBasicAuthentication())
		s.NoError(e)
		_, e = subject.Do(req)
		s.NoError(e)
		s.Regexp("^Basic ", hdr.Get("Authorization"))
	})

	s.Run("explicit AuthPreference=bearer", func() {
		hdr, req, svr := s.authPreferenceMockServer()
		subject, e := s.subject.DupWith(WithTestServer(svr), WithFavouringBearerAuthentication())
		s.NoError(e)
		_, e = subject.Do(req)
		s.NoError(e)
		s.Regexp("^Bearer ", hdr.Get("Authorization"))
	})

	s.Run("implicit", func() {
		hdr, req, svr := s.authPreferenceMockServer()
		subject, e := s.subject.DupWith(WithTestServer(svr))
		s.NoError(e)
		_, e = subject.Do(req)
		s.NoError(e)
		s.Regexp("^Bearer ", hdr.Get("Authorization"))
	})
}
