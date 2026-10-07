#====================
AUTHOR         ?= The sacloud/monitoring-suite-api-go Authors
COPYRIGHT_YEAR ?= 2022-2025

BIN            ?= monitoring-suite-api-go
GO_FILES       ?= $(shell find . -name '*.go')

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools

gen:
	go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yml ./openapi/openapi.json
