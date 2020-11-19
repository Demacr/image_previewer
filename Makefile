build:
	docker-compose -f build/docker-compose.yml build

build-push:
	docker-compose -f build/docker-compose.yml build
	docker-compose -f build/docker-compose.yml push

build-local:
	go build  -o image_previewer ./cmd/image_previewer/...

test-unit:
	go test -race ./... -count=100 -cover

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

lint:
	golangci-lint run ./...

clean:
	rm -rf cache/
	rm image_previewer

.PHONY: build build-local build-push test-unit test-integration test run run-local lint clean
