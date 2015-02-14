all: build test

deps:
	go get github.com/codegangsta/cli

build: deps
	go build hostess.go
	go build cmd/main.go

test: build
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

install: build test
	cp main /usr/sbin/hostess