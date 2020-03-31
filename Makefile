.DEFAULT_GOAL := build

build:
	go build

pi:
	env GOOS=linux GOARCH=arm GOARM=5 go build

small:
	go build -ldflags="-s -w"
	upx --brute nas

small_pi:
	env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags="-s -w"
	upx --brute nas

run:
	go run webserver.go