SHELL := /bin/bash

# I try to make this makefile generic enough so other go project can just use
# it as is. Any project specific info is extracted from the real projct assests,
# like go.mod or git, so when a change is needed, it is changed in one place.
GO_VERSION := $(shell awk '/^go / {print $$2}' go.mod)
# RPM doesn't support dashes '-' in the version string
RPM_VERSION = $(shell git describe --always --tags HEAD | sed 's/-/~/g' || true)
GO_MODULE := $(shell awk '/^module / {print $$2}' go.mod)
GO_MODULE_BASE := $(notdir $(GO_MODULE))
GIT_VERSION = $(shell git describe --always --tags HEAD)
GIT_COMMIT = $(shell git rev-parse HEAD)
GIT_COMMIT_DATE = $(shell git log --pretty=%ct  -1)
GIT_TREE_STATE = $(shell git diff --exit-code --quiet && echo clean || echo dirty)
ldflags := -X $(GO_MODULE)/internal/version.Version=$(GIT_VERSION)
ldflags += -X $(GO_MODULE)/internal/version.Commit=$(GIT_COMMIT)
ldflags += -X $(GO_MODULE)/internal/version.CommitDate=$(GIT_COMMIT_DATE)
ldflags += -X $(GO_MODULE)/internal/version.TreeState=$(GIT_TREE_STATE)

INSTALL_PATH = bin/

# Print out only the variables declared in this makefile(not any of the builtins).
# Will be used by other tools like github workflows or any other build tool.
# Note - the ldflags value is long with spaces and isn't a valid bash expression
# unless it is quoted. I did't add quotes to the 'echo' section here because it
# creates other problems when using valus in github action. So to source all
# vars in one-shot just ignore ldflags with a grep or specifically surround it
# single quotes. e.g.  `make print-vars | grep -v ldflags=`
print-vars:
	@$(foreach v,$(sort $(.VARIABLES)), \
		$(if $(filter-out environment% default automatic, $(origin $v)), \
			$(if $(filter-out .%,$(v)), \
				echo $v=$($v);)))

tarball:
	git archive --format=tar.gz --prefix=$(GO_MODULE_BASE)-$(RPM_VERSION)/ -o $(GO_MODULE_BASE)-$(RPM_VERSION).tar.gz HEAD

install:
	go install -ldflags="$(ldflags)" .

fmt:
	formatted="$(shell go fmt ./...)"
	test -z "$$formatted"

vet:
	go vet ./...

build: fmt vet
	go build -v -o $(INSTALL_PATH) -ldflags="$(ldflags)" ./...

rpm: tarball
	mv -v $(GO_MODULE_BASE)-$(RPM_VERSION).tar.gz ~/rpmbuild/SOURCES
	rpmbuild -ba pq.spec --verbose
