GIT_BRANCH=$(shell git branch --show-current)
GIT_TAG=$(shell git describe --tags)
ifneq ($(findstring $(GIT_BRANCH),master),)
GIT_VERSION=$(GIT_TAG)
else ifneq ($(findstring develop,$(GIT_BRANCH)),)
GIT_VERSION=$(GIT_TAG).dev
else ifneq ($(findstring feature,$(GIT_BRANCH)),)
GIT_VERSION=$(GIT_TAG).dev
else ifneq ($(findstring bugfix,$(GIT_BRANCH)),)
GIT_VERSION=$(GIT_TAG).dev
else ifneq ($(findstring release,$(GIT_BRANCH)),)
GIT_VERSION=$(GIT_TAG).$(subst release/,test.,$(GIT_BRANCH))
else ifneq ($(findstring hotfix,$(GIT_BRANCH)),)
GIT_VERSION=$(GIT_TAG).fix
else
GIT_VERSION=$(GIT_TAG).$(GIT_BRANCH)
endif
GIT_STATUS=$(shell git status --porcelain)
GO_VERSION=$(shell go version)
BUILD_TIME=$(shell date +%FT%T%z)
LDFLAGS=-s -w -X 'tools/flagTool.version=$(GIT_VERSION)' -X 'tools/flagTool.goVersion=$(GO_VERSION)' -X 'tools/flagTool.buildTime=$(BUILD_TIME)'

#build:
#ifneq ($(GIT_STATUS),)
#	@git status
#else
#	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"
#	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"
#endif

#build:
#ifneq ($(GIT_STATUS),)
#	@git status
#else
#	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"
#	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/ -trimpath -ldflags "$(LDFLAGS)"
#endif

vet:
	go vet ./...

cover:
	go test -coverprofile cover.out ./...
	go tool cover -html=cover.out

clean:
	go clean -x
	rm -f cover.out
	rm -rf bin

test:
	@echo $(GIT_TAG)
	@echo $(GIT_VERSION)
	@echo $(GIT_BRANCH)
	@echo $(GIT_STATUS)
	@echo $(GO_VERSION)
	@echo $(BUILD_TIME)