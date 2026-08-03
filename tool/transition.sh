#!/bin/bash
[[  -z "${BASH_VERSION}" ]] && echo "run with bash" && exit 1

set -euxo pipefail

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
    # ["makefile"]="makefiles"
    ["monitoring-suite-api-go"]="api/monitoring-suite"
    ["nosql-api-go"]="api/nosql"
    ["object-storage-api-go"]="api/object-storage"
    ["packages-go"]="common/packages"
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

__dir__="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
__repo__="$(git -C "${__dir__}" rev-parse --show-toplevel)"
__wd__="${__repo__}/../transition-workspace"

git -C "${__repo__}" fetch --all --recurse-submodules=yes --progress --jobs "${#map[@]}"
git -C "${__repo__}" worktree add "${__wd__}" main --detach

for repo in "${!map[@]}"
do
    path="${map[${repo}]}"
    git -C "${__wd__}" checkout "${repo}"

    before=$(git -C "${__wd__}" rev-parse HEAD)
    git -C "${__wd__}" subtree merge --no-squash --prefix "${path}" \
        "refs/tags/${repo}/shyouhei/right-before-the-eol" \
        "${repo}"
    after=$(git -C "${__wd__}" rev-parse HEAD)

    if [ "$before" != "$after" ]
    then
        git -C "${__wd__}" commit --amend --no-edit --signoff --gpg-sign
        git -C "${__wd__}" checkout "${path}"
        git -C "${__wd__}" merge --no-ff --signoff --gpg-sign "${repo}"
    fi
done

git -C "${__wd__}" checkout main
git -C "${__wd__}" merge --no-ff --no-edit --signoff --gpg-sign $(
    for repo in "${!map[@]}"
    do
        echo "${map[${repo}]}"
    done
)

git -C "${__repo__}" worktree remove "${__wd__}"
