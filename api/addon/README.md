# addon-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/addon-api-go.svg)](https://pkg.go.dev/github.com/sacloud/addon-api-go)
[![Tests](https://github.com/sacloud/addon-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/addon-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/addon-api-go)](https://goreportcard.com/report/github.com/sacloud/addon-api-go)

さくらのクラウド「Add-on」APIのGoクライアントライブラリ

## 概要

このライブラリは、さくらのクラウド「Add-on」APIをGo言語から利用するためのクライアントです。
OpenAPI仕様から自動生成された型安全なAPIクライアントと、それをラップして使い勝手を向上させたクライアントを提供します。

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。


## インストール

```bash
go get github.com/sacloud/addon-api-go
```

## 使い方

### 事前準備

SDKを利用開始するにはさくらのクラウド「Add-on」自体を利用開始する必要があります。執筆時点ではさくらのクラウドの利用申請とは別途Add-onの利用申請も必要です。申請手順の詳細に関しては[Add-onマニュアル](https://manual.sakura.ad.jp/cloud/add-on/index.html)をご参照ください。

またSDKを利用するにはAPIキーが必要ですが、これには「Add-on」のアクセスレベルを付与してください。アクセスレベルに関しては[アクセスレベル マニュアル](https://manual.sakura.ad.jp/cloud/controlpanel/access-level.html)をご参照ください。

### SDK

```go
package main

import (
    "context"
    "os"

    "github.com/sacloud/addon-api-go"
    "github.com/sacloud/saclient-go"
)

var theClient saclient.Client

func main() {
	err := theClient.FlagSet(flag.PanicOnError).Parse(os.Args[1:])
    if err != nil {
        // エラーハンドリング
	}

    err = theClient.SetEnviron(os.Environ())
    if err != nil {
        // エラーハンドリング
	}

    client, err := addon.NewClient(&theClient)
    if err != nil {
        // エラーハンドリング
    }

    // 例: データレーク一覧取得
    ctx := context.Background()
    lakes, err := addon.NewDataLakeOp(client).List(ctx)
    if err != nil {
        // エラーハンドリング
    }
}
```

APIの詳細は[GoDoc](https://pkg.go.dev/github.com/sacloud/addon-api-go)や`apis/v1/`配下の型定義を参照してください。

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

`openapi/openapi.json`は[Add-on API](https://manual.sakura.ad.jp/api/cloud/addon/)からダウンロードしたものを一部加工しています。

## 開発

ビルドやテストはMakefile経由で実行できます。

```bash
make
make test
```

## ライセンス

Copyright (C) 2025- The sacloud/addon-api-go Authors.
このプロジェクトは[Apache 2.0 License](LICENSE)の下で公開されています。