build:
	go build -o cmd/bot/bot cmd/bot/*.go

build-ubuntu:
	GOOS=linux GOARCH=amd64 go build -o cmd/bot/bot-amd64-linux cmd/bot/*.go

run: build
	./cmd/bot/bot
