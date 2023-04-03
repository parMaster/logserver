.PHONY: build, run, test
build:
	go build -v ./cmd/logserver

run: 
	go build -v ./cmd/logserver
	./logserver

test:
	go test -v ./...

.DEFAULT_GOAL: build