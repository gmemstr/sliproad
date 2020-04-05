.DEFAULT_GOAL := build

build:
	go build

pi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o nas-arm

small:
	go build -ldflags="-s -w"
	upx --brute nas -9 --no-progress

small_pi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -o nas-arm -ldflags="-s -w"
	upx --brute nas-arm -9 --no-progress

run:
	go run webserver.go

test:
	go test ./... -cover

dist: clean small small_pi
	mkdir build
	mv nas* build

clean:
	rm -rf build