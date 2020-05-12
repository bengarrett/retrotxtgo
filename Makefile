APP=retrotxtgo
VERSION=$(shell cat ./VERSION)

COMMIT=$(shell git rev-list --abbrev-commit -1 HEAD)
CMD=github.com/bengarrett/retrotxtgo/lib/cmd
DATE=$(shell date -u +%H:%M:%S/%Y-%m-%d)
HEADCNT=$(shell git rev-list --count HEAD)
FLAGS=-ldflags "-X ${CMD}.GoBuildVer=${VERSION} -X ${CMD}.GoBuildGitCommit=${COMMIT} -X ${CMD}.GoBuildGitCount=${HEADCNT} -X ${CMD}.GoBuildDate=${DATE}"
BLD=go build ${FLAGS} -v -o artifacts/${APP}-v${VER}.${HEADCNT}

.PHONY: run
## run: runs the app with race conditions
run:
	go run -race ${FLAGS} main.go version 

.PHONY: ver
## ver: builds the app and display the version information
ver: 
	${BLD}
	artifacts/${APP}-v${VER}.${HEADCNT} version

.PHONY: clean
## clean: removes any existing app builds
clean:
	rm -rf artifacts/${APP}*

.PHONY: build-linux
## build-linux: build the app for Linux 64-bit
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${BLD}-linux

.PHONY: build-mac
## build-mac: build the app for macOS 64-bit
build-mac:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 ${BLD}-mac

.PHONY: build-pi
## build-pi: build the app for Linux 32-bit ARM 8 commonly in use on Raspberry Pis
build-pi:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=8 ${BLD}-arm8

.PHONY: build-windows
## build-windows: build the app for macOS 64-bit
build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 ${BLD}-win

.PHONY: test
## test: run all app tests
test:
	go test ./...

.PHONY: lint
## lint: performance and best practice analysts
lint:
	golangci-lint run

.PHONY: release
## release: build the app for all listed platforms
# make release -j4 for 4 multiple threads
release: clean build-mac build-linux build-pi build-windows

# release:
# 	git tag -a $(VERSION) -m "Release" || true
# 	git push origin $(VERSION)
# 	goreleaser --rm-dist

# image:
# 	docker build -t cirocosta/l7 .

.PHONY: mod
## mod: prune any no-longer-needed dependencies
mod:
	go mod tidy

.PHONY: help
## help: display this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

# sources:
# https://danishpraka.sh/2019/12/07/using-makefiles-for-go.html