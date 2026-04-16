## API定義(openapi.yaml)について

オリジナルの定義ファイルは以下のサイトで公開されています。
[https://manual.sakura.ad.jp/sakura-apprun-api/spec.html](https://manual.sakura.ad.jp/sakura-apprun-api/spec.html)

公開されている定義ファイルのままでは https://github.com/deepmap/oapi-codegen でコード生成した際にコンパイルエラーが出るため、手作業で修正しています。
修正は以下のように行っています。

- オリジナルの定義ファイルをダウンロード、`original-openapi.json`として保存
- `make gen`を実行することで`original-openapi.json`から`original-openapi.yaml`へ変換
- `original-openapi.yaml`をコピー/編集し`openapi.yaml`を作成

`original-openapi.yaml`については生成される対象なため`.gitignore`に登録されています。
今後オリジナルの定義ファイルが更新された場合は`original-openapi.yaml`と`openapi.yaml`のdiffを取り、適宜`openapi.yaml`へ反映するようにします。
