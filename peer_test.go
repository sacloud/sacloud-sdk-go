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
	"os"
	"testing"

	. "github.com/sacloud/cloudhsm-api-go"
	v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
	"github.com/sacloud/packages-go/testutil"
	"github.com/stretchr/testify/require"
)

func newTestCloudHSMPeerClient(resp interface{}, status ...int) *v1.Client {
	return newTestClient(resp, status...)
}

func TestCloudHSMPeerOp_List(t *testing.T) {
	assert := require.New(t)
	expected := v1.CloudHSMPeerList{
		Peers: []v1.CloudHSMPeer{TemplateCloudHSMPeer},
	}
	client := newTestCloudHSMPeerClient(expected)
	api, err := NewPeerOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()
	peers, err := api.List(ctx)

	assert.NoError(err)
	assert.NotNil(peers)
	assert.Equal(1, len(peers))
}

func TestCloudHSMPeerOp_Create(t *testing.T) {
	assert := require.New(t)
	client := newTestCloudHSMPeerClient(nil, http.StatusNoContent)
	api, err := NewPeerOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	err = api.Create(ctx, CloudHSMPeerCreateParams{
		RouterID:  "peer-2",
		SecretKey: "secret-key-2",
	})
	assert.NoError(err)
}

func TestCloudHSMPeerOp_Create_422(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Invalid request body.")
	client := newTestCloudHSMPeerClient(expected, http.StatusUnprocessableEntity)
	api, err := NewPeerOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	err = api.Create(ctx, CloudHSMPeerCreateParams{})
	assert.Error(err)
	assert.ErrorContains(err, "invalid")
}

func TestCloudHSMPeerOp_Delete(t *testing.T) {
	assert := require.New(t)
	client := newTestCloudHSMPeerClient(nil, http.StatusNoContent)
	api, err := NewPeerOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	err = api.Delete(ctx, "12345")
	assert.NoError(err)
}

func TestCloudHSMPeerOp_Delete_400(t *testing.T) {
	assert := require.New(t)
	expected := newErrorResponse("Not found")
	client := newTestCloudHSMPeerClient(expected, http.StatusNotFound)
	api, err := NewPeerOp(client, &TemplateCloudHSM)
	assert.NoError(err)
	ctx := context.Background()

	err = api.Delete(ctx, "0")
	assert.Error(err)
	assert.ErrorContains(err, "not found")
}

func TestCloudHSMPeerIntegrated(t *testing.T) {
	assert := require.New(t)
	client := newIntegratedClient(t)

	testutil.PreCheckEnvsFunc("SAKURA_CLOUDHSM_ID")(t)

	ctx := context.Background()
	hsm, err := NewCloudHSMOp(client).Read(ctx, os.Getenv("SAKURA_CLOUDHSM_ID"))
	assert.NoError(err)
	assert.NotNil(hsm)
	assert.Equal(v1.AvailabilityEnumAvailable, hsm.GetAvailability())
	api, err := NewPeerOp(client, hsm)
	assert.NoError(err)

	// There is no way to know what ID is assigned to the created peer,
	// unless we list them all.
	// First save the existing peers...
	peers, err := api.List(ctx)
	assert.NoError(err)
	assert.NotNil(peers)
	existingPeerIDs := []string{}
	for _, p := range peers {
		existingPeerIDs = append(existingPeerIDs, p.GetID())
	}

	// Create
	err = api.Create(ctx, CloudHSMPeerCreateParams{
		RouterID:  "peer-integrated",
		SecretKey: "secret-integrated",
	})
	assert.NoError(err)

	// List again
	peers, err = api.List(ctx)
	assert.NoError(err)
	assert.NotNil(peers)
	assert.NotEmpty(peers)
	newPeerIDs := []string{}
	for _, p := range peers {
		newPeerIDs = append(newPeerIDs, p.GetID())
	}

	// find
	var createdPeerID string
	for _, i := range newPeerIDs {
		found := false
		for _, j := range existingPeerIDs {
			if i == j {
				found = true
				break
			}
		}
		if !found {
			createdPeerID = i
			break
		}
	}

	assert.NotEmpty(createdPeerID)

	// Delete
	t.Cleanup(func() {
		err := api.Delete(ctx, createdPeerID)
		assert.NoError(err)
	})
}
