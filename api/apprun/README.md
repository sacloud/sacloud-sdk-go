# sacloud-sdk-go/api/apprun

Go言語向けのさくらのクラウド AppRun APIライブラリ

AppRun共用型ドキュメント: https://manual.sakura.ad.jp/cloud/manual-sakura-apprun.html
AppRun共用型 APIドキュメント: https://manual.sakura.ad.jp/api/cloud/portal/?api=apprun-shared-api

## 概要
sacloud-sdk-go/api/apprunはさくらのクラウド AppRun共用型 APIをGo言語から利用するためのAPIライブラリです。

利用イメージ:

```go
package main

import (
	"context"
	"fmt"

	"github.com/sacloud/sacloud-sdk-go/api/apprun"
	v1 "github.com/sacloud/sacloud-sdk-go/api/apprun/apis/v1"
	"github.com/sacloud/sacloud-sdk-go/common/saclient"
)

func main() {
	// デフォルトでusacloud互換プロファイル or 環境変数(SAKURA_ACCESS_TOKEN{_SECRET})が利用される
	var theClient saclient.Client
	client, err := apprun.NewClient(&theClient)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()

	// アプリケーションを作成
	appOp := apprun.NewApplicationOp(client)
	application, err := appOp.Create(ctx, &v1.CreateApplicationBody{
		Name:           "example-app1",
		TimeoutSeconds: 100,
		Port:           80,
		MinScale:       0,
		MaxScale:       1,
		Components: []v1.CreateApplicationBodyComponentsItem{
			{
				Name:      "component1",
				MaxCPU:    v1.CreateApplicationBodyComponentsItemMaxCPU05,
				MaxMemory: v1.CreateApplicationBodyComponentsItemMaxMemory1Gi,
				DeploySource: v1.CreateApplicationBodyComponentsItemDeploySource{
					ContainerRegistry: v1.NewOptCreateApplicationBodyComponentsItemDeploySourceContainerRegistry(
						v1.CreateApplicationBodyComponentsItemDeploySourceContainerRegistry{
							Image:    "sakura-oss-dev.sakuracr.jp/test:latest",
							Server:   v1.NewOptNilString("sakura-oss-dev.sakuracr.jp"),
							Username: v1.NewOptNilString("test-user"),
							Password: v1.NewOptNilString("test-password"),
						},
					),
				},
				Probe: v1.NewOptNilCreateApplicationBodyComponentsItemProbe(
					v1.CreateApplicationBodyComponentsItemProbe{
						HTTPGet: v1.NewOptNilCreateApplicationBodyComponentsItemProbeHTTPGet(
							v1.CreateApplicationBodyComponentsItemProbeHTTPGet{
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

	// アプリケーションバージョンを取得
	versionOp := apprun.NewVersionOp(client)
	versions, err := versionOp.List(ctx, application.ID, &v1.ListApplicationVersionsParams{})
	if err != nil {
		panic(err)
	}

	// アプリケーションの削除
	defer func() {
		if err := appOp.Delete(ctx, application.ID); err != nil {
			panic(err)
		}
	}()

	v := versions.Data[0]
	fmt.Println(v.Name)
}
```

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

sacloud-sdk-go/api/apprunはv0.8からogenベースの実装となっています。oapi-codegenベースの実装を使いたい場合にはv0.7系を使ってください。ただし新機能は追加されないため、新規APIを利用したい場合には移行が必要となります。

## License

Copyright (C) 2021-2026 The sacloud/sacloud-sdk-go Authors.
このプロジェクトは[Apache 2.0 License](LICENSE)の下で公開されています。
