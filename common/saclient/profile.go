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

package saclient

import (
	"bytes"
	"crypto/sha256"
	"crypto/x509"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"iter"
	"maps"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/goccy/go-yaml"
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

	// The directory where profiles are stored
	Dir() string
}

// A (loaded) profile
type Profile struct {
	dir string `yaml:"-"`

	// Name (intuitive)
	Name string `yaml:"-"`

	Version int32 `yaml:"version,omitempty"`

	// Profile contents V1
	Credentials struct {
		AccessToken            option[string] `yaml:"access_token,omitempty"`
		AccessTokenSecret      option[string] `yaml:"access_token_secret,omitempty"`
		ServicePrincipalID     option[string] `yaml:"service_principal_id,omitempty"`
		ServicePrincipalKeyKID option[string] `yaml:"service_principal_key_kid,omitempty"`
		PrivateKey             option[string] `yaml:"private_key,omitempty"`
		PrivateKeyPEMPath      option[string] `yaml:"private_key_path,omitempty"`
	} `yaml:"credentials,omitempty"`
	Endpoints map[string]string `yaml:"endpoints,omitempty"`
	Cli       struct {
		ArgumentMatchMode  option[string] `yaml:"argument_match_mode,omitempty"`
		DefaultOutputType  option[string] `yaml:"default_output_type,omitempty"`
		DefaultQueryDriver option[string] `yaml:"default_query_driver,omitempty"`
		NoColor            option[bool]   `yaml:"no_color,omitempty"`
		ProcessTimeoutSec  option[int64]  `yaml:"process_timeout_sec,omitempty"`
	} `yaml:"cli,omitempty"`
	Go struct {
		APIRootURL           option[string]   `yaml:"api_root_url,omitempty"`
		AcceptLanguage       option[string]   `yaml:"accept_language,omitempty"`
		DefaultZone          option[string]   `yaml:"default_zone,omitempty"`
		FakeMode             option[bool]     `yaml:"fake_mode,omitempty"`
		FakeStorePath        option[string]   `yaml:"fake_store_path,omitempty"`
		HTTPRequestRateLimit option[int64]    `yaml:"http_request_rate_limit,omitempty"`
		HTTPRequestTimeout   option[int64]    `yaml:"http_request_timeout,omitempty"`
		RetryMax             option[int64]    `yaml:"retry_max,omitempty"`
		RetryWaitMax         option[int64]    `yaml:"retry_wait_max,omitempty"`
		RetryWaitMin         option[int64]    `yaml:"retry_wait_min,omitempty"`
		StatePollingInterval option[int64]    `yaml:"state_polling_interval,omitempty"`
		StatePollingTimeout  option[int64]    `yaml:"state_polling_timeout,omitempty"`
		TraceMode            option[string]   `yaml:"trace_mode,omitempty"`
		Zone                 option[string]   `yaml:"zone,omitempty"`
		Zones                option[[]string] `yaml:"zones,omitempty"`
	} `yaml:"go,omitempty"`

	// Profile contents V0
	//
	// This is intentionally untyped to allow arbitrary key/values
	Attributes map[string]any `yaml:"-"`
}

type ProfileOp struct {
	// profile directory (historically `~/.usacloud` but configurable via env var)
	// Note that this is not necessarily existing at process startup.
	dir string
}

var _ ProfileAPI = (*ProfileOp)(nil)

const (
	profileConfigNameV1YAML = "config.yaml"
	profileConfigNameV0JSON = "config.json"
)

var profileConfigOrder = []string{
	profileConfigNameV1YAML,
	profileConfigNameV0JSON,
}

var profileV1CredentialsKeyMap = map[string]string{
	"access_token":              "AccessToken",
	"access_token_secret":       "AccessTokenSecret",
	"service_principal_id":      "ServicePrincipalID",
	"service_principal_key_kid": "ServicePrincipalKeyKID",
	"private_key":               "PrivateKey",
	"private_key_path":          "PrivateKeyPEMPath",
}

