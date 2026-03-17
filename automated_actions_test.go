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

func TestAutomatedActionsAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET",
		"SAKURA_SERVICE_PRINCIPAL_ID", "SAKURA_SERVICE_PRINCIPAL_KEY_ID", "SAKURA_SIMPLE_NOTIFICATION_GROUP_ID")(t)

	ctx := t.Context()

	var theClient saclient.Client
	client, err := securitycontrol.NewClient(&theClient)
	require.NoError(t, err)

	aaOp := securitycontrol.NewAutomatedActionsOp(client)
	list, err := aaOp.List(ctx, v1.AutomatedActionsListParams{PageSize: v1.NewOptInt(10)})
	require.NoError(t, err)
	require.NotNil(t, list)

	respCreate, err := aaOp.Create(ctx, &v1.AutomatedActionInput{
		Name:        "example-action",
		Description: v1.NewOptString("example automated action"),
		IsEnabled:   true,
		Action: v1.ActionDefinition{
			OneOf: v1.ActionDefinitionSum{
				Type: v1.ActionDefinitionSumType("simpleNotification"),
				ActionDefinitionSimpleNotification: v1.ActionDefinitionSimpleNotification{
					ActionType: v1.ActionDefinitionSimpleNotificationActionType("simpleNotification"),
					ActionParameter: v1.SakuraSimpleNotification{
						ServicePrincipalId:  os.Getenv("SAKURA_SERVICE_PRINCIPAL_ID"),
						NotificationGroupId: os.Getenv("SAKURA_SIMPLE_NOTIFICATION_GROUP_ID"),
					},
				},
			},
		},
		ExecutionCondition: `event.evaluationResult.status == "REJECTED" && event.evaluationResult.severity in ["HIGH", "CRITICAL"]`,
	})
	defer func() {
		err := aaOp.Delete(ctx, string(respCreate.AutomatedActionId))
		require.NoError(t, err)
	}()

	require.NoError(t, err)
	require.NotNil(t, respCreate)
	assert.NotEmpty(t, respCreate.AutomatedActionId)

	respRead, err := aaOp.Read(ctx, string(respCreate.AutomatedActionId))
	require.NoError(t, err)
	require.NotNil(t, respRead)
	assert.Equal(t, "example-action", respRead.Name)
	assert.Equal(t, "example automated action", respRead.Description.Value)
	assert.True(t, respRead.IsEnabled)
	assert.Equal(t, "simpleNotification", string(respRead.Action.OneOf.Type))
	assert.Equal(t, `event.evaluationResult.status == "REJECTED" && event.evaluationResult.severity in ["HIGH", "CRITICAL"]`, respRead.ExecutionCondition)

	respUpdate, err := aaOp.Update(ctx, string(respCreate.AutomatedActionId), &v1.AutomatedActionInput{
		Name:        "example-action",
		Description: v1.NewOptString("updated automated action"),
		IsEnabled:   false,
		Action: v1.ActionDefinition{
			OneOf: v1.ActionDefinitionSum{
				Type: v1.ActionDefinitionSumType("simpleNotification"),
				ActionDefinitionSimpleNotification: v1.ActionDefinitionSimpleNotification{
					ActionType: v1.ActionDefinitionSimpleNotificationActionType("simpleNotification"),
					ActionParameter: v1.SakuraSimpleNotification{
						ServicePrincipalId:  os.Getenv("SAKURA_SERVICE_PRINCIPAL_ID"),
						NotificationGroupId: os.Getenv("SAKURA_SIMPLE_NOTIFICATION_GROUP_ID"),
					},
				},
			},
		},
		ExecutionCondition: `event.evaluationResult.status == "REJECTED"`,
	})
	require.NoError(t, err)
	require.NotNil(t, respUpdate)
	assert.Equal(t, "updated automated action", respUpdate.Description.Value)
	assert.Equal(t, `event.evaluationResult.status == "REJECTED"`, respUpdate.ExecutionCondition)

	// List
	respList, err := aaOp.List(ctx, v1.AutomatedActionsListParams{PageSize: v1.NewOptInt(10)})
	require.NoError(t, err)
	require.NotNil(t, respList)
	assert.Equal(t, len(list.Items)+1, len(respList.Items))
}
