#====================
AUTHOR         ?= The sacloud/eventbus-api-go Authors
COPYRIGHT_YEAR ?= 2025-2026

BIN            ?= eventbus-api-go
GO_FILES       ?= $(shell find . -name '*.go')

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools

.PHONY: gen
gen:
	go tool ogen --config ogen-config.yaml --target ./apis/v1 --package v1 --clean ./openapi/openapi.json
