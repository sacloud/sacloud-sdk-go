# sacloud-sdk-go/api/simplemq

Go言語向けのさくらのクラウド SimpleMQ APIライブラリ

SimpleMQ APIドキュメント: https://manual.sakura.ad.jp/api/cloud/simplemq/

## 概要

sacloud-sdk-go/api/simplemqはさくらのクラウド SimpleMQ APIをGo言語から利用するためのAPIライブラリです。

キューの作成や削除などリソース管理のためのQueue APIと、メッセージの送受信を行うMessage APIに分かれています。

## 使い方

[queue_test.go](./queue_test.go) / [message_test.go](./message_test.go) を参照。

> [!WARNING]
> v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## TODO

- Testの追加

## License

Copyright (C) 2025- The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](LICENSE).
