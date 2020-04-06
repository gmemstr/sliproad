.DEFAULT_GOAL := build
NAS_VERSION := 2.0.0

make_build_dir:
	mkdir -p build/{bin,assets,tars}

build: make_build_dir
	go build -o build/bin/nas

pi: make_build_dir
	env GOOS=linux GOARCH=arm GOARM=5 go build -o build/bin/nas-arm

small: make_build_dir
	go build -o build/bin/nas -ldflags="-s -w"
	upx --brute build/bin/nas -9 --no-progress

small_pi: make_build_dir
	env GOOS=linux GOARCH=arm GOARM=5 go build -o build/bin/nas-arm -ldflags="-s -w"
	upx --brute build/bin/nas-arm -9 --no-progress

run:
	go run webserver.go

test:
	go test ./... -cover

dist: clean make_build_dir small small_pi
	cp -r assets/* build/assets
	tar -czf build/tars/nas-$(NAS_VERSION)-arm.tar.gz build/assets build/bin/nas-arm README.md LICENSE
	tar -czf build/tars/nas-$(NAS_VERSION)-x86.tar.gz build/assets build/bin/nas README.md LICENSE

clean:
	rm -rf build