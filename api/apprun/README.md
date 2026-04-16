# sacloud/apprun-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/apprun-api-go.svg)](https://pkg.go.dev/github.com/sacloud/apprun-api-go)
[![Tests](https://github.com/sacloud/apprun-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/apprun-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/apprun-api-go)](https://goreportcard.com/report/github.com/sacloud/apprun-api-go)

Go言語向けのさくらのクラウド AppRun APIライブラリ

AppRun APIドキュメント: https://manual.sakura.ad.jp/sakura-apprun-api/spec.html

## 概要
sacloud/apprun-api-goはさくらのクラウド AppRun APIをGo言語から利用するためのAPIライブラリです。

利用イメージ:

```go
package main

import (
	"context"
	"fmt"

	"github.com/sacloud/apprun-api-go"
	v1 "github.com/sacloud/apprun-api-go/apis/v1"
)

func main() {
	// デフォルトでusacloud互換プロファイル or 環境変数(SAKURACLOUD_ACCESS_TOKEN{_SECRET})が利用される
	client := &apprun.Client{}

	ctx := context.Background()

	// アプリケーションを作成
	appOp := apprun.NewApplicationOp(client)
	application, err := appOp.Create(ctx, &v1.PostApplicationBody{
		Name:           "example-app1",
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

	// アプリケーションバージョンを取得
	versionOp := apprun.NewVersionOp(client)
	versions, err := versionOp.List(ctx, *application.Id, &v1.ListApplicationVersionsParams{})
	if err != nil {
		panic(err)
	}

	// アプリケーションの削除
	defer func() {
		if err := appOp.Delete(ctx, *application.Id); err != nil {
			panic(err)
		}
	}()

	v := (*versions.Data)[0]
	fmt.Println(*v.Name)
}
```

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## License

`apprun-api-go` Copyright (C) 2022-2023 The sacloud/apprun-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).
