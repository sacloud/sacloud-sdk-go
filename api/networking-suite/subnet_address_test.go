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

func TestAddressesAPI(t *testing.T) {
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := networkingsuite.NewClient(&theClient)
	require.NoError(t, err)
	ctx := t.Context()
	zone := getZone(t, os.Getenv("SAKURA_ENDPOINTS_NETWORKING_SUITE"))

	subnetGroupsOp := networkingsuite.NewSubnetGroupsOp(client)
	subnetGroup, err := subnetGroupsOp.Create(ctx, v1.CreateSubnetGroup{
		Name:                 testutil.RandomName("address-test-", 8, testutil.CharSetAlphaNum),
		Description:          "networking-suite address test",
		IPv4AddressRangeCIDR: "192.168.10.0/24",
		Region:               v1.Region{Code: zone[:len(zone)-1]},
	})
	require.NoError(t, err)
	subnetGroupSRN, err := srn.Parse(subnetGroup.SRN)
	require.NoError(t, err)
	defer func() {
		if err := subnetGroupsOp.Delete(ctx, subnetGroupSRN); err != nil {
			t.Errorf("failed to delete test subnet group: %v", err)
		}
	}()

	subnetsOp := networkingsuite.NewSubnetsOp(client)
	subnet, err := subnetsOp.Create(ctx, v1.CreateSubnet{
		Name:                 testutil.RandomName("address-test-", 8, testutil.CharSetAlphaNum),
		Description:          "networking-suite address test",
		IPv4AddressRangeCIDR: "192.168.10.0/28",
		SubnetGroup:          v1.SakuraResourceNameRef{SRN: subnetGroup.SRN},
		Zone:                 v1.Zone{Code: zone},
	})
	require.NoError(t, err)
	subnetSRN, err := srn.Parse(subnet.SRN)
	require.NoError(t, err)
	defer func() {
		if err := subnetsOp.Delete(ctx, subnetSRN); err != nil {
			t.Errorf("failed to delete test subnet: %v", err)
		}
	}()

	addrOp := networkingsuite.NewAddressesOp(client)
	address, err := addrOp.Create(ctx, networkingsuite.CreateAddressParams{SubnetSRN: subnetSRN})
	require.NoError(t, err)
	assert.Equal(t, "SERVICE_ENDPOINT_GATEWAY_ADDRESS", address.AddressType)
	assert.Equal(t, "IPv4", address.IPVersion)
	assert.NotEmpty(t, address.IPAddress)
	assert.Equal(t, subnet.SRN, address.Subnet.SRN)

	addressSRN, err := srn.Parse(address.SRN)
	require.NoError(t, err)

	defer func() {
		if err := addrOp.Delete(ctx, addressSRN); err != nil {
			t.Errorf("failed to delete test address: %v", err)
		}
		require.NoError(t, err)
	}()

	read, err := addrOp.Read(ctx, addressSRN)
	require.NoError(t, err)
	assert.Equal(t, address.SRN, read.SRN)

	list, err := addrOp.List(ctx, networkingsuite.ListAddressesParams{SubnetSRN: subnetSRN})
	require.NoError(t, err)
	found := false
	for _, listedAddress := range list {
		if listedAddress.SRN == address.SRN {
			found = true
			break
		}
	}
	assert.True(t, found, "created address not found in list")
}
