SHELL := /bin/bash

GO_VERSION := $(shell awk '/^go / {print $$2}' go.mod)
GIT_VERSION := $(shell git describe --always --tags HEAD)
GIT_COMMIT := $(shell git rev-parse HEAD)
GIT_COMMIT_DATE := $(shell git log --pretty=%ct  -1)
GIT_TREE_STATE := $(shell git diff --exit-code --quiet && echo clean || echo dirty)

# Print out only the variables declared in this makefile. Will be used
# by other tools like github workflows or any other build tool.
# note: the ldflags value is long with spaces and isn't a valid bash expression
# unless it is quoted. I did't add quotes to the 'echo' section here because it
# creates other problems when using valus in github action. So to source all
# vars in one-shot just ignore ldflags with a grep or specifically surround it
# single quotes. e.g.  `make print-vars | grep -v ldflags=`
print-vars:
	@$(foreach v,$(sort $(.VARIABLES)), \
		$(if $(filter-out environment% default automatic, $(origin $v)), \
			$(if $(filter-out .%,$(v)), \
				echo $v=$($v);)))

install:
	go install -ldflags="$(ldflags)" .

fmt:
	go fmt ./...

vet:
	go vet ./...

ldflags := -X github.com/rgolangh/pq/internal/version.Version=$(GIT_VERSION)
ldflags += -X github.com/rgolangh/pq/internal/version.Commit=$(GIT_COMMIT)
ldflags += -X github.com/rgolangh/pq/internal/version.CommitDate=$(GIT_COMMIT_DATE)
ldflags += -X github.com/rgolangh/pq/internal/version.TreeState=$(GIT_TREE_STATE)

build: fmt vet
	go build -v -o bin/ -ldflags="$(ldflags)" ./...

