# sacloud/dedicated-storage-api-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/dedicated-storage-api-go.svg)](https://pkg.go.dev/github.com/sacloud/dedicated-storage-api-go)
[![Tests](https://github.com/sacloud/dedicated-storage-api-go/workflows/Tests/badge.svg)](https://github.com/sacloud/dedicated-storage-api-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/dedicated-storage-api-go)](https://goreportcard.com/report/github.com/sacloud/dedicated-storage-api-go)

Go言語向け　さくらのクラウド　専有ストレージ　APIライブラリ

## 概要

`sacloud/dedicated-storage-api-go`は、さくらのクラウドの専有ストレージAPIをGo言語から利用するためのライブラリです。

Note: このライブラリは専有ストレージ関連のAPIのみを扱います。ディスクの作成や専有ホストの操作はサポートしていないため必要に応じて [sacloud/iaas-api-go](https://github.com/sacloud/iaas-api-go)と組み合わせてご利用ください。

利用例: [example_test.go](./example_test.go)

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## ogenによるコード生成

以下のコマンドを実行

```
$ make gen
```

## License

`dedicated-storage-api-go` Copyright (C) 2022-2025 The sacloud/dedicated-storage-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).