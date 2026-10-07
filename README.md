# sacloud-sdk-go

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

さくらのクラウド Go SDK

## 概要

sacloud-sdk-goは、[さくらのクラウド](https://cloud.sakura.ad.jp/)の各種サービスをGo言語から利用するためのSDK群をまとめたモノレポジトリです。

このリポジトリには、さくらのクラウドの各サービスに対応したAPIクライアントライブラリと、それらを利用した高レベルなサービスライブラリが含まれています。

> [!WARNING]
> 本プロジェクトの一部パッケージはv1.0に達するまでの間、互換性のない形で変更される可能性があります。

## 構成

### APIクライアントライブラリ (`api/`)

さくらのクラウド各サービスのREST APIを操作するための低レベルクライアントライブラリです。

| パッケージ | 説明 | importパス |
|-----------|------|-------------|
| `api/iaas` | IaaS API (サーバー、ディスク、ネットワークなど) | `github.com/sacloud/sacloud-sdk-go/api/iaas` |
| `api/webaccel` | ウェブアクセラレータ API | `github.com/sacloud/sacloud-sdk-go/api/webaccel` |
| `api/iam` | IAM API (認証・認可) | `github.com/sacloud/sacloud-sdk-go/api/iam` |
| `api/addon` | Add-on API (データレーク、WAFなど) | `github.com/sacloud/sacloud-sdk-go/api/addon` |
| `api/apigw` | API Gateway API | `github.com/sacloud/sacloud-sdk-go/api/apigw` |
| `api/apprun` | AppRun API | `github.com/sacloud/sacloud-sdk-go/api/apprun` |
| `api/apprun-dedicated` | AppRun Dedicated API | `github.com/sacloud/sacloud-sdk-go/api/apprun-dedicated` |
| `api/cloudhsm` | CloudHSM API | `github.com/sacloud/sacloud-sdk-go/api/cloudhsm` |
| `api/dedicated-storage` | 専用ストレージ API | `github.com/sacloud/sacloud-sdk-go/api/dedicated-storage` |
| `api/eventbus` | EventBus API | `github.com/sacloud/sacloud-sdk-go/api/eventbus` |
| `api/kms` | Key Management Service API | `github.com/sacloud/sacloud-sdk-go/api/kms` |
| `api/monitoring-suite` | 監視サービス API | `github.com/sacloud/sacloud-sdk-go/api/monitoring-suite` |
| `api/networking-suite` | ネットワークスイート API | `github.com/sacloud/sacloud-sdk-go/api/networking-suite` |
| `api/nosql` | NoSQL API | `github.com/sacloud/sacloud-sdk-go/api/nosql` |
| `api/object-storage` | オブジェクトストレージ API | `github.com/sacloud/sacloud-sdk-go/api/object-storage` |
| `api/secretmanager` | Secret Manager API | `github.com/sacloud/sacloud-sdk-go/api/secretmanager` |
| `api/security-control` | セキュリティコントロール API | `github.com/sacloud/sacloud-sdk-go/api/security-control` |
| `api/service-endpoint-gateway` | Service Endpoint Gateway API | `github.com/sacloud/sacloud-sdk-go/api/service-endpoint-gateway` |
| `api/simple-notification` | シンプル通知 API | `github.com/sacloud/sacloud-sdk-go/api/simple-notification` |
| `api/simplemq` | SimpleMQ API | `github.com/sacloud/sacloud-sdk-go/api/simplemq` |
| `api/workflows` | Workflows API | `github.com/sacloud/sacloud-sdk-go/api/workflows` |

### 高レベルサービスライブラリ (`service/`)

APIクライアントライブラリをラップし、より使いやすいインターフェースを提供する高レベルライブラリです。

| パッケージ | 説明 | importパス |
|-----------|------|-------------|
| `service/iaas` | IaaS向け高レベルAPI | `github.com/sacloud/sacloud-sdk-go/service/iaas` |
| `service/webaccel` | ウェブアクセラレータ高レベルAPI | `github.com/sacloud/sacloud-sdk-go/service/webaccel` |

### 内部パッケージ (`internal/`)

SDKの内部で使用される共有パッケージです。Goの`internal`制約により、外部モジュールから直接importできません。

| パッケージ | 説明 |
|-----------|------|
| `internal/api-client` | APIクライアントの共通実装 |
| `internal/go-http` | HTTP通信の共通実装 |
| `internal/services` | サービスレイヤーの共通実装 |

### 共通パッケージ (`common/`)

SDK外からも利用可能な共有パッケージです。

| パッケージ | 説明 | importパス |
|-----------|------|-------------|
| `common/packages` | 汎用パッケージ群 | `github.com/sacloud/sacloud-sdk-go/common/packages` |
| `common/saclient` | クライアント認証・設定の共通実装 | `github.com/sacloud/sacloud-sdk-go/common/saclient` |


### SRNパッケージ (`srn/`)

Sakura resource name(SRN)を扱うパッケージ。

| パッケージ | 説明 | importパス |
|-----------|------|-------------|
| `srn` | SRN 構造体と関連メソッドを提供 | `github.com/sacloud/sacloud-sdk-go/srn` |


## インストール

このリポジトリは単一のGoモジュールです。

### ルートモジュール

```bash
go get github.com/sacloud/sacloud-sdk-go
```

### サブパッケージの例

```bash
# IAM APIクライアントの場合
go get github.com/sacloud/sacloud-sdk-go/api/iam

# IaaS APIクライアントの場合
go get github.com/sacloud/sacloud-sdk-go/api/iaas

# 高レベルサービスライブラリ
go get github.com/sacloud/sacloud-sdk-go/service/iaas
```

## 使用方法

各パッケージの詳しい使い方は、それぞれのディレクトリにある README.mdを参照してください。

### 簡単な使用例 (IaaS)

```go
package main

import (
    "context"
    "log"

    "github.com/sacloud/sacloud-sdk-go/api/iaas"
    "github.com/sacloud/sacloud-sdk-go/common/saclient"
)

func main() {
    // クライアント作成 (環境変数やプロファイルから自動設定)
    var sa saclient.Client
    client := iaas.NewClientFromSaclient(&sa)

    // サーバー一覧取得
    ctx := context.Background()
    servers, err := iaas.NewServerOp(client).Find(ctx, &iaas.FindServerRequest{})
    if err != nil {
        log.Fatal(err)
    }

    for _, server := range servers.Servers {
        log.Printf("Server: %s (%s)", server.Name, server.ID)
    }
}
```

## 認証情報

APIを利用するには、さくらのクラウドの認証情報が必要です。`common/saclient` が環境変数、プロファイル、コマンドラインフラグなどを統合的に読み込みます。

主な方法は次の 2 つです。

### [usacloud](https://github.com/sacloud/usacloud)プロファイルの利用

インタラクティブな環境では、usacloudを使用して設定ファイルを作成することをおすすめします。

```bash
# usacloudでプロファイル作成
usacloud config create --name production

# 作成したプロファイルを使用
export SAKURA_PROFILE=production
```

### 環境変数による設定

```bash
export SAKURA_SERVICE_PRINCIPAL_ID=sub
export SAKURA_SERVICE_PRINCIPAL_KEY_KID=kid
export SAKURA_PRIVATE_KEY=-----BEGIN PRIVATE KEY-----...
```

その他の設定項目（リトライ、エンドポイント上書き、ゾーン指定など）については [common/saclient/README.md](common/saclient/README.md) を参照してください。

## 関連プロジェクト

- [usacloud](https://github.com/sacloud/usacloud): さくらのクラウド用CLIツール
- [terraform-provider-sakuracloud](https://github.com/sacloud/terraform-provider-sakuracloud): Terraformプロバイダー
- [libsacloud](https://github.com/sacloud/libsacloud): 本SDKの前身となるライブラリ（v2系まで）

## 開発

### 前提条件

- Go 1.25 以上
- Make

### ビルドとテスト

各サブディレクトリで個別にビルド・テストが可能です。

```bash
# IaaS APIクライアントのテスト
cd api/iaas
make test

# 全パッケージのテストは各ディレクトリで順次実行
cd api/iaas && make test
cd api/webaccel && make test
# ...
```

### モノレポ全体のテスト

ルートには横断的な Makefile はないため、以下のように各ディレクトリで順次実行してください。

```bash
for d in api/*/ service/*/ common/*/ internal/*/ srn; do
    if [ -f "$d/Makefile" ]; then
        make -C "$d" test
    fi
done
```

### 共通ターゲット

各 Makefile で利用できる主なターゲットです。

| ターゲット | 内容 |
|-----------|------|
| `make test` | ユニットテストを実行 |
| `make testacc` | 受け入れテストを実行 |
| `make build` | ビルド（single.mkをincludeしているパッケージ） |
| `make fmt` | `gofmt` でフォーマット |
| `make goimports` | `gosimports`でimport整理 |
| `make lint` | Go / 文章のlint |
| `make set-license` | ライセンスヘッダーを追加・更新 |
| `make go-licenses-check` | 依存ライブラリのライセンスをチェック |
| `make vulncheck` | `govulncheck` で既知の脆弱性をチェック |

## ライセンス

`sacloud/sacloud-sdk-go` Copyright (C) 2026- [The sacloud/sacloud-sdk-go Authors](AUTHORS).

This project is published under [Apache 2.0 License](LICENSE).
