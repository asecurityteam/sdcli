<a id="markdown-sdcli---security-development-repository-tools" name="sdcli---security-development-repository-tools"></a>
# SDCLI - Security Development Repository Tools

*Status: Incubation*

<!-- TOC -->

- [SDCLI - Security Development Repository Tools](#sdcli---security-development-repository-tools)
    - [Overview](#overview)
    - [Usage](#usage)
        - [For Shells Other than `bash`](#for-shells-other-than-bash)
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
  yaml
      lint #runs yamllint against all yamls in current directory
  version # lists the versions of the installed languages and applications in SDCLI
```

<a id="markdown-usage" name="usage"></a>
## Usage

The project is delivered as a docker image that contains our tooling:

```bash
docker pull asecurityteam/sdcli:v3
```

With the image installed you call it like (omit the first `--mount` if on Mac):

```bash
export cwd=$(pwd)
export project_path=${cwd#"${GOPATH}/src/"}
docker run -ti \
    # If Linux, mount and configure SSH inside the container.
    --mount src="${SSH_AUTH_SOCK}",target="/ssh-agent",type="bind" \
    --env SSH_AUTH_SOCK=/ssh-agent \
    # Mount the current project directory to a patch inside the container.
    --mount src="$(pwd -L)",target="/go/src/${project_path}",type="bind" \
    # Adjust the container workspace to the newly mounted project.
    -w "/go/src/${project_path}" \
    # Run a command.
    asecurityteam/sdcli:v3 go test
```

To make this easier, you can add this function to your .bashrc file (omit the first `--mount` if on Mac):

```bash
sdcli() {
    local cwd
    local gopath
    cwd="$(pwd)"
    gopath="${GOPATH:-~/go}"
    # Remove gopath from the front of the directory path. The resulting
    # path is used to construct a mount point inside the container. For
    # go projects this results in them being placed within the gopath
    # of the container. 
    local project_path=${cwd#"${gopath}/src/"}
    docker run -ti --rm \
        --mount src="${SSH_AUTH_SOCK}",target="/ssh-agent",type="bind" \
        --env "SSH_AUTH_SOCK=/ssh-agent" \
        --mount src="$(pwd -L)",target="/go/src/${project_path}",type="bind" \
        -w "/go/src/${project_path}" \
        asecurityteam/sdcli:v3 "$@"
}
```

which will enable you to call the container like:

```bash
sdcli go test
```
For python tooling, you can call the container with:

```bash
export cwd=$(pwd)
export project_path=${cwd#"${GOPATH}/src/"}
docker run -ti \
    # If Linux, mount and configure SSH inside the container.
    --mount src="${SSH_AUTH_SOCK}",target="/ssh-agent",type="bind" \
    --env SSH_AUTH_SOCK=/ssh-agent \
    # Mount the current project directory to a patch inside the container.
    --mount src="$(pwd -L)",target="/go/src/${project_path}",type="bind" \
    # Adjust the container workspace to the newly mounted project.
    -w "/go/src/${project_path}" \
    # Run a command.
    asecurityteam/sdcli:v3 go lint
```

Or, if you've already added the sdcli bash function to your .bashrc file, you can simply type:

```bash
sdcli go lint
```

<a id="markdown-for-shells-other-than-bash" name="for-shells-other-than-bash"></a>
### For Shells Other than `bash`

In fish shell, you create a `~/.config/fish/functions/sdcli.fish` file with 755
permissions having contents:

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
  docker run --rm \
    --mount src="$cwd",target="/go/src/$project_path",type="bind" \
    -w "/go/src/$project_path" \
    asecurityteam/sdcli:v3 $argv
end
```

Some commands are interactive, but if you run `fish` or shells other than `bash`, you
might see "no TTY for interactive shell" or seemingly inexplicable "project_name [New
Project]: Aborted!".  No worries!  Just run in non-interactive mode by specifying all
args on the command line, like:

```bash
sdcli repo go create -- project_name="myproject" project_description="description" --no-input
```

Or start the Docker image with `/bin/bash` as the entrypoint and run `/usr/bin/sdcli $args`
from within (be sure to set `$cwd` and `$project_path` first):

```bash
docker run -it \
    --entrypoint "/bin/bash" \
    --mount src="$cwd",target="/go/src/$project_path",type="bind" \
    -w "/go/src/$project_path" \
    asecurityteam/sdcli:v3
```

<a id="markdown-adding-commands" name="adding-commands"></a>
## Adding Commands

The top-level `sdcli` script will dispatch commands by accumulating all the
arguments and joining them with an `_` character. For example, `sdcli my feature`
will be converted to `sdcli_my_feature` and executed. To add a new command, drop
an executable file in the `./commands` directory and name it according to how you
want the script to be called.
