// users_test.go
package apigw_test

import (
	"context"
	"testing"

	apigw "github.com/sacloud/apigw-api-go"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.Nil(t, err)

	ctx := context.Background()
	userOp := apigw.NewUserOp(client)

	// Create
	userReq := &v1.UserDetail{
		Name:     "testuser",
		CustomID: v1.NewOptString("testuser-custom-id"),
		Tags:     v1.Tags{"Test"},
	}
	createdUser, err := userOp.Create(ctx, userReq)
	require.Nil(t, err)

	gotUser, err := userOp.Read(ctx, createdUser.ID.Value)
	assert.Nil(t, err)
	assert.Equal(t, createdUser.ID.Value, gotUser.ID.Value)

	userReq.CustomID.SetTo("testuser-custom-id-updated")
	err = userOp.Update(ctx, userReq, createdUser.ID.Value)
	assert.Nil(t, err)

	users, err := userOp.List(ctx)
	assert.Nil(t, err)
	assert.Greater(t, len(users), 0)
	assert.Equal(t, "testuser-custom-id-updated", users[0].CustomID.Value)

	err = userOp.Delete(ctx, createdUser.ID.Value)
	assert.Nil(t, err)
}

func TestUserExtraAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.Nil(t, err)

	ctx := context.Background()

	userOp := apigw.NewUserOp(client)
	createdUser, err := userOp.Create(ctx, &v1.UserDetail{
		Name:     "extrauser",
		CustomID: v1.NewOptString("extrauser-custom-id"),
		Tags:     v1.Tags{"Test"},
	})
	require.Nil(t, err)
	defer func() { _ = userOp.Delete(ctx, createdUser.ID.Value) }()

	groupOp := apigw.NewGroupOp(client)
	createdGroup, err := groupOp.Create(ctx, &v1.Group{
		Name: v1.NewOptName("test-group"), Tags: []string{"Test"}})
	require.Nil(t, err)
	defer func() { _ = groupOp.Delete(ctx, createdGroup.ID.Value) }()

	userExtraOp := apigw.NewUserExtraOp(client, createdUser.ID.Value)

	// ListGroup: 登録した User が所属している Group の一覧を取得します。とドキュメントにあるが全て返ってくるので要確認
	/*
		groups, err := userExtraOp.ListGroup(ctx)
		fmt.Println(groups)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(groups))
	*/

	err = userExtraOp.UpdateGroup(ctx, string(createdGroup.Name.Value), true)
	assert.Nil(t, err)
	err = userExtraOp.UpdateGroup(ctx, createdGroup.ID.Value.String(), false)
	assert.Nil(t, err)

	err = userExtraOp.UpdateAuth(ctx, v1.UserAuthentication{
		BasicAuth: v1.NewOptBasicAuth(v1.BasicAuth{UserName: "test-user", Password: "test-password"}),
	})
	assert.Nil(t, err)

	gotAuth, err := userExtraOp.ReadAuth(ctx)
	assert.Nil(t, err)
	assert.Equal(t, "test-user", gotAuth.BasicAuth.Value.UserName)
}
