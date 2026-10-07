#====================
AUTHOR         ?= The sacloud/iam-api-go Authors
COPYRIGHT_YEAR ?= 2025-

BIN            ?= iam-api-go
GO_FILES       ?= $(shell find . -name '*.go')

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools

PHONY: ogen
ogen:
	go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yml ./openapi/openapi.json