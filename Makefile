.DEFAULT_GOAL := build

fmt: 
	go fmt ./...
.PHONY:fmt

vet: fmt
	go vet ./...
.PHONY:vet

build: vet
	go build -o $(CURDIR)/bin/myhttp .
.PHONY:build

test:
	go test ./...
.PHONY:test