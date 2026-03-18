# Changelog

## [v0.3.3](https://github.com/sacloud/saclient-go/compare/v0.3.2...v0.3.3) - 2026-03-18
- [chore] tagpr label settings by @shyouhei in https://github.com/sacloud/saclient-go/pull/83
- `SAKURA_SERVICE_PRINCIPAL_KEY_KID` environment variable by @shyouhei in https://github.com/sacloud/saclient-go/pull/85
- Permission check on Windows by @shyouhei in https://github.com/sacloud/saclient-go/pull/86
- go: bump github.com/hashicorp/terraform-plugin-framework from 1.18.0 to 1.19.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/87

## [v0.4.0](https://github.com/sacloud/saclient-go/compare/v0.3.1...v0.4.0) - 2026-03-10
- [chore] rename functions by @shyouhei in https://github.com/sacloud/saclient-go/pull/66
- refactor hoist out withSettingHeader by @shyouhei in https://github.com/sacloud/saclient-go/pull/67
- Refactor introduce resultOption[T].decompose() by @shyouhei in https://github.com/sacloud/saclient-go/pull/68
- new `WithDefaultTimeout()` option by @shyouhei in https://github.com/sacloud/saclient-go/pull/69
- Refactor hoist out `config.hasSome()` by @shyouhei in https://github.com/sacloud/saclient-go/pull/72
- Refactor split `prioritizedParameterValue2()` by @shyouhei in https://github.com/sacloud/saclient-go/pull/73
- [feat] trace retries by @shyouhei in https://github.com/sacloud/saclient-go/pull/77
- [Feat] `WithZone()` by @shyouhei in https://github.com/sacloud/saclient-go/pull/81
- update golang @1.26.1 by @shyouhei in https://github.com/sacloud/saclient-go/pull/82
- go: bump github.com/golang-jwt/jwt/v5 from 5.3.0 to 5.3.1 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/71
- go: bump github.com/hashicorp/terraform-plugin-framework from 1.17.0 to 1.18.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/80
- ci: bump actions/setup-go from 6.2.0 to 6.3.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/78
- ci: bump Songmu/tagpr from 1.10.0 to 1.17.1 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/79

## [v0.3.1](https://github.com/sacloud/saclient-go/compare/v0.3.0...v0.3.1) - 2026-02-03
- Prefer bearer authentication by @shyouhei in https://github.com/sacloud/saclient-go/pull/64
- ci: bump actions/setup-go from 6.1.0 to 6.2.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/58
- go: bump github.com/sacloud/api-client-go from 0.3.4 to 0.3.5 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/63

## [v0.3.0](https://github.com/sacloud/saclient-go/compare/v0.2.7...v0.3.0) - 2026-01-28
- feat: add configurable service endpoints support by @yamamoto-febc in https://github.com/sacloud/saclient-go/pull/57
- ci: bump actions/checkout from 6.0.1 to 6.0.2 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/60

## [v0.2.7](https://github.com/sacloud/saclient-go/compare/v0.2.6...v0.2.7) - 2026-01-15
- prefer dynamic configurations over others by @shyouhei in https://github.com/sacloud/saclient-go/pull/55

## [v0.2.6](https://github.com/sacloud/saclient-go/compare/v0.2.5...v0.2.6) - 2026-01-08
- go 1.25.0 by @shyouhei in https://github.com/sacloud/saclient-go/pull/52
- treat null value in profile as unspecified key by @shyouhei in https://github.com/sacloud/saclient-go/pull/54

## [v0.2.5](https://github.com/sacloud/saclient-go/compare/v0.2.4...v0.2.5) - 2025-12-24
- explicit `bytes.Clone` to avoid ownership glitch by @shyouhei in https://github.com/sacloud/saclient-go/pull/45
- SAKURA_ environment variables by @shyouhei in https://github.com/sacloud/saclient-go/pull/46
- Add authorization header related runtime options by @shyouhei in https://github.com/sacloud/saclient-go/pull/49
- go: bump github.com/sacloud/api-client-go from 0.3.3 to 0.3.4 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/50
- ci: bump Songmu/tagpr from 1.9.0 to 1.10.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/48
- fix: handle EOF from gzip.NewReader gracefully by @yamamoto-febc in https://github.com/sacloud/saclient-go/pull/51

