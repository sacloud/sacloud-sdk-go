#====================
AUTHOR         ?= The sacloud/service-endpoint-gateway-api-go Authors
COPYRIGHT_YEAR ?= 2026-

BIN            ?= service-endpoint-gateway-api-go 
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
ogen:
	go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yaml ./openapi/openapi.json
