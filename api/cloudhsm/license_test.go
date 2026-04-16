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
	"net/http"
	"testing"

	. "github.com/sacloud/cloudhsm-api-go"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func newTestLicenseClient(resp interface{}, status ...int) *v1.Client {
	return newTestClient(resp, status...)
}

func TestLicenseOp_List(t *testing.T) {
	assert := require.New(t)
	expected := v1.PaginatedCloudHSMSoftwareLicenseList{
		Count:    1,
		From:     v1.NewOptInt(0),
		Total:    v1.NewOptInt(1),
		Licenses: []v1.CloudHSMSoftwareLicense{TemplateLicense},
	}
	client := newTestLicenseClient(expected)
	api := NewLicenseOp(client)
	ctx := context.Background()
	licenses, err := api.List(ctx)

	assert.NoError(err)
	assert.NotNil(licenses)
	assert.Equal(1, len(licenses))
}

func TestLicenseOp_Read(t *testing.T) {
	assert := require.New(t)
	client := newTestLicenseClient(TemplateWrappedLicense)
	api := NewLicenseOp(client)
	ctx := context.Background()

	res, err := api.Read(ctx, "12345")
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(TemplateWrappedLicense.GetLicense().Value, *res)
}

func TestLicenseOp_Read_404(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("No License matches the given query.")
	client := newTestLicenseClient(expected, http.StatusNotFound)
	api := NewLicenseOp(client)
	ctx := context.Background()

	license, err := api.Read(ctx, "99999")
	assert.Nil(license)
	assert.Error(err)
	assert.ErrorContains(err, "not found")
}

func TestLicenseOp_Create(t *testing.T) {
	assert := require.New(t)
	client := newTestLicenseClient(TemplateWrappedCreateLicense, http.StatusCreated)
	api := NewLicenseOp(client)
	ctx := context.Background()

	res, err := api.Create(ctx, CloudHSMSoftwareLicenseCreateParams{
		Name:        "Test License",
		Description: ref("This is a test license"),
		Tags: []string{
			"tag1",
			"tag2",
		},
	})
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(TemplateWrappedCreateLicense.GetLicense().Value, *res)
}

func TestLicenseOp_Create_422(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Invalid request body.")
	client := newTestLicenseClient(expected, http.StatusUnprocessableEntity)
	api := NewLicenseOp(client)
	ctx := context.Background()

	license, err := api.Create(ctx, CloudHSMSoftwareLicenseCreateParams{})
	assert.Nil(license)
	assert.Error(err)
	assert.ErrorContains(err, "invalid")
}

func TestLicenseOp_Update(t *testing.T) {
	assert := require.New(t)
	client := newTestLicenseClient(TemplateWrappedLicense)
	api := NewLicenseOp(client)
	ctx := context.Background()

	res, err := api.Update(ctx, "12345", CloudHSMSoftwareLicenseUpdateParams{
		Description: "Updated Description",
		Name:        "Updated Name",
		Tags: []string{
			"tag1",
			"tag2",
		},
	})
	assert.NoError(err)
	assert.NotNil(res)
	assert.Equal(TemplateWrappedLicense.GetLicense().Value, *res)
}

func TestLicenseOp_Update_400(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Invalid update parameters.")
	client := newTestLicenseClient(expected, http.StatusUnprocessableEntity)
	api := NewLicenseOp(client)
	ctx := context.Background()

	license, err := api.Update(ctx, "0", CloudHSMSoftwareLicenseUpdateParams{})
	assert.Nil(license)
	assert.Error(err)
	assert.ErrorContains(err, "invalid")
}

func TestLicenseOp_Delete(t *testing.T) {
	assert := require.New(t)
	client := newTestLicenseClient(nil, http.StatusNoContent)
	api := NewLicenseOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "12345")
	assert.NoError(err)
}

func TestLicenseOp_Delete_400(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Not found")
	client := newTestLicenseClient(expected, http.StatusNotFound)
	api := NewLicenseOp(client)
	ctx := context.Background()

	err := api.Delete(ctx, "0")
	assert.Error(err)
	assert.ErrorContains(err, "not found")
}

func TestLicenseIntegrated(t *testing.T) {
	assert := require.New(t)
	client := newIntegratedClient(t)
	api := NewLicenseOp(client)
	ctx := context.Background()

	// Create
	created, err := api.Create(ctx, CloudHSMSoftwareLicenseCreateParams{
		Name:        testutil.RandomName("test-license-", 16, testutil.CharSetAlphaNum),
		Description: ref(testutil.Random(128, testutil.CharSetAlphaNum)),
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
	licenses, err := api.List(ctx)
	assert.NoError(err)
	assert.NotNil(licenses)
	assert.NotEmpty(licenses)

	// Update
	newDesc := "updated integration test License"
	updateReq := CloudHSMSoftwareLicenseUpdateParams{
		Description: newDesc,
		Name:        read.GetName(),
	}
	updated, err := api.Update(ctx, created.GetID(), updateReq)
	assert.NoError(err)
	assert.NotNil(updated)
	assert.Equal(newDesc, updated.GetDescription())
}
