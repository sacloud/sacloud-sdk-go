// Copyright 2022-2025 The sacloud/services Authors
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

// Package examples
//
// パラメータの例示のための汎用的な値を提供する
package examples

import (
	"strings"
)

var (
	Id          = "123456789012"
	Name        = "example"
	Description = "example"
	Tags        = []string{"tag1", "tag2=value"}

	IpAddress        = "192.0.2.11"
	IpAddresses      = []string{"192.0.2.21", "192.0.2.22"}
	VirtualIpAddress = "192.0.2.101"
	NetworkMaskLen   = 24
	DefaultRoute     = "192.0.2.1"

	MacAddress   = "00:00:5E:00:53:01"
	MacAddresses = []string{"00:00:5E:00:53:11", "00:00:5E:00:53:12"}

	SlackNotifyWebhooksURL = "https://hooks.slack.com/services/XXXXXXXXX/XXXXXXXXX/XXXXXXXXXXXXXXXXXXXXXXXX"

	ScriptContent = `#!/bin/bash

...`

	PublicKeyContent = "ssh-rsa ..."

	FilePath = "/path/to/file"
)

func OptionsString(opts []string) string {
	return strings.Join(opts, " | ")
}