var profileV1CLIKeyMap = map[string]string{
	"argument_match_mode":  "ArgumentMatchMode",
	"default_output_type":  "DefaultOutputType",
	"default_query_driver": "DefaultQueryDriver",
	"no_color":             "NoColor",
	"process_timeout_sec":  "ProcessTimeoutSec",
}

var profileV1GoKeyMap = map[string]string{
	"api_root_url":            "APIRootURL",
	"accept_language":         "AcceptLanguage",
	"default_zone":            "DefaultZone",
	"fake_mode":               "FakeMode",
	"fake_store_path":         "FakeStorePath",
	"http_request_rate_limit": "HTTPRequestRateLimit",
	"http_request_timeout":    "HTTPRequestTimeout",
	"retry_max":               "RetryMax",
	"retry_wait_max":          "RetryWaitMax",
	"retry_wait_min":          "RetryWaitMin",
	"state_polling_interval":  "StatePollingInterval",
	"state_polling_timeout":   "StatePollingTimeout",
	"trace_mode":              "TraceMode",
	"zone":                    "Zone",
	"zones":                   "Zones",
}

var profileV0ToV1SectionMap = map[string]v0ProfileField{
	"AccessToken":            {section: "credentials", key: "access_token"},
	"AccessTokenSecret":      {section: "credentials", key: "access_token_secret"},
	"ServicePrincipalID":     {section: "credentials", key: "service_principal_id"},
	"ServicePrincipalKeyID":  {section: "credentials", key: "service_principal_key_kid"},
	"ServicePrincipalKeyKID": {section: "credentials", key: "service_principal_key_kid"},
	"PrivateKey":             {section: "credentials", key: "private_key"},
	"PrivateKeyPEMPath":      {section: "credentials", key: "private_key_path"},
	"ArgumentMatchMode":      {section: "cli", key: "argument_match_mode"},
	"DefaultOutputType":      {section: "cli", key: "default_output_type"},
	"DefaultQueryDriver":     {section: "cli", key: "default_query_driver"},
	"NoColor":                {section: "cli", key: "no_color"},
	"ProcessTimeoutSec":      {section: "cli", key: "process_timeout_sec"},
	"APIRootURL":             {section: "go", key: "api_root_url"},
	"AcceptLanguage":         {section: "go", key: "accept_language"},
	"DefaultZone":            {section: "go", key: "default_zone"},
	"FakeMode":               {section: "go", key: "fake_mode"},
	"FakeStorePath":          {section: "go", key: "fake_store_path"},
	"HTTPRequestRateLimit":   {section: "go", key: "http_request_rate_limit"},
	"HTTPRequestTimeout":     {section: "go", key: "http_request_timeout"},
	"RetryMax":               {section: "go", key: "retry_max"},
	"RetryWaitMax":           {section: "go", key: "retry_wait_max"},
	"RetryWaitMin":           {section: "go", key: "retry_wait_min"},
	"StatePollingInterval":   {section: "go", key: "state_polling_interval"},
	"StatePollingTimeout":    {section: "go", key: "state_polling_timeout"},
	"TraceMode":              {section: "go", key: "trace_mode"},
	"Zone":                   {section: "go", key: "zone"},
	"Zones":                  {section: "go", key: "zones"},
}

type v0ProfileField struct {
	section string
	key     string
}

// profileV1Document is the YAML representation of a profile, used for import/export.
type profileV1Document struct {
	Version     int            `yaml:"version"`
	Credentials map[string]any `yaml:"credentials,omitempty"`
	Endpoints   map[string]any `yaml:"endpoints,omitempty"`
	Cli         map[string]any `yaml:"cli,omitempty"`
	Go          map[string]any `yaml:"go,omitempty"`
	Extra       map[string]any `yaml:",inline,omitempty"`
}

// Creates a profile operator
func NewProfileOp(envp []string) (*ProfileOp, error) {
	dir, err := lookupProfileDir(envp)

	if err != nil {
		return nil, err
	}

	return &ProfileOp{dir}, nil
}

