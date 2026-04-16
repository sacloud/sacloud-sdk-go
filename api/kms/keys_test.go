package kms_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	kms "github.com/sacloud/kms-api-go"
	v1 "github.com/sacloud/kms-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
)

func TestKeyAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := kms.NewClient(&theClient)
	require.NoError(t, err, "failed to create client")

	ctx := context.Background()
	keyOp := kms.NewKeyOp(client)

	resCreate, err := keyOp.Create(ctx, v1.CreateKey{
		Name:        "key gen from go",
		Description: v1.NewOptString("key gen from go client"),
		KeyOrigin:   v1.KeyOriginEnumGenerated,
		Tags:        []string{"tag1", "tag2"},
	})
	require.NoError(t, err, "failed to create key")
	assert.Equal(t, "key gen from go", resCreate.Name)

	defer func() {
		err = keyOp.Delete(ctx, resCreate.ID)
		require.NoError(t, err, "failed to delete key")
	}()

	resList, err := keyOp.List(ctx)
	assert.NoError(t, err, "failed to list keys")

	found := false
	for _, key := range resList {
		if key.ID == resCreate.ID {
			found = true
			assert.Equal(t, "key gen from go client", key.Description)
		}
	}
	assert.True(t, found, "created key not found in list")

	updated, err := keyOp.Update(ctx, resCreate.ID, v1.Key{
		Name:        "key gen from go 2",
		Description: "key gen from go client 2",
		KeyOrigin:   v1.KeyOriginEnumGenerated,
		Tags:        []string{"Test"},
	})
	assert.NoError(t, err, "failed to update key")
	assert.Equal(t, "key gen from go 2", updated.Name)
	assert.Equal(t, "key gen from go client 2", updated.Description)
	assert.Equal(t, []string{"Test"}, updated.Tags)
	assert.Equal(t, v1.KeyStatusEnumActive, updated.Status)
	assert.Equal(t, 0, updated.LatestVersion.Value)

	plain := []byte("hello world!")
	cipher, err := keyOp.Encrypt(ctx, resCreate.ID, plain, v1.KeyEncryptAlgoEnumAes256Gcm)
	assert.NoError(t, err, "failed to encrypt data")

	decrypted, err := keyOp.Decrypt(ctx, resCreate.ID, cipher)
	assert.NoError(t, err, "failed to decrypt data")
	assert.Equal(t, plain, decrypted)

	rotated, err := keyOp.Rotate(ctx, resCreate.ID)
	assert.NoError(t, err, "failed to rotate key")
	assert.Equal(t, 1, rotated.LatestVersion.Value)

	err = keyOp.ChangeStatus(ctx, resCreate.ID, v1.ChangeKeyStatusStatusSuspended)
	assert.NoError(t, err, "failed to change key status")

	read, err := keyOp.Read(ctx, resCreate.ID)
	assert.NoError(t, err, "failed to read key for Rotate / ChangeStatus")
	assert.Equal(t, v1.KeyStatusEnumSuspended, read.Status)
	assert.Equal(t, 1, read.LatestVersion.Value)

	err = keyOp.ScheduleDestruction(ctx, resCreate.ID, 100)
	assert.Error(t, err, "schedule destruction: longer pending days must be an error")

	err = keyOp.ScheduleDestruction(ctx, resCreate.ID, 10)
	assert.NoError(t, err, "failed to schedule destruction")

	read, err = keyOp.Read(ctx, resCreate.ID)
	assert.NoError(t, err, "failed to read key for ScheduleDestruction")
	assert.Equal(t, v1.KeyStatusEnumPendingDestruction, read.Status)
}
