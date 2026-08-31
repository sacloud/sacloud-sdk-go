// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

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
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET", "SAKURA_NETWORKING_SUITE_INTERFACE_SRN")(t)

	var theClient saclient.Client
	client, err := networkingsuite.NewClient(&theClient)
	require.NoError(t, err)

	ctx := t.Context()
	subnetGroupsOp := networkingsuite.NewSubnetGroupsOp(client)
	zone := getZone(t, os.Getenv("SAKURA_ENDPOINTS_NETWORKING_SUITE"))
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
		InterfaceSRN: srn.MustParse(os.Getenv("SAKURA_NETWORKING_SUITE_INTERFACE_SRN")),
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
