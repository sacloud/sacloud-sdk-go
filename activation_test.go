// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol_test

import (
	"os"
	"testing"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	securitycontrol "github.com/sacloud/security-control-api-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActivationAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_SERVICE_PRINCIPAL_ID", "SAKURA_SERVICE_PRINCIPAL_KEY_KID")(t)

	ctx := t.Context()
	var theClient saclient.Client
	client, err := securitycontrol.NewClient(&theClient)
	require.NoError(t, err)

	actOp := securitycontrol.NewActivationOp(client)
	respRead, err := actOp.Read(ctx)
	if saclient.IsNotFoundError(err) { // アクティベーションされてないアカウントではReadで404が返るため、それを利用してCreateをテスト
		created, err := actOp.Create(ctx, os.Getenv("SAKURA_SERVICE_PRINCIPAL_ID"))
		require.NoError(t, err)
		require.NotNil(t, created)
		assert.True(t, created.IsActive)

		respUpdate, err := actOp.Update(ctx, os.Getenv("SAKURA_SERVICE_PRINCIPAL_ID"), false)
		require.NoError(t, err)
		require.NotNil(t, respUpdate)
		assert.False(t, respUpdate.IsActive)
	} else {
		require.NoError(t, err)
		require.NotNil(t, respRead)
		assert.True(t, respRead.AutomatedActionLimit > 0)

		respUpdate, err := actOp.Update(ctx, os.Getenv("SAKURA_SERVICE_PRINCIPAL_ID"), false)
		require.NoError(t, err)
		require.NotNil(t, respUpdate)
		assert.False(t, respUpdate.IsActive)
	}

	respUpdate, err := actOp.Update(ctx, os.Getenv("SAKURA_SERVICE_PRINCIPAL_ID"), true)
	require.NoError(t, err)
	require.NotNil(t, respUpdate)
	assert.True(t, respUpdate.IsActive)

	res, err := actOp.Read(ctx)
	require.NoError(t, err)
	require.NotNil(t, res)
	assert.True(t, res.IsActive)
	assert.True(t, res.AutomatedActionLimit > 0)
}
