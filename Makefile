.PHONY: build
build:
	go build -v ./cmd/logserver

.PHONY: run
run: 
	go build -v ./cmd/logserver
	./logserver

.PHONY: test
test:
	go test -v ./...

.DEFAULT_GOAL := build