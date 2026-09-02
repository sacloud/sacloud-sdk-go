#====================
AUTHOR         ?= The sacloud/saclient-go Authors
COPYRIGHT_YEAR ?= 2022-2025

GO_FILES       ?= $(shell find . -name '*.go')

SACLOUD_SDK_GO_ROOT := $(abspath $(dir $(lastword $(MAKEFILE_LIST)))/../..)
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/common.mk
include $(SACLOUD_SDK_GO_ROOT)/makefiles/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
