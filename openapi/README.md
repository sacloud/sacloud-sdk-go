openapi.jsonは https://manual.sakura.ad.jp/api/cloud/monitoring-suite/ からダウンロードできるJSONに以下の変更を加えたものです

```console
$ jq 'del(.paths.[].[].requestBody.content.["application/x-www-form-urlencoded", "multipart/form-data"])' openapi.json
```