// Copyright 2025- The sacloud/apigw-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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

func TestServiceAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN",
		"SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_TEST_HOST")(t)

	var theClient saclient.Client
	client, err := apigw.NewClient(&theClient)
	require.NoError(t, err)

	subReq, err := getSvcSubRequest()
	require.Nil(t, err)

	ctx := context.Background()
	serviceOp := apigw.NewServiceOp(client)

	// Create a service for testing
	serviceReq := v1.ServiceDetailRequest{
		Name:         "test-service",
		Host:         os.Getenv("SAKURA_TEST_HOST"),
		Port:         v1.NewOptInt(80),
		Protocol:     "http",
		Subscription: subReq,
	}
	created, err := serviceOp.Create(ctx, &serviceReq)
	require.NoError(t, err)
	defer func() { _ = serviceOp.Delete(ctx, created.ID.Value) }()

	serviceUpd := v1.ServiceDetail{
		Name:     "test-service-updated",
		Host:     os.Getenv("SAKURA_TEST_HOST"),
		Port:     v1.NewOptInt(80),
		Protocol: "http",
	}
	err = serviceOp.Update(ctx, &serviceUpd, created.ID.Value)
	assert.NoError(t, err)

	updated, err := serviceOp.Read(ctx, created.ID.Value)
	assert.NoError(t, err)
	assert.Equal(t, "test-service-updated", string(updated.Name))

	services, err := serviceOp.List(ctx)
	assert.NoError(t, err)
	assert.Greater(t, len(services), 0)

	err = serviceOp.Delete(ctx, created.ID.Value)
	assert.NoError(t, err)
}
