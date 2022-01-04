VERSION=v1.2.3
GO_VERSION=$(shell go version)
GIT_COMMIT=$(shell git rev-parse HEAD)
BUILD_TIME=$(shell date +"%Y-%m-%d %H:%M:%S")
LDFLAGS=-s -w -X 'tools/flagTool.version=${VERSION}' -X 'tools/flagTool.goVersion=${GO_VERSION}' -X 'tools/flagTool.gitCommit=${GIT_COMMIT}' -X 'tools/flagTool.buildTime=${BUILD_TIME}'

vet:
	go vet ./...

cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

clean:
	go clean -x
	rm -f cover.out

build:
	go build -trimpath -ldflags "${LDFLAGS}"

build-gcc-linux:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc go build -trimpath -ldflags "${LDFLAGS}"

build-gcc-windows:
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -trimpath -ldflags "${LDFLAGS}"

build-gcc-arm:
	CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -trimpath -ldflags "${LDFLAGS}"

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags "${LDFLAGS}"

build-windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags "${LDFLAGS}"

build-arm:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -trimpath -ldflags "${LDFLAGS}"