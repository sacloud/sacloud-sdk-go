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
	"encoding/json"
	"fmt"
	"os"
	"testing"

	. "github.com/sacloud/http-client-go"
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
	} else if home, ok := os.LookupEnv("HOME"); !ok {
		s.T().Fatal("$HOME is not set")
	} else {
		s.dir = dir
		s.home = home
		if err := os.Setenv("HOME", s.dir); err != nil {
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
		os.WriteFile(dir+"/.usacloud/current", []byte("usacloud"), 0o600)
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
	if err := os.Setenv("HOME", s.home); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ProfileTestSuite) TearDownSubTest() {
	_ = s.op.Create(&Profile{
		Name: "usacloud",
		Attributes: map[string]any{
			"Zone":              "usacloud",
			"PrivateKeyPEMPath": s.dir + "/usamin.pem",
		},
	})
}

func (s *ProfileTestSuite) TestProfileOp_usacloud() {
	s.op = NewProfileOp(os.Environ())
	s.NotNil(s.op)
	op := s.op

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

	s.Run("GetCurrentName", func() {
		name, err := op.GetCurrentName()
		s.NoError(err)
		s.Equal("usacloud", name) // not set yet
	})

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
}

func (s *ProfileTestSuite) TestProfileOp_XDG() {
	s.op = NewProfileOp([]string{"XDG_CONFIG_HOME=" + s.dir + "/.config"})
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

func (s *ProfileTestSuite) TestProfile_GetCacheFilePath() {
	op := NewProfileOp(os.Environ())
	s.NotNil(s.op)

	subject, err := op.Read("usacloud")
	s.NoError(err)
	s.NotNil(subject)

	fmt.Printf("%#+v", subject)

	path, err := subject.GetCacheFilePath()
	s.NoError(err)
	s.NotEmpty(path)
	s.Equal(s.dir+"/.usacloud/usacloud/cache/5f20028ef6763408a4dd438db2b0e3a6e7455b82195335f04204b0662345a132.json", path)
}
