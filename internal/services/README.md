# services

`sacloud-sdk-go/internal/services` — さくらのクラウド高レベルAPIライブラリ向け基盤

さくらのクラウドAPIライブラリをラップし、各リソースに対する統一的な操作インターフェースを提供します。

## 実装プロジェクト

- IaaS: [service/iaas](https://github.com/sacloud/sacloud-sdk-go/tree/main/service/iaas)
- オブジェクトストレージ: [sacloud/object-storage-service-go](https://github.com/sacloud/object-storage-service-go)
- 専用サーバPHY: [sacloud/phy-service-go](https://github.com/sacloud/phy-service-go)
- ウェブアクセラレータ: [service/webaccel](https://github.com/sacloud/sacloud-sdk-go/tree/main/service/webaccel)

Note: IaaSとオブジェクトストレージは現在sacloud/servicesへの対応作業中です。

## License

Copyright (C) 2022-2025 The sacloud/sacloud-sdk-go Authors.
This project is published under [Apache 2.0 License](LICENSE).
