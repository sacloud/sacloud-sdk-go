# Agent Instructions for `api/`

このドキュメントは `api/` 配下の各 API クライアントパッケージが従うべき**全ての規約**を定める。ディレクトリ構造、Go パッケージ命名、OpenAPI / ogen 設定、クライアント・エラー・バージョン定義、Op レイヤー命名、テスト、Makefile、ライセンスヘッダまでを対象とする。新規パッケージを追加する際、および既存パッケージをリファクタリングする際はここを参照すること。

プロジェクト全体の構成は [`../AGENTS.md`](../AGENTS.md) を参照。本ドキュメントは `api/` レイヤー固有の規約に特化する。

## 目次

1. [対象と除外](#対象と除外)
2. [背景 — なぜ Op レイヤーが存在するか](#背景--なぜ-op-レイヤーが存在するか)
3. [パッケージ構成](#パッケージ構成)
4. [Go モジュール](#go-モジュール)
5. [Go パッケージ命名](#go-パッケージ命名)
6. [ライセンスヘッダ](#ライセンスヘッダ)
7. [コメント言語](#コメント言語)
8. [OpenAPI スキーマ (`openapi/`)](#openapi-スキーマ-openapi)
9. [ogen 設定 (`ogen-config.yaml`)](#ogen-設定-ogen-configyaml)
10. [生成コード (`apis/v1/`)](#生成コード-apisv1)
11. [`client.go`](#clientgo)
12. [`version.go`](#versiongo)
13. [`error.go`](#errorgo)
14. [Op レイヤー命名規則](#op-レイヤー命名規則)
15. [エラーハンドリング詳細](#エラーハンドリング詳細)
16. [ファイル構成](#ファイル構成)
17. [集約パッケージ (iam / apprun-dedicated)](#集約パッケージ-iam--apprun-dedicated)
18. [テスト規約](#テスト規約)
19. [Makefile と `includes/`](#makefile-と-includes)
20. [ドキュメント (`README.md`, `CHANGELOG.md`, `AUTHORS`, `LICENSE`)](#ドキュメント-readmemd-changelogmd-authors-license)
21. [正準例](#正準例)
22. [現状の逸脱 (リファクタリング TODO)](#現状の逸脱-リファクタリング-todo)
23. [新規パッケージ追加チェックリスト](#新規パッケージ追加チェックリスト)

---

## 対象と除外

**対象 (17 パッケージ、OpenAPI + ogen ベース):**

`addon`, `apigw`, `apprun`, `apprun-dedicated`, `cloudhsm`, `dedicated-storage`, `eventbus`, `iam`, `kms`, `monitoring-suite`, `nosql`, `object-storage`, `secretmanager`, `security-control`, `service-endpoint-gateway`, `simple-notification`, `simplemq`, `workflows`

**対象外:**

- `api/iaas` — OpenAPI 仕様を持たず、手書き API クライアント。
- `api/webaccel` — 同じく手書き。`openapi/` も `apis/v1/` も持たない。

除外パッケージは独自の規約で書かれており、本ドキュメントの規約を強制しない。

## 背景 — なぜ Op レイヤーが存在するか

さくらのクラウドの OpenAPI 仕様は、サービスごとに operationId の命名がバラバラで、生成される関数名も統一されていない。たとえば同じ「一覧取得」でも下記のように揺れている。

| サービス | ogen 生成メソッド名 |
| --- | --- |
| `api/apigw` | `GetCertificates` |
| `api/secretmanager` | `SecretmanagerVaultsSecretsList` |
| `api/addon` | `ListAZ0501` |
| `api/iam` | `OrganizationPasswordPolicyGet` |

これを直接叩かせると SDK 利用者の開発体験が崩壊するため、各パッケージは `apis/v1/` 配下の ogen 生成コードの上に**手書きの薄いラッパ層 (= 「Op レイヤー」)** を置き、`List` / `Read` / `Create` / `Update` / `Delete` といった統一的な動詞で呼び出せるようにしている。

## パッケージ構成

```
api/<service>/
├── openapi/                 OpenAPI スキーマ (JSON / YAML)
├── apis/                    ogen 生成コード (手を入れない)
│   └── v1/
│       └── oas_*.go
├── <resource>.go            Op レイヤー (本ドキュメントの対象)
├── <resource>_test.go       Op レイヤーのテスト
├── client.go                NewClient コンストラクタと UserAgent / エンドポイント定数
├── error.go                 NewError / NewAPIError とパッケージ固有 Error 型
├── version.go               Version 定数
├── ogen-config.yaml         ogen 設定 (アップストリーム管理の場合は省略可)
├── Makefile                 includes/ を取り込むだけの薄いラッパ
├── includes/go/common.mk    モノレポ共通 Make レシピ (symlink 相当)
├── includes/go/single.mk
├── README.md                日本語概要 / 使い方
├── CHANGELOG.md             リリース履歴
├── LICENSE                  Apache 2.0 全文
├── AUTHORS                  著作権表示
├── go.mod / go.sum
```

集約パッケージ (サブリソースが多い場合) は `apis/<subresource>/<subresource>.go` にサブパッケージを切り、ルート `api.go` で再エクスポートする。後述「[集約パッケージ](#集約パッケージ-iam--apprun-dedicated)」参照。

## Go モジュール

各パッケージは独立した Go モジュール。モノレポだが `go.work` で束ねる構成。

- **モジュールパス**: `github.com/sacloud/<service>-api-go`
  - 例: `github.com/sacloud/apigw-api-go`, `github.com/sacloud/secretmanager-api-go`
- **Go バージョン**: `go 1.25.x` (最新安定版に追従)
- **toolchain**: 具体版でピン留め (例: `toolchain go1.25.7`)
- **ツール宣言**: `tool github.com/ogen-go/ogen/cmd/ogen`
- **直接依存 (共通)**:
  - `github.com/sacloud/saclient-go` — 認証 / エンドポイント解決 / HTTP クライアント augmentation
  - `github.com/sacloud/packages-go` — テストユーティリティ共通実装 (利用する場合)
  - `github.com/ogen-go/ogen` — ogen ランタイム
  - `github.com/go-faster/errors`, `github.com/go-faster/jx` — ogen 生成コードの依存
  - `github.com/stretchr/testify` — テスト

- **internal/ への依存は禁止**: `api/*` は公開モジュールなので、モノレポ内の `internal/` パッケージから import してはならない。内部パッケージの成果は `github.com/sacloud/saclient-go` などの公開モジュールとして公開される。

## Go パッケージ命名

ディレクトリ名にハイフンが含まれる場合は**ハイフンを除去して全小文字**とする。名前が長くなりすぎる場合は略称を採用してよい。

| ディレクトリ | Go パッケージ名 |
| --- | --- |
| `addon` | `addon` |
| `apigw` | `apigw` |
| `apprun` | `apprun` |
| `cloudhsm` | `cloudhsm` |
| `eventbus` | `eventbus` |
| `iam` | `iam` |
| `kms` | `kms` |
| `nosql` | `nosql` |
| `secretmanager` | `secretmanager` |
| `simplemq` | `simplemq` |
| `workflows` | `workflows` |
| `dedicated-storage` | `dedicatedstorage` |
| `monitoring-suite` | `monitoringsuite` |
| `object-storage` | `objectstorage` |
| `security-control` | `securitycontrol` |
| `simple-notification` | `simplenotification` |
| `service-endpoint-gateway` | `seg` (略称) |

**既存逸脱 (新規では真似しない):**

- `api/apprun-dedicated` → `package apprun_dedicated` (アンダースコア入り)。Go のスタイルガイド上はハイフン除去が望ましい。

## ライセンスヘッダ

すべての `.go` ファイル (ogen 生成物を除く) には以下のいずれかのヘッダを先頭に置く。現行は 2 形式が混在している。

### 形式 A (Apache 2.0 全文ボイラープレート) — 多数派

```go
// Copyright <年> The sacloud/<service>-api-go authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
```

### 形式 B (SPDX 短縮形)

```go
// Copyright <年> The sacloud/<service>-api-go authors
// SPDX-License-Identifier: Apache-2.0
```

### 正規

**新規ファイルは形式 B (SPDX) を推奨**。短く、ツール連携 (REUSE, scancode, addlicense 等) も容易。形式 A を使う既存パッケージが大半だが、モノレポ統合後に順次 SPDX 化していく方針。

年次表記は `YYYY-` (起点年のみ、継続中) または `YYYY-YYYY` (区間)。`Makefile` の `COPYRIGHT_YEAR` と一致させる。

## コメント言語

- **godoc (exported シンボルの doc コメント)**: 原則**日本語**。`DefaultAPIRootURL デフォルトのAPIルートURL` のように、シンボル名の後ろに日本語説明を書く Go 慣習を踏襲。
- **非公開の実装コメント / TODO / NOTE**: 日本語・英語いずれも可。読み手 (日本語圏エンジニア) を意識。
- **README.md**: 日本語。
- **コミットメッセージ / PR タイトル**: 英語または日本語、既存 git log に合わせる。
- **コード内の識別子**: 英語 (Go の一般慣習通り)。

## OpenAPI スキーマ (`openapi/`)

- 配置: `api/<service>/openapi/openapi.json` または `openapi.yaml`。
- `.json` と `.yaml` は好み。新規は YAML 推奨 (差分が読みやすい)。
- 複数仕様を持つ例: `api/simplemq/openapi/queue.yaml` と `api/simplemq/openapi/message.yaml` (Queue API と Message API が別々)。生成先も `apis/v1/queue` / `apis/v1/message` に分かれる。

## ogen 設定 (`ogen-config.yaml`)

ファイル名は `ogen-config.yaml` (拡張子 `.yaml`) を正規とする。`.yml` の既存パッケージ (`addon`, `apprun-dedicated` ほか) は移行対象。

標準的な内容:

```yaml
generator:
  features:
    enable:
      - 'paths/client'
      - 'client/request/validation'
      - 'debug/example_tests'
    disable_all: true
```

`disable_all: true` の上で必要な機能だけ `enable` する方針。コメントで「ogen 未対応のため OpenAPI 側で回避した」箇所を明記する。参考: [`apigw/ogen-config.yaml`](apigw/ogen-config.yaml)。

## 生成コード (`apis/v1/`)

- **直接編集しない**。変更は `openapi/` を編集し、`make gen` (`go tool ogen` を呼ぶ) で再生成する。
- `simplemq` のように仕様が分かれる場合は `apis/v1/<name>/` にサブディレクトリを切る。

## `client.go`

各パッケージは必ず `client.go` に以下を定義する。

### 定数

```go
const (
    // DefaultAPIRootURL デフォルトのAPIルートURL
    DefaultAPIRootURL = "https://secure.sakura.ad.jp/cloud/api/<path>/1.0/"

    // ServiceKey SDKの種別を示すキー、プロファイルでのエンドポイント取得に利用する
    ServiceKey = "<service>"
)
```

- `DefaultAPIRootURL`: さくらのクラウドの本番エンドポイント。
- `ServiceKey`: `saclient.ClientAPI.EndpointConfig()` の `Endpoints[ServiceKey]` で上書き可能にするためのキー。
  - 命名は snake_case (例: `"apigw"`, `"secretmanager"`, `"simple_mq_queue"`)。

### UserAgent

```go
// UserAgent APIリクエスト時のユーザーエージェント
var UserAgent = fmt.Sprintf(
    "<service>-api-go/%s (%s/%s; +https://github.com/sacloud/<service>-api-go)",
    Version,
    runtime.GOOS,
    runtime.GOARCH,
)
```

### コンストラクタ

```go
func NewClient(client saclient.ClientAPI) (*v1.Client, error) {
    endpointConfig, err := client.EndpointConfig()
    if err != nil {
        return nil, NewError("unable to load endpoint configuration", err)
    }
    endpoint := DefaultAPIRootURL
    if ep, ok := endpointConfig.Endpoints[ServiceKey]; ok && ep != "" {
        endpoint = ep
    }
    return NewClientWithAPIRootURL(client, endpoint)
}

func NewClientWithAPIRootURL(client saclient.ClientAPI, apiRootURL string) (*v1.Client, error) {
    dupable, ok := client.(saclient.ClientOptionAPI)
    if !ok {
        return nil, NewError("client does not implement saclient.ClientOptionAPI", nil)
    }
    augmented, err := dupable.DupWith(
        saclient.WithUserAgent(UserAgent),
    )
    if err != nil {
        return nil, err
    }
    return v1.NewClient(apiRootURL, v1.WithClient(augmented))
}
```

正準例: [`apigw/client.go`](apigw/client.go)。

**複数エンドポイントを持つ場合** (例: Queue と Message を両方扱う `simplemq`) は、`NewQueueClient` / `NewMessageClient` のようにリソースごとに関数を用意し、`DefaultXxxAPIRootURL` / `ServiceKeyXxx` を個別定義する。

## `version.go`

```go
package <service>

const Version = "X.Y.Z"
```

- 単一の `Version` 定数のみを置く。他のコードは書かない。
- `client.go` の `UserAgent` 生成でこれを参照する。
- リリース時は `CHANGELOG.md` と併せて更新。

## `error.go`

パッケージ固有の `Error` 型と 2 つのヘルパーを必ず定義する。

```go
type Error struct {
    msg string
    err error
}

func (e *Error) Error() string {
    if e.msg != "" {
        if e.err != nil {
            return "<service>: " + e.msg + ": " + e.err.Error()
        }
        return "<service>: " + e.msg
    }
    return "<service>: " + e.err.Error()
}

func (e *Error) Unwrap() error {
    return e.err
}

func NewError(msg string, err error) *Error {
    return &Error{msg: msg, err: err}
}

func NewAPIError(method string, code int, err error) *Error {
    return &Error{msg: method, err: saclient.NewError(code, "", err)}
}
```

- パッケージ名プレフィックス (`"<service>: "`) はパッケージ名と一致させる。
- `Unwrap()` を実装することで `errors.As` / `errors.Is` による判定を可能にする。
- 正準例: [`apigw/error.go`](apigw/error.go), [`simplemq/error.go`](simplemq/error.go)。

## Op レイヤー命名規則

新規パッケージ追加時の標準を以下に示す。

### 1. インターフェース名

- 形式: `XxxAPI` (単数形 / PascalCase)
- 例: `CertificateAPI`, `QueueAPI`, `SecretAPI`, `ContractAPI`
- サービス名を冠さない (`ApigwCertificateAPI` とはしない)。

### 2. 実装構造体名

- 形式: `xxxOp` (先頭小文字 / `Op` サフィックス)
- 例: `certificateOp`, `queueOp`, `secretOp`, `contractOp`
- **exported にしない**。コンストラクタ経由でのみ生成させる。
- 直後に `var _ XxxAPI = (*xxxOp)(nil)` を置く。

### 3. コンストラクタ

- 形式: `NewXxxOp(client *v1.Client) XxxAPI`
- 返り値はインターフェース型。`*xxxOp` を直接返さない。
- 追加の依存 (例: `vaultId` や `zone`) は第 2 引数以降で受ける。
- 初期化に失敗しうる場合のみ `(XxxAPI, error)` を返してよい (例: `cloudhsm.NewClientOp`)。

### 4. クライアントフィールド

- 形式: `client *v1.Client` (名前付きフィールド)
- **埋め込みにしない** (`struct{ *v1.Client }` は不可)。
- 追加フィールドは名前付きで並べる。

```go
type secretOp struct {
    client  *v1.Client
    vaultId string
}
```

### 5. レシーバ名

- 形式: `op *xxxOp`
- 1 文字の `o` や `s` は使わない (grep しやすさのため)。

### 6. メソッド名 (CRUD)

| 動詞 | 用途 |
| --- | --- |
| `List` | 複数件取得 |
| `Read` | 単一件取得 (`Get` は使わない) |
| `Create` | 作成 |
| `Update` | 更新 |
| `Delete` | 削除 |

ドメイン固有の操作は動詞始まりで付け足す。横断で広く使われているもの:

- `Status` / `ReadStatus` (状態取得)
- `PowerOn` / `Shutdown` / `Reset` (電源操作)
- `Rotate` / `RotateAPIKey` (キーローテーション)
- `Cancel` (ジョブキャンセル)
- `Apply` (設定反映)

同リソースの複数バリエーションは `List<サブリソース>` / `Read<サブリソース>` で合成 (例: `ListDiskSnapshots`, `ReadInterface`)。

### 7. メソッドシグネチャ

- 第 1 引数は必ず `ctx context.Context`。
- ID は `Read` / `Update` / `Delete` で第 2 引数。
- リクエストボディ (POST / PUT) は ID より後ろ。
- 返り値は `(*v1.X, error)` または `([]v1.X, error)`。値型で返さない。

```go
Read(ctx context.Context, id string) (*v1.Secret, error)
Update(ctx context.Context, id string, req v1.UpdateSecretRequest) (*v1.Secret, error)
Delete(ctx context.Context, id string) error
```

### 8. インターフェース assertion

```go
var _ XxxAPI = (*xxxOp)(nil)
```

構造体定義の直後に置く。

## エラーハンドリング詳細

Op メソッドは ogen 生成レスポンスを型スイッチし、HTTP ステータスを拾って `NewAPIError("<Resource>.<Method>", status, err)` で包む。

```go
func (op *certificateOp) List(ctx context.Context) ([]v1.Certificate, error) {
    res, err := op.client.GetCertificates(ctx)
    if err != nil {
        return nil, NewAPIError("Certificate.List", 0, err)
    }
    switch p := res.(type) {
    case *v1.GetCertificatesOK:
        return p.Apigw.Certificates, nil
    case *v1.GetCertificatesBadRequest:
        return nil, NewAPIError("Certificate.List", 400, errors.New(p.Message.Value))
    // ...
    }
    return nil, NewAPIError("Certificate.List", 0, nil)
}
```

- `method` 部分は `"<Resource>.<Method>"` 形式 (例: `"Queue.List"`, `"Certificate.Create"`)。
- トランスポートエラー (ogen から直接 error が返る) はステータス `0`。
- 型スイッチの default 節は必ず書き、予期しないレスポンスでも `NewAPIError(..., 0, nil)` を返す。
- ジェネリクス版ヘルパー (`ErrorFromDecodedResponse[T, E]`, `createAPIError`) を定義している既存パッケージもあるが、**新規は採用しない**。統一方針未決。

## ファイル構成

- **1 リソース 1 ファイル**。
- ファイル名は**単数形**を用いる (`certificate.go`, `queue.go`, `secret.go`)。
- テストは `<resource>_test.go` として同じ階層に置く。
- サブパッケージ化する場合は `apis/<subresource>/<subresource>.go`。

## 集約パッケージ (iam / apprun-dedicated)

サブリソースが多く単一ファイルに収まらない API は以下のパターンを用いる。

### ルート `api.go`

```go
type AuthAPI = auth.AuthAPI
type FolderAPI = folder.FolderAPI
// ...

var NewAuthOp = auth.NewAuthOp
var NewFolderOp = folder.NewFolderOp
// ...
```

### サブパッケージ `apis/auth/auth.go`

```go
package auth

type AuthAPI interface { /* ... */ }
type authOp struct { client *v1.Client }
func NewAuthOp(client *v1.Client) AuthAPI { return &authOp{client} }
```

参考: [`iam/api.go`](iam/api.go), [`iam/apis/auth/auth.go`](iam/apis/auth/auth.go)。

## テスト規約

### ユニットテスト

- `<resource>_test.go` として Op レイヤーと同じディレクトリに配置。
- パッケージ名は `<service>` (internal) または `<service>_test` (external) どちらも可。外部テストの方がパブリック API の契約を検証できるため推奨。

### 受け入れテスト (acceptance)

- ファイル名: `acceptance_test.go` または `<resource>_acceptance_test.go`。
- ビルドタグ: `//go:build acctest` (スペースなしの新形式、および `// +build acctest` の旧形式を両方書く) 。

  ```go
  //go:build acctest
  // +build acctest

  package apprun_test
  ```

- 実 API を叩くため、環境変数でゲート:
  - `skipIfNoAPIKey(t)` のような helper を用意して API キー未設定時は `t.Skip`。
- 実行: `make testacc` (`TESTACC=1 go test ... --tags=acctest`)。

### テストユーティリティ

- 共通フィクスチャは `helper_test.go` に集約。
- `newTestClient()` のような factory を提供。
- `testdata/` ディレクトリに JSON / YAML のフィクスチャを置く。

### Example テスト

- `example_test.go` で godoc 用の Example 関数を提供可能。強制ではない。

## Makefile と `includes/`

各パッケージの `Makefile` は以下の薄い形に揃える。

```makefile
#====================
AUTHOR         ?= The sacloud/<service>-api-go Authors
COPYRIGHT_YEAR ?= <year-range>

BIN            ?= <service>-api-go
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
include includes/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
```

- `includes/go/common.mk` / `single.mk` はモノレポの `makefiles/` からコピー or シンボリックリンクされた共通レシピ。
- 標準ターゲット:
  - `make test` — `go test ./... -v -timeout=120m -parallel=8 -race`
  - `make testacc` — `TESTACC=1 go test ... --tags=acctest`
  - `make gen` — ogen でコード再生成
  - `make fmt`, `make goimports`, `make lint` — フォーマット・静的解析
  - `make vulncheck` — `govulncheck`
  - `make dev-tools` — gosimports / addlicense / golangci-lint / actionlint / textlint のインストール

## ドキュメント (`README.md`, `CHANGELOG.md`, `AUTHORS`, `LICENSE`)

各パッケージ直下に以下の 4 ファイルを必ず置く。

- **`README.md`** — 日本語で概要、認証方法、使い方、コード生成手順を記述。
- **`CHANGELOG.md`** — セマンティックバージョニングでリリース履歴。
- **`AUTHORS`** — 著作権者表記。
- **`LICENSE`** — Apache 2.0 全文。モノレポ直下の `LICENSE` と内容一致。

## 正準例

| ケース | ファイル |
| --- | --- |
| `client.go` / `NewClient` | [`apigw/client.go`](apigw/client.go) |
| `version.go` | [`apigw/version.go`](apigw/version.go) |
| `error.go` | [`apigw/error.go`](apigw/error.go) |
| 基本 CRUD の Op | [`apigw/certificates.go`](apigw/certificates.go) |
| CRUD + ドメイン固有メソッド | [`simplemq/queue.go`](simplemq/queue.go) |
| 追加依存フィールド (`vaultId`) | [`secretmanager/secrets.go`](secretmanager/secrets.go) |
| 集約 + 再エクスポート | [`iam/api.go`](iam/api.go), [`iam/apis/auth/auth.go`](iam/apis/auth/auth.go) |
| 複数 OpenAPI 仕様 + 複数 Client | [`simplemq/client.go`](simplemq/client.go) |

---

## 現状の逸脱 (リファクタリング TODO)

以下は本ドキュメントの規約に準拠していない既存実装のリスト。**新規に真似しないこと**。破壊的変更を伴うものは個別タスクとして修正する。

### 構造体の公開性 / 命名

| パッケージ | 現状 | あるべき形 |
| --- | --- | --- |
| `api/addon` | `type AIOp struct{ *v1.Client }` — exported かつクライアント埋め込み | `type aiOp struct { client *v1.Client }` (非公開 + 名前付き) |
| `api/cloudhsm` | `type ClientOp struct { ... }` — exported | `type clientOp struct { ... }` (ただしリソース名 `Client` との衝突注意) |
| `api/service-endpoint-gateway` | `type ServiceEndpointGatewayOp struct { ... }` — exported | `type serviceEndpointGatewayOp struct { ... }` |
| `api/simple-notification` | `type DestinationOp struct { ... }` — exported | `type destinationOp struct { ... }` |

### レシーバ名

| パッケージ | 現状 | あるべき形 |
| --- | --- | --- |
| `api/simple-notification` | `func (o *DestinationOp) ...` | `func (op *destinationOp) ...` |

### ファイル名 (単数形への統一)

| パッケージ | 現状 | あるべき形 |
| --- | --- | --- |
| `api/apigw` | `certificates.go`, `domains.go`, `groups.go` ほか | `certificate.go`, `domain.go`, `group.go` |
| `api/apprun` | `applications.go`, `traffics.go`, `users.go`, `versions.go` | 単数形 |
| `api/eventbus` | `schedules.go` | `schedule.go` |
| `api/kms` | `keys.go` | `key.go` |
| `api/nosql` | `instances.go` | `instance.go` |
| `api/object-storage` | `accounts.go`, `buckets.go` ほか | 単数形 |
| `api/secretmanager` | `secrets.go`, `vaults.go` | `secret.go`, `vault.go` |

### Go パッケージ名

| パッケージ | 現状 | あるべき形 |
| --- | --- | --- |
| `api/apprun-dedicated` | `package apprun_dedicated` (アンダースコア) | `package apprundedicated` (ハイフン除去のみ) |

### ライセンスヘッダ形式の混在

| パッケージ | 現状 | あるべき形 |
| --- | --- | --- |
| 多数 (apigw, addon, iam ほか) | Apache 2.0 全文ボイラープレート | 新規は SPDX (`SPDX-License-Identifier: Apache-2.0`)。既存は段階的移行。 |

### ogen-config の拡張子

| パッケージ | 現状 | あるべき形 |
| --- | --- | --- |
| `api/addon`, `api/apprun-dedicated`, `api/iam` ほか | `ogen-config.yml` | `ogen-config.yaml` |

### エラーヘルパーの命名統一

| パッケージ | 現状ヘルパー名 | 備考 |
| --- | --- | --- |
| `api/addon` | `ErrorFromDecodedResponse` (exported, ジェネリクス) | exported である点が他と不整合 |
| `api/iam` (common) | `ErrorFromDecodedResponse` (exported, ジェネリクス) | 同上 |
| `api/monitoring-suite` | `errorFromDecodedResponse` (unexported) | 命名揺れ |
| `api/kms` | `createAPIError` | `NewAPIError` との混在 |
| `api/secretmanager` | `createAPIError` | 同上 |

統一方針は未決。本ドキュメントでの正規形は `NewAPIError(method, status, err)`。ジェネリクス版ヘルパー採用可否は別途議論。

### エラー呼び出しの一貫性

| パッケージ | 現状 | あるべき形 |
| --- | --- | --- |
| `api/simplemq` (`queue.go`) | `Delete` 内で `NewError("Delete", err)` と `NewAPIError("Queue.Delete", ...)` が混在 | `NewAPIError("Queue.Delete", ...)` に統一 |

### その他

| パッケージ | 状態 | 補足 |
| --- | --- | --- |
| `api/secretmanager` | `SecretAPI.Read` が未実装 (コメントアウト) | アップストリーム API 未対応。実装待ち。 |
| `api/nosql` | `GetVersion` / `GetParameters` / `SetParameters` を使用 | リソース特性上妥当な例外。`Read` 改名検討時はセマンティクスを再検討すること。 |

---

## 新規パッケージ追加チェックリスト

### ディレクトリ構造

- [ ] `api/<service>/` を作成。
- [ ] `openapi/` に OpenAPI スキーマを配置 (アップストリーム管理なら省略可)。
- [ ] `ogen-config.yaml` を作成 (拡張子は `.yaml`)。
- [ ] `make gen` で `apis/v1/` を生成。**手編集禁止**。

### Go モジュール

- [ ] モジュールパス `github.com/sacloud/<service>-api-go`。
- [ ] `go 1.25.x` 以上。toolchain をピン留め。
- [ ] ogen パッケージでは `tool github.com/ogen-go/ogen/cmd/ogen` を宣言。
- [ ] `internal/` への依存を持たない。

### Go パッケージ名

- [ ] ディレクトリ名からハイフンを除去した全小文字を使う (アンダースコア・略称・混合命名を避ける)。

### 必須ファイル

- [ ] `client.go` — `DefaultAPIRootURL`, `ServiceKey`, `UserAgent`, `NewClient`, `NewClientWithAPIRootURL`。
- [ ] `version.go` — `const Version = "X.Y.Z"` のみ。
- [ ] `error.go` — `Error` 型, `NewError`, `NewAPIError`。
- [ ] `README.md`, `CHANGELOG.md`, `AUTHORS`, `LICENSE`。
- [ ] `Makefile` — `includes/go/common.mk` と `includes/go/single.mk` を include。

### Op レイヤー

- [ ] リソースごとに `<resource>.go` (単数形ファイル名)。
- [ ] `XxxAPI` インターフェースと `xxxOp` (小文字) 構造体を定義。
- [ ] `var _ XxxAPI = (*xxxOp)(nil)` を構造体定義直後に置く。
- [ ] コンストラクタは `NewXxxOp(client *v1.Client) XxxAPI`。
- [ ] CRUD メソッドは `List` / `Read` / `Create` / `Update` / `Delete`。`Get` は使わない。
- [ ] 第 1 引数は `ctx context.Context`。
- [ ] 各 Op メソッドは `NewAPIError("<Resource>.<Method>", status, err)` でエラーを包む。

### ヘッダ・コメント

- [ ] 新規 `.go` ファイルは SPDX 形式ライセンスヘッダを付ける。
- [ ] exported シンボルには日本語 godoc を書く。

### テスト

- [ ] ユニットテストを `<resource>_test.go` に配置。
- [ ] 受け入れテストは `//go:build acctest` でゲート。
- [ ] 外部 API を叩く場合は API キー未設定時に `t.Skip`。

### 最終確認

- [ ] `make fmt` / `make goimports` / `make lint` が通る。
- [ ] `make test` が通る。
- [ ] `CHANGELOG.md` にエントリを追加。
- [ ] ルート `go.work` にモジュールが登録されている。
