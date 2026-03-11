// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package loadbalancer_test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/loadbalancer"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	super "github.com/sacloud/packages-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, v interface{ Encode(*jx.Encoder) }, s ...int) (assert *require.Assertions, api LoadBalancerAPI) {
	cid := v1.ClusterID(uuid.New())
	asgid := v1.AutoScalingGroupID(uuid.New())
	assert = require.New(t)
	c, e := apprun_test.NewTestClient(v, s...)
	assert.NoError(e)
	api = NewLoadBalancerOp(c, cid, asgid)

	return
}

func TestList(t *testing.T) {
	next := v1.LoadBalancerID(uuid.New())
	var expected v1.ListLoadBalancersResponse
	expected.SetFake()
	expected.NextCursor.SetTo(next)
	expected.SetLoadBalancers(make([]v1.ReadLoadBalancerSummary, 3))
	for i := 0; i < len(expected.GetLoadBalancers()); i++ {
		expected.LoadBalancers[i] = apprun_test.FakeLoadBalancer()
	}
	assert, api := setup(t, &expected)

	actual, cursor, err := api.List(t.Context(), 10, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.NotNil(cursor)
	assert.Equal(expected.GetLoadBalancers(), actual)
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
	var expected v1.CreateLoadBalancerResponse
	expected.SetFake()
	assert, api := setup(t, &expected)

	actual, err := api.Create(t.Context(), CreateParams{
		Name:             "test-lb",
		ServiceClassPath: "/is1a/load-balancer/1/core-2-2",
		NameServers:      []v1.IPv4{"133.242.0.3"},
		Interfaces: []LoadBalancerInterface{{
			InterfaceIndex:  0,
			Upstream:        "shared",
			IpPool:          []v1.IpRange{},
			NetmaskLen:      saclient.Ptr(int16(24)),
			DefaultGateway:  saclient.Ptr("192.168.1.1"),
			Vip:             saclient.Ptr("203.0.113.1"),
			VirtualRouterID: saclient.Ptr(int16(1)),
			PacketFilterID:  saclient.Ptr("filter-id"),
		}},
	})

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetLoadBalancer(), *actual)
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
	var expected v1.GetLoadBalancerResponse
	fake := apprun_test.FakeLoadBalancer()
	id := fake.LoadBalancerID
	detail := apprun_test.FakeLoadBalancerDetail()
	detail.SetLoadBalancerID(id)
	expected.SetLoadBalancer(detail)

	assert, api := setup(t, &expected)

	actual, err := api.Read(t.Context(), id)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(id, actual.LoadBalancerID)
}

func TestRead_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.LoadBalancerID(uuid.New())
	actual, err := api.Read(t.Context(), id)

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestDelete(t *testing.T) {
	assert, api := setup(t, nil, http.StatusNoContent)

	id := v1.LoadBalancerID(uuid.New())
	err := api.Delete(t.Context(), id)

	assert.NoError(err)
}

func TestDelete_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	id := v1.LoadBalancerID(uuid.New())
	err := api.Delete(t.Context(), id)

	assert.Error(err)
	assert.True(saclient.IsNotFoundError(err))
}

func TestListNode(t *testing.T) {
	var expected v1.ListLoadBalancerNodesResponse
	expected.SetFake()
	expected.SetLoadBalancerNodes(make([]v1.ReadLoadBalancerNodeSummary, 3))
	for i := 0; i < len(expected.GetLoadBalancerNodes()); i++ {
		expected.LoadBalancerNodes[i] = apprun_test.FakeLoadBalancerNodeSummary()
	}
	assert, api := setup(t, &expected)

	lbID := v1.LoadBalancerID(uuid.New())
	actual, err := api.ListNodes(t.Context(), lbID, 10, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetLoadBalancerNodes(), actual)
}

func TestListNode_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	lbID := v1.LoadBalancerID(uuid.New())
	actual, err := api.ListNodes(t.Context(), lbID, 0, nil)

	assert.Error(err)
	assert.Nil(actual)
	assert.False(saclient.IsNotFoundError(err))
}

func TestReadNode(t *testing.T) {
	var expected v1.GetLoadBalancerNodeResponse
	fake := apprun_test.FakeLoadBalancerNode()
	nodeID := fake.LoadBalancerNodeID
	expected.SetLoadBalancerNode(fake)

	assert, api := setup(t, &expected)

	lbID := v1.LoadBalancerID(uuid.New())
	actual, err := api.ReadNode(t.Context(), lbID, nodeID)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(nodeID, actual.LoadBalancerNodeID)
}

func TestReadNode_failed(t *testing.T) {
	var expected v1.Error
	expected.SetFake()
	expected.SetStatus(http.StatusNotFound)
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	lbID := v1.LoadBalancerID(uuid.New())
	nodeID := v1.LoadBalancerNodeID(uuid.New())
	actual, err := api.ReadNode(t.Context(), lbID, nodeID)

	assert.Error(err)
	assert.Nil(actual)
	assert.True(saclient.IsNotFoundError(err))
}

