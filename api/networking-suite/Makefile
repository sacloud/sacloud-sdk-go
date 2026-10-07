#====================
AUTHOR         ?= The sacloud/sacloud-sdk-go Authors
COPYRIGHT_YEAR ?= 2026

BIN            ?= networking-suite-api-go
GO_FILES       ?= $(shell find . -name '*.go')

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
