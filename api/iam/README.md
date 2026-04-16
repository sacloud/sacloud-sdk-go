# iam-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/iam-api-go.svg)](https://pkg.go.dev/github.com/sacloud/iam-api-go)
[![Tests](https://github.com/sacloud/iam-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/iam-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/iam-api-go)](https://goreportcard.com/report/github.com/sacloud/iam-api-go)

さくらのクラウド「IAM」APIのGoクライアントライブラリ

## 概要

このライブラリは、さくらのクラウド「IAM」APIをGo言語から利用するためのクライアントです。
OpenAPI仕様から自動生成された型安全なAPIクライアントと、それをラップして使い勝手を向上させたクライアントを提供します。

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。


## インストール

```bash
go get github.com/sacloud/iam-api-go
```

## 使い方

```go
package main

import (
    "context"
    "os"

    "github.com/sacloud/iam-api-go"
    "github.com/sacloud/iam-api-go/apis/user"
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

    client, err := iam.NewClient(&theClient)
    if err != nil {
        // エラーハンドリング
    }

    // 例: ユーザー一覧取得
    ctx := context.Background()
    users, err := iam.NewUserOp(client).List(ctx, user.ListParams{})
    if err != nil {
        // エラーハンドリング
    }
}
```

APIの詳細は[GoDoc](https://pkg.go.dev/github.com/sacloud/iam-api-go)や`apis/v1/`配下の型定義を参照してください。

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
export SAKURA_SERVICE_PRINCIPAL_ID=sub
export SAKURA_SERVICE_PRINCIPAL_KEY_ID=kid
export SAKURA_PRIVATE_KEY=-----BEGIN PRIVATE KEY-----...

./a.out --profile=production
```

のように指定できます。

なおIAMの各機能を利用するには各機能を利用できる権限を付与したサービスプリンシパルが必要です。
さくらのクラウド マニュアルの [サービスプリンシパル](https://manual.sakura.ad.jp/cloud/controlpanel/service-principal.html) 内「ロールの付与」、および [IAMポリシー](https://manual.sakura.ad.jp/cloud/controlpanel/iam-policy.html) 内「設定方法」をご参照ください。

単体テストを流すのに必要な権限に関しては [./doc/testing.md](./doc/testing.md) もご参照ください。

## 開発

ビルドやテストはMakefile経由で実行できます。

```bash
make
make test
```

## ライセンス

Copyright (C) 2025- The sacloud/iam-api-go Authors.
このプロジェクトは[Apache 2.0 License](LICENSE)の下で公開されています。
