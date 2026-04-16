# api-client-goでのHTTP/APIクライアントの共通化

- URL: https://github.com/sacloud/api-client-go/pull/8
- Author: @yamamoto-febc

## 概要

IaaS/オブジェクトストレージ/PHYでHTTP/APIクライアント周りの処理を共通化したい。

libsacloudでは`helper/api`や`sacloud/profile`パッケージを用いて`sacloud.Client`を組み立てる実装となっている。  
HTTPクライアントレベルでは`sacloud/go-http`で共通化されているが、環境変数からのAPIキーの読み込みやプロファイルのサポートなどは
libsacloudで実装されておりPHYやオブジェクトストレージで利用するのが困難。  

このためapi-client-goでこれらを実装し共通で利用可能にする。

## 実装

api-client-go側で各種設定値を保持し、HTTPクライアントとして`HttpRequestDoer`を返すところまでを実装する。  
各クライアント側プロジェクトから参照されることによる循環参照を防ぐためにapi-client-go配下で独立したモジュールとして実装する。

### 主要コンポーネント

#### `Options`

APIトークンやシークレットなどのAPI/HTTPクライアントの動作に必要な値やオプション値を保持する。

```go
// Options sacloudhttp.Clientを作成する際のオプション
type Options struct {
    // AccessToken APIキー:トークン
    AccessToken string
    // AccessTokenSecret APIキー:シークレット
    AccessTokenSecret string

    // ...
}
```

プロファイルからOptionsを生成した場合はOptions内部にprofile.ConfigValueを保持する。  
clientパッケージで利用されないルートURLなどの設定は各アプリケーション側でここから参照して利用する。  

#### `Factory`

Optionsを保持し、HttpRequestDoerを生成する。
内部的にhttp.Clientを保持し、Transportのカスタマイズなどを担当する。  
HttpRequestDoerの実装としてsacloud/go-httpを用いる。

```go
// Factory client.HttpRequestDoerを作成して返すファクトリー
type Factory struct {
	options    *Options
	httpClient *http.Client // Transportの初期化を1度だけ行うためにoptionsのHttpClientの参照をここにコピーして保持しておく

	once sync.Once
}
```

#### HttpRequestDoer

HTTPリクエストを担当するインターフェース

```go
// HttpRequestDoer API/HTTPクライアントインターフェース
type HttpRequestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}
```

### 利用イメージ

```go
	// 環境変数/プロファイルを読み込んでオプションを組み立てる
	opt, _ := client.DefaultOption()
	// オプションからファクトリー生成
	clientFactory := client.NewFactory(opt)
	// ファクトリーからHttpRequestDoerを生成
	doer := clientFactory.NewHttpRequestDoer()

	// doerを用いてHTTPリクエスト
    resp, err := doer.Do(http.NewRequest("GET", url, nil))
	// ...
```

## 改訂履歴

- 2022/3/8: 初版作成
- 2022/3/11: sacloud-goから切り出し