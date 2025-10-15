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

package client_test

import (
	"encoding/pem"
	"fmt"
	"os"
	"runtime"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	. "github.com/sacloud/http-client-go"
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
	"api-client-go/v%s (%s/%s; +https://github.com/sacloud/http-client-go)",
	Version,
	runtime.GOOS,
	runtime.GOARCH,
)

type ClientTestSuite struct {
	suite.Suite
	XDG_CONFIG_HOME *string
	subject         *Client
}

func TestClientTestSuite(t *testing.T) { suite.Run(t, new(ClientTestSuite)) }

//nolint:errcheck,gosec
func (s *ClientTestSuite) SetupSuite() {
	// Note that `s.T().TempDir()` is removed every time after a _test_, not afrer a suite.
	if dir, ok := os.LookupEnv("XDG_CONFIG_HOME"); ok {
		s.XDG_CONFIG_HOME = &dir
	}
	dir, _ := os.MkdirTemp(os.TempDir(), "profile_test")
	os.Setenv("XDG_CONFIG_HOME", dir)

	// create sample profiles
	os.MkdirAll(dir+"/usacloud/usacloud", 0o700)

	os.WriteFile(dir+"/usacloud/current", []byte("usacloud"), 0o600)
	os.WriteFile(dir+"/usacloud/usacloud/config.json",
		[]byte(fmt.Sprintf(`{
			"Zone":"usacloud",
			"PrivateKeyPEMPath":"%s/usacloud/usacloud/usamin.pem"
		}`, dir)),
		0o600,
	)
	fp, _ := os.OpenFile(dir+"/usacloud/usacloud/usamin.pem", os.O_WRONLY|os.O_CREATE, 0o600)
	defer fp.Close()
	pem.Encode(fp, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: []byte("dummy"),
	})
}

//nolint:errcheck,gosec
func (s *ClientTestSuite) TearDownSuite() {
	if s.XDG_CONFIG_HOME != nil {
		os.Setenv("XDG_CONFIG_HOME", *s.XDG_CONFIG_HOME)
	} else {
		os.Unsetenv("XDG_CONFIG_HOME")
	}
}

func (s *ClientTestSuite) SetupTest() {
	s.subject = new(Client)
}

func (s *ClientTestSuite) TestCLI() {
	e := s.subject.FlagSet().Parse([]string{
		"--secret=bar",
		"--token=foo",
		"--trace=error",
		"--zones=foo,\", bar\"",
	})
	s.NoError(e)
	e = s.subject.Populate()
	s.NoError(e)
	s.Equal(map[string]any{
		"AccessToken":         "foo",
		"AccessTokenSecret":   "bar",
		"APIRequestRateLimit": int64(5),
		"APIRequestTimeout":   int64(300),
		"PrivateKeyPEMPath":   os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
		"RetryMax":            int64(10),
		"RetryWaitMax":        int64(64),
		"RetryWaitMin":        int64(1),
		"TraceMode":           "error",
		"UserAgent":           ua,
		"Zone":                "usacloud",
		"Zones": []string{
			"foo",
			", bar",
		},
	}, s.subject.JSON())
}

func (s *ClientTestSuite) TestEnviron() {
	e := s.subject.SetEnviron([]string{
		"SAKURACLOUD_ACCESS_TOKEN_SECRET=bar",
		"SAKURACLOUD_ACCESS_TOKEN=foo",
		"SAKURACLOUD_API_REQUEST_RATE_LIMIT=20",
		"SAKURACLOUD_API_REQUEST_TIMEOUT=30",
		"SAKURACLOUD_API_ROOT_URL=https://api.example.com",
		"SAKURACLOUD_RETRY_MAX=3",
		"SAKURACLOUD_RETRY_WAIT_MAX=7",
		"SAKURACLOUD_RETRY_WAIT_MIN=5",
		"SAKURACLOUD_ZONE=foo",
		"SAKURACLOUD_ZONES=foo,\", bar\"",
		"SAKURACLOUD_TRACE=error",
		"XDG_CONFIG_HOME=" + os.Getenv("XDG_CONFIG_HOME"),
	})
	s.NoError(e)
	e = s.subject.Populate()
	s.NoError(e)
	s.Equal(map[string]any{
		"AccessToken":         "foo",
		"AccessTokenSecret":   "bar",
		"APIRequestRateLimit": int64(20),
		"APIRequestTimeout":   int64(30),
		"APIRootURL":          "https://api.example.com",
		"PrivateKeyPEMPath":   os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
		"RetryMax":            int64(3),
		"RetryWaitMax":        int64(7),
		"RetryWaitMin":        int64(5),
		"TraceMode":           "error",
		"UserAgent":           ua,
		"Zone":                "foo",
		"Zones": []string{
			"foo",
			", bar",
		},
	}, s.subject.JSON())
}

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
		"APIRequestTimeout":   int64(30),
		"APIRootURL":          "https://api.example.com",
		"DefaultZone":         "foo",
		"PrivateKeyPEMPath":   os.Getenv("XDG_CONFIG_HOME") + "/usacloud/usacloud/usamin.pem",
		"RetryMax":            int64(3),
		"RetryWaitMax":        int64(7),
		"RetryWaitMin":        int64(5),
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

	e := s.subject.FlagSet().Parse([]string{
		"--secret=bar",
		"--token=foo",
	})
	s.NoError(e)
	e = s.subject.Populate()
	s.NoError(e)
	s.Nil(s.subject.Profile())
	s.Equal("bar", s.subject.JSON()["AccessTokenSecret"])
}
