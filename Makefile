#====================
AUTHOR         ?= The sacloud/monitoring-suite-api-go Authors
COPYRIGHT_YEAR ?= 2022-2025

BIN            ?= monitoring-suite-api-go
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
