// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol_test

import (
	"testing"

	"github.com/sacloud/saclient-go"
	. "github.com/sacloud/security-control-api-go"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	assert := require.New(t)

	var theClient saclient.Client
	actual, err := NewClient(&theClient)
	assert.NoError(err)
	assert.NotNil(actual)
}
