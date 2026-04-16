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
	ctx := t.Context()

	limits, err := api.ResourceLimits(ctx)
	require.NoError(t, err)
	require.NotNil(t, limits)
	require.Equal(t, res.GetLogs(), limits.GetLogs())
	require.Equal(t, res.GetMetrics(), limits.GetMetrics())
	require.Equal(t, res.GetAlerts(), limits.GetAlerts())
	require.Equal(t, res.GetDashboards(), limits.GetDashboards())
}

func TestManagementOp_ResourceLimits_400(t *testing.T) {
	expected := newErrorResponse(400, "insufficient privileges")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewManagementOp(client)
	ctx := t.Context()

	limits, err := api.ResourceLimits(ctx)
	require.Nil(t, limits)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient privileges")
}

func TestManagementOp_ReadProvisioning(t *testing.T) {
	var res v1.Provisioning
	res.SetFake()
	client := newTestClient(res)
	api := NewManagementOp(client)
	ctx := t.Context()

	prov, err := api.ReadProvisioning(ctx)
	require.NoError(t, err)
	require.NotNil(t, prov)
	require.Equal(t, res.GetLogs(), prov.GetLogs())
	require.Equal(t, res.GetMetrics(), prov.GetMetrics())
}

func TestManagementOp_ReadProvisioning_400(t *testing.T) {
	expected := newErrorResponse(400, "insufficient privileges")
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewManagementOp(client)
	ctx := t.Context()

	prov, err := api.ReadProvisioning(ctx)
	require.Nil(t, prov)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient privileges")
}

func TestManagementOp_CreateProvisioning(t *testing.T) {
	var req ProvisioningCreateParam
	var res v1.Provisioning
	res.SetFake()
	req.Logs = ref(res.GetLogs())
	req.Metrics = ref(res.GetMetrics())
	client := newTestClient(res, http.StatusCreated)
	api := NewManagementOp(client)
	ctx := t.Context()

	prov, err := api.CreateProvisioning(ctx, req)
	require.NoError(t, err)
	require.NotNil(t, prov)
	require.Equal(t, res.GetLogs(), prov.GetLogs())
	require.Equal(t, res.GetMetrics(), prov.GetMetrics())
}

func TestManagementOp_CreateProvisioning_400(t *testing.T) {
	expected := newErrorResponse(400, "insufficient privileges")
	req := ProvisioningCreateParam{}
	client := newTestClient(expected, http.StatusBadRequest)
	api := NewManagementOp(client)
	ctx := t.Context()

	prov, err := api.CreateProvisioning(ctx, req)
	require.Nil(t, prov)
	require.Error(t, err)
	require.ErrorContains(t, err, "insufficient privileges")
}

func TestManagementIntegrated(t *testing.T) {
	client, err := IntegratedClient(t)
	require.NoError(t, err)
	api := NewManagementOp(client)
	ctx := t.Context()

	// ResourceLimits
	limits, err := api.ResourceLimits(ctx)
	require.NoError(t, err)
	require.NotNil(t, limits)

	// CreateProvisioning
	created, err := api.CreateProvisioning(ctx, ProvisioningCreateParam{})
	require.NoError(t, err)
	require.NotNil(t, created)

	// ReadProvisioning
	prov, err := api.ReadProvisioning(ctx)
	require.NoError(t, err)
	require.NotNil(t, prov)
}
