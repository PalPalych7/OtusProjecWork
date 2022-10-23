BIN := "./bin/mainService"
DOCKER_IMG="banner_rotation:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/mainService

run: build
	$(BIN) -config ./configs/config.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

build-compose:
	docker-compose build

run-compose: build-compose
	docker-compose up -d --build

#version: build
#	$(BIN) version
#
test:
	go test -race ./internal/... 

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.41.1

lint: install-lint-deps
	CGO_ENABLED=0 golangci-lint run ./...


grpc_generate:
	protoc --proto_path=api --go_out=pb/ --go-grpc_out=pb/ api/*.proto

.PHONY: build run build-img run-img version test lint