func (this *ProfileOp) List() ([]string, error) {
	if stat, err := os.Stat(this.dir); err != nil {
		return []string{}, nil // This is when e.g. the first invocation
	} else if !stat.IsDir() {
		return nil, NewErrorf("failed to open %+v", this.dir)
	} else {
		seen := make(map[string]struct{})

		for _, name := range profileConfigOrder {
			glob := filepath.Join(this.dir, "*", name)
			ent, err := filepath.Glob(glob)
			if err != nil {
				return nil, Wrapf(err, "failed to open %+v", this.dir)
			}

			for _, p := range ent {
				if stat, err := os.Stat(p); err == nil && stat.Mode().IsRegular() {
					base := filepath.Base(filepath.Dir(p))
					seen[base] = struct{}{}
				}
			}
		}

		ret := make([]string, 0, len(seen))
		for name := range seen {
			ret = append(ret, name)
		}

		slices.Sort(ret)
		return ret, nil
	}
}

func (this *ProfileOp) Read(name string) (*Profile, error) {
	for _, configName := range profileConfigOrder {
		n := filepath.Join(name, configName)

		profile, err := this.open(n, os.O_RDONLY, func(fp *os.File) (*Profile, error) {
			profile, err := decodeProfile(fp)
			if err != nil {
				return nil, err
			}

			profile.dir = this.dir
			profile.Name = name
			return profile, nil
		})

		if err == nil {
			return profile, nil
		}

		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		return nil, err
	}

	path := filepath.Join(name, profileConfigNameV0JSON)
	return nil, Wrapf(&os.PathError{Op: "open", Path: path, Err: os.ErrNotExist}, "failed to open %+v", path)
}

func (this *ProfileOp) Create(p *Profile) error {
	if _, err := this.Read(p.Name); err == nil {
		path := filepath.Join(p.Name, profileConfigNameV1YAML)
		return Wrapf(&os.PathError{Op: "open", Path: path, Err: os.ErrExist}, "failed to open %+v", path)
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	configName := profileConfigNameV0JSON
	if p.Version == 1 {
		configName = profileConfigNameV1YAML
	}
	n := filepath.Join(p.Name, configName)
	_, err := this.open(n, os.O_WRONLY|os.O_CREATE|os.O_TRUNC|os.O_EXCL, func(fp *os.File) (*Profile, error) {
		if err := encodeProfile(fp, p); err != nil {
			return nil, Wrapf(err, "failed to serialize %+v", p.Pathname())
		}

		return p, nil
	})
	return err
}

func (this *ProfileOp) Update(p *Profile) (*Profile, error) {
	var current map[string]any
	version := int32(-1)

	for _, configName := range profileConfigOrder {
		_, err := this.open(filepath.Join(p.Name, configName), os.O_RDONLY, func(fp *os.File) (*Profile, error) {
			profile, err := decodeProfile(fp)
			if err != nil {
				return nil, err
			}

			current = profile.Attributes
			version = profile.Version
			return nil, nil
		})
		if err == nil {
			break
		}
		if errors.Is(err, os.ErrNotExist) {
			continue
		}

		return nil, err
	}

	if version == -1 {
		return nil, Wrapf(&os.PathError{Op: "open", Path: p.Name, Err: os.ErrNotExist}, "failed to open profile %+v", p.Name)
	}

	var ret *Profile
	writePath := filepath.Join(p.Name, profileConfigNameV0JSON)
	// 既存の実装ではProfile.Attributes経由の差分更新を行なっていたので、versionが0の場合は過去の挙動を維持する。
	// v1以降は新しいProfileをそのまま使う。基本的にはReadした結果のProfileのフィールドを変更してUpdateに渡す挙動を想定する。
	if p.Version == 1 {
		ret = p
		writePath = filepath.Join(p.Name, profileConfigNameV1YAML)
	} else {
		ret = newProfile(this.dir, p.Name, version, merge(current, p.Attributes))
		if ret.Version == 1 {
			writePath = filepath.Join(p.Name, profileConfigNameV1YAML)
		}
	}

	_, err := this.open(writePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, func(fp *os.File) (*Profile, error) {
		if err := encodeProfile(fp, ret); err != nil {
			return nil, Wrapf(err, "failed to serialize %+v", p.Pathname())
		}

		return ret, nil
	})
	if err != nil {
		return nil, err
	}

	return ret, nil
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

	_, err := this.open("current", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, func(fp *os.File) (*Profile, error) {
		if _, err := fp.WriteString(name); err != nil {
			return nil, Wrapf(err, "failed to write %+v", fp.Name())
		} else {
			return nil, nil
		}
	})
	return err
}

func (this *ProfileOp) Dir() string { return this.dir }

// Calculated pathname of the configuration file
func (this *Profile) Pathname() string {
	if this.Version == 1 {
		return filepath.Join(this.dir, this.Name, profileConfigNameV1YAML)
	} else {
		return filepath.Join(this.dir, this.Name, profileConfigNameV0JSON)
	}
}

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
		// This call to MkdirAll leads to `G703` path traversal problem,
		// which is in fact true.   This is the only place in this entire
		// library where caution is needed.  `G703` is not a false positive,
		// but it's not a problem either.
		//
		// The intention of this function is exactly to create files (given
		// the `O_CREATE` flag) under this exact directory.  This is where
		// base directory is set and all following operations are forced to
		// be under this directory, via os.OpenRoot.
		//
		// #nosec G703 -- This is intentional
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
		return zero, Wrapf(&os.PathError{Op: "open", Path: n, Err: err}, "failed to open %+v", n)
	}
	defer func() { _ = file.Close() }()

	return callback(file)
}

