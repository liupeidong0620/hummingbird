DIR = bin
NAME = hummingbird
MODULE = github.com/liupeidong0620/hummingbird

TAGS = ""
BUILD_FLAGS = "-v"

BUILD_COMMIT  = $(shell git rev-parse --short HEAD)
VERSION = $(shell git describe --tags || echo "$(BUILD_COMMIT)")
BUILD_TIME = $(shell date -u '+%FT%TZ')

LDFLAGS += -w -s -buildid=
LDFLAGS += -X "$(MODULE)/version.Version=$(VERSION)"
LDFLAGS += -X "$(MODULE)/version.BuildTime=$(BUILD_TIME)"
LDFLAGS += -X "$(MODULE)/version.SoftName=$(NAME)"
GO_BUILD = CGO_ENABLED=0 go build $(BUILD_FLAGS) -ldflags '$(LDFLAGS)' -tags '$(TAGS)' -trimpath

all: linux-amd64 darwin-amd64 windows-amd64

darwin-amd64:
	GOARCH=amd64 GOOS=darwin $(GO_BUILD) -o $(DIR)/$(NAME)-$@

linux-amd64:
	GOARCH=amd64 GOOS=linux $(GO_BUILD) -o $(DIR)/$(NAME)-$@

windows-amd64:
	GOARCH=amd64 GOOS=windows $(GO_BUILD) -o $(DIR)/$(NAME)-$@.exe

clean:
	rm $(DIR)/*
