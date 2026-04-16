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

func TestProcessConfigurationAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_SIMPLE_NOTIFICATION_GROUP_ID")(t)

	var theClient saclient.Client
	client, err := eventbus.NewClient(&theClient)
	require.NoError(t, err)

	ctx := context.Background()
	pcOp := eventbus.NewProcessConfigurationOp(client)
	groupId := os.Getenv("SAKURA_SIMPLE_NOTIFICATION_GROUP_ID")

	resCreate, err := pcOp.Create(ctx, v1.CreateCommonServiceItemRequest{
		CommonServiceItem: v1.CreateCommonServiceItemRequestCommonServiceItem{
			Name: "SDK Test", Description: v1.NewOptNilString("SDK Testの概要"), Tags: []string{"test"},
			Settings: eventbus.CreateSimpleNotificationSettings(groupId, "メッセージ"),
		},
	})
	require.NoError(t, err)
	pcId := resCreate.ID

	resList, err := pcOp.List(ctx)
	assert.NoError(t, err)
	found := false
	for _, pc := range resList {
		if pc.ID == resCreate.ID {
			found = true
			assert.Equal(t, "SDK Test", pc.Name)
			assert.Equal(t, []string{"test"}, pc.Tags)
		}
	}
	assert.True(t, found, "Created ProcessConfiguration not found in list")

	_, err = pcOp.Update(ctx, pcId, v1.UpdateCommonServiceItemRequest{
		CommonServiceItem: v1.UpdateCommonServiceItemRequestCommonServiceItem{
			Name: v1.NewOptString("SDK Test 2"), Description: v1.NewOptNilString("SDK Test 2の概要"), Tags: []string{"test2"},
			Settings: v1.NewOptSettings(eventbus.CreateSimpleNotificationSettings(groupId, "メッセージ2")),
		},
	})
	assert.NoError(t, err)

	resRead, err := pcOp.Read(ctx, pcId)
	assert.NoError(t, err)
	assert.Equal(t, "SDK Test 2", resRead.Name)
	assert.Equal(t, []string{"test2"}, resRead.Tags)

	err = pcOp.UpdateSecret(ctx, pcId, v1.SetSecretRequest{
		Secret: v1.NewSacloudAPISecretSetSecretRequestSecret(v1.SacloudAPISecret{
			AccessToken:       os.Getenv("SAKURA_ACCESS_TOKEN"),
			AccessTokenSecret: os.Getenv("SAKURA_ACCESS_TOKEN_SECRET"),
		}),
	})
	assert.NoError(t, err)

	err = pcOp.Delete(ctx, pcId)
	require.NoError(t, err)
}
