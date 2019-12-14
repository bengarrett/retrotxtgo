APP=retrotxtgo
COMMIT=$(shell git rev-list --abbrev-commit -1 HEAD)
CMD=github.com/bengarrett/retrotxtgo/cmd
DATE=$(shell date -u +%H:%M:%S/%Y-%m-%d)
HEADCNT=$(shell git rev-list --count HEAD)
FLAGS=-ldflags "-X ${CMD}.GoBuildGitCommit=${COMMIT} -X ${CMD}.GoBuildGitCount=${HEADCNT} -X ${CMD}.GoBuildDate=${DATE}"
BLD=go build ${FLAGS} -v -o ${APP}

.PHONY: build
build: clean
	go build -o ${APP} main.go

.PHONY: run
run:
	go run -race ${FLAGS} main.go version 

.PHONY: clean
clean:
	rm -rf ${APP}*

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${BLD}-linux

.PHONY: build-mac
build-mac:
	GOOS=darwin GOARCH=amd64 ${BLD}-mac

.PHONY: build-pi
build-pi:
	GOOS=linux GOARCH=arm GOARM=5 ${BLD}-arm5

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=amd64 ${BLD}-win