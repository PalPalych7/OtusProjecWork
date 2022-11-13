BIN := "./bin/mainService"
BIN_SS := "./bin/statSender"
DOCKER_IMG="banner_rotation:develop"
DOCKER_IMG_SS="banner_stat_sender:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/mainService

run: build
	$(BIN) -config ./configs/config.toml

build_ss:
	go build -v -o $(BIN_SS) -ldflags "$(LDFLAGS)" ./cmd/statSender

run_ss: build_ss
	$(BIN_SS) -config ./configs/statSenderConfig.toml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/mainService/Dockerfile .

build-img_ss:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_SS) \
		-f build/statSender/Dockerfile .


run-img: build-img
	docker run $(DOCKER_IMG)

run-img_ss: build-img_ss
	docker run $(DOCKER_IMG_SS)


build-compose: build-img build-img_ss 
	docker-compose build

test:
	go test -race -count 100 ./internal/... 

integration_test:
	go test -tags integration ./integration_tests/...

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.46.2

lint: install-lint-deps
	CGO_ENABLED=0 golangci-lint run ./... --config=./.golangci.yml

up: build-compose
	docker-compose up -d --build

down:
	docker-compose down

.PHONY: build run build-img run-img test lint