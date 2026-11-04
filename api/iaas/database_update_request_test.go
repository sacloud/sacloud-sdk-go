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
	"bytes"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDatabaseUpdateRequestOmitEmptyPasswords(t *testing.T) {
	tests := []struct {
		name                string
		password            string
		replicaPassword     string
		wantPasswordPresent bool
		wantReplicaPresent  bool
	}{
		{name: "omitted"},
		{
			name:                "explicit passwords",
			password:            "new-password",
			replicaPassword:     "new-replica-password",
			wantPasswordPresent: true,
			wantReplicaPresent:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request, err := (&DatabaseOp{}).transformUpdateArgs(1, &DatabaseUpdateRequest{
				CommonSetting: &DatabaseSettingCommon{
					DefaultUser:     "default",
					UserPassword:    tt.password,
					ReplicaPassword: tt.replicaPassword,
				},
			})
			require.NoError(t, err)

			payload, err := json.Marshal(request)
			require.NoError(t, err)
			require.Equal(t, tt.wantPasswordPresent, bytes.Contains(payload, []byte(`"UserPassword"`)))
			require.Equal(t, tt.wantReplicaPresent, bytes.Contains(payload, []byte(`"ReplicaPassword"`)))
		})
	}
}
