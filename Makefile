all: build test

deps:
	go get github.com/codegangsta/cli

build: deps
	go build

test: build
	@cd lib && go test -coverprofile=../coverage.out
	@cd lib && go tool cover -html=../coverage.out -o ../coverage.html

install: build test
	cp main /usr/sbin/hostess