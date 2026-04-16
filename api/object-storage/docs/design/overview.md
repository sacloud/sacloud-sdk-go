# object-storage-api-go

- URL: https://github.com/sacloud/object-storage-api-go/pull/1
- Parent: https://github.com/sacloud/sacloud-go/pull/1
- Author: @yamamoto-febc

## 概要

[さくらのクラウド オブジェクトストレージAPI](https://manual.sakura.ad.jp/cloud/objectstorage/api.html)をGo言語から利用するためのライブラリを提供する。

さくらのクラウドオブジェクトストレージAPIは専用サーバPHYと同じくOpenAPI v3.0でのAPI定義が公開されている。

- [さくらのクラウド オブジェクトストレージAPI OpenAPIでの定義](https://manual.sakura.ad.jp/cloud/objectstorage/api/api-json.html)

このAPI定義を用いてGo言語向けのコードを生成する。  
またより簡単に使うためのラップしたコードやテスト用モックサーバも提供する。

## やること/やらないこと

### やること

- OpenAPIでのAPI定義から生成したコードの提供
- 生成したコードをラップするコードの提供
- テスト用モックサーバの提供
- テスト用モックサーバを操作するためのCLIの提供

### やらないこと

- Amazon S3互換APIで出来る範囲の機能の実装
  作成済みバケット/オブジェクトの操作といったS3互換APIで出来る機能は対象外とする。

## 実装

- コード生成には[ogen](https://github.com/ogen-go/ogen)を用いる

## 改訂履歴

- 2022/2/26: 初版作成
- 2026/2/10: ogenベースにしたので実装の項目を更新
