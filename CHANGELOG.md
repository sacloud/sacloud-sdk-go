# Changelog

## [v0.2.0](https://github.com/sacloud/sacloud-sdk-go/compare/v0.1.0...v0.2.0) - 2026-09-07

### 🚀 New Features
- toolchain go1.27.0 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/234
### 📦 Dependency Updates
- go: bump github.com/stretchr/testify from 1.11.1 to 1.12.0 in /api/iaas/trace/otel by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/227
- go: bump github.com/stretchr/testify from 1.11.1 to 1.12.1 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/226
- go: bump golang.org/x/crypto from 0.54.0 to 0.55.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/231
- go: bump github.com/sacloud/sacloud-sdk-go from 0.0.1 to 0.1.0 in /api/iaas/trace/otel by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/245
- go: bump google.golang.org/grpc from 1.83.0 to 1.83.1 in /api/iaas/trace/otel by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/250
- ci: bump github/codeql-action/upload-sarif from 4.37.6 to 4.37.9 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/251
- go: bump github.com/minio/minio-go/v7 from 7.2.1 to 7.3.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/236
- ci: bump Songmu/tagpr from 1.20.1 to 1.20.2 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/248
- go: bump github.com/ogen-go/ogen from 1.23.0 to 1.24.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/230
- go: bump github.com/hashicorp/terraform-exec from 0.25.2 to 0.25.3 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/243
- go: bump github.com/jlaffaye/ftp from 0.2.2 to 0.2.4 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/244
### Other Changes
- common/saclient: allow 1 min time window by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/233
- apprun: Update OpenAPI to v1.5.0 by @repeatedly in https://github.com/sacloud/sacloud-sdk-go/pull/238
- common/saclient: add `WithAPIRequestRateLimit` by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/240
- Add Networking Suite API by @repeatedly in https://github.com/sacloud/sacloud-sdk-go/pull/241
- chore: standardize includes/go/common.mk and includes/go/single.mk across all modules by @tokuhirom in https://github.com/sacloud/sacloud-sdk-go/pull/246
- chore: centralize Makefile recipes by @tokuhirom in https://github.com/sacloud/sacloud-sdk-go/pull/247
- fix: avoid sharing partially built fake servers by @tokuhirom in https://github.com/sacloud/sacloud-sdk-go/pull/249

## [v0.1.0](https://github.com/sacloud/sacloud-sdk-go/compare/v0.0.1...v0.1.0) - 2026-08-20

### 🚀 New Features
- feat: re-generate V1 codes using ogen v1.23.0 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/215
### 📦 Dependency Updates
- ci: bump golangci/golangci-lint-action from 9.2.1 to 9.3.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/189
- ci: bump sacloud/textlint-action from 0.1.0 to 0.1.1 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/200
- ci: bump dorny/paths-filter from 4.0.1 to 4.0.2 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/199
- ci: bump github/codeql-action/upload-sarif from 4.36.0 to 4.37.5 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/204
- go: bump go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace from 0.68.0 to 0.69.0 in /api/iaas/trace/otel by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/197
- go: bump go.opentelemetry.io/otel/exporters/stdout/stdouttrace from 1.43.0 to 1.44.0 in /api/iaas/trace/otel by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/195
- go: bump github.com/ogen-go/ogen from 1.20.3 to 1.23.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/192
- go: bump github.com/jlaffaye/ftp from 0.2.0 to 0.2.2 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/207
- go: bump github.com/go-playground/validator/v10 from 10.30.2 to 10.30.3 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/208
- go: bump github.com/go-faster/errors from 0.7.1 to 0.8.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/209
- ci: bump github/codeql-action/upload-sarif from 4.37.5 to 4.37.6 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/214
- ci: bump dorny/paths-filter from 4.0.2 to 4.0.3 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/219
- go: bump the otel group across 2 directories with 7 updates by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/225
### Other Changes
- [doc] repository no longer under transition by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/186
- add: expose DefaultZone via EndpointConfig() by @yamamoto-febc in https://github.com/sacloud/sacloud-sdk-go/pull/201
- [CI] get rid of `lint-def` by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/202
- toolchain go1.26.5 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/188
- [chore] delete unnecessary replace by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/203
- [CI] speed up by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/205
- [chore] delete unnecessary scripts by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/206
- Fix fake ProxyLB certificate mapping by @tokuhirom in https://github.com/sacloud/sacloud-sdk-go/pull/216
- api/iaas/trace/otel: use sacloud-sdk-go v0.0.1 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/217
- fix: prevent fake database parameter panic without Conf by @tokuhirom in https://github.com/sacloud/sacloud-sdk-go/pull/218
- iaas: add Ubuntu 26.04 to ostype by @yamamoto-febc in https://github.com/sacloud/sacloud-sdk-go/pull/220
- `README.md` fix after transition by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/223
- api/eventbus: fix Provider.Class filter query injection by @tokuhirom in https://github.com/sacloud/sacloud-sdk-go/pull/222
- [CI] dependabot groups for otel updates by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/221
- chore: categorize GitHub auto-generated release notes by @tokuhirom in https://github.com/sacloud/sacloud-sdk-go/pull/224
- s/スイッチ\+ルータ/ルータ+スイッチ/g by @tokuhirom in https://github.com/sacloud/sacloud-sdk-go/pull/229

