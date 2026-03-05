// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite

import (
	"testing"

	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

var theClient saclient.Client

func TestNewClient(t *testing.T) {
	c, err := NewClient(&theClient)
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestNewClientWithApiUrl_OtherZone(t *testing.T) {
	c, err := NewClientWithApiUrl("https://secure.sakura.ad.jp/cloud/zone/tk1b/api/monitoring/1.0/", &theClient)
	require.NoError(t, err)
	require.NotNil(t, c)
}
