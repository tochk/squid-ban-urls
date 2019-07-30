build:
	qtc
	go fmt ./...
	GOOS=linux GOARCH=amd64 go build -o bin/squid_ban_urls