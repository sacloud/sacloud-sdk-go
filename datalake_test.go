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

func setupDataLake(t *testing.T, v encodable, s ...int) (*require.Assertions, DataLakeAPI) {
	assert := require.New(t)
	client := newTestClient(v, s...)
	api := NewDataLakeOp(client)
	return assert, api
}

func TestNewDataLakeOp(t *testing.T) {
	assert, api := setupDataLake(t, nil, http.StatusAccepted)
	assert.NotNil(api)
}

func TestDataLakeOp_List_200(t *testing.T) {
	assert, api := setupDataLake(t, &MockListResourcesResponse)

	result, err := api.List(t.Context())
	assert.NoError(err)
	assert.NotNil(result)
	assert.Len(result, 1)
	assert.Equal(&MockResourceGroupResource, &result[0])
}

func TestDataLakeOp_List_400(t *testing.T) {
	assert, api := setupDataLake(t, &MockErrorResponse, http.StatusBadRequest)

	result, err := api.List(t.Context())
	assert.Error(err)
	assert.Nil(result)
	assert.False(saclient.IsNotFoundError(err))
}

func TestDataLakeOp_List_404(t *testing.T) {
	assert, api := setupDataLake(t, nil, http.StatusNotFound)

	result, err := api.List(t.Context())
	assert.Error(err)
	assert.Nil(result)
	assert.True(saclient.IsNotFoundError(err))
}

func TestDataLakeOp_Create_202(t *testing.T) {
	assert, api := setupDataLake(t, &MockPostDeploymentResponse, http.StatusAccepted)

	result, err := api.Create(t.Context(), DataLakeCreateParams{
		Location:    "japaneast",
		Performance: v1.DataLakePerformance1,
		Redundancy:  v1.DataLakeRedundancy1,
	})
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(&MockPostDeploymentResponse, result)
}

func TestDataLakeOp_Create_400(t *testing.T) {
	assert, api := setupDataLake(t, &MockErrorResponse, http.StatusBadRequest)

	result, err := api.Create(t.Context(), DataLakeCreateParams{
		Location:    "invalid-location",
		Performance: v1.DataLakePerformance1,
		Redundancy:  v1.DataLakeRedundancy1,
	})
	assert.Error(err)
	assert.Nil(result)
	assert.False(saclient.IsNotFoundError(err))
}

func TestDataLakeOp_Read_200(t *testing.T) {
	assert, api := setupDataLake(t, &MockResourceResponse)

	result, err := api.Read(t.Context(), "test-rg")
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(&MockResourceResponse, result)
}

func TestDataLakeOp_Read_400(t *testing.T) {
	assert, api := setupDataLake(t, &MockErrorResponse, http.StatusBadRequest)

	result, err := api.Read(t.Context(), "invalid-location")
	assert.Error(err)
	assert.Nil(result)
	assert.False(saclient.IsNotFoundError(err))
}

func TestDataLakeOp_Read_404(t *testing.T) {
	assert, api := setupDataLake(t, nil, http.StatusNotFound)

	result, err := api.Read(t.Context(), "nonexistent-rg")
	assert.Error(err)
	assert.Nil(result)
}

func TestDataLakeOp_Delete_204(t *testing.T) {
	assert, api := setupDataLake(t, nil, http.StatusNoContent)

	err := api.Delete(t.Context(), "test-rg")
	assert.NoError(err)
}

func TestDataLakeOp_Delete_400(t *testing.T) {
	assert, api := setupDataLake(t, &MockErrorResponse, http.StatusBadRequest)

	err := api.Delete(t.Context(), "nonexistent-rg")
	assert.Error(err)
	assert.False(saclient.IsNotFoundError(err))
}

func TestDataLakeOp_Delete_404(t *testing.T) {
	assert, api := setupDataLake(t, nil, http.StatusNotFound)

	err := api.Delete(t.Context(), "nonexistent-rg")
	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestDataLakeOp_Status_200(t *testing.T) {
	assert, api := setupDataLake(t, &MockDeploymentStatus)

	result, err := api.Status(t.Context(), "test-rg", "test-deployment")
	assert.NoError(err)
	assert.NotNil(result)
	assert.Equal(&MockDeploymentStatus, result)
}

func TestDataLakeOp_Status_400(t *testing.T) {
	assert, api := setupDataLake(t, &MockErrorResponse, http.StatusBadRequest)

	result, err := api.Status(t.Context(), "test-rg", "nonexistent-deployment")
	assert.Error(err)
	assert.Nil(result)
	assert.False(saclient.IsNotFoundError(err))
}

func TestDataLakeOp_Status_404(t *testing.T) {
	assert, api := setupDataLake(t, nil, http.StatusNotFound)

	result, err := api.Status(t.Context(), "test-rg", "nonexistent-deployment")
	assert.Error(err)
	assert.Nil(result)
	assert.True(saclient.IsNotFoundError(err))
}
