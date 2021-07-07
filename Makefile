ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

.PHONY: clean
clean: go-clean
	rm -rf ./bin
	rm -f ./testing/server.jar

.PHONY: build
build: bin/mineshaft

.PHONY: test
test: build
	./bin/mineshaft -f ./.mineshaft.toml

##
##
## Golang
##
##
# All the go files in the ROOT_DIR
GO_FILES := $(shell find ${ROOT_DIR} -type f -name '*.go' ! -name '*_test.go')

# All the go files
GOLANG ?= $(shell which go)
GO_ENV :=
GO_BUILD_FLAGS :=
GO_BUILD ?= ${GO_ENV} ${GOLANG} build ${GO_BUILD_FLAGS}

.PHONY: go-clean
go-clean:
	${GOLANG} clean -testcache
	${GOLANG} mod tidy

.PHONY: go-test
go-test:
	${GOLANG} test -v ./...

bin/mineshaft: ${GO_FILES} go.mod
	cd ${ROOT_DIR}/cmd/mineshaft && ${GO_BUILD} -o ${ROOT_DIR}/bin/mineshaft .
