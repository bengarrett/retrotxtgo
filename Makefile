APP=retrotxtgo
VERSION=$(shell cat ./VERSION)

COMMIT=$(shell git rev-list --abbrev-commit -1 HEAD)
CMD=github.com/bengarrett/retrotxtgo/lib/cmd
DATE=$(shell date -u +%H:%M:%S/%Y-%m-%d)
HEADCNT=$(shell git rev-list --count HEAD)
FLAGS=-ldflags "-X ${CMD}.GoBuildVer=${VERSION} -X ${CMD}.GoBuildGitCommit=${COMMIT} -X ${CMD}.GoBuildGitCount=${HEADCNT} -X ${CMD}.GoBuildDate=${DATE}"
BLD=go build ${FLAGS} -v -o artifacts/${APP}-v${VER}.${HEADCNT}

.PHONY: run
run:
	go run -race ${FLAGS} main.go version 

.PHONY: ver
ver: 
	${BLD}
	artifacts/${APP}-v${VER}.${HEADCNT} version

.PHONY: clean
clean:
	rm -rf artifacts/${APP}*

.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${BLD}-linux

.PHONY: build-mac
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ${BLD}-mac

.PHONY: build-pi
build-pi:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=5 ${BLD}-arm5

.PHONY: build-windows
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ${BLD}-win

.PHONY: test
test:
	go test ./...

.PHONY: lint
lint:
	golangci-lint run

# make release -j4 for 4 multiple threads
.PHONY: release
release: clean build-mac build-linux build-pi build-windows

.PHONY: mod
mod:
	go mod tidy