
# monitoring-suite-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/monitoring-suite-api-go.svg)](https://pkg.go.dev/github.com/sacloud/monitoring-suite-api-go)
[![Tests](https://github.com/sacloud/monitoring-suite-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/monitoring-suite-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/monitoring-suite-api-go)](https://goreportcard.com/report/github.com/sacloud/monitoring-suite-api-go)

さくらのクラウド「モニタリングスイート」APIのGoクライアントライブラリ

## 概要

このライブラリは、さくらのクラウド「モニタリングスイート」APIをGo言語から利用するためのクライアントです。
OpenAPI仕様から自動生成された型安全なAPIクライアントと、それをラップして使い勝手を向上させたクライアントを提供します。

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。


## インストール

```bash
go get github.com/sacloud/monitoring-suite-api-go
```

## 使い方

```go
import (
    "context"

    "github.com/sacloud/monitoring-suite-api-go"
)

ctx := context.Background()
client, err := monitoringsuite.NewClient()
if err != nil {
	// エラーハンドリング
}
// 例: アラートプロジェクト一覧取得
projects, err := NewAlertProjectOp(client).List(ctx, 32768, 0)
```

APIの詳細は[GoDoc](https://pkg.go.dev/github.com/sacloud/monitoring-suite-api-go)や`apis/v1/`配下の型定義を参照してください。

## OpenAPI仕様について

`openapi/openapi.json`は[モニタリングスイート API ドキュメント](https://manual.sakura.ad.jp/api/cloud/monitoring-suite/)からダウンロードしたものを一部加工しています。

```console
$ jq 'del(.paths.[].[].requestBody.content.["application/x-www-form-urlencoded", "multipart/form-data"])' openapi.json
```

## 開発

ビルドやテストはMakefile経由で実行できます。

```bash
make
make test
```

## ライセンス

Copyright (C) 2022-2025 The sacloud/monitoring-suite-api-go Authors.
このプロジェクトは[Apache 2.0 License](LICENSE)の下で公開されています。