# Changelog

## [v1.26.0](https://github.com/sacloud/iaas-api-go/compare/v1.25.0...v1.26.0) - 2026-03-10
- tagpr: use custom major/minorLabels by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/448
- go: bump github.com/sacloud/saclient-go from 0.2.7 to 0.3.1 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/450
- feat: update constant value for simple-notification by @sugyan in https://github.com/sacloud/iaas-api-go/pull/452
- Remove trailing slash by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/453
- go: bump golang.org/x/crypto from 0.47.0 to 0.48.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/451

## [v1.25.0](https://github.com/sacloud/iaas-api-go/compare/v1.24.2...v1.25.0) - 2026-01-27
- docs: typo by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/443
- feat: added ProxySSL field to EnhancedLB.Servers by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/445
- Add string representations for load balancer plans by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/446
- go: bump github.com/sacloud/api-client-go from 0.3.4 to 0.3.5 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/442
- go: bump github.com/sacloud/saclient-go from 0.2.6 to 0.2.7 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/440
- go: bump golang.org/x/crypto from 0.46.0 to 0.47.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/439

## [v1.24.2](https://github.com/sacloud/iaas-api-go/compare/v1.24.1...v1.24.2) - 2026-01-08
- github.com/sacloud/saclient-go v0.2.6 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/436

## [v1.24.2](https://github.com/sacloud/iaas-api-go/compare/v1.24.1...v1.24.2) - 2026-01-08
- github.com/sacloud/saclient-go v0.2.6 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/436

## [v1.24.1](https://github.com/sacloud/iaas-api-go/compare/v1.24.0...v1.24.1) - 2025-12-22
- feat: support ConfidentialVM field for query.FindServerPlan() by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/433

## [v1.24.0](https://github.com/sacloud/iaas-api-go/compare/v1.23.1...v1.24.0) - 2025-12-22
- feat: Confidential VM by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/427
- sacloud/api-client-go v0.3.4 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/429
- feat: use of `saclient-go` by @shyouhei in https://github.com/sacloud/iaas-api-go/pull/430
- fix: api.NewCaller() to respect profile values by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/432
- feat: Dedicated Storage API support for disks and private hosts by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/431

## [v1.23.1](https://github.com/sacloud/iaas-api-go/compare/v1.23.0...v1.23.1) - 2025-12-11
- fix: add almalinux8 and rockylinux8 back to ostype by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/425

## [v1.23.0](https://github.com/sacloud/iaas-api-go/compare/v1.22.0...v1.23.0) - 2025-12-10
- Database: Disk Encryption by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/420
- Database: MonitoringSuite by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/422
- ci: bump actions/checkout from 5 to 6 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/417
- go: bump github.com/sacloud/packages-go from 0.0.11 to 0.0.12 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/418
- go: bump golang.org/x/crypto from 0.44.0 to 0.46.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/419
- Database: Continuous Backup by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/424

## [v1.22.0](https://github.com/sacloud/iaas-api-go/compare/v1.21.1...v1.22.0) - 2025-11-13
- feat: Add monitoring suite integration for SimpleMonitor, DNS, GSLB, ELB, and VPNRouter by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/414
- go: bump golang.org/x/crypto from 0.43.0 to 0.44.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/413

## [v1.21.1](https://github.com/sacloud/iaas-api-go/compare/v1.21.0...v1.21.1) - 2025-11-05
- Revert "Send signals to MonitoringSuite" by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/411

## [v1.21.0](https://github.com/sacloud/iaas-api-go/compare/v1.20.0...v1.21.0) - 2025-11-05
- Send signals to MonitoringSuite by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/407
- Handle cases where Interfaces may contain null values in DB/NFS/LoadBalancer by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/409

## [v1.20.0](https://github.com/sacloud/iaas-api-go/compare/v1.19.0...v1.20.0) - 2025-10-23
- Update ostype: AlmaLinux 9/10 & Rocky Linux 9/10 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/405

## [v1.19.0](https://github.com/sacloud/iaas-api-go/compare/v1.18.0...v1.19.0) - 2025-10-10
- feat(koukaryoku-vrt): add GPUModel field by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/403

## [v1.18.0](https://github.com/sacloud/iaas-api-go/compare/v1.17.4...v1.18.0) - 2025-10-09
- go: bump golang.org/x/crypto from 0.42.0 to 0.43.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/401
- Disk: BYOK by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/400

## [v1.17.4](https://github.com/sacloud/iaas-api-go/compare/v1.17.3...v1.17.4) - 2025-10-08
- Always include password field in ContainerRegistryUser by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/398

## [v1.17.3](https://github.com/sacloud/iaas-api-go/compare/v1.17.2...v1.17.3) - 2025-10-07
- Handle optional BinlogUsedSizeKiB and DelayTimeSec in DB monitor by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/396

## [v1.17.2](https://github.com/sacloud/iaas-api-go/compare/v1.17.1...v1.17.2) - 2025-09-25
- go: bump golang.org/x/crypto from 0.41.0 to 0.42.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/394
- ci: bump actions/setup-go from 5 to 6 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/393