func decodeProfile(fp *os.File) (*Profile, error) {
	buf, err := io.ReadAll(fp)
	if err != nil {
		return nil, Wrapf(err, "failed to read %+v", fp.Name())
	}

	trimmed := bytes.TrimSpace(buf)
	if len(trimmed) == 0 {
		return nil, NewErrorf("failed to parse %+v", fp.Name())
	}

	if isJSONProfile(trimmed, fp.Name()) {
		var attrs map[string]any
		if err := json.Unmarshal(trimmed, &attrs); err != nil {
			return nil, Wrapf(err, "failed to parse %+v", fp.Name())
		}
		return newProfile("", "", 0, attrs), nil
	}

	var doc profileV1Document
	if err := yaml.Unmarshal(trimmed, &doc); err != nil {
		return nil, Wrapf(err, "failed to parse %+v", fp.Name())
	}

	if doc.Version != 1 {
		return nil, NewErrorf("missing or unsupported version field %+v", fp.Name())
	}

	return importProfileV1(doc), nil
}

func encodeProfile(w io.Writer, p *Profile) error {
	if p.Version == 0 {
		return json.NewEncoder(w).Encode(p.Attributes)
	}

	doc := exportProfileV1(p)
	enc := yaml.NewEncoder(w)
	defer func() { _ = enc.Close() }()
	return enc.Encode(doc)
}

// v0ベースの設定から新しいプロファイルを作るのに使う。v1用の設定も同時に行う
func newProfile(dir, name string, version int32, attrs map[string]any) *Profile {
	p := &Profile{
		dir:        dir,
		Name:       name,
		Version:    version,
		Attributes: attrs,
	}
	p.populateV1Fields()
	return p
}

