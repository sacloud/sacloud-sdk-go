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
		os.MkdirAll(dir+"/.usacloud/usacloud", 0o700)
		os.MkdirAll(dir+"/.usacloud/broken", 0o700)
		os.MkdirAll(dir+"/.config/usacloud/xdg", 0o700)

		os.WriteFile(dir+"/.usacloud/usacloud/config.json", []byte(`{"Zone":"usacloud"}`), 0o600)
		os.WriteFile(dir+"/.usacloud/broken/config.json", []byte("偶因狂疾成殊類 災患相仍不可逃"), 0o600)
		os.WriteFile(dir+"/.usacloud/current", []byte("usacloud"), 0o600)
		os.WriteFile(dir+"/.config/usacloud/xdg/config.json", []byte(`{"Zone":"xdg"}`), 0o600)
	}
}

func (s *ProfileTestSuite) TearDownSuite() {
	if err := os.Setenv("HOME", s.home); err != nil {
		s.T().Fatal(err)
	}
}

func (s *ProfileTestSuite) TearDownSubTest() {
	_ = s.op.Create(&Profile{
		Name:       "usacloud",
		Attributes: map[string]any{"Zone": "usacloud"},
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
			s.Equal(map[string]any{"Zone": "usacloud"}, profile.Attributes)
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
