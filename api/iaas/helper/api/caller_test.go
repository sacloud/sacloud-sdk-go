// Copyright 2022-2025 The sacloud/iaas-api-go Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package api

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	client "github.com/sacloud/api-client-go"
	"github.com/sacloud/api-client-go/profile"
	"github.com/stretchr/testify/require"
)

func initTestProfileDir() func() {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	os.Setenv("SAKURACLOUD_PROFILE_DIR", wd) //nolint
	profileDir := filepath.Join(wd, ".usacloud")
	if _, err := os.Stat(profileDir); err == nil {
		os.RemoveAll(profileDir) //nolint
	}

	return func() {
		os.RemoveAll(profileDir) //nolint
	}
}

func Test_defaultCallerOption(t *testing.T) {
	type args struct {
		options *client.Options
	}
	tests := []struct {
		name           string
		profiles       map[string]*profile.ConfigValue
		currentProfile string
		args           args // Note: currentProfileが指定されてたらoptionはそちらから読まれる
		want           *CallerOptions
	}{
		{
			name: "minimum",
			args: args{
				options: &client.Options{},
			},
			want: &CallerOptions{
				Options:       &client.Options{},
				APIRootURL:    "",
				DefaultZone:   "",
				Zones:         nil,
				TraceAPI:      false,
				FakeMode:      false,
				FakeStorePath: "",
			},
		},
		{
			name: "with profile",
			profiles: map[string]*profile.ConfigValue{
				"test-profile": {
					AccessToken:       "token",
					AccessTokenSecret: "secret",
					Zone:              "is1b",
					APIRootURL:        "https://api.example.com/",
				},
			},
			currentProfile: "test-profile",
			args:           args{},
			want: &CallerOptions{
				Options: &client.Options{
					AccessToken:       "token",
					AccessTokenSecret: "secret",
				},
				APIRootURL: "https://api.example.com/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer initTestProfileDir()()

			options := tt.args.options

			for profileName, profileValue := range tt.profiles {
				if err := profile.Save(profileName, profileValue); err != nil {
					t.Fatal(err)
				}
			}

			if tt.currentProfile != "" {
				if err := profile.SetCurrentName(tt.currentProfile); err != nil {
					t.Fatal(err)
				}

				// カレントプロファイルが指定されているときはそこから読む
				opts, err := client.DefaultOption()
				if err != nil {
					t.Fatal(err)
				}
				options = opts
			}

			got := defaultCallerOption(options)
			require.EqualValues(t, got.AccessToken, tt.want.AccessToken)
			require.EqualValues(t, got.AccessTokenSecret, tt.want.AccessTokenSecret)
			require.EqualValues(t, got.APIRootURL, tt.want.APIRootURL)
			require.EqualValues(t, got.DefaultZone, tt.want.DefaultZone)
			require.EqualValues(t, got.Zones, tt.want.Zones)
		})
	}
}
