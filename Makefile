all: build test

deps:
	go get golang.org/x/tools/cmd/cover
	go get github.com/codegangsta/cli

build: deps
	go build

test: build
	@go test -coverprofile=../coverage.out
	@go tool cover -html=../coverage.out -o ../coverage.html

gox:
	go get github.com/mitchellh/gox
	gox -build-toolchain

build-all: test
	@which gox || make gox
	gox -arch="amd64" -os="darwin" -os="linux"

install: build test
	cp main /usr/sbin/hostess
