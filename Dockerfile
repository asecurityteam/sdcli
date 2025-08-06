FROM golang:1.24.5-bookworm AS base

ENV LANG=en_US.UTF-8
ENV DEBIAN_FRONTEND=noninteractive

COPY apt.conf /etc/apt/apt.conf.d/99-sdcli-local

#########################################

FROM base AS system_deps

# Install apt dependencies
RUN apt-get update && \
    apt-get install \
    apt-transport-https \
    ca-certificates \
    curl \
    make \
    git \
    bc \
    jq \
    yamllint \  
    locales\  
    unzip && \
    apt-get upgrade && \
    rm -rf /var/lib/apt/lists/* && \
    # Generate required locale
    sed -i '/en_US.UTF-8/s/^# //' /etc/locale.gen && \
    locale-gen && \
    update-locale LANG=en_US.UTF-8

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
RUN apt-get update && apt-get install nodejs && rm -rf /var/lib/apt/lists/*

#########################################

FROM js_deps AS ssh_deps

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
ENV DOCKER_PACKAGE_VERSION=5:27.5.1-1~debian.12~bookworm
ENV COMPOSE_PLUGIN_PACKAGE_VERSION=2.33.1-1~debian.12~bookworm
ENV COMPOSE_PACKAGE_VERSION=1.29.2-3
# comes from curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o - > docker-archive-keyring.gpg
ADD docker-archive-keyring.gpg /usr/share/keyrings/
ADD docker-apt.list /etc/apt/sources.list.d/docker.list
# we need cli only, not deamon
RUN apt-get update && apt-get install docker-ce-cli=${DOCKER_PACKAGE_VERSION} \
    docker-compose-plugin=${COMPOSE_PLUGIN_PACKAGE_VERSION} \
    docker-compose=${COMPOSE_PACKAGE_VERSION} \
    && rm -rf /var/lib/apt/lists/*

#########################################

FROM docker_cli_deps

COPY ./commands/* /usr/bin/

USER sdcli
ENTRYPOINT [ "sdcli" ]