## [v0.0.1](https://github.com/sacloud/sacloud-sdk-go/compare/v0.0.0...v0.0.1) - 2026-08-03

- feat: initial implementation by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/1
- feat: .github/copilot-instructions.md by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/2
- feat: .github/dependabot.yml by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/3
- feat: .github/workflows/tests.yaml by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/38
- [CI] fix CI failures by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/66
- pull upstream @ 2026/04/27 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/67
- initial implementation of `package sacloud` by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/68
- internal/go-http: merge `go.mod` into toplevel by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/69
- internal/packages: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/73
- internal/packages/e2e: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/115
- internal/services: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/116
- internal/api-client: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/122
- internal/saclient: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/123
- common/saclient: added by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/125
- api/dedicated-storage: merge go.mod into root go.mod  by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/124
- [CI][chore] dependabot target sync by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/127
- api/simple-notification: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/128
- api/apprun-dedicated: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/129
- api/iam: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/130
- api/iaas: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/131
- api/service-endpoint-gateway: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/132
- api/apprun: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/133
- api/simplemq: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/134
- api/workflows: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/135
- api/security-control: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/136
- api/secretmanager: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/137
- api/object-storage: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/138
- api/monitoring-suite: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/139
- api/eventbus: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/140
- api/cloudhsm: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/141
- api/addon: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/142
- api/nosql: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/143
- api/kms: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/144
- api/apigw: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/154
- api/webaccel: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/155
- service/iaas: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/156
- service/webaccel: merge go.mod into root go.mod by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/157
- pull upstream @ 2026/05/25 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/158
- CI: GitHub Actions update by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/159
- pull upstream @ 2026/05/27 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/160
- CI: give up lint-go at push by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/162
- CI: minimal permission by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/161
- toolchain: go1.26.3 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/168
- common/packages: added by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/166
- Add srn module for SRN support by @repeatedly in https://github.com/sacloud/sacloud-sdk-go/pull/179
- pull upstream @ 2026/08/03 by @shyouhei in https://github.com/sacloud/sacloud-sdk-go/pull/183
- ci: bump actions/setup-go from 6.3.0 to 7.0.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/40
- go: bump golang.org/x/crypto from 0.51.0 to 0.54.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/167
- ci: bump Songmu/tagpr from 1.18.3 to 1.20.1 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/174
- go: bump github.com/minio/minio-go/v7 from 7.0.98 to 7.2.1 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/170
- go: bump go.opentelemetry.io/otel/sdk from 1.43.0 to 1.44.0 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/169
- go: bump go.opentelemetry.io/otel/sdk from 1.43.0 to 1.44.0 in /api/iaas/trace/otel by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/172
- ci: bump actions/checkout from 6.0.2 to 7.0.1 by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/176
- ci: bump golang/govulncheck-action from 31f7c5463448f83528bd771c2d978d940080c9fd to 032d45514ae346b1db93c04b0c90b841c370344f by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/178
- feat: add WithoutProfile option to disable profile loading by @yamamoto-febc in https://github.com/sacloud/sacloud-sdk-go/pull/184
- go: bump google.golang.org/grpc from 1.79.3 to 1.82.1 in /api/iaas/trace/otel by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/182
- go: bump go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc from 1.19.0 to 1.44.0 in /api/iaas/trace/otel by @dependabot[bot] in https://github.com/sacloud/sacloud-sdk-go/pull/173

## [v0.0.0](https://github.com/sacloud/sacloud-sdk-go/commits/v0.0.0) - 2026-04-16
