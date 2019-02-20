# SDCLI - Security Development Repository Tools

This project helps us automate working with our repositories and implement
our repository contract. Current feature set:

```bash
sdcli
    -   go # go project related tools
        -   dep # Install build dependencies. Implements the `dep` contract.
        -   lint # Run all static analysis. Implements the `lint` contract.
        -   test # Run unit tests and record coverage. Implements the `test` contract.
        -   integration # Run integraiton tests and record coverage. Implements `integration`.
        -   coverage # Generate a coverage report. Implements `coverage`.
```

## Usage

The project is delivered as a docker image that contains our tooling:

```bash
docker pull docker.atl-paas.net/sox/asecurityteam/sdcli:v1
```

With the image installed you call it like:

```bash
docker run -ti \
    # Mount and configure SSH inside the container.
    --mount src=${SSH_AUTH_SOCK},target=/ssh-agent,type=bind \
    --env SSH_AUTH_SOCK=/ssh-agent \
    # Mount the current project directory to a patch inside the container.
    --mount src="$(pwd)",target=/go/src/bitbucket.org/asecurityteam/vpc-scheduler,type=bind \
    # Adjust the container workspace to the newly mounted project.
    -w /go/src/bitbucket.org/asecurityteam/vpc-scheduler \
    # Run a command.
    docker.atl-paas.net/sox/asecurityteam/sdcli:v1 go test
```

To make this easier, you can add this function to your .bashrc file:

```bash
sdcli() {
    local cwd="$(pwd)"
    local gopath="${GOPATH}"
    if [[ "${gopath}" == "" ]]; then
        gopath=~/go # default gopath since 1.8
    fi
    # Remove gopath from the front of the directory path. The resulting
    # path is used to construct a mount point inside the container. For
    # go projects this results in them being placed within the gopath
    # of the container. Other languages, such as Python, will still get
    # placed within the gopath but should be agnostic to this fact since
    # they can be placed anywhere.
    local project_path=${cwd#"${gopath}/src/"}
    docker run -ti \
        --mount src="${SSH_AUTH_SOCK}",target="/ssh-agent",type="bind" \
        --env "SSH_AUTH_SOCK=/ssh-agent" \
        --mount src="$(pwd)",target="/go/src/${project_path}",type="bind" \
        -w "/go/src/${project_path}" \
        docker.atl-paas.net/sox/asecurityteam/sdcli:v1 $@
}
```

which will enable you to call the container like:

```bash
sdcli go test
```

## Adding Commands

To add a new language or overarching feature pack to the project then edit
the `./commands/sdcli` file and add a new case the redirects to a new executable:

```bash
case ${PACKAGE} in
    go)
        sdcli_go $@
        ;;
    # New feature
    newfeature)
        sdcli_newfeature $@
    *)
        echo "Unknown package ${PACKAGE}"
        exit 1
        ;;
esac
```

Then create a new file with the name `sdcli_<package>` that will executed commands.

```bash
#!/usr/bin/env bash

# Capture the target command and pop it from the args
COMMAND="$1"
shift

case ${COMMAND} in
    cmd1)
        sdcli_newfeature_cmd1 $@
        ;;
    *)
        echo "Unknown newfeature command ${COMMAND}"
        exit 1
        ;;
esac
```

From here, each command is a separate executable named `sdcli_<package>_<command>`
that performs the function.
