// Copyright 2026- The sacloud/sacloud-sdk-go Authors
// SPDX-License-Identifier: Apache-2.0

package networkingsuite_test

import (
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	networkingsuite "github.com/sacloud/sacloud-sdk-go/api/networking-suite"
	v1 "github.com/sacloud/sacloud-sdk-go/api/networking-suite/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/packages/envvar"
	"github.com/sacloud/sacloud-sdk-go/common/packages/testutil"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
	"github.com/sacloud/sacloud-sdk-go/srn"
)

func TestSubnetGroupsAPI(t *testing.T) {
	// SAKURA_ENDPOINTS_NETWORKING_SUITE for testing
	testutil.PreCheckEnvsFunc("SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")(t)

	var theClient saclient.Client
	client, err := networkingsuite.NewClient(&theClient)
	require.NoError(t, err)

	ctx := t.Context()
	subnetGroupsOp := networkingsuite.NewSubnetGroupsOp(client)

	zone := getZone(t, os.Getenv("SAKURA_ENDPOINTS_NETWORKING_SUITE"))
	resCreate, err := subnetGroupsOp.Create(ctx, v1.CreateSubnetGroup{
		Name:                 "subnet group from go",
		Description:          "subnet group from go client",
		IPv4AddressRangeCIDR: "192.168.0.0/24",
		Region:               v1.Region{Code: zone[:len(zone)-1]},
	})
	require.NoError(t, err)
	assert.Equal(t, "subnet group from go", resCreate.Name)

	s, err := srn.Parse(resCreate.SRN)
	require.NoError(t, err)

	defer func() {
		err := subnetGroupsOp.Delete(ctx, s)
		if err != nil {
			t.Fatalf("unexpected error on subnet group delete: %v", err)
		}
	}()

	resList, err := subnetGroupsOp.List(ctx)
	assert.NoError(t, err)

	found := false
	for _, sg := range resList {
		if sg.SRN == resCreate.SRN {
			require.Equal(t, "subnet group from go client", sg.Description)
			found = true
		}
	}
	assert.True(t, found, "created subnet group not found in list")

	_, err = subnetGroupsOp.Update(ctx, s, v1.UpdateSubnetGroup{
		Name:        "subnet group from go 2",
		Description: "subnet group from go client 2",
	})
	assert.NoError(t, err)

	resRead, err := subnetGroupsOp.Read(ctx, s)
	assert.NoError(t, err)
	assert.Equal(t, "subnet group from go 2", resRead.Name)
	assert.Equal(t, "subnet group from go client 2", resRead.Description)
}

func getZone(t *testing.T, urlStr string) string {
	if urlStr == "" {
		return envvar.StringFromEnv("SAKURA_ZONE", "is1c")
	} else {
		u, err := url.ParseRequestURI(urlStr)
		if err != nil {
			t.Fatalf("failed to parse URL: %v", err)
		}
		parts := strings.Split(u.Path, "/")
		if len(parts) < 4 {
			t.Fatalf("unexpected URL path format: %s", u.Path)
		}
		return parts[3]
	}
}
