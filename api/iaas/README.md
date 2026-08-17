# sacloud-sdk-go/api/iaas

Go言語向けのさくらのクラウドIaaS APIライブラリ

## 概要

sacloud-sdk-go/api/iaasは[sacloud/libsacloud v2](https://github.com/sacloud/libsacloud)の後継プロジェクトで、さくらのクラウド APIのうちのIaaS部分を担当します。

概要/設計/実装方針: [docs/design/overview.md](docs/design/overview.md)

### libsacloudとiaas-api-goのバージョン対応表

| libsacloud | iaas-api-go | Note/Status                       |
|------------|-------------|-----------------------------------|
| v1         | -           | libsacloud v1系はiaas-api-goへの移植対象外 |
| v2         | v1          | 開発中                               |
| v3(未リリース)  | v2          | 未リリース/未着手                         |


### 関連プロジェクト

- [service/iaas](../../service/iaas): `github.com/sacloud/sacloud-sdk-go/api/iaas` を用いた高レベルAPIライブラリ
- [common/saclient](../../common/saclient): sacloudプロダクト向けHTTP/APIクライアントライブラリ(環境変数やプロファイルの処理など)
- [common/packages](../../common/packages): sacloudプロダクト向けの汎用パッケージ群

## License

Copyright (C) 2025- The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](LICENSE).
