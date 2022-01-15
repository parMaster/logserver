.PHONY: build
build:
	go build -v ./cmd/logserver

.PHONY: run
run: 
	go build -v ./cmd/logserver
	./logserver.exe


.DEFAULT_GOAL := build