## [v1.17.1](https://github.com/sacloud/iaas-api-go/compare/v1.17.0...v1.17.1) - 2025-09-04
- ubuntu-20.04 -> 22.04, 24.04, latest by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/383
- go: bump golang.org/x/crypto from 0.40.0 to 0.41.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/382
- InstanceStatusの判定を小文字で統一 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/391
- go: bump github.com/stretchr/testify from 1.10.0 to 1.11.1 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/389
- go: bump github.com/sacloud/api-client-go from 0.3.2 to 0.3.3 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/387
- textlint: ignore CHANGELOG.md by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/392
- ci: bump actions/checkout from 4 to 5 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/385

## [v1.17.0](https://github.com/sacloud/iaas-api-go/compare/v1.16.1...v1.17.0) - 2025-08-06
- go: bump golang.org/x/crypto from 0.39.0 to 0.40.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/375
- go: bump github.com/sacloud/api-client-go from 0.3.0 to 0.3.2 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/374
- golangci-lint v2 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/378
- goreleaser -> tagpr by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/379
- feat: is1c by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/381

## [v1.16.1](https://github.com/sacloud/iaas-api-go/compare/v1.16.0...v1.16.1) - 2025-07-14
- docomo回線提供終了対応 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/373
- go: bump golang.org/x/crypto from 0.38.0 to 0.39.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/372
- Generators improvements by @ldez in https://github.com/sacloud/iaas-api-go/pull/377

## [v1.16.0](https://github.com/sacloud/iaas-api-go/compare/v1.15.0...v1.16.0) - 2025-06-05
- go: bump github.com/sacloud/api-client-go from 0.2.10 to 0.3.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/371
- Remove SSH key pair generation func by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/370

## [v1.15.0](https://github.com/sacloud/iaas-api-go/compare/v1.14.0...v1.15.0) - 2025-05-08
- go: bump github.com/sacloud/go-http from 0.1.8 to 0.1.9 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/355
- Copyright 2025 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/357
- go: bump golang.org/x/crypto from 0.31.0 to 0.33.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/359
- go: bump github.com/sacloud/packages-go from 0.0.10 to 0.0.11 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/356
- シンプル通知 β版 対応 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/360
- go: bump github.com/fsnotify/fsnotify from 1.8.0 to 1.9.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/364
- go: bump golang.org/x/crypto from 0.33.0 to 0.38.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/367
- Remove ubuntu2004 from ostype by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/368
- Added debian12 to ostype by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/369

## [v1.14.0](https://github.com/sacloud/iaas-api-go/compare/v1.13.0...v1.14.0) - 2024-12-20
- go: bump github.com/stretchr/testify from 1.9.0 to 1.10.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/352
- ci: bump goreleaser/goreleaser-action from 5 to 6 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/317
- goreleaser v2に対応 by @hekki in https://github.com/sacloud/iaas-api-go/pull/354
- go: bump github.com/fsnotify/fsnotify from 1.7.0 to 1.8.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/353
- go: bump golang.org/x/crypto from 0.25.0 to 0.31.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/351

## [v1.13.0](https://github.com/sacloud/iaas-api-go/compare/v1.12.0...v1.13.0) - 2024-12-18
- OSType更新 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/346
- go: bump golang.org/x/crypto from 0.22.0 to 0.25.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/324
- go: bump github.com/huandu/xstrings from 1.4.0 to 1.5.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/318
- ExternalPermission: supports AppRun and KoukaryokuDOK by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/350

## [v1.12.0](https://github.com/sacloud/iaas-api-go/compare/v1.11.2...v1.12.0) - 2024-04-05
- go: bump github.com/sacloud/go-http from 0.1.7 to 0.1.8 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/275
- update dependencies by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/295
- go: bump golang.org/x/crypto from 0.16.0 to 0.21.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/294
- go: bump github.com/stretchr/testify from 1.8.4 to 1.9.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/293
- Feature: Disk Encryption by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/279
- go: bump golang.org/x/crypto from 0.21.0 to 0.22.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/310

## [v1.11.2](https://github.com/sacloud/iaas-api-go/compare/v1.11.1...v1.11.2) - 2023-12-08
- go: bump github.com/stretchr/testify from 1.8.3 to 1.8.4 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/219
- go: bump golang.org/x/crypto from 0.9.0 to 0.10.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/223
- AMDプラン対応 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/238
- go: bump golang.org/x/crypto from 0.10.0 to 0.12.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/237
- go 1.21 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/242
- GitHub ActionsでのCIパフォーマンスの改善 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/244
- ci: bump actions/checkout from 3 to 4 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/246
- ci: bump goreleaser/goreleaser-action from 4 to 5 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/247
- go: bump github.com/sacloud/go-http from 0.1.6 to 0.1.7 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/250
- go: bump golang.org/x/crypto from 0.12.0 to 0.13.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/251
- go: bump github.com/sacloud/api-client-go from 0.2.8 to 0.2.9 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/252
- Update trace/otel package by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/260
- trace/otel: exampleをOTLP/gRPCを使うように修正 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/265
- go: bump golang.org/x/crypto from 0.13.0 to 0.15.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/266
- go: bump github.com/fsnotify/fsnotify from 1.6.0 to 1.7.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/259
- ci: bump actions/setup-go from 4 to 5 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/272
- go: bump golang.org/x/crypto from 0.15.0 to 0.16.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/271
- Update dependencies by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/274

