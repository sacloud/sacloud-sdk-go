#====================
AUTHOR         ?= The sacloud/apprun-api-go authors
COPYRIGHT_YEAR ?= 2021-2024

BIN            ?= dist/sacloud-apprun-fake-server
GO_ENTRY_FILE  ?= cmd/sacloud-apprun-fake-server/*.go
BUILD_LDFLAGS  ?=

include includes/go/common.mk
include includes/go/single.mk
#====================

default: gen $(DEFAULT_GOALS)

.PHONY: tools
tools: dev-tools
	npm i -g @redocly/cli@latest
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.4.1

.PHONY: clean-all
clean-all:
	find . -type f -name "*_gen.go" -delete
	rm apis/v1/spec/original-openapi.yaml
	rm apis/v1/spec/openapi.json

.PHONY: gen
gen: _gen fmt goimports set-license

.PHONY: _gen
_gen: apis/v1/spec/original-openapi.yaml apis/v1/spec/openapi.json apis/v1/zz_types_gen.go apis/v1/zz_client_gen.go apis/v1/zz_server_gen.go
	go generate ./...

apis/v1/spec/original-openapi.yaml: apis/v1/spec/original-openapi.json
	redocly bundle apis/v1/spec/original-openapi.json -o apis/v1/spec/original-openapi.yaml

apis/v1/spec/openapi.json: apis/v1/spec/openapi.yaml
	redocly bundle apis/v1/spec/openapi.yaml -o apis/v1/spec/openapi.json

apis/v1/zz_types_gen.go: apis/v1/spec/openapi.yaml apis/v1/spec/codegen/types.yaml
	oapi-codegen --old-config-style -config apis/v1/spec/codegen/types.yaml apis/v1/spec/openapi.yaml

apis/v1/zz_client_gen.go: apis/v1/spec/openapi.yaml apis/v1/spec/codegen/client.yaml
	oapi-codegen --old-config-style -config apis/v1/spec/codegen/client.yaml apis/v1/spec/openapi.yaml

apis/v1/zz_server_gen.go: apis/v1/spec/openapi.yaml apis/v1/spec/codegen/server.yaml
	oapi-codegen -config apis/v1/spec/codegen/server.yaml apis/v1/spec/openapi.yaml

lint-def:
	# NOTE: 上流側のOpenAPI定義ではtag周りの警告が出るため "-F warn"の指定を外しておく
	docker run --rm -v $$PWD:$$PWD -w $$PWD stoplight/spectral:latest lint apis/v1/spec/openapi.yaml
