# iaas-service-go

- URL: https://github.com/sacloud/iaas-service-go/pull/1
- Author: @yamamoto-febc

## 概要

さくらのクラウドではいくつかの性質の異なるAPIが提供されている。  

- IaaS API
- 請求関連
- ウェブアクセラレータ API
- オブジェクトストレージ API
- 専用サーバPHY API

従来はIaas API(sacloud/libsacloud)としてこれらの一部のAPI向けのAPIライブラリを提供していた。

さくらのクラウドのサービス拡大と共に専用サーバPHYやオブジェクトストレージのAPIなどが公開されたが、
これらのAPIはOpenAPIによるAPI定義が提供されていたりHTTP レスポンスステータスコードの扱いが異なったりと
Iaas APIとはライブラリで処理すべき内容が異なっている。

このままIaaS API側で実装すると(違いを吸収するための)共通部分の実装の肥大化が予想されるため、APIライブラリを以下のように分ける。  

- [IaaS API: sacloud/iaas-api-go](https://github.com/sacloud/iaas-api-go)  
  以下のAPIも含む  
    - 請求 API
    - ウェブアクセラレータ API
- [オブジェクトストレージ API: sacloud/object-storage-api-go](https://github.com/sacloud/object-storage-api-go)
- [専用サーバPHY API: sacloud/phy-go](https://github.com/sacloud/phy-go)

各ライブラリは単体でも利用可能とするが、Usacloudに対するlibsacloudのhelper/serviceパッケージのように、各APIに対し統一的にCRUD+L操作が行えるようにしたい。  

そこで、高レベルAPIライブラリとしてsacloud/xxx-service-goを作成し、
libsacloudのhelper/serviceパッケージの移植やその他高レベル操作の実装/集約などを行いたい。  

- [IaaS: sacloud/iaas-service-go](https://github.com/sacloud/iaas-service-go)  
- [オブジェクトストレージ: sacloud/object-storage-service-go](https://github.com/sacloud/object-storage-service-go)
- [専用サーバPHY: sacloud/phy-service-go](https://github.com/sacloud/phy-service-go)
 
以降はIaaS向けであるsacloud/iaas-service-goについて記載する。

## やること/やらないこと

### やること

- libsacloudからの移植
  - helper/builderとhelper/serviceの統合
  - 統合したhelper/serviceの移植
  - 汎用的な処理をsacloud/packages-goに切り出し
  - クライアント周り(helper/apiやprofile)をsacloud/api-client-goへ切り出し

## 実装

- libsacloudから必要部分を移植  
- 下位ライブラリ(IaaS/オブジェクトストレージ/PHY)を意識せず使えるようなAPIの設計

## 改訂履歴

- 2022/2/26: 初版作成
- 2022/3/7: パッケージ構成追加
- 2022/3/23: 各プラットフォームごとにxxx-service-goリポジトリとして独立させるように変更