func importProfileV1(doc profileV1Document) *Profile {
	p := &Profile{
		Version:    1,
		Attributes: maps.Clone(doc.Extra),
	}
	if p.Attributes == nil {
		p.Attributes = make(map[string]any)
	}

	// 互換性のためv0向けのAttributesフィールドも更新する
	mergeMappedSection(doc.Credentials, profileV1CredentialsKeyMap, p.Attributes)
	mergeMappedSection(doc.Cli, profileV1CLIKeyMap, p.Attributes)
	mergeMappedSection(doc.Go, profileV1GoKeyMap, p.Attributes)
	if len(doc.Endpoints) > 0 {
		p.Attributes["Endpoints"] = maps.Clone(doc.Endpoints)
	}

	p.Credentials.AccessToken = profileStringOption(doc.Credentials, "access_token")
	p.Credentials.AccessTokenSecret = profileStringOption(doc.Credentials, "access_token_secret")
	p.Credentials.ServicePrincipalID = profileStringOption(doc.Credentials, "service_principal_id")
	p.Credentials.ServicePrincipalKeyKID = profileStringOption(doc.Credentials, "service_principal_key_kid")
	p.Credentials.PrivateKey = profileStringOption(doc.Credentials, "private_key")
	p.Credentials.PrivateKeyPEMPath = profileStringOption(doc.Credentials, "private_key_path")

	if len(doc.Endpoints) > 0 {
		p.Endpoints = make(map[string]string, len(doc.Endpoints))
		for key, value := range doc.Endpoints {
			if endpoint, ok := value.(string); ok {
				p.Endpoints[key] = endpoint
			}
		}
	}

	p.Cli.ArgumentMatchMode = profileStringOption(doc.Cli, "argument_match_mode")
	p.Cli.DefaultOutputType = profileStringOption(doc.Cli, "default_output_type")
	p.Cli.DefaultQueryDriver = profileStringOption(doc.Cli, "default_query_driver")
	p.Cli.NoColor = profileBoolOption(doc.Cli, "no_color")
	p.Cli.ProcessTimeoutSec = profileInt64Option(doc.Cli, "process_timeout_sec")

	p.Go.APIRootURL = profileStringOption(doc.Go, "api_root_url")
	p.Go.AcceptLanguage = profileStringOption(doc.Go, "accept_language")
	p.Go.DefaultZone = profileStringOption(doc.Go, "default_zone")
	p.Go.FakeMode = profileBoolOption(doc.Go, "fake_mode")
	p.Go.FakeStorePath = profileStringOption(doc.Go, "fake_store_path")
	p.Go.HTTPRequestRateLimit = profileInt64Option(doc.Go, "http_request_rate_limit")
	p.Go.HTTPRequestTimeout = profileInt64Option(doc.Go, "http_request_timeout")
	p.Go.RetryMax = profileInt64Option(doc.Go, "retry_max")
	p.Go.RetryWaitMax = profileInt64Option(doc.Go, "retry_wait_max")
	p.Go.RetryWaitMin = profileInt64Option(doc.Go, "retry_wait_min")
	p.Go.StatePollingInterval = profileInt64Option(doc.Go, "state_polling_interval")
	p.Go.StatePollingTimeout = profileInt64Option(doc.Go, "state_polling_timeout")
	p.Go.TraceMode = profileStringOption(doc.Go, "trace_mode")
	p.Go.Zone = profileStringOption(doc.Go, "zone")
	p.Go.Zones = profileStringsOption(doc.Go, "zones")

	return p
}

func exportProfileV1(p *Profile) profileV1Document {
	doc := profileV1Document{Version: 1}
	if p == nil {
		return doc
	}

	credentials := map[string]any{}
	setProfileOption(credentials, "access_token", p.Credentials.AccessToken)
	setProfileOption(credentials, "access_token_secret", p.Credentials.AccessTokenSecret)
	setProfileOption(credentials, "service_principal_id", p.Credentials.ServicePrincipalID)
	setProfileOption(credentials, "service_principal_key_kid", p.Credentials.ServicePrincipalKeyKID)
	setProfileOption(credentials, "private_key", p.Credentials.PrivateKey)
	setProfileOption(credentials, "private_key_path", p.Credentials.PrivateKeyPEMPath)
	if len(credentials) > 0 {
		doc.Credentials = credentials
	}

	if len(p.Endpoints) > 0 {
		doc.Endpoints = make(map[string]any, len(p.Endpoints))
		for key, value := range p.Endpoints {
			doc.Endpoints[key] = value
		}
	}

	cli := map[string]any{}
	setProfileOption(cli, "argument_match_mode", p.Cli.ArgumentMatchMode)
	setProfileOption(cli, "default_output_type", p.Cli.DefaultOutputType)
	setProfileOption(cli, "default_query_driver", p.Cli.DefaultQueryDriver)
	setProfileOption(cli, "no_color", p.Cli.NoColor)
	setProfileOption(cli, "process_timeout_sec", p.Cli.ProcessTimeoutSec)
	if len(cli) > 0 {
		doc.Cli = cli
	}

	g := map[string]any{}
	setProfileOption(g, "api_root_url", p.Go.APIRootURL)
	setProfileOption(g, "accept_language", p.Go.AcceptLanguage)
	setProfileOption(g, "default_zone", p.Go.DefaultZone)
	setProfileOption(g, "fake_mode", p.Go.FakeMode)
	setProfileOption(g, "fake_store_path", p.Go.FakeStorePath)
	setProfileOption(g, "http_request_rate_limit", p.Go.HTTPRequestRateLimit)
	setProfileOption(g, "http_request_timeout", p.Go.HTTPRequestTimeout)
	setProfileOption(g, "retry_max", p.Go.RetryMax)
	setProfileOption(g, "retry_wait_max", p.Go.RetryWaitMax)
	setProfileOption(g, "retry_wait_min", p.Go.RetryWaitMin)
	setProfileOption(g, "state_polling_interval", p.Go.StatePollingInterval)
	setProfileOption(g, "state_polling_timeout", p.Go.StatePollingTimeout)
	setProfileOption(g, "trace_mode", p.Go.TraceMode)
	setProfileOption(g, "zone", p.Go.Zone)
	setProfileOption(g, "zones", p.Go.Zones)
	if len(g) > 0 {
		doc.Go = g
	}

	// Go実装で使われないフィールドはAttributesに残っているため、それをExtraにコピーすることでフィールドを維持する。
	for key, value := range p.Attributes {
		if value != nil && !isV0ProfileKey(key) {
			if doc.Extra == nil {
				doc.Extra = make(map[string]any)
			}
			doc.Extra[key] = value
		}
	}

	return doc
}

