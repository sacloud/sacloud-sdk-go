// users_test.go
package apigw_test

import (
	"context"
	"testing"

	apigw "github.com/sacloud/sacloud-sdk-go/api/apigw"
	v1 "github.com/sacloud/sacloud-sdk-go/api/apigw/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/stretchr/testify/assert"
)

func TestUserAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	assert.NoError(t, err)

	ctx := context.Background()
	userOp := apigw.NewUserOp(client)

	// Create
	userReq := &v1.UserDetail{
		Name:     "testuser",
		CustomID: v1.NewOptString("testuser-custom-id"),
		Tags:     v1.Tags{"Test"},
	}
	createdUser, err := userOp.Create(ctx, userReq)
	assert.NoError(t, err)

	gotUser, err := userOp.Read(ctx, createdUser.ID.Value)
	assert.NoError(t, err)
	assert.Equal(t, createdUser.ID.Value, gotUser.ID.Value)

	userReq.CustomID.SetTo("testuser-custom-id-updated")
	err = userOp.Update(ctx, userReq, createdUser.ID.Value)
	assert.NoError(t, err)

	users, err := userOp.List(ctx)
	assert.NoError(t, err)
	assert.Greater(t, len(users), 0)
	assert.Equal(t, "testuser-custom-id-updated", users[0].CustomID.Value)

	err = userOp.Delete(ctx, createdUser.ID.Value)
	assert.NoError(t, err)
}

func TestUserExtraAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	assert.NoError(t, err)

	ctx := context.Background()

	userOp := apigw.NewUserOp(client)
	createdUser, err := userOp.Create(ctx, &v1.UserDetail{
		Name:     "extrauser",
		CustomID: v1.NewOptString("extrauser-custom-id"),
		Tags:     v1.Tags{"Test"},
	})
	assert.NoError(t, err)
	defer func() { _ = userOp.Delete(ctx, createdUser.ID.Value) }()

	groupOp := apigw.NewGroupOp(client)
	createdGroup, err := groupOp.Create(ctx, &v1.Group{
		Name: v1.NewOptName("test-group"), Tags: []string{"Test"}})
	assert.NoError(t, err)
	defer func() { _ = groupOp.Delete(ctx, createdGroup.ID.Value) }()

	userExtraOp := apigw.NewUserExtraOp(client, createdUser.ID.Value)

	// ListGroup: 登録した User が所属している Group の一覧を取得します。とドキュメントにあるが全て返ってくるので要確認
	/*
		groups, err := userExtraOp.ListGroup(ctx)
		fmt.Println(groups)
		assert.NoError(t, err)
		assert.Equal(t, 0, len(groups))
	*/

	err = userExtraOp.UpdateGroup(ctx, string(createdGroup.Name.Value), true)
	assert.NoError(t, err)
	err = userExtraOp.UpdateGroup(ctx, createdGroup.ID.Value.String(), false)
	assert.NoError(t, err)

	err = userExtraOp.UpdateAuth(ctx, v1.UserAuthentication{
		BasicAuth: v1.NewOptBasicAuth(v1.BasicAuth{UserName: "test-user", Password: "test-password"}),
	})
	assert.NoError(t, err)

	gotAuth, err := userExtraOp.ReadAuth(ctx)
	assert.NoError(t, err)
	assert.Equal(t, "test-user", gotAuth.BasicAuth.Value.UserName)
}
