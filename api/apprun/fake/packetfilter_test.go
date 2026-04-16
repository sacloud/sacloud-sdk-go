// Copyright 2021-2024 The sacloud/apprun-api-go authors
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

package fake

import (
	"encoding/json"
	"testing"

	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestEngine_PacketFilter(t *testing.T) {
	t.Run("get and update packet filter", func(t *testing.T) {
		engine := NewEngine()
		req := postApplicationBody()
		createdApp, err := engine.CreateApplication(req)
		require.NoError(t, err)

		// 初期状態ではパケットフィルタは存在しない
		_, err = engine.ReadPacketFilter(createdApp.Id)
		require.Error(t, err)

		enabled := true
		settings := []v1.PacketFilterSetting{
			{
				FromIp:             "192.0.2.0",
				FromIpPrefixLength: 24,
			},
		}
		updated, err := engine.UpdatePacketFilter(createdApp.Id, &v1.PatchPacketFilter{
			IsEnabled: &enabled,
			Settings:  &settings,
		})

		require.NoError(t, err)

		respJson, err := json.Marshal(updated)
		require.NoError(t, err)

		expectedJSON := `
		{
			"is_enabled": true,
			"settings": [
				{
					"from_ip": "192.0.2.0",
					"from_ip_prefix_length": 24
				}
			]
		}`
		require.JSONEq(t, expectedJSON, string(respJson))

		got, err := engine.ReadPacketFilter(createdApp.Id)
		require.NoError(t, err)

		respJson, err = json.Marshal(got)
		require.NoError(t, err)
		require.JSONEq(t, expectedJSON, string(respJson))
	})
}
