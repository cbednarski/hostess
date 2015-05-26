all: build test

deps:
	go get github.com/golang/lint/golint
	go get github.com/stretchr/testify/assert
	go get

build: deps
	go build cmd/hostess/hostess.go

test:
	go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o coverage.html
	go vet
	golint

gox:
	go get github.com/mitchellh/gox
	gox -build-toolchain

build-all: test
	which gox || make gox
	gox -arch="amd64" -os="darwin" -os="linux" github.com/cbednarski/hostess/cmd/hostess

install: hostess
	cp hostess /usr/local/bin/hostess

clean:
	rm -f ./hostess
	rm -f ./hostess_*
	rm -f ./coverage.*
