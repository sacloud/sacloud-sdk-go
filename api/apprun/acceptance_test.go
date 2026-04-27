// Copyright 2021-2026 The sacloud/apprun-api-go authors
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

//go:build acctest
// +build acctest

package apprun_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
	"github.com/stretchr/testify/require"
)

const appName = "app-for-acceptance"

// TestUserAPI ユーザー関連APIの操作テスト
// ユーザーは削除もできないため、2回目以降は既にユーザーが存在する状態でのテストとなってしまうことに注意する。
func TestUserAPI(t *testing.T) {
	skipIfNoAPIKey(t)

	ctx := context.Background()
	client, err := newAPIClient()
	require.NoError(t, err)
	userOp := apprun.NewUserOp(client)

	// Create
	_, err = userOp.Create(ctx)
	if err != nil {
		if saclient.IsConflictError(err) {
			t.Log("user already exists, ignoring conflict error and continuing test")
		} else {
			t.Fatal(err)
		}
	}

	// Read
	res, err := userOp.Read(ctx)
	require.NoError(t, err)
	require.GreaterOrEqual(t, res.Limit.ApplicationCount, 0)
}

// TestApplicationAPI アプリケーションの一連の操作テスト
// 以下のシナリオでテストを行う
//   - アプリケーションを作成
//   - アプリケーションの一覧を取得
//   - アプリケーションを更新
//   - アプリケーションが更新できたかどうかを確認
//   - アプリケーションのステータスを取得
//   - アプリケーションを削除
func TestApplicationAPI(t *testing.T) {
	skipIfNoAPIKey(t)

	if err := cleanupTestApplication(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	client, err := newAPIClient()
	require.NoError(t, err)
	appOp := apprun.NewApplicationOp(client)

	// Create
	created, err := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:                   appName,
		TimeoutSeconds:         100,
		Port:                   80,
		MinScale:               0,
		MaxScale:               1,
		ScaleTargetConcurrency: v1.NewOptInt(100),
		Components: []v1.PostApplicationBodyComponentsItem{
			{
				Name:      "component1",
				MaxCPU:    v1.PostApplicationBodyComponentsItemMaxCPU05,
				MaxMemory: v1.PostApplicationBodyComponentsItemMaxMemory1Gi,
				DeploySource: v1.PostApplicationBodyComponentsItemDeploySource{
					ContainerRegistry: v1.NewOptPostApplicationBodyComponentsItemDeploySourceContainerRegistry(
						v1.PostApplicationBodyComponentsItemDeploySourceContainerRegistry{
							Image: "sakura-oss-dev.sakuracr.jp/test:latest",
						},
					),
				},
				Probe: v1.NewOptNilPostApplicationBodyComponentsItemProbe(
					v1.PostApplicationBodyComponentsItemProbe{
						HTTPGet: v1.NewOptNilPostApplicationBodyComponentsItemProbeHTTPGet(
							v1.PostApplicationBodyComponentsItemProbeHTTPGet{
								Path: "/",
								Port: 80,
							},
						),
					},
				),
			},
		},
	})
	require.NoError(t, err)

	// Read
	application, err := appOp.Read(ctx, created.ID)
	require.NoError(t, err)
	require.Equal(t, application.Name, appName)

	// Update
	timeoutUpdated := 20
	appOp.Update(ctx, application.ID, &v1.PatchApplicationBody{
		TimeoutSeconds: v1.NewOptInt(timeoutUpdated),
	})

	// Read
	application, err = appOp.Read(ctx, application.ID)
	require.NoError(t, err)
	require.Equal(t, application.TimeoutSeconds, timeoutUpdated)

	// Read Status
	// ヘルスチェックが完了するまでタイムラグがあるため暫く待つ
	time.Sleep(30 * time.Second)

	res, err := appOp.ReadStatus(ctx, application.ID)
	require.NoError(t, err)
	require.Equal(t, res.Status, v1.HandlerGetApplicationOnlyStatusStatusHealthy)

	// Delete
	err = appOp.Delete(ctx, application.ID)
	require.NoError(t, err)
}

