<a id="markdown-sdcli---security-development-repository-tools" name="sdcli---security-development-repository-tools"></a>
# SDCLI - Security Development Repository Tools

*Status: Incubation*

<!-- TOC -->

- [SDCLI - Security Development Repository Tools](#sdcli---security-development-repository-tools)
    - [Overview](#overview)
    - [Usage](#usage)
    - [Adding Commands](#adding-commands)

<!-- /TOC -->

<a id="markdown-overview" name="overview"></a>
## Overview

This project is a docker image that we use as a CLI for repository management.
It bundles our tools related to creating a new project, running any tests or
CI automation, and keep our repositories consistent in how they build.

The current feature set is:

```bash
sdcli
  go
      dep # install go project dependencies
      lint # run our standard go linter
      test # run unit tests
      integration # run integration tests
      coverage # generate a coverage report
  repo
      all # generic repo tools
          add-oss # add license and contributing files
          audit-contract # verify the repo implements the contract
      go # go repo tools
          add-docker # add a Dockerfile
          add-layout # render the standard layout
          add-lint # add linter configuration
          add-travis # add a travisci configuraiton
          create # generate a full go project
```

<a id="markdown-usage" name="usage"></a>
## Usage

The project is delivered as a docker image that contains our tooling:

```bash
docker pull asecurityteam/sdcli:v1
```

With the image installed you call it like:

```bash
docker run -ti \
    # Mount and configure SSH inside the container.
    --mount src=${SSH_AUTH_SOCK},target=/ssh-agent,type=bind \
    --env SSH_AUTH_SOCK=/ssh-agent \
    # Mount the current project directory to a patch inside the container.
    --mount src="$(pwd)",target=/go/src/github.com/asecurityteam/go-vpcflow,type=bind \
    # Adjust the container workspace to the newly mounted project.
    -w /go/src/github.com/asecurityteam/go-vpcflow \
    # Run a command.
    asecurityteam/sdcli:v1 go test
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
        asecurityteam/sdcli:v1 "$@"
}
```

which will enable you to call the container like:

```bash
sdcli go test
```

<a id="markdown-adding-commands" name="adding-commands"></a>
## Adding Commands

The top-level `sdcli` script will dispatch commands by accumulating all the
arguments and joining them with an `_` character. For example, `sdcli my feature`
will be converted to `sdcli_my_feature` and executed. To add a new command, drop
an executable file in the `./commands` directory and name it according to how you
want the script to be called.
