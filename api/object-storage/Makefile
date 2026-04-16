#
# Copyright 2022-2025 The sacloud/object-storage-api-go Authors
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

#====================
AUTHOR         ?= The sacloud/object-storage-api-go authors
COPYRIGHT_YEAR ?= 2022-2026

include includes/go/common.mk
#====================

default: gen $(DEFAULT_GOALS)
tools: dev-tools

.PHONY: tools
tools: dev-tools
	npm install -g @redocly/cli
	go get -tool github.com/ogen-go/ogen/cmd/ogen@latest

.PHONY: clean-all
clean-all:
	find . -type f -name "*_gen.go" -delete
	rm -f apis/v2

.PHONY: gen
gen: _gen fmt goimports set-license

.PHONY: _gen
_gen:
	go tool ogen -package v1 -target apis/v2 -clean -config ogen-config.yaml openapi/openapi.json

.PHONY: lint-def
lint-def:
	docker run --rm -v $$PWD:$$PWD -w $$PWD stoplight/spectral:latest lint -F warn openapi/openapi.yaml
