# makefiles

このリポジトリ内で共通利用するMakefileです。

- `go/common.mk`: Goモジュール向けの共通ゴール
- `go/single.mk`: 単一バイナリをビルドするモジュール向けの追加ゴール

## Usage

各モジュールのMakefileから、リポジトリルートを基準にインクルードします。
`SACLOUD_SDK_GO_ROOT` の相対パスは、モジュールの階層に応じて調整してください。

```makefile
# 必要に応じて変数定義
AUTHOR         ?= The sacloud/example Authors
COPYRIGHT_YEAR ?= 2026
BIN            ?= example
DEFAULT_GOALS  ?= fmt set-license go-licenses-check goimports lint vulncheck test build

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/single.mk

# ゴールを追加
default: $(DEFAULT_GOALS)
tools: dev-tools # toolsゴールはsacloudプロダクト向け日次CIを行うプロジェクトでは必須
```

## License

`sacloud/makefile` Copyright (C) 2022-2026 The sacloud/makefile Authors.

This project is published under [Apache 2.0 License](LICENSE).
