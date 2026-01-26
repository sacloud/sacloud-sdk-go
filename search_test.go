// Copyright 2025- The sacloud/addon-api-go Authors
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

package addon_test

import (
	"net/http"
	"testing"

	. "github.com/sacloud/addon-api-go"
	v1 "github.com/sacloud/addon-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setupSearch(t *testing.T, v encodable, s ...int) (*require.Assertions, SearchAPI) {
	assert := require.New(t)
	client := newTestClient(v, s...)
	api := NewSearchOp(client)
	return assert, api
}

func TestNewSearchOp(t *testing.T) {
	assert, api := setupSearch(t, nil, http.StatusAccepted)
	assert.NotNil(api)
}

func TestSearchOp_List_200(t *testing.T) {
	assert, api := setupSearch(t, &MockListResourcesResponse)

	result, err := api.List(t.Context())
	assert.NoError(err)
	assert.NotNil(result)
	assert.Len(result, 1)
	assert.Equal(&MockResourceGroupResource, &result[0])
}

func TestSearchOp_List_400(t *testing.T) {
	assert, api := setupSearch(t, &MockErrorResponse, http.StatusBadRequest)

	result, err := api.List(t.Context())
	assert.Error(err)
	assert.Nil(result)
	assert.False(saclient.IsNotFoundError(err))
}

func TestSearchOp_List_404(t *testing.T) {
	assert, api := setupSearch(t, nil, http.StatusNotFound)

	result, err := api.List(t.Context())
	assert.Error(err)
	assert.Nil(result)
	assert.True(saclient.IsNotFoundError(err))
}

func TestSearchOp_Create_202(t *testing.T) {
	assert, api := setupSearch(t, &MockPostDeploymentResponse, http.StatusAccepted)

	result, err := api.Create(t.Context(), SearchCreateParams{
		Location:       "japaneast",
		PartitionCount: 1,
		ReplicaCount:   3,
		Sku:            v1.SearchSku1,
	})
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(&MockPostDeploymentResponse, result)
}

func TestSearchOp_Create_400(t *testing.T) {
	assert, api := setupSearch(t, &MockErrorResponse, http.StatusBadRequest)

	result, err := api.Create(t.Context(), SearchCreateParams{
		Location:       "invalid-location",
		PartitionCount: 1,
		ReplicaCount:   3,
		Sku:            v1.SearchSku1,
	})
	assert.Error(err)
	assert.Nil(result)
	assert.False(saclient.IsNotFoundError(err))
}

func TestSearchOp_Read_200(t *testing.T) {
	assert, api := setupSearch(t, &MockResourceResponse)

	result, err := api.Read(t.Context(), "test-rg")
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(&MockResourceResponse, result)
}

func TestSearchOp_Read_400(t *testing.T) {
	assert, api := setupSearch(t, &MockErrorResponse, http.StatusBadRequest)

	result, err := api.Read(t.Context(), "invalid-location")
	assert.Error(err)
	assert.Nil(result)
	assert.False(saclient.IsNotFoundError(err))
}

func TestSearchOp_Read_404(t *testing.T) {
	assert, api := setupSearch(t, nil, http.StatusNotFound)

	result, err := api.Read(t.Context(), "nonexistent-rg")
	assert.Error(err)
	assert.Nil(result)
}

func TestSearchOp_Delete_204(t *testing.T) {
	assert, api := setupSearch(t, nil, http.StatusNoContent)

	err := api.Delete(t.Context(), "test-rg")
	assert.NoError(err)
}

func TestSearchOp_Delete_400(t *testing.T) {
	assert, api := setupSearch(t, &MockErrorResponse, http.StatusBadRequest)

	err := api.Delete(t.Context(), "nonexistent-rg")
	assert.Error(err)
	assert.False(saclient.IsNotFoundError(err))
}

func TestSearchOp_Delete_404(t *testing.T) {
	assert, api := setupSearch(t, nil, http.StatusNotFound)

	err := api.Delete(t.Context(), "nonexistent-rg")
	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestSearchOp_Status_200(t *testing.T) {
	assert, api := setupSearch(t, &MockDeploymentStatus)

	result, err := api.Status(t.Context(), "test-rg", "test-deployment")
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(&MockDeploymentStatus, result)
}

func TestSearchOp_Status_400(t *testing.T) {
	assert, api := setupSearch(t, &MockErrorResponse, http.StatusBadRequest)

	result, err := api.Status(t.Context(), "test-rg", "nonexistent-deployment")
	assert.Error(err)
	assert.Nil(result)
	assert.False(saclient.IsNotFoundError(err))
}

func TestSearchOp_Status_404(t *testing.T) {
	assert, api := setupSearch(t, nil, http.StatusNotFound)

	result, err := api.Status(t.Context(), "test-rg", "nonexistent-deployment")
	assert.Error(err)
	assert.Nil(result)
}

func TestSearchOp_Integrated(t *testing.T) {
	assert, client := IntegraetdClient(t)
	api := NewSearchOp(client)

	// Create
	result, err := api.Create(t.Context(), SearchCreateParams{
		Location:       "japaneast",
		PartitionCount: 1,
		ReplicaCount:   1, // SearchSku1 の場合レプリカは1つまで
		Sku:            v1.SearchSku1,
	})
	assert.NoError(err)
	assert.NotNil(result)

	rg, ok := result.ResourceGroupName.Get()
	assert.True(ok)
	assert.NotEmpty(rg)

	dn, ok := result.DeploymentName.Get()
	assert.True(ok)
	assert.NotEmpty(dn)

	defer func() {
		// Delete
		err := api.Delete(t.Context(), rg)
		assert.NoError(err)
	}()

	// List
	list, err := api.List(t.Context())
	assert.NoError(err)
	assert.NotEmpty(list)

	// Status
	status, err := api.Status(t.Context(), rg, dn)
	assert.NoError(err)
	assert.NotNil(status)

	// Read
	// このレスポンスはタイミング依存
	// プロビジョニングが完了していたりしていなかったりする
	read, err := api.Read(t.Context(), rg)
	if saclient.IsNotFoundError(err) {
		assert.Nil(read)
	} else {
		assert.NoError(err)
		assert.NotNil(read)
	}
}