// TestPacketFilterAPI アプリケーションのパケットフィルタの一連の操作テスト
// 以下のシナリオでテストを行う
//   - アプリケーションを作成
//   - パケットフィルタの作成
//   - パケットフィルタの取得
//   - アプリケーションを削除
func TestPacketFilterAPI(t *testing.T) {
	skipIfNoAPIKey(t)

	if err := cleanupTestApplication(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	client, err := newAPIClient()
	require.NoError(t, err)
	appOp := apprun.NewApplicationOp(client)
	pfOp := apprun.NewPacketFilterOp(client)

	// Application Create
	application, _ := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:                   appName,
		TimeoutSeconds:         100,
		Port:                   80,
		MinScale:               0,
		MaxScale:               1,
		ScaleTargetConcurrency: v1.NewOptInt(100),
		Components: []v1.PostApplicationBodyComponentsItem{
			{
				Name:      "component1",
				MaxCPU:    v1.PostApplicationBodyComponentsItemMaxCPU05,
				MaxMemory: v1.PostApplicationBodyComponentsItemMaxMemory1Gi,
				DeploySource: v1.PostApplicationBodyComponentsItemDeploySource{
					ContainerRegistry: v1.NewOptPostApplicationBodyComponentsItemDeploySourceContainerRegistry(
						v1.PostApplicationBodyComponentsItemDeploySourceContainerRegistry{
							Image: "sakura-oss-dev.sakuracr.jp/test:latest",
						},
					),
				},
				Probe: v1.NewOptNilPostApplicationBodyComponentsItemProbe(
					v1.PostApplicationBodyComponentsItemProbe{
						HTTPGet: v1.NewOptNilPostApplicationBodyComponentsItemProbeHTTPGet(
							v1.PostApplicationBodyComponentsItemProbeHTTPGet{
								Path: "/",
								Port: 80,
							},
						),
					},
				),
			},
		},
	})

	// Update PacketFilter
	enabled := true
	settings := []v1.PatchPacketFilterSettingsItem{
		{
			FromIP:             "192.0.2.0",
			FromIPPrefixLength: 24,
		},
	}
	updated, err := pfOp.Update(ctx, application.ID, &v1.PatchPacketFilter{
		IsEnabled: v1.NewOptBool(enabled),
		Settings:  settings,
	})
	require.NoError(t, err)
	require.Equal(t, updated.IsEnabled, true)
	require.Equal(t, len(updated.Settings), 1)
	require.Equal(t, (updated.Settings)[0].FromIP, "192.0.2.0")
	require.Equal(t, (updated.Settings)[0].FromIPPrefixLength, 24)

	read, err := pfOp.Read(ctx, application.ID)
	require.NoError(t, err)
	require.Equal(t, read.IsEnabled, true)
	require.Equal(t, len(read.Settings), 1)
	require.Equal(t, (read.Settings)[0].FromIP, "192.0.2.0")
	require.Equal(t, (read.Settings)[0].FromIPPrefixLength, 24)

	// Delete Application
	appOp.Delete(ctx, application.ID)
}

// TestVersionAPI アプリケーションバージョンの一連の操作テスト
// 以下のシナリオでテストを行う
//   - アプリケーションを作成
//   - アプリケーションを更新
//   - アプリケーションバージョンの一覧を取得
//   - アプリケーションバージョンを削除
//   - アプリケーションバージョンを確認し、削除できていることを確認
//   - アプリケーションを削除
func TestVersionAPI(t *testing.T) {
	skipIfNoAPIKey(t)

	if err := cleanupTestApplication(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	client, err := newAPIClient()
	require.NoError(t, err)
	appOp := apprun.NewApplicationOp(client)
	versionOp := apprun.NewVersionOp(client)

	// Application Create
	application, _ := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:                   appName,
		TimeoutSeconds:         100,
		Port:                   80,
		MinScale:               0,
		MaxScale:               1,
		ScaleTargetConcurrency: v1.NewOptInt(100),
		Components: []v1.PostApplicationBodyComponentsItem{
			{
				Name:      "component1",
				MaxCPU:    v1.PostApplicationBodyComponentsItemMaxCPU05,
				MaxMemory: v1.PostApplicationBodyComponentsItemMaxMemory1Gi,
				DeploySource: v1.PostApplicationBodyComponentsItemDeploySource{
					ContainerRegistry: v1.NewOptPostApplicationBodyComponentsItemDeploySourceContainerRegistry(
						v1.PostApplicationBodyComponentsItemDeploySourceContainerRegistry{
							Image: "sakura-oss-dev.sakuracr.jp/test:latest",
						},
					),
				},
				Probe: v1.NewOptNilPostApplicationBodyComponentsItemProbe(
					v1.PostApplicationBodyComponentsItemProbe{
						HTTPGet: v1.NewOptNilPostApplicationBodyComponentsItemProbeHTTPGet(
							v1.PostApplicationBodyComponentsItemProbeHTTPGet{
								Path: "/",
								Port: 80,
							},
						),
					},
				),
			},
		},
	})

	// Update Application
	timeoutSeconds := 10
	appOp.Update(ctx, application.ID, &v1.PatchApplicationBody{
		TimeoutSeconds: v1.NewOptInt(timeoutSeconds),
	})

	// List Version
	versions, err := versionOp.List(ctx, application.ID, &v1.ListApplicationVersionsParams{})
	require.NoError(t, err)
	require.Equal(t, len(versions.Data), 2)

	// Delete Version
	err = versionOp.Delete(ctx, application.ID, versions.Data[1].ID)
	require.NoError(t, err)

	// List Version
	versions, err = versionOp.List(ctx, application.ID, &v1.ListApplicationVersionsParams{})
	require.NoError(t, err)
	require.Equal(t, len(versions.Data), 1)

	status, err := versionOp.ReadStatus(ctx, application.ID, versions.Data[0].ID)
	require.NoError(t, err)
	// タイミングによってはDeployingの可能性もあるため、HealthyかDeployingのどちらかであればテスト成功とする
	require.Contains(t, []string{string(v1.HandlerGetApplicationVersionOnlyStatusStatusHealthy), string(v1.HandlerGetApplicationVersionOnlyStatusStatusDeploying)}, string(status.Status))

	// Delete Application
	appOp.Delete(ctx, application.ID)
}

