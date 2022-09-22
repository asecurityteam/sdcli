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
	docker run --rm -i -v "$(DIR):/mnt:ro" koalaman/shellcheck:v0.8.0 commands/*

test:
	docker build -t local/test/sdcli .
	docker build -t local/test/sdclitests test
	docker run -i local/test/sdclitests
	docker tag local/test/sdcli sdcli

integration: ;

coverage: ;

doc: ;

build:
	docker build -t $(ARTIFACT) .

run:
	docker run -i $(ARTIFACT)
