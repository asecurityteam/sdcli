.PHONY : dep lint test integration coverage doc build run deploy

DOCKER_TAG := $(shell cat version)
DIR := $(shell pwd -L)
ENVIRON := dev
IMAGE_NAME := golang
VERSION := $(shell cat version)
IMAGE_PATH := /asecurityteam
ifeq ($(ENVIRON),prod)
IMAGE_PATH := /sox/asecurityteam
endif
REGISTRY := docker.atl-paas.net
REGISTRY_USER := $(shell whoami)
REGISTRY_PWD := ""
ARTIFACT := $(REGISTRY)$(IMAGE_PATH)/$(IMAGE_NAME)

dep: ;

lint: ;

test: ;

integration: ;

coverage: ;

doc: ;

build:
ifeq ($(REGISTRY),docker-proxy.services.atlassian.com)
	docker login -u=$(REGISTRY_USER) -p=$(REGISTRY_PWD) $(REGISTRY)
endif
	docker build -t $(ARTIFACT):$(VERSION) .

run:
	docker run -ti $(ARTIFACT):$(VERSION)

deploy: build
ifeq ($(REGISTRY),docker-proxy.services.atlassian.com)
	docker login -u=$(REGISTRY_USER) -p=$(REGISTRY_PWD) $(REGISTRY)
endif
	docker push $(ARTIFACT):$(VERSION)
	docker tag $(ARTIFACT):$(VERSION) $(ARTIFACT):latest
	docker push $(ARTIFACT):latest
