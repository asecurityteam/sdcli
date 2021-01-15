FROM golang:1.13.8 AS BASE

ENV APT_MAKE_VERSION=4.2.1-1.2 \
    APT_GCC_VERSION=4:8.3.0-1 \
    APT_GIT_VERSION=1:2.20.1-2+deb10u1 \
    GOLANGCI_VERSION=v1.33.0 \
    LANG=C.UTF-8

#########################################

FROM BASE AS SYSTEM_DEPS

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

FROM SYSTEM_DEPS AS GO_DEPS

# Install dep
RUN curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Install gocov tools
RUN go get github.com/axw/gocov/... && \
    go install github.com/axw/gocov/gocov && \
    go get github.com/AlekSi/gocov-xml && \
    go install github.com/AlekSi/gocov-xml && \
    go get github.com/wadey/gocovmerge && \
    go install github.com/wadey/gocovmerge

# Install lint
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b ${GOPATH}/bin ${GOLANGCI_VERSION}

#########################################

FROM GO_DEPS AS JS_DEPS

# Install NPM
RUN curl -sfL https://deb.nodesource.com/setup_12.x | bash - && \
    apt-get install -y nodejs

#########################################

FROM JS_DEPS AS PYTHON_DEPS

RUN apt-get install -y locales python3-distutils
RUN curl https://bootstrap.pypa.io/get-pip.py | python3
RUN pip3 install -U setuptools cookiecutter
RUN sed -i 's/^# *\(en_US.UTF-8\)/\1/' /etc/locale.gen \
    && locale-gen
RUN pip3 install -U flake8

RUN pip3 install coverage
RUN pip3 install pytest
RUN pip3 install pytest-cov
RUN pip3 install pipenv
RUN pip3 install oyaml
RUN pip3 install python-slugify
RUN pip3 install --upgrade git+git://github.com/asecurityteam/ccextender
RUN pip3 install yamllint

#########################################

FROM PYTHON_DEPS AS SSH_DEPS

# Install the bitbucket SSH host
RUN mkdir -p /home/sdcli/.ssh
RUN echo 'bitbucket.org ssh-rsa AAAAB3NzaC1yc2EAAAABIwAAAQEAubiN81eDcafrgMeLzaFPsw2kNvEcqTKl/VqLat/MaB33pZy0y3rJZtnqwR2qOOvbwKZYKiEO1O6VqNEBxKvJJelCq0dTXWT5pbO2gDXC6h6QDXCaHo6pOHGPUy+YBaGQRGuSusMEASYiWunYN0vCAI8QaXnWMXNMdFP3jHAJH0eDsoiGnLPBlBp4TNm6rYI74nMzgz3B9IikW4WVK+dc8KZJZWYjAuORU3jc1c/NPskD2ASinf8v3xnfXeukU0sJ5N6m5E8VLjObPEO+mN2t/FZTMZLiFqPWc/ALSqnMnnhwrNi2rbfg/rd/IpL8Le3pSBne8+seeFVBoGqzHM9yXw==' >> /home/sdcli/.ssh/known_hosts

#########################################

FROM SSH_DEPS AS USER_DEPS

# Create a non-root user to avoid permissions issues when
# modifying files on the mounted host directories.
RUN groupadd -r sdcli -g 1000 \
    && useradd --no-log-init -r -g sdcli -u 1000 sdcli \
    && chown -R sdcli:sdcli /opt \
    && chown -R sdcli:sdcli /go \
    && chown -R sdcli:sdcli /home/sdcli \
    && chown -R sdcli:sdcli /usr/local

#########################################

FROM USER_DEPS

USER sdcli

RUN mkdir -p /home/sdcli/oss-templates/

COPY ./oss-templates/ /home/sdcli/oss-templates/

COPY ./commands/* /usr/bin/

ENTRYPOINT [ "sdcli" ]
