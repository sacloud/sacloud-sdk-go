# cloudhsm-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/cloudhsm-api-go.svg)](https://pkg.go.dev/github.com/sacloud/cloudhsm-api-go)
[![Tests](https://github.com/sacloud/cloudhsm-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/cloudhsm-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/cloudhsm-api-go)](https://goreportcard.com/report/github.com/sacloud/cloudhsm-api-go)

さくらのクラウド「クラウドHSM」APIのGoクライアントライブラリ

## 概要

このライブラリは、さくらのクラウド「クラウドHSM」APIをGo言語から利用するためのクライアントです。
OpenAPI仕様から自動生成された型安全なAPIクライアントと、それをラップして使い勝手を向上させたクライアントを提供します。

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。


## インストール

```bash
go get github.com/sacloud/cloudhsm-api-go
```

## 使い方

### 事前準備

SDKを利用開始するにはさくらのクラウド「クラウドHSM」自体を利用開始する必要があります。執筆時点ではさくらのクラウドの利用申請とは別途クラウドHSMの利用申請も必要です。申請手順の詳細に関しては[クラウドHSMマニュアル](https://manual.sakura.ad.jp/cloud/appliance/cloudhsm/index.html)をご参照ください。

### 認証情報

APIを実行するには認証が必要です。インタラクティブな環境の場合おすすめは [`usacloud`](https://github.com/sacloud/usacloud) を使って設定ファイルを作成することです。たとえば

```sh
usacloud config create --name is1a
```

にて作成したプロファイル `is1a` があるとすると、SDKとしては、

```golang
import (
    "github.com/sacloud/cloudhsm-api-go"
    "github.com/sacloud/saclient-go"
)

var theClient saclient.Client

func main() {
    _ = theClient.SetEnviron([]string{"SAKURA_PROFILE=is1a"})
    client, err := cloudhsm.NewClient(&theClient)

    // 以下略
}
```

のようにして読み込むことができます。

一方でCI環境のようにファイルに書き出すのが適切ではない場合、環境変数経由で

```golang
import (
    "os"

    "github.com/sacloud/cloudhsm-api-go"
    "github.com/sacloud/saclient-go"
)

var theClient saclient.Client

func main() {
    _ = theClient.SetEnviron([]string{
        "SAKURA_ZONE=is1a",
        "SAKURA_SERVICE_PRINCIPAL_ID=something",
        // 他、　os.Environ()から必要な環境変数を追加
    })
    client, err := cloudhsm.NewClient(&theClient)

    // 以下略
}
```

のように指定できます。

### SDK

```go
package main

import (
    "context"

    "github.com/sacloud/cloudhsm-api-go"
    v1 "github.com/sacloud/cloudhsm-api-go/apis/v1"
)

func Logic(ctx context.Context, client *v1.Client) {
    // 例: ライセンス一覧取得
    licenses, err := cloudhsm.NewLicenseOp(client).List(ctx)
    if err != nil {
        // エラーハンドリング
    }
}
```

APIの詳細は[GoDoc](https://pkg.go.dev/github.com/sacloud/cloudhsm-api-go)や`apis/v1/`配下の型定義を参照してください。

## OpenAPI仕様について

`openapi/openapi.json`は[KMS/SecretManager/CloudHSM API](https://manual.sakura.ad.jp/api/cloud/security-encryption/)からダウンロードしたものを一部加工しています。

## 開発

ビルドやテストはMakefile経由で実行できます。

```bash
make
make test
```

## ライセンス

Copyright (C) 2022-2025 The sacloud/cloudhsm-api-go Authors.
このプロジェクトは[Apache 2.0 License](LICENSE)の下で公開されています。