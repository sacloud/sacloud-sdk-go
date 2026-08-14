# sacloud-sdk-go/api/service-endpoint-gateway

さくらのクラウド サービスエンドポイントゲートウェイ Go言語向け APIライブラリ

マニュアル: https://manual.sakura.ad.jp/cloud/network/switch/seg.html

## 概要

sacloud-sdk-go/api/service-endpoint-gatewayはさくらのクラウド サービスエンドポイントゲートウェイ APIをGo言語から利用するためのAPIライブラリです。

> [!NOTE]
> このライブラリはサービスエンドポイントゲートウェイ関連のAPIのみを扱います。サーバおよびスイッチの作成や操作はサポートしていないため必要に応じて [sacloud/sacloud-sdk-go/api/iaas](https://github.com/sacloud/sacloud-sdk-go/tree/main/api/iaas)と組み合わせてご利用ください。

## 利用イメージ
利用例: [example_test.go](./example_test.go)

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## License

Copyright (C) 2026- The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](LICENSE).