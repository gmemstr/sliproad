.DEFAULT_GOAL := build
SLIPROAD_VERSION := 2.0.0
# Workaround for CircleCI Docker image and mkdir.
SHELL := /bin/bash

make_build_dir:
	mkdir -p build/{bin,assets,tars}

build: make_build_dir
	go build -o build/bin/sliproad

pi: make_build_dir
	env GOOS=linux GOARCH=arm GOARM=5 go build -o build/bin/sliproad-arm

small: make_build_dir
	go build -o build/bin/sliproad -ldflags="-s -w"
	upx --brute build/bin/sliproad -9 --no-progress

small_pi: make_build_dir
	env GOOS=linux GOARCH=arm GOARM=5 go build -o build/bin/sliproad-arm -ldflags="-s -w"
	upx --brute build/bin/sliproad-arm -9 --no-progress

run:
	go run webserver.go

test:
	go test ./... -cover

dist: clean make_build_dir small small_pi
	cp -r assets/* build/assets
	tar -czf build/tars/sliproad-$(SLIPROAD_VERSION)-arm.tar.gz build/assets build/bin/sliproad-arm README.md LICENSE
	tar -czf build/tars/sliproad-$(SLIPROAD_VERSION)-x86.tar.gz build/assets build/bin/sliproad README.md LICENSE

clean:
	rm -rf build
