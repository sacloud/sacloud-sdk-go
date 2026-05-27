// Copyright 2022-2025 The sacloud/iaas-api-go Authors
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

package plans

import (
	"context"
	"testing"

	"github.com/sacloud/iaas-api-go"
	"github.com/sacloud/iaas-api-go/testutil"
	"github.com/sacloud/iaas-api-go/types"
	"github.com/stretchr/testify/require"
)

func TestChangeProxyLBPlan_PreservesAllSettings(t *testing.T) {
	ctx := context.Background()
	caller := testutil.SingletonAPICaller()

	elbOp := iaas.NewProxyLBOp(caller)

	// ProxyLBを各種設定値を含めて作成
	createReq := &iaas.ProxyLBCreateRequest{
		Plan:   types.ProxyLBPlans.CPS100,
		Region: types.ProxyLBRegions.IS1,
		Name:   testutil.ResourceName("proxy-lb-plan-change-test"),
		HealthCheck: &iaas.ProxyLBHealthCheck{
			Protocol:  types.ProxyLBProtocols.TCP,
			DelayLoop: 10,
		},
		Timeout: &iaas.ProxyLBTimeout{
			InactiveSec: 10,
		},
		BackendHttpKeepAlive: &iaas.ProxyLBBackendHttpKeepAlive{
			Mode: types.ProxyLBBackendHttpKeepAlive.Aggressive,
		},
		MonitoringSuiteLog: &iaas.MonitoringSuiteLog{
			Enabled: true,
		},
		OriginGuard: &iaas.ProxyLBOriginGuard{
			Token: "test-token",
		},
		StrictRule: &iaas.ProxyLBStrictRule{
			Enabled: true,
		},
	}

	elb, err := elbOp.Create(ctx, createReq)
	require.NoError(t, err)
	require.NotNil(t, elb)

	// リソースのクリーンアップを確実に実行するため、t.Cleanupで登録
	// プラン変更後IDが変化する可能性があるため、変数で追跡
	cleanupID := elb.ID
	t.Cleanup(func() {
		_ = elbOp.Delete(ctx, cleanupID)
	})

	// 各設定が初期状態で保持されていることを確認
	require.NotNil(t, elb.BackendHttpKeepAlive)
	require.Equal(t, types.ProxyLBBackendHttpKeepAlive.Aggressive, elb.BackendHttpKeepAlive.Mode)
	require.NotNil(t, elb.MonitoringSuiteLog)
	require.True(t, elb.MonitoringSuiteLog.Enabled)
	require.NotNil(t, elb.OriginGuard)
	require.Equal(t, "test-token", elb.OriginGuard.Token)
	require.NotNil(t, elb.StrictRule)
	require.True(t, elb.StrictRule.Enabled)

	// プラン変更実行
	changed, err := ChangeProxyLBPlan(ctx, caller, elb.ID, types.ProxyLBPlans.CPS500.Int())
	require.NoError(t, err)
	require.NotNil(t, changed)

	// プラン変更後にIDが変化する可能性があるため、クリーンアップ対象IDを更新
	cleanupID = changed.ID

	// プランが変更されていることを確認
	require.Equal(t, types.ProxyLBPlans.CPS500, changed.Plan)

	// 各設定値が保持されていることを確認
	require.NotNil(t, changed.BackendHttpKeepAlive)
	require.Equal(t, types.ProxyLBBackendHttpKeepAlive.Aggressive, changed.BackendHttpKeepAlive.Mode)
	require.NotNil(t, changed.MonitoringSuiteLog)
	require.True(t, changed.MonitoringSuiteLog.Enabled)
	require.NotNil(t, changed.OriginGuard)
	require.Equal(t, "test-token", changed.OriginGuard.Token)
	require.NotNil(t, changed.StrictRule)
	require.True(t, changed.StrictRule.Enabled)
}