func setProfileOption[T any](dst map[string]any, key string, value option[T]) {
	if value, ok := value.Get(); ok {
		dst[key] = value
	}
}

func isV0ProfileKey(key string) bool {
	if key == "Endpoints" {
		return true
	}
	_, ok := profileV0ToV1SectionMap[key]
	return ok
}

func (this *Profile) populateV1Fields() {
	if this.Attributes == nil {
		return
	}

	this.Credentials.AccessToken = profileStringOption(this.Attributes, "AccessToken")
	this.Credentials.AccessTokenSecret = profileStringOption(this.Attributes, "AccessTokenSecret")
	this.Credentials.ServicePrincipalID = profileStringOption(this.Attributes, "ServicePrincipalID")
	this.Credentials.ServicePrincipalKeyKID = profileStringOption(this.Attributes, "ServicePrincipalKeyKID")
	this.Credentials.PrivateKey = profileStringOption(this.Attributes, "PrivateKey")
	this.Credentials.PrivateKeyPEMPath = profileStringOption(this.Attributes, "PrivateKeyPEMPath")

	if endpoints, ok := this.Attributes["Endpoints"].(map[string]any); ok {
		this.Endpoints = make(map[string]string, len(endpoints))
		for key, value := range endpoints {
			if endpoint, ok := value.(string); ok {
				this.Endpoints[key] = endpoint
			}
		}
	}

	this.Cli.ArgumentMatchMode = profileStringOption(this.Attributes, "ArgumentMatchMode")
	this.Cli.DefaultOutputType = profileStringOption(this.Attributes, "DefaultOutputType")
	this.Cli.DefaultQueryDriver = profileStringOption(this.Attributes, "DefaultQueryDriver")
	this.Cli.NoColor = profileBoolOption(this.Attributes, "NoColor")
	this.Cli.ProcessTimeoutSec = profileInt64Option(this.Attributes, "ProcessTimeoutSec")

	this.Go.APIRootURL = profileStringOption(this.Attributes, "APIRootURL")
	this.Go.AcceptLanguage = profileStringOption(this.Attributes, "AcceptLanguage")
	this.Go.DefaultZone = profileStringOption(this.Attributes, "DefaultZone")
	this.Go.FakeMode = profileBoolOption(this.Attributes, "FakeMode")
	this.Go.FakeStorePath = profileStringOption(this.Attributes, "FakeStorePath")
	this.Go.HTTPRequestRateLimit = profileInt64Option(this.Attributes, "HTTPRequestRateLimit")
	this.Go.HTTPRequestTimeout = profileInt64Option(this.Attributes, "HTTPRequestTimeout")
	this.Go.RetryMax = profileInt64Option(this.Attributes, "RetryMax")
	this.Go.RetryWaitMax = profileInt64Option(this.Attributes, "RetryWaitMax")
	this.Go.RetryWaitMin = profileInt64Option(this.Attributes, "RetryWaitMin")
	this.Go.StatePollingInterval = profileInt64Option(this.Attributes, "StatePollingInterval")
	this.Go.StatePollingTimeout = profileInt64Option(this.Attributes, "StatePollingTimeout")
	this.Go.TraceMode = profileStringOption(this.Attributes, "TraceMode")
	this.Go.Zone = profileStringOption(this.Attributes, "Zone")
	this.Go.Zones = profileStringsOption(this.Attributes, "Zones")
}

