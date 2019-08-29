.PHONY: all build

GO111MODULE?=on
GOPROXY?=https://proxy.golang.org

all: build

build:
	go build -o repair main.go
