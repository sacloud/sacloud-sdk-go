// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package cluster_test

import (
	"net/http"
	"os"
	"testing"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/cluster"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	super "github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v interface{ Encode(*jx.Encoder) }, s ...int) (assert *require.Assertions, api ClusterAPI) {
	assert = require.New(t)
	c, e := apprun_test.NewTestClient(v, s...)
	assert.NoError(e)
	api = NewClusterOp(c)

	return
}

func TestList(t *testing.T) {
	next := v1.ClusterID(uuid.New())
	var expected v1.ListClusterResponse
	expected.SetFake()
	expected.NextCursor.SetTo(next)
	expected.SetClusters(make([]v1.ReadClusterSummary, 3))
	for i := 0; i < len(expected.GetClusters()); i++ {
		var cluster v1.ReadClusterSummary
		cluster.SetFake()
		expected.Clusters[i] = cluster
	}
	assert, api := setup(t, &expected)

	actual, cursor, err := api.List(t.Context(), 10, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.NotNil(cursor)
	assert.Len(actual, 3)
	assert.Equal(expected.GetClusters()[0].GetClusterID(), actual[0].ClusterID)
}

func TestList_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, cursor, err := api.List(t.Context(), 0, nil)

	assert.Error(err)
	assert.Nil(actual)
	assert.Nil(cursor)
	assert.False(saclient.IsNotFoundError(err))
}

func TestCreate(t *testing.T) {
	var expected v1.CreateClusterResponse
	var cluster v1.CreatedCluster
	cluster.SetFake()
	expected.SetCluster(cluster)
	assert, api := setup(t, &expected)

	params := CreateParams{
		Name:               "test-cluster",
		LetsEncryptEmail:   nil,
		Ports:              []v1.CreateLoadBalancerPort{},
		ServicePrincipalID: "sp-123456789", // 12 chars
	}
	actual, err := api.Create(t.Context(), params)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(cluster.GetClusterID(), actual.ClusterID)
}

func TestCreate_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusConflict)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	params := CreateParams{
		Name:               "test-cluster",
		ServicePrincipalID: "sp-123456789", // 12 chars
	}
	actual, err := api.Create(t.Context(), params)

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}

func TestRead(t *testing.T) {
	var expected v1.GetClusterResponse
	var cluster v1.ReadClusterDetail
	cluster.SetFake()
	id := cluster.GetClusterID()
	expected.SetCluster(cluster)

	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), id)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(id, actual.ClusterID)
	assert.Equal(cluster.GetName(), actual.Name)
}

func TestRead_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.ClusterID(uuid.New())
	actual, err := api.Read(t.Context(), id)

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestUpdate(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	id := v1.ClusterID(uuid.New())
	params := UpdateParams{
		ServicePrincipalID: "sp-678901234", // 12 chars
	}
	err := api.Update(t.Context(), id, params)

	assert.NoError(err)
}

func TestUpdate_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.ClusterID(uuid.New())
	params := UpdateParams{
		ServicePrincipalID: "sp-678901234", // 12 chars
	}
	err := api.Update(t.Context(), id, params)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	id := v1.ClusterID(uuid.New())
	err := api.Delete(t.Context(), id)

	assert.NoError(err)
}

func TestDelete_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.ClusterID(uuid.New())
	err := api.Delete(t.Context(), id)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestIntegrated(t *testing.T) {
	assert, client := apprun_test.IntegratedClient(t)
	api := NewClusterOp(client)

	assert.NotNil(api)

	t.Run("Create", func(t *testing.T) {
		spid := os.Getenv("SAKURA_APPRUN_DEDICATED_SERVICE_PRINCIPAL_ID")
		name := super.RandomName("test-", 15, super.CharSetAlphaNum)
		params := CreateParams{
			Name:               name,
			ServicePrincipalID: spid,
			LetsEncryptEmail:   nil, // :TODO: to be tested
			Ports: []v1.CreateLoadBalancerPort{
				{
					Port:     443,
					Protocol: v1.CreateLoadBalancerPortProtocolHTTPS,
				},
			},
		}
		cluster, err := api.Create(t.Context(), params)
		assert.NoError(err)
		assert.NotNil(cluster)

		cid := cluster.ClusterID

		defer t.Run("Delete", func(t *testing.T) {
			err := api.Delete(t.Context(), cid)
			assert.NoError(err)
		})

		t.Run("List", func(t *testing.T) {
			listed := apprun_test.RepeatedList(func(cursor *v1.ClusterID) (clusters []ClusterDetail, next *v1.ClusterID) {
				clusters, next, err := api.List(t.Context(), 10, cursor)
				assert.NoError(err)
				return
			})

			assert.NotEmpty(listed)
		})

		t.Run("Read", func(t *testing.T) {
			actual, err := api.Read(t.Context(), cid)
			assert.NoError(err)
			assert.NotNil(actual)
			assert.Equal(cid, actual.ClusterID)
			assert.Equal(name, actual.Name)
			assert.Equal(spid, actual.ServicePrincipalID)
		})

		t.Run("Update", func(t *testing.T) {
			// :TODO: to be tested with different SPID
			err := api.Update(t.Context(), cid, UpdateParams{
				ServicePrincipalID: spid,
			})
			assert.NoError(err)
		})
	})
}
