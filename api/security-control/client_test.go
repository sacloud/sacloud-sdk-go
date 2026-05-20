// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol_test

import (
	"testing"

	. "github.com/sacloud/sacloud-sdk-go/api/security-control"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	assert := require.New(t)

	var theClient saclient.Client
	actual, err := NewClient(&theClient)
	assert.NoError(err)
	assert.NotNil(actual)
}
