# Changelog

## [v1.24.0](https://github.com/sacloud/iaas-service-go/compare/v1.23.0...v1.24.0) - 2026-04-23
- fix: adapt to iaas-api-go v1.27.0 breaking changes by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/228
- fix: pass zone to NFS plan query functions by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/230

## [v1.23.0](https://github.com/sacloud/iaas-service-go/compare/v1.22.3...v1.23.0) - 2026-04-16
- use jlaffaye/ftp instead of sacloud/ftps by @tokuhirom in https://github.com/sacloud/iaas-service-go/pull/223
- mark authstatus package as deprecated & use iam-api-go auth context API by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/225
- go: bump golang.org/x/crypto from 0.49.0 to 0.50.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/222

## [v1.22.3](https://github.com/sacloud/iaas-service-go/compare/v1.22.2...v1.22.3) - 2026-03-10
- iaas-api-go v1.26.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/219

## [v1.22.2](https://github.com/sacloud/iaas-service-go/compare/v1.22.1...v1.22.2) - 2026-02-12
- bill: remove ExternalPermission check for bill APIs by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/215
- go: bump golang.org/x/crypto from 0.47.0 to 0.48.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/214
- tagpr: use custom major/minorLabels by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/217

## [v1.22.1](https://github.com/sacloud/iaas-service-go/compare/v1.22.0...v1.22.1) - 2026-02-04
- ftps: configure TLS versions and ServerName in tls.Config by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/212

## [v1.22.0](https://github.com/sacloud/iaas-service-go/compare/v1.21.1...v1.22.0) - 2026-01-28
- go: bump golang.org/x/crypto from 0.46.0 to 0.47.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/208
- sacloud/iaas-api-go@v1.25.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/210

## [v1.22.0](https://github.com/sacloud/iaas-service-go/compare/v1.21.1...v1.22.0) - 2026-01-28
- go: bump golang.org/x/crypto from 0.46.0 to 0.47.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/208
- sacloud/iaas-api-go@v1.25.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/210

## [v1.21.1](https://github.com/sacloud/iaas-service-go/compare/v1.21.0...v1.21.1) - 2026-01-08
- sacloud/iaas-api-go v1.24.2 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/206

## [v1.21.0](https://github.com/sacloud/iaas-service-go/compare/v1.20.1...v1.21.0) - 2025-12-22
- feat: dedicated storage API by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/203
- feat: Confidential VM by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/205

## [v1.20.1](https://github.com/sacloud/iaas-service-go/compare/v1.20.0...v1.20.1) - 2025-12-11
- sacloud/iaas-api-go@v1.23.1 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/201

## [v1.20.0](https://github.com/sacloud/iaas-service-go/compare/v1.19.0...v1.20.0) - 2025-12-10
- Database: MonitoringSuite,ContinuousBackup,DiskEncryption by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/198
- go: bump golang.org/x/crypto from 0.44.0 to 0.46.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/197
- ci: bump actions/checkout from 5 to 6 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/195
- sacloud/iaas-api-go v1.23.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/200

## [v1.19.0](https://github.com/sacloud/iaas-service-go/compare/v1.18.1...v1.19.0) - 2025-11-13
- Send signals to MonitoringSuite by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/192

## [v1.18.1](https://github.com/sacloud/iaas-service-go/compare/v1.18.0...v1.18.1) - 2025-11-05
- Revert "Send signals to MonitoringSuite" by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/188
- iaas-api-go v1.21.1 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/190

## [v1.18.0](https://github.com/sacloud/iaas-service-go/compare/v1.17.0...v1.18.0) - 2025-11-05
- iaas-api-go v1.21.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/185
- Send signals to MonitoringSuite by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/187

## [v1.17.0](https://github.com/sacloud/iaas-service-go/compare/v1.16.0...v1.17.0) - 2025-10-24
- iaas-api-go v1.20.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/182

