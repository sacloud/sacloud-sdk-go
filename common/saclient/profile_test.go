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

package saclient_test

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	. "github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/suite"
)

type ProfileTestSuite struct {
	suite.Suite
	home string
	dir  string
	op   *ProfileOp
}

func TestProfileTestSuite(t *testing.T) { suite.Run(t, new(ProfileTestSuite)) }

//nolint:errcheck,gosec
func (s *ProfileTestSuite) SetupSuite() {
	// Note that `s.T().TempDir()` is removed every time after a _test_, not afrer a suite.
	if dir, err := os.MkdirTemp(os.TempDir(), "profile_test"); err != nil {
		s.T().Fatal(err)
	} else {
		s.dir = dir
		home, _ := os.UserHomeDir()
		s.home = home

		var envkey string
		switch runtime.GOOS {
		case "windows":
			envkey = "USERPROFILE"
		default:
			envkey = "HOME"
		}
		if err := os.Setenv(envkey, s.dir); err != nil {
			s.T().Fatal(err)
		}

		// create sample profiles
		sane := map[string]any{
			"Zone":              "usacloud",
			"PrivateKeyPEMPath": dir + "/usamin.pem",
		}
		buf, _ := json.MarshalIndent(sane, "", "  ")
		os.MkdirAll(dir+"/.usacloud/usacloud", 0o700)
		os.MkdirAll(dir+"/.usacloud/broken", 0o700)
		os.MkdirAll(dir+"/.config/usacloud/xdg", 0o700)

		os.WriteFile(dir+"/.usacloud/usacloud/config.json", buf, 0o600)
		os.WriteFile(dir+"/.usacloud/broken/config.json", []byte("偶因狂疾成殊類 災患相仍不可逃"), 0o600)
		// os.WriteFile(dir+"/.usacloud/current", []byte("usacloud"), 0o600) // created during tests
		os.WriteFile(dir+"/.config/usacloud/xdg/config.json", []byte(`{"Zone":"xdg"}`), 0o600)

		os.WriteFile(dir+"/usamin.pem", []byte(`
-----BEGIN RSA PRIVATE KEY-----
MIICXQIBAAKBgQC/AfcvlUlhcPpDD/1HqWBWvGQZxdER+fy6jbm1BlhVT156hjZi
UUwetUMrVGuy+bYE50j+qJB2VYKIhUTIUYCJg/AlruszlmydV0dOWPpSsLMXA5XU
GhoijaZY9l8vsbGN3n0QJ313GvFQQ+GrP1PzmRbpK686weAwtCx+PXYPQwIDAQAB
AoGAEvG69nk0AfoWmDgpwsXFzFR7CSNZjRLiQg50cMPkVvG8SSKumim+Bv2rX8zL
scCakPnvf3JwgYwRmkC9hbCvssfQK2o0Zzc6zPa560TxXYK5rADTfMXqeLnF6nFZ
sKLlE5vxyv2XD6zDcc1K2q25ARYMeWOGQ2WfuMYexBd36EECQQD0va3JquOaPQI7
2yRXNumv2fRwYohnJxOymu4vKZp11R0gTGljsv7y8I+mcVDJnJy27t9a7tUSLS4F
G1FMId0LAkEAx8t39aRzchpUoJYl9KmigFQ5AS6qAmDqdGIOBFQ5hf6HErukbRBd
2q+tNXAKF62ecXR3dlaS54CpSXkQVxlJqQJBANJD1/hIEk0kFzQ3nSw06GaFmcWo
UcpVv02WYAYy9xo/I0vpei4GzZUI6lG0TxU3sUhVR53HTVXVbRFEG/+NpGsCQQCi
qPilOJn0z5MOmq+UHXd7WxZ96+vlu9mlnx8iTx/2A18c1T/su2Jt5JDz7J+K34Mb
g2KvKZS4fXtVoga3opLhAkAtR4iVtxGi3NxOw0XrTXClzJD1e357/MrSDQ09gdRG
sP9Knwr9WVBtRYPRFjC3YccLTwoQnjVcF1qJN6ybMvnS
-----END RSA PRIVATE KEY-----
		`), 0o600)
	}
}

