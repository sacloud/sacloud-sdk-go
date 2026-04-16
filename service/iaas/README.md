# iaas-service-go

[![Go Reference](https://pkg.go.dev/badge/github.com/sacloud/iaas-service-go.svg)](https://pkg.go.dev/github.com/sacloud/iaas-service-go)
[![Tests](https://github.com/sacloud/iaas-service-go/workflows/Tests/badge.svg)](https://github.com/sacloud/iaas-service-go/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/sacloud/iaas-service-go)](https://goreportcard.com/report/github.com/sacloud/iaas-service-go)

さくらのクラウド高レベルAPIライブラリ  

## 概要

iaas-service-goは[sacloud/libsacloud v2](https://github.com/sacloud/libsacloud)の後継プロジェクトで、さくらのクラウド APIのうちのIaaS部分を担当します。  
[sacloud/iaas-api-go](https://github.com/sacloud/iaas-api-go)を用いた高レベルAPIを提供します。  

概要/設計/実装方針: [docs/overview.md](https://github.com/sacloud/iaas-service-go/blob/main/docs/design/overview.md)

### libsacloudとiaas-service-goのバージョン対応表

| libsacloud | iaas-api-go | Note/Status                       |
|------------|-------------|-----------------------------------|
| v1         | -           | libsacloud v1系はiaas-api-goへの移植対象外 |
| v2         | v1          | 開発中                               |
| v3(未リリース)  | v2          | 未リリース/未着手                         |

## License

`sacloud/iaas-service-go` Copyright (C) 2022-2025 [The sacloud/iaas-service-go Authors](AUTHORS).

This project is published under [Apache 2.0 License](LICENSE.txt).
