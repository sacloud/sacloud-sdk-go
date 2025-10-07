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
	"encoding/json"
	"io"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// Profile is a set of options, named.
type ProfileAPI interface {
	// List the currently known names of profiles
	List() ([]string, error)

	// Load one of profiles, by name
	Read(name string) (*Profile, error)

	// Create a new profile
	Create(profile *Profile) error

	// Update (merge) an existing profile, updated contents returned
	Update(profile *Profile) (*Profile, error)

	// Delete an existing profile
	Delete(name string) error

	// Get the default profile
	GetCurrentName() (string, error)

	// Set the default profile
	SetCurrentName(name string) error
}

// A (loaded) profile
type Profile struct {
	dir string

	// Name (intuitive)
	Name string

	// Profile contents
	//
	// This is intentionally untyped to allow arbitrary key/values
	Attributes map[string]any
}

type ProfileOp struct {
	// profile directory (historically `~/.usacloud` but configurable via env var)
	// Note that this is not necessarily existing at process startup.
	dir string
}

// Creates a profile operator
func NewProfileOp(envp []string) *ProfileOp { return &ProfileOp{lookupProfileDir(envp)} }

func (this *ProfileOp) List() ([]string, error) {
	var ret []string
	glob := filepath.Join(this.dir, "*", "config.json")

	if stat, err := os.Stat(this.dir); err != nil {
		return []string{}, nil // This is when e.g. the first invocation
	} else if !stat.IsDir() {
		return nil, NewErrorf("failed to open %+v", this.dir)
	} else if ent, err := filepath.Glob(glob); err != nil {
		return nil, Wrapf(err, "failed to open %+v", this.dir)
	} else {
		for _, Profile := range ent {
			if stat, err := os.Stat(Profile); err != nil {
				// skip
			} else if stat.IsDir() {
				// skip
			} else {
				dir := filepath.Dir(Profile)
				name := filepath.Base(dir)
				ret = append(ret, name)
			}
		}
	}

	slices.Sort(ret) // stabilize return order (easy test)
	return ret, nil
}

func (this *ProfileOp) Read(name string) (*Profile, error) {
	n := filepath.Join(name, "config.json")

	return this.open(n, os.O_RDONLY, func(fp *os.File) (*Profile, error) {
		var attrs map[string]any

		if buf, err := io.ReadAll(fp); err != nil {
			return nil, Wrapf(err, "failed to read %+v", fp.Name())
		} else if err := json.Unmarshal(buf, &attrs); err != nil {
			return nil, Wrapf(err, "failed to parse %+v", fp.Name())
		} else {
			return &Profile{this.dir, name, attrs}, nil
		}
	},
	)
}

func (this *ProfileOp) Create(p *Profile) error {
	n := filepath.Join(p.Name, "config.json")
	_, err := this.open(n, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, func(fp *os.File) (*Profile, error) {
		if buf, err := json.MarshalIndent(p.Attributes, "", "  "); err != nil {
			return nil, Wrapf(err, "failed to serialize %+v", p.Pathname())
		} else if _, err := fp.Write(buf); err != nil {
			return nil, Wrapf(err, "failed to write %+v", p.Pathname())
		} else {
			return p, nil
		}
	},
	)
	return err
}

func (this *ProfileOp) Update(p *Profile) (*Profile, error) {
	n := filepath.Join(p.Name, "config.json")

	return this.open(n, os.O_RDWR, func(fp *os.File) (*Profile, error) {
		// :TODO: need flock for possible concurrent writes
		// :TODO: better have backup mechanism
		var attrs map[string]any

		if buf, err := io.ReadAll(fp); err != nil {
			return nil, Wrapf(err, "failed to read %+v", fp.Name())
		} else if err := json.Unmarshal(buf, &attrs); err != nil {
			return nil, Wrapf(err, "failed to parse %+v", fp.Name())
		}

		ret := &Profile{this.dir, p.Name, deepMerge(attrs, p.Attributes)}

		if buf, err := json.MarshalIndent(ret.Attributes, "", "  "); err != nil {
			return nil, Wrapf(err, "failed to serialize %+v", p.Pathname())
		} else if _, err := fp.Seek(0, 0); err != nil {
			return nil, Wrapf(err, "failed to seek %+v", p.Pathname())
		} else if err := fp.Truncate(0); err != nil {
			return nil, Wrapf(err, "failed to truncate %+v", p.Pathname())
		} else if _, err := fp.Write(buf); err != nil {
			return nil, Wrapf(err, "failed to write %+v", p.Pathname())
		} else {
			return ret, nil
		}
	})
}