func (s *ProfileTestSuite) TearDownSuite() {
	if s.home == "" {
		return
	}

	var envkey string
	switch runtime.GOOS {
	case "windows":
		envkey = "USERPROFILE"
	default:
		envkey = "HOME"
	}
	if err := os.Setenv(envkey, s.home); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ProfileTestSuite) TearDownSubTest() {
	if s.op == nil {
		return
	}

	_ = s.op.Create(&Profile{
		Name: "usacloud",
		Attributes: map[string]any{
			"Zone":              "usacloud",
			"PrivateKeyPEMPath": s.dir + "/usamin.pem",
		},
	})
}

func (s *ProfileTestSuite) TestProfileOp_usacloud() {
	var err error
	s.op, err = NewProfileOp(os.Environ())
	s.NoError(err)
	s.NotNil(s.op)
	op := s.op

	s.Run("SetCurrentName", func() {
		s.Run("on success", func() {
			err := op.SetCurrentName("usacloud")
			s.NoError(err)
		})

		s.Run("not found", func() {
			err := op.SetCurrentName("not-found")
			s.Error(err)

			var e1 *Error
			s.ErrorAs(err, &e1)
		})
	})

	s.Run("GetCurrentName", func() {
		name, err := op.GetCurrentName()
		s.NoError(err)
		s.Equal("usacloud", name) // not set yet
	})

	s.Run("List", func() {
		names, err := op.List()
		s.NoError(err)
		s.Equal([]string{"broken", "usacloud"}, names)
	})

	s.Run("Read", func() {
		s.Run("found sane", func() {
			profile, err := op.Read("usacloud")
			s.NoError(err)
			s.NotNil(profile)
			s.Equal("usacloud", profile.Name)
			s.Equal(map[string]any{
				"Zone":              "usacloud",
				"PrivateKeyPEMPath": s.dir + "/usamin.pem",
			}, profile.Attributes)
		})

		s.Run("found broken", func() {
			profile, err := op.Read("broken")
			s.Nil(profile)
			var e1 *Error
			var e2 *json.SyntaxError
			s.ErrorAs(err, &e1)
			s.ErrorAs(err, &e2)
		})

		s.Run("not found", func() {
			profile, err := op.Read("not-found")
			s.Nil(profile)
			var e1 *Error
			var e2 *os.PathError
			s.ErrorAs(err, &e1)
			s.ErrorAs(err, &e2)
		})
	})

	s.Run("Create", func() {
		s.Run("on success", func() {
			err := op.Create(&Profile{
				Name:       "new-profile",
				Attributes: map[string]any{"Zone": "new-profile"},
			})
			s.NoError(err)
		})

		s.Run("on conflict", func() {
			err := op.Create(&Profile{
				Name:       "usacloud",
				Attributes: map[string]any{"Zone": "new-profile"},
			})
			s.Error(err)

			var e1 *Error
			var e2 *os.PathError
			s.ErrorAs(err, &e1)
			s.ErrorAs(err, &e2)
		})

		s.Run("on a malicious name", func() {
			err := op.Create(&Profile{
				Name:       "../../../../../../etc/passwd",
				Attributes: map[string]any{"Zone": "new-profile"},
			})
			s.Error(err)

			var e1 *Error
			var e2 *os.PathError
			s.ErrorAs(err, &e1)
			s.ErrorAs(err, &e2)
		})
	})

	s.Run("Update", func() {
		s.Run("on success", func() {
			profile, err := op.Update(&Profile{
				Name: "usacloud",
				Attributes: map[string]any{
					"Zone":      "updated",
					"Arbitrary": []string{"values", "can", "be", "set"},
				},
			})
			s.NoError(err)
			s.NotNil(profile)
			s.Equal("usacloud", profile.Name)
			s.Equal("updated", profile.Attributes["Zone"])
		})

		s.Run("not found", func() {
			profile, err := op.Update(&Profile{
				Name: "not-found",
				Attributes: map[string]any{
					"Zone": "updated",
				},
			})
			s.Nil(profile)
			s.Error(err)

			var e1 *Error
			var e2 *os.PathError
			s.ErrorAs(err, &e1)
			s.ErrorAs(err, &e2)
		})

		s.Run("on repeated update", func() {
			// Run #1
			profile, err := op.Update(&Profile{
				Name: "usacloud",
				Attributes: map[string]any{
					"Zone":      "updated",
					"Arbitrary": []string{"values", "can", "be", "set", "again"},
				},
			})
			s.NoError(err)
			s.NotNil(profile)

			// Run #2
			profile, err = op.Update(&Profile{
				Name: "usacloud",
				Attributes: map[string]any{
					"Zone":      "updated",
					"Arbitrary": []string{"shortened", "values"},
				},
			})
			s.NoError(err)
			s.NotNil(profile)

			// verify
			profile, err = op.Read("usacloud")
			s.NoError(err)
			s.NotNil(profile)
			s.Equal("usacloud", profile.Name)
			attrs, ok := profile.Attributes["Arbitrary"].([]any)
			s.True(ok)
			s.Len(attrs, 2)
			s.Equal("shortened", attrs[0].(string))
			s.Equal("values", attrs[1].(string))
		})
	})

	s.Run("Delete", func() {
		s.Run("on success", func() {
			err := op.Delete("usacloud")
			s.NoError(err)
		})

		s.Run("already gone", func() {
			err := op.Delete("not-found")
			s.NoError(err)
		})
	})
}

