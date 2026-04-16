# sacloud-sdk-go

[![Go Version](https://img.shields.io/badge/go-1.25+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](LICENSE)

さくらのクラウド Go SDK

## 概要

sacloud-sdk-goは、[さくらのクラウド](https://cloud.sakura.ad.jp/)の各種サービスをGo言語から利用するためのSDK群をまとめたモノレポジトリです。

このリポジトリには、さくらのクラウドの各サービスに対応したAPIクライアントライブラリと、それらを利用した高レベルなサービスライブラリが含まれています。

> [!WARNING]
> このレポジトリは旧URLからの移行作業中です。この記述がある間はまだ実運用は推奨されません。

> [!WARNING]
> 本プロジェクトの一部パッケージはv1.0に達するまでの間、互換性のない形で変更される可能性があります。

## 構成

### APIクライアントライブラリ (`api/`)

さくらのクラウド各サービスのREST APIを操作するための低レベルクライアントライブラリです。

| パッケージ | 説明 | Goパッケージ |
|-----------|------|-------------|
| `api/iaas` | IaaS API (サーバー、ディスク、ネットワークなど) | `github.com/sacloud/iaas-api-go` |
| `api/webaccel` | ウェブアクセラレータ API | `github.com/sacloud/webaccel-api-go` |
| `api/iam` | IAM API (認証・認可) | `github.com/sacloud/iam-api-go` |
| `api/addon` | Add-on API (データレーク、WAFなど) | `github.com/sacloud/addon-api-go` |
| `api/apigw` | API Gateway API | `github.com/sacloud/apigw-api-go` |
| `api/apprun` | AppRun API | `github.com/sacloud/apprun-api-go` |
| `api/apprun-dedicated` | AppRun Dedicated API | `github.com/sacloud/apprun-dedicated-api-go` |
| `api/cloudhsm` | CloudHSM API | `github.com/sacloud/cloudhsm-api-go` |
| `api/dedicated-storage` | 専用ストレージ API | `github.com/sacloud/dedicated-storage-api-go` |
| `api/eventbus` | EventBus API | `github.com/sacloud/eventbus-api-go` |
| `api/kms` | Key Management Service API | `github.com/sacloud/kms-api-go` |
| `api/monitoring-suite` | 監視サービス API | `github.com/sacloud/monitoring-suite-api-go` |
| `api/nosql` | NoSQL API | `github.com/sacloud/nosql-api-go` |
| `api/object-storage` | オブジェクトストレージ API | `github.com/sacloud/object-storage-api-go` |
| `api/secretmanager` | Secret Manager API | `github.com/sacloud/secretmanager-api-go` |
| `api/security-control` | セキュリティコントロール API | `github.com/sacloud/security-control-api-go` |
| `api/service-endpoint-gateway` | Service Endpoint Gateway API | `github.com/sacloud/service-endpoint-gateway-api-go` |
| `api/simple-notification` | シンプル通知 API | `github.com/sacloud/simple-notification-api-go` |
| `api/simplemq` | SimpleMQ API | `github.com/sacloud/simplemq-api-go` |
| `api/workflows` | Workflows API | `github.com/sacloud/workflows-api-go` |

### 高レベルサービスライブラリ (`service/`)

APIクライアントライブラリをラップし、より使いやすいインターフェースを提供する高レベルライブラリです。

| パッケージ | 説明 | Goパッケージ |
|-----------|------|-------------|
| `service/iaas` | IaaS向け高レベルAPI | `github.com/sacloud/iaas-service-go` |
| `service/webaccel` | ウェブアクセラレータ高レベルAPI | `github.com/sacloud/webaccel-service-go` |

### 内部パッケージ (`internal/`)

SDKの内部で使用される共有パッケージです。

| パッケージ | 説明 |
|-----------|------|
| `internal/api-client` | APIクライアントの共通実装 |
| `internal/go-http` | HTTP通信の共通実装 |
| `internal/packages` | 汎用パッケージ群 |
| `internal/saclient` | クライアント認証・設定の共通実装 |
| `internal/services` | サービスレイヤーの共通実装 |

## ワークスペース構成

このリポジトリはGoワークスペース(`go.work`)を使用して管理されています。

```
go.work
├── api/               # APIクライアントライブラリ
├── service/           # 高レベルサービスライブラリ
├── internal/          # 内部共有パッケージ
└── makefiles/         # 共有Makefile
```

## インストール

各パッケージは個別にインストールできます。

### IaaS APIクライアントの場合

```bash
go get github.com/sacloud/iaas-api-go
```

### ウェブアクセラレータ APIクライアントの場合

```bash
go get github.com/sacloud/webaccel-api-go
```

### IAM APIクライアントの場合

```bash
go get github.com/sacloud/iam-api-go
```

## 使用方法

基本的な使用方法は各パッケージのREADME.mdを参照してください。

### 簡単な使用例 (IaaS)

```go
package main

import (
    "context"
    "log"

    "github.com/sacloud/iaas-api-go"
    "github.com/sacloud/saclient-go"
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

APIを利用するには、さくらのクラウドのアクセストークンが必要です。

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

## ライセンス

`sacloud/sacloud-sdk-go` Copyright (C) 2026- [The sacloud/sacloud-sdk-go Authors](AUTHORS).

This project is published under [Apache 2.0 License](LICENSE).