// TestTrafficAPI アプリケーショントラフィックの一連の操作テスト
// 以下のシナリオでテストを行う
//   - アプリケーションを作成
//   - アプリケーションを更新
//   - アプリケーショントラフィックを変更
//   - アプリケーショントラフィックを確認
//   - アプリケーションを削除
func TestTrafficAPI(t *testing.T) {
	skipIfNoAPIKey(t)

	if err := cleanupTestApplication(); err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	client, err := newAPIClient()
	require.NoError(t, err)
	appOp := apprun.NewApplicationOp(client)
	versionOp := apprun.NewVersionOp(client)
	trafficOp := apprun.NewTrafficOp(client)

	// Application Create
	application, _ := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:           appName,
		TimeoutSeconds: 100,
		Port:           80,
		MinScale:       0,
		MaxScale:       1,
		Components: []v1.PostApplicationBodyComponentsItem{
			{
				Name:      "component1",
				MaxCPU:    v1.PostApplicationBodyComponentsItemMaxCPU05,
				MaxMemory: v1.PostApplicationBodyComponentsItemMaxMemory1Gi,
				DeploySource: v1.PostApplicationBodyComponentsItemDeploySource{
					ContainerRegistry: v1.NewOptPostApplicationBodyComponentsItemDeploySourceContainerRegistry(
						v1.PostApplicationBodyComponentsItemDeploySourceContainerRegistry{
							Image: "sakura-oss-dev.sakuracr.jp/test:latest",
						},
					),
				},
				Probe: v1.NewOptNilPostApplicationBodyComponentsItemProbe(
					v1.PostApplicationBodyComponentsItemProbe{
						HTTPGet: v1.NewOptNilPostApplicationBodyComponentsItemProbeHTTPGet(
							v1.PostApplicationBodyComponentsItemProbeHTTPGet{
								Path: "/",
								Port: 80,
							},
						),
					},
				),
			},
		},
	})

	// Update Application
	timeoutSeconds := 10
	appOp.Update(ctx, application.ID, &v1.PatchApplicationBody{
		TimeoutSeconds: v1.NewOptInt(timeoutSeconds),
	})

	// Update Application Traffic
	versions, _ := versionOp.List(ctx, application.ID, &v1.ListApplicationVersionsParams{})

	v0IsLatestVersion := true
	v0Percent := 90

	v1Name := versions.Data[1].Name
	v1Percent := 10

	trafficBody := v1.PutTrafficsBody{
		v1.NewPutTrafficsBodyItem0PutTrafficsBodyItem(v1.PutTrafficsBodyItem0{
			IsLatestVersion: v0IsLatestVersion,
			Percent:         v0Percent,
		}),
		v1.NewPutTrafficsBodyItem1PutTrafficsBodyItem(v1.PutTrafficsBodyItem1{
			VersionName: v1Name,
			Percent:     v1Percent,
		}),
	}

	_, err = trafficOp.Update(ctx, application.ID, &trafficBody)
	require.NoError(t, err)

	// List Application Traffic
	traffics, err := trafficOp.List(ctx, application.ID)
	require.NoError(t, err)

	require.Equal(t, len(traffics.Data), 2)
	require.Equal(t, traffics.Data[0].IsLatestVersion, v0IsLatestVersion)
	require.Equal(t, traffics.Data[0].Percent, v0Percent)
	require.Equal(t, traffics.Data[1].VersionName, v1Name)
	require.Equal(t, traffics.Data[1].Percent, v1Percent)

	// Delete Application
	appOp.Delete(ctx, application.ID)
}

// skipIfNoEnv 指定の環境変数のいずれかが空の場合はt.SkipNow()する
func skipIfNoEnv(t *testing.T, envs ...string) {
	var emptyEnvs []string
	for _, env := range envs {
		if os.Getenv(env) == "" {
			emptyEnvs = append(emptyEnvs, env)
		}
	}
	if len(emptyEnvs) > 0 {
		for _, env := range emptyEnvs {
			t.Logf("environment variable %q is not set", env)
		}
		t.SkipNow()
	}
}

func skipIfNoAPIKey(t *testing.T) {
	skipIfNoEnv(t, "SAKURA_ACCESS_TOKEN", "SAKURA_ACCESS_TOKEN_SECRET")
}

func cleanupTestApplication() error {
	ctx := context.Background()
	client, err := newAPIClient()
	if err != nil {
		return err
	}

	appOp := apprun.NewApplicationOp(client)
	apps, err := appOp.List(ctx, &v1.ListApplicationsParams{})
	if err != nil {
		return err
	}
	if apps.Data == nil {
		return nil
	}

	for _, app := range apps.Data {
		if app.Name == appName {
			if err := appOp.Delete(ctx, app.ID); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func newAPIClient() (*v1.Client, error) {
	var theClient saclient.Client
	return apprun.NewClient(&theClient)
}