## [v1.11.1](https://github.com/sacloud/iaas-api-go/compare/v1.11.0...v1.11.1) - 2023-06-09
- helper/power: モバイルゲートウェイでの強制終了APIの無効化 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/222

## [v1.11.0](https://github.com/sacloud/iaas-api-go/compare/v1.10.0...v1.11.0) - 2023-05-23
- go: bump github.com/stretchr/testify from 1.8.2 to 1.8.3 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/211
- go: bump golang.org/x/crypto from 0.8.0 to 0.9.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/204
- Enhanced DBでのMariaDB対応 & 東京リージョン対応 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/212
- Enhanced DBでのアクセス元IPアドレス制限機能 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/213
- update dependencies by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/214
- sacloud/packages-go@v0.0.9 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/215

## [v1.10.0](https://github.com/sacloud/iaas-api-go/compare/v1.9.2...v1.10.0) - 2023-04-17
- EDBのJSON構造変更対応 by @chibiegg in https://github.com/sacloud/iaas-api-go/pull/197
- EDBのconfigエンドポイントの追加 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/199
- go: bump golang.org/x/crypto from 0.7.0 to 0.8.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/195

## [v1.9.2](https://github.com/sacloud/iaas-api-go/compare/v1.9.1...v1.9.2) - 2023-04-06
- DNSのHTTPS と SVCB レコードへの対応 by @kazeburo in https://github.com/sacloud/iaas-api-go/pull/190

## [v1.9.1](https://github.com/sacloud/iaas-api-go/compare/v1.9.0...v1.9.1) - 2023-03-28
- fix: ELB.BackendHttpKeepAliveの省略対応 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/188

## [v1.9.0](https://github.com/sacloud/iaas-api-go/compare/v1.8.3...v1.9.0) - 2023-03-27
- go: bump github.com/stretchr/testify from 1.8.1 to 1.8.2 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/176
- go: bump golang.org/x/crypto from 0.0.0-20220214200702-86341886e292 to 0.7.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/178
- go 1.20 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/180
- ci: bump actions/setup-go from 3 to 4 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/181
- ProxyLB: BackendHttpKeepAlive by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/184
- AutoScale: Disabled by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/185
- AutoScale: ScheduleTrigger by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/186
- Update dependencies by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/187

## [v1.8.3](https://github.com/sacloud/iaas-api-go/compare/v1.8.2...v1.8.3) - 2023-01-30
- HTTPエラーログ出力が有効な場合にAPIログ出力を無効化 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/165

## [v1.8.2](https://github.com/sacloud/iaas-api-go/compare/v1.8.1...v1.8.2) - 2023-01-30
- sacloud/api-client-go@v0.2.5 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/164

## [v1.8.1](https://github.com/sacloud/iaas-api-go/compare/v1.8.0...v1.8.1) - 2023-01-18
- VPCルータ: ping by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/158

## [v1.8.0](https://github.com/sacloud/iaas-api-go/compare/v1.7.1...v1.8.0) - 2023-01-17
- copyright: 2023 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/154
- ELB: SourceIPs by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/157

## [v1.7.1](https://github.com/sacloud/iaas-api-go/compare/v1.7.0...v1.7.1) - 2022-12-19
- データベースアプライアンス冗長化プランでの型不一致エラーの修正 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/153

## [v1.7.0](https://github.com/sacloud/iaas-api-go/compare/v1.6.2...v1.7.0) - 2022-12-14
- オートスケール: トラフィック量トリガー by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/145
- go: bump github.com/huandu/xstrings from 1.3.3 to 1.4.0 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/143
- go: bump github.com/sacloud/packages-go from 0.0.6 to 0.0.7 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/147
- VPCRouter: DHGroup by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/150
- ci: bump goreleaser/goreleaser-action from 3 to 4 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/149
- helper: cleanup.DeleteServer() by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/151

## [v1.6.2](https://github.com/sacloud/iaas-api-go/compare/v1.6.1...v1.6.2) - 2022-12-01
- MIRACLE LINUX9 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/142

## [v1.6.1](https://github.com/sacloud/iaas-api-go/compare/v1.6.0...v1.6.1) - 2022-11-25
- go: bump github.com/stretchr/testify from 1.8.0 to 1.8.1 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/130
- go: bump github.com/huandu/xstrings from 1.3.2 to 1.3.3 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/131
- go: bump github.com/sacloud/api-client-go from 0.2.3 to 0.2.4 by @dependabot[bot] in https://github.com/sacloud/iaas-api-go/pull/135
- データベースアプライアンスの/statusで特定のフィールドの値が文字列/数値混在する問題への対応 by @yamamoto-febc in https://github.com/sacloud/iaas-api-go/pull/137
