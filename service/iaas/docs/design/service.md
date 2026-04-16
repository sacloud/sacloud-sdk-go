# iaas-service-go/service

- URL: https://github.com/sacloud/iaas-service-go/pull/21
- Author: @yamamoto-febc

## 概要

libsacloudのhelper/serviceパッケージ+helper/builderを母体とした、
各種クライアントから利用可能な操作をまとめた高レベル操作を実装する`service`パッケージを提供する。

### `service`パッケージ

クライアントから利用可能な操作をまとめ、統一的なインターフェースを提供する。  
serviceパッケージには以下のような役割を持たせる。

- 複雑なAPI呼び出しを隠蔽するシンプルなインターフェース
- 操作対象のリソースが異なっても同じ操作感でCRUD+L操作ができる仕組み
- コード生成を念頭に置いたメタデータ

#### 複雑なAPI呼び出しを隠蔽し再利用可能な形で公開

例えばサーバを作成したい場合、実際には以下のようなAPIを用いる必要がある。

- サーバプランの検索/参照
- サーバ作成
- ディスクプランの検索/参照
- ソースアーカイブ/ソースディスクの検索/参照
- ディスク作成
- ディスク作成完了まで待機
- ディスクの修正
- サーバとディスクの接続
- サーバへのNIC追加
- サーバのNICへのパケットフィルタ接続
- サーバのNICをスイッチに接続
- サーバの電源投入(+cloud-init)

これらの一連の操作を再利用可能な形で実装/公開する役割をserviceパッケージに持たせる。

#### 統一的にCRUD+L操作ができるインターフェース

例えばIaaSでリソースを検索する場合、リソースごとに以下のような点が異なる。

- ゾーンの指定の有無(グローバルリソースか？)
- 検索対象にできるフィールド

リソースごとにfuncの引数が異なるとクライアント側のコード生成が煩雑になる。  
このため、リソースが違っても統一的にCRUD+L操作が行えるインターフェースを提供する。

例:

```go
/* 
 *  Note:実装イメージであり実際のコードとは異なります 
 */

package xxx // リソース種別ごとにパッケージを切る

type Service interface {
	// 基本的なCRUD+L操作
    Create(context.Context, *CreateRequest) (interface{}, error)
    Update(context.Context, *UpdateRequest) (interface{}, error)
    Read(context.Context, *ReadRequest) (interface{}, error)
    Delete(context.Context, *DeleteRequest) error
	List(context.Context, *ListRequest) ([]interface{}, error)
	
	// さらにリソースごとに個別操作を持つ場合もある
}
```

- リソースごとにパッケージを切る
- xxxRequest(xxxには操作名)という構造体でゾーン指定有無などのリソースごとの違いを吸収

#### コード生成を念頭においたメタデータの提供

TODO 取り組むときに追記する

## やること/やらないこと

### やること

TODO: 必要に応じて追記

### やらないこと

- メタデータは次回以降のアップデート時に実装し、libsacloudからの移植時には提供しない
 
TODO: 必要に応じて追記

## 実装

- libsacloudからhelper/serviceを移植

## 改訂履歴

- 2022/3/15: 初版作成