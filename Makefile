build:
	go build hostess.go

test: build
	./hostess add domain ip

install: build
	cp hostess /usr/sbin/hostess