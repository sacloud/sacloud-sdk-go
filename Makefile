#====================
AUTHOR         ?= The sacloud/simple-notification-api-go Authors
COPYRIGHT_YEAR ?= 2022-2026

BIN            ?= simple-notification-api-go 
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
include includes/go/single.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
ogen:
	ogen -package v1 -target apis/v1 -clean -config ogen-config.yaml ./openapi/openapi.yaml
