#====================
AUTHOR         ?= The sacloud/apprun-dedicated-api-go Authors
COPYRIGHT_YEAR ?= 2026-

BIN            ?= apprun-dedicated-api-go
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
include includes/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
