# sacloud/service-endpoint-gateway-api-go
[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/service-endpoint-gateway-api-go.svg)](https://pkg.go.dev/github.com/sacloud/service-endpoint-gateway-api-go)
[![Tests](https://github.com/sacloud/service-endpoint-gateway-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/service-endpoint-gateway-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/service-endpoint-gateway-api-go)](https://goreportcard.com/report/github.com/sacloud/service-endpoint-gateway-api-go)

さくらのクラウド サービスエンドポイントゲートウェイ Go言語向け APIライブラリ

マニュアル: https://manual.sakura.ad.jp/cloud/network/switch/seg.html


## 概要
sacloud/service-endpoint-gateway-api-goはさくらのクラウド サービスエンドポイントゲートウェイ APIをGo言語から利用するためのAPIライブラリです。

Note: このライブラリはサービスエンドポイントゲートウェイ関連のAPIのみを扱います。サーバおよびスイッチの作成や操作はサポートしていないため必要に応じて [sacloud/iaas-api-go](https://github.com/sacloud/iaas-api-go)と組み合わせてご利用ください。

## 利用イメージ
利用例: [example_test.go](./example_test.go)
:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

# Licence
`service-endpoint-gateway-api-go` Copyright (C) 2026- The sacloud/service-endpoint-gateway-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).