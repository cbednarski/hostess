RELEASE_VERSION=$(shell git describe --tags)
PREFIX ?= /usr/local

test:
	go test ./...
	go vet ./...

install:
	go build -o bin/hostess .
	install -C bin/hostess -t ${PREFIX}/bin

release: test
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${RELEASE_VERSION}" -o bin/hostess_windows_amd64.exe .
	GOOS=darwin  GOARCH=amd64 go build -ldflags "-X main.Version=${RELEASE_VERSION}" -o bin/hostess_macos_amd64 .
	GOOS=linux   GOARCH=amd64 go build -ldflags "-X main.Version=${RELEASE_VERSION}" -o bin/hostess_linux_amd64 .
	GOOS=linux   GOARCH=arm   go build -ldflags "-X main.Version=${RELEASE_VERSION}" -o bin/hostess_linux_arm .

clean:
	rm -rf ./bin/

.PHONY: install test release clean
