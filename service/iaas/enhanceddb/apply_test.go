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

package enhanceddb

import (
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/sacloud/iaas-service-go/enhanceddb/builder"
	sacloudTestUtil "github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func TestEnhancedDBService_convertApplyRequest(t *testing.T) {
	caller := testutil.SingletonAPICaller()
	name := testutil.ResourceName("container-registry-service")
	dbName := sacloudTestUtil.Random(10, sacloudTestUtil.CharSetAlpha)
	password := sacloudTestUtil.Random(16, sacloudTestUtil.CharSetAlpha)

	cases := []struct {
		in     *ApplyRequest
		expect *builder.Builder
	}{
		{
			in: &ApplyRequest{
				Name:         name,
				Description:  "desc",
				Tags:         types.Tags{"tag1", "tag2"},
				DatabaseName: dbName,
				Password:     password,
				SettingsHash: "aaaaaaaa",
			},
			expect: &builder.Builder{
				ID:           0,
				Name:         name,
				Description:  "desc",
				Tags:         types.Tags{"tag1", "tag2"},
				IconID:       0,
				DatabaseName: dbName,
				Password:     password,
				SettingsHash: "aaaaaaaa",
				Client:       iaas.NewEnhancedDBOp(caller),
			},
		},
	}

	for _, tc := range cases {
		builder, err := tc.in.Builder(caller)
		require.NoError(t, err)
		require.EqualValues(t, tc.expect, builder)
	}
}
