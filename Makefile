args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

MODULE=github.com/bluecolor/tractor

.PHONY: build run migrate start

build:
	go build -o bin/tractor $(MODULE)/cmd/tractor

test:
	go test ./...

server-start:
	go run cmd/tractor/main.go server start

