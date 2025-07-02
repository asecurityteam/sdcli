FROM golang:1.24.4-bullseye AS base

ENV APT_MAKE_VERSION=4.3-4.1 \
    APT_GCC_VERSION=4:10.2.1-1 \
    APT_GIT_VERSION=1:2.30.2-1+deb11u4 \
    LANG=C.UTF-8

#########################################

FROM base AS system_deps

# Install apt dependencies
RUN apt-get update && \
    apt-get install -y \
    apt-transport-https \
    ca-certificates \
    curl \
    make=${APT_MAKE_VERSION} \
    gcc=${APT_GCC_VERSION} \
    git=${APT_GIT_VERSION} \
    bc \
    jq \
    unzip && \
    apt-get upgrade -y

#########################################

FROM system_deps AS go_deps
# https://marcofranssen.nl/manage-go-tools-via-go-modules
ADD golang/* /go-tools/
ADD defaults/.golangci.yaml /defaults/.golangci.yaml
WORKDIR /go-tools
RUN go mod download && grep _ tools.go | awk -F'"' '{print $2}' | xargs -tI % go install % && cd .. && rm /go-tools/* && rmdir /go-tools
WORKDIR /

#########################################

FROM go_deps AS js_deps

# Install NPM
ADD nodesource.gpg /usr/share/keyrings/
ADD nodesource-apt.list /etc/apt/sources.list.d/nodesource.list
RUN apt-get -y update && apt-get install -y nodejs


#########################################

FROM js_deps AS python_deps

ENV PIPENV_VENV_IN_PROJECT=1

RUN apt-get install -y locales python3-distutils python3-pip
RUN sed -i 's/^# *\(en_US.UTF-8\)/\1/' /etc/locale.gen \
    && locale-gen
RUN python3 -mpip install -U pipenv==2024.1.0
ADD python/* /python/
WORKDIR /python/
# this allows to use advanced features of pipenv while still using pip to install actual requirements globally
RUN pipenv requirements > requirements.txt && python3 -m pip install -U -r requirements.txt && rm /python/* && rmdir /python
WORKDIR /

#########################################

FROM python_deps AS ssh_deps

# Install the bitbucket SSH host
RUN mkdir -p /home/sdcli/.ssh
RUN echo 'bitbucket.org ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==' >> /home/sdcli/.ssh/known_hosts

#########################################

FROM ssh_deps AS user_deps

# Create a non-root user to avoid permissions issues when
# modifying files on the mounted host directories.
RUN groupadd -r sdcli -g 1000 \
    && useradd --no-log-init -r -g sdcli -u 1000 sdcli \
    && chown -R sdcli:sdcli /opt \
    && chown -R sdcli:sdcli /go \
    && chown -R sdcli:sdcli /home/sdcli \
    && chown -R sdcli:sdcli /usr/local

#########################################

FROM user_deps AS docker_cli_deps
# https://docs.docker.com/engine/install/debian/
ENV DOCKER_PACKAGE_VERSION=5:27.3.1-1~debian.11~bullseye
ENV COMPOSE_PLUGIN_PACKAGE_VERSION=2.29.7-1~debian.11~bullseye
ENV COMPOSE_PACKAGE_VERSION=1.25.0-1
# comes from curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o - > docker-archive-keyring.gpg
ADD docker-archive-keyring.gpg /usr/share/keyrings/
ADD docker-apt.list /etc/apt/sources.list.d/docker.list
# we need cli only, not deamon
RUN apt-get update && apt-get -y install docker-ce-cli=${DOCKER_PACKAGE_VERSION} docker-compose-plugin=${COMPOSE_PLUGIN_PACKAGE_VERSION} docker-compose=${COMPOSE_PACKAGE_VERSION} \
    && rm -rf /var/lib/apt/lists/*

#########################################

FROM docker_cli_deps

RUN mkdir -p /home/sdcli/oss-templates/
COPY ./commands/* /usr/bin/

USER sdcli
ENTRYPOINT [ "sdcli" ]
