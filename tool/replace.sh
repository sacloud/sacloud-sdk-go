#!/bin/bash
set -e
set -o pipefail

declare -A map=(
    ["addon-api-go"]="api/addon"
    ["api-client-go"]="internal/api-client"
    ["apigw-api-go"]="api/apigw"
    ["apprun-api-go"]="api/apprun"
    ["apprun-dedicated-api-go"]="api/apprun-dedicated"
    ["cloudhsm-api-go"]="api/cloudhsm"
    ["dedicated-storage-api-go"]="api/dedicated-storage"
    ["eventbus-api-go"]="api/eventbus"
    ["go-http"]="internal/go-http"
    ["iaas-api-go"]="api/iaas"
    ["iaas-service-go"]="service/iaas"
    ["iam-api-go"]="api/iam"
    ["kms-api-go"]="api/kms"
    ["monitoring-suite-api-go"]="api/monitoring-suite"
    ["nosql-api-go"]="api/nosql"
    ["object-storage-api-go"]="api/object-storage"
    ["packages-go"]="internal/packages"
    ["saclient-go"]="common/saclient"
    ["secretmanager-api-go"]="api/secretmanager"
    ["security-control-api-go"]="api/security-control"
    ["service-endpoint-gateway-api-go"]="api/service-endpoint-gateway"
    ["services"]="internal/services"
    ["simple-notification-api-go"]="api/simple-notification"
    ["simplemq-api-go"]="api/simplemq"
    ["webaccel-api-go"]="api/webaccel"
    ["webaccel-service-go"]="service/webaccel"
    ["workflows-api-go"]="api/workflows"
)

for repo in "${!map[@]}"
do
    path="${map[${repo}]}"

    find "${path}" -type f -name "go.mod" -print0 | \
    xargs -0 -I {} bash -c \
    "echo -n touching {}...; \
    sed -i'' 's|^module github.com/sacloud/${repo}|module github.com/sacloud/sacloud-sdk-go/${path}|g' {}; \
    echo done"

    find "${path}" -type f -name "*.go" -print0 | \
    xargs -0 -I {} bash -c \
    "echo -n touching {} ...; \
    sed -i'' 's|\"github.com/sacloud/${repo}|\"github.com/sacloud/sacloud-sdk-go/${path}|g' {}; \
    echo done"

    pushd "${path}"
    go mod tidy
    popd
done