func (s *ProfileTestSuite) TestProfileOp_XDG() {
	var err error
	s.op, err = NewProfileOp([]string{"XDG_CONFIG_HOME=" + s.dir + "/.config"})
	s.NoError(err)
	s.NotNil(s.op)
	op := s.op

	s.Run("List", func() {
		names, err := op.List()
		s.NoError(err)
		s.Equal([]string{"xdg"}, names)
	})

	s.Run("Read", func() {
		s.Run("found sane", func() {
			profile, err := op.Read("xdg")
			s.NoError(err)
			s.NotNil(profile)
			s.Equal("xdg", profile.Name)
			s.Equal(map[string]any{"Zone": "xdg"}, profile.Attributes)
		})
	})

	s.Run("Create", func() {
		s.Run("on success", func() {
			err := op.Create(&Profile{
				Name:       "new-profile",
				Attributes: map[string]any{"Zone": "new-profile"},
			})
			s.NoError(err)

			// must not exist
			_, err = os.Stat(s.dir + "/.usacloud/new-profile/config.json")
			var e1 *os.PathError
			s.ErrorAs(err, &e1)
		})
	})
}

func (s *ProfileTestSuite) TestProfileOp_V1YAML() {
	dir := s.T().TempDir()

	s.Require().NoError(os.MkdirAll(dir+"/default", 0o700))
	s.Require().NoError(os.MkdirAll(dir+"/legacy", 0o700))

	v1 := `version: 1
credentials:
  access_token: v1-token
  access_token_secret: v1-secret
  service_principal_id: spid
  service_principal_key_kid: kid
go:
  zone: is1a
  zones:
    - is1a
    - tk1a
endpoints:
  iam: https://example.invalid/iam
dotnet: {foobar: preserved}
`

	v0 := `{"AccessToken":"legacy-token","Zone":"tk1a"}`

	s.Require().NoError(os.WriteFile(dir+"/default/config.yaml", []byte(v1), 0o600))
	s.Require().NoError(os.WriteFile(dir+"/legacy/config.json", []byte(v0), 0o600))

	op, err := NewProfileOp([]string{"SAKURA_PROFILE_DIR=" + dir})
	s.Require().NoError(err)
	s.Require().NotNil(op)

	s.Run("List", func() {
		names, err := op.List()
		s.NoError(err)
		s.Equal([]string{"default", "legacy"}, names)
	})

	s.Run("read v1 yaml", func() {
		profile, err := op.Read("default")
		s.NoError(err)
		s.NotNil(profile)
		s.EqualValues(1, profile.Version)
		s.Equal("v1-token", profile.Attributes["AccessToken"])
		s.Equal("v1-secret", profile.Attributes["AccessTokenSecret"])
		s.Equal("spid", profile.Attributes["ServicePrincipalID"])
		s.Equal("kid", profile.Attributes["ServicePrincipalKeyKID"])
		s.Equal("is1a", profile.Attributes["Zone"])

		eps, ok := profile.Attributes["Endpoints"].(map[string]any)
		s.True(ok)
		s.Equal("https://example.invalid/iam", eps["iam"])
		s.Equal("v1-token", profile.Credentials.AccessToken.MustGet())
		s.Equal("kid", profile.Credentials.ServicePrincipalKeyKID.MustGet())
		s.Equal("is1a", profile.Go.Zone.MustGet())
		s.Equal([]string{"is1a", "tk1a"}, profile.Go.Zones.MustGet())
		s.Equal("https://example.invalid/iam", profile.Endpoints["iam"])
	})

	s.Run("fallback to v0 json", func() {
		profile, err := op.Read("legacy")
		s.NoError(err)
		s.NotNil(profile)
		s.EqualValues(0, profile.Version)
		s.Equal("legacy-token", profile.Attributes["AccessToken"])
		s.Equal("tk1a", profile.Attributes["Zone"])
		s.Equal("legacy-token", profile.Credentials.AccessToken.MustGet())
		s.Equal("tk1a", profile.Go.Zone.MustGet())
	})

	s.Run("create writes v1 yaml", func() {
		created := &Profile{Name: "created", Version: 1}
		created.Credentials.AccessToken.SetTo("created-token")
		created.Credentials.AccessTokenSecret.SetTo("created-secret")
		created.Go.Zone.SetTo("is1b")
		err := op.Create(created)
		s.NoError(err)

		_, err = os.Stat(dir + "/created/config.yaml")
		s.NoError(err)

		_, err = os.Stat(dir + "/created/config.json")
		var e1 *os.PathError
		s.ErrorAs(err, &e1)

		profile, err := op.Read("created")
		s.NoError(err)
		s.NotNil(profile)
		s.Equal("created-token", profile.Attributes["AccessToken"])
		s.Equal("is1b", profile.Attributes["Zone"])
	})

	s.Run("create writes v0 json", func() {
		err := op.Create(&Profile{
			Name: "created-v0",
			Attributes: map[string]any{
				"AccessToken": "created-v0-token",
				"Zone":        "tk1a",
			},
		})
		s.NoError(err)

		contents, err := os.ReadFile(filepath.Clean(dir + "/created-v0/config.json"))
		s.NoError(err)
		s.NotContains(string(contents), `"Version"`)

		profile, err := op.Read("created-v0")
		s.NoError(err)
		s.EqualValues(0, profile.Version)
		s.Equal("created-v0-token", profile.Credentials.AccessToken.MustGet())
	})

	s.Run("update preserves source format", func() {
		profile, err := op.Update(&Profile{
			Name:       "default",
			Attributes: map[string]any{"Zone": "is1b"},
		})
		s.NoError(err)
		s.EqualValues(1, profile.Version)
		s.Equal("is1b", profile.Go.Zone.MustGet())

		contents, err := os.ReadFile(filepath.Clean(dir + "/default/config.yaml"))
		s.NoError(err)
		s.Contains(string(contents), "dotnet:")
		s.Contains(string(contents), "foobar: preserved")

		profile, err = op.Update(&Profile{
			Name:       "legacy",
			Attributes: map[string]any{"Zone": "tk1b"},
		})
		s.NoError(err)
		s.EqualValues(0, profile.Version)
		s.Equal("tk1b", profile.Go.Zone.MustGet())

		contents, err = os.ReadFile(filepath.Clean(dir + "/legacy/config.json"))
		s.NoError(err)
		s.True(json.Valid(contents))
	})

	s.Run("update v1 writes an updated profile", func() {
		updated, err := op.Read("default")
		s.NoError(err)
		updated.Credentials.AccessToken.SetTo("replacement-token")
		updated.Go.Zone.SetTo("tk1b")

		profile, err := op.Update(updated)
		s.NoError(err)
		s.Same(updated, profile)

		contents, err := os.ReadFile(filepath.Clean(dir + "/default/config.yaml"))
		s.NoError(err)
		s.Contains(string(contents), "dotnet:")
		s.Contains(string(contents), "access_token_secret: v1-secret")

		profile, err = op.Read("default")
		s.NoError(err)
		s.EqualValues(1, profile.Version)
		s.Equal("replacement-token", profile.Credentials.AccessToken.MustGet())
		s.Equal("tk1b", profile.Go.Zone.MustGet())
		s.Equal("v1-secret", profile.Credentials.AccessTokenSecret.MustGet())
	})
}