## [v1.16.0](https://github.com/sacloud/iaas-service-go/compare/v1.15.0...v1.16.0) - 2025-10-10
- feat(koukaryoku-vrt): add GPUModel field (#403) by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/180

## [v1.15.0](https://github.com/sacloud/iaas-service-go/compare/v1.14.2...v1.15.0) - 2025-10-09
- feat: BYOK by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/178

## [v1.14.2](https://github.com/sacloud/iaas-service-go/compare/v1.14.1...v1.14.2) - 2025-10-08
- iaas-api-go v1.17.4 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/175

## [v1.14.1](https://github.com/sacloud/iaas-service-go/compare/v1.14.0...v1.14.1) - 2025-10-07
- github.com/sacloud/iaas-api-go v1.17.3 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/173

## [v1.14.0](https://github.com/sacloud/iaas-service-go/compare/v1.13.0...v1.14.0) - 2025-09-25
- ubuntu-20.04 -> 22.04, 24.04, latest by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/164
- go: bump golang.org/x/crypto from 0.40.0 to 0.41.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/163
- go: bump github.com/stretchr/testify from 1.10.0 to 1.11.1 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/168
- ci: bump actions/setup-go from 5 to 6 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/169
- ci: bump actions/checkout from 4 to 5 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/166
- go: bump golang.org/x/crypto from 0.41.0 to 0.42.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/171
- iaas-api-go@v1.17.2 and fixes validate tag by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/172

## [v1.13.0](https://github.com/sacloud/iaas-service-go/compare/v1.12.1...v1.13.0) - 2025-08-06
- go: bump golang.org/x/crypto from 0.39.0 to 0.40.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/157
- golangci-lint v2 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/159
- goreleaser -> tagpr by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/160
- iaas-api-go v1.17.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/162

## [v1.12.1](https://github.com/sacloud/iaas-service-go/compare/v1.12.0...v1.12.1) - 2025-07-08
- fixes missing SSHKeys assignment in edit request by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/156

## [v1.12.0](https://github.com/sacloud/iaas-service-go/compare/v1.11.1...v1.12.0) - 2025-06-12
- Ubuntu 20.04を廃止 & latestと24.04を追加 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/154

## [v1.11.1](https://github.com/sacloud/iaas-service-go/compare/v1.11.0...v1.11.1) - 2025-06-06
- GoReleaser v2 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/153

## [v1.11.0](https://github.com/sacloud/iaas-service-go/compare/v1.10.0...v1.11.0) - 2025-06-06
- Copyright 2025 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/145
- golangci-lint v1.64.8 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/146
- Go 1.23 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/147
- go: bump golang.org/x/crypto from 0.22.0 to 0.25.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/144
- go: bump github.com/sacloud/packages-go from 0.0.10 to 0.0.11 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/148
- go: bump golang.org/x/crypto from 0.38.0 to 0.39.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/150
- ci: bump goreleaser/goreleaser-action from 5 to 6 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/143
- iaas-api-go v1.16.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/152

## [v1.10.0](https://github.com/sacloud/iaas-service-go/compare/v1.9.2...v1.10.0) - 2024-04-05
- update dependencies by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/138
- go: bump golang.org/x/crypto from 0.16.0 to 0.21.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/137
- go: bump github.com/stretchr/testify from 1.8.4 to 1.9.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/136
- Feature: added EncryptionAlgorithm field to Disk by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/140
- go: bump golang.org/x/crypto from 0.21.0 to 0.22.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/139

## [v1.9.2](https://github.com/sacloud/iaas-service-go/compare/v1.9.1...v1.9.2) - 2023-12-08
- go: bump github.com/stretchr/testify from 1.8.3 to 1.8.4 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/115
- go: bump golang.org/x/crypto from 0.9.0 to 0.10.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/117
- AMDプラン by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/120
- go: bump golang.org/x/crypto from 0.10.0 to 0.12.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/119
- 逆引きレコードのクリアを行えるようにバリデーションルールを修正 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/121
- go 1.21 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/122
- GitHub ActionsでのCIパフォーマンスの改善 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/123
- ci: bump actions/checkout from 3 to 4 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/124
- ci: bump goreleaser/goreleaser-action from 4 to 5 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/125
- go: bump golang.org/x/crypto from 0.12.0 to 0.13.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/126
- go: bump golang.org/x/crypto from 0.13.0 to 0.15.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/128
- ci: bump actions/setup-go from 4 to 5 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/130
- go: bump golang.org/x/crypto from 0.15.0 to 0.16.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/129
- github.com/sacloud/iaas-api-go v1.11.2 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/131

## [v1.9.1](https://github.com/sacloud/iaas-service-go/compare/v1.9.0...v1.9.1) - 2023-06-09
- iaas-api-go@v1.11.1 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/116

## [v1.9.0](https://github.com/sacloud/iaas-service-go/compare/v1.8.2...v1.9.0) - 2023-05-24
- go: bump golang.org/x/crypto from 0.8.0 to 0.9.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/112
- go: bump github.com/stretchr/testify from 1.8.2 to 1.8.3 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/113
- Enhanced DBでタイプ/リージョン/接続元ネットワークを指定可能に by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/114

## [v1.8.2](https://github.com/sacloud/iaas-service-go/compare/v1.8.1...v1.8.2) - 2023-04-19
- FPTサーバ接続時にTLS 1.2を利用 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/111

## [v1.8.1](https://github.com/sacloud/iaas-service-go/compare/v1.8.0...v1.8.1) - 2023-04-19
- 説明文が空文字でも更新できるように修正 by @y2aiskni in https://github.com/sacloud/iaas-service-go/pull/109
- setupにおけるNIC接続後の待ち時間を20秒に増やす by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/110

## [v1.8.0](https://github.com/sacloud/iaas-service-go/compare/v1.7.0...v1.8.0) - 2023-04-17
- go: bump github.com/sacloud/iaas-api-go from 1.9.0 to 1.9.1 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/104
- go: bump golang.org/x/crypto from 0.7.0 to 0.8.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/106
- go: bump github.com/sacloud/iaas-api-go from 1.9.1 to 1.9.2 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/105
- Iaas api go/v1.10.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/108

## [v1.7.0](https://github.com/sacloud/iaas-service-go/compare/v1.6.1...v1.7.0) - 2023-03-27
- go: bump github.com/stretchr/testify from 1.8.1 to 1.8.2 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/96
- go: bump golang.org/x/crypto from 0.0.0-20220214200702-86341886e292 to 0.7.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/97
- go 1.20 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/99
- ci: bump actions/setup-go from 3 to 4 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/100
- iaas-api-go v1.9 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/103

## [v1.6.1](https://github.com/sacloud/iaas-service-go/compare/v1.6.0...v1.6.1) - 2023-03-06
- go: bump github.com/sacloud/iaas-api-go from 1.8.0 to 1.8.1 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/93
- go: bump github.com/sacloud/iaas-api-go from 1.8.1 to 1.8.3 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/94
- fix: disk/CreateRequest by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/98

## [v1.6.0](https://github.com/sacloud/iaas-service-go/compare/v1.5.0...v1.6.0) - 2023-01-17
- go: bump github.com/sacloud/iaas-api-go from 1.7.0 to 1.7.1 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/90
- copyright: 2023 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/91
- iaas-api-go v1.8 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/92

## [v1.5.0](https://github.com/sacloud/iaas-service-go/compare/v1.4.0...v1.5.0) - 2022-12-15
- go: bump github.com/stretchr/testify from 1.8.0 to 1.8.1 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/83
- go: bump github.com/sacloud/iaas-api-go from 1.6.0 to 1.6.1 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/84
- go: bump github.com/sacloud/packages-go from 0.0.6 to 0.0.7 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/86
- ci: bump goreleaser/goreleaser-action from 3 to 4 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/87
- go: bump github.com/sacloud/iaas-api-go from 1.6.1 to 1.7.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/88
- iaas-api-go v1.7: オートスケールのリクエストパラメータ修正 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/89

## [v1.4.0](https://github.com/sacloud/iaas-service-go/compare/v1.3.2...v1.4.0) - 2022-10-24
- go: bump github.com/sacloud/iaas-api-go from 1.4.1 to 1.5.1 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/81
- iaas-api-go v1.6.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/82

## [v1.3.2](https://github.com/sacloud/iaas-service-go/compare/v1.3.1...v1.3.2) - 2022-10-05
- go: bump github.com/sacloud/iaas-api-go from 1.3.1 to 1.3.2 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/74
- go 1.19 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/76
- sacloud/makefile v0.0.7 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/77
- go: bump github.com/sacloud/packages-go from 0.0.5 to 0.0.6 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/78
- github.com/sacloud/iaas-api-go v1.4.1 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/79

## [v1.3.1](https://github.com/sacloud/iaas-service-go/compare/v1.3.0...v1.3.1) - 2022-08-22
- iaas-api-go@v1.3.1 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/73

## [v1.3.0](https://github.com/sacloud/iaas-service-go/compare/v1.2.0...v1.3.0) - 2022-08-19
- go: bump github.com/sacloud/packages-go from 0.0.3 to 0.0.4 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/68
- go: bump github.com/sacloud/packages-go from 0.0.4 to 0.0.5 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/69
- VPCルータでのスケジュールメンテナンス機能 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/70
- update: iaas-api-go by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/71
- iaas-api-go@v1.3.0 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/72

## [v1.2.0](https://github.com/sacloud/iaas-service-go/compare/v1.1.3...v1.2.0) - 2022-07-06
- ディスク作成時にIsSSHKeysEphemeralを指定するとエラーになる問題の修正 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/57
- Go 1.18 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/58
- linterの見直し & エラー修正 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/59
- iaas-api-go@be21866b8f84c2ae8e115459363effe8313a51a1 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/60
- iaas-api-go@98f8c80efc496a402d451518e4cbdd68ff3582bf by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/61
- sacloud/go-template@v0.0.2 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/62
- go: bump github.com/stretchr/testify from 1.7.1 to 1.7.5 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/63
- sacloud/go-template@v0.0.5  by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/66
- go: bump github.com/stretchr/testify from 1.7.5 to 1.8.0 by @dependabot[bot] in https://github.com/sacloud/iaas-service-go/pull/64
- sacloud/iaas-api-go@v1.2.1 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/67

## [v1.1.3](https://github.com/sacloud/iaas-service-go/compare/v1.1.2...v1.1.3) - 2022-06-06
- AutoScaleのUpdateで現在の値が反映されない問題を修正 by @yamamoto-febc in https://github.com/sacloud/iaas-service-go/pull/55
