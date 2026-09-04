// Copyright 2026 The sacloud/iaas-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package fake

import (
	"testing"

	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
	"github.com/stretchr/testify/require"
)

func TestCloneServerDoesNotShareInterfaces(t *testing.T) {
	source := &iaas.Server{
		Interfaces: []*iaas.InterfaceView{{ID: types.ID(1)}},
	}

	cloned := cloneServer(source)
	source.Interfaces = append(source.Interfaces, &iaas.InterfaceView{ID: types.ID(2)})
	source.Interfaces[0].ID = types.ID(3)

	require.Len(t, cloned.Interfaces, 1)
	require.Equal(t, types.ID(1), cloned.Interfaces[0].ID)
}