func TestIntegrated(t *testing.T) {
	assert, client := apprun_test.IntegratedClient(t)
	cid, d1 := apprun_test.IntegratedCluster(t.Context(), assert, client)
	defer d1()

	aid, d2 := apprun_test.IntegratedAsg(t.Context(), assert, client, cid)
	defer d2()

	api := NewLoadBalancerOp(client, cid, aid)
	assert.NotNil(api)

	t.Run("Create", func(t *testing.T) {
		lbName := super.RandomName("test-", 15, super.CharSetAlphaNum)
		lb, err := api.Create(t.Context(), CreateParams{
			Name:             lbName,
			ServiceClassPath: "cloud/apprun/dedicated/lb/1vcpu_2gb_1", // :FIXME: there is no way to find a minimal class
			NameServers:      []v1.IPv4{"133.242.0.3"},
			Interfaces: []LoadBalancerInterface{{
				InterfaceIndex:  0,
				Upstream:        "shared",
				IpPool:          []v1.IpRange{},
				NetmaskLen:      nil,
				DefaultGateway:  nil,
				Vip:             nil,
				VirtualRouterID: nil,
				PacketFilterID:  nil,
			}},
		})
		assert.NoError(err)
		assert.NotNil(lb)

		lbID := lb.LoadBalancerID

		defer t.Run("Delete", func(t *testing.T) {
			tkr := time.NewTicker(16 * time.Second)
			defer tkr.Stop()

			nodes, err := api.ListNodes(t.Context(), lbID, 10, nil)
			assert.NoError(err)

			// a load balancer node needs to be deleted _after_ it finished provisioned
			// need to wait
		loop1:
			for {
				select {
				case <-t.Context().Done():
					assert.Fail("timeout while waiting for load balancer to be provisioned")
					break loop1
				case <-tkr.C:
					var provisioning []v1.ReadLoadBalancerNodeSummary
					for _, node := range nodes {
						current, err := api.ReadNode(t.Context(), lbID, node.LoadBalancerNodeID)
						switch {
						case saclient.IsNotFoundError(err):
							continue
						case err != nil:
							assert.NoError(err)
						default:
							switch current.Status {
							case v1.LoadBalancerNodeStatusCreating, v1.LoadBalancerNodeStatusStarting:
								fmt.Printf("waiting for LB node provisioning: %v\n", uuid.UUID(current.LoadBalancerNodeID).String())
								provisioning = append(provisioning, node)
							case v1.LoadBalancerNodeStatusHealthy:
								break loop1
							default:
								assert.Failf("unexpected load balancer node status", "status=%s", current.Status)
							}
						}
					}
				}
			}

			err = api.Delete(t.Context(), lbID)
			assert.NoError(err)

		loop2:
			for {
				select {
				case <-t.Context().Done():
					assert.Fail("timeout while waiting for load balancer deletion")
					break loop2
				case <-tkr.C:
					actual, err := api.Read(t.Context(), lbID)

					switch {
					case saclient.IsNotFoundError(err):
						break loop2
					case err != nil:
						assert.NoError(err)
					default:
						assert.NotNil(actual)
						fmt.Printf("waiting for LB deletion: %v\n", uuid.UUID(actual.LoadBalancerID).String())
					}
				}
			}
		})

		t.Run("List", func(t *testing.T) {
			list := apprun_test.RepeatedList(func(cursor *v1.LoadBalancerID) (res []v1.ReadLoadBalancerSummary, next *v1.LoadBalancerID) {
				res, next, err := api.List(t.Context(), 10, cursor)
				assert.NoError(err)
				return
			})
			assert.NotEmpty(list)
		})

		t.Run("Read", func(t *testing.T) {
			actual, err := api.Read(t.Context(), lbID)
			assert.NoError(err)
			assert.NotNil(actual)
			assert.Equal(lbID, actual.LoadBalancerID)
			assert.Equal(lbName, actual.Name)
		})

		t.Run("ListNode", func(t *testing.T) {
			nodes := apprun_test.RepeatedList(func(cursor *v1.LoadBalancerID) (res []v1.ReadLoadBalancerNodeSummary, next *v1.LoadBalancerID) {
				res, err := api.ListNodes(t.Context(), lbID, 10, cursor)
				assert.NoError(err)
				return
			})
			assert.NotEmpty(nodes)
		})

		t.Run("ReadNode", func(t *testing.T) {
			// First list nodes to get a node ID
			nodes := apprun_test.RepeatedList(func(cursor *v1.LoadBalancerID) (res []v1.ReadLoadBalancerNodeSummary, next *v1.LoadBalancerID) {
				res, err := api.ListNodes(t.Context(), lbID, 10, cursor)
				assert.NoError(err)
				return
			})
			assert.NotEmpty(nodes)

			// Read the first node
			nodeID := nodes[0].LoadBalancerNodeID
			actual, err := api.ReadNode(t.Context(), lbID, nodeID)
			assert.NoError(err)
			assert.NotNil(actual)
			assert.Equal(nodeID, actual.LoadBalancerNodeID)
		})
	})
}
