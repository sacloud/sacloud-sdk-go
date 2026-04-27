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

package apprun_test

import (
	"context"
	"fmt"

	apprun "github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
)

const defaultServerURL = "https://secure.sakura.ad.jp/cloud/api/apprun/1.0/apprun/api"

var serverURL = defaultServerURL

// Example_userAPI ユーザーAPIの利用例
func Example_userAPI() {
	var theClient saclient.Client
	client, err := apprun.NewClientWithAPIRootURL(&theClient, serverURL)
	if err != nil {
		panic(err)
	}

	// ユーザー情報の取得
	userOp := apprun.NewUserOp(client)
	res, err := userOp.Read(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(res.Limit.ApplicationCount >= 0)
	// output:
	// true
}

// Example_applicationAPI アプリケーションAPIの利用例
func Example_applicationAPI() {
	var theClient saclient.Client
	client, err := apprun.NewClientWithAPIRootURL(&theClient, serverURL)
	if err != nil {
		panic(err)
	}

	// アプリケーションの作成
	ctx := context.Background()
	appOp := apprun.NewApplicationOp(client)

	created, err := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:           "example-app",
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
							Image: "apprun-test.sakuracr.jp/apprun/test1:latest",
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
	if err != nil {
		panic(err)
	}

	// アプリケーションの参照
	application, err := appOp.Read(ctx, created.ID)
	if err != nil {
		panic(err)
	}

	// アプリケーションの削除
	err = appOp.Delete(ctx, application.ID)
	if err != nil {
		panic(err)
	}

	fmt.Println(application.Name)
	// output:
	// example-app
}

// Example_versionAPI アプリケーションバージョンAPIの利用例
func Example_versionAPI() {
	var theClient saclient.Client
	client, err := apprun.NewClientWithAPIRootURL(&theClient, serverURL)
	if err != nil {
		panic(err)
	}

	// アプリケーションの作成
	ctx := context.Background()
	appOp := apprun.NewApplicationOp(client)
	versionOp := apprun.NewVersionOp(client)

	application, err := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:           "example-app",
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
							Image: "apprun-test.sakuracr.jp/apprun/test1:latest",
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
	if err != nil {
		panic(err)
	}

	// アプリケーションの更新
	timeoutSeconds := 10
	_, err = appOp.Update(ctx, application.ID, &v1.PatchApplicationBody{
		TimeoutSeconds: v1.NewOptInt(timeoutSeconds),
	})
	if err != nil {
		panic(err)
	}

	// バージョン一覧の取得
	versions, err := versionOp.List(ctx, application.ID, &v1.ListApplicationVersionsParams{})
	if err != nil {
		panic(err)
	}
	if len(versions.Data) != 2 {
		fmt.Println(len(versions.Data))
		panic("ListVersions failed")
	}

	d0 := versions.Data[0]
	d1 := versions.Data[1]

	// バージョンの削除
	err = versionOp.Delete(ctx, application.ID, d0.ID)
	if err != nil {
		panic(err)
	}

	// バージョンの参照
	version, err := versionOp.Read(ctx, application.ID, d1.ID)
	if err != nil {
		panic(err)
	}

	fmt.Printf("version status: %s", version.Status)
	// output:
	// version status: Healthy
}

// Example_trafficAPI アプリケーショントラフィックAPIの利用例
func Example_trafficAPI() {
	var theClient saclient.Client
	client, err := apprun.NewClientWithAPIRootURL(&theClient, serverURL)
	if err != nil {
		panic(err)
	}

	// アプリケーションの作成
	ctx := context.Background()
	appOp := apprun.NewApplicationOp(client)
	versionOp := apprun.NewVersionOp(client)
	trafficOp := apprun.NewTrafficOp(client)

	application, err := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:           "example-app",
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
							Image: "apprun-test.sakuracr.jp/apprun/test1:latest",
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
	if err != nil {
		panic(err)
	}

	// アプリケーションの更新
	timeoutSeconds := 10
	_, err = appOp.Update(ctx, application.ID, &v1.PatchApplicationBody{
		TimeoutSeconds: v1.NewOptInt(timeoutSeconds),
	})
	if err != nil {
		panic(err)
	}

	// バージョン一覧の取得
	versions, err := versionOp.List(ctx, application.ID, &v1.ListApplicationVersionsParams{})
	if err != nil {
		panic(err)
	}

	// トラフィック分散を更新
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
	if err != nil {
		panic(err)
	}

	// トラフィック分散を取得
	traffics, err := trafficOp.List(ctx, application.ID)
	if err != nil {
		panic(err)
	}

	for _, data := range traffics.Data {
		if data.IsLatestVersion {
			fmt.Printf("is_latest_version: %t, percent: %d", data.IsLatestVersion, data.Percent)
		}
	}
	// output:
	// is_latest_version: true, percent: 90
}
