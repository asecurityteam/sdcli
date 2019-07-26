.PHONY : dep lint test integration coverage run

IMAGE_NAME := {{cookiecutter.project_slug}}
DIR := $(shell pwd -L)
VERSION := $(shell cat version)
GOPATH := ${GOPATH}
ifeq ($(GOPATH),)
	GOPATH := ${HOME}/go
endif
PROJECT_PATH := $(subst $(GOPATH)/src/,,$(DIR))
REGISTRY := registry.hub.docker.com
REGISTRY_USER := $(shell whoami)
REGISTRY_PWD := ""
IMAGE_PATH := /asecurityteam
ARTIFACT := $(REGISTRY)$(IMAGE_PATH)/$(IMAGE_NAME)



dep:
	docker run -ti \
        --mount src="$(DIR)",target="/go/src/$(PROJECT_PATH)",type="bind" \
        -w "/go/src/$(PROJECT_PATH)" \
        registry.hub.docker.com/asecurityteam/sdcli:v1 python dep

lint:
	docker run -ti \
        --mount src="$(DIR)",target="/go/src/$(PROJECT_PATH)",type="bind" \
        -w "/go/src/$(PROJECT_PATH)" \
        registry.hub.docker.com/asecurityteam/sdcli:v1 python lint

test:
	docker run -ti \
        --mount src="$(DIR)",target="/go/src/$(PROJECT_PATH)",type="bind" \
        -w "/go/src/$(PROJECT_PATH)" \
        registry.hub.docker.com/asecurityteam/sdcli:v1 python test

integration: ;

coverage:
	docker run -ti \
        --mount src="$(DIR)",target="/go/src/$(PROJECT_PATH)",type="bind" \
        -w "/go/src/$(PROJECT_PATH)" \
        registry.hub.docker.com/asecurityteam/sdcli:v1 python coverage

doc: ;

docker-login: ;

build:
	python3 setup.py bdist_wheel

run: ;

deploy: build
	python3 -m twine upload dist/*

clean: ;
