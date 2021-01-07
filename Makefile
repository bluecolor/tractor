OUT := bin/tractor
PKG := github.com/bluecolor/tractor
VERSION := $(shell git describe --always --long --dirty)
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/)
INPUT_PLUGINS_PATH=bin/plugins/input
OUTPUT_PLUGINS_PATH=bin/plugins/output


all: run

build:
	go build -i -v -o ${OUT} -ldflags="-X github.com/bluecolor/tractor/cmd.version=${VERSION}" ${PKG}

build-plugins:
    go build -buildmode=plugin -o ${INPUT_PLUGINS_PATH}/oracle.so plugins/input/oracle/main.go
    go build -buildmode=plugin -o ${OUTPUT_PLUGINS_PATH}/oracle.so plugins/output/oracle/main.go

test:
	@go test -short ${PKG_LIST}

vet:
	@go vet ${PKG_LIST}

lint:
	@for file in ${GO_FILES} ;  do \
		golint $$file ; \
	done

static: vet lint
	go build -i -v -o ${OUT}-v${VERSION} -tags netgo -ldflags="-extldflags \"-static\" -w -s -X cmd.version=${VERSION}" ${PKG}

run: server
	./${OUT}

clean:
	-@rm ${OUT} ${OUT}-v*

.PHONY: run server static vet lint