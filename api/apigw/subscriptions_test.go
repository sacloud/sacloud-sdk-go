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
	"os"
	"testing"

	apigw "github.com/sacloud/apigw-api-go"
	v1 "github.com/sacloud/apigw-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getSvcSubRequest() (v1.ServiceSubscriptionRequest, error) {
	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	if err != nil {
		return v1.ServiceSubscriptionRequest{}, err
	}

	ctx := context.Background()
	subOp := apigw.NewSubscriptionOp(client)

	list, err := subOp.List(ctx)
	if err != nil || len(list) == 0 {
		return v1.ServiceSubscriptionRequest{}, err
	}

	return v1.ServiceSubscriptionRequest{
		ID: list[0].ID.Value,
	}, nil
}

func TestSubscriptionAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.Nil(t, err)

	ctx := context.Background()
	subscriptionOp := apigw.NewSubscriptionOp(client)

	plans, err := subscriptionOp.ListPlans(ctx)
	if err != nil {
		t.Error(err.Error())
	}
	require.Nil(t, err)
	require.Greater(t, len(plans), 0)

	if os.Getenv("ENABLE_SAKURA_APIGW_SUBSCRIPTION_TEST") != "1" {
		return
	}

	err = subscriptionOp.Create(ctx, plans[0].ID.Value, "sdk-test")
	require.Nil(t, err)

	subs, err := subscriptionOp.List(ctx)
	require.Nil(t, err)

	status, err := subscriptionOp.Read(ctx, subs[0].ID.Value)
	assert.Nil(t, err)
	assert.Equal(t, string(status.Name.Value), "sdk-test")

	err = subscriptionOp.Delete(ctx, subs[0].ID.Value)
	assert.Nil(t, err)
}
