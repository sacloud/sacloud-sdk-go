# Changelog

## [v0.5.1](https://github.com/sacloud/simplemq-api-go/compare/v0.5.0...v0.5.1) - 2026-04-13
- Remove api-client-go dependency by @repeatedly in https://github.com/sacloud/simplemq-api-go/pull/80
- ci: bump Songmu/tagpr from 1.10.0 to 1.18.1 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/79
- ci: bump actions/setup-go from 6.1.0 to 6.4.0 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/76

## [v0.5.0](https://github.com/sacloud/simplemq-api-go/compare/v0.4.0...v0.5.0) - 2026-03-09
- CIでのE2Eテストの廃止 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/56
- OpenAPI定義の追加 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/58
- Configure Queue API to use is1a zone by @yamamoto-febc in https://github.com/sacloud/simplemq-api-go/pull/60
- go toolchainを上げる by @a2not in https://github.com/sacloud/simplemq-api-go/pull/61
- Make API root URL configurable via saclient-go endpoints by @yamamoto-febc in https://github.com/sacloud/simplemq-api-go/pull/63
- Align SimpleMQ service keys and endpoint env vars to snake_case by @yamamoto-febc in https://github.com/sacloud/simplemq-api-go/pull/66
- Fix queue list filter query encoding by @yamamoto-febc in https://github.com/sacloud/simplemq-api-go/pull/72

## [v0.4.0](https://github.com/sacloud/simplemq-api-go/compare/v0.3.0...v0.4.0) - 2026-01-19
- feat: use of `saclient-go` by @shyouhei in https://github.com/sacloud/simplemq-api-go/pull/45
- ci: bump actions/setup-go from 6.0.0 to 6.1.0 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/37
- ci: bump Songmu/tagpr from 1.9.0 to 1.10.0 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/43
- ci: bump actions/checkout from 5.0.0 to 6.0.1 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/42
- go: bump github.com/ogen-go/ogen from 1.16.0 to 1.18.0 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/41
- Revert "Release for v0.4.0" by @a2not in https://github.com/sacloud/simplemq-api-go/pull/48

## [v0.3.0](https://github.com/sacloud/simplemq-api-go/compare/v0.2.0...v0.3.0) - 2025-11-26
- Update github settings by @a2not in https://github.com/sacloud/simplemq-api-go/pull/26
- chore(textlint): ignore CHANGELOG.md file by @a2not in https://github.com/sacloud/simplemq-api-go/pull/33
- ci: bump actions/setup-go from 5.5.0 to 6.0.0 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/29
- go: bump github.com/stretchr/testify from 1.10.0 to 1.11.1 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/24
- go: bump github.com/sacloud/api-client-go from 0.3.2 to 0.3.3 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/23
- メッセージAPIの429エラーを追加 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/34
- go: bump github.com/go-faster/jx from 1.1.0 to 1.2.0 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/35
- go: bump github.com/ogen-go/ogen from 1.14.0 to 1.16.0 by @dependabot[bot] in https://github.com/sacloud/simplemq-api-go/pull/32
- chore: update Go to 1.25 and golangci-lint to v2.6.2 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/36

## [v0.2.0](https://github.com/sacloud/simplemq-api-go/compare/v0.1.3...v0.2.0) - 2025-07-04
- Migrate from Example to Test by @repeatedly in https://github.com/sacloud/simplemq-api-go/pull/14
- Use simplemq.Error instead of common error by @repeatedly in https://github.com/sacloud/simplemq-api-go/pull/15
- Rename to Read by @repeatedly in https://github.com/sacloud/simplemq-api-go/pull/18
- Use api-client-go's APIError by @repeatedly in https://github.com/sacloud/simplemq-api-go/pull/20

## [v0.1.3](https://github.com/sacloud/simplemq-api-go/compare/v0.1.2...v0.1.3) - 2025-06-10
- textlint修正 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/10
- nil check追加 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/12
- OpenAPIの更新に合わせたコード生成 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/13

## [v0.1.2](https://github.com/sacloud/simplemq-api-go/compare/v0.1.1...v0.1.2) - 2025-06-05
- govulncheckによる脆弱性チェックを追加 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/7
- IconとTagsを作成・更新時で指定可能にする by @a2not in https://github.com/sacloud/simplemq-api-go/pull/9
- go: bump github.com/ogen-go/ogen from 1.13.0 to 1.14.0 by @dependabot in https://github.com/sacloud/simplemq-api-go/pull/6

## [v0.1.1](https://github.com/sacloud/simplemq-api-go/compare/v0.1.0...v0.1.1) - 2025-06-03
- クライアント初期化時にapi-client-goのオプションをサポート by @a2not in https://github.com/sacloud/simplemq-api-go/pull/4

## [v0.1.0](https://github.com/sacloud/simplemq-api-go/commits/v0.1.0) - 2025-06-02
- Simple MQ SDKの初期実装 by @a2not in https://github.com/sacloud/simplemq-api-go/pull/2
