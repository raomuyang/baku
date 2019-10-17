.PHONY: test compile

test:
	go test -race ./...
	go vet ./...


compile:
	mkdir -p target
	GOOS=windows GOARCH=amd64 go build  -o target/`go run main.go -v`-win.exe
	GOOS=linux GOARCH=amd64 go build  -o target/`go run main.go -v`-linux
	GOOS=darwin GOARCH=amd64 go build  -o target/`go run main.go -v`-mac
	chmod a+x target/*
