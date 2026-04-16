// Copyright 2025- The sacloud/kms-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

func TestGroupAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.Nil(t, err)

	ctx := context.Background()
	groupOp := apigw.NewGroupOp(client)

	orig, err := groupOp.Create(ctx, &v1.Group{Name: v1.NewOptName("test-group"), Tags: []string{"Test"}})
	require.Nil(t, err)

	err = groupOp.Update(ctx, &v1.Group{Name: v1.NewOptName("test-group-updated"), Tags: []string{"SDK"}}, orig.ID.Value)
	assert.Nil(t, err)

	got, err := groupOp.Read(ctx, orig.ID.Value)
	assert.Nil(t, err)
	assert.Equal(t, orig.ID.Value, got.ID.Value)
	assert.Equal(t, v1.Name("test-group-updated"), got.Name.Value)

	groups, err := groupOp.List(ctx)
	assert.Nil(t, err)
	assert.Greater(t, len(groups), 0)

	err = groupOp.Delete(ctx, orig.ID.Value)
	assert.Nil(t, err)
}
