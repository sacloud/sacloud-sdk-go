# OpenAPI定義

シンプルMQのOpenAPI定義は [さくらのクラウド API ポータル](https://manual.sakura.ad.jp/api/cloud/portal/) にて公開されています。

- キュー管理API: https://manual.sakura.ad.jp/api/cloud/portal/openapis/simplemq-sacloud-api.yaml
- メッセージ送受信API: https://manual.sakura.ad.jp/api/cloud/portal/openapis/simplemq-api.yaml

取得した定義はそれぞれ `queue.yaml`、`message.yaml` として配置し、`make gen` でクライアントを再生成する。

## ogenによる生成コードの修正

キュー一覧の取得APIにおいて、詳細欄には記述があるもののOpenAPI定義として表現できないクエリパラメータを必要とするので、patchファイルで生成後の修正を管理している。

修正内容については [patchファイル](../patch/01_list_filter.patch) を参照のこと。
