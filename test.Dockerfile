FROM local/test/sdcli

RUN mkdir -p /go/src/github.com/asecurityteam/sdcli-test

RUN cd /go/src/github.com/asecurityteam/sdcli-test && \
    sdcli repo go create -- \
            --no-input \
            project_name=sdcli-test \
            project_slug=sdcli-test \
            project_description="A test project" \
            project_namespace="asecurityteam"

RUN cd /go/src/github.com/asecurityteam/sdcli-test && \
    sdcli repo all audit-contract
