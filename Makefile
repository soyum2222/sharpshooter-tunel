build:
	go mod tidy
	go build -o ./sharpshoot-client ./client/main.go
	go build -o ./sharpshoot-server ./server/main.go