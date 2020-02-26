all: build test

deps:
	go get golang.org/x/lint/golint
	go get golang.org/x/tools/cmd/cover
	go get

build: deps
	go build cmd/hostess/hostess.go

test:
	go test -coverprofile=coverage.out; go tool cover -html=coverage.out -o coverage.html
	go vet $(PACKAGES)
	golint $(PACKAGES)

build-all: test

	echo FIXME
	exit 1

install:
	go install .

clean:
	rm -f ./hostess
	rm -f ./hostess_*
	rm -f ./coverage.*

.PHONY: all deps build test gox build-all install clean
