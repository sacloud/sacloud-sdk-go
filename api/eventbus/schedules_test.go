package eventbus_test

import (
	"context"
	"os"
	"testing"
	"time"

	eventbus "github.com/sacloud/eventbus-api-go"
	v1 "github.com/sacloud/eventbus-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduleAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_SIMPLE_NOTIFICATION_GROUP_ID")(t)

	var theClient saclient.Client
	client, err := eventbus.NewClient(&theClient)
	require.NoError(t, err)

	ctx := context.Background()
	pcOp := eventbus.NewProcessConfigurationOp(client)
	schedOp := eventbus.NewScheduleOp(client)
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

	resCreate, err := schedOp.Create(ctx, v1.CreateCommonServiceItemRequest{
		CommonServiceItem: v1.CreateCommonServiceItemRequestCommonServiceItem{
			Name: "SDK Test", Description: v1.NewOptNilString("SDK Testの概要"),
			Settings: v1.NewScheduleSettingsSettings(v1.ScheduleSettings{
				ProcessConfigurationID: pcId,
				RecurringStep:          v1.NewOptInt(5),
				RecurringUnit:          v1.NewOptScheduleSettingsRecurringUnit(v1.ScheduleSettingsRecurringUnitMin),
				StartsAt:               v1.NewInt64ScheduleSettingsStartsAt(time.Now().UnixMilli()),
			}),
		},
	})
	require.NoError(t, err)
	schedId := resCreate.ID

	resList, err := schedOp.List(ctx)
	assert.NoError(t, err)
	found := false
	for _, sched := range resList {
		if sched.ID == resCreate.ID {
			found = true
			assert.Equal(t, "SDK Test", sched.Name)
			assert.Equal(t, v1.NewOptScheduleSettingsRecurringUnit(v1.ScheduleSettingsRecurringUnitMin), sched.Settings.ScheduleSettings.RecurringUnit)
		}
	}
	assert.True(t, found, "Created Schedule not found in list")

	_, err = schedOp.Update(ctx, schedId, v1.UpdateCommonServiceItemRequest{
		CommonServiceItem: v1.UpdateCommonServiceItemRequestCommonServiceItem{
			Name: v1.NewOptString("SDK Test 2"), Description: v1.NewOptNilString("SDK Test 2の概要"), Tags: []string{"tag1", "tag2"},
			Settings: v1.NewOptSettings(v1.NewScheduleSettingsSettings(v1.ScheduleSettings{
				ProcessConfigurationID: pcId,
				RecurringStep:          v1.NewOptInt(1),
				RecurringUnit:          v1.NewOptScheduleSettingsRecurringUnit(v1.ScheduleSettingsRecurringUnitHour),
				StartsAt:               v1.NewInt64ScheduleSettingsStartsAt(time.Now().UnixMilli()),
			})),
		},
	})
	assert.NoError(t, err)

	resRead, err := schedOp.Read(ctx, schedId)
	assert.NoError(t, err)
	assert.Equal(t, "SDK Test 2", resRead.Name)
	assert.Equal(t, v1.NewOptScheduleSettingsRecurringUnit(v1.ScheduleSettingsRecurringUnitHour), resRead.Settings.ScheduleSettings.RecurringUnit)
	assert.Equal(t, []string{"tag1", "tag2"}, resRead.Tags)

	err = schedOp.Delete(ctx, schedId)
	require.NoError(t, err)
}
