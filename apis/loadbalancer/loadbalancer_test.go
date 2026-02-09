// Copyright 2026- The sacloud/apprun-dedicated-api-go authors
// SPDX-License-Identifier: Apache-2.0

package loadbalancer_test

import (
	"net/http"
	"testing"

	"github.com/go-faster/jx"
	"github.com/google/uuid"
	. "github.com/sacloud/apprun-dedicated-api-go/apis/loadbalancer"
	v1 "github.com/sacloud/apprun-dedicated-api-go/apis/v1"
	apprun_test "github.com/sacloud/apprun-dedicated-api-go/testutil"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

func setup(
	t *testing.T,
	v interface{ Encode(*jx.Encoder) },
	s ...int,
) (
	assert *require.Assertions,
	api LoadBalancerAPI,
) {
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
	actual, err := api.ListNode(t.Context(), lbID, 10, nil)

	assert.NoError(err)
	assert.NotNil(actual)
	assert.Equal(expected.GetLoadBalancerNodes(), actual)
}

func TestListNode_failed(t *testing.T) {
	expected := apprun_test.Fake403Error()
	assert, api := setup(t, &expected, int(expected.GetStatus()))

	lbID := v1.LoadBalancerID(uuid.New())
	actual, err := api.ListNode(t.Context(), lbID, 0, nil)

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
