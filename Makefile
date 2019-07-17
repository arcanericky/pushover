VERSION_INJECT=main.versionText
SRCS=*.go cmd/*.go
PACKAGE=cmd/*.go
export GO111MODULE=on

EXECUTABLE=bin/pushover

LINUX=$(EXECUTABLE)-linux
DARWIN=$(EXECUTABLE)-darwin
WINDOWS=$(EXECUTABLE)-windows

LINUX_AMD64=$(LINUX)-amd64
DARWIN_AMD64=$(DARWIN)-amd64
WINDOWS_AMD64=$(WINDOWS)-amd64.exe

LINUX_386=$(LINUX)-386
WINDOWS_386=$(WINDOWS)-386.exe

.PHONY: all test clean

all: test build

build: linux-amd64 windows-amd64 darwin-amd64 linux-386 windows-386

linux-amd64: $(LINUX_AMD64)

windows-amd64: $(WINDOWS_AMD64)

darwin-amd64: $(DARWIN_AMD64)

linux-386: $(LINUX_386)

windows-386: $(WINDOWS_386)

test:
	go test -race -coverprofile=coverage.txt -covermode=atomic ./...
	go tool cover -html=coverage.txt -o coverage.html

# AMD64 Versions
$(WINDOWS_AMD64): $(SRCS)
	GOOS=windows GOARCH=amd64 go build -o $@ -ldflags "-s -w -X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" $(PACKAGE)

$(LINUX_AMD64): $(SRCS)
	GOOS=linux GOARCH=amd64 go build -o $@ -ldflags "-s -w -X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" $(PACKAGE)

$(DARWIN_AMD64): $(SRCS)
	GOOS=darwin GOARCH=amd64 go build -o $@ -ldflags "-s -w -X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" $(PACKAGE)

# 386 Versions
$(WINDOWS_386): $(SRCS)
	GOOS=windows GOARCH=386 go build -o $@ -ldflags "-s -w -X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" $(PACKAGE)

$(LINUX_386): $(SRCS)
	GOOS=linux GOARCH=386 go build -o $@ -ldflags "-s -w -X $(VERSION_INJECT)=$(shell sh scripts/get-version.sh)" $(PACKAGE)

clean:
	rm -rf bin
