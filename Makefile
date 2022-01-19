args = `arg="$(filter-out $@,$(MAKECMDGOALS))" && echo $${arg:-${1}}`

MODULE=github.com/bluecolor/tractor

.PHONY: build run migrate start

build:
	go build -o bin/tractor $(MODULE)/cmd/tractor

test:
	go test -v ./...

run_incsv_outcsv:
	go run cmd/tractor/main.go run \
		--config /home/ceyhun/projects/tractor/config/examples/incsv_outcsv.yml \
		--mapping Demo

run_incsv_outjson:
	go run cmd/tractor/main.go run \
		--config /home/ceyhun/projects/tractor/config/examples/incsv_outjson.yml \
		--mapping Demo

