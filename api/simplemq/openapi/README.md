# OpenAPI定義

シンプルMQのOpenAPI定義は下記にて公開されています。

- キュー管理API: https://manual.sakura.ad.jp/api/cloud/simplemq/sacloud/
- メッセージ送受信API: https://manual.sakura.ad.jp/api/cloud/simplemq/

## ogenによる生成コードの修正

キュー一覧の取得APIにおいて、詳細欄には記述があるもののOpenAPI定義として表現できないクエリパラメータを必要とするので、patchファイルで生成後の修正を管理している。

修正内容については [patchファイル](../patch/01_list_filter.patch) を参照のこと。
