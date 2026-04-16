# sacloud/simplemq-api-go

Go言語向けのさくらのクラウド SimpleMQ APIライブラリ

SimpleMQ APIドキュメント: https://manual.sakura.ad.jp/api/cloud/simplemq/

## 概要

sacloud/simplemq-api-goはさくらのクラウド SimpleMQ APIをGo言語から利用するためのAPIライブラリです。

キューの作成や削除などリソース管理のためのQueue APIと、メッセージの送受信を行うMessage APIに分かれています。

## 使い方

[queue_test.go](./queue_test.go) / [message_test.go](./message_test.go) を参照。

:warning:  v1.0に達するまでは互換性のない形で変更される可能性がありますのでご注意ください。

## TODO

- Testの追加

## License

`simplemq-api-go` Copyright (C) 2025- The sacloud/simplemq-api-go authors.
This project is published under [Apache 2.0 License](LICENSE).
