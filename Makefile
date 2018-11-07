cwd := $(shell pwd -L)
VERSION := v1

.PHONY : dep lint test integration coverage doc build run deploy clean

.builds/:
	mkdir -p .builds/

.builds/bin/: | .builds/
	mkdir -p .builds/bin/
.builds/pkg/: | .builds/
	mkdir -p .builds/pkg/
.builds/src/: | .builds/
	mkdir -p .builds/src/
.builds/gopath: | .builds/bin/ .builds/pkg/ .builds/src/
	touch .builds/gopath

.builds/bin/dep: | .builds/gopath
	rm -rf .builds/src/github.com/golang/dep/
	mkdir -p .builds/src/github.com/golang/dep/
	git clone --single-branch --branch="v0.4.1" https://github.com/golang/dep.git .builds/src/github.com/golang/dep
	cd .builds/src/github.com/golang/dep/cmd/dep/ && GOPATH=$(cwd)/.builds/ go get -u
	cd .builds/src/github.com/golang/dep/cmd/dep/ && GOPATH=$(cwd)/.builds/ go build -o dep .
	mv .builds/src/github.com/golang/dep/cmd/dep/dep .builds/bin/
	chmod +x .builds/bin/dep

.builds/bin/golangci-lint: | .builds/bin/dep
	rm -rf .builds/src/github.com/golangci/golangci-lint/
	mkdir -p .builds/src/github.com/golangci/golangci-lint/
	git clone --single-branch --branch="v1.10.2" https://github.com/golangci/golangci-lint.git .builds/src/github.com/golangci/golangci-lint
	cd .builds/src/github.com/golangci/golangci-lint/ && GOPATH=$(cwd)/.builds/ $(cwd)/.builds/bin/dep ensure
	cd .builds/src/github.com/golangci/golangci-lint/cmd/golangci-lint/ && GOPATH=$(cwd)/.builds/ go build -o golangci-lint .
	mv .builds/src/github.com/golangci/golangci-lint/cmd/golangci-lint/golangci-lint .builds/bin/
	chmod +x .builds/bin/golangci-lint

.builds/bin/gocov: | .builds/gopath
	rm -rf .builds/src/github.com/axw/gocov/
	mkdir -p .builds/src/github.com/axw/gocov/
	git clone --single-branch --branch="master" https://github.com/axw/gocov.git .builds/src/github.com/axw/gocov
	cd .builds/src/github.com/axw/gocov/gocov/ && GOPATH=$(cwd)/.builds/ go get -u
	cd .builds/src/github.com/axw/gocov/gocov/ && GOPATH=$(cwd)/.builds/ go build -o gocov .
	mv .builds/src/github.com/axw/gocov/gocov/gocov .builds/bin/
	chmod +x .builds/bin/gocov

.builds/bin/gocovmerge: | .builds/gopath
	rm -rf .builds/src/github.com/wadey/gocovmerge/
	mkdir -p .builds/src/github.com/wadey/gocovmerge/
	git clone --single-branch --branch="master" https://github.com/wadey/gocovmerge.git .builds/src/github.com/wadey/gocovmerge
	cd .builds/src/github.com/wadey/gocovmerge/ && GOPATH=$(cwd)/.builds/ go get -u
	cd .builds/src/github.com/wadey/gocovmerge/ && GOPATH=$(cwd)/.builds/ go build -o gocovmerge .
	mv .builds/src/github.com/wadey/gocovmerge/gocovmerge .builds/bin/
	chmod +x .builds/bin/gocovmerge

.builds/bin/gocov-xml: | .builds/gopath
	rm -rf .builds/src/github.com/AlekSi/gocov-xml/
	mkdir -p .builds/src/github.com/AlekSi/gocov-xml/
	git clone --single-branch --branch="master" https://github.com/AlekSi/gocov-xml.git .builds/src/github.com/AlekSi/gocov-xml
	cd .builds/src/github.com/AlekSi/gocov-xml/ && GOPATH=$(cwd)/.builds/ go get -u
	cd .builds/src/github.com/AlekSi/gocov-xml/ && GOPATH=$(cwd)/.builds/ go build -o gocov-xml .
	mv .builds/src/github.com/AlekSi/gocov-xml/gocov-xml .builds/bin/
	chmod +x .builds/bin/gocov-xml

.builds/all: | .builds/bin/dep .builds/bin/golangci-lint .builds/bin/gocov .builds/bin/gocovmerge .builds/bin/gocov-xml
	touch .builds/all

dep: .builds/all
	PATH="${PATH}:$(cwd)/.builds/bin" dep ensure

lint: .builds/all
	PATH="${PATH}:$(cwd)/.builds/bin" golangci-lint run --config .golangci.yaml ./...

test: .builds/all
	mkdir -p .coverage
	go test -v -cover -coverpkg=./... -coverprofile=.coverage/unit.cover.out ./...
	PATH="${PATH}:$(cwd)/.builds/bin" gocov convert .coverage/unit.cover.out | PATH="${PATH}:$(cwd)/.builds/bin" gocov-xml > .coverage/unit.xml

integration: .builds/all
	mkdir -p .coverage
	go test -tags="integration" -v -cover -coverpkg=./... -coverprofile=.coverage/integration.cover.out ./tests/
	PATH="${PATH}:$(cwd)/.builds/bin" gocov convert .coverage/integration.cover.out | PATH="${PATH}:$(cwd)/.builds/bin" gocov-xml > .coverage/integration.xml

coverage: .builds/all
	mkdir -p .coverage
	PATH="${PATH}:$(cwd)/.builds/bin" gocovmerge .coverage/*.cover.out > .coverage/combined.cover.out
	PATH="${PATH}:$(cwd)/.builds/bin" gocov convert .coverage/combined.cover.out | PATH="${PATH}:$(cwd)/.builds/bin" gocov-xml > .coverage/combined.xml
	go tool cover -func .coverage/combined.cover.out

doc:
	godoc -http ':8080'

build:
	mkdir -p .builds
	GOOS=linux go build -o .builds/linux .
	GOOS=darwin go build -o .builds/osx .

run:
	go run main.go

deploy:
	curl -X POST -H "Authorization: Token ${STATLAS_TOKEN}" -T ".builds/linux" https://statlas.prod.atl-paas.net/security-development/sdcli/$(VERSION)/linux/sdcli
	curl -X POST -H "Authorization: Token ${STATLAS_TOKEN}" -T ".builds/osx" https://statlas.prod.atl-paas.net/security-development/sdcli/$(VERSION)/osx/sdcli

clean:
	rm -rf .builds/