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

With the image installed you call it like (omit the first `--mount` if on Mac):

```bash
export cwd=$(pwd)
export project_path=${cwd#"${GOPATH}/src/"}
docker run -ti \
    # If Linux, mount and configure SSH inside the container.
    --mount src=${SSH_AUTH_SOCK},target=/ssh-agent,type=bind \
    --env SSH_AUTH_SOCK=/ssh-agent \
    # Mount the current project directory to a patch inside the container.
    -v "$(pwd -L):/go/src/${project_path}" \
    # Adjust the container workspace to the newly mounted project.
    -w "/go/src/${project_path}" \
    # Run a command.
    asecurityteam/sdcli:v1 go test
```

To make this easier, you can add this function to your .bashrc file (omit the first `--mount` if on Mac):

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
    docker run --rm \
        --mount src="${SSH_AUTH_SOCK}",target="/ssh-agent",type="bind" \
        --env "SSH_AUTH_SOCK=/ssh-agent" \
        -v "$(pwd -L):/go/src/${project_path}" \
        -w "/go/src/${project_path}" \
        asecurityteam/sdcli:v1 $@
}
```

which will enable you to call the container like:

```bash
sdcli go test
```

## For Shells Other than `bash`

In fish shell, you create a `~/.config/fish/functions/sdcli.fish` file with 755 permissions having contents:

```bash
function sdcli
  set cwd (pwd)
  set gopath "$GOPATH"
  if test -z "$gopath"
    set gopath ~/go # default gopath since 1.8
  end
  # Remove gopath from the front of the directory path. The resulting
  # path is used to construct a mount point inside the container. For
  # go projects this results in them being placed within the gopath
  # of the container. Other languages, such as Python, will still get
  # placed within the gopath but should be agnostic to this fact since
  # they can be placed anywhere.
  set project_path (string replace "$gopath/src/" "" $cwd)
  docker run --rm -v "$cwd:/go/src/$project_path" \
    -w "/go/src/$project_path" \
    asecurityteam/sdcli:v1 $argv
end
```

Some commands are interactive, but if you run `fish` or shells other than
`bash`, you might see "no TTY for interactive shell" or seemingly
inexplicable "project_name [New Project]: Aborted!".  No worries!  Just run in non-interactive mode by
specifying all args on the command line, like:

```bash
sdcli repo go create -- project_name="myproject" project_description="description" --no-input
```

Or start the Docker image with `/bin/bash` as the entrypoint and run `/usr/bin/sdcli $args` from within
(be sure to set `$cwd` and `$project_path` first):

```bash
docker run -it \
    --entrypoint "/bin/bash" -v "$cwd:/go/src/$project_path" \
    -w "/go/src/$project_path" \
    asecurityteam/sdcli:v1
```

<a id="markdown-adding-commands" name="adding-commands"></a>
## Adding Commands

The top-level `sdcli` script will dispatch commands by accumulating all the
arguments and joining them with an `_` character. For example, `sdcli my feature`
will be converted to `sdcli_my_feature` and executed. To add a new command, drop
an executable file in the `./commands` directory and name it according to how you
want the script to be called.
