.PHONY: all dep build

all: build

dep:
	dep ensure -v

build: dep
	go build -o repair main.go
