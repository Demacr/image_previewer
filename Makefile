build:
	go build  -o image_previewer ./cmd/image_previewer/...

test:
	go test -race ./... -count=10 -cover

run:
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

.PHONY: build
