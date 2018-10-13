BIN = mkrg

all: clean build

build: deps
	go build -o build/$(BIN) ./cmd/...

install: deps
	go install ./...

deps:
	go get -d -v ./...

test: build
	go test -v ./...

lint: lintdeps build
	golint -set_exit_status ./...

lintdeps:
	go get -d -v -t ./...
	command -v golint >/dev/null || go get -u golang.org/x/lint/golint

clean:
	rm -rf build
	go clean

.PHONY: build install deps test lint lintdeps clean
