// Copyright 2026 The sacloud/iaas-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"context"
	"testing"

	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
	"github.com/stretchr/testify/require"
)

func TestDatabaseOpGetParameterMetaInfo(t *testing.T) {
	tests := []struct {
		name string
		conf *iaas.DatabaseRemarkDBConfCommon
		want []*iaas.DatabaseParameterMeta
	}{
		{
			name: "without configuration",
			want: fakeDatabaseParameterMetaForMariaDB,
		},
		{
			name: "PostgreSQL",
			conf: &iaas.DatabaseRemarkDBConfCommon{
				DatabaseName: types.RDBMSTypesPostgreSQL.String(),
			},
			want: fakeDatabaseParameterMetaForPostgreSQL,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op := NewDatabaseOp()
			zone := "test-zone-database-parameter-meta-info"
			id := types.ID(len(tt.name))
			ds().Put(ResourceDatabase, zone, id, &iaas.Database{
				ID:   id,
				Conf: tt.conf,
			})

			parameter, err := op.GetParameter(context.Background(), zone, id)
			require.NoError(t, err)
			require.Equal(t, tt.want, parameter.MetaInfo)
		})
	}
}
