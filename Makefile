.DEFAULT_GOAL := build

build:
	go build

pi:
	env GOOS=linux GOARCH=arm GOARM=5 go build

run:
	go run webserver.go