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
  import saht "github.com/sacloud/http-client-go"
  var GlobalClient saht.Client
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
  