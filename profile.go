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
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"io"
	"iter"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

var _ ProfileAPI = (*ProfileOp)(nil)

// Creates a profile operator
func NewProfileOp(envp []string) *ProfileOp { return &ProfileOp{lookupProfileDir(envp)} }

func (this *ProfileOp) List() ([]string, error) {
	glob := filepath.Join(this.dir, "*", "config.json")

	if stat, err := os.Stat(this.dir); err != nil {
		return []string{}, nil // This is when e.g. the first invocation

	} else if !stat.IsDir() {
		return nil, NewErrorf("failed to open %+v", this.dir)

	} else if ent, err := filepath.Glob(glob); err != nil {
		return nil, Wrapf(err, "failed to open %+v", this.dir)

	} else {
		exists := func(p string) bool {
			_, err := os.Stat(p)
			return err == nil
		}
		isRegular := func(p string) bool {
			stat, _ := os.Stat(p)
			return stat.Mode().IsRegular()
		}

		q := slices.Values(ent)
		w := selectSeq(q, isRegular)
		e := selectSeq(w, exists)
		r := mapSeq(e, filepath.Dir)
		t := mapSeq(r, filepath.Base)
		y := slices.Sorted(t) // stabilize return order (easy test)

		return y, nil
	}
}

func (this *ProfileOp) Read(name string) (*Profile, error) {
	n := filepath.Join(name, "config.json")

	return this.open(n, os.O_RDONLY, func(fp *os.File) (*Profile, error) {
		var attrs map[string]any
		dec := json.NewDecoder(fp)

		if err := dec.Decode(&attrs); err != nil {
			return nil, Wrapf(err, "failed to parse %+v", fp.Name())

		} else {
			return &Profile{this.dir, name, attrs}, nil
		}
	})
}

func (this *ProfileOp) Create(p *Profile) error {
	n := filepath.Join(p.Name, "config.json")
	_, err := this.open(n, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, func(fp *os.File) (*Profile, error) {
		enc := json.NewEncoder(fp)
		enc.SetIndent("", "  ")

		if err := enc.Encode(p.Attributes); err != nil {
			return nil, Wrapf(err, "failed to serialize %+v", p.Pathname())

		} else {
			return p, nil
		}
	})
	return err
}

func (this *ProfileOp) Update(p *Profile) (*Profile, error) {
	n := filepath.Join(p.Name, "config.json")

	return this.open(n, os.O_RDWR, func(fp *os.File) (*Profile, error) {
		// :TODO: need flock for possible concurrent writes
		// :TODO: better have backup mechanism
		var attrs map[string]any
		dec := json.NewDecoder(fp)
		enc := json.NewEncoder(fp)
		enc.SetIndent("", "  ")

		if err := dec.Decode(&attrs); err != nil {
			return nil, Wrapf(err, "failed to parse %+v", fp.Name())
		}

		ret := &Profile{this.dir, p.Name, deepMerge(attrs, p.Attributes)}

		if _, err := fp.Seek(0, 0); err != nil {
			return nil, Wrapf(err, "failed to seek %+v", p.Pathname())

		} else if err := fp.Truncate(0); err != nil {
			return nil, Wrapf(err, "failed to truncate %+v", p.Pathname())

		} else if err := enc.Encode(ret.Attributes); err != nil {
			return nil, Wrapf(err, "failed to serialize %+v", p.Pathname())

		} else {
			return ret, nil
		}
	})
}

func (this *ProfileOp) Delete(name string) error {
	if _, err := os.Stat(this.dir); os.IsNotExist(err) {
		return nil // already gone, nothing to do

	} else if err != nil {
		return Wrapf(err, "failed to stat directory %+v", this.dir)
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

func (this *Profile) Get(k string) (any, bool) {
	if this == nil {
		return nil, false

	} else {
		v, ok := this.Attributes[k]
		return v, ok
	}
}

func (this *Profile) Set(k string, v any) {
	if this == nil {
		return

	}
	if this.Attributes == nil {
		this.Attributes = map[string]any{}
	}
	this.Attributes[k] = v
}

func (this *Profile) Keys() iter.Seq[string] {
	//nolint:gocritic
	if this == nil {
		return nonceSeq[string]()

	} else if this.Attributes == nil {
		return nonceSeq[string]()

	} else {
		return maps.Keys(this.Attributes)
	}
}

func (this *Profile) GetCacheFilePath(path *string, verbatim *string) (string, error) {
	var err error

	//nolint:gocritic
	if this == nil {
		return "", NewErrorf("nil profile")

	} else if path != nil && verbatim != nil {
		return "", NewErrorf("only one of path or verbatim can be set")

	} else if path == nil && verbatim == nil {
		// try obtaining from PrivateKeyPEMPath
		if str, ok := this.Get("PrivateKeyPEMPath"); !ok {
			return "", NewErrorf("neither path nor verbatim is given")

		} else if s, ok := str.(string); !ok {
			return "", NewErrorf("invalid PrivateKeyPEMPath: %T", str)

		} else {
			path = &s
		}
	}

	var bytes []byte

	//nolint:gosec // This `os.ReadFile` does not reveal any secret info
	if verbatim != nil {
		bytes = []byte(*verbatim)

	} else if bytes, err = os.ReadFile(*path); err != nil {
		return "", Wrapf(err, "failed to read PrivateKeyPEMPath")
	}

	if k, err := jwt.ParseRSAPrivateKeyFromPEM(bytes); err != nil {
		return "", Wrapf(err, "failed to parse PEM: %+v", path)

	} else if asn1, err := x509.MarshalPKIXPublicKey(&k.PublicKey); err != nil {
		return "", Wrapf(err, "failed to marshal public key: %+v", path)

	} else {
		sum := sha256.Sum256(asn1)
		base := hex.EncodeToString(sum[:])
		name := base + ".json"
		return filepath.Join(this.dir, this.Name, "cache", name), nil
	}
}

func (this *ProfileOp) open(
	n string,
	mode int,
	callback func(*os.File) (*Profile, error),
) (*Profile, error) {
	return openFileAt(this.dir, n, mode, callback)
}

// wrapper of OS `openat(2)`
func openFileAt[
	T any,
](
	dir string,
	n string,
	mode int,
	callback func(*os.File) (T, error),
) (
	ret T,
	err error,
) {
	var zero T

	if (mode & os.O_CREATE) != 0 {
		if err := os.MkdirAll(dir, 0o700); err != nil {
			return zero, Wrapf(err, "failed to create directory %+v", dir)
		}
	}

	root, err := os.OpenRoot(dir)
	if err != nil {
		return zero, Wrapf(err, "failed to open directory %+v", dir)
	}
	defer func() { _ = root.Close() }()

	if (mode & os.O_CREATE) != 0 {
		dirname := filepath.Dir(n)
		if err := root.MkdirAll(dirname, 0o700); err != nil {
			return zero, Wrapf(err, "failed to create directory %+v", dirname)
		}
	}

	file, err := root.OpenFile(n, mode, 0o600)
	if err != nil {
		return zero, Wrapf(err, "failed to open %+v", n)
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
	i := slices.Values(envp)
	j := intoSeq2(i, func(e string) (string, string, bool) { return strings.Cut(e, "=") })
	_, v, ok := findFirst(j, func(k, _ string) bool { return k == key })
	return v, ok
}
