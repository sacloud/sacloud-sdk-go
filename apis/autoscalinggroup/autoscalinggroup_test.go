// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package autoscalinggroup_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/autoscalinggroup"
	"github.com/sacloud/apprun-dedicated-api-go/apis/service_class"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	"github.com/sacloud/apprun-dedicated-api-go/apis/workernode"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	super "github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v interface{ Encode(*jx.Encoder) }, s ...int) (assert *require.Assertions, api AutoScalingGroupAPI) {
	cid := v1.ClusterID(uuid.New())
	assert = require.New(t)
	c, e := apprun_test.NewTestClient(v, s...)
	assert.NoError(e)
	api = NewAutoScalingGroupOp(c, cid)

	return
}

func TestList(t *testing.T) {
	next := v1.AutoScalingGroupID(uuid.New())
	var expected v1.ListAutoScalingGroupResponse
	expected.SetFake()
	expected.NextCursor.SetTo(next)
	expected.SetAutoScalingGroups(make([]v1.ReadAutoScalingGroupDetail, 3))
	for i := 0; i < len(expected.GetAutoScalingGroups()); i++ {
		expected.AutoScalingGroups[i] = apprun_test.FakeAutoScalingGroup()
	}
	assert, api := setup(t, &expected)

	actual, cursor, err := api.List(t.Context(), 10, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.NotNil(cursor)
	assert.Equal(expected.GetAutoScalingGroups(), actual)
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
	var expected v1.CreateAutoScalingGroupResponse
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Create(t.Context(), CreateParams{
		Name:                   "test-asg",
		Zone:                   "is1a",
		NameServers:            []v1.IPv4{"133.242.0.3"},
		WorkerServiceClassPath: "/is1a/server/1/core-2-2",
		MinNodes:               1,
		MaxNodes:               3,
		Interfaces: []NodeInterface{{
			InterfaceIndex: 0,
			Upstream:       "shared",
			IpPool:         []v1.IpRange{},
			NetmaskLen:     saclient.Ptr(int16(24)),
			DefaultGateway: saclient.Ptr("192.168.1.1"),
			PacketFilterID: saclient.Ptr("filter-id"),
			ConnectsToLB:   true,
		}},
	})

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetAutoScalingGroup(), *actual)
}

func TestCreate_failed(t *testing.T) {
	expected := apprun_test.Fake400Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	actual, err := api.Create(t.Context(), CreateParams{})

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}

func TestRead(t *testing.T) {
	var expected v1.GetAutoScalingGroupResponse
	fake := apprun_test.FakeAutoScalingGroup()
	id := fake.AutoScalingGroupID
	expected.SetAutoScalingGroup(fake)

	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), id)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(id, actual.AutoScalingGroupID)
}

func TestRead_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.AutoScalingGroupID(uuid.New())
	actual, err := api.Read(t.Context(), id)

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	id := v1.AutoScalingGroupID(uuid.New())
	err := api.Delete(t.Context(), id)

	assert.NoError(err)
}

func TestDelete_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.AutoScalingGroupID(uuid.New())
	err := api.Delete(t.Context(), id)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestIntegrated(t *testing.T) {
	assert, client := apprun_test.IntegratedClient(t)
	cid, deleter := apprun_test.IntegratedCluster(t.Context(), assert, client)
	defer deleter()

	api := NewAutoScalingGroupOp(client, cid)
	assert.NotNil(api)

	svc := service_class.NewServiceClassOp(client)
	assert.NotNil(svc)

	t.Run("Create", func(t *testing.T) {
		asgName := super.RandomName("test-", 15, super.CharSetAlphaNum)
		asg, err := api.Create(t.Context(), CreateParams{
			Name:                   asgName,
			Zone:                   "is1a",
			NameServers:            []v1.IPv4{"133.242.0.3"},
			WorkerServiceClassPath: "cloud/apprun/dedicated/worker/1vcpu_2gb", // :FIXME: there is no way to find a minimal class
			MinNodes:               1,
			MaxNodes:               1,
			Interfaces: []NodeInterface{{
				InterfaceIndex: 0,
				Upstream:       "shared",
				IpPool:         []v1.IpRange{},
				NetmaskLen:     nil,
				DefaultGateway: nil,
				PacketFilterID: nil,
				ConnectsToLB:   true,
			}},
		})
		assert.NoError(err)
		assert.NotNil(asg)

		asgID := asg.AutoScalingGroupID

		defer t.Run("Delete", func(t *testing.T) {
			// worker node deletion may take some time
			nodeop := workernode.NewWorkerNodeOp(client, cid, asgID)
			tkr := time.NewTicker(16 * time.Second)
			defer tkr.Stop()

			nodes := apprun_test.RepeatedList(func(cursor *v1.WorkerNodeID) (res []workernode.WorkerNodeDetail, next *v1.WorkerNodeID) {
				res, next, err := nodeop.List(t.Context(), 10, cursor)
				assert.NoError(err)
				return
			})
			if len(nodes) != 0 {
				for _, node := range nodes {
					err := nodeop.Update(t.Context(), node.WorkerNodeID, true)
					assert.NoError(err)
					fmt.Printf("requested: to drain %v\n", uuid.UUID(node.WorkerNodeID).String())
				}
			loop1:
				for {
					var alive []workernode.WorkerNodeDetail
					for _, node := range nodes {
						node, err := nodeop.Read(t.Context(), node.WorkerNodeID)
						assert.NoError(err)
						if node.Creating {
							fmt.Printf("waiting for creation: %v\n", uuid.UUID(node.WorkerNodeID).String())
							alive = append(alive, *node)
						}
					}
					if len(alive) == 0 {
						break loop1
					} else {
						select {
						case <-t.Context().Done():
							assert.Fail("timeout while waiting for worker nodes to be drained")
							break loop1
						case <-tkr.C:
							continue loop1
						}
					}
				}
			}
			err := api.Delete(t.Context(), asgID)
			assert.NoError(err)

		loop2:
			for {
				select {
				case <-t.Context().Done():
					assert.Fail("timeout while waiting for worker nodes to be drained")
					break loop2
				case <-tkr.C:
					actual, err := api.Read(t.Context(), asgID)

					switch {
					case saclient.IsNotFoundError(err):
						break loop2
					case err != nil:
						assert.NoError(err)
					default:
						assert.NotNil(actual)
						fmt.Printf("waiting for deletion: %v\n", uuid.UUID(actual.AutoScalingGroupID).String())
					}
				}
			}
		})

		t.Run("List", func(t *testing.T) {
			list := apprun_test.RepeatedList(func(cursor *v1.AutoScalingGroupID) (res []v1.ReadAutoScalingGroupDetail, next *v1.AutoScalingGroupID) {
				res, next, err := api.List(t.Context(), 10, cursor)
				assert.NoError(err)
				return
			})
			assert.NotEmpty(list)
		})

		t.Run("Read", func(t *testing.T) {
			actual, err := api.Read(t.Context(), asgID)
			assert.NoError(err)
			assert.NotNil(actual)
			assert.Equal(asgID, actual.AutoScalingGroupID)
			assert.Equal(asgName, actual.Name)
		})
	})
}
