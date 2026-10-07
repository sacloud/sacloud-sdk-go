#====================
AUTHOR         ?= The sacloud/service-endpoint-gateway-api-go Authors
COPYRIGHT_YEAR ?= 2026-

BIN            ?= service-endpoint-gateway-api-go 
GO_FILES       ?= $(shell find . -name '*.go')

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
ogen:
	go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yaml ./openapi/openapi.json
