BIN := "./bin/mainService"
BIN_SS := "./bin/statSender"
DOCKER_IMG="banner_rotation:develop"
DOCKER_IMG_SS="banner_stat_sender:develop"
DOCKER_IMG_INT_TESTS="integration_tests:develop"

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build_main_service:
	go build -v -o $(BIN) -ldflags "$(LDFLAGS)" ./cmd/mainService

run_main_service: build_main_service
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

build-img_int_tests:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG_INT_TESTS) \
		-f integration_tests/Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

runrun-img_ss: build-img_ss
	docker run $(DOCKER_IMG_SS)

run-img_int_tests: build-img_int_tests
	docker run $(DOCKER_IMG_INT_TESTS)


build: build-img build-img_ss 
	docker-compose build

test:
	go test -race -count 100 ./internal/... 

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.46.2

lint: install-lint-deps
	CGO_ENABLED=0 golangci-lint run ./... --config=./.golangci.yml

up: build
	docker-compose up -d postgres_db
	docker-compose up -d rabbitmq
	docker-compose up -d mainSevice
	docker-compose up -d statSender

down:
	docker-compose down

run: up

integration_test:
	docker-compose up integraton_tests
#	go test -tags integration ./integration_tests/...


.PHONY: build run build-img run-img test lint