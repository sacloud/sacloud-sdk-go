# Changelog

## [v0.4.0](https://github.com/sacloud/workflows-api-go/compare/v0.3.0...v0.4.0) - 2026-03-23
- go oldstableバージョンでのテスト追加 by @a2not in https://github.com/sacloud/workflows-api-go/pull/42
- プロファイル/環境変数による接続先切り替えに対応 by @to-shoji in https://github.com/sacloud/workflows-api-go/pull/50
- fix(test): retry execution delete to handle 409 errors by @a2not in https://github.com/sacloud/workflows-api-go/pull/51
- OpenAPI定義の更新 by @a2not in https://github.com/sacloud/workflows-api-go/pull/52
- retrigger tagpr by @to-shoji in https://github.com/sacloud/workflows-api-go/pull/54

## [v0.3.0](https://github.com/sacloud/workflows-api-go/compare/v0.2.0...v0.3.0) - 2026-02-10
- OpenAPI定義ファイルの追加 by @a2not in https://github.com/sacloud/workflows-api-go/pull/36
- chore(openapi): OpenAPI定義の更新を反映 by @a2not in https://github.com/sacloud/workflows-api-go/pull/40
- go: bump github.com/sacloud/saclient-go from 0.2.7 to 0.3.0 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/38
- ci: bump Songmu/tagpr from 1.10.0 to 1.15.0 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/39

## [v0.2.0](https://github.com/sacloud/workflows-api-go/compare/v0.1.0...v0.2.0) - 2026-01-23
- OpenAPI specの更新に追従 by @a2not in https://github.com/sacloud/workflows-api-go/pull/22
- feat(workflow): ListWorkflowSuggest APIの実装 by @a2not in https://github.com/sacloud/workflows-api-go/pull/27
- ci: bump Songmu/tagpr from 1.9.0 to 1.10.0 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/20
- ci: bump actions/checkout from 6.0.0 to 6.0.1 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/16
- feat(subscription): Subscription APIの実装 by @a2not in https://github.com/sacloud/workflows-api-go/pull/28
- `SAKURA_` で始まる環境変数をテストに利用する by @a2not in https://github.com/sacloud/workflows-api-go/pull/33
- ci: bump actions/setup-go from 6.1.0 to 6.2.0 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/29
- go: bump github.com/sacloud/saclient-go from 0.2.6 to 0.2.7 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/31
- ci: bump actions/checkout from 6.0.1 to 6.0.2 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/32
- go: bump github.com/ogen-go/ogen from 1.17.0 to 1.18.0 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/17

## [v0.1.0](https://github.com/sacloud/workflows-api-go/commits/v0.1.0) - 2025-12-09
- Initialize the repository from go-template into Workflows by @a2not in https://github.com/sacloud/workflows-api-go/pull/7
- ci: bump Songmu/tagpr from 1.8.4 to 1.9.0 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/4
- ci: bump actions/setup-go from 5.5.0 to 6.1.0 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/3
- ogen generate from openapi.yaml by @a2not in https://github.com/sacloud/workflows-api-go/pull/9
- ci: bump actions/checkout from 5.0.0 to 6.0.0 by @dependabot[bot] in https://github.com/sacloud/workflows-api-go/pull/11
- initial implementation of `client.go` by @a2not in https://github.com/sacloud/workflows-api-go/pull/10
- feat: add `WorkflowAPI` interface and implementation by @a2not in https://github.com/sacloud/workflows-api-go/pull/12
- feat: add `RevisionAPI` interface and implementation by @a2not in https://github.com/sacloud/workflows-api-go/pull/13
- feat: add `ExecutionAPI` interface and implementation by @a2not in https://github.com/sacloud/workflows-api-go/pull/14
