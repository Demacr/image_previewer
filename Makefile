build:
	docker-compose -f build/docker-compose.yml build

build-push:
	docker-compose -f build/docker-compose.yml build
	docker-compose -f build/docker-compose.yml push

build-local:
	go build  -o image_previewer ./cmd/image_previewer/...

test-unit:
	go test -race ./... -count=10 -cover

test-integration:
	docker-compose -f deployments/docker-compose.test.yml up -d
	./test/test.sh
	export RC=$?
	docker-compose -f deployments/docker-compose.test.yml down
	exit ${RC}

test: test-unit test-integration

run:
	docker-compose -f deployments/docker-compose.yml up

run-local:
	./image_previewer

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.30.0

lint: install-lint-deps
	golangci-lint run ./...

lint-local:
	golangci-lint run ./...

clean:
	rm -rf cache/
	rm image_previewer

.PHONY: build-local
