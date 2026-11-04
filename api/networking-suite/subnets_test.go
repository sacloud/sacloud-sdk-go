// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/sacloud/sacloud-sdk-go/api/iaas"
	"github.com/sacloud/sacloud-sdk-go/api/iaas/types"
	networkingsuite "github.com/sacloud/sacloud-sdk-go/api/networking-suite"
	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/sacloud/sacloud-sdk-go/srn"
)

func TestSubnetAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := networkingsuite.NewClient(&theClient)
	require.NoError(t, err)

	ctx := t.Context()
	subnetGroupsOp := networkingsuite.NewSubnetGroupsOp(client)

	zone := getZone(t, os.Getenv("SAKURA_ENDPOINTS_NETWORKING_SUITE"))
	sg, err := subnetGroupsOp.Create(ctx, v1.CreateSubnetGroup{
		Name:                 "subnet group from go",
		Description:          "subnet group from go client",
		IPv4AddressRangeCIDR: "192.168.0.0/24",
		Region:               v1.Region{Code: zone[:len(zone)-1]},
	})
	require.NoError(t, err)
	assert.Equal(t, "subnet group from go", sg.Name)

	sgs, err := srn.Parse(sg.SRN)
	require.NoError(t, err)

	defer func() {
		err := subnetGroupsOp.Delete(ctx, sgs)
		if err != nil {
			t.Fatalf("unexpected error on delete: %v", err)
		}
	}()

	subnetOp := networkingsuite.NewSubnetsOp(client)

	subnet, err := subnetOp.Create(ctx, v1.CreateSubnet{
		Name:                 "subnet from go",
		Description:          "subnet from go client",
		IPv4AddressRangeCIDR: "192.168.0.0/28",
		SubnetGroup:          v1.SakuraResourceNameRef{SRN: sg.SRN},
		Zone:                 v1.Zone{Code: zone},
	})
	require.NoError(t, err)

	ss, err := srn.Parse(subnet.SRN)
	require.NoError(t, err)

	defer func() {
		err := subnetOp.Delete(ctx, ss)
		if err != nil {
			t.Fatalf("unexpected error on subnet delete: %v", err)
		}
	}()

	resList, err := subnetOp.List(ctx, sgs)
	assert.NoError(t, err)

	found := false
	for _, s := range resList {
		if s.SRN == subnet.SRN {
			require.Equal(t, "subnet from go client", s.Description)
			found = true
		}
	}
	assert.True(t, found, "created subnet not found in list")

	_, err = subnetOp.Update(ctx, ss, v1.UpdateSubnet{
		Name:        "subnet from go 2",
		Description: "subnet from go client 2",
	})
	assert.NoError(t, err)

	resRead, err := subnetOp.Read(ctx, ss)
	assert.NoError(t, err)
	assert.Equal(t, "subnet from go 2", resRead.Name)
	assert.Equal(t, "subnet from go client 2", resRead.Description)
}

