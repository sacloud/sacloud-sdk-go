// Copyright 2026- The security-control-api-go Authors
// SPDX-License-Identifier: Apache-2.0

package securitycontrol_test

import (
	"os"
	"testing"

	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	securitycontrol "github.com/sacloud/security-control-api-go"
	v1 "github.com/sacloud/security-control-api-go/apis/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEvaluationRulesAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_SERVICE_PRINCIPAL_ID", "SAKURA_SERVICE_PRINCIPAL_KEY_KID")(t)

	ctx := t.Context()
	var theClient saclient.Client
	client, err := securitycontrol.NewClient(&theClient)
	require.NoError(t, err)

	erOp := securitycontrol.NewEvaluationRulesOp(client)
	rules, err := erOp.List(ctx, v1.EvaluationRulesListParams{PageSize: v1.NewOptInt(10)})
	require.NoError(t, err)
	require.NotNil(t, rules)
	assert.NotEmpty(t, rules.Items)

	respRead, err := erOp.Read(ctx, "server-no-public-ip")
	require.NoError(t, err)
	require.NotNil(t, respRead)
	assert.Equal(t, "server-no-public-ip", string(respRead.Rule.OneOf.Type))

	zones := []string{"is1a", "is1b", "tk1a", "tk1b"}
	respUpdate, err := erOp.Update(ctx, "server-no-public-ip", securitycontrol.SetupEvaluationRuleInput(&securitycontrol.EvaluationRuleInputParams{
		ID:                 "server-no-public-ip",
		ServicePrincipalID: os.Getenv("SAKURA_SERVICE_PRINCIPAL_ID"),
		Targets:            zones,
		Enabled:            !respRead.IsEnabled,
	}))
	require.NoError(t, err)
	require.NotNil(t, respUpdate)
	assert.Equal(t, !respRead.IsEnabled, respUpdate.IsEnabled)
	assert.Equal(t, "server-no-public-ip", string(respUpdate.Rule.OneOf.Type))
	if respUpdate.IsEnabled {
		assert.Equal(t, []string{"resource-viewer"}, respUpdate.IamRolesRequired)
		assert.Equal(t, zones, respUpdate.Rule.OneOf.ServerNoPublicIP.Parameter.Value.Zones)
	} else {
		assert.Empty(t, respUpdate.Rule.OneOf.ServerNoPublicIP.Parameter.Value.Zones)
	}
}
