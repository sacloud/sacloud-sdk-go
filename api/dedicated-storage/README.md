# sacloud-sdk-go/api/dedicated-storage

Go言語向け　さくらのクラウド　専有ストレージ　APIライブラリ

## 概要

`sacloud-sdk-go/api/dedicated-storage`は、さくらのクラウドの専有ストレージAPIをGo言語から利用するためのライブラリです。

Note: このライブラリは専有ストレージ関連のAPIのみを扱います。ディスクの作成や専有ホストの操作はサポートしていないため必要に応じて [sacloud-sdk-go/api/iaas](https://pkg.go.dev/github.com/sacloud/sacloud-sdk-go/api/iaas) や [sacloud-sdk-go/service/iaas](https://pkg.go.dev/github.com/sacloud/sacloud-sdk-go/service/iaas) と組み合わせてご利用ください。

利用例: [example_test.go](./example_test.go)

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## インストール

```bash
go get github.com/sacloud/sacloud-sdk-go/api/dedicated-storage
```

## ogenによるコード生成

以下のコマンドを実行

```
$ make gen
```

## License

Copyright (C) 2025- The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](LICENSE).