func (s *ProfileTestSuite) TestProfileOp_FULLV1YAML() {
	dir := s.T().TempDir()
	s.Require().NoError(os.MkdirAll(dir+"/v0", 0o700))
	s.Require().NoError(os.MkdirAll(dir+"/v1", 0o700))

	v0 := `{
  "APIRootURL": "https://secure.sakura.ad.jp/cloud/zone",
  "AcceptLanguage": "en-US,en;q=0.9",
  "AccessToken": "<your-access-token>",
  "AccessTokenSecret": "<your-access-secret>",
  "ArgumentMatchMode": "exact",
  "DefaultOutputType": "table",
  "DefaultQueryDriver": "jq",
  "DefaultZone": "is1a",
  "Endpoints": {
    "iam": "https://secure.sakura.ad.jp/cloud-test/api/iam/1.0"
  },
  "FakeMode": false,
  "FakeStorePath": "~/.usacloud/fake_store.json",
  "HTTPRequestRateLimit": 5,
  "HTTPRequestTimeout": 300,
  "NoColor": false,
  "ProcessTimeoutSec": 7200,
  "RetryMax": 0,
  "RetryWaitMax": 64,
  "RetryWaitMin": 1,
  "StatePollingInterval": 0,
  "StatePollingTimeout": 0,
  "TraceMode": "HTTP",
  "Zone": "is1a",
  "Zones": ["is1a", "is1b", "tk1a", "tk1b", "tk1v"]
}`
	//nolint:gosec
	v1 := `version: 1
credentials:
  access_token: <your-access-token>
  access_token_secret: <your-access-token-secret>
  service_principal_id: <your-service-principal-id>
  service_principal_key_kid: <your-service-principal-kid>
  private_key: '-----BEGIN RSA PRIVATE KEY-----...'
  private_key_path: /path/to/key
endpoints:
  iam: https://secure.sakura.ad.jp/cloud-test/api/iam/1.0
cli:
  argument_match_mode: exact
  default_output_type: table
  default_query_driver: jq
  no_color: true
  process_timeout_sec: 7200
go:
  api_root_url: https://secure.sakura.ad.jp/cloud/zone
  accept_language: en-US,en;q=0.9
  default_zone: is1a
  fake_mode: false
  fake_store_path: ~/.usacloud/fake_store.json
  http_request_rate_limit: 5
  http_request_timeout: 300
  retry_max: 0
  retry_wait_max: 64
  retry_wait_min: 1
  state_polling_interval: 0
  state_polling_timeout: 0
  trace_mode: HTTP
  zone: is1a
  zones:
    - is1a
    - is1b
    - tk1a
    - tk1b
    - tk1v
dotnet:
  str_field: foobar
  num_field: 123
  bool_field: true
`

	s.Require().NoError(os.WriteFile(dir+"/v0/config.json", []byte(v0), 0o600))
	s.Require().NoError(os.WriteFile(dir+"/v1/config.yaml", []byte(v1), 0o600))

	op, err := NewProfileOp([]string{"SAKURA_PROFILE_DIR=" + dir})
	s.Require().NoError(err)

	s.Run("read v0 full profile", func() {
		profile, err := op.Read("v0")
		s.NoError(err)
		s.EqualValues(0, profile.Version)
		s.Equal("<your-access-token>", profile.Credentials.AccessToken.MustGet())
		s.Equal("<your-access-secret>", profile.Credentials.AccessTokenSecret.MustGet())
		s.Equal("exact", profile.Cli.ArgumentMatchMode.MustGet())
		s.Equal("table", profile.Cli.DefaultOutputType.MustGet())
		s.Equal("jq", profile.Cli.DefaultQueryDriver.MustGet())
		s.False(profile.Cli.NoColor.MustGet())
		s.Equal(int64(7200), profile.Cli.ProcessTimeoutSec.MustGet())
		s.Equal("https://secure.sakura.ad.jp/cloud/zone", profile.Go.APIRootURL.MustGet())
		s.Equal("en-US,en;q=0.9", profile.Go.AcceptLanguage.MustGet())
		s.Equal("is1a", profile.Go.DefaultZone.MustGet())
		s.False(profile.Go.FakeMode.MustGet())
		s.Equal("~/.usacloud/fake_store.json", profile.Go.FakeStorePath.MustGet())
		s.Equal(int64(5), profile.Go.HTTPRequestRateLimit.MustGet())
		s.Equal(int64(300), profile.Go.HTTPRequestTimeout.MustGet())
		s.Equal(int64(0), profile.Go.RetryMax.MustGet())
		s.Equal(int64(64), profile.Go.RetryWaitMax.MustGet())
		s.Equal(int64(1), profile.Go.RetryWaitMin.MustGet())
		s.Equal(int64(0), profile.Go.StatePollingInterval.MustGet())
		s.Equal(int64(0), profile.Go.StatePollingTimeout.MustGet())
		s.Equal("HTTP", profile.Go.TraceMode.MustGet())
		s.Equal("is1a", profile.Go.Zone.MustGet())
		s.Equal([]string{"is1a", "is1b", "tk1a", "tk1b", "tk1v"}, profile.Go.Zones.MustGet())
		s.Equal("https://secure.sakura.ad.jp/cloud-test/api/iam/1.0", profile.Endpoints["iam"])
		s.Equal("<your-access-token>", profile.Attributes["AccessToken"])
		s.Equal("<your-access-secret>", profile.Attributes["AccessTokenSecret"])
		s.Equal("exact", profile.Attributes["ArgumentMatchMode"])
		s.Equal("table", profile.Attributes["DefaultOutputType"])
		s.Equal("jq", profile.Attributes["DefaultQueryDriver"])
		s.Equal(false, profile.Attributes["NoColor"])
		s.Equal(float64(7200), profile.Attributes["ProcessTimeoutSec"])
		s.Equal("https://secure.sakura.ad.jp/cloud/zone", profile.Attributes["APIRootURL"])
		s.Equal("en-US,en;q=0.9", profile.Attributes["AcceptLanguage"])
		s.Equal("is1a", profile.Attributes["DefaultZone"])
		s.Equal(false, profile.Attributes["FakeMode"])
		s.Equal(float64(5), profile.Attributes["HTTPRequestRateLimit"])
		s.Equal(float64(300), profile.Attributes["HTTPRequestTimeout"])
		s.Equal(float64(0), profile.Attributes["RetryMax"])
		s.Equal(float64(64), profile.Attributes["RetryWaitMax"])
		s.Equal(float64(1), profile.Attributes["RetryWaitMin"])
		s.Equal("HTTP", profile.Attributes["TraceMode"])
		s.Equal("is1a", profile.Attributes["Zone"])
		s.Equal([]any{"is1a", "is1b", "tk1a", "tk1b", "tk1v"}, profile.Attributes["Zones"])
		endpoints, ok := profile.Attributes["Endpoints"].(map[string]any)
		s.True(ok)
		s.Equal("https://secure.sakura.ad.jp/cloud-test/api/iam/1.0", endpoints["iam"])
	})

	s.Run("read v1 full profile", func() {
		profile, err := op.Read("v1")
		s.NoError(err)
		s.EqualValues(1, profile.Version)
		s.Equal("<your-access-token>", profile.Credentials.AccessToken.MustGet())
		s.Equal("<your-access-token-secret>", profile.Credentials.AccessTokenSecret.MustGet())
		s.Equal("<your-service-principal-id>", profile.Credentials.ServicePrincipalID.MustGet())
		s.Equal("<your-service-principal-kid>", profile.Credentials.ServicePrincipalKeyKID.MustGet())
		s.Equal("-----BEGIN RSA PRIVATE KEY-----...", profile.Credentials.PrivateKey.MustGet())
		s.Equal("/path/to/key", profile.Credentials.PrivateKeyPEMPath.MustGet())
		s.Equal("https://secure.sakura.ad.jp/cloud-test/api/iam/1.0", profile.Endpoints["iam"])
		s.Equal("exact", profile.Cli.ArgumentMatchMode.MustGet())
		s.Equal("table", profile.Cli.DefaultOutputType.MustGet())
		s.Equal("jq", profile.Cli.DefaultQueryDriver.MustGet())
		s.True(profile.Cli.NoColor.MustGet())
		s.Equal(int64(7200), profile.Cli.ProcessTimeoutSec.MustGet())
		s.Equal("https://secure.sakura.ad.jp/cloud/zone", profile.Go.APIRootURL.MustGet())
		s.Equal("en-US,en;q=0.9", profile.Go.AcceptLanguage.MustGet())
		s.Equal("is1a", profile.Go.DefaultZone.MustGet())
		s.False(profile.Go.FakeMode.MustGet())
		s.Equal("~/.usacloud/fake_store.json", profile.Go.FakeStorePath.MustGet())
		s.Equal(int64(5), profile.Go.HTTPRequestRateLimit.MustGet())
		s.Equal(int64(300), profile.Go.HTTPRequestTimeout.MustGet())
		s.Equal(int64(0), profile.Go.RetryMax.MustGet())
		s.Equal(int64(64), profile.Go.RetryWaitMax.MustGet())
		s.Equal(int64(1), profile.Go.RetryWaitMin.MustGet())
		s.Equal(int64(0), profile.Go.StatePollingInterval.MustGet())
		s.Equal(int64(0), profile.Go.StatePollingTimeout.MustGet())
		s.Equal("HTTP", profile.Go.TraceMode.MustGet())
		s.Equal("is1a", profile.Go.Zone.MustGet())
		s.Equal([]string{"is1a", "is1b", "tk1a", "tk1b", "tk1v"}, profile.Go.Zones.MustGet())
		s.Equal("<your-access-token>", profile.Attributes["AccessToken"])
		s.Equal("<your-access-token-secret>", profile.Attributes["AccessTokenSecret"])
		s.Equal("<your-service-principal-id>", profile.Attributes["ServicePrincipalID"])
		s.Equal("<your-service-principal-kid>", profile.Attributes["ServicePrincipalKeyKID"])
		s.Equal("-----BEGIN RSA PRIVATE KEY-----...", profile.Attributes["PrivateKey"])
		s.Equal("/path/to/key", profile.Attributes["PrivateKeyPEMPath"])
		s.Equal("exact", profile.Attributes["ArgumentMatchMode"])
		s.Equal("table", profile.Attributes["DefaultOutputType"])
		s.Equal("jq", profile.Attributes["DefaultQueryDriver"])
		s.Equal(true, profile.Attributes["NoColor"])
		s.Equal(uint64(7200), profile.Attributes["ProcessTimeoutSec"])
		s.Equal("https://secure.sakura.ad.jp/cloud/zone", profile.Attributes["APIRootURL"])
		s.Equal("en-US,en;q=0.9", profile.Attributes["AcceptLanguage"])
		s.Equal("is1a", profile.Attributes["DefaultZone"])
		s.Equal(false, profile.Attributes["FakeMode"])
		s.Equal(uint64(5), profile.Attributes["HTTPRequestRateLimit"])
		s.Equal(uint64(300), profile.Attributes["HTTPRequestTimeout"])
		s.Equal(uint64(0), profile.Attributes["RetryMax"])
		s.Equal(uint64(64), profile.Attributes["RetryWaitMax"])
		s.Equal(uint64(1), profile.Attributes["RetryWaitMin"])
		s.Equal("HTTP", profile.Attributes["TraceMode"])

		dotnet, ok := profile.Attributes["dotnet"].(map[string]any)
		s.True(ok)
		s.Equal("foobar", dotnet["str_field"])
		s.Equal(uint64(123), dotnet["num_field"])
		s.Equal(true, dotnet["bool_field"])
		s.Equal("is1a", profile.Attributes["Zone"])
	})

	s.Run("update v1 full profile", func() {
		profile, err := op.Read("v1")
		s.NoError(err)
		profile.Credentials.AccessToken.SetTo("updated-access-token")
		profile.Go.Zone.SetTo("tk1a")

		_, err = op.Update(profile)
		s.NoError(err)

		profile, err = op.Read("v1")
		s.NoError(err)
		s.Equal("updated-access-token", profile.Credentials.AccessToken.MustGet())
		s.Equal("updated-access-token", profile.Attributes["AccessToken"])
		s.Equal("tk1a", profile.Go.Zone.MustGet())
		s.Equal("tk1a", profile.Attributes["Zone"])
		s.Equal("<your-access-token-secret>", profile.Attributes["AccessTokenSecret"])

		dotnet, ok := profile.Attributes["dotnet"].(map[string]any)
		s.True(ok)
		s.Equal("foobar", dotnet["str_field"])
	})
}

