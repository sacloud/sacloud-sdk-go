# sacloud/apprun-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/apprun-api-go.svg)](https://pkg.go.dev/github.com/sacloud/apprun-api-go)
[![Tests](https://github.com/sacloud/apprun-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/apprun-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/apprun-api-go)](https://goreportcard.com/report/github.com/sacloud/apprun-api-go)

Go言語向けのさくらのクラウド AppRun APIライブラリ

AppRun共用型ドキュメント: https://manual.sakura.ad.jp/cloud/manual-sakura-apprun.html
AppRun共用型 APIドキュメント: https://manual.sakura.ad.jp/api/cloud/portal/?api=apprun-shared-api

## 概要
sacloud/apprun-api-goはさくらのクラウド AppRun共用型 APIをGo言語から利用するためのAPIライブラリです。

利用イメージ:

```go
package main

import (
	"context"
	"fmt"

	"github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
	"github.com/sacloud/saclient-go"
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
	application, err := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:           "example-app1",
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

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

apprun-api-goはv0.8からogenベースの実装となっています。oapi-codegenベースの実装を使いたい場合にはv0.7系を使ってください。ただし新機能は追加されないため、新規APIを利用したい場合には移行が必要となります。

## License

`apprun-api-go` Copyright (C) 2021-2026 The sacloud/apprun-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).
