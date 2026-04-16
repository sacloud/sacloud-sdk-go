// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package apprun_dedicated_test

import (
	"testing"

	. "github.com/sacloud/apprun-dedicated-api-go"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func TestNewClient(t *testing.T) {
	assert := require.New(t)

	var theClient saclient.Client
	actual, err := NewClient(&theClient)
	assert.NoError(err)
	assert.NotNil(actual)
}
