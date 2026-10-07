package secretmanager_test

import (
	"context"
	"os"
	"sort"
	"testing"

	sm "github.com/sacloud/sacloud-sdk-go/api/secretmanager"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVaultAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_KMS_KEY_ID")(t)

	client, err := sm.NewClient(&theClient)
	require.NoError(t, err)

	ctx := context.Background()
	keyId := os.Getenv("SAKURA_KMS_KEY_ID")
	vaultOp := sm.NewVaultOp(client)

	resCreate, err := vaultOp.Create(ctx, sm.CreateVaultParams{
		Name:        "vault from go",
		Description: new("vault from go client"),
		KmsKeyID:    keyId,
		Tags:        []string{"App", "Vault"},
	})
	require.NoError(t, err)
	assert.Equal(t, "vault from go", resCreate.Name)
	vaultID, ok := resCreate.ID.Get()
	require.True(t, ok)

	resList, err := vaultOp.List(ctx, new(10), new(0))
	assert.NoError(t, err)

	sort.Slice(resList, func(i, j int) bool { return resList[i].ID < resList[j].ID })
	found := false
	for _, vault := range resList {
		if vault.ID == vaultID {
			require.Equal(t, "vault from go client", vault.Description)
			found = true
		}
	}
	assert.True(t, found, "created vault not found in list")

	_, err = vaultOp.Update(ctx, vaultID, sm.UpdateVaultParams{
		Name:        "vault from go 2",
		Description: new("vault from go client 2"),
		Tags:        []string{"Test"},
	})
	assert.NoError(t, err)

	resRead, err := vaultOp.Read(ctx, vaultID)
	assert.NoError(t, err)
	assert.Equal(t, "vault from go 2", resRead.Name)
	assert.Equal(t, "vault from go client 2", resRead.Description)

	err = vaultOp.Delete(ctx, vaultID)
	require.NoError(t, err)
}
