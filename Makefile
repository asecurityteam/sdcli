.PHONY : dep lint test integration coverage doc build run deploy

DIR := $(shell pwd -L)
ENVIRON := dev
IMAGE_NAME := sdcli
VERSION := $(shell cat version)
IMAGE_PATH := /asecurityteam
REGISTRY := registry.hub.docker.com
REGISTRY_USER := $(shell whoami)
REGISTRY_PWD := ""
ARTIFACT := $(REGISTRY)$(IMAGE_PATH)/$(IMAGE_NAME)

dep: ;

lint:
	docker run --rm -i -v "$(DIR):/mnt:ro" koalaman/shellcheck:v0.6.0 commands/*

test:
	./tests/run_tests.sh

integration: ;

coverage: ;

doc: ;

build:
	docker build -t $(ARTIFACT) .

run:
	docker run -ti $(ARTIFACT)
