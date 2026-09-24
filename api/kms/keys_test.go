package kms_test

import (
	"context"
	"testing"

	kms "github.com/sacloud/sacloud-sdk-go/api/kms"
	v1 "github.com/sacloud/sacloud-sdk-go/api/kms/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeyAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := kms.NewClient(&theClient)
	require.NoError(t, err, "failed to create client")

	ctx := context.Background()
	keyOp := kms.NewKeyOp(client)

	resCreate, err := keyOp.Create(ctx, kms.CreateParams{
		Name:        "key gen from go",
		Description: new("key gen from go client"),
		Tags:        []string{"tag1", "tag2"},
	})
	require.NoError(t, err, "failed to create key")
	assert.Equal(t, "key gen from go", resCreate.Name)

	defer func() {
		err = keyOp.Delete(ctx, resCreate.ID)
		require.NoError(t, err, "failed to delete key")
	}()

	resList, err := keyOp.List(ctx, nil, nil)
	assert.NoError(t, err, "failed to list keys")

	found := false
	for _, key := range resList {
		if key.ID == resCreate.ID {
			found = true
			assert.Equal(t, "key gen from go client", key.Description)
		}
	}
	assert.True(t, found, "created key not found in list")

	updated, err := keyOp.Update(ctx, resCreate.ID, kms.UpdateParams{
		Name:        "key gen from go 2",
		Description: new("key gen from go client 2"),
		Tags:        []string{"Test"},
	})
	assert.NoError(t, err, "failed to update key")
	assert.Equal(t, "key gen from go 2", updated.Name)
	assert.Equal(t, "key gen from go client 2", updated.Description)
	assert.Equal(t, []string{"Test"}, updated.Tags)
	assert.Equal(t, v1.KeyStatusActive, updated.Status)
	assert.Equal(t, 0, updated.LatestVersion)

	plain := []byte("hello world!")
	cipher, err := keyOp.Encrypt(ctx, resCreate.ID, plain, v1.EncryptionRequestAlgoAes256Gcm)
	assert.NoError(t, err, "failed to encrypt data")

	decrypted, err := keyOp.Decrypt(ctx, resCreate.ID, cipher)
	assert.NoError(t, err, "failed to decrypt data")
	assert.Equal(t, plain, decrypted)

	rotated, err := keyOp.Rotate(ctx, resCreate.ID)
	assert.NoError(t, err, "failed to rotate key")
	assert.Equal(t, 1, rotated.LatestVersion)

	_, err = keyOp.ChangeStatus(ctx, resCreate.ID, v1.ChangeKeyStateRequestStatusSuspended)
	assert.NoError(t, err, "failed to change key status")

	read, err := keyOp.Read(ctx, resCreate.ID)
	assert.NoError(t, err, "failed to read key for Rotate / ChangeStatus")
	assert.Equal(t, v1.KeyStatusSuspended, read.Status)
	assert.Equal(t, 1, read.LatestVersion)

	_, err = keyOp.ScheduleDestruction(ctx, resCreate.ID, 100)
	assert.Error(t, err, "schedule destruction: longer pending days must be an error")

	_, err = keyOp.ScheduleDestruction(ctx, resCreate.ID, 10)
	assert.NoError(t, err, "failed to schedule destruction")

	read, err = keyOp.Read(ctx, resCreate.ID)
	assert.NoError(t, err, "failed to read key for ScheduleDestruction")
	assert.Equal(t, v1.KeyStatusPendingDestruction, read.Status)
}
