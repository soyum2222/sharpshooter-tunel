build:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -o ./sharpshooter-client-linux-amd64 ./client/main.go
	GOOS=linux GOARCH=amd64 go build -o ./sharpshooter-server-linux-amd64 ./server/main.go
	GOOS=linux GOARCH=arm   go build -o ./sharpshooter-client-linux-arm ./client/main.go
	GOOS=linux GOARCH=arm   go build -o ./sharpshooter-server-linux-arm ./server/main.go
