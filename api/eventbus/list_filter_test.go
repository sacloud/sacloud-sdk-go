package eventbus_test

import (
	"context"
	"os"
	"testing"
	"time"

	eventbus "github.com/sacloud/sacloud-sdk-go/api/eventbus"
	v1 "github.com/sacloud/sacloud-sdk-go/api/eventbus/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFilterByProviderClass(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_SIMPLE_NOTIFICATION_GROUP_ID")(t)

	var theClient saclient.Client
	client, err := eventbus.NewClient(&theClient)
	require.NoError(t, err)

	ctx := context.Background()
	pcOp := eventbus.NewProcessConfigurationOp(client)
	schedOp := eventbus.NewScheduleOp(client)
	triggerOp := eventbus.NewTriggerOp(client)
	groupId := os.Getenv("SAKURA_SIMPLE_NOTIFICATION_GROUP_ID")

	pc, err := pcOp.Create(ctx, v1.CreateCommonServiceItemRequest{
		CommonServiceItem: v1.CreateCommonServiceItemRequestCommonServiceItem{
			Name:        "SDK List Filter Test",
			Description: v1.NewOptNilString("SDK List Filter Test"),
			Tags:        []string{"test"},
			Settings:    eventbus.CreateSimpleNotificationSettings(groupId, "メッセージ"),
		},
	})
	require.NoError(t, err)
	defer func() {
		_ = pcOp.Delete(ctx, pc.ID)
	}()

	schedule, err := schedOp.Create(ctx, v1.CreateCommonServiceItemRequest{
		CommonServiceItem: v1.CreateCommonServiceItemRequestCommonServiceItem{
			Name:        "SDK List Filter Test",
			Description: v1.NewOptNilString("SDK List Filter Test"),
			Settings: v1.NewScheduleSettingsSettings(v1.ScheduleSettings{
				ProcessConfigurationID: pc.ID,
				RecurringStep:          v1.NewOptInt(5),
				RecurringUnit:          v1.NewOptScheduleSettingsRecurringUnit(v1.ScheduleSettingsRecurringUnitMin),
				StartsAt:               v1.NewInt64ScheduleSettingsStartsAt(time.Now().UnixMilli()),
			}),
		},
	})
	require.NoError(t, err)
	defer func() {
		_ = schedOp.Delete(ctx, schedule.ID)
	}()

	trigger, err := triggerOp.Create(ctx, v1.CreateCommonServiceItemRequest{
		CommonServiceItem: v1.CreateCommonServiceItemRequestCommonServiceItem{
			Name:        "SDK List Filter Test",
			Description: v1.NewOptNilString("SDK List Filter Test"),
			Settings: v1.NewTriggerSettingsSettings(v1.TriggerSettings{
				ProcessConfigurationID: pc.ID,
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
	defer func() {
		_ = triggerOp.Delete(ctx, trigger.ID)
	}()

	pcList, err := pcOp.List(ctx)
	require.NoError(t, err)
	for _, item := range pcList {
		assert.True(t, item.Settings.IsProcessConfigurationSettings(),
			"ProcessConfiguration.List returned non-ProcessConfiguration item: ID=%s", item.ID)
	}

	scheduleList, err := schedOp.List(ctx)
	require.NoError(t, err)
	for _, item := range scheduleList {
		assert.True(t, item.Settings.IsScheduleSettings(),
			"Schedule.List returned non-Schedule item: ID=%s", item.ID)
	}

	triggerList, err := triggerOp.List(ctx)
	require.NoError(t, err)
	for _, item := range triggerList {
		assert.True(t, item.Settings.IsTriggerSettings(),
			"Trigger.List returned non-Trigger item: ID=%s", item.ID)
	}
}