func TestInterfaceConnectionAPI(t *testing.T) {
	if !testutil.IsAccTest() {
		t.Skip("set TESTACC=1 to run the acceptance test")
	}

	var theClient saclient.Client
	client, err := networkingsuite.NewClient(&theClient)
	require.NoError(t, err)
	endpointConfig, err := theClient.EndpointConfig()
	require.NoError(t, err)

	ctx := t.Context()
	networkingSuiteEndpoint := os.Getenv("SAKURA_ENDPOINTS_NETWORKING_SUITE")
	if networkingSuiteEndpoint == "" {
		networkingSuiteEndpoint = endpointConfig.Endpoints["networking_suite"]
	}
	zone := endpointConfig.Zone
	if networkingSuiteEndpoint != "" {
		zone = getZone(t, networkingSuiteEndpoint)
	}
	if zone == "" {
		zone = getZone(t, "")
	}

	iaasClient := iaas.NewClientFromSaclient(&theClient)
	serverOp := iaas.NewServerOp(iaasClient)
	server, err := serverOp.Create(ctx, zone, &iaas.ServerCreateRequest{
		CPU:        1,
		MemoryMB:   1024,
		Generation: types.PlanGenerations.G100,
		ConnectedSwitches: []*iaas.ConnectedSwitch{
			nil,
		},
		Name: testutil.RandomName("networking-suite-", 8, testutil.CharSetAlphaNum),
	})
	require.NoError(t, err)
	defer func() {
		if err := serverOp.Delete(ctx, zone, server.ID); err != nil {
			t.Errorf("failed to delete test server: %v", err)
		}
	}()

	server, err = serverOp.Read(ctx, zone, server.ID)
	require.NoError(t, err)
	require.NotEmpty(t, server.Interfaces)
	interfaceID := server.Interfaces[0].ID
	require.False(t, interfaceID.IsEmpty())
	_, err = iaas.NewInterfaceOp(iaasClient).Read(ctx, zone, interfaceID)
	require.NoError(t, err)

	interfaceSRN := srn.SRN{
		Version:  1,
		Location: "sakura-" + zone,
		Resource: "sakura.iaas.interface",
		ID:       interfaceID.String(),
	}
	require.NoError(t, interfaceSRN.Validate())

	subnetGroupsOp := networkingsuite.NewSubnetGroupsOp(client)
	sg, err := subnetGroupsOp.Create(ctx, v1.CreateSubnetGroup{
		Name:                 "for InterfaceConnection",
		Description:          "for InterfaceConnection test",
		IPv4AddressRangeCIDR: "192.168.0.0/24",
		Region:               v1.Region{Code: zone[:len(zone)-1]},
	})
	require.NoError(t, err)
	defer func() {
		err := subnetGroupsOp.Delete(ctx, srn.MustParse(sg.SRN))
		if err != nil {
			t.Fatalf("unexpected error on delete: %v", err)
		}
	}()

	subnetOp := networkingsuite.NewSubnetsOp(client)
	subnet, err := subnetOp.Create(ctx, v1.CreateSubnet{
		Name:                 "for InterfaceConnection",
		Description:          "for InterfaceConnection test",
		IPv4AddressRangeCIDR: "192.168.0.0/28",
		SubnetGroup:          v1.SakuraResourceNameRef{SRN: sg.SRN},
		Zone:                 v1.Zone{Code: zone},
	})
	require.NoError(t, err)
	defer func() {
		err := subnetOp.Delete(ctx, srn.MustParse(subnet.SRN))
		if err != nil {
			t.Fatalf("unexpected error on subnet delete: %v", err)
		}
	}()

	ip := "192.168.0.10"
	icOp := networkingsuite.NewInterfaceConnectionOp(client)
	ic, err := icOp.Create(ctx, networkingsuite.CreateInterfaceConnectionParams{
		InterfaceSRN: interfaceSRN,
		SubnetSRN:    srn.MustParse(subnet.SRN),
		IPAddress:    &ip,
	})
	require.NoError(t, err)
	defer func() {
		err = icOp.Delete(ctx, srn.MustParse(ic.SRN))
		if err != nil {
			t.Fatalf("unexpected error on interface connection delete: %v", err)
		}
	}()

	connectionSRN, err := srn.Parse(ic.SRN)
	require.NoError(t, err)

	server, err = serverOp.Read(ctx, zone, server.ID)
	require.NoError(t, err)
	var connectedInterface *iaas.InterfaceView
	for _, iface := range server.Interfaces {
		if iface.ID == interfaceID {
			connectedInterface = iface
			break
		}
	}
	require.NotNil(t, connectedInterface, "created NIC should be present in the server response")
	require.Equal(t, types.StringID(connectionSRN.ID), connectedInterface.NetworkingSuiteInterfaceConnectionID)

	aOp := networkingsuite.NewAddressesOp(client)
	listRes, err := aOp.List(ctx, networkingsuite.ListAddressesParams{
		SubnetSRN: srn.MustParse(subnet.SRN),
	})
	require.NoError(t, err)
	require.Len(t, listRes, 1)

	addr := listRes[0]
	assert.True(t, srn.IsSRN(addr.SRN))
	assert.Equal(t, "EPHEMERAL_ADDRESS", addr.AddressType)
	assert.Equal(t, "IPv4", addr.IPVersion)
	assert.Equal(t, ip, addr.IPAddress)
	assert.Equal(t, subnet.SRN, addr.Subnet.SRN)
}
