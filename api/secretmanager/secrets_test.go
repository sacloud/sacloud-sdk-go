package secretmanager_test

import (
	"context"
	"os"
	"sort"
	"strconv"
	"testing"

	sm "github.com/sacloud/sacloud-sdk-go/api/secretmanager"
	v1 "github.com/sacloud/sacloud-sdk-go/api/secretmanager/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var theClient saclient.Client

func TestSecretAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_KMS_KEY_ID")(t)

	client, err := sm.NewClient(&theClient)
	require.NoError(t, err)

	ctx := context.Background()
	keyId := os.Getenv("SAKURA_KMS_KEY_ID")
	vaultOp := sm.NewVaultOp(client)

	vault, err := vaultOp.Create(ctx, v1.CreateVault{
		Name:        "vault for secret test",
		Description: v1.NewOptString("vault for secret test"),
		KmsKeyID:    keyId,
		Tags:        []string{"Test"},
	})
	require.NoError(t, err)

	defer func() {
		_ = vaultOp.Delete(ctx, vault.ID)
	}()

	secOp := sm.NewSecretOp(client, vault.ID)

	for i := range 2 {
		resCreate, err := secOp.Create(ctx, v1.CreateSecret{
			Name:  "Sec1",
			Value: "SecretValue" + strconv.Itoa(i),
		})
		require.NoError(t, err)
		require.Equal(t, i+1, resCreate.LatestVersion)
	}
	resCreate, err := secOp.Create(ctx, v1.CreateSecret{
		Name:  "Sec2",
		Value: "SV22",
	})
	require.NoError(t, err)
	assert.Equal(t, 1, resCreate.LatestVersion)

	resList, err := secOp.List(ctx)
	assert.NoError(t, err)

	sort.Slice(resList, func(i, j int) bool { return resList[i].Name < resList[j].Name })
	require.Len(t, resList, 2)
	assert.Equal(t, "Sec1", resList[0].Name)
	assert.Equal(t, 2, resList[0].LatestVersion)
	assert.Equal(t, "Sec2", resList[1].Name)
	assert.Equal(t, 1, resList[1].LatestVersion)

	for i := range 2 {
		resUn, err := secOp.Unveil(ctx, v1.Unveil{
			Name:    "Sec1",
			Version: v1.NewOptNilInt(i + 1),
		})
		assert.NoError(t, err)
		assert.Equal(t, "SecretValue"+strconv.Itoa(i), resUn.Value)
	}

	err = secOp.Delete(ctx, v1.DeleteSecret{Name: "Sec1"})
	require.NoError(t, err)
	err = secOp.Delete(ctx, v1.DeleteSecret{Name: "Sec2"})
	require.NoError(t, err)
}
