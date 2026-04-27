#!/bin/bash
set -e
set -o pipefail

WD="sacloud-sdk-go"
initial=$([[ ! -d "${WD}" ]]; echo "${?}")
map=(
    "addon-api-go:api/addon"
    "api-client-go:internal/api-client"
    "apigw-api-go:api/apigw"
    "apprun-api-go:api/apprun"
    "apprun-dedicated-api-go:api/apprun-dedicated"
    "cloudhsm-api-go:api/cloudhsm"
    "dedicated-storage-api-go:api/dedicated-storage"
    "eventbus-api-go:api/eventbus"
    "go-http:internal/go-http"
    "iaas-api-go:api/iaas"
    "iaas-service-go:service/iaas"
    "iam-api-go:api/iam"
    "kms-api-go:api/kms"
    "makefile:makefiles"
    "monitoring-suite-api-go:api/monitoring-suite"
    "nosql-api-go:api/nosql"
    "object-storage-api-go:api/object-storage"
    "packages-go:internal/packages"
    "saclient-go:internal/saclient"
    "secretmanager-api-go:api/secretmanager"
    "security-control-api-go:api/security-control"
    "service-endpoint-gateway-api-go:api/service-endpoint-gateway"
    "services:internal/services"
    "simple-notification-api-go:api/simple-notification"
    "simplemq-api-go:api/simplemq"
    "webaccel-api-go:api/webaccel"
    "webaccel-service-go:service/webaccel"
    "workflows-api-go:api/workflows"
)

if [[ "${initial}" -eq 0 ]]
then
    mkdir -p "${WD}"
    git -C "${WD}" init -b main
    (cd "${WD}" && go work init)
    git -C "${WD}" add go.work
    git -C "${WD}" commit --no-edit --signoff --gpg-sign -m "Initial implementation of go.work"

    for item in "${map[@]}"
    do
        IFS=":" read -r repo path <<< "$item"
        git -C "${WD}" remote add "${repo}" "git@github.com:sacloud/${repo}.git"
    done
fi

git -C "${WD}" fetch --all --recurse-submodules=yes --progress --jobs "${#map[@]}"

for item in "${map[@]}"
do
    IFS=":" read -r repo path <<< "${item}"
    if [[ "${initial}" -eq 0 ]]
    then
        git -C "${WD}" checkout -b "${repo}" "main"
        git -C "${WD}" subtree add --no-squash --prefix "${path}" "${repo}" main
        git -C "${WD}" commit --amend --no-edit --signoff --gpg-sign
    else
        git -C "${WD}" checkout "${repo}"
        before=$(git -C "${WD}" rev-parse HEAD)
        git -C "${WD}" subtree pull --no-squash --prefix "${path}" "${repo}" main
        after=$(git -C "${WD}" rev-parse HEAD)

        if [ "$before" != "$after" ]
        then
            git -C "${WD}" commit --amend --no-edit --signoff --gpg-sign
        fi
    fi
done

git -C "${WD}" checkout main
git -C "${WD}" merge --no-ff --no-edit --signoff --gpg-sign $(
    for item in "${map[@]}"
    do
        IFS=":" read -r repo path <<< "${item}"
        echo "${repo}"
    done
)

if [[ "${initial}" -eq 0 ]]
then
    (cd "${WD}" && go work use -r .)
    git -C "${WD}" add go.work
    git -C "${WD}" commit --signoff --gpg-sign -m "go work use -r ."
fi