func profileStringOption(attrs map[string]any, key string) option[string] {
	value, ok := attrs[key].(string)
	return option[string]{some: value, set: ok}
}

func profileBoolOption(attrs map[string]any, key string) option[bool] {
	value, ok := attrs[key].(bool)
	return option[bool]{some: value, set: ok}
}

func profileInt64Option(attrs map[string]any, key string) option[int64] {
	value, ok := attrs[key]
	if !ok {
		return option[int64]{}
	}

	var ret int64
	switch value := value.(type) {
	case int:
		ret = int64(value)
	case int64:
		ret = value
	case uint64:
		ret = int64(value) //nolint:gosec
	case float64:
		ret = int64(value)
	default:
		return option[int64]{}
	}
	return option[int64]{some: ret, set: true}
}

func profileStringsOption(attrs map[string]any, key string) option[[]string] {
	value, ok := attrs[key]
	if !ok {
		return option[[]string]{}
	}

	switch values := value.(type) {
	case []string:
		return option[[]string]{some: slices.Clone(values), set: true}
	case []any:
		ret := make([]string, 0, len(values))
		for _, value := range values {
			str, ok := value.(string)
			if !ok {
				return option[[]string]{}
			}
			ret = append(ret, str)
		}
		return option[[]string]{some: ret, set: true}
	default:
		return option[[]string]{}
	}
}

func isJSONProfile(trimmed []byte, pathname string) bool {
	ext := strings.ToLower(filepath.Ext(pathname))
	if ext == ".json" {
		return true
	}

	return len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[')
}

func mergeMappedSection(section map[string]any, keymap map[string]string, dst map[string]any) {
	for k, v := range section {
		if dstKey, ok := keymap[k]; ok {
			dst[dstKey] = v
		} else {
			dst[k] = v
		}
	}
}

func merge(dst, src map[string]any) map[string]any {
	ret := make(map[string]any)
	maps.Copy(ret, dst)
	maps.Copy(ret, src)
	return ret
}

func lookupProfileDir(envp []string) (string, error) {
	if v, ok := lookupEnv(envp, "SAKURA_PROFILE_DIR"); ok {
		return filepath.Clean(v), nil
	} else if v, ok := lookupEnv(envp, "SAKURACLOUD_PROFILE_DIR"); ok {
		return filepath.Clean(v), nil
	} else if v, ok := lookupEnv(envp, "USACLOUD_PROFILE_DIR"); ok {
		return filepath.Clean(v), nil // backward compat
	} else if v, ok := lookupEnv(envp, "XDG_CONFIG_HOME"); ok {
		// if, and only if `~/.config/usacloud` exists, take it.
		ret := filepath.Join(v, "usacloud")
		if stat, err := os.Stat(ret); err == nil && stat.IsDir() { // #nosec G703 -- this is in fact secure
			return filepath.Clean(ret), nil
		}
	}

	// fallback to '~/.usacloud'
	home, err := os.UserHomeDir()

	if err != nil {
		// :ESOTERIC: $HOME not set
		return "", err
	} else if home == "" {
		// :UNLIKELY: current implementation of `os.UserHomeDir`
		// does not return empty string without error.  But in case
		// it changes its mind in the future, we want to be defensive here.
		// We want to avoid creating global toplevel `/.usacloud`
		return "", NewErrorf("unable to determine home directory")
	}

	return filepath.Join(home, ".usacloud"), nil
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