## [v0.2.4](https://github.com/sacloud/saclient-go/compare/v0.2.3...v0.2.4) - 2025-12-17
- migration from `api-client-go` by @shyouhei in https://github.com/sacloud/saclient-go/pull/40
- propagate `APIRequestTimeout` from config by @shyouhei in https://github.com/sacloud/saclient-go/pull/41
- handling of non-valued environment variables by @shyouhei in https://github.com/sacloud/saclient-go/pull/42
- follow terraform parameter precedence by @shyouhei in https://github.com/sacloud/saclient-go/pull/43

## [v0.2.3](https://github.com/sacloud/saclient-go/compare/v0.2.2...v0.2.3) - 2025-12-15
- ProfileAPI.Update(): avoid merging nested arrays by @shyouhei in https://github.com/sacloud/saclient-go/pull/38

## [v0.2.2](https://github.com/sacloud/saclient-go/compare/v0.2.1...v0.2.2) - 2025-12-11
- create the current file if it does not exist when executing ProfileOp::SetCurrentName() by @yamamoto-febc in https://github.com/sacloud/saclient-go/pull/36
- ci: bump actions/checkout from 6.0.0 to 6.0.1 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/34
- go: bump github.com/hashicorp/terraform-plugin-framework from 1.16.1 to 1.17.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/35
- avoid caching token response from a fake endpoint by @shyouhei in https://github.com/sacloud/saclient-go/pull/33

## [v0.2.1](https://github.com/sacloud/saclient-go/compare/v0.2.0...v0.2.1) - 2025-12-05
- ignore CHANGELOG.md by @shyouhei in https://github.com/sacloud/saclient-go/pull/25
- align profile priority with other parameters by @shyouhei in https://github.com/sacloud/saclient-go/pull/24
- empty string in config means "not set" by @shyouhei in https://github.com/sacloud/saclient-go/pull/26
- ci: bump actions/checkout from 5.0.0 to 6.0.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/29
- ci: bump actions/setup-go from 6.0.0 to 6.1.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/30
- `ProfileAPI.Update()` by @shyouhei in https://github.com/sacloud/saclient-go/pull/31
- publicize Middleware by @shyouhei in https://github.com/sacloud/saclient-go/pull/27
- go: bump github.com/sacloud/packages-go from 0.0.11 to 0.0.12 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/32

## [v0.2.0](https://github.com/sacloud/saclient-go/commits/v0.2.0) - 2025-11-18
- ci: bump Songmu/tagpr from 1.8.4 to 1.9.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/3
- ci: bump actions/setup-go from 5.5.0 to 6.0.0 by @dependabot[bot] in https://github.com/sacloud/saclient-go/pull/4
- initial repository structure by @shyouhei in https://github.com/sacloud/saclient-go/pull/5
- initial implementation of `Error` by @shyouhei in https://github.com/sacloud/saclient-go/pull/6
- initial implementation of `Profile` by @shyouhei in https://github.com/sacloud/saclient-go/pull/7
- initial implementation of `parameter` by @shyouhei in https://github.com/sacloud/saclient-go/pull/9
- initial implementation of `HttpRequestDoer` by @shyouhei in https://github.com/sacloud/saclient-go/pull/10
- example request CLI by @shyouhei in https://github.com/sacloud/saclient-go/pull/11
- README by @shyouhei in https://github.com/sacloud/saclient-go/pull/12
- add WithUserAgent by @shyouhei in https://github.com/sacloud/saclient-go/pull/13
- add ClientAPI.ProfileOp  by @shyouhei in https://github.com/sacloud/saclient-go/pull/14
- fix test failures when there is ~/.usacloud by @shyouhei in https://github.com/sacloud/saclient-go/pull/15
- typo fix by @shyouhei in https://github.com/sacloud/saclient-go/pull/16
- configrable FlagSet error handling by @shyouhei in https://github.com/sacloud/saclient-go/pull/17
- --trace is boolean by @shyouhei in https://github.com/sacloud/saclient-go/pull/18
- dedicated http.Client by @shyouhei in https://github.com/sacloud/saclient-go/pull/19
- go @ 1.25.3 by @shyouhei in https://github.com/sacloud/saclient-go/pull/20
- Rename package by @shyouhei in https://github.com/sacloud/saclient-go/pull/21
- add `Ptr()` by @shyouhei in https://github.com/sacloud/saclient-go/pull/22
- calculation of AuthPreference by @shyouhei in https://github.com/sacloud/saclient-go/pull/23
