.PHONY: test exec cover devtools env cloc

GOMOD:=github.com/sunmxt/arbiter

PROJECT_ROOT:=$(shell pwd)
BUILD_DIR:=build
REVISION:=$(shell git rev-parse HEAD || cat REVISION)

ifeq ($(USE_GLOBAL_GOPATH),)
export GOPATH:=$(PROJECT_ROOT)/$(BUILD_DIR)
endif

export PATH:=$(GOPATH)/bin:$(PATH)

COVERAGE_DIR:=coverage

all: cover

$(COVERAGE_DIR):
	mkdir -p $(COVERAGE_DIR)

build:
	mkdir build

bin: build/bin
	mkdir bin

env:
	@echo "export PROJECT_ROOT=\"$(PROJECT_ROOT)\""
	@echo "export GOPATH=\"\$${PROJECT_ROOT}/build\""
	@echo "export PATH=\"\$${PROJECT_ROOT}/bin:\$${PATH}\""

cover: coverage test
	go tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html

test: coverage
	go test -v -coverprofile=$(COVERAGE_DIR)/coverage.out -cover ./
	go tool cover -func=$(COVERAGE_DIR)/coverage.out

cloc:
	cloc . --exclude-dir=ci,mocks

devtools: $(GOPATH)/bin/gopls $(GOPATH)/bin/goimports

exec:
	$(CMD)

build/bin: bin build
	test -d build/bin || ln -s $$(pwd)/bin build/bin

$(GOPATH)/bin/gopls: bin
	go get -u golang.org/x/tools/gopls

$(GOPATH)/bin/goimports: bin
	go get -u golang.org/x/tools/cmd/goimports