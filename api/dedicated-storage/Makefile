#====================
AUTHOR         ?= The sacloud/dedicated-storage-api-go Authors
COPYRIGHT_YEAR ?= 2022-2025

BIN            ?= dedicated-storage-api-go
GO_FILES       ?= $(shell find . -name '*.go')

include includes/go/common.mk
#====================

default: $(DEFAULT_GOALS)
tools: dev-tools
	go get -tool github.com/ogen-go/ogen/cmd/ogen@latest

gen: 
	go tool ogen -package v1 -target apis/v1 -clean -config ogen-config.yaml ./openapi/openapi.yaml

.PHONY: lint-def
lint-def:
	@echo "running lint-def..."
	@docker run --rm -v $$PWD:$$PWD -w $$PWD stoplight/spectral:latest lint -F warn ./oepnapi/openapi.yaml