package eventbus_test

import (
	"context"
	"os"
	"testing"

	eventbus "github.com/sacloud/eventbus-api-go"
	v1 "github.com/sacloud/eventbus-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTriggerAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_SIMPLE_NOTIFICATION_GROUP_ID")(t)

	var theClient saclient.Client
	client, err := eventbus.NewClient(&theClient)
	require.NoError(t, err)

	ctx := context.Background()
	pcOp := eventbus.NewProcessConfigurationOp(client)
	triggerOp := eventbus.NewTriggerOp(client)
	groupId := os.Getenv("SAKURA_SIMPLE_NOTIFICATION_GROUP_ID")

	pc, err := pcOp.Create(ctx, v1.CreateCommonServiceItemRequest{
		CommonServiceItem: v1.CreateCommonServiceItemRequestCommonServiceItem{
			Name: "SDK Test", Description: v1.NewOptNilString("SDK Testの概要"), Tags: []string{"tag1", "tag2"},
			Settings: eventbus.CreateSimpleNotificationSettings(groupId, "メッセージ"),
		},
	})
	require.NoError(t, err)
	pcId := pc.ID
	defer func() {
		_ = pcOp.Delete(ctx, pcId)
	}()

	resCreate, err := triggerOp.Create(ctx, v1.CreateCommonServiceItemRequest{
		CommonServiceItem: v1.CreateCommonServiceItemRequestCommonServiceItem{
			Name: "SDK Test", Description: v1.NewOptNilString("SDK Testの概要"),
			Settings: v1.NewTriggerSettingsSettings(v1.TriggerSettings{
				ProcessConfigurationID: pcId,
				Source:                 "//eventbus.sakura.ad.jp/test",
				Types:                  v1.NewOptNilStringArray([]string{"test.instance.created"}),
				Conditions: v1.NewOptNilTriggerSettingsConditionsItemArray([]v1.TriggerSettingsConditionsItem{
					v1.NewTriggerConditionEqTriggerSettingsConditionsItem(v1.TriggerConditionEq{
						Key:    "key1",
						Op:     v1.TriggerConditionEqOpEq,
						Values: []string{"value1"},
					}),
				}),
			}),
		},
	})
	require.NoError(t, err)
	triggerId := resCreate.ID

	resList, err := triggerOp.List(ctx)
	assert.NoError(t, err)
	found := false
	for _, trigger := range resList {
		if trigger.ID == resCreate.ID {
			found = true
			assert.Equal(t, "SDK Test", trigger.Name)
			assert.Equal(t, "//eventbus.sakura.ad.jp/test", trigger.Settings.TriggerSettings.Source)
			assert.Equal(t, []string{"test.instance.created"}, trigger.Settings.TriggerSettings.Types.Value)
			assert.Equal(t, []v1.TriggerSettingsConditionsItem{
				v1.NewTriggerConditionEqTriggerSettingsConditionsItem(v1.TriggerConditionEq{
					Key:    "key1",
					Op:     v1.TriggerConditionEqOpEq,
					Values: []string{"value1"},
				}),
			}, trigger.Settings.TriggerSettings.Conditions.Value)
		}
	}
	assert.True(t, found, "Created Trigger not found in list")

	_, err = triggerOp.Update(ctx, triggerId, v1.UpdateCommonServiceItemRequest{
		CommonServiceItem: v1.UpdateCommonServiceItemRequestCommonServiceItem{
			Name: v1.NewOptString("SDK Test 2"), Description: v1.NewOptNilString("SDK Test 2の概要"), Tags: []string{"tag1", "tag2"},
			Settings: v1.NewOptSettings(v1.NewTriggerSettingsSettings(v1.TriggerSettings{
				ProcessConfigurationID: pcId,
				Source:                 "//eventbus.sakura.ad.jp/test-updated",
			})),
		},
	})
	assert.NoError(t, err)

	resRead, err := triggerOp.Read(ctx, triggerId)
	assert.NoError(t, err)
	assert.Equal(t, "SDK Test 2", resRead.Name)
	assert.Equal(t, "//eventbus.sakura.ad.jp/test-updated", resRead.Settings.TriggerSettings.Source)
	assert.Equal(t, ([]string)(nil), resRead.Settings.TriggerSettings.Types.Value)
	assert.Equal(t, ([]v1.TriggerSettingsConditionsItem)(nil), resRead.Settings.TriggerSettings.Conditions.Value)
	assert.Equal(t, []string{"tag1", "tag2"}, resRead.Tags)

	err = triggerOp.Delete(ctx, triggerId)
	require.NoError(t, err)
}
