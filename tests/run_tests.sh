#!/usr/bin/env bash

RENDERDIR="$(mktemp -d)"
chmod 766 -R "${RENDERDIR}"
docker build -t local/test/sdcli .
docker run \
    --rm \
    -i \
    --mount src="${RENDERDIR}",target=/go/src/github.com/asecurityteam/sdcli-test,type=bind \
    -w="/go/src/github.com/asecurityteam/sdcli-test" \
    local/test/sdcli \
        repo go create -- \
            --no-input \
            project_name=sdcli-test \
            project_slug=sdcli-test \
            project_description="A test project" \
            project_namespace="asecurityteam"
docker run \
    --rm \
    -i \
    --mount src="${RENDERDIR}",target=/go/src/github.com/asecurityteam/sdcli-test,type=bind \
    -w="/go/src/github.com/asecurityteam/sdcli-test" \
    local/test/sdcli \
        repo all audit-contract