func (s *ProfileTestSuite) TestProfileOp_UnsetHOME() {
	runWithoutEnv(s, "TestProfileTestSuite/TestProfileOp_UnsetHOME", func() {
		err := os.Unsetenv("HOME")
		s.NoError(err)

		// SetupSuite sets USERPROFILE on Windows;
		// unset it explicitly to ensure NewProfileOp fails without home dir.
		if runtime.GOOS == "windows" {
			_ = os.Unsetenv("USERPROFILE")
			_ = os.Unsetenv("HOMEDRIVE")
			_ = os.Unsetenv("HOMEPATH")
		}

		s.Run("Barely no env", func() {
			op, err := NewProfileOp(os.Environ())
			s.Nil(op)
			s.Error(err)
		})

		s.Run("with SAKURA_PROFILE_DIR", func() {
			s.T().Setenv("SAKURA_PROFILE_DIR", s.dir+"/.usacloud")

			op, err := NewProfileOp(os.Environ())
			s.NoError(err)
			s.NotNil(op)

			names, err := op.List()
			s.NoError(err)
			s.Equal([]string{"broken", "usacloud"}, names)
		})
	})
}

func (s *ProfileTestSuite) TestProfile_GetCacheFilePath() {
	op, err := NewProfileOp(os.Environ())
	s.NoError(err)
	s.NotNil(op)

	subject, err := op.Read("usacloud")
	s.NoError(err)
	s.NotNil(subject)

	path, err := subject.GetCacheFilePath(nil, nil)
	s.NoError(err)
	s.NotEmpty(path)
	expected := filepath.Join(s.dir, ".usacloud", "usacloud", "cache", "5f20028ef6763408a4dd438db2b0e3a6e7455b82195335f04204b0662345a132.json")
	s.Equal(expected, path)
}

const randomKey = "jD6XUaOeniloPCB9X2ydjznzgOVgz1Bn"

//nolint:gosec // this is only a test
func runWithoutEnv(s *ProfileTestSuite, name string, yield func()) {
	if os.Getenv(randomKey) == "child" {
		yield()
	} else {
		cmd := exec.CommandContext(s.T().Context(), os.Args[0], "-test.run="+name)
		cmd.Env = []string{
			randomKey + "=child",
			"HOME=",
			"USERPROFILE=",
			"HOMEDRIVE=",
			"HOMEPATH=",
		}

		out, err := cmd.CombinedOutput()
		s.NoError(err)
		s.True(strings.HasPrefix(strings.TrimSpace(string(out)), "PASS"))
	}
}
