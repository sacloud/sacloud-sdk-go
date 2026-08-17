# sacloud-sdk-go/service/iaas

さくらのクラウド高レベルAPIライブラリ

## 概要

sacloud-sdk-go/service/iaasは[sacloud/libsacloud v2](https://github.com/sacloud/libsacloud)の後継プロジェクトで、さくらのクラウド APIのうちのIaaS部分を担当します。
[sacloud-sdk-go/api/iaas](../../api/iaas)を用いた高レベルAPIを提供します。

概要/設計/実装方針: [docs/design/overview.md](docs/design/overview.md)

### libsacloudとsacloud-sdk-go/service/iaasのバージョン対応表

| libsacloud | iaas-api-go | Note/Status                       |
|------------|-------------|-----------------------------------|
| v1         | -           | libsacloud v1系はiaas-api-goへの移植対象外 |
| v2         | v1          | 開発中                               |
| v3(未リリース)  | v2          | 未リリース/未着手                         |

## License

Copyright (C) 2022-2025 The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](LICENSE.txt).
