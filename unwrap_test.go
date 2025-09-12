// Copyright 2022-2025 The sacloud/monitoring-suite-api-go Authors
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
	"testing"

	. "github.com/sacloud/monitoring-suite-api-go"
	v1 "github.com/sacloud/monitoring-suite-api-go/apis/v1"
	"github.com/stretchr/testify/require"
)

func TestUnwrap_WrappedDashboardProject(t *testing.T) {
	var dp v1.DashboardProject

	_, err := Unwrap(&dp, &TemplateDashboardProject)
	require.NoError(t, err)
	require.Equal(t, TemplateDashboardProject.ID, dp.ID)
	require.Equal(t, TemplateDashboardProject.Name, dp.Name)
	require.Equal(t, TemplateDashboardProject.Description, dp.Description)
	require.Equal(t, TemplateDashboardProject.Tags, dp.Tags)
}

func TestUnwrap_WrappedLogRouting(t *testing.T) {
	var lr v1.LogRouting

	_, err := Unwrap(&lr, &TemplateWrappedLogRouting)
	require.NoError(t, err)
	require.Equal(t, TemplateLogRouting.ID, lr.ID)
	require.Equal(t, TemplateLogRouting.ResourceID, lr.ResourceID)
	require.Equal(t, TemplateLogRouting.Publisher.Code, lr.Publisher.Code)
}

func TestUnwrap_WrappedLogStorage(t *testing.T) {
	var lt v1.LogStorage

	_, err := Unwrap(&lt, &TemplateWrappedLogStorage)
	require.NoError(t, err)
	require.Equal(t, TemplateLogStorage.ID, lt.ID)
	require.Equal(t, TemplateLogStorage.Name, lt.Name)
	require.Equal(t, TemplateLogStorage.Description, lt.Description)
	require.Equal(t, TemplateLogStorage.Tags, lt.Tags)
}

func TestUnwrap_WrappedMetricsRouting(t *testing.T) {
	var mr v1.MetricsRouting

	_, err := Unwrap(&mr, &TemplateWrappedMetricsRouting)
	require.NoError(t, err)
	require.Equal(t, TemplateMetricsRouting.ID, mr.ID)
	require.Equal(t, TemplateMetricsRouting.ResourceID, mr.ResourceID)
	require.Equal(t, TemplateMetricsRouting.Publisher.Code, mr.Publisher.Code)
}

func TestUnwrap_WrappedMetricsStorage(t *testing.T) {
	var mt v1.MetricsStorage

	_, err := Unwrap(&mt, &TemplateWrappedMetricsStorage)
	require.NoError(t, err)
	require.Equal(t, TemplateMetricsStorage.ID, mt.ID)
	require.Equal(t, TemplateMetricsStorage.Name, mt.Name)
	require.Equal(t, TemplateMetricsStorage.Description, mt.Description)
	require.Equal(t, TemplateMetricsStorage.Tags, mt.Tags)
}
