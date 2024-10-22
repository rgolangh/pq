SHELL: /bin/bash

install:
	go install -ldflags="$(ldflags)" .

fmt:
	go fmt ./...

vet:
	go vet ./...

ldflags := -X github.com/rgolangh/pq/internal/version.Version=$(shell git describe --all HEAD)
ldflags += -X github.com/rgolangh/pq/internal/version.Commit=$(shell git rev-parse HEAD)
ldflags += -X github.com/rgolangh/pq/internal/version.CommitDate=$(shell date --iso-8601=seconds)
ldflags += -X github.com/rgolangh/pq/internal/version.TreeState=$(shell git diff --exit-code --quiet && echo clean || echo dirty)

build: fmt vet
	go build -o bin/ -ldflags="$(ldflags)" ./...
