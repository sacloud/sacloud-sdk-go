// Copyright 2025- The sacloud/cloudhsm-api-go Authors
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

package cloudhsm_test

import (
	"context"
	"fmt"
	"math/rand/v2"
	"net/http"
	"testing"

	. "github.com/sacloud/cloudhsm-api-go"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func newTestCloudHSMClient(resp interface{}, status ...int) *v1.Client {
	return newTestClient(resp, status...)
}

func TestCloudHSMOp_List(t *testing.T) {
	assert := require.New(t)
	expected := v1.PaginatedCloudHSMList{
		Count:     1,
		From:      v1.NewOptInt(0),
		Total:     v1.NewOptInt(1),
		CloudHSMs: []v1.CloudHSM{TemplateCloudHSM},
	}
	client := newTestCloudHSMClient(expected)
	api := NewCloudHSMOp(client)
	ctx := context.Background()
	cloudhsms, err := api.List(ctx)

	assert.NoError(err)
	assert.NotNil(cloudhsms)
	assert.Equal(1, len(cloudhsms))
}

func TestCloudHSMOp_Read(t *testing.T) {
	assert := require.New(t)
	client := newTestCloudHSMClient(TemplateWrappedCloudHSM)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	res, err := api.Read(ctx, "12345")
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(TemplateWrappedCloudHSM.GetCloudHSM(), *res)
}

func TestCloudHSMOp_Read_404(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("No CloudHSM matches the given query.")
	client := newTestCloudHSMClient(expected, http.StatusNotFound)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	cloudhsm, err := api.Read(ctx, "99999")
	assert.Nil(cloudhsm)
	assert.Error(err)
	assert.ErrorContains(err, "not found")
}

func TestCloudHSMOp_Create(t *testing.T) {
	assert := require.New(t)
	client := newTestCloudHSMClient(TemplateWrappedCreateCloudHSM, http.StatusCreated)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	res, err := api.Create(ctx, CloudHSMCreateParams{
		Name:        "Test HSM",
		Description: ref("This is a test HSM"),
		Tags: []string{
			"tag1",
			"tag2",
		},
	})
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(TemplateWrappedCreateCloudHSM.GetCloudHSM(), *res)
}

func TestCloudHSMOp_Create_422(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Invalid request body.")
	client := newTestCloudHSMClient(expected, http.StatusUnprocessableEntity)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	cloudhsm, err := api.Create(ctx, CloudHSMCreateParams{})
	assert.Nil(cloudhsm)
	assert.Error(err)
	assert.ErrorContains(err, "invalid")
}

func TestCloudHSMOp_Update(t *testing.T) {
	assert := require.New(t)
	client := newTestCloudHSMClient(TemplateWrappedCloudHSM)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	res, err := api.Update(ctx, "12345", CloudHSMUpdateParams{
		Description: ref("Updated Description"),
		Name:        "Updated Name",
		Tags: []string{
			"tag1",
			"tag2",
		},
	})
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(TemplateWrappedCloudHSM.GetCloudHSM(), *res)
}

func TestCloudHSMOp_Update_400(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Invalid update parameters.")
	client := newTestCloudHSMClient(expected, http.StatusUnprocessableEntity)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	cloudhsm, err := api.Update(ctx, "0", CloudHSMUpdateParams{})
	assert.Nil(cloudhsm)
	assert.Error(err)
	assert.ErrorContains(err, "invalid")
}

func TestCloudHSMOp_Delete(t *testing.T) {
	assert := require.New(t)
	client := newTestCloudHSMClient(nil, http.StatusNoContent)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "12345")
	assert.NoError(err)
}

func TestCloudHSMOp_Delete_400(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Not found")
	client := newTestCloudHSMClient(expected, http.StatusNotFound)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "0")
	assert.Error(err)
	assert.ErrorContains(err, "not found")
}

//nolint:gosec // no security issue here
func TestCloudHSMIntegrated(t *testing.T) {
	assert := require.New(t)
	client := newIntegratedClient(t)
	api := NewCloudHSMOp(client)
	ctx := context.Background()

	// Create
	created, err := api.Create(ctx, CloudHSMCreateParams{
		Name:        testutil.RandomName("test-cloudhsm-", 16, testutil.CharSetAlphaNum),
		Description: ref(testutil.Random(128, testutil.CharSetAlphaNum)),
		// This IP address is arbitrary, but recommended to be in the private range.
		Ipv4NetworkAddress: fmt.Sprintf("172.%d.%d.0", rand.Uint32N(31), rand.Uint32N(255)),
		Ipv4PrefixLength:   28,
	})
	assert.NoError(err)
	assert.NotNil(created)

	// Delete
	t.Cleanup(func() {
		err := api.Delete(ctx, created.GetID())
		assert.NoError(err)
	})

	// Read
	read, err := api.Read(ctx, created.GetID())
	assert.NoError(err)
	assert.NotNil(read)
	assert.Equal(created.GetID(), read.GetID())
	assert.Equal(created.GetName(), read.GetName())

	// List
	cloudhsms, err := api.List(ctx)
	assert.NoError(err)
	assert.NotNil(cloudhsms)
	assert.NotEmpty(cloudhsms)

	// Update
	newDesc := "updated integration test CloudHSM"
	updateReq := CloudHSMUpdateParams{
		Description:        ref(newDesc),
		Name:               read.GetName(),
		Ipv4NetworkAddress: read.Ipv4NetworkAddress,
		Ipv4PrefixLength:   read.Ipv4PrefixLength,
	}
	updated, err := api.Update(ctx, created.GetID(), updateReq)
	assert.NoError(err)
	assert.NotNil(updated)
	assert.Equal(newDesc, updated.GetDescription().Or("failure"))
}
