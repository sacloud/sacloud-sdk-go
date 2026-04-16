# serviceインターフェース

- From: https://github.com/sacloud/iaas-service-go/pull/21
- Author: @yamamoto-febc

## 概要

iaas-service-goでのserviceへの取り組みを一歩進め、iaas以外に対してもserviceを拡張したい。  
拡張にあたり、serviceに求めるものや備えるべきインターフェース、共通で必要になるツールなどについてドキュメント化する。  

### serviceの目的

クライアントから利用可能な操作をまとめ、統一的なインターフェースを提供する。  
serviceパッケージには以下のような役割を持たせる。

- 複雑なAPI呼び出しを隠蔽するシンプルなインターフェース
- 操作対象のリソースが異なっても同じ操作感でCRUD+L操作ができる仕組み
- コード生成を念頭に置いたメタデータ(カタログ)の提供

### serviceが備えるべきインターフェース

serviceは以下の操作を持つ。

- Create: リソースを作成し、作成したリソースの情報を返す
- Read: Idを元にリソースを参照する
- Update: リソースを更新し、更新したリソースの情報を返す
- Delete: リソースを削除する
- List: リソース一覧を取得
- Action: 戻り値の不要な、CRUD+L以外の操作(電源操作など)

リソースごとに全てを備える必要はなく、任意で実装する。

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
	Boot(context.Context, *BootRequest) error
	Shutdown(context.Context, *ShutdownRequest) error
}
```

### serviceが提供するメタデータ(カタログ)

以下のようなメタデータを提供する。

- リソースの種類(プラットフォーム単位)
  - サポートする操作(リソース単位)
    - キー(Id)となる項目
    - パラメータ/戻り値の型情報

### servicesが提供するツール

TODO: 必要に応じて追記

## やること/やらないこと

### やること

TODO: 必要に応じて追記

### やらないこと

TODO: 必要に応じて追記

## 改訂履歴

- 2022/3/24: 初版作成