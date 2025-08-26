// Copyright 2025- The sacloud/monitoring-suite-api-go Authors
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

package monitoringsuite_test

import (
	"context"
	"net/http"
	"testing"

	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestManagementOp_ResourceLimits(t *testing.T) {
	var res v1.ResourcesLimits
	res.SetFake()
	client := newTestClient(res)
	api := NewManagementOp(client)
	ctx := context.Background()

	limits, err := api.ResourceLimits(ctx)
	require.NoError(t, err)
	require.NotNil(t, limits)
	require.Equal(t, res.GetLogs(), limits.GetLogs())
	require.Equal(t, res.GetMetrics(), limits.GetMetrics())
	require.Equal(t, res.GetAlerts(), limits.GetAlerts())
	require.Equal(t, res.GetDashboards(), limits.GetDashboards())
}

func TestManagementOp_ResourceLimits_400(t *testing.T) {
	expected := ErrorResponse{
		Code:    "bad_request",
		Message: "insufficient privileges",
		IsOk:    false,
		Status:  400,
	}
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewManagementOp(client)
	ctx := context.Background()

	limits, err := api.ResourceLimits(ctx)
	require.Nil(t, limits)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient privileges")
}

func TestManagementOp_ProvisioningRead(t *testing.T) {
	var res v1.Provisioning
	res.SetFake()
	client := newTestClient(res)
	api := NewManagementOp(client)
	ctx := context.Background()

	prov, err := api.ProvisioningRead(ctx)
	require.NoError(t, err)
	require.NotNil(t, prov)
	require.Equal(t, res.GetLogs(), prov.GetLogs())
	require.Equal(t, res.GetMetrics(), prov.GetMetrics())
}

func TestManagementOp_ProvisioningRead_400(t *testing.T) {
	expected := ErrorResponse{
		Code:    "bad_request",
		Message: "insufficient privileges",
		IsOk:    false,
		Status:  400,
	}
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewManagementOp(client)
	ctx := context.Background()

	prov, err := api.ProvisioningRead(ctx)
	require.Nil(t, prov)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient privileges")
}

func TestManagementOp_ProvisioningCreate(t *testing.T) {
	var req v1.ProvisioningCreate
	var res v1.Provisioning
	res.SetFake()
	req.SetLogs(v1.NewOptProvisioningExist(res.GetLogs()))
	req.SetMetrics(v1.NewOptProvisioningExist(res.GetMetrics()))
	client := newTestClient(res, http.StatusCreated)
	api := NewManagementOp(client)
	ctx := context.Background()

	prov, err := api.ProvisioningCreate(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, prov)
	require.Equal(t, res.GetLogs(), prov.GetLogs())
	require.Equal(t, res.GetMetrics(), prov.GetMetrics())
}

func TestManagementOp_ProvisioningCreate_400(t *testing.T) {
	expected := ErrorResponse{
		Code:    "bad_request",
		Message: "insufficient privileges",
		IsOk:    false,
		Status:  400,
	}
	req := v1.ProvisioningCreate{}
	req.SetFake()
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewManagementOp(client)
	ctx := context.Background()

	prov, err := api.ProvisioningCreate(ctx, req)
	require.Nil(t, prov)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient privileges")
}
