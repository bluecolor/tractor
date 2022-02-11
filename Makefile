args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

MODULE=github.com/bluecolor/tractor

.PHONY: build run migrate start

build:
	go build -o bin/tractor $(MODULE)/cmd/tractor

test:
	go test -v ./... -count=1 -failfast
db-drop:
	go run cmd/tractor/main.go db drop

db-migrate:
	go run cmd/tractor/main.go db migrate

db-seed:
	go run cmd/tractor/main.go db seed

db-reset:
	go run cmd/tractor/main.go db reset

server-start:
	go run cmd/tractor/main.go server start

server-reset:
	go run cmd/tractor/main.go db reset
	go run cmd/tractor/main.go server start
