#====================
AUTHOR         ?= The sacloud/eventbus-api-go Authors
COPYRIGHT_YEAR ?= 2025-2026

BIN            ?= eventbus-api-go
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
include includes/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools

.PHONY: gen
gen:
	go tool ogen --config ogen-config.yaml --target ./apis/v1 --package v1 --clean ./openapi/openapi.json
