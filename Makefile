build:
	go build hostess.go
	go build cmd/main.go

test: build
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

install: build
	cp hostess /usr/sbin/hostess