FROM golang:1.17.7-buster AS BASE

ENV APT_MAKE_VERSION=4.2.1-1.2 \
    APT_GCC_VERSION=4:8.3.0-1 \
    APT_GIT_VERSION=1:2.20.1-2+deb10u3 \
    GOLANGCI_VERSION=v1.39.0 \
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
    jq && \
    apt-get upgrade -y

#########################################

FROM system_deps AS go_deps

# Install gocov tools
RUN go get github.com/axw/gocov/... && \
    go install github.com/axw/gocov/gocov@latest && \
    go get github.com/AlekSi/gocov-xml && \
    go install github.com/AlekSi/gocov-xml@latest && \
    go get github.com/wadey/gocovmerge && \
    go install github.com/wadey/gocovmerge@latest

# Install lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b ${GOPATH}/bin ${GOLANGCI_VERSION}

#########################################

FROM go_deps AS js_deps

# Install NPM
RUN curl -sfL https://deb.nodesource.com/setup_12.x | bash - && \
    apt-get install -y nodejs

#########################################

FROM js_deps AS python_deps

RUN apt-get install -y locales python3-distutils python3-pip
RUN sed -i 's/^# *\(en_US.UTF-8\)/\1/' /etc/locale.gen \
    && locale-gen
RUN pip3 install -U setuptools cookiecutter flake8 coverage pytest pytest-cov pipenv oyaml python-slugify yamllint

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
ENV DOCKER_PACKAGE_VERSION=5:20.10.7~3-0~debian-buster
ENV COMPOSE_PACKAGE_VERSION=1.21.0-3
# comes from curl -fsSL https://download.docker.com/linux/debian/gpg | gpg --dearmor -o - > docker-archive-keyring.gpg
ADD docker-archive-keyring.gpg /usr/share/keyrings/
ADD docker-apt.list /etc/apt/sources.list.d/docker.list
# we need cli only, not deamon
RUN apt-get update && apt-get -y install docker-ce-cli=${DOCKER_PACKAGE_VERSION} docker-compose=${COMPOSE_PACKAGE_VERSION} && rm -rf /var/lib/apt/lists/*

#########################################

FROM docker_cli_deps

RUN mkdir -p /home/sdcli/oss-templates/

COPY ./oss-templates/ /home/sdcli/oss-templates/

COPY ./commands/* /usr/bin/

USER sdcli
ENTRYPOINT [ "sdcli" ]
