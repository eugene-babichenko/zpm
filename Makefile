TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(TAG:v%=%)
ifneq ($(COMMIT), $(TAG_COMMIT))
	VERSION := $(VERSION)-next-$(COMMIT)-$(DATE)
endif
ifeq ($(VERSION),)
	VERSION := $(COMMIT)-$(DATA)
endif
ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

FLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(FLAGS) -o zpm-$(VERSION) main.go

run:
	go run $(FLAGS) main.go

install:
	go install $(FLAGS)

test:
	GO111MODULE=on go test -race -covermode=atomic -coverpkg=./... ./...

