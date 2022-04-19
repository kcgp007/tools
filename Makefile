GIT_BRANCH=$(shell git branch --show-current)
ifneq ($(findstring $(GIT_BRANCH),master),)
GIT_TAG=$(subst -,+,$(shell git describe --tags))
else ifneq ($(findstring $(GIT_BRANCH),develop),)
GIT_TAG=$(subst -,+,$(shell git describe --tags)).dev
else ifneq ($(findstring $(GIT_BRANCH),feature),)
GIT_TAG=$(subst -,+,$(shell git describe --tags)).dev
else ifneq ($(findstring $(GIT_BRANCH),bugfix),)
GIT_TAG=$(subst -,+,$(shell git describe --tags)).dev
else ifneq ($(findstring $(GIT_BRANCH),release),)
GIT_TAG=$(subst -,+,$(shell git describe --tags)).test
else ifneq ($(findstring $(GIT_BRANCH),hotfix),)
GIT_TAG=$(subst -,+,$(shell git describe --tags)).test
else
GIT_TAG=$(subst -,+,$(shell git describe --tags)).$(GIT_BRANCH)
endif
GIT_STATUS=$(shell git status --porcelain)
GO_VERSION=$(shell go version)
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-s -w -X 'tools/flagTool.version=$(GIT_TAG)' -X 'tools/flagTool.goVersion=$(GO_VERSION)' -X 'tools/flagTool.buildTime=$(BUILD_TIME)'

vet:
	go vet ./...

list:
	go list ./...

install:
	go mod tidy

update:
	go mod tidy
	go get -u all

cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

clean:
	go clean -x
	rm -f cover.out

version:
	@echo $(GIT_TAG)
	@echo $(GIT_BRANCH)
	@echo $(GIT_STATUS)
	@echo $(GO_VERSION)
	@echo $(BUILD_TIME)

build:
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"
	CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"

build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"
	CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"