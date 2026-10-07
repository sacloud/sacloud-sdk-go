#====================
AUTHOR         ?= The sacloud/cloudhsm-api-go Authors
COPYRIGHT_YEAR ?= 2022-2025

BIN            ?= cloudhsm-api-go
GO_FILES       ?= $(shell find . -name '*.go')

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
