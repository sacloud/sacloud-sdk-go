// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package limit_test

import (
	"testing"

	"github.com/go-faster/jx"
	. "github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/limit"
	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/apis/v1"
	apprun_test "github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, response interface{ Encode(*jx.Encoder) }, status ...int) (*require.Assertions, LimitAPI) {
	assert := require.New(t)
	client, err := apprun_test.NewTestClient(response, status...)
	assert.NoError(err)
	return assert, NewLimitOp(client)
}

func TestRead(t *testing.T) {
	var expected v1.ReadLimitsResponse
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context())

	assert.NoError(err)
	assert.Equal(expected.GetLimit(), *actual)
}

func TestRead_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, err := api.Read(t.Context())

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}
