# デザイン

このライブラリは過去にいくつか作成されたさくらのクラウド向けHTTPクライアントを念頭に置きながら概念を整理して新たに作り直したものです。

# 大枠

- `Client` という単一の構造体を提供しており、すべての操作はこれに対する操作である
- `Client` 構造体を作成する方法として以下がある
    - `~/.usacloud` というディレクトリから設定を読み込むという方法
    - 環境変数を使って設定をする方法
    - コマンドライン引数を使って設定をする方法
- `Client` 構造体には機能として以下の機能がある
    - HTTPRequestDoerを実装
    - 設定に応じてリクエストの際に認証を実施
    - レスポンスに応じて自動でリトライを試みる
    - クライアントサイドでレートリミットを実施
    - トレース出力機能(デバッグで便利)
    - モック機能(テストで便利)

# API

- いろいろ考えた結果この`Client`は一気呵成に生成することはできず、セットアップが必要
- 基本的な流れとしては、

  ```golang
  import "github.com/sacloud/saclient-go"
  var GlobalClient saclient.Client
  ```

  をどこかで作っておき、

  ```golang
  GlobalClient.SetEnviron(os.Envirion())
  GlobalClient.FlagSet().Parse(os.Args)
  GlobalClient.SettingsFromTerraformProvider(&model)
  ```

  のうちの任意の組み合わせで設定を読み込んでから、

  ```golang
  GlobalClient.Populate()
  ```

- このようにして生成した`Client`構造体には以下のメソッドを提供する。

    ```golang
    func (*Client) Do (req *http.Request) (*http.Response, error)
    ```

  このメソッドにより各種SDKとの連携を可能にする。

# 複数インスタンス

- 基本的に環境変数などはプロセスグローバルな資源であるからクライアントもプロセスグローバルであるのが自然
- しかしながらAPIエンドポイントごとに異なる認証方式などある現状、複数個のクライアントを別々に持てた方が便利
- そこで以下のようにする

  - まずプロセスグローバルなテンプレートを作成

      ```golang
      import "github.com/sacloud/saclient-go"
      var TemplateClient saclient.Client
      ```

  - 通常通り環境変数などから初期化(上記参照)
  - これをコピーして必要な変更を加えつつ新しいクライアントを作成

      ```golang
      var basicCient saclient.Client = TemplateClient.DupWith(WithFavouringRFC7523())
      var bearerClient saclient.Client = TemplateClient.DupWith(WithBearerToken(tok))
      ```

  - これらクライアントをSDKに渡す

# エラー

- 内部で生じたエラーは構造体`Error`にて表現
- このライブラリを設計した時点でGoのエラー処理は混沌としており正解がない(新しいエラーライブラリを作りました、みたいなテックブログが各社に一個ずつある状態)
- しかもあらゆる著名なエラーライブラリが更新止まっている
- オープンソースライブラリであることを鑑みるにあまり強い思想で作るのは違うと思った
- golang stdlib だけに依存することにする
- スタックトレースなどは取れないがしょうがない

# 設定ファイル

- 設定ファイルの場所は
  - 環境変数`SAKURACLOUD_PROFILE_DIR`があればそれを採用
  - 環境変数`USACLOUD_PROFILE_DIR`があればそれを採用
  - `XDG_CONFIG_HOME`配下に`usacloud`があればそれを採用
  - どれもなければ`~/.usacloud`を採用

- 設定ファイルにはAPIアクセスと直接は関係しないusacloud固有の設定なども含めることが過去から可能。
  このため設定ファイルの内容に関しては「JSONとして読み書きできる」以上のことは求めないものとする

# 認証

- 以下の3種の認証をサポートする
  - さくらのクラウド ホームから「APIキー」機能で作成したAPIキー
  - さくらのクラウド ホームから「サービスプリンシパル」に登録した公開鍵を利用するパターン
  - さくらのクラウド「シンプルMQ」でキューを作成すると付随してくるトークン

## APIキーの場合

- 以下のいずれかの方法で指定(以下の順で探す)
  - 環境変数 `SAKURACLOUD_ACCESS_TOKEN` / `SAKURACLOUD_ACCESS_TOKEN_SECRET` で指定
  - それがなければ、設定ファイルに記載の`AccessToken` / `AccessTokenSecret`で指定
  - それもなければ、コマンドラインから `--token` / `--secret` で指定
  - それもなければterraform provider block内に記載の `AccessToken` / `AccessToken`

## 公開鍵の場合

- 前提としてさくらのクラウドに登録した公開鍵のペアになる秘密鍵をローカルファイルに保存する
  - 現状PEMのみのサポート、将来拡充するかも
- その上で以下のいずれかの方法で指定(以下の順で探す)
  - 環境変数 `SAKURACLOUD_PRIVATE_KEY_PATH` で指定
  - それがなければ、設定ファイルに記載の`PrivateKeyPEMPath`で指定
  - それもなければ、コマンドラインから `--token` / `--secret` で指定
  - それもなければterraform provider block内に記載の `PrivateKeyPEMPath`
- CI環境などローカルファイルに書き出すことが«手間、もしくはセキュリティ上の懸念»により難しい場合は環境変数`SAKURACLOUD_PRIVATE_KEY`に生のPEMを指定することも可能

## MQ のトークンの場合

- 設定ファイルや環境変数で指定する方法はない
- `client = client.DupWith(WithBearerToken(key))` などとしてプログラム側から指定