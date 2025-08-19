#====================
AUTHOR         ?= The sacloud/go-template Authors
COPYRIGHT_YEAR ?= 2022-2025

BIN            ?= go-template
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
include includes/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
