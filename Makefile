linux:
	go mod tidy
	GOOS=linux GOARCH=amd64 go build -gcflags 'all=-trimpath=$(GOPATH)' -asmflags 'all=-trimpath=$(GOPATH)' -o ./sharpshooter-client-linux-amd64 ./client/main.go
	GOOS=linux GOARCH=amd64 go build -gcflags 'all=-trimpath=$(GOPATH)' -asmflags 'all=-trimpath=$(GOPATH)' -o ./sharpshooter-server-linux-amd64 ./server/main.go
	GOOS=linux GOARCH=arm   go build -gcflags 'all=-trimpath=$(GOPATH)' -asmflags 'all=-trimpath=$(GOPATH)' -o ./sharpshooter-client-linux-arm ./client/main.go
	GOOS=linux GOARCH=arm   go build -gcflags 'all=-trimpath=$(GOPATH)' -asmflags 'all=-trimpath=$(GOPATH)' -o ./sharpshooter-server-linux-arm ./server/main.go

win:
	go mod tidy
	GOOS=windows GOARCH=amd64 go build -gcflags 'all=-trimpath=$(GOPATH)' -asmflags 'all=-trimpath=$(GOPATH)' -o ./sharpshooter-client-win-amd64.exe ./client/main.go
	GOOS=windows GOARCH=amd64 go build -gcflags 'all=-trimpath=$(GOPATH)' -asmflags 'all=-trimpath=$(GOPATH)' -o ./sharpshooter-server-win-amd64.exe ./server/main.go

mac:
	go mod tidy
	GOOS=darwin GOARCH=amd64 go build -gcflags 'all=-trimpath=$(GOPATH)' -asmflags 'all=-trimpath=$(GOPATH)' -o ./sharpshooter-client-darwin-amd64 ./client/main.go
