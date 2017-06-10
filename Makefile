VERSION = v0.0.1
GIT_COMMIT = $(shell git rev-parse --short HEAD)
LDFLAGS = -ldflags "-X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION) -s -w"

all: build
docker: docker-build


build:
	go build $(LDFLAGS) -o redgo main.go

docker-build:
	docker run --rm -v `pwd`:/go/src/github.com/wupeaking/redgo golang:1.8-alpine go build $(LDFLAGS) -o /go/src/github.com/wupeaking/redgo/redgo /go/src/github.com/wupeaking/redgo/main.go
	docker build -t redgo:$(VERSION) .

