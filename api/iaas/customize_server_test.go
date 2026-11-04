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

package iaas

import (
	"encoding/json"
	"testing"

	"github.com/sacloud/sacloud-sdk-go/api/iaas/mapconv"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/naked"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
	"github.com/sacloud/sacloud-sdk-go/common/packages/size"
	"github.com/stretchr/testify/require"
)

func TestServer_BandWidthAt(t *testing.T) {
	cases := []struct {
		msg    string
		in     *Server
		index  int
		expect int
	}{
		{
			msg:    "no NICs",
			in:     &Server{},
			expect: -1,
		},
		{
			msg: "invalid NIC index",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Unknown,
					},
				},
			},
			index:  1,
			expect: -1,
		},
		{
			msg: "unknown upstream type",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Unknown,
					},
				},
			},
			expect: -1,
		},
		{
			msg: "shared",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Shared,
					},
				},
			},
			expect: 100,
		},
		{
			msg: "switch with private host",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Switch,
					},
				},
				PrivateHostID: 1,
			},
			expect: 0,
		},
		{
			msg: "memory less than 32GB",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Switch,
					},
				},
				MemoryMB: 31 * size.GiB,
			},
			expect: 1000,
		},
		{
			msg: "switch with memory 32GB",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Switch,
					},
				},
				MemoryMB: 32 * size.GiB,
			},
			expect: 2000,
		},
		{
			msg: "memory less than 128GB",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Switch,
					},
				},
				MemoryMB: 127 * size.GiB,
			},
			expect: 2000,
		},
		{
			msg: "switch with memory 128GB",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Switch,
					},
				},
				MemoryMB: 128 * size.GiB,
			},
			expect: 5000,
		},
		{
			msg: "memory less than 224GB",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Switch,
					},
				},
				MemoryMB: 223 * size.GiB,
			},
			expect: 5000,
		},
		{
			msg: "switch with memory 224GB",
			in: &Server{
				Interfaces: []*InterfaceView{
					{
						UpstreamType: types.UpstreamNetworkTypes.Switch,
					},
				},
				MemoryMB: 224 * size.GiB,
			},
			expect: 10000,
		},
	}

	for _, tc := range cases {
		actual := tc.in.BandWidthAt(tc.index)
		require.Equal(t, tc.expect, actual, tc.msg)
	}
}

func TestServer_NetworkingSuiteInterfaceConnectionID(t *testing.T) {
	tests := []struct {
		name     string
		response string
		want     types.ID
	}{
		{
			name:     "connected",
			response: `{"Interfaces":[{"NetworkingSuiteInterfaceConnection":{"ID":123456789}}]}`,
			want:     types.ID(123456789),
		},
		{
			name:     "not connected",
			response: `{"Interfaces":[{}]}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var source naked.Server
			require.NoError(t, json.Unmarshal([]byte(tt.response), &source))

			var server Server
			require.NoError(t, mapconv.ConvertFrom(&source, &server))
			require.Len(t, server.Interfaces, 1)
			require.Equal(t, tt.want, server.Interfaces[0].NetworkingSuiteInterfaceConnectionID)
		})
	}
}
