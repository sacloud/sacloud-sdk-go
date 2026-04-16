package secretmanager_test

import (
	"context"
	"os"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sacloud/packages-go/testutil"
	sm "github.com/sacloud/secretmanager-api-go"
	v1 "github.com/sacloud/secretmanager-api-go/apis/v1"
)

func TestVaultAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_KMS_KEY_ID")(t)

	client, err := sm.NewClient(&theClient)
	require.NoError(t, err)

	ctx := context.Background()
	keyId := os.Getenv("SAKURA_KMS_KEY_ID")
	vaultOp := sm.NewVaultOp(client)

	resCreate, err := vaultOp.Create(ctx, v1.CreateVault{
		Name:        "vault from go",
		Description: v1.NewOptString("vault from go client"),
		KmsKeyID:    keyId,
		Tags:        []string{"App", "Vault"},
	})
	require.NoError(t, err)
	assert.Equal(t, "vault from go", resCreate.Name)

	resList, err := vaultOp.List(ctx)
	assert.NoError(t, err)

	sort.Slice(resList, func(i, j int) bool { return resList[i].ID < resList[j].ID })
	found := false
	for _, vault := range resList {
		if vault.ID == resCreate.ID {
			require.Equal(t, "vault from go client", vault.Description.Value)
			found = true
		}
	}
	assert.True(t, found, "created vault not found in list")

	_, err = vaultOp.Update(ctx, resCreate.ID, v1.Vault{
		Name:        "vault from go 2",
		Description: v1.NewOptString("vault from go client 2"),
		KmsKeyID:    keyId,
		Tags:        []string{"Test"},
	})
	assert.NoError(t, err)

	resRead, err := vaultOp.Read(ctx, resCreate.ID)
	assert.NoError(t, err)
	assert.Equal(t, "vault from go 2", resRead.Name)
	assert.Equal(t, "vault from go client 2", resRead.Description.Value)

	err = vaultOp.Delete(ctx, resCreate.ID)
	require.NoError(t, err)
}
