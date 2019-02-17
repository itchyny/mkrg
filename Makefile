BIN := mkrg

.PHONY: all
all: clean build

.PHONY: build
build: deps
	go build -o build/$(BIN) ./cmd/...

.PHONY: install
install: deps
	go install ./...

.PHONY: deps
deps:
	go get -d -v ./...

.PHONY: test
test: build
	go test -v ./...

.PHONY: lint
lint: lintdeps build
	golint -set_exit_status ./...

.PHONY: lintdeps
lintdeps:
	go get -d -v -t ./...
	command -v golint >/dev/null || go get -u golang.org/x/lint/golint

.PHONY: clean
clean:
	rm -rf build
	go clean
