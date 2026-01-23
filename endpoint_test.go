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
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	. "github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/suite"
)

type endpointTest struct {
	suite.Suite
	testServer *httptest.Server
	envBackups []string
}

func (this *endpointTest) SetupTest() {
	this.testServer = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
}

func (this *endpointTest) TearDownTest() {
	if this.testServer != nil {
		this.testServer.Close()
	}
	for _, env := range this.envBackups {
		_ = os.Unsetenv(env)
	}
	this.envBackups = nil
}

func (this *endpointTest) setEnv(key, value string) {
	this.envBackups = append(this.envBackups, key)
	_ = os.Setenv(key, value)
}

func (this *endpointTest) TestEndpointConfig_NoEndpoints() {
	this.setEnv("SAKURA_ACCESS_TOKEN", "test-token")
	this.setEnv("SAKURA_ACCESS_TOKEN_SECRET", "test-secret")

	var client Client
	err := client.SetEnviron(os.Environ())
	this.NoError(err)

	api, err := client.DupWith(WithTestServer(this.testServer))
	this.NoError(err)

	err = api.Populate()
	this.NoError(err)

	cfg, err := api.EndpointConfig()
	this.NoError(err)
	this.NotNil(cfg)
	this.Nil(cfg.Endpoints)
}

func (this *endpointTest) TestEndpointConfig_FromEnvVariables() {
	// Set environment variables
	this.setEnv("SAKURA_ACCESS_TOKEN", "test-token")
	this.setEnv("SAKURA_ACCESS_TOKEN_SECRET", "test-secret")
	this.setEnv("SAKURA_ENDPOINTS_IAAS", "https://secure.sakura.ad.jp/cloud/zone")
	this.setEnv("SAKURA_ENDPOINTS_IAM", "https://secure.sakura.ad.jp/cloud/api/iam/1.0")

	var client Client
	err := client.SetEnviron(os.Environ())
	this.NoError(err)

	api, err := client.DupWith(WithTestServer(this.testServer))
	this.NoError(err)

	err = api.Populate()
	this.NoError(err)

	cfg, err := api.EndpointConfig()
	this.NoError(err)
	this.NotNil(cfg)
	this.NotNil(cfg.Endpoints)
	this.Equal("https://secure.sakura.ad.jp/cloud/zone", cfg.Endpoints["iaas"])
	this.Equal("https://secure.sakura.ad.jp/cloud/api/iam/1.0", cfg.Endpoints["iam"])
}

func (this *endpointTest) TestEndpointConfig_EnvOverridesProfile() {
	defer func() {
		_ = os.RemoveAll(os.TempDir() + "/test-profiles-override")
		_ = os.Unsetenv("XDG_CONFIG_HOME")
	}()

	profileDir := os.TempDir() + "/test-profiles-override"
	_ = os.MkdirAll(profileDir+"/default", 0700)

	profileContent := `{
  "AccessToken": "test-token",
  "AccessTokenSecret": "test-secret",
  "Endpoints": {
    "iaas": "https://profile.example.com/cloud/zone"
  }
}`

	this.setEnv("XDG_CONFIG_HOME", profileDir)
	_ = os.WriteFile(profileDir+"/default/config.json", []byte(profileContent), 0600)

	this.setEnv("SAKURA_ENDPOINTS_IAAS", "https://env.example.com/cloud/zone")
	this.setEnv("SAKURA_ACCESS_TOKEN", "test-token")
	this.setEnv("SAKURA_ACCESS_TOKEN_SECRET", "test-secret")

	var client Client
	err := client.SetEnviron(os.Environ())
	this.NoError(err)

	api, err := client.DupWith(WithTestServer(this.testServer))
	this.NoError(err)

	err = api.Populate()
	this.NoError(err)

	cfg, err := api.EndpointConfig()
	this.NoError(err)
	this.NotNil(cfg)
	this.NotNil(cfg.Endpoints)
	// Environment variable should override profile
	this.Equal("https://env.example.com/cloud/zone", cfg.Endpoints["iaas"])
}

func TestEndpointSuite(t *testing.T) {
	suite.Run(t, new(endpointTest))
}
