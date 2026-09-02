#====================
AUTHOR         ?= The sacloud/simplemq-api-go Authors
COPYRIGHT_YEAR ?= 2022-2025

BIN            ?= simplemq-api-go
GO_FILES       ?= $(shell find . -name '*.go')

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools

.PHONY: gen
gen:
	go tool ogen --config ogen-config.yaml --target ./apis/v1/queue --package queue --clean ./openapi/queue.yaml
	go tool ogen --config ogen-config.yaml --target ./apis/v1/message --package message --clean ./openapi/message.yaml
	patch -p1 < patch/01_list_filter.patch
