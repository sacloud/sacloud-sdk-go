# apprun-dedicated-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/apprun-dedicated-api-go.svg)](https://pkg.go.dev/github.com/sacloud/apprun-dedicated-api-go)
[![Tests](https://github.com/sacloud/apprun-dedicated-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/apprun-dedicated-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/apprun-dedicated-api-go)](https://goreportcard.com/report/github.com/sacloud/apprun-dedicated-api-go)

さくらのクラウド「AppRun 専有型」APIのGoクライアントライブラリ

## 概要

このライブラリは、さくらのクラウド「AppRun 専有型」APIをGo言語から利用するためのクライアントです。
OpenAPI仕様から自動生成された型安全なAPIクライアントと、それをラップして使い勝手を向上させたクライアントを提供します。

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。


## インストール

```bash
go get github.com/sacloud/apprun-dedicated-api-go
```

## 使い方


```go
package main

import (
    "context"
    "os"

    "github.com/sacloud/apprun-dedicated-api-go"
    "github.com/sacloud/saclient-go"
)

var theClient saclient.Client

func main() {
    ctx := context.Background()
    err := theClient.FlagSet(flag.PanicOnError).Parse(os.Args[1:])
    if err != nil {
        // エラーハンドリング
    }

    err = theClient.SetEnviron(os.Environ())
    if err != nil {
        // エラーハンドリング
    }

    client, err := apprun_dedicated.NewClient(&theClient)
    if err != nil {
        // エラーハンドリング
    }

    // 例: クラスタ一覧取得
    clusters, err := apprun_dedicated.NewClusterOp(client).List(ctx, nil, nil)
    if err != nil {
        // エラーハンドリング
    }
}
```

APIの詳細は[GoDoc](https://pkg.go.dev/github.com/sacloud/apprun-dedicated-api-go)や`apis/v1/`配下の型定義を参照してください。

### 認証情報

APIを実行するには認証が必要です。インタラクティブな環境の場合おすすめは [`usacloud`](https://github.com/sacloud/usacloud) を使って設定ファイルを作成することです。たとえば

```sh
usacloud config create --name production
```

にて作成したプロファイル `production` があるとすると、上記の`main`関数をもつバイナリは

```sh
./a.out --profile=production
```

のようにして設定を読み込むことができます。

一方でCI環境のようにファイルに書き出すのが適切ではない場合、環境変数経由で

```sh
export SAKURA_ACCESS_TOKEN=TOKEN
export SAKURA_ACCESS_TOKEN_SECRET=SECRET

./a.out
```

のように指定できます。

## OpenAPI仕様について

`openapi/openapi.json`は[AppRun専有型 API](https://manual.sakura.ad.jp/api/cloud/apprun-dedicated/)からダウンロードしたものです。

## 開発

ビルドやテストはMakefile経由で実行できます。

```bash
make
make test
```

## ライセンス

Copyright (C) 2022-2026 The sacloud/apprun-dedicated-api-go Authors.
このプロジェクトは[Apache 2.0 License](LICENSE)の下で公開されています。
