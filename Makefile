all: test install

# Note: this command will be slightly messy UNLESS we are on a tag, which is
# what we want.
RELEASE_VERSION=$(shell git describe --tags)

install:
	go install .

deps:
	go install golang.org/x/lint/golint
	go install golang.org/x/tools/cmd/cover

test: deps
	go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o coverage.html
	go vet $(PACKAGES)
	golint $(PACKAGES)

release: test
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.Version=${RELEASE_VERSION}" -o bin/hostess_windows_amd64.exe .
	GOOS=darwin  GOARCH=amd64 go build -ldflags "-X main.Version=${RELEASE_VERSION}" -o bin/hostess_macos_amd64 .
	GOOS=linux   GOARCH=amd64 go build -ldflags "-X main.Version=${RELEASE_VERSION}" -o bin/hostess_linux_amd64 .
	GOOS=linux   GOARCH=arm   go build -ldflags "-X main.Version=${RELEASE_VERSION}" -o bin/hostess_linux_arm .

clean:
	rm -f ./coverage.*
	rm -rf ./bin/

.PHONY: all install deps test release clean
