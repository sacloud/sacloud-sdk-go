// Copyright 2021-2024 The sacloud/apprun-api-go authors
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
	client := &apprun.Client{
		APIRootURL: serverURL, // 省略可能
		Saclient:   &theClient,
	}

	// ユーザー情報の取得
	userOp := apprun.NewUserOp(client)
	res, err := userOp.Read(context.Background())
	if err != nil {
		panic(err)
	}

	fmt.Println(res.StatusCode)
	// output:
	// 200
}

// Example_applicationAPI アプリケーションAPIの利用例
func Example_applicationAPI() {
	client := &apprun.Client{
		APIRootURL: serverURL, // 省略可能
	}

	// アプリケーションの作成
	ctx := context.Background()
	appOp := apprun.NewApplicationOp(client)

	application, err := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:           "example-app",
		TimeoutSeconds: 100,
		Port:           80,
		MinScale:       0,
		MaxScale:       1,
		Components: []v1.PostApplicationBodyComponent{
			{
				Name:      "component1",
				MaxCpu:    "0.5",
				MaxMemory: "1Gi",
				DeploySource: v1.PostApplicationBodyComponentDeploySource{
					ContainerRegistry: &v1.PostApplicationBodyComponentDeploySourceContainerRegistry{
						Image: "apprun-test.sakuracr.jp/apprun/test1:latest",
					},
				},
				Probe: &v1.PostApplicationBodyComponentProbe{
					HttpGet: &v1.PostApplicationBodyComponentProbeHttpGet{
						Path: "/",
						Port: 80,
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// アプリケーションの参照
	application, err = appOp.Read(ctx, application.Id)
	if err != nil {
		panic(err)
	}

	// アプリケーションの削除
	err = appOp.Delete(ctx, application.Id)
	if err != nil {
		panic(err)
	}

	fmt.Println(application.Name)
	// output:
	// example-app
}

// Example_versionAPI アプリケーションバージョンAPIの利用例
func Example_versionAPI() {
	client := &apprun.Client{
		APIRootURL: serverURL, // 省略可能
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
		Components: []v1.PostApplicationBodyComponent{
			{
				Name:      "component1",
				MaxCpu:    "0.5",
				MaxMemory: "1Gi",
				DeploySource: v1.PostApplicationBodyComponentDeploySource{
					ContainerRegistry: &v1.PostApplicationBodyComponentDeploySourceContainerRegistry{
						Image: "apprun-test.sakuracr.jp/apprun/test1:latest",
					},
				},
				Probe: &v1.PostApplicationBodyComponentProbe{
					HttpGet: &v1.PostApplicationBodyComponentProbeHttpGet{
						Path: "/",
						Port: 80,
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// アプリケーションの更新
	timeoutSeconds := 10
	_, err = appOp.Update(ctx, application.Id, &v1.PatchApplicationBody{
		TimeoutSeconds: &timeoutSeconds,
	})
	if err != nil {
		panic(err)
	}

	// バージョン一覧の取得
	versions, err := versionOp.List(ctx, application.Id, &v1.ListApplicationVersionsParams{})
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
	err = versionOp.Delete(ctx, application.Id, d0.Id)
	if err != nil {
		panic(err)
	}

	// バージョンの参照
	version, err := versionOp.Read(ctx, application.Id, d1.Id)
	if err != nil {
		panic(err)
	}

	fmt.Printf("version status: %s", version.Status)
	// output:
	// version status: Healthy
}

// Example_trafficAPI アプリケーショントラフィックAPIの利用例
func Example_trafficAPI() {
	client := &apprun.Client{
		APIRootURL: serverURL, // 省略可能
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
		Components: []v1.PostApplicationBodyComponent{
			{
				Name:      "component1",
				MaxCpu:    "0.5",
				MaxMemory: "1Gi",
				DeploySource: v1.PostApplicationBodyComponentDeploySource{
					ContainerRegistry: &v1.PostApplicationBodyComponentDeploySourceContainerRegistry{
						Image: "apprun-test.sakuracr.jp/apprun/test1:latest",
					},
				},
				Probe: &v1.PostApplicationBodyComponentProbe{
					HttpGet: &v1.PostApplicationBodyComponentProbeHttpGet{
						Path: "/",
						Port: 80,
					},
				},
			},
		},
	})
	if err != nil {
		panic(err)
	}

	// アプリケーションの更新
	timeoutSeconds := 10
	_, err = appOp.Update(ctx, application.Id, &v1.PatchApplicationBody{
		TimeoutSeconds: &timeoutSeconds,
	})
	if err != nil {
		panic(err)
	}

	// バージョン一覧の取得
	versions, err := versionOp.List(ctx, application.Id, &v1.ListApplicationVersionsParams{})
	if err != nil {
		panic(err)
	}

	// トラフィック分散を更新
	v0IsLatestVersion := true
	v0Percent := 90

	v1Name := versions.Data[1].Name
	v1Percent := 10

	v0tr := &v1.Traffic{}
	if err := v0tr.FromTrafficWithLatestVersion(v1.TrafficWithLatestVersion{
		IsLatestVersion: v0IsLatestVersion,
		Percent:         v0Percent,
	}); err != nil {
		panic(err)
	}

	v1tr := &v1.Traffic{}
	if err := v1tr.FromTrafficWithVersionName(v1.TrafficWithVersionName{
		VersionName: v1Name,
		Percent:     v1Percent,
	}); err != nil {
		panic(err)
	}

	_, err = trafficOp.Update(ctx, application.Id, &[]v1.Traffic{*v0tr, *v1tr})
	if err != nil {
		panic(err)
	}

	// トラフィック分散を取得
	traffics, err := trafficOp.List(ctx, application.Id)
	if err != nil {
		panic(err)
	}

	for _, data := range traffics.Data {
		withLatest, _ := data.AsTrafficWithLatestVersion()
		if withLatest.IsLatestVersion {
			fmt.Printf("is_latest_version: %t, percent: %d", withLatest.IsLatestVersion, withLatest.Percent)
		}
	}
	// output:
	// is_latest_version: true, percent: 90
}
