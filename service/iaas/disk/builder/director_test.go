// Copyright 2022-2025 The sacloud/iaas-service-go Authors
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

package disk

import (
	"testing"

	"github.com/sacloud/iaas-api-go/ostype"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/require"
)

func TestDiskDirector_Builder(t *testing.T) {
	cases := []struct {
		name string
		in   *Director
		out  Builder
	}{
		{
			name: "blank disk",
			in:   &Director{},
			out:  &BlankBuilder{},
		},
		{
			name: "connect disk",
			in: &Director{
				DiskID: types.ID(1),
			},
			out: &ConnectedDiskBuilder{
				ID: types.ID(1),
			},
		},
		{
			name: "from archive",
			in: &Director{
				SourceArchiveID: types.ID(1),
			},
			out: &FromDiskOrArchiveBuilder{
				SourceArchiveID: types.ID(1),
			},
		},
		{
			name: "from disk",
			in: &Director{
				SourceDiskID: types.ID(1),
			},
			out: &FromDiskOrArchiveBuilder{
				SourceDiskID: types.ID(1),
			},
		},
		{
			name: "unix",
			in: &Director{
				OSType: ostype.Ubuntu,
				EditParameter: &EditRequest{
					HostName: "example",
				},
			},
			out: &FromUnixBuilder{
				OSType: ostype.Ubuntu,
				EditParameter: &UnixEditRequest{
					HostName: "example",
				},
			},
		},
	}

	for _, tt := range cases {
		builder := tt.in.Builder()
		require.Equal(t, tt.out, builder, tt.name)
	}
}