func (this *ProfileOp) Delete(name string) error {
	if _, err := os.Stat(this.dir); err != nil {
		return nil // already gone, nothing to do
	}

	root, err := os.OpenRoot(this.dir)
	if err != nil {
		return Wrapf(err, "failed to open directory %+v", this.dir)
	}
	defer func() { _ = root.Close() }()

	return root.RemoveAll(name)
}

func (this *ProfileOp) GetCurrentName() (string, error) {
	var ret string
	_, err := this.open("current", os.O_RDONLY, func(fp *os.File) (*Profile, error) {
		if buf, err := io.ReadAll(fp); err != nil {
			return nil, Wrapf(err, "failed to read %+v", fp.Name())
		} else {
			ret = strings.TrimSpace(string(buf))
			return nil, nil
		}
	})

	return ret, err
}

func (this *ProfileOp) SetCurrentName(name string) error {
	if list, err := this.List(); err != nil {
		return err
	} else if !slices.Contains(list, name) {
		return NewErrorf("invalid profile name: %+v", name)
	}

	_, err := this.open("current", os.O_WRONLY|os.O_TRUNC, func(fp *os.File) (*Profile, error) {
		if _, err := fp.WriteString(name); err != nil {
			return nil, Wrapf(err, "failed to write %+v", fp.Name())
		} else {
			return nil, nil
		}
	})
	return err
}

// Calculated pathname of the configuration file
func (this *Profile) Pathname() string { return filepath.Join(this.dir, this.Name, "config.json") }

func (this *ProfileOp) open(
	n string,
	mode int,
	callback func(*os.File) (*Profile, error),
) (*Profile, error) {
	if (mode & os.O_CREATE) != 0 {
		if err := os.MkdirAll(this.dir, 0o700); err != nil {
			return nil, Wrapf(err, "failed to create directory %+v", this.dir)
		}
	}

	root, err := os.OpenRoot(this.dir)
	if err != nil {
		return nil, Wrapf(err, "failed to open directory %+v", this.dir)
	}
	defer func() { _ = root.Close() }()

	if (mode & os.O_CREATE) != 0 {
		dirname := filepath.Dir(n)
		if err := root.MkdirAll(dirname, 0o700); err != nil {
			return nil, Wrapf(err, "failed to create directory %+v", dirname)
		}
	}

	file, err := root.OpenFile(n, mode, 0o600)
	if err != nil {
		return nil, Wrapf(err, "failed to open %+v", n)
	}
	defer func() { _ = file.Close() }()

	return callback(file)
}

func deepMerge(dst, src map[string]any) map[string]any {
	ret := make(map[string]any)
	maps.Copy(ret, dst)
	for k, v := range src {
		switch v := v.(type) {
		case map[string]any:
			if ov, ok := ret[k]; !ok {
				ret[k] = v
			} else if ov, ok := ov.(map[string]any); !ok {
				ret[k] = v
			} else {
				ret[k] = deepMerge(ov, v)
			}
		case []any:
			if ov, ok := ret[k]; !ok {
				ret[k] = v
			} else if ov, ok := ov.([]any); !ok {
				ret[k] = v
			} else {
				ret[k] = append(ov, v...)
			}
		default:
			ret[k] = v
		}
	}
	return ret
}

func lookupProfileDir(envp []string) string {
	if v, ok := lookupEnv(envp, "SAKURACLOUD_PROFILE_DIR"); ok {
		return filepath.Clean(v)
	} else if v, ok := lookupEnv(envp, "USACLOUD_PROFILE_DIR"); ok {
		return filepath.Clean(v) // backward compat
	} else if v, ok := lookupEnv(envp, "XDG_CONFIG_HOME"); ok {
		// if, and only if `~/.config/usacloud` exists, take it.
		ret := filepath.Join(v, "usacloud")
		if stat, err := os.Stat(ret); err == nil && stat.IsDir() {
			return filepath.Clean(ret)
		}
	}

	// fallback to '~/.usacloud'
	if home, err := os.UserHomeDir(); err == nil && home != "" {
		return filepath.Join(home, ".usacloud")
	} else {
		// :ESOTERIC: $HOME not set
		//
		// @shyouhei doesn't think it's worth returning an error here.
		// Is there any mitigation we can do then?
		panic("unable to determine profile directory")
	}
}

// lookupEnv searches for an environment variable in the provided slice
// and returns its value and a boolean indicating if it was found.
// The key must not be empty.
func lookupEnv(envp []string, key string) (string, bool) {
	n := len(key)
	for _, env := range envp {
		//nolint:gocritic
		if len(env) <= n {
			continue
		} else if env[:n] != key {
			continue
		} else if env[n] != '=' {
			continue
		} else {
			return env[n+1:], true
		}
	}
	return